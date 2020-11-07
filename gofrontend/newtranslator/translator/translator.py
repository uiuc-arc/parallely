from ..antlrgenerated.GoParser import GoParser
from newtranslator.antlrgenerated.GoLexer import GoLexer
from ..antlrgenerated import GoParserVisitor
from argparse import ArgumentParser
from antlr4 import CommonTokenStream
from antlr4 import InputStream
import sys
from antlr4.error.ErrorListener import ErrorListener
from antlr4 import *
from antlr4 import Lexer

def isFunction(listidentifiers, fname):
    return fname == '.'.join([i.getText() for i in listidentifiers])

def noTerminatorBetween(TokenOffset):
    return True

def lineTerminatorAhead():
    possibleIndexEosToken = self.getCurrentToken().tokenIndex - 1;

    if possibleIndexEosToken == -1:
        return True

    ahead = self._input.get(possibleIndexEosToken)
    if ahead.channel != Lexer.HIDDEN:
        return False

    if ahead.type == GoLexer.TERMINATOR:
        return True

    if ahead.type == GoLexer.WS:
        possibleIndexEosToken = self.getCurrentToken().tokenIndex - 2;
        ahead = self._input.get(possibleIndexEosToken)

    text = ahead.text
    token_type = ahead.type

def isDynamic(stream, token):
    rightTokens = [t.text for t in stream.getHiddenTokensToRight(token.stop.tokenIndex)]
    if u'/*@dynamic*/' in rightTokens:
        return True
    return False

def translateFunction(func, stream):
    print("Visiting Function: " + func.IDENTIFIER().getText())

    for decl in func.block().statementList().statement():
        if isinstance(decl, GoParser.StmtdecContext):
            print(decl.getText(), isDynamic(stream, decl))
        
def getThreadSet(func):
    thread_names = []
    is_group = []
    for stat in func.block().statementList().statement():
        if isinstance(stat, GoParser.StmtfunctionContext):
            if isFunction(stat.IDENTIFIER(), "diesel.LaunchProcess"):
                thread_names.append(stat.arguments().expressionList().expression()[0].getText())
                is_group.append(0)
            if isFunction(stat.IDENTIFIER(), "diesel.LaunchProcessGroup"):
                thread_names.append(stat.arguments().expressionList().expression()[1].getText())
                is_group.append(stat.arguments().expressionList().expression()[0].getText())             
    return (thread_names, is_group)

def lineTerminatorAhead():
        return True

def main(program_str, args):
    input_stream = InputStream(program_str)
    lexer = GoLexer(input_stream)
    stream = CommonTokenStream(lexer)
    parser = GoParser(stream)
    # parser._listeners = [MyErrorListener()]
    parser._interp.predictionMode = PredictionMode.SLL

    tree = parser.sourceFile()

    functions = {}
    spawned = []
    for func in tree.functionDecl():
        functions[func.IDENTIFIER().getText()] = func.block()
        if func.IDENTIFIER().getText()=="main":
            spawned = getThreadSet(func)
    print("Main thread launching: " + str(spawned))

    for func in tree.functionDecl():
        # print(func.IDENTIFIER().getText())
        if func.IDENTIFIER().getText() in spawned[0]:
            translateFunction(func, stream)
            # print(func.getText())

if __name__ == '__main__':
    parser = ArgumentParser()
    parser.add_argument("-f", dest="programfile",
                        help="File containing the code", required=True)
    args = parser.parse_args()
    programfile = open(args.programfile, 'r')
    program_str = programfile.read()
    main(program_str, args)
