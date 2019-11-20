from GoLangParser import GoLangParser
from GoLangLexer import GoLangLexer
from GoLangVisitor import GoLangVisitor
from argparse import ArgumentParser
from antlr4 import CommonTokenStream
from antlr4 import InputStream
import sys
from antlr4.error.ErrorListener import ErrorListener
from antlr4 import *


# class MyErrorListener(ErrorListener):
#     def syntaxError(self, recognizer, offendingSymbol, line, column, msg, e):
#         print str(line) + ":" + str(column) + ": syntax ERROR, " + str(msg)
#         print "Terminating Translation"
#         sys.exit()

#     def reportAmbiguity(self, recognizer, dfa, startIndex, stopIndex, exact, ambigAlts, configs):
#         print "Ambiguity ERROR, " + str(configs)
#         sys.exit()

#     def reportAttemptingFullContext(self, recognizer, dfa, startIndex, stopIndex, conflictingAlts, configs):
#         print "Attempting full context ERROR, " + str(configs)
#         sys.exit()

#     def reportContextSensitivity(self, recognizer, dfa, startIndex, stopIndex, prediction, configs):
#         print "Context ERROR, " + str(configs)
#         sys.exit()

def isFunction(listidentifiers, fname):
    return fname == '.'.join([i.getText() for i in listidentifiers])

def getThreadSet(func):
    thread_names = []
    for stat in func.block().statementList().statement():
        # print stat.getText(), type(stat)
        if isinstance(stat, GoLangParser.SmtfunctionContext):
            if isFunction(stat.IDENTIFIER(), "parallely.LaunchThread"):
                # print "===: ", stat.getText(), stat.IDENTIFIER()
                thread_names.append((stat.arguments().expressionList().expression()[0].getText(), 0))
            if isFunction(stat.IDENTIFIER(), "parallely.LaunchThreadGroup"):
                thread_names.append((stat.arguments().expressionList().expression()[0].getText(),
                                     stat.arguments().expressionList().expression()[1].getText()))
    return thread_names


def main(program_str, args):
    input_stream = InputStream(program_str)
    lexer = GoLangLexer(input_stream)
    stream = CommonTokenStream(lexer)
    parser = GoLangParser(stream)
    # parser._listeners = [MyErrorListener()]
    parser._interp.predictionMode = PredictionMode.SLL

    tree = parser.sourceFile()

    functions = {}
    for func in tree.functionDecl():
        functions[func.IDENTIFIER().getText()] = func.block()
        if func.IDENTIFIER().getText()=="main":
            spawned_threads = getThreadSet(func)
            print spawned_threads

    for 

if __name__ == '__main__':
    parser = ArgumentParser()
    parser.add_argument("-f", dest="programfile",
                        help="File containing the code", required=True)
    args = parser.parse_args()
    programfile = open(args.programfile, 'r')
    program_str = programfile.read()
    main(program_str, args)
