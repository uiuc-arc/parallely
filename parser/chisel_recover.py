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
        aff1Copy = copy.deepcopy(aff1)
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
# merge two intervals
def mergeInt(int1, int2):
    lower = min(int1[0], int2[0])
    upper = max(int1[1], int2[1])
    return (lower, upper)
# check if interval contains another
def containsInt(intOuter, intInner):
    return intOuter[0]<=intInner[0] and intOuter[1]>=intInner[1]

class intervalAnalysis(ParallelyVisitor):
    def __init__(self, func_specs):
        self.func_specs = func_specs

    def visitLiteral(self, ctx):
        value = float(ctx.getText())
        ctx.interval = (value,value)

    def visitFliteral(self, ctx):
        value = float(ctx.getText())
        ctx.interval = (value,value)

    def visitVariable(self, ctx):
        ctx.interval = copy.deepcopy(var_int[ctx.getText()])

    def visitSubExpAndGetIntervals(self, ctx):
        self.visit(ctx.expression(0))
        self.visit(ctx.expression(1))
        int1 = ctx.expression(0).interval
        int2 = ctx.expression(1).interval
        return (int1, int2)

    def visitAdd(self, ctx):
        int1, int2 = visitSubExpAndGetIntervals(self, ctx)
        ctx.interval = addInt(int1, int2)

    def visitMinus(self, ctx):
        int1, int2 = visitSubExpAndGetIntervals(self, ctx)
        ctx.interval = subInt(int1, int2)

    def visitMultiply(self, ctx):
        int1, int2 = visitSubExpAndGetIntervals(self, ctx)
        ctx.interval = mulInt(int1, int2)

    def visitDivide(self, ctx):
        int1, int2 = visitSubExpAndGetIntervals(self, ctx)
        ctx.interval = divInt(int1, int2)

    def visitGreater(self, ctx):
        int1, int2 = visitSubExpAndGetIntervals(self, ctx)
        if not (int1[1]<int2[0] or int1[0]>int2[1]):
            raise Exception("Unclear outcome of boolean check!")

    def visitLess(self, ctx):
        int1, int2 = visitSubExpAndGetIntervals(self, ctx)
        if not (int1[1]<int2[0] or int1[0]>int2[1]):
            raise Exception("Unclear outcome of boolean check!")

    def visitGEQ(self, ctx):
        int1, int2 = visitSubExpAndGetIntervals(self, ctx)
        if not (int1[1]<int2[0] or int1[0]>int2[1]):
            raise Exception("Unclear outcome of boolean check!")

    def visitLEQ(self, ctx):
        int1, int2 = visitSubExpAndGetIntervals(self, ctx)
        if not (int1[1]<int2[0] or int1[0]>int2[1]):
            raise Exception("Unclear outcome of boolean check!")

    def visitEqual(self, ctx):
        int1, int2 = visitSubExpAndGetIntervals(self, ctx)
        if not (int1[0]==int1[1] and int2[0]==int2[1] and int1[0]==int2[0]):
            raise Exception("Unclear outcome of boolean check!")

    def visitSelect(self, ctx):
        self.visit(ctx.expression(0))

    def visitAndexp(self, ctx):
        self.visit(ctx.expression(0))
        self.visit(ctx.expression(1))

    def processProbassignment(self, ctx):
        var = ctx.var(0).getText()
        self.visit(ctx.expression(0))
        var_int[var] = ctx.expression(0).interval

    def processExpassignment(self, ctx):
        var = ctx.var(0).getText()
        self.visit(ctx.expression(0))
        var_int[var] = ctx.expression(0).interval

    def processFunction(self, ctx):
        outputs = [obj.getText() for obj in ctx.var()[:-1]]
        func = ctx.var()[-1].getText()
        args = ctx.expression()
        funcspec = self.func_specs[func]
        # check arguments to make sure they are within the function's domain
        for i, argSpec in enumerate(funcspec[1]):
            self.visit(args[i])
            arg_int = args[i].interval
            if not containsInt(argSpec, arg_int):
                raise Exception("Function call argument interval is not contained within function parameter interval!")
        for i, output in enumerate(outputs):
            var_int[output] = funcspec[2][i][:2]

    def processRecover(self, ctx):
        # treats tcr block as ite block
        # get statements
        ifStmts = ctx.trys()
        elseStmts = ctx.recovers()
        # backup current intervals (deep copy)
        var_int_current = copy.deepcopy(var_int)
        # analyze if branch
        self.analyze(ifStmts, var_int)
        # store if branch exit intervals (deep copy)
        var_int_if = copy.deepcopy(var_int)
        # restore to current intervals (shallow copy)
        var_int = var_int_current
        # analyze else branch
        self.analyze(elseStmts, var_int)
        # merge if branch exit intervals into else branch exit intervals
        for var in var_int:
            var_int[var] = merge(var_int[var], var_int_if[var])

    def processDec(self, ctx):
        var = ctx.var(0).getText()
        var_int[var] = (float('-inf'),float('inf'))

    def processIf(self, ctx):
        # get statements
        ifStmts = ctx.ifs()
        elseStmts = ctx.elses()
        # backup current intervals (deep copy)
        var_int_current = copy.deepcopy(var_int)
        # analyze if branch
        self.analyze(ifStmts, var_int)
        # store if branch exit intervals (deep copy)
        var_int_if = copy.deepcopy(var_int)
        # restore to current intervals (shallow copy)
        var_int = var_int_current
        # analyze else branch
        self.analyze(elseStmts, var_int)
        # merge if branch exit intervals into else branch exit intervals
        for var in var_int:
            var_int[var] = merge(var_int[var], var_int_if[var])

    def analyze(self, statements, var_int):
        self.var_int = var_int
        for i, statement in enumerate(statements):
            if isinstance(statement, ParallelyParser.ProbassignmentContext):
                self.processProbassignment(statement)
            # elif isinstance(statement, ParallelyParser.CastContext):
            #     self.processCast(statement)
            # elif isinstance(statement, ParallelyParser.ApproximateContext):
            #     self.processApproximate(statement)
            elif isinstance(statement, ParallelyParser.ExpassignmentContext):
                self.processExpassignment(statement)
            elif isinstance(statement, ParallelyParser.FuncContext):
                self.processFunction(statement)
            elif isinstance(statement, ParallelyParser.RecoverContext):
                self.processRecover(statement)
            elif isinstance(statement, ParallelyParser.SingledeclarationContext):
                self.processDec(statement)
            elif isinstance(statement, ParallelyParser.IfContext):
                self.processIf(statement)
            else:
                print "Unable to process the statement :", statement.getText()
                exit(-1)
        return spec

