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
            if type1[0] == 'approx' or type2[0] == 'approx':
                return ('approx', type1[1])
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
    def visitSkipstatement(self, ctx):
        return True

    def visitSeqcomposition(self, ctx):
        # print ctx.getText()
        type1 = self.visit(ctx.getChild(0))
        type2 = self.visit(ctx.getChild(2))
        return (type1 and type2)

    def visitBlock(self, ctx):
        return self.visit(ctx.getChild(1))

    def visitExpassignment(self, ctx):
        # print ctx.getText()
        var_type = self.typecontext[ctx.VAR().getText()]
        expr_type = self.visit(ctx.expression())
        if (var_type == expr_type):
            return True
        if (var_type[1] == expr_type[1]) and (var_type[0] == 'approx'):
            return True
        else:
            print "Type Error : {}, {}, {}".format(ctx.getText(),
                                                   var_type, expr_type)
            return False

    def visitBoolassignment(self, ctx):
        var_type = self.typecontext[ctx.VAR().getText()]
        expr_type = self.visit(ctx.expression())
        if (var_type == expr_type):
            return True
        if (var_type[1] == expr_type[1]) and (var_type[0] == 'approx'):
            return True
        else:
            print "Type Error : {}, {}, {}".format(ctx.getText(),
                                                   var_type, expr_type)
            return False

    def visitIf(self, ctx):
        guardtype = self.visit(ctx.getChild(1))
        if guardtype != ('precise', 'bool'):
            print "Type Error precise boolean expected. ", ctx.getText()
            return False
        then_type = self.visit(ctx.getChild(3))
        else_type = self.visit(ctx.getChild(5))
        return (then_type and else_type)

    def visitSend(self, ctx):
        # At some point check if the first element is a pid
        var_type = self.typecontext[ctx.getChild(6).getText()]
        sent_qual = ctx.getChild(4).getChild(0).getText()
        sent_type = ctx.getChild(4).getChild(1).getText()
        if var_type == (sent_qual, sent_type):
            return True
        else:
            print "Type Error : {}".format(ctx.getText())
            return False

    def visitReceive(self, ctx):
        # At some point check if the first element is a pid
        var_type = self.typecontext[ctx.getChild(0).getText()]
        rec_qual = ctx.getChild(6).getChild(0).getText()
        rec_type = ctx.getChild(6).getChild(1).getText()
        if var_type == (rec_qual, rec_type):
            return True
        else:
            print "Type Error : {}".format(ctx.getText())
            return False

    def visitSingleprogram(self, ctx):
        return self.visit(ctx.statement())

    def visitParcomposition(self, ctx):
        print ctx.getChild(0).getText(), ctx.getChild(2).getText()
        type1 = self.visit(ctx.getChild(0))
        type2 = self.visit(ctx.getChild(2))
        return (type1 and type2)

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

    # parser.buildParseTrees = True

    tree = parser.program()

    visitor = parallelyTypeChecker()
    visitor.visit(tree)
    # print(Trees.toStringTree(tree, None, parser))


if __name__ == '__main__':
    programfile = open(sys.argv[1])
    program_str = programfile.readline()
    main(program_str)
