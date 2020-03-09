from ParallelyParser import ParallelyParser
from ParallelyVisitor import ParallelyVisitor
from ParallelyLexer import ParallelyLexer
from antlr4 import InputStream
from antlr4 import CommonTokenStream

from argparse import ArgumentParser
import time


class dieselStaticOptimizer(ParallelyVisitor):
    def __init__(self, debug, annotate):
        self.statement_lists = {}
        self.msgcontext = {}
        self.globaldecs = {}
        self.declarations = []
        self.grouped_list = {}
        self.isProcessGroup = {}
        self.debug = debug
        self.annotate = annotate

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

    def debugMsg(self, msg):
        if self.debug:
            print("[Debug - TypeChecker] " + msg)

    def getAnnotationStr(self, annotatedstatement):
        # print annotatedstatement.annotation[0].text
        return ','.join([i.text for i in annotatedstatement.annotation])

    def getAnnotationList(self, annotatedstatement):
        # print annotatedstatement.annotation[0].text
        return [i.text for i in annotatedstatement.annotation]

    def getProcessFromVarName(self, varname):
        return varname.split('$')[-1]

    def getLocation(self, annotatedstatement):
        if isinstance(annotatedstatement, ParallelyParser.AnnotatedContext):
            alist = self.getAnnotationList(annotatedstatement)
            return alist[0]
        else:
            return -1

    def isValid(self, relfactor):
        if len(set([self.getProcessFromVarName(var) for var in relfactor])) == 1:
            return True
        return False

    def attemptOpt(self, statements, relfactor, prob):
        # checkstatement = annotated_checkstatement.statement()
        # relfactor = set([i.getText() for i in checkstatement.rel_factor])
        # prob = float(checkstatement.probability().getText())
        removed_statements = 0

        new_relfactor = relfactor
        new_location = -1
        new_s = False

        for annotatedstatement in statements[::-1]:
            if isinstance(annotatedstatement, ParallelyParser.AnnotatedContext):
                statement = annotatedstatement.statement()
            else:
                statement = annotatedstatement
            current_relfactor = set(relfactor)
            if isinstance(statement, ParallelyParser.ExpassignmentContext):
                success, relfactor, prob = self.processExpassignment(statement,
                                                                     current_relfactor, prob)
            elif isinstance(statement, ParallelyParser.IfContext):
                success, relfactor, prob = self.processCondition(statement,
                                                                 current_relfactor, prob)
            else:
                break

            # If we were table to move beyond this instruction
            if success:
                new_relfactor = relfactor

                new_prob = prob
                new_s = True
                new_location = self.getLocation(annotatedstatement)
                # err_msg = "{} in line {} => to check({}, {}) in line {}"
                # print(err_msg.format(checkstatement.getText(),
                #                      self.getAnnotationStr(annotated_checkstatement),
                #                      ', '.join(relfactor),
                #                      prob, new_location))
                removed_statements += 1
                # print success, relfactor, prob, self.isValid(relfactor)
            else:
                break

        print "------------------------------------------"
        print new_s, new_relfactor, new_prob, new_location
        print "------------------------------------------"        
        return new_s, new_relfactor, new_prob, new_location

    # Simple substituion because only one variable is assigned here
    def substitute(self, relfactor, substituted, substituition):
        relfactor.remove(substituted)
        relfactor.update(substituition)

        # We are treating the conditionals a non-deterministic choise
    def processCondition(self, statement, relfactor, prob):
        if_branch = statement.ifs
        else_branch = statement.elses

        s_if, relfactor_if, prob_if, location_if = self.attemptOpt(if_branch,
                                                                   relfactor, prob)
        s_else, relfactor_else, prob_else, location_else = self.attemptOpt(else_branch,
                                                                           relfactor, prob)

        if s_if and s_else:
            print(relfactor_if, relfactor_else, prob_if, prob_else, relfactor_if.union(relfactor_else),
                  min(prob_if, prob_else))
            return True, relfactor_if.union(relfactor_else), min(prob_if, prob_else)
            exit(-1)
        # newspec = []
        # for i, spec_part in enumerate(spec):
        #     s1_data = if_spec[i].jointreliability
        #     s2_data = else_spec[i].jointreliability
        #     all_data = s1_data | s2_data  # Set union. Magic !!!!
        #     new_mult = min(if_spec[i].multiplicative, else_spec[i].multiplicative)
        #     newConstraint = Constraint(spec_part.limit,
        #                                spec_part.condition,
        #                                new_mult,
        #                                all_data)
        #     newspec.append(newConstraint)
        # return newspec

    def processExpassignment(self, ctx, relfactor, prob):
        assigned_var = ctx.var().getText()

        if assigned_var in relfactor:
            vars_list = set(self.visit(ctx.expression()))
            # print vars_list
            # if isinstance(ctx.expression().VAR(), ParallelyParser.LocalvariableContext):
            #     vars_list = set([ctx.expression().VAR().getText()])
            # else:
            #     vars_list = set([i.getText() for i in ctx.expression().var()])
            self.substitute(relfactor, assigned_var, vars_list)
            return True, relfactor, prob
        else:
            return True, relfactor, prob
        # print new_spec, assigned_var, vars_list
        return False, relfactor, prob

    def rewriteProgram(self, tree, outfile):
        statement_list = tree.statement()
        # Assumption: Only work on annotated stuff
        for i, annotatedstatement in enumerate(statement_list):
            statement = annotatedstatement.statement()
            # print i, statement.getText(), type(statement)
            # if isinstance(statement, ParallelyParser.SpeccheckarrayContext):
            #     print statement.getText()
            if isinstance(statement, ParallelyParser.SpeccheckContext):
                # checkstatement = annotated_checkstatement.statement()
                relfactor = set([var.getText() for var in statement.rel_factor])
                prob = float(statement.probability().getText())

                print("Found a check function {} on line : {}".format(statement.getText(),
                                                                      statement.start.line))
                new_s, new_relfactor, new_prob, new_location = self.attemptOpt(statement_list[:i],
                                                                               relfactor, prob)
                if new_s:
                    msg = "Move: {} in line {} => to check({}, {}) in line {}"
                    print(msg.format(statement.getText(),
                                     self.getAnnotationStr(annotatedstatement),
                                     ', '.join(new_relfactor),
                                     new_prob, new_location))


if __name__ == '__main__':
    parser = ArgumentParser()
    parser.add_argument("-f", dest="programfile",
                        help="File containing sequentialized code", required=True)
    parser.add_argument("-o", dest="outfile",
                        help="File to output the optimized code", required=True)
    parser.add_argument("-d", "--debug", action="store_true",
                        help="Print debug info")
    parser.add_argument("-g", "--annotate", action="store_true",
                        help="annotate with debug info")

    args = parser.parse_args()
    programfile = open(args.programfile, 'r')
    outfile = open(args.outfile, 'w')
    program_str = programfile.read()

    input_stream = InputStream(program_str)
    lexer = ParallelyLexer(input_stream)
    stream = CommonTokenStream(lexer)
    parser = ParallelyParser(stream)

    tree = parser.sequentialprogram()

    print ("Parsed the sequentialized program...")

    # Sequentialization
    start2 = time.time()
    optimizer = dieselStaticOptimizer(args.debug, args.annotate)
    optimizer.rewriteProgram(tree, outfile)
    end2 = time.time()
    print "Time for optimization :", end2 - start2
