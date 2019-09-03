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
from unroller import unrollRepeat
from antlr4.error.ErrorListener import ErrorListener
from collections import namedtuple

key_error_msg = "Type error detected: Undeclared variable (probably : {})"

Constraint = namedtuple('Constraint', "limit condition multiplicative jointreliability")


# class MyErrorListener( ErrorListener ):
#     def __init__(self):
#         super(MyErrorListener, self).__init__()

#     def syntaxError(self, recognizer, offendingSymbol, line, column, msg, e):
#         raise Exception("Parsing Syntax error!! : ", e)

#     def reportAmbiguity(self, recognizer, dfa, startIndex, stopIndex, exact, ambigAlts, configs):
#         raise Exception("Ambiguious Syntax error!! : ")

#     def reportAttemptingFullContext(self, recognizer, dfa, startIndex, stopIndex, conflictingAlts, configs):
#         raise Exception("AttemptingFullContext", conflictingAlts)

#     def reportContextSensitivity(self, recognizer, dfa, startIndex, stopIndex, prediction, configs):
#         raise Exception("Oh no!!")


class CalculatePSuccess(ParallelyVisitor):
    def visitProbassignment(self, ctx):
        try:
            p = float(ctx.probability().getText())
        except ValueError:
            print "The probabilities have to be numbers: ", ctx.getText()
            exit(-1)
        return p

    def visitIf(self, ctx):
        prob_ifs = 1
        prob_elses = 1
        for statements in ctx.ins:
            temp = self.vist(statement)
            if temp:
                prob_ifs = prob_ifs * temp
        for statements in ctx.elses:
            temp = self.vist(statement)
            if temp:
                prob_ifs = prob_ifs * temp
        return min(prob_ifs, prob_elses)

    def calc(self, statements):
        pass_prob = 1
        for statement in statements:
            prob_temp = self.visit(statement)
            if prob_temp:
                pass_prob = pass_prob * prob_temp
        return pass_prob