class chiselGenerator(ParallelyVisitor):
    def __init__(self, func_specs, checker_specs):
        self.func_specs = func_specs
        self.checker_specs = checker_specs

    def visitLiteral(self, ctx):
        return ({}, set())

    def visitFliteral(self, ctx):
        return ({}, set())

    def visitVariable(self, ctx):
        var = ctx.getText()
        return ({var:1}, {var})

    def visitAdd(self, ctx):
        op1 = self.visit(ctx.expression(0))
        op2 = self.visit(ctx.expression(1))
        aff = addAff(op1[0], op2[0])
        varset = op1[1].union(op2[1])
        return (aff, varset)

    def visitMinus(self, ctx):
        op1 = self.visit(ctx.expression(0))
        op2 = self.visit(ctx.expression(1))
        aff = addAff(op1[0], op2[0])
        varset = op1[1].union(op2[1])
        return (aff, varset)

    def visitMultiply(self, ctx):
        op1 = self.visit(ctx.expression(0))
        op2 = self.visit(ctx.expression(1))
        int1 = ctx.expression(0).interval
        int2 = ctx.expression(1).interval
        maxop1 = max(abs(int1[0]), abs(int1[1]))
        maxop2 = max(abs(int2[0]), abs(int2[1]))
        aff = addAff(multaff(op1[0],maxop2), multaff(op2[0],maxop1))
        varset = op1[1].union(op2[1])
        return (aff, varset)

    def visitDivide(self, ctx):
        op1 = self.visit(ctx.expression(0))
        op2 = self.visit(ctx.expression(1))
        int1 = ctx.expression(0).interval
        int2 = ctx.expression(1).interval
        # TODO what to do here?
        raise Exception('Division not implemented!')
        # varset = op1[1].union(op2[1])
        # return (aff, varset)

    def visitSelect(self, ctx):
        return self.visit(ctx.expression(0))

    def processProbassignment(self, ctx, spec):
        var = ctx.var(0).getText()
        expData = self.visit(ctx.expression(0))
        newspec = []
        for constraint in spec:
            if any([(var in comparison[1]) for comparison in constraint.jointreliability]):
                newmultiplicative = constraint.multiplicative * float(ctx.probability(0).getText)
                newjointreliability = []
                for comparison in constraint.jointreliability:
                    newRHS = replaceAff(comparison[1], var, expData[0])
                    newjointreliability.append([comparison[0],newRHS])
                newspec.append(Constraint(constraint.limit, constraint.condition, newmultiplicative, newjointreliability)
            else:
                newspec.append(constraint)
        return newspec

    def processExpassignment(self, ctx, spec):
        var = ctx.var(0).getText()
        expData = self.visit(ctx.expression(0))
        newspec = []
        for constraint in spec:
            if any([(var in comparison[1]) for comparison in constraint.jointreliability]):
                newjointreliability = []
                for comparison in constraint.jointreliability:
                    newRHS = replaceAff(comparison[1], var, expData[0])
                    newjointreliability.append([comparison[0],newRHS])
                newspec.append(Constraint(constraint.limit, constraint.condition, constraint.multiplicative, newjointreliability)
            else:
                newspec.append(constraint)
        return newspec

    def processFunction(self, ctx, spec):
        func = ctx.var()[-1].getText()
        func_spec = self.func_specs[func]
        func_rel = func_spec[0]
        outputs = [var.getText() for var in ctx.var()[:-1]]
        maxerrors = [func_spec[1][i][2] for i in range(len(outputs))]
        newspec = []
        # new output constraints
        for constraint in spec:
            if any([any([(output in comparison[1]) for output in outputs]) for comparison in constraint.jointreliability]):
                newmultiplicative = constraint.multiplicative * func_rel
                newjointreliability = []
                for comparison in constraint.jointreliability:
                    newRHS = comparison[1]
                    for i, output in enumerate(outputs):
                        newRHS = replaceAff(newRHS, output, {1:maxerrors[i]})
                    newjointreliability.append([comparison[0],newRHS])
                newspec.append(Constraint(constraint.limit, constraint.condition, newmultiplicative, newjointreliability)
            else:
                newspec.append(constraint)
        # new input constraints
        newjointreliability = []
        for i, arg in enumerate(func_spec[2]):
            argData = self.visit(ctx.expression(i))
            newjointreliability.append([arg[2],argData])
        newspec.append(Constraint(1.0, "<=", 1.0, newjointreliability))
        return newspec

    def processRecover(self, ctx, spec):
        #TODO make it actually do recovery stuff
        ifspec = self.processspec(ctx.ifs(), spec)
        elsespec = self.processspec(ctx.elses(), spec)
        return ifspec + elsespec

    def processDec(self, ctx, spec):
        var = ctx.var(0).getText()
        newspec = []
        for constraint in spec:
            if any([(var in comparison[1]) for comparison in constraint.jointreliability]):
                newjointreliability = []
                for comparison in constraint.jointreliability:
                    newRHS = replaceAff(comparison[1], var, {1:float('inf')})
                    newjointreliability.append([comparison[0],newRHS])
                newspec.append(Constraint(constraint.limit, constraint.condition, newmultiplicative, newjointreliability)
            else:
                newspec.append(constraint)
        return newspec

    def processIf(self, ctx, spec):
        ifspec = self.processspec(ctx.ifs(), spec)
        elsespec = self.processspec(ctx.elses(), spec)
        return ifspec + elsespec

    def processspec(self, statements, spec):
        reversed_statements = statements[::-1]
        for i, statement in enumerate(reversed_statements):
            if isinstance(statement, ParallelyParser.ProbassignmentContext):
                spec = self.processProbassignment(statement, spec)
            # elif isinstance(statement, ParallelyParser.CastContext):
            #     spec = self.processCast(statement, spec)
            # elif isinstance(statement, ParallelyParser.ApproximateContext):
            #     spec = self.processApproximate(statement, spec)
            elif isinstance(statement, ParallelyParser.ExpassignmentContext):
                spec = self.processExpassignment(statement, spec)
            elif isinstance(statement, ParallelyParser.FuncContext):
                spec = self.processFunction(statement, spec)
            elif isinstance(statement, ParallelyParser.RecoverContext):
                spec = self.processRecover(statement, spec)
            elif isinstance(statement, ParallelyParser.SingledeclarationContext):
                spec = self.processDec(statement, spec)
            elif isinstance(statement, ParallelyParser.IfContext):
                spec = self.processIf(statement, spec)
            else:
                print "Unable to process the statement :", statement.getText()
                exit(-1)
        return spec


# Takes in a .seq file performs the chisel accuracy analysis
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
            var_list.append([float(maxerr.getText()),{var.getText():1.0}])
        chisel_spec.append(Constraint(temp_limit, "<=", temp_mult, var_list))

    initial_var_int = {}
    for interval_spec in spec_str.varchiselspec():
        var = interval_spec.var(0).getText()
        interval = tuple(float(num.getText()) for num in interval_spec.interval(0).FLOAT())
        initial_var_int[var] = interval

    func_specs = {}
    for func_spec in spec_str.funcchiselspec():
        func = func_spec.var(0).getText()
        numIn = int(func_spec.INT(0).getText())
        numOut = int(func_spec.INT(1).getText())
        rel = float(func_spec.FLOAT(0).getText())
        arg_list = []
        ret_list = []
        for i, maxerr_obj, interval_obj in zip(range(numIn+numOut), func_spec.FLOAT()[1:], func_spec.interval()):
            maxerr = float(maxerr_obj.getText())
            interval = [float(num.getText()) for num in interval_obj.FLOAT()]
            combined = tuple(interval + [maxerr])
            if i < numIn:
                arg_list.append(combined)
            else:
                ret_list.append(combined)
        func_specs[func] = (rel, arg_list, ret_list)

    checker_specs = {}
    for checker_spec in spec_str.checkerchiselspec():
        func = func_spec.var(0).getText()
        ret_list = []
        for exp_obj, var_obj in zip(spec_str.expression(), spec_str.var()[1:]):
            var = var_obj.getText() #unused - we use positional return vals
            ret_list.append(exp_obj)
        checker_specs[func] = ret_list

    intervalAnalysisInstance = intervalAnalysis(func_specs)
    intervalAnalysisInstance.analyze(tree.program(0).statement(), initial_var_int)

    chisel = chiselGenerator(func_specs, checker_specs)
    result_spec = chisel.processspec(tree.program(0).statement(), chisel_spec)

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
