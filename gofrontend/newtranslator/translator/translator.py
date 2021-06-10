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
from antlr4 import PredictionMode


# Helper function to print error message and exit
def EXITWITHERROR(msg):
    print("[ERROR] " + msg)
    exit(-1)

# Constants
ENDLINE = ";\n"

# Check if functions are the same by combining the identifiers
def isFunction(listidentifiers, fname):
    return fname == '.'.join([i.getText() for i in listidentifiers])


def isDynamic(stream, token):
    rightTokens = [t.text for t in stream.getHiddenTokensToRight(token.stop.tokenIndex)]
    if u'/*@dynamic*/' in rightTokens:
        return True
    return False


def translateDeclaration(dec):
    if dec.varDecl():
        if dec.varDecl().varSpec():
            if (len(dec.varDecl().varSpec()) > 1):
                EXITWITHERROR("Only one declaration per line allowed")
            varname = dec.varDecl().varSpec()[0].identifierList().getText()
            basetype = dec.varDecl().varSpec()[0].type_()

            if basetype:
                # checks if arraytype. This works because we dont support pointer, etc. here
                # We should change parally syntax to match golang more
                if basetype.typeLit():
                    arraysize = basetype.typeLit().arrayType().arrayLength().getText()
                    arraytype = basetype.typeLit().arrayType().elementType().getText()
                    return "precise {}[{}] {};\n".format(arraytype, arraysize, varname)
                else:
                    vartype = dec.varDecl().varSpec()[0].type_().typeName().getText()
                    return "precise {} {};\n".format(vartype, varname)
            else:
                # Hack: fix later
                try:
                    typestr = dec.varDecl().getText()
                    # if "process" in typestr:
                    content = typestr.split('=[]process')
                    # print(content, "{} = {};\n".format(content[0][3:], content[-1]))
                    return "{} = {};\n".format(content[0][3:], content[-1])
                except Exception as e:
                    print(e)
                    EXITWITHERROR("[Error] Unable to translate declaration: " + dec.getText())
        else:
            EXITWITHERROR("[Error] Unable to translate declaration: " + dec.getText())