class relyGenerator(ParallelyVisitor):
    def __init__(self):
        self.typecontext = {}
        self.processgroups = {}

    def visitSequential(self, ctx):
        type1 = self.visit(ctx.getChild(0))
        type2 = self.visit(ctx.getChild(2))
        return (type1 and type2)

    def visitSkipstatement(self, ctx):
        return True

    def visitLiteral(self, ctx):
        return []

    def visitFliteral(self, ctx):
        return []

    def visitVariable(self, ctx):
        return [ctx.getText()]

    def visitMultiply(self, ctx):
        e_1 = self.visit(ctx.expression(0))
        e_2 = self.visit(ctx.expression(1))
        return e_1 + e_2

    def visitSelect(self, ctx):
        return self.visit(ctx.expression())

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
        new_spec = []
        print "Updating spec", assigned_var, vars_list, multiplicatives, spec
        for spec_part in spec:
            temp_joinedrel = set(spec_part.jointreliability)
            temp_mul = spec_part.multiplicative * multiplicatives
            if assigned_var in spec_part.jointreliability:
                temp_joinedrel.remove(assigned_var)
                temp_joinedrel.update(vars_list)
            new_spec.append(Constraint(spec_part.limit, spec_part.condition, temp_mul, temp_joinedrel))
        print "--->", new_spec
        return new_spec

    def processExpassignment(self, ctx, spec):
        assigned_var = ctx.var().getText()
        vars_list = self.visit(ctx.expression())
        new_spec = self.updateSpec(spec, ctx, assigned_var, vars_list, 1)
        # print new_spec, assigned_var, vars_list
        return new_spec

    # Assumes that functions are implemented precisely
    # Only errors in the date propogate
    def processfunction(self, ctx, spec):
        assigned_var = ctx.var().getText()
        vars_list = []

        expression_list = ctx.expression()
        for expr in expression_list:
            temp_vars_list = self.visit(expr)
            if not temp_vars_list:
                continue
            else:
                vars_list.extend(temp_vars_list)
        new_spec = self.updateSpec(spec, ctx, assigned_var, vars_list, 1)
        # print new_spec, assigned_var, vars_list
        return new_spec

    def processAExpassignment(self, ctx, spec):
        assigned_var = ctx.var().getText()
        vars_list = []
        expression_list = ctx.expression()
        for expr in expression_list:
            temp_vars_list = self.visit(expr)
            if not temp_vars_list:
                continue
            else:
                vars_list.extend(temp_vars_list)
        new_spec = self.updateSpec(spec, ctx, assigned_var, vars_list, 1)
        return new_spec

    def processALoad(self, ctx, spec):
        assigned_var = ctx.var(0).getText()
        array_var = ctx.var(1).getText()
        vars_list = [array_var]
        expression_list = ctx.expression()
        for expr in expression_list:
            temp_vars_list = self.visit(expr)
            if not temp_vars_list:
                continue
            else:
                vars_list.extend(temp_vars_list)

        new_spec = self.updateSpec(spec, ctx, assigned_var, vars_list, 1)
        # print new_spec, assigned_var, vars_list
        return new_spec

    def processCast(self, ctx, spec):
        assigned_var = ctx.var(0).getText()
        # print assigned_var, spec
        new_spec = self.updateSpec(spec, ctx, assigned_var, [0], 1)
        # print "[Debug] Cast : ", new_spec, assigned_var
        return new_spec

    def processDec(self, ctx, spec):
        assigned_var = ctx.var().getText()
        print assigned_var
        new_spec = self.updateSpec(spec, ctx, assigned_var, [], 1)
        return new_spec

    def processProbassignment(self, ctx, spec):
        try:
            p = float(ctx.probability().getText())
        except ValueError:
            print "Probabilities have to be numbers"
            exit(-1)
        e_1 = self.visit(ctx.expression(0))
        # e_2 = self.visit(ctx.expression(1))
        assigned_var = ctx.var().getText()

        # if isinstance(p, ParallelyParser.VarprobContext):
        #     new_items = e_1 + e_2 + [assigned_var, p.getText()]
        #     return self.updateSpec(spec, ctx, assigned_var, new_items, [])
        # if isinstance(p, ParallelyParser.FloatprobContext):
        new_items = e_1 + [assigned_var]
        return self.updateSpec(spec, ctx, assigned_var,
                               new_items, p)

    def processRecover(self, ctx, spec):
        ps1calculator = CalculatePSuccess()
        ps1 = ps1calculator.calc(ctx.trys)
        spec_try = self.processspec(ctx.trys, spec)
        spec_recover = self.processspec(ctx.recovers, spec)
        ps2 = spec_recover[0].multiplicative / spec[0].multiplicative

        newspec = []
        for i, spec_part in enumerate(spec):
            s1_data = spec_try[i].jointreliability
            s2_data = spec_recover[i].jointreliability
            all_data = s1_data | s2_data  # Set union. Magic !!!!
            new_mult = ps1 * spec_part.multiplicative + (1-ps1) * spec_recover[i].multiplicative
            newConstraint = Constraint(spec_part.limit,
                                       spec_part.condition,
                                       new_mult,
                                       all_data)
            newspec.append(newConstraint)
            # print "==============================>", ps1, ps2, all_data
        return newspec

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
        reversed_statements = statements[::-1]
        for i, statement in enumerate(reversed_statements):
            print "Processing : {} ({}/{})".format(statement.getText(), i, len(reversed_statements))
            # self.visit(statement)
            if isinstance(statement, ParallelyParser.CastContext):
                spec = self.processCast(statement, spec)
            elif isinstance(statement, ParallelyParser.ProbassignmentContext):
                spec = self.processProbassignment(statement, spec)
            elif isinstance(statement, ParallelyParser.ExpassignmentContext):
                spec = self.processExpassignment(statement, spec)
            elif isinstance(statement, ParallelyParser.ArrayassignmentContext):
                spec = self.processAExpassignment(statement, spec)
            elif isinstance(statement, ParallelyParser.ArrayloadContext):
                spec = self.processALoad(statement, spec)
            elif isinstance(statement, ParallelyParser.FuncContext):
                spec = self.processfunction(statement, spec)
            elif isinstance(statement, ParallelyParser.RecoverContext):
                spec = self.processRecover(statement, spec)
            elif isinstance(statement, ParallelyParser.SingledeclarationContext):
                spec = self.processDec(statement, spec)
            elif isinstance(statement, ParallelyParser.ArraydeclarationContext):
                spec = self.processDec(statement, spec)
            # elif isinstance(statement, ParallelyParser.IfContext):
            #     if_branch = statement.ifs
            #     else_branch = statement.elses
            #     b_cond = statement.var().getText()
            #     if_spec = self.processspec(if_branch, spec)
            #     else_spec = self.processspec(else_branch, spec)
            #     # print if_spec, else_spec
            #     for spec_part in if_spec:
            #         spec_part[2].add(b_cond)
            #     for spec_part in else_spec:
            #         spec_part[2].add(b_cond)
            #     spec = if_spec + else_spec
            else:
                print "Unable to process the statement :", statement.getText()
                exit(-1)
        return spec

    # String manipulation for now.
    # Parser later
    # def generateRelyCondition(self, ctx, spec):
    #     disjoints = spec.split('^')
    #     for pred in disjoints:
    #         r_1, r_2 = pred.split('<=')
    #         rs = r_2[3:-2].split(',')
    #         print r_2, rs
    #         rs_cleaned = [r.strip() for r in rs]
    #         self.spec.append((r_1, set(), set(rs_cleaned)))

    #     statements = ctx.statement()[::-1]
    #     spec = self.processspec(statements, self.spec)
    #     return spec


