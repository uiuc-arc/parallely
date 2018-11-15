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

key_error_msg = "Type error detected: Undeclared variable (probably : {})"


class parallelyTypeChecker(ParallelyVisitor):
    def __init__(self):
        self.typecontext = {}
        self.processgroups = {}

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

    def visitArraydec(self, ctx):
        decl_type = (ctx.fulltype().typequantifier().getText(),
                     ctx.fulltype().getChild(1).getText())
        self.typecontext[ctx.var().getText()] = decl_type

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

    def visitForloop(self, ctx):
        type_checked = self.visit(ctx.statement())
        return type_checked

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

    def visitGroupedprogram(self, ctx):
        self.typecontext = {}
        self.visit(ctx.declaration())
        try:
            typechecked = self.visit(ctx.statement())
        except KeyError, keyerror:
            print key_error_msg.format(keyerror)
            typechecked = False
        if not typechecked:
            print "Process {} failed typechecker".format(ctx.processset().getText())
        self.typecontext = {}
        return typechecked

    def visitParcomposition(self, ctx):
        # print ctx.getChild(0).getText(), ctx.getChild(2).getText()
        type1 = self.visit(ctx.getChild(0))
        type2 = self.visit(ctx.getChild(2))
        return type1 and type2

    def visitSingle(self, ctx):
        # Read the declarations and build up the type table
        self.visit(ctx.globaldec())
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
        self.globaldecs = {}
        self.grouped_list = {}

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
        return self.statement_lists

    def visitGroupedprogram(self, ctx):
        pid = ctx.VAR()
        pgroup = ctx.GLOBALVAR().getText()
        decs = self.flattenStatement(ctx.declaration())
        statements = self.flattenStatement(ctx.statement())
        print [x.getText() for x in statements]
        # self.statement_lists[pid.getText()] = decs + statements
        self.grouped_list[pgroup] = (pid, decs, statements)
        return self.grouped_list

    def visitParcomposition(self, ctx):
        self.visit(ctx.getChild(0))
        self.visit(ctx.getChild(2))
        return self.statement_lists

    def visitMultipledeclaration(self, ctx):
        self.visit(ctx.getChild(0))
        self.visit(ctx.getChild(2))
        return self.statement_lists

    def visitSingleglobaldec(self, ctx):
        ind = ctx.GLOBALVAR().getText()
        members = ctx.processid()
        self.globaldecs[ind] = members

    def visitSingle(self, ctx):
        self.visit(ctx.globaldec())
        self.visit(ctx.parallelprogram())
        return self.statement_lists

    def rewriteStatement(self, pid, statement, statement_list, msgcontext):
        if isinstance(statement, ParallelyParser.DeclarationContext):
            dec_type_q = statement.fulltype().typequantifier().getText()
            dec_type_t = statement.fulltype().getChild(1).getText()
            dec_name = statement.var().getText()
            newdec = " ".join([dec_type_q, dec_type_t, dec_name])
            return True, newdec, msgcontext
        if isinstance(statement, ParallelyParser.SendContext):
            rec = statement.processid().getText()
            sent_type_q = statement.fulltype().typequantifier().getText()
            sent_type_t = statement.fulltype().getChild(1).getText()
            sent_var = statement.var().getText()
            my_key = (rec, pid, sent_type_q, sent_type_t)
            if my_key in msgcontext.keys():
                msgcontext[my_key].append((sent_var,))
            else:
                msgcontext[my_key] = [(sent_var,)]
            return True, '', msgcontext
        if isinstance(statement, ParallelyParser.CondsendContext):
            rec = statement.processid().getText()
            sent_type_q = statement.fulltype().typequantifier().getText()
            sent_type_t = statement.fulltype().getChild(1).getText()
            guard_var = statement.var()[0].getText()
            sent_var = statement.var()[1].getText()
            my_key = (rec, pid, sent_type_q, sent_type_t)
            if my_key in msgcontext.keys():
                msgcontext[my_key].append((sent_var, guard_var))
            else:
                msgcontext[my_key] = [(sent_var, guard_var)]
            return True, '', msgcontext
        if isinstance(statement, ParallelyParser.CondreceiveContext):
            guard_var = statement.var()[0].getText()
            assigned_var = statement.var()[1].getText()
            sender = statement.processid().getText()
            sent_type_q = statement.fulltype().typequantifier().getText()
            sent_type_t = statement.fulltype().getChild(1).getText()
            # Dont have to do this?
            assign_symbol = statement.getChild(1).getText()
            my_key = (pid, sender, sent_type_q, sent_type_t)
            # print my_key
            if my_key in msgcontext.keys():
                if len(msgcontext[my_key]) > 0:
                    # If the top is not a guarded expression
                    if len(msgcontext[my_key][0]) != 2:
                        return False, '', msgcontext
                    rec_val, rec_guard = msgcontext[my_key].pop(0)
                    # Working with strings feel weird. Fix Later
                    out_format = "{} = 1 [{}] 0;\n{}={} [{}] {}"
                    rewrite = out_format.format(guard_var, rec_guard,
                                                assigned_var, rec_val,
                                                rec_guard,
                                                assigned_var)
                    return True, rewrite, msgcontext
                else:
                    return False, '', msgcontext
            else:
                return False, '', msgcontext
        if isinstance(statement, ParallelyParser.ReceiveContext):
            assigned_var = statement.var().getText()
            sender = statement.processid().getText()
            sent_type_q = statement.fulltype().typequantifier().getText()
            sent_type_t = statement.fulltype().getChild(1).getText()
            # Dont have to do this?
            assign_symbol = statement.getChild(1).getText()
            my_key = (pid, sender, sent_type_q, sent_type_t)
            # print my_key
            if my_key in msgcontext.keys():
                if len(msgcontext[my_key]) > 0:
                    # If the top is a guarded expression exit
                    if len(msgcontext[my_key][0]) != 1:
                        return False, '', msgcontext
                    rec_val = msgcontext[my_key].pop(0)[0]
                    # Working with strings feel weird. Fix Later
                    rewrite = "{}{}{}".format(assigned_var,
                                              assign_symbol,
                                              rec_val)
                    return True, rewrite, msgcontext
                else:
                    return False, '', msgcontext
            else:
                return False, '', msgcontext
        if isinstance(statement, ParallelyParser.ForloopContext):
            # Do the renaming step later.
            # For now assuming that the variable groups have the same iterator
            out_template = "for {} in {} do {{\n{}\n}}"

            my_statements = self.flattenStatement(statement.statement())
            target_group = statement.GLOBALVAR().getText()
            group_statements = self.grouped_list[target_group]
            group_var = group_statements[0].getText()

            tmp_statements = {}
            tmp_statements[pid] = my_statements
            tmp_statements[group_var] = group_statements[2]
            output = self.rewritePair(pid, group_var, tmp_statements,
                                      copy.deepcopy(self.msgcontext))
            success, result, remaining = output
            # print success, result, remaining
            if success:
                if remaining:
                    remaing_group = (group_statements[0], group_statements[1],
                                     remaining)
                    self.grouped_list[target_group] = remaing_group
                else:
                    self.grouped_list.pop(group_var, None)
                final_res = out_template.format(group_var,
                                                target_group, result)
                return True, final_res, msgcontext
            else:
                return False, '', msgcontext
        if isinstance(statement, ParallelyParser.IfContext):
            out_template = "if {} then {{{}}} else {{{}}}"
            bool_var = statement.var().getText()
            if_state = statement.statement(0).getText()
            then_state = statement.statement(1).getText()
            result = out_template.format(bool_var, if_state, then_state)
            return True, result, msgcontext
        else:
            result = statement.getText()
            return True, result, msgcontext

    def rewritePair(self, pid1, pid2, statement, msgcontext):
        temp_statement = {}
        temp_statement[pid1] = statement[pid1]
        for i in range(len(statement[pid2])):
            temp_statement[pid2] = statement[pid2][:i + 1]
            success = self.doRewriteProgram(pid1, temp_statement, msgcontext)
            if success[0]:
                return True, success[1], statement[pid2][i + 1:]
        return False, '', statement[pid2]

    def doRewriteProgram(self, pidin, statements, msgcontext):
        statements = copy.deepcopy(statements)
        my_msgcontext = copy.deepcopy(msgcontext)
        rewritten_statements = []
        while(True):
            changed = False
            for pid in statements.keys():
                # If all statements from a pid is removed
                # Congruence rule
                if len(statements[pid]) == 0:
                    statements.pop(pid, None)
                    break
                statement = statements[pid][0]
                # print statement.getText()
                output = self.rewriteStatement(pid,
                                               statement,
                                               statements,
                                               my_msgcontext)
                success, result, my_msgcontext = output
                if success:
                    statements[pid].pop(0)
                    if result != '':
                        rewritten_statements.append(result)
                    changed = True
                # print success, result, my_msgcontext, statements, pid

            # If no rewrite is possible
            if not changed:
                break
        # print statements, pid, pid in statements.keys()

        keys_removed = []
        for key in my_msgcontext:
            if len(my_msgcontext[key]) == 0:
                keys_removed.append(key)
        for key in keys_removed:
            my_msgcontext.pop(key, None)

        if pidin in statements.keys() or my_msgcontext != msgcontext:
            # print "===================="
            # print "Partial rewrite failed"
            # print "Current State : ",
            # print statements, my_msgcontext
            # print "===================="
            return False, ";\n".join(rewritten_statements), statements
        else:
            # print "===================="
            # print "YAY!!!"
            # print "Current State : "
            # print statements, my_msgcontext
            # print "===================="
            return True, ";\n".join(rewritten_statements), statements

    def rewriteProgram(self, tree, outfile):
        # Build the statement lists
        statements = self.visit(tree)
        msgcontext = {}
        # print self.statement_lists
        # print [x.getText() for x in self.statement_lists['1']]

        print '----------------------------------------'
        print 'Starting the rewriting process'
        print '----------------------------------------'

        temp = []
        for key in self.globaldecs:
            array_str = ', '.join([a.getText() for a in self.globaldecs[key]])
            temp.append("{}={{{}}}".format(key, array_str))
        global_decs_str = '\n'.join(temp)

        # rewritten_string = ""
        rewritten_statements = []

        while(True):
            changed = False
            for pid in statements.keys():
                # If all statements from a pid is removed
                # Congruence rule
                if len(statements[pid]) == 0:
                    statements.pop(pid, None)
                    break
                statement = self.statement_lists[pid][0]
                # print statement.getText()
                output = self.rewriteStatement(pid,
                                               statement,
                                               statements,
                                               msgcontext)
                success, result, msgcontext = output
                if success:
                    statements[pid].pop(0)
                    if result != '':
                        rewritten_statements.append(result)
                    changed = True

            # If no rewrite is possible
            if not changed:
                break
        if statements:
            print "Rewriting failed to completely sequentialize"
            print "Current State :"
            print statements
            print msgcontext
        else:
            print "Rewriting Successful"
            rewritten_string = global_decs_str + ";\n" + ";\n".join(rewritten_statements)
            outfile.write(rewritten_string)


