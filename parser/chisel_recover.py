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
from argparse import ArgumentParser
import time
from unroller import unrollRepeat
from antlr4.error.ErrorListener import ErrorListener
# from antlr4.PredictionContext import PredictionMode
import json

# from antlr4.
from collections import namedtuple

key_error_msg = "Type error detected: Undeclared variable (probably : {})"

Constraint = namedtuple('Constraint', "limit condition multiplicative jointreliability")

# nature of data: dict mapping variable to coefficient (for constant use 1 as variable)
# add two affine expressions, adding up terms of common variables
def addAff(aff1, aff2):
    keys = set(aff1.keys()).union(set(aff2.keys()))
    result = {}
    for key in keys:
        coeff = aff1.get(key,0.0)+aff2.get(key,0.0)
        if coeff != 0.0:
            result[key] = coeff
    return result
# multiply affine expression by constant
def multAff(aff, n):
    if n == 0.0:
        return {}
    else:
        return {key:val*n for key,val in aff.items()}
# replace variable with affine expression
def replaceAff(aff1, var, aff2):
    if var in aff1:
        coeff = aff1[var]
        aff2Scaled = multAff(aff2, coeff)
        aff1Copy = dict(aff1)
        aff1Copy.pop(var)
        return addAff(aff1Copy, aff2Scaled)
    else:
        return aff1

# nature of data: pairs of lower bound / upper bound
# add two intervals
def addInt(int1, int2):
    return (int1[0]+int2[0], int1[1]+int2[1])
# sub two intervals
def subInt(int1, int2):
    return (int1[0]-int2[1], int1[1]-int2[0])
# mul two intervals
def mulInt(int1, int2):
    vals = [int1[0]*x for x in int2] + [int1[1]*x for x in int2]
    return (min(vals), max(vals))
# div two intervals
def divInt(int1, int2):
    inf = float('inf')
    tempInt = (0.0,0.0)
    tempInt[0] = -inf if int2[1]==0.0 else 1.0/int2[1]
    tempInt[1] = inf if int2[0]==0.0 else 1.0/int2[0]
    if tempInt[0] > tempInt[1]:
        tempInt = (-inf,inf)
    return mulInt(int1, tempInt)

class MyErrorListener( ErrorListener ):
    def __init__(self):
        super(MyErrorListener, self).__init__()

    def syntaxError(self, recognizer, offendingSymbol, line, column, msg, e):
        print "Syntax Error: ", line, msg
        raise Exception("Parsing Syntax error!! : ", e)

    def reportAmbiguity(self, recognizer, dfa, startIndex, stopIndex, exact, ambigAlts, configs):
        raise Exception("Ambiguious Syntax error!! : ")

    def reportAttemptingFullContext(self, recognizer, dfa, startIndex, stopIndex, conflictingAlts, configs):
        raise Exception("AttemptingFullContext", conflictingAlts)

    def reportContextSensitivity(self, recognizer, dfa, startIndex, stopIndex, prediction, configs):
        raise Exception("Oh no!!")