# Takes in a .seq file performs the rely reliability analysis
def main(program_str, spec, skiprename):
    input_stream = InputStream(program_str)
    # lexer = ParallelyLexer(input_stream)
    # stream = CommonTokenStream(lexer)
    # parser = ParallelyParser(stream)

    start = time.time()
    print "Unrolling Repeat statements?: ", (not skiprename)
    if not skiprename:
        i = 0
        while(True):
            print "unrolling {} deep".format(i)

            lexer = ParallelyLexer(input_stream)
            stream = CommonTokenStream(lexer)
            parser = ParallelyParser(stream)
            # parser.addErrorListener(MyErrorListener())
            try:
                tree = parser.parallelprogram()
            except Exception as e:
                print "Parsing Error!!!"
                print e
                exit(-1)
            unroller = unrollRepeat(stream)
            walker = ParseTreeWalker()
            walker.walk(unroller, tree)
            input_stream = InputStream(unroller.rewriter.getDefaultText())
            if not unroller.replacedone:
                input_stream = InputStream(unroller.rewriter.getDefaultText())
                break
                # if debug:
            i = i + 1
            # print "----------------------------------------"
            # print "Intermediate step. Writing to _DEBUG_UNROLLED_.txt"
            debug_file = open("_DEBUG_UNROLLED_{}.txt".format(i), 'w')
            debug_file.write(input_stream.strdata)
            debug_file.close()
            # print "----------------------------------------"

        print "----------------------------------------"
        print "Intermediate step. Writing to _DEBUG_UNROLLED_.txt"
        debug_file = open("_DEBUG_UNROLLED_.txt", 'w')
        debug_file.write(input_stream.strdata)
        debug_file.close()
        print "----------------------------------------"

    start2 = time.time()

    # print input_stream.strdata
    lexer = ParallelyLexer(input_stream)
    stream = CommonTokenStream(lexer)
    parser = ParallelyParser(stream)
    tree = parser.parallelprogram()
    start3 = time.time()

    # if len(tree.program()) > 1:
    #     print "Needs to be a sequential program for rely analysis"
    #     exit(-1)

    # spec = rely.generateRelyCondition(tree, spec.read())

    # Processing the spec
    spec_input_stream = InputStream(spec)
    spec_lexer = ParallelyLexer(spec_input_stream)
    spec_stream = CommonTokenStream(spec_lexer)
    spec_parser = ParallelyParser(spec_stream)

    spec_str = spec_parser.relyspec()
    print spec_str.getText()

    rely_spec = []
    for constraint_str in spec_str.singlerelyspec():
        temp_limit = float(constraint_str.FLOAT(0).getText())
        if constraint_str.FLOAT(1):
            temp_mult = float(constraint_str.FLOAT(1))
        else:
            temp_mult = 1
        var_list = []
        for var in constraint_str.VAR():
            var_list.append(var.getText())
        rely_spec.append(Constraint(temp_limit, "<=", temp_mult, var_list))

    print rely_spec

    rely = relyGenerator()
    result_spec = rely.processspec(tree.program(0).statement(), rely_spec)

    # if the variable declaration is found the reliability is 1
    decs = tree.program(0).declaration()
    result_spec = rely.processspec(decs, result_spec)
    end = time.time()

    print '----------------------------------------'
    print result_spec
    print '----------------------------------------'
    print "Analysis time :", end - start, end - start2, end - start3


if __name__ == '__main__':
    sys.setrecursionlimit(15000)
    programfile = open(sys.argv[1], 'r')
    spec = open(sys.argv[2], 'r').read()
    unroll = (sys.argv[3] == '1')
    program_str = programfile.read()
    main(program_str, spec, unroll)
