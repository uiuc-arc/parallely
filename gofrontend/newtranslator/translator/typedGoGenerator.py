from ..antlrgenerated.GoParser import GoParser
from newtranslator.antlrgenerated.GoLexer import GoLexer
from ..antlrgenerated import GoParserVisitor
from ..antlrgenerated import GoParserListener
from argparse import ArgumentParser
from antlr4 import CommonTokenStream
from antlr4 import InputStream
from antlr4 import PredictionMode
from antlr4 import ParseTreeWalker
import TokenStreamRewriter

send_str = {
    ("int", 0): "parallely.SendInt({}, {}, {});\n",
    ("int32", 0): "parallely.SendInt32({}, {}, {});\n",
    ("int64", 0): "parallely.SendInt64({}, {}, {});\n",
    ("float32", 0): "parallely.SendFloat32({}, {}, {});\n",
    ("float64", 0): "parallely.SendFloat64({}, {}, {});\n",
    ("int", 1): "parallely.SendIntArray({}[:], {}, {});\n",
    ("int32", 1): "parallely.SendInt32Array({}[:], {}, {});\n",
    ("int64", 1): "parallely.SendInt64Array({}[:], {}, {});\n",
    ("float32", 1): "parallely.SendFloat32Array({}[:], {}, {});\n",
    ("float64", 1): "parallely.SendFloat64Array({}[:], {}, {});\n"
}


cond_send_str = {
    ("int", 0): "parallely.Condsend({}, {}, {}, {});\n",
    ("int32", 0): "parallely.CondsendInt32({}, {}, {}, {});\n",
    ("int64", 0): "parallely.CondsendInt64({}, {}, {}, {});\n",
    ("float32", 0): "parallely.CondsendFloat32({}, {}, {}, {});\n",
    ("float64", 0): "parallely.CondsendFloat64({}, {}, {}, {});\n",
    ("int", 1): "parallely.CondsendIntArray({}, {}[:], {}, {});\n",
    ("int", 1): "parallely.CondsendIntArray({}, {}[:], {}, {});\n",
    ("float64", 1): "parallely.CondsendFloat64Array({}, {}[:], {}, {});\n"
}

rec_str = {
    ("int", 0): "parallely.ReceiveInt(&{}, {}, {});\n",
    ("int32", 0): "parallely.ReceiveInt32(&{}, {}, {});\n",
    ("int64", 0): "parallely.ReceiveInt64(&{}, {}, {});\n",
    ("float32", 0): "parallely.ReceiveFloat32(&{}, {}, {});\n",
    ("float64", 0): "parallely.ReceiveFloat64(&{}, {}, {});\n",
    ("int", 1): "parallely.ReceiveIntArray({}[:], {}, {});\n",
    ("int32", 1): "parallely.ReceiveInt32Array({}[:], {}, {});\n",
    ("int64", 1): "parallely.ReceiveInt64Array({}[:], {}, {});\n",
    ("float32", 1): "parallely.ReceiveFloat32Array({}[:], {}, {});\n",
    ("float64", 1): "parallely.ReceiveFloat64Array({}[:], {}, {});\n"
}

cond_rec_str = {
    ("int", 0): "parallely.Condreceive(&{}, &{}, {}, {});\n",
    ("int32", 0): "parallely.CondreceiveInt32(&{}, &{}, {}, {});\n",
    ("int64", 0): "parallely.CondreceiveInt64(&{}, &{}, {}, {});\n",
    ("int", 1): "parallely.CondreceiveIntArray(&{}, {}[:], {}, {});\n",
    ("float32", 1): "parallely.CondreceiveFloat32(&{}, &{}, {}, {});\n",
    ("float64", 1): "parallely.CondreceiveFloat64Array(&{}, {}[:], {}, {});\n"
}

rand_str = {
    ("int"): "{} = parallely.Randchoice(float32({}), {}, {});\n",
    ("float64"): "{} = parallely.RandchoiceFloat64(float32({}), {}, {});\n",
}


# Helper function to print error message and exit
def EXITWITHERROR(msg):
    print("[ERROR] " + msg)
    exit(-1)


# Check if functions are the same by combining the identifiers
def isFunction(listidentifiers, fname):
    return fname == '.'.join([i.getText() for i in listidentifiers])


class BooleanTranslator(GoParserListener.GoParserListener):
    def __init__(self, rewriter):
        self.rewriter = rewriter

    def exitAssignment(self, ctx):
         if isinstance(ctx.expressionList(1).expression()[0], GoParser.BooleanexpContext):
            bexpr = ctx.expressionList(1).expression()[0]
            convert_str = "parallely.ConvBool({})".format(bexpr.getText())
            self.rewriter.delete(self.rewriter.DEFAULT_PROGRAM_NAME,
                                 bexpr.start.tokenIndex,
                                 bexpr.stop.tokenIndex)
            self.rewriter.insertAfter(bexpr.stop.tokenIndex, convert_str)