class chiselGenerator(ParallelyVisitor):
    def __init__(self, checker_spec, ifs):
        self.typecontext = {}
        self.processgroups = {}
        self.checker_spec = checker_spec
        self.ifs = ifs

    def visitSequential(self, ctx):
        type1 = self.visit(ctx.getChild(0))
        type2 = self.visit(ctx.getChild(2))
        return (type1 and type2)

    def visitSkipstatement(self, ctx):
        return True

    def visitLiteral(self, ctx):
        return ([],0.0)

    def visitFliteral(self, ctx):
        return ([],0.0)

    def visitVariable(self, ctx):
        return ([(ctx.getText(),1.0)],0.0)

    def visitMultiply(self, ctx):
        e_1 = self.visit(ctx.expression(0))
        e_2 = self.visit(ctx.expression(1))
        maxval = 10.0
        for pair in e_1[0]:
          pair[1] *= maxval
        e_1[1] *= maxval
        for pair in e_2[0]:
          pair[1] *= maxval
        e_2[1] *= maxval
        return (e_1[0]+e_2[0],e_1[1]+e_2[1])

    def visitSelect(self, ctx):
        return self.visit(ctx.expression())

    def visitDivide(self, ctx):
        e_1 = self.visit(ctx.expression(0))
        e_2 = self.visit(ctx.expression(1))
        maxval = 10.0
        for pair in e_1[0]:
          pair[1] /= maxval
        e_1[1] /= maxval
        for pair in e_2[0]:
          pair[1] /= maxval
        e_2[1] /= maxval
        return (e_1[0]+e_2[0],e_1[1]+e_2[1])

    def visitAdd(self, ctx):
        e_1 = self.visit(ctx.expression(0))
        e_2 = self.visit(ctx.expression(1))
        return (e_1[0]+e_2[0],e_1[1]+e_2[1])

    def visitMinus(self, ctx):
        e_1 = self.visit(ctx.expression(0))
        e_2 = self.visit(ctx.expression(1))
        return (e_1[0]+e_2[0],e_1[1]+e_2[1])

    def updateSpec(self, spec, ctx, assigned_var, vars_list, multiplicatives):
        new_spec = []
        # print "Updating spec", assigned_var, vars_list, multiplicatives, spec
        for spec_part in spec:
            temp_mul = spec_part.multiplicative * multiplicatives
            for errConstraint in spec_part.jointreliability:
                for varDelta in errConstraint:
                  if varDelta[0] == assigned_var:
                    multiplier = varDelta[1]
                    #
                
                
                
            if assigned_var in spec_part.jointreliability:
                temp_joinedrel.remove(assigned_var)
                temp_joinedrel.update(vars_list)
            new_spec.append(Constraint(spec_part.limit, spec_part.condition, temp_mul, temp_joinedrel))
        # print "--->", new_spec
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
        # print assigned_var
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
        new_items = e_1
        return self.updateSpec(spec, ctx, assigned_var,
                               new_items, p)

    def processRecover(self, ctx, spec):

        if self.ifs:
            newspec = []
            spec_try = self.processspec(ctx.trys, spec)
            spec_recover = self.processspec(ctx.recovers, spec)
            for i, spec_part in enumerate(spec):
                s1_data = spec_try[i].jointreliability
                s2_data = spec_recover[i].jointreliability
                all_data = s1_data | s2_data  # Set union!!!!

                new_mult = min(spec_try[i].multiplicative, spec_recover[i].multiplicative)
                newConstraint = Constraint(spec_part.limit,
                                           spec_part.condition,
                                           new_mult,
                                           all_data)
                newspec.append(newConstraint)
            return newspec

        ps1calculator = CalculatePSuccess()
        ps1 = ps1calculator.calc(ctx.trys)
        spec_try = self.processspec(ctx.trys, spec)
        spec_recover = self.processspec(ctx.recovers, spec)
        # ps2 = spec_recover[0].multiplicative / spec[0].multiplicative

        checker_f = ctx.check.getText()

        checker_f_spec = {"TP": 1, "TN": 1}

        # print checker_f, self.checker_spec
        if checker_f in self.checker_spec:
            checker_f_spec = self.checker_spec[checker_f]

        # print checker_f_spec

        newspec = []
        for i, spec_part in enumerate(spec):
            s1_data = spec_try[i].jointreliability
            s2_data = spec_recover[i].jointreliability
            all_data = s1_data | s2_data  # Set union!!!!

            # print s1_data
            # print s2_data
            # print all_data

            # Calculate the new multiplication
            temp1 = ps1 * checker_f_spec['TN'] * spec_part.multiplicative
            temp2 = ps1 * (1 - checker_f_spec['TN']) * spec_recover[i].multiplicative
            temp3 = (1 - ps1) * spec_recover[i].multiplicative * checker_f_spec['TP']

            new_mult = temp1 + temp2 + temp3
            # print new_mult
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

    # We are treating the conditionals a non-deterministic choise
    def processIf(self, statement, spec):
        if_branch = statement.ifs
        else_branch = statement.elses
        # b_cond = statement.var().getText()
        if_spec = self.processspec(if_branch, spec)
        else_spec = self.processspec(else_branch, spec)
        # print if_spec, else_spec
        newspec = []
        for i, spec_part in enumerate(spec):
            s1_data = if_spec[i].jointreliability
            s2_data = else_spec[i].jointreliability
            all_data = s1_data | s2_data  # Set union. Magic !!!!
            new_mult = min(if_spec[i].multiplicative, else_spec[i].multiplicative)
            newConstraint = Constraint(spec_part.limit,
                                       spec_part.condition,
                                       new_mult,
                                       all_data)
            newspec.append(newConstraint)
        return newspec

    def processWhile(self, statement, spec):
        unrelVarsFinder = SimpleFindUnrelVars()
        unrelVars = unrelVarsFinder.simpleUnrelVars(statement.body)

        unrelspec = []
        relspec = []
        for i, spec_part in enumerate(spec):
            # if unrel modified var in joint rel, rel becomes 0
            if unrelVars.intersection(spec_part.jointreliability):
                newConstraint = Constraint(spec_part.limit,
                                           spec_part.condition,
                                           0,
                                           set())
                unrelspec.append(newConstraint)
            else:
                relspec.append(spec_part)
        #TODO improve fixed point analysis
        for _ in range(10): #TODO remove 10 iteration safety limit
            newrelspec = self.processspec(statement.body, relspec)
            if False: #set(newrelspec)==set(relspec): #TODO intelligent stopping the fixpoint analysis
                break
            else:
                relspec=newrelspec

        return relspec+unrelspec

    def processspec(self, statements, spec):
        reversed_statements = statements[::-1]
        for i, statement in enumerate(reversed_statements):
            # self.visit(statement)
            if isinstance(statement, ParallelyParser.CastContext):
                spec = self.processCast(statement, spec)
            elif isinstance(statement, ParallelyParser.ProbassignmentContext):
                spec = self.processProbassignment(statement, spec)
            elif isinstance(statement, ParallelyParser.ApproximateContext):
                spec = self.processApproximate(statement, spec)
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
            elif isinstance(statement, ParallelyParser.IfContext):
                spec = self.processIf(statement, spec)
            elif isinstance(statement, ParallelyParser.WhileContext):
                spec = self.processWhile(statement, spec)
            else:
                print "Unable to process the statement :", statement.getText()
                exit(-1)
            # print "Processed : {} :> {} ({}/{})".format(statement.getText(), spec, i, len(reversed_statements))
        return spec


