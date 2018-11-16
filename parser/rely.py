import sys
from antlr4 import CommonTokenStream
from antlr4 import InputStream
import TokenStreamRewriter
from antlr4 import ParseTreeWalker
from antlr4 import *
from ParallelyLexer import ParallelyLexer
from ParallelyParser import ParallelyParser
from antlr4.tree.Trees import Trees
from ParallelyVisitor import ParallelyVisitor
from ParallelyListener import ParallelyListener
import copy
import time

key_error_msg = "Type error detected: Undeclared variable (probably : {})"


class unrollLoops(ParallelyListener):
    def __init__(self, stream):
        self.globalvariables = {}
        self.rewriter = TokenStreamRewriter.TokenStreamRewriter(stream)
        self.done = set()

    def enterSingleglobaldec(self, ctx):
        if not ctx.GLOBALVAR():
            return
        var_name = ctx.GLOBALVAR().getText()
        var_values = [v.getText() for v in ctx.processid()]
        self.globalvariables[var_name] = var_values

    def enterForloop(self, ctx):
        var_group = ctx.GLOBALVAR().getText()
        concrete_vars = self.globalvariables[var_group]
        # statements = ctx.statement().getText()
        orig_variable = ctx.VAR().getText()
        edited = ''

        cs = ctx.statement().start.getInputStream()
        statements = cs.getText(ctx.statement().start.start,
                                ctx.statement().stop.stop)
        print '-------------------------------'
        print statements
        print '-------------------------------'

        # removing the code for process groups
        self.rewriter.delete(self.rewriter.DEFAULT_PROGRAM_NAME,
                             ctx.start.tokenIndex, ctx.stop.tokenIndex)
        for var in concrete_vars:
            # Including the _s to be safe. Still can screw up a lot
            # Deadline mode
            new_version = statements.replace("_" + orig_variable, "_" + var)
            edited += new_version + ";\n"
        self.rewriter.insertAfter(ctx.stop.tokenIndex, edited)