class FunctionTranslator(GoParserVisitor.GoParserVisitor):
    def __init__(self, stream, tid, rewriter):
        self.statements = []
        self.typemap = {}
        self.tid = tid
        self.newDeclarations = []
        self.rewriter = rewriter

    def replaceString(self, ctx, new_str):
        self.rewriter.delete(self.rewriter.DEFAULT_PROGRAM_NAME,
                             ctx.start.tokenIndex,
                             ctx.stop.tokenIndex)
        self.rewriter.insertAfter(ctx.stop.tokenIndex, new_str)

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
        new_send_str = send_str[(sentType[0], sentType[1])].format(ctx.variable.text, self.tid, ctx.rec.getText())
        self.replaceString(ctx, new_send_str)

    def visitCondsend(self, ctx):
        sentType = self.typemap[ctx.variable.text]
        new_send_str = cond_send_str[(sentType[0], sentType[1])].format(ctx.cond.text,
                                                                        ctx.variable.text,
                                                                        self.tid, ctx.rec.getText())
        self.replaceString(ctx, new_send_str)

    def visitRec(self, ctx):
        recType = self.typemap[ctx.variable.text]
        new_rec_str = rec_str[(recType[0], recType[1])].format(ctx.variable.text,
                                                               self.tid, ctx.sender.getText())
        self.replaceString(ctx, new_rec_str)

    def visitNoisyrec(self, ctx):
        recType = self.typemap[ctx.variable.text]
        new_rec_str = rec_str[(recType[0], recType[1])].format(ctx.variable.text, self.tid, ctx.sender.getText())
        self.replaceString(ctx, new_rec_str)

    def visitCondrec(self, ctx):
        recType = self.typemap[ctx.variable.text]
        new_cond_rec_str = cond_rec_str[(recType[0], recType[1])].format(ctx.signal.text, ctx.variable.text,
                                                                         self.tid, ctx.sender.getText())
        self.replaceString(ctx, new_cond_rec_str)

    def visitProbc(self, ctx):
        new_rand_str = rand_str[self.typemap[ctx.variable.text][0]].format(ctx.variable.text,
                                                                           ctx.p3.getText(),
                                                                           ctx.p1.getText(),
                                                                           ctx.p2.getText())
        self.replaceString(ctx, new_rand_str)

    def visitForStmt(self, ctx):
        self.translateBlock(ctx.block())

    def visitSimpleif(self, ctx):
        self.translateBlock(ctx.block(0))
        if ctx.block(1):
            self.translateBlock(ctx.block(1))

    def translateBlock(self, block):
        for statement in block.statementList().statement():
            if isinstance(statement, GoParser.StmtdecContext):
                self.translateDeclaration(statement)
            elif isinstance(statement, GoParser.StmtsimpleContext):
                self.visit(statement)
            elif isinstance(statement, GoParser.StmtforContext):
                self.visit(statement)
            elif isinstance(statement, GoParser.StmtifContext):
                self.visit(statement)

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
                arraytype = basetype.typeLit().arrayType().elementType().getText()
                self.typemap[varname] = (arraytype, 1, isapprox)
            else:
                vartype = varspec[0].type_().typeName().getText()
                self.typemap[varname] = (vartype, 0, isapprox)


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
                                     stat.arguments().expressionList().expression()[2].getText()[1:-1]))
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

    functions = {}
    spawned = []
    for func in tree.functionDecl():
        functions[func.IDENTIFIER().getText()] = func.block()
        if func.IDENTIFIER().getText() == "main":
            spawned = getThreadSet(func)
    rewriter = TokenStreamRewriter.TokenStreamRewriter(stream)

    for func in tree.functionDecl():
        if func.IDENTIFIER().getText() in spawned[0]:
            findex = spawned[0].index(func.IDENTIFIER().getText())
            if not spawned[2][findex]:
                tid = spawned[1][findex]
            else:
                tid = spawned[1][findex][1]
            translator = FunctionTranslator(stream, tid, rewriter)
            translator.translateBlock(func.block())

    btrans = BooleanTranslator(rewriter)
    walker = ParseTreeWalker()
    walker.walk(btrans, tree)

    outfile = open(args.outfile, 'w')
    outfile.write(rewriter.getDefaultText())


if __name__ == '__main__':
    parser = ArgumentParser()
    parser.add_argument("-f", dest="programfile",
                        help="File containing the code", required=True)
    parser.add_argument("-o", dest="outfile",
                        help="Output File", required=False)
    args = parser.parse_args()
    main(args)