class FunctionTranslator(GoParserVisitor.GoParserVisitor):
    def __init__(self):
        self.statements = []
        self.typemap = {}
        self.newDeclarations = []

    def visitAssignment(self, ctx):
        return ctx.getText() + ENDLINE

    def getTypeString(self, typeDef):
        typestr = ""
        if typeDef[2]:
            typestr += "approx "
        else:
            typestr += "precise "
        typestr += typeDef[0]
        if typeDef[1]:
            typestr += "[]"
        return typestr

    def visitNewsend(self, ctx):
        sentType = self.typemap[ctx.variable.text]
        typeStr = self.getTypeString(sentType)
        send_tmpl = "send({}, {}, {});\n".format(ctx.rec.getText(), typeStr, ctx.variable.text)
        return send_tmpl

    # def visitNoisysend(self, ctx):
    #     sentType = self.typemap[ctx.variable.text]
    #     typeStr = self.getTypeString(sentType)
    #     rel = ctx.NCHAN().getText().split('=')[-1][:-2]
    #     send_tmpl = "_temp = {} [{}] (-1);\nsend({}, {}, {});\n".format(ctx.variable.text, rel,
    #                                                                     ctx.rec.getText(),
    #                                                                     typeStr, ctx.variable.text)
    #     self.newDeclarations.append("approx int _temp;\n")
    #     # self.statements.insert(0, "approx int _temp;\n")
    #     return send_tmpl

    def visitCondsend(self, ctx):
        sentType = self.typemap[ctx.variable.text]
        typeStr = self.getTypeString(sentType)
        send_tmpl = "condsend({}, {}, {}, {});\n".format(ctx.cond.text, ctx.rec.getText(),
                                                         typeStr, ctx.variable.text)
        return send_tmpl

    def visitNoisyrec(self, ctx):
        recType = self.typemap[ctx.variable.text]
        typestr = self.getTypeString(recType)
        rel = ctx.NCHAN().getText().split('=')[-1][:-2]
        rec_str = "{0} = receive({1}, {2});\n{0}={0}[{3}](-1);\n".format(ctx.variable.text,
                                                                         ctx.sender.getText(),
                                                                         typestr, rel)
        return rec_str

    def visitRec(self, ctx):
        recType = self.typemap[ctx.variable.text]
        typestr = self.getTypeString(recType)
        rec_str = "{} = receive({}, {});\n".format(ctx.variable.text, ctx.sender.getText(), typestr)
        return rec_str

    def visitCondrec(self, ctx):
        recType = self.typemap[ctx.variable.text]
        typestr = self.getTypeString(recType)
        rec_str = "{}, {} = condreceive({}, {});\n".format(ctx.signal.text,
                                                           ctx.variable.text,
                                                           ctx.sender.getText(), typestr)
        return rec_str

    def visitProbc(self, ctx):
        rec_str = "{} = {} [{}] {};\n".format(ctx.variable.text,
                                              ctx.p1.getText(), ctx.p3.getText(),
                                              ctx.p2.getText())
        return rec_str

    def visitForStmt(self, ctx):
        # If range is used we assume it is over a set of threads
        # We can do more involved checks, but not needed for our benchmarks
        if ctx.rangeClause():
            # Do we need to add a check for nested loops here? Will lead to problems
            threadGroup = ctx.rangeClause().expression().getText()
            enumeratorVar = ctx.rangeClause().identifierList().IDENTIFIER(0).getText()
            assignedVar = ctx.rangeClause().identifierList().IDENTIFIER(1).getText()

            transltatedBlock = self.translateBlock(ctx.block())
            if not (None in transltatedBlock):
                translatedString = ''.join(transltatedBlock)
                if enumeratorVar != "_":
                    forloop_template_str = "{3} = 0;\nfor {0} in {1} do {{\n{2}{3}={3}+1;\n}};\n"
                    return forloop_template_str.format(assignedVar, threadGroup,
                                                       translatedString, enumeratorVar)
                else:
                    forloop_template_str = "for {0} in {1} do {{\n{2}}};\n"
                    return forloop_template_str.format(assignedVar, threadGroup, translatedString)
            else:
                return ""
        # Other loops require bounds on the number of iterations to perform reliability/accuracy analysis
        # Assumption: Loop goes from 0-N.
        # Assume that the developer annotates the loop with the maximum number of iterations
        else:
            try:
                maxiterations = int(ctx.MAXITER().getText().split('=')[-1][:-2])
            except Exception as e:
                EXITWITHERROR("Loops require bounds on maxiterations: \n{}\n{}".format(ctx.getText(), str(e)))
            loopvar = ctx.forClause().inc.incDecStmt().expression().getText()
            transltatedBlock = self.translateBlock(ctx.block())
            translatedString = ''.join(transltatedBlock)
            repeat_tmpl = "{2}=0;\nrepeat {0} {{\n{1}\n{2}={2}+1;\n}};\n"
            return repeat_tmpl.format(maxiterations, translatedString, loopvar)

    def visitSimpleif(self, ctx):
        # "IF var THEN '{' (ifs+=statement ';')+ '}' ELSE '{' (elses+=statement ';')+ '}'"
        try:
            if int(ctx.cond.getText().split('!=')[-1]) == 0:
                if ctx.block(1):
                    transltatedIfBlock = ''.join(self.translateBlock(ctx.block(0)))
                    transltatedElseBlock = ''.join(self.translateBlock(ctx.block(1)))
                    return "if {} then {{\n {} }}\n else {{\n {} }};\n".format(
                        ctx.cond.getText().split('!=')[0],
                        transltatedIfBlock,
                        transltatedElseBlock)
                else:
                    transltatedIfBlock = ''.join(self.translateBlock(ctx.block(0)))
                    return "if {} then {{\n {} }};\n".format(ctx.cond.getText().split('!=')[0],
                                                             transltatedIfBlock)
            else:
                EXITWITHERROR("Only x != 0 type conditions supported: " + ctx.getText())
        except Exception as e:
            print(e)
            EXITWITHERROR("Only x != 0 type conditions supported: " + ctx.getText())

    def translateBlock(self, block):
        translatedStatements = []
        for statement in block.statementList().statement():
            if isinstance(statement, GoParser.StmtdecContext):
                translatedStatements.append(self.translateDeclaration(statement))
            elif isinstance(statement, GoParser.StmtsimpleContext):
                translation = self.visit(statement)
                if translation:
                    translatedStatements.append(translation)
                else:
                    print("[WARNING] Unable to translate simpleStmt: " + statement.getText())
            elif isinstance(statement, GoParser.StmtforContext):
                translatedStatements.append(self.visit(statement))
            elif isinstance(statement, GoParser.StmtifContext):
                translatedStatements.append(self.visit(statement))
            # Parallely does not support defer statements. THis is used to wait for threads to finish.
            elif isinstance(statement, GoParser.StmtdeferContext):
                continue
            else:
                    print("[WARNING] Unable to translate: " + statement.getText())
        return translatedStatements

    def translateDeclaration(self, dec):
        if dec.declaration().varDecl() or dec.declaration().approxVarDecl():
            isapprox = (dec.declaration().approxVarDecl() is not None)
            if isapprox:
                typequal = "approx"
                varspec = dec.declaration().approxVarDecl().varSpec()
            else:
                typequal = "precise"
                varspec = dec.declaration().varDecl().varSpec()

            if (len(varspec) > 1):
                EXITWITHERROR("Only 1 declaration per line allowed")

            varname = varspec[0].identifierList().getText()
            basetype = varspec[0].type_()

            # checks if arraytype. This works because we dont support pointer, etc. here
            # We should change parally syntax to match golang more
            if basetype.typeLit():
                arraysize = basetype.typeLit().arrayType().arrayLength().getText()
                arraytype = basetype.typeLit().arrayType().elementType().getText()
                self.typemap[varname] = (arraytype, 1, isapprox)
                return "{} {}[{}] {};\n".format(typequal, arraytype, arraysize, varname)
            else:
                vartype = varspec[0].type_().typeName().getText()
                self.typemap[varname] = (vartype, 0, isapprox)
                return "{} {} {};\n".format(typequal, vartype, varname)