# Takes in a .seq file performs the chisel reliability analysis
def main(program_str, spec, skiprename, checker_spec, ifs):
    input_stream = InputStream(program_str)
    # lexer = ParallelyLexer(input_stream)
    # stream = CommonTokenStream(lexer)
    # parser = ParallelyParser(stream)

    start = time.time()
    # print "Unrolling Repeat statements?: ", (not skiprename)
    replacement = 0
    replacement_map = {}
    if not skiprename:
        i = 0
        while(True):
            # print "unrolling {} deep".format(i)

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

            unroller = unrollRepeat(stream, replacement, replacement_map)
            walker = ParseTreeWalker()
            walker.walk(unroller, tree)
            input_stream = InputStream(unroller.rewriter.getDefaultText())
            replacement = unroller.replacement
            # print unroller.replacement, unroller.dummymap
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

        # print "----------------------------------------"
        # print "Intermediate step. Writing to _DEBUG_UNROLLED_.txt"

        unroller = unrollRepeat(stream, replacement - 1, replacement_map)
        new_program = unroller.replace_dummies(input_stream.strdata)

        # print new_program[:100]

        debug_file = open("_DEBUG_UNROLLED_.txt", 'w')
        debug_file.write(new_program)
        debug_file.close()
        # print "----------------------------------------"
    else:
        new_program = input_stream.strdata

    start2 = time.time()

    # print input_stream.strdata
    input_stream = InputStream(new_program)
    lexer = ParallelyLexer(input_stream)
    stream = CommonTokenStream(lexer)
    parser = ParallelyParser(stream)
    parser.addErrorListener(MyErrorListener())
    parser._interp.predictionMode = PredictionMode.SLL

    try:
        tree = parser.parallelprogram()
    except Exception as e:
        print "Parsing Error!!!"
        print e
        exit(-1)

    start3 = time.time()

    # Processing the spec
    spec_input_stream = InputStream(spec)
    spec_lexer = ParallelyLexer(spec_input_stream)
    spec_stream = CommonTokenStream(spec_lexer)
    spec_parser = ParallelyParser(spec_stream)

    spec_str = spec_parser.chiselspec()
    # print spec_str.getText()

    chisel_spec = []
    for constraint_str in spec_str.singlechiselspec():
        temp_limit = float(constraint_str.FLOAT(0).getText())
        temp_mult = 1.0
        var_list = []
        for var, maxerr in zip(constraint_str.var(), constraint_str.FLOAT()[1:]):
            var_list.append((float(maxerr.getText()),{var.getText():1.0}))
        chisel_spec.append(Constraint(temp_limit, "<=", temp_mult, var_list))

    initial_var_int = {}
    for interval_spec in spec_str.varchiselspec():
        var = interval_spec.var(0).getText()
        interval = tuple(float(num.getText()) for num in interval_spec.interval(0).FLOAT())
        initial_var_int[var] = interval

    func_specs = {}
    for func_spec in spec_str.funcchiselspec():
        func = func_spec.var(0).getText()
        out_delta = float(func_spec.FLOAT(0).getText())
        out_interval = [float(num.getText()) for num in func_spec.interval(0).FLOAT()]
        out_rel = float(func_spec.FLOAT(1).getText())
        arg_dict = {}
        for var_obj, maxerr_obj, var_interval_obj in zip(func_spec.var()[1:], func_spec.FLOAT()[2:], func_spec.interval()[1:]):
            var = var_obj.getText()
            maxerr = float(maxerr_obj.getText())
            var_interval = [float(num.getText()) for num in var_interval_obj.FLOAT()]
            arg_dict[var] = tuple(var_interval + [maxerr])
        func_specs[func] = (tuple(out_interval + [out_delta]), out_rel, arg_dict)

    chisel = chiselGenerator(checker_spec, ifs)
    result_spec = chisel.processspec(tree.program(0).statement(), chisel_spec)

    # if the variable declaration is found the reliability is 1
    decs = tree.program(0).declaration()
    result_spec = chisel.processspec(decs, result_spec)
    end = time.time()

    # print '----------------------------------------'
    print result_spec
    # print '----------------------------------------'
    print "Analysis time Total: {}, Unroll: {}, chisel: {}".format(end - start, start2 - start, end - start3)

    return result_spec


if __name__ == '__main__':
    sys.setrecursionlimit(15000)

    parser = ArgumentParser()
    parser.add_argument("-f", dest="programfile",
                        help="File containing the code", required=True)
    parser.add_argument("-s", dest="spec",
                        help="File to output the sequential code", required=True)
    parser.add_argument("-du", "--dontunroll", action="store_true",
                        help="Unroll the loops", default=False)
    parser.add_argument("-ifs", "--tryisif", action="store_true",
                        help="Treat try block as an if", default=False)
    parser.add_argument("-func", dest="functionspec",
                        help="specification of functions", default=None)
    parser.add_argument("-d", "--debug", action="store_true",
                        help="Print debug info")

    args = parser.parse_args()

    programfile = open(args.programfile, 'r')
    spec = open(args.spec, 'r').read()
    dontunroll = args.dontunroll

    # print sys.argv
    checker_spec = {}

    if args.functionspec:
        checker_spec = json.loads(open(args.functionspec, 'r').read())
        # print "Using the checker functions :", checker_spec

    program_str = programfile.read()
    main(program_str, spec, dontunroll, checker_spec, args.tryisif)

