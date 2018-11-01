import sys
from antlr4 import CommonTokenStream
from antlr4 import InputStream
from antlr4 import *
from ParallelyLexer import ParallelyLexer
from ParallelyParser import ParallelyParser
from antlr4.tree.Trees import Trees
from ParallelyVisitor import ParallelyVisitor

key_error_msg = "Type error detected: Undeclared variable (probably : {})"

class parallelyTypeChecker(ParallelyVisitor):
    def __init__(self):
        self.typecontext = {}

    def baseTypesEqual(self, type1, type2, ctx):
        if not (type1[1] == type2[1]):
            print "Type error : ", ctx.getText(), type1, type2
            exit(-1)
        else:
            return type1

    ########################################
    # Expression type checking
    ########################################
    def visitLiteral(self, ctx):
        # return (ParallelyLexer.PRECISETYPE, ParallelyLexer.INTTYPE)
        return ("precise", "int")

    def visitVariable(self, ctx):
        return self.typecontext[ctx.getText()]

    def visitMultiply(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.baseTypesEqual(type1, type2, ctx)

    def visitAdd(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.baseTypesEqual(type1, type2, ctx)

    def visitMinus(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.baseTypesEqual(type1, type2, ctx)

    def visitDivide(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.baseTypesEqual(type1, type2, ctx)

    ########################################
    # Boolean expressions type checking
    ########################################
    def visitTrue(self, ctx):
        return ("precise", "bool")

    def visitFalse(self, ctx):
        return ("precise", "bool")

    def visitBoolvariable(self, ctx):
        return self.typecontext[ctx.getText()]

    def visitEqual(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.baseTypesEqual(type1, type2, ctx)

    def visitGreater(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.baseTypesEqual(type1, type2, ctx)

    def visitLess(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.baseTypesEqual(type1, type2, ctx)

    def visitAnd(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.baseTypesEqual(type1, type2, ctx)

    def visitOr(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.baseTypesEqual(type1, type2, ctx)

    def visitNot(self, ctx):
        type1 = self.visit(ctx.expression(0))
        if type1[1] == 'bool':
            return type1

    ########################################
    # Declaration type checking
    ########################################
    def visitSingledeclaration(self, ctx):
        decl_type = (ctx.fulltype().typequantifier().getText(),
                     ctx.fulltype().getChild(1).getText())
        self.typecontext[ctx.VAR().getText()] = decl_type

    def visitMultipledeclaration(self, ctx):
        self.visit(ctx.getChild(0))
        self.visit(ctx.getChild(2))

    ########################################
    # Statement type checking
    ########################################
    
    def visitSingleprogram(self, ctx):
        return self.visit(ctx.statement())

    def visitExpassignment(self, ctx):
        var_type = self.typecontext[ctx.VAR().getText()]
        expr_type = self.visit(ctx.expression())
        return var_type == expr_type

    def visitProgram(self, ctx):
        print ctx.getText()
        # Read the declarations and build up the type table
        self.visit(ctx.declaration())
        print self.typecontext

        try:
            typechecked = self.visit(ctx.parallelprogram())
            print typechecked
        except KeyError, keyerror:
            print key_error_msg.format(keyerror)


def main(program_str):
    input_stream = InputStream(program_str)
    lexer = ParallelyLexer(input_stream)
    stream = CommonTokenStream(lexer)
    parser = ParallelyParser(stream)
    tree = parser.program()

    visitor = parallelyTypeChecker()
    visitor.visit(tree)
    # print(Trees.toStringTree(tree, None, parser))


if __name__ == '__main__':
    programfile = open(sys.argv[1])
    program_str = programfile.readline()
    main(program_str)