def getThreadSet(func):
    func_names = []
    thread_names = []
    is_group = []
    for stat in func.block().statementList().statement():
        if isinstance(stat, GoParser.StmtfunctionContext):
            if isFunction(stat.IDENTIFIER(), "parallely.LaunchThread"):
                func_names.append(stat.arguments().expressionList().expression()[1].getText())
                thread_names.append(stat.arguments().expressionList().expression()[0].getText())
                is_group.append(0)
            if isFunction(stat.IDENTIFIER(), "parallely.LaunchThreadGroup"):
                func_names.append(stat.arguments().expressionList().expression()[1].getText())
                thread_names.append((stat.arguments().expressionList().expression()[0].getText(),
                                     stat.arguments().expressionList().expression()[2].getText()))
                is_group.append(1)
    return (func_names, thread_names, is_group)


def main(args):
    programfile = open(args.programfile, 'r')
    program_str = programfile.read()

    input_stream = InputStream(program_str)
    lexer = GoLexer(input_stream)
    stream = CommonTokenStream(lexer)
    parser = GoParser(stream)
    # parser._listeners = [MyErrorListener()]
    parser._interp.predictionMode = PredictionMode.SLL

    tree = parser.sourceFile()

    globaldecstrings = []
    for globaldec in tree.declaration():
        globaldecstrings.append(translateDeclaration(globaldec))

    functions = {}
    spawned = []
    for func in tree.functionDecl():
        functions[func.IDENTIFIER().getText()] = func.block()
        if func.IDENTIFIER().getText() == "main":
            spawned = getThreadSet(func)
    print("Main thread launches: " + str(spawned))

    translatedFuncs = {}
    for func in tree.functionDecl():
        translator = FunctionTranslator()
        if func.IDENTIFIER().getText() in spawned[0]:
            print("------------------------------------------------------------")
            print("Translating Function: " + func.IDENTIFIER().getText())
            statementStr = translator.translateBlock(func.block())
            t_str = ''.join(statementStr)
            translatedFuncs[func.IDENTIFIER().getText()] = t_str

    translatedThreads = []
    for i, tid in enumerate(spawned[1]):
        if not spawned[2][i]:
            translatedThreads.append("{}:[\nprecise int _temp;\n{}]".format(
                tid,
                translatedFuncs[spawned[0][i]]))
        else:
            translatedThreads.append("{} in {}:[\nprecise int _temp;\nprecise int {};\n{}]".format(
                tid[1][1:-1], tid[0], tid[1][1:-1],
                translatedFuncs[spawned[0][i]]))

    outfile = open(args.outfile, 'w')
    outfile.write(''.join(globaldecstrings) + "\n||\n".join(translatedThreads))


if __name__ == '__main__':
    parser = ArgumentParser()
    parser.add_argument("-f", dest="programfile",
                        help="File containing the code", required=True)
    parser.add_argument("-o", dest="outfile",
                        help="Output File", required=False)
    args = parser.parse_args()
    main(args)
