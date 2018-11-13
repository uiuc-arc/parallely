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

    def boolBaseTypesEqual(self, type1, type2, ctx):
        if not (type1[1] == type2[1]):
            print "Type error : ", ctx.getText(), type1, type2
            exit(-1)
        else:
            if type1[0] == 'approx' or type2[0] == 'approx':
                return ('approx', 'bool')
            return ('precise', 'bool')

    ########################################
    # Expression type checking
    ########################################
    def visitLiteral(self, ctx):
        # return (ParallelyLexer.PRECISETYPE, ParallelyLexer.INTTYPE)
        return ("precise", "int")

    def visitVariable(self, ctx):
        return self.typecontext[ctx.getText()]

    def visitVar(self, ctx):
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

    def visitProb(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        if not (type1[1] == type2[1]):
            print "Type error : ", ctx.getText(), type1, type2
            exit(-1)
        else:
            return ('approx', type1[1])

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
        return self.boolBaseTypesEqual(type1, type2, ctx)

    def visitGreater(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.boolBaseTypesEqual(type1, type2, ctx)

    def visitLess(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.boolBaseTypesEqual(type1, type2, ctx)

    def visitAnd(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.boolBaseTypesEqual(type1, type2, ctx)

    def visitOr(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.boolBaseTypesEqual(type1, type2, ctx)

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
        self.typecontext[ctx.var().getText()] = decl_type

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

    # Removed blocks from the grammer. Still keeping this here
    def visitBlock(self, ctx):
        return self.visit(ctx.getChild(1))

    def visitExpassignment(self, ctx):
        # print ctx.getText()
        var_type = self.typecontext[ctx.var().getText()]
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
        var_type = self.typecontext[ctx.var().getText()]
        expr_type = self.visit(ctx.boolexpression())
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
        if guardtype[0] != ('precise'):
            print "Type Error precise boolean expected. ", ctx.getText()
            return False
        then_type = self.visit(ctx.statement(0))
        else_type = self.visit(ctx.statement(1))
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

    def visitCondsend(self, ctx):
        variables = ctx.var()
        guard = self.typecontext[variables[0].getText()]
        var_type = self.typecontext[variables[1].getText()]

        if guard[0] != 'approx':
            err = "Type Error : {} has to be approx".format(
                variables[0].getText())
            print err
            return False

        if var_type[0] != 'approx':
            err = "Type Error : {} has to be approx".format(
                variables[1].getText())
            print err
            return False

        sent_qual = ctx.fulltype().getChild(0).getText()
        sent_type = ctx.fulltype().getChild(1).getText()
        if var_type == (sent_qual, sent_type):
            return True
        else:
            print "Type Error : {}".format(ctx.getText())
            return False

    def visitCondreceive(self, ctx):
        # At some point check if the first element is a pid
        variables = ctx.var()
        signal = self.typecontext[variables[0].getText()]
        var_type = self.typecontext[variables[1].getText()]

        if signal[0] != 'approx':
            err = "Type Error : {} has to be approx".format(
                variables[0].getText())
            print err
            return False

        if var_type[0] != 'approx':
            err = "Type Error : {} has to be approx".format(
                variables[1].getText())
            print err
            return False

        rec_qual = ctx.fulltype().getChild(0).getText()
        rec_type = ctx.fulltype().getChild(1).getText()
        if var_type == (rec_qual, rec_type):
            return True
        else:
            print "Type Error : {}".format(ctx.getText())
            return False

    def visitSingleprogram(self, ctx):
        self.typecontext = {}
        self.visit(ctx.declaration())
        try:
            typechecked = self.visit(ctx.statement())
        except KeyError, keyerror:
            print key_error_msg.format(keyerror)
            typechecked = False

        if not typechecked:
            print "Process {} failed typechecker".format(ctx.processid().getText())
        self.typecontext = {}
        return typechecked

    def visitParcomposition(self, ctx):
        print "Parallel Program : ", ctx.getText()
        # print ctx.getChild(0).getText(), ctx.getChild(2).getText()
        type1 = self.visit(ctx.getChild(0))
        type2 = self.visit(ctx.getChild(2))
        return type1 and type2

    def visitSingle(self, ctx):
        print "Single Program : ", ctx.getText()

        # Read the declarations and build up the type table
        # self.visit(ctx.declaration())
        # print self.typecontext
        typechecked = self.visit(ctx.parallelprogram())
        if typechecked:
            print "Type checker passed"
        else:
            print "Type checker failed. Please check"
            exit(-1)


class parallelySequentializer(ParallelyVisitor):
    def __init__(self):
        self.statement_lists = {}
        self.msgcontext = {}

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

    def visitSingleprogram(self, ctx):
        pid = ctx.processid()
        decs = self.flattenStatement(ctx.declaration())
        statements = self.flattenStatement(ctx.statement())
        # print [x.getText() for x in statements]
        self.statement_lists[pid.getText()] = decs + statements

    def visitParcomposition(self, ctx):
        self.visit(ctx.getChild(0))
        self.visit(ctx.getChild(2))

    def visitSingle(self, ctx):
        self.visit(ctx.parallelprogram())
        return self.statement_lists

    def rewriteStatement(self, pid, statement, outfile):
        rewrite_template = "{};\n"
        if isinstance(statement, ParallelyParser.DeclarationContext):
            dec_type_q = statement.fulltype().typequantifier().getText()
            dec_type_t = statement.fulltype().getChild(1).getText()
            dec_name = statement.var().getText()
            newdec = " ".join([dec_type_q, dec_type_t, dec_name]) + ";\n"
            outfile.write(newdec)
            return True
        if isinstance(statement, ParallelyParser.SendContext):
            rec = statement.processid().getText()
            sent_type_q = statement.fulltype().typequantifier().getText()
            sent_type_t = statement.fulltype().getChild(1).getText()
            sent_var = statement.var().getText()
            my_key = (rec, pid, sent_type_q, sent_type_t)
            if my_key in self.msgcontext.keys():
                self.msgcontext[my_key].append(sent_var)
            else:
                self.msgcontext[my_key] = [sent_var]
            return True
        if isinstance(statement, ParallelyParser.ReceiveContext):
            assigned_var = statement.var().getText()
            sender = statement.processid().getText()
            sent_type_q = statement.fulltype().typequantifier().getText()
            sent_type_t = statement.fulltype().getChild(1).getText()
            # Dont have to do this?
            assign_symbol = statement.getChild(1).getText()
            my_key = (pid, sender, sent_type_q, sent_type_t)
            # print my_key
            if my_key in self.msgcontext.keys():
                if len(self.msgcontext[my_key]) > 0:
                    rec_val = self.msgcontext[my_key].pop(0)
                    # Working with strings feel weird. Fix Later
                    rewrite = "{}{}{}".format(assigned_var,
                                              assign_symbol,
                                              rec_val)
                    outfile.write(rewrite_template.format(rewrite))
                    return True
                else:
                    return False
            else:
                return False
        else:
            outfile.write(rewrite_template.format(statement.getText()))
            return True

    def rewriteProgram(self, tree, outfile):
        # Build the statement lists
        self.visit(tree)
        # print self.statement_lists
        # print [x.getText() for x in self.statement_lists['1']]

        print '----------------------------------------'

        while(True):
            changed = False
            for pid in self.statement_lists.keys():
                # If all statements from a pid is removed
                # Congruence rule
                if len(self.statement_lists[pid]) == 0:
                    self.statement_lists.pop(pid, None)
                    break
                first_statement = self.statement_lists[pid][0]
                success = self.rewriteStatement(pid, first_statement, outfile)
                if success:
                    self.statement_lists[pid].pop(0)
                    changed = True

            # If no rewrite is possible
            if not changed:
                break
        if self.statement_lists:
            print "Rewriting failed to completely sequentialize"
            print "Current State :"
            print self.statement_lists
            print self.msgcontext
        else:
            print "Rewriting Successful"


class VariableRenamer(ParallelyListener):
    def __init__(self, stream):
        self.current_process = None
        self.rewriter = TokenStreamRewriter.TokenStreamRewriter(stream)
        self.done = set()

    def enterSingleprogram(self, ctx):
        print ctx.getText()
        self.current_process = ctx.processid()

    # def enterVariable(self, ctx):
    #     new_name = "_" + self.current_process.getText()
    #     # self.rewriter.insertBeforeIndex(ctx.start.tokenIndex, new_name)
    #     self.rewriter.insertAfterToken(ctx.stop, new_name)

    def enterVar(self, ctx):
        new_name = "_" + self.current_process.getText()
        # self.rewriter.insertBeforeIndex(ctx.start.tokenIndex, new_name)
        self.rewriter.insertAfterToken(ctx.stop, new_name)


class UnrollGroups(ParallelyListener):
    def __init__(self, stream):
        self.current_process = None
        self.rewriter = TokenStreamRewriter.TokenStreamRewriter(stream)
        self.done = set()

    def enterGroupedprogram(self, ctx):
        new_process_str = "{}:[{};{}]"
        processes = ctx.processset().processid()
        print [p.getText() for p in processes]

        # removing the code for process groups
        self.rewriter.delete(self.rewriter.DEFAULT_PROGRAM_NAME,
                             ctx.start.tokenIndex, ctx.stop.tokenIndex)

        # Adding the unrolled text for the new procs
        new_procs = []
        for process in processes:
            processCode = new_process_str.format(process.getText(),
                                                 ctx.declaration().getText(),
                                                 ctx.statement().getText())
            print ctx.parentCtx.getTokens(ctx) # .getTokens(ctx.start, ctx.stop)
            new_procs.append(processCode)
        edited = "||".join(new_procs)
        self.rewriter.insertAfter(ctx.stop.tokenIndex, edited)

    # def enterVariable(self, ctx):
    #     new_name = "_" + self.current_process.getText()
    #     # self.rewriter.insertBeforeIndex(ctx.start.tokenIndex, new_name)
    #     self.rewriter.insertAfterToken(ctx.stop, new_name)

    # def enterVar(self, ctx):
    #     new_name = "_" + self.current_process.getText()
    #     # self.rewriter.insertBeforeIndex(ctx.start.tokenIndex, new_name)
    #     self.rewriter.insertAfterToken(ctx.stop, new_name)


def main(program_str, outfile):
    input_stream = InputStream(program_str)
    lexer = ParallelyLexer(input_stream)
    stream = CommonTokenStream(lexer)
    parser = ParallelyParser(stream)

    # # Unroll process groups for easy analysis
    # tree = parser.program()
    # renamer = UnrollGroups(stream)
    # walker = ParseTreeWalker()
    # walker.walk(renamer, tree)

    # print stream.getText()

    # Rename all the variables to var_pid
    # input_stream = InputStream(renamer.rewriter.getDefaultText())

    # lexer = ParallelyLexer(input_stream)
    # stream = CommonTokenStream(lexer)

    tree = parser.program()
    renamer = VariableRenamer(stream)
    walker = ParseTreeWalker()
    walker.walk(renamer, tree)

    # Run type checker on the renamed version
    input_stream = InputStream(renamer.rewriter.getDefaultText())
    lexer = ParallelyLexer(input_stream)
    stream = CommonTokenStream(lexer)
    parser = ParallelyParser(stream)
    tree = parser.program()
    typechecker = parallelyTypeChecker()
    typechecker.visit(tree)

    # Sequentialization
    sequentializer = parallelySequentializer()
    sequentializer.rewriteProgram(tree, outfile)


if __name__ == '__main__':
    programfile = open(sys.argv[1], 'r')
    outfile = open(sys.argv[2], 'w')
    program_str = programfile.read()
    main(program_str, outfile)