class VariableRenamer(ParallelyListener):
    def __init__(self, stream):
        self.current_process = None
        self.rewriter = TokenStreamRewriter.TokenStreamRewriter(stream)
        self.done = set()

    def enterSingleprogram(self, ctx):
        self.current_process = ctx.processid()

    def enterGroupedprogram(self, ctx):
        self.current_process = ctx.getChild(0)

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
            print ctx.parentCtx.getTokens(ctx)
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


class unrollRepeat(ParallelyListener):
    def __init__(self, stream):
        self.rewriter = TokenStreamRewriter.TokenStreamRewriter(stream)

    def enterRepeat(self, ctx):
        cs = ctx.statement().start.getInputStream()
        statements = cs.getText(ctx.statement().start.start,
                                ctx.statement().stop.stop)
        rep_variable = int(ctx.INT().getText())
        edited = ''
        # removing the code for process groups
        self.rewriter.delete(self.rewriter.DEFAULT_PROGRAM_NAME,
                             ctx.start.tokenIndex, ctx.stop.tokenIndex)
        for var in range(rep_variable):
            edited += statements + ";\n"
        self.rewriter.insertAfter(ctx.stop.tokenIndex, edited)


def main(program_str, outfile):
    input_stream = InputStream(program_str)
    lexer = ParallelyLexer(input_stream)
    stream = CommonTokenStream(lexer)
    parser = ParallelyParser(stream)

    # # Unroll process groups for easy analysis?
    # For now not doing this
    # Damages the readability of the code
    tree = parser.program()
    unroller = unrollRepeat(stream)
    walker = ParseTreeWalker()
    walker.walk(unroller, tree)

    # # Rename all the variables to var_pid
    input_stream = InputStream(unroller.rewriter.getDefaultText())
    print '----------------------------------------'
    print "After unrolling repeats"
    print input_stream
    print '----------------------------------------'

    lexer = ParallelyLexer(input_stream)
    stream = CommonTokenStream(lexer)
    parser = ParallelyParser(stream)
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
