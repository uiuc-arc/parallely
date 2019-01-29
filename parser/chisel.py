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

    def visitVariable(self, ctx):
        return [ctx.getText()]

    def visitAnd(self, ctx):
        e_1 = self.visit(ctx.expression(0))
        e_2 = self.visit(ctx.expression(1))
        return e_1 + e_2

    def visitAdd(self, ctx):
        e_1 = self.visit(ctx.expression(0))
        e_2 = self.visit(ctx.expression(1))
        return e_1 + e_2

    def processProbassignment(self, ctx, spec):
        print "in processprobassn", spec
        p = ctx.probability()
        e_1 = self.visit(ctx.expression(0))
        e_2 = self.visit(ctx.expression(1))
        assigned_var = ctx.var().getText()
        new_spec = list(spec) #copy it
        if spec[1] == assigned_var:
            new_spec[1] = ''
            new_spec[2] = 255 #assume this is max of interval
        print new_spec, assigned_var, p.getText(), e_1, e_2
        return new_spec
        
    def processCast(self, ctx, spec):
        assigned_var = ctx.var(0).getText()
        assignee_var = ctx.var(1).getText()
        cast_type = ctx.fulltype().getText()
        new_spec = list(spec) #copy it
        if spec[1] == assigned_var:
            new_spec[1] = assignee_var
            if cast_type == 'approxfloat32':
                new_spec[2] += 2**(-18) #assume this is conversion error
        print new_spec, assigned_var, assignee_var, cast_type
        return new_spec

    def processExpassignment(self, ctx, spec):
        assigned_var = ctx.var().getText()
        vars_list = self.visit(ctx.expression())
        new_spec = list(spec) #copy it
        if spec[1] == assigned_var:
            #assert(len(vars_list)==1)
            new_spec[1] = vars_list[0]
            new_spec[3] *= len(vars_list)
        print new_spec, assigned_var, vars_list
        return new_spec
    
    def processspec(self, statements, spec):
        for statement in statements:
            self.visit(statement)
            if isinstance(statement, ParallelyParser.CastContext):
                spec = self.processCast(statement, spec)
            elif isinstance(statement, ParallelyParser.ExpassignmentContext):
                spec = self.processExpassignment(statement, spec)
            elif isinstance(statement, ParallelyParser.ProbassignmentContext):
                spec = self.processProbassignment(statement, spec)
            else:
                continue
        return spec
        
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

    # String manipulation for now.
    # Parser later
    def generateRelyCondition(self, ctx, spec):
        e_b, e_e = spec.split('>=')
        e_b = e_b.strip()
        e_e = e_e.strip()[2:-1]
        print e_b, e_e
        self.spec = [e_b,e_e,0.0,1.0]

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