class relyGenerator(ParallelyVisitor):
    def __init__(self):
        self.typecontext = {}
        self.processgroups = {}
        self.spec = []

    def visitSequential(self, ctx):
        type1 = self.visit(ctx.getChild(0))
        type2 = self.visit(ctx.getChild(2))
        return (type1 and type2)

    def visitSkipstatement(self, ctx):
        return True

    def visitLiteral(self, ctx):
        return []

    def visitVariable(self, ctx):
        return [ctx.getText()]

    def visitMultiply(self, ctx):
        e_1 = self.visit(ctx.expression(0))
        e_2 = self.visit(ctx.expression(1))
        return e_1 + e_2

    def visitDivide(self, ctx):
        e_1 = self.visit(ctx.expression(0))
        e_2 = self.visit(ctx.expression(1))
        return e_1 + e_2

    def visitAdd(self, ctx):
        e_1 = self.visit(ctx.expression(0))
        e_2 = self.visit(ctx.expression(1))
        return e_1 + e_2

    def visitMinus(self, ctx):
        e_1 = self.visit(ctx.expression(0))
        e_2 = self.visit(ctx.expression(1))
        return e_1 + e_2

    def visitProb(self, ctx):
        e_1 = self.visit(ctx.expression(0))
        e_2 = self.visit(ctx.expression(1))
        return e_1 + e_2

    def updateSpec(self, spec, ctx, assigned_var, vars_list, multiplicatives):
        # print "Updating spec", assigned_var, vars_list, multiplicatives, spec
        for spec_part in self.spec:
            if assigned_var in spec_part[2]:
                spec_part[2].remove(assigned_var)
                spec_part[2].update(vars_list)
            if len(multiplicatives) > 0:
                spec_part[1].extend(multiplicatives)
        return spec

    def processExpassignment(self, ctx, spec):
        assigned_var = ctx.var().getText()
        vars_list = self.visit(ctx.expression())
        new_spec = self.updateSpec(spec, ctx, assigned_var, vars_list, [])
        print new_spec, assigned_var, vars_list
        return new_spec

    def processProbassignment(self, ctx, spec):
        p = ctx.probability()
        e_1 = self.visit(ctx.expression(0))
        e_2 = self.visit(ctx.expression(1))
        assigned_var = ctx.var().getText()

        if isinstance(p, ParallelyParser.VarprobContext):
            new_items = e_1 + e_2 + [assigned_var, p.getText()]
            return self.updateSpec(spec, ctx, assigned_var, new_items, [])
        if isinstance(p, ParallelyParser.FloatprobContext):
            new_items = e_1 + e_2 + [assigned_var]
            return self.updateSpec(spec, ctx, assigned_var,
                                   new_items, [float(p.getText())])

    def flattenStatement(self, ctx):
        if isinstance(ctx, ParallelyParser.MultipledeclarationContext):
            first_half = self.flattenStatement(ctx.getChild(0))
            second_half = self.flattenStatement(ctx.getChild(2))
            return first_half + second_half
        if isinstance(ctx, ParallelyParser.SeqcompositionContext):
            first_half = self.flattenStatement(ctx.getChild(0))
            second_half = self.flattenStatement(ctx.getChild(2))
            return first_half + second_half
        else:
            return [ctx]

    def processspec(self, statements, spec):
        for statement in statements:
            # print "Processing : ", statement.getText(), spec
            self.visit(statement)
            if isinstance(statement, ParallelyParser.ProbassignmentContext):
                spec = self.processProbassignment(statement, spec)
            elif isinstance(statement, ParallelyParser.ExpassignmentContext):
                spec = self.processExpassignment(statement, spec)
            elif isinstance(statement, ParallelyParser.IfContext):
                if_branch = self.flattenStatement(statement.statement(0))
                else_branch = self.flattenStatement(statement.statement(1))
                b_cond = statement.var().getText()
                if_spec = self.processspec(if_branch, spec)
                else_spec = self.processspec(else_branch, spec)
                # print if_spec, else_spec
                for spec_part in if_spec:
                    spec_part[2].add(b_cond)
                for spec_part in else_spec:
                    spec_part[2].add(b_cond)
                spec = if_spec + else_spec
            else:
                continue
        return spec

    # String manipulation for now.
    # Parser later
    def generateRelyCondition(self, ctx, spec):
        disjoints = spec.split('^')
        for pred in disjoints:
            r_1, r_2 = pred.split('<=')
            rs = r_2[3:-2].split(',')
            print r_2, rs
            rs_cleaned = [r.strip() for r in rs]
            self.spec.append([r_1, [], set(rs_cleaned)])

        statements = self.flattenStatement(ctx.statement())[::-1]
        spec = self.processspec(statements, self.spec)
        print '----------------------------------------'
        print spec
        print '----------------------------------------'


# Takes in a .seq file performs the rely reliability analysis
def main(program_str, spec):
    input_stream = InputStream(program_str)
    lexer = ParallelyLexer(input_stream)
    stream = CommonTokenStream(lexer)
    parser = ParallelyParser(stream)

    # # Unroll process groups for easy analysis?
    # For now not doing this
    # Damages the readability of the code
    start = time.time()
    tree = parser.sequentialprogram()
    renamer = unrollLoops(stream)
    walker = ParseTreeWalker()
    walker.walk(renamer, tree)

    print "----------------------------------------"
    print "Intermediate step"
    print renamer.rewriter.getDefaultText()
    print "----------------------------------------"

    # print renamer.rewriter.getDefaultText()

    start2 = time.time()
    input_stream = InputStream(renamer.rewriter.getDefaultText())

    lexer = ParallelyLexer(input_stream)
    stream = CommonTokenStream(lexer)
    parser = ParallelyParser(stream)
    tree = parser.sequentialprogram()

    start3 = time.time()
    rely = relyGenerator()
    rely.generateRelyCondition(tree, spec.read())
    end = time.time()

    print "Analysis time :", end - start, end - start2, end - start3


if __name__ == '__main__':
    sys.setrecursionlimit(15000)
    programfile = open(sys.argv[1], 'r')
    spec = open(sys.argv[2], 'r')
    program_str = programfile.read()
    main(program_str, spec)
