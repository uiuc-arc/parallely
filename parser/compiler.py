# flake8: noqa E501

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
import time
from argparse import ArgumentParser
import random

key_error_msg = "Type error detected: Undeclared variable (probably : {})"


class parallelyTypeChecker(ParallelyVisitor):
    def __init__(self, debug):
        self.typecontext = {}
        self.processgroups = {}
        self.debug = debug

    def baseTypesEqual(self, type1, type2, ctx):
        # Deadline mode. Fix!
        if not (type1[1][:3] == type2[1][:3]):
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

    def visitFliteral(self, ctx):
        return ("precise", "float")

    def visitVariable(self, ctx):
        return self.typecontext[ctx.getText()]

    def visitVar(self, ctx):
        return self.typecontext[ctx.getText()]

    def visitSelect(self, ctx):
        return self.visit(ctx.expression())

    def visitMultiply(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        res = self.baseTypesEqual(type1, type2, ctx)
        return res

    def visitAdd(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.baseTypesEqual(type1, type2, ctx)

    def visitMinus(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        res = self.baseTypesEqual(type1, type2, ctx)
        return res

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
        return self.baseTypesEqual(type1, type2, ctx)

    def visitGreater(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.baseTypesEqual(type1, type2, ctx)

    def visitLess(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.baseTypesEqual(type1, type2, ctx)

    def visitGeq(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.baseTypesEqual(type1, type2, ctx)

    def visitLeq(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.baseTypesEqual(type1, type2, ctx)

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
    def getType(self, fulltype):
        if (fulltype.typequantifier()):
            decl_type = (fulltype.typequantifier().getText(),
                         fulltype.getChild(1).getText())
        else:
            decl_type = (fulltype.fulltype().typequantifier().getText(),
                         fulltype.fulltype().getChild(1).getText() + "[]")
        return decl_type

    def visitSingledeclaration(self, ctx):
        decl_type = self.getType(ctx.fulltype())
        self.typecontext[ctx.var().getText()] = decl_type

    def visitArraydec(self, ctx):
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

    def visitArrayload(self, ctx):
        var_type = self.typecontext[ctx.var(0).getText()]
        array_type = self.typecontext[ctx.var(1).getText()]

        # print ctx.expression()
        for expr in ctx.expression():
            expr_type = self.visit(expr)
            if expr_type[0] != 'precise':
                return False

        # Deadline day
        if (var_type[1] == array_type[1][:-2]) or (var_type[0] == 'approx'):
            return True

    def visitExpassignment(self, ctx):
        # print ctx.getText()
        var_type = self.typecontext[ctx.var().getText()]
        expr_type = self.visit(ctx.expression())
        if (var_type == expr_type):
            return True
        if (var_type[1][:3] == expr_type[1][:3]) or (var_type[0] == 'approx'):
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
        sent_qual, sent_type = self.getType(ctx.getChild(4))
        # print sent_qual, sent_type
        # sent_qual = ctx.getChild(4).getChild(0).getText()
        # sent_type = ctx.getChild(4).getChild(1).getText()
        if var_type == (sent_qual, sent_type):
            return True
        else:
            print "Type Error : {}".format(ctx.getText())
            return False

    def visitReceive(self, ctx):
        # At some point check if the first element is a pid
        var_type = self.typecontext[ctx.getChild(0).getText()]
        rec_qual, rec_type = self.getType(ctx.getChild(6))
        # rec_type = ctx.getChild(6).getChild(1).getText()
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

    def visitCast(self, ctx):
        type1 = self.visit(ctx.var(0))
        type2 = ctx.fulltype().getText()
        if ''.join(list(type1)) == type2.strip():
            return True
        else:
            return False

    def visitForloop(self, ctx):
        # TODO : Check if the variable was declared in global scope
        all_typechecked = True
        for statement in ctx.statement():
            typechecked = self.visit(statement)
            if self.debug:
                print "[Debug] ", statement.getText(), typechecked
            all_typechecked = typechecked and all_typechecked
        return all_typechecked

    def visitRepeat(self, ctx):
        all_typechecked = True
        for statement in ctx.statement():
            typechecked = self.visit(statement)
            if self.debug:
                print "[Debug] ", statement.getText(), typechecked
            all_typechecked = typechecked and all_typechecked
        return all_typechecked

    def visitRepeatvar(self, ctx):
        var_type = self.typecontext[ctx.var().getText()]
        if not var_type[0] == 'precise':
            print "Type error: only precise int allowed in a repeat statement: ", ctx.getText()
            exit(-1)
        
        all_typechecked = True
        for statement in ctx.statement():
            typechecked = self.visit(statement)
            if self.debug:
                print "[Debug] ", statement.getText(), typechecked
            all_typechecked = typechecked and all_typechecked
        return all_typechecked

    def visitSingle(self, ctx):
        self.typecontext = {}
        pid = ctx.processid().getText()

        try:
            for declaration in ctx.declaration():
                self.visit(declaration)
        except Exception:
            print "Type error in declarations : ", pid
            exit(-1)

        all_typechecked = True
        try:
            for statement in ctx.statement():
                if self.debug:
                    print "[Debug - checking] ", statement.getText()
                typechecked = self.visit(statement)
                if self.debug:
                    print "[Debug - checked] ", statement.getText(), typechecked
                all_typechecked = typechecked and all_typechecked
        except KeyError, keyerror:
            print key_error_msg.format(keyerror)
            typechecked = False

        if not all_typechecked:
            print "Process {} failed typechecker".format(pid)
        self.typecontext = {}
        return all_typechecked

    def visitParcomposition(self, ctx):
        # Does nothing for now.
        # Only sets of procs allowed in these declarations
        if ctx.globaldec():
            self.visit(ctx.globaldec())

        all_type_checked = True
        for current_program in ctx.program():
            type_checked = self.visit(current_program)
            all_type_checked = type_checked and all_type_checked

        return all_type_checked


class parallelySequentializer(ParallelyVisitor):
    def __init__(self, debug):
        self.statement_lists = {}
        self.msgcontext = {}
        self.globaldecs = {}
        self.declarations = []
        self.grouped_list = {}
        self.isProcessGroup = {}
        self.debug = debug

    def isGroup(self, pid):
        if isinstance(pid, ParallelyParser.NamedpContext):
            return (False, pid.getText())
        elif isinstance(pid, ParallelyParser.VariablepContext):
            return (True, pid.VAR().getText(),)
        else:
            return (True, pid.GLOBALVAR().getText(), pid.VAR().getText())

    def visitSingle(self, ctx):
        pid = ctx.processid()
        is_group = self.isGroup(pid)
        self.isProcessGroup[is_group[1]] = is_group
        decs = ctx.declaration()
        statements = ctx.statement()
        self.declarations.extend(decs)
        print pid.getText(), is_group, len(statements)
        if not is_group[0]:
            self.statement_lists[pid.getText()] = statements
        else:
            self.statement_lists[pid.GLOBALVAR().getText()] = statements

    def visitMultipleglobaldec(self, ctx):
        self.visit(ctx.getChild(0))
        self.visit(ctx.getChild(2))

    def visitSingleglobaldec(self, ctx):
        ind = ctx.GLOBALVAR().getText()
        members = ctx.processid()
        self.globaldecs[ind] = members

    def getType(self, fulltype):
        if (fulltype.typequantifier()):
            decl_type = (fulltype.typequantifier().getText(),
                         fulltype.getChild(1).getText())
        else:
            decl_type = (fulltype.fulltype().typequantifier().getText(),
                         fulltype.fulltype().getChild(1).getText() + "[]")
        return decl_type

    def getDecString(self, dec):
        dec_type_q, dec_type_t = self.getType(dec.fulltype())
        # dec_type_t = statement.fulltype().getChild(1).getText()
        dec_name = dec.var().getText()
        newdec = "{} {} {};".format(dec_type_q, dec_type_t, dec_name)
        return newdec

    def appendIfExists(self, key, dict_in, val):
        if key in dict_in.keys():
            dict_in[key].append(val)
        else:
            dict_in[key] = [val]

    def handleSend(self, pid, statement, statement_list, msgcontext, seq_prefix):
        rec = statement.processid().getText()
        sent_type_q, sent_type_t = self.getType(statement.fulltype())
        # = statement.fulltype().getChild(1).getText()
        sent_var = statement.var().getText()
        my_key = (rec, pid, sent_type_q, sent_type_t)
        # print "1=======> ", sent_var, my_key, msgcontext
        self.appendIfExists(my_key, msgcontext, (sent_var,))
        # print "2=======> ", msgcontext
        new_statement_list = dict(statement_list)
        new_statement_list[pid].pop(0)
        return True, seq_prefix, dict(msgcontext), new_statement_list

    def handleCondSend(self, pid, statement, statement_list, msgcontext, seq_prefix):
        rec = statement.processid().getText()
        sent_type_q = statement.fulltype().typequantifier().getText()
        sent_type_t = statement.fulltype().getChild(1).getText()
        guard_var = statement.var()[0].getText()
        sent_var = statement.var()[1].getText()
        my_key = (rec, pid, sent_type_q, sent_type_t)
        self.appendIfExists(my_key, msgcontext, (sent_var, guard_var))
        statement_list[pid].pop(0)
        return True, seq_prefix, msgcontext, statement_list

    def handleReceive(self, pid, statement, statement_list, msgcontext, seq_prefix):
        assigned_var = statement.var().getText()
        sender = statement.processid().getText()

        if sender not in statement_list.keys():
            # print ("[ERROR] Receiving from an unknown sender!")
            # print ("[ERROR] ", statement.getText())
            return False, seq_prefix, msgcontext, statement_list

        sent_type_q, sent_type_t = self.getType(statement.fulltype())
        assign_symbol = statement.getChild(1).getText()
        my_key = (pid, sender, sent_type_q, sent_type_t)

        # If the msgcontext is empty or top is a guarded expression exit
        if my_key in msgcontext.keys() and len(msgcontext[my_key]) > 0:
            rec_val = msgcontext[my_key].pop(0)[0]
            rewrite = "{}{}{}".format(assigned_var, assign_symbol, rec_val)
            seq_prefix.append(rewrite)
            statement_list[pid].pop(0)
            return True, seq_prefix, msgcontext, statement_list
        else:
            return False, seq_prefix, msgcontext, statement_list

    def handleCondReceive(self, pid, statement, statement_list, msgcontext, seq_prefix):
        guard_var = statement.var()[0].getText()
        assigned_var = statement.var()[1].getText()
        sender = statement.processid().getText()
        sent_type_q = statement.fulltype().typequantifier().getText()
        sent_type_t = statement.fulltype().getChild(1).getText()

        # print "====>", self.isProcessGroup, sender
        if sender in self.isProcessGroup and self.isProcessGroup[sender][0]:
            print ("[ERROR] Receiving from a group is not supported yet!")
            # self.rewrite_statements(seq_prefix, msgcontext, remaining_statements)
            exit(-1)
        else:

            my_key = (pid, sender, sent_type_q, sent_type_t)
            if my_key in msgcontext.keys() and len(msgcontext[my_key]) > 0:
                rec_val, rec_guard = msgcontext[my_key].pop(0)
                out_format = "{} = 1 [{}] 0;\n{}={} [{}] {}"
                rewrite = out_format.format(guard_var, rec_guard,
                                            assigned_var, rec_val,
                                            rec_guard,
                                            assigned_var)
                seq_prefix.append(rewrite)
                statement_list[pid].pop(0)
                return True, seq_prefix, msgcontext, statement_list
            else:
                return False, seq_prefix, msgcontext, statement_list

    def handleIf(self, pid, statement, statement_list, msgcontext, seq_prefix):
        out_template = "if {} then {{{}}} else {{{}}}"
        bool_var = statement.var().getText()

        if_state = statement.ifs
        ifstart = if_state[0].start.getInputStream()
        ifstatements = ifstart.getText(if_state[0].start.start,
                                       if_state[-1].stop.stop) + ';\n'

        else_state = statement.elses
        elsestart = else_state[0].start.getInputStream()
        elsestatements = elsestart.getText(else_state[0].start.start,
                                           else_state[-1].stop.stop) + ';\n'

        result = out_template.format(bool_var, ifstatements, elsestatements)
        statement_list[pid].pop(0)
        seq_prefix.append(result)
        return True, seq_prefix, msgcontext, statement_list

    def isEmptyMsgContext(self, msg_context):
        for key in msg_context.keys():
            if len(msg_context[key]) > 0:
                return False
        return True

    def handleFor(self, pid, statement, statement_list, msgcontext, seq_prefix):
        # TODO: *** Do the renaming step ***
        # For now assuming that the variable groups have the same iterator
        out_template = "for {} in {} do {{\n{}\n}}"

        my_statements = statement.statement()
        target_group = statement.GLOBALVAR().getText()
        group_statements = statement_list[target_group]

        group_var = self.isProcessGroup[target_group][2]

        limit = 0

        # print "$$$$$ ", target_group

        while True:
            if limit > len(group_statements):
                if self.debug:
                    print "Giving up : ", seq_prefix, msgcontext, statement_list
                return False, seq_prefix, msgcontext, statement_list

            tmp_statements = {}
            tmp_statements[pid] = list(my_statements)
            tmp_statements[group_var] = list(group_statements[:len(group_statements) - limit])
            tmp_msgcontext = dict(msgcontext)

            if self.debug:
                print "Attempting to rewrite: ", tmp_statements, tmp_msgcontext

            output = self.rewrite_statements([], tmp_msgcontext, tmp_statements)
            if self.isEmptyMsgContext(output[1]) and (pid not in output[2]):
                break
            if self.isEmptyMsgContext(output[1]) and (pid in output[2]) and len(output[2][pid]) == 0:
                break
            limit += 1

            # print "--------------------------------------------"
            # print len(group_statements), limit, output
            # print "--------------------------------------------"

        # Entire process was rewritten
        if limit == 0:
            statement_list[pid].pop(0)
            statement_list.pop(target_group, None)

            rewrite = out_template.format(group_var, target_group, ';\n'.join(output[0]))
            seq_prefix.append(rewrite)
        # Only part of the process was rewritten
        else:
            statement_list[pid].pop(0)
            statement_list[target_group] = group_statements[len(group_statements) - limit:]

            rewrite = out_template.format(group_var, target_group, ';\n'.join(output[0]))
            seq_prefix.append(rewrite)

        return True, seq_prefix, msgcontext, statement_list

    def rewriteOneStep(self, pid, statement, statement_list, msgcontext, seq_prefix):
        if isinstance(statement, ParallelyParser.SendContext):
            return self.handleSend(pid, statement, statement_list, msgcontext, seq_prefix)
        if isinstance(statement, ParallelyParser.ReceiveContext):
            return self.handleReceive(pid, statement, statement_list, msgcontext, seq_prefix)
        if isinstance(statement, ParallelyParser.CondsendContext):
            return self.handleCondSend(pid, statement, statement_list, msgcontext, seq_prefix)
        if isinstance(statement, ParallelyParser.CondreceiveContext):
            return self.handleCondReceive(pid, statement, statement_list, msgcontext, seq_prefix)
        if isinstance(statement, ParallelyParser.ForloopContext):
            return self.handleFor(pid, statement, statement_list, msgcontext, seq_prefix)
        if isinstance(statement, ParallelyParser.IfContext):
            return self.handleIf(pid, statement, statement_list, msgcontext, seq_prefix)
        if isinstance(statement, ParallelyParser.RepeatvarContext):
            out_template = "repeat {} {{{}}}"
            bool_var = statement.VAR().getText()
            # if_state = statement.statement().getText()
            cs = statement.statement().start.getInputStream()
            statements = cs.getText(statement.statement().start.start,
                                    statement.statement().stop.stop)

            result = out_template.format(bool_var, statements)
            return True, result, msgcontext
        else:
            result = statement.getText()
            statement_list[pid].pop(0)
            seq_prefix.append(result)
            return True, seq_prefix, msgcontext, statement_list

    def rewritePair(self, pid1, pid2, statement, msgcontext):
        temp_statement = {}
        temp_statement[pid1] = statement[pid1]
        for i in range(len(statement[pid2])):
            temp_statement[pid2] = statement[pid2][:i + 1]
            success = self.doRewriteProgram(pid1, temp_statement, msgcontext)
            if success[0]:
                return True, success[1], statement[pid2][i + 1:]
        return False, '', statement[pid2]

    def isGroupedProcess(self, pid):
        # pid created in the renaming process (I assume)
        if pid not in self.isProcessGroup:
            return False
        else:
            return self.isProcessGroup[pid][0]

    def rewrite_statements(self, seq_prefix, msgcontext, remaining_statements):
        if self.debug:
            print "[Debug] ", remaining_statements, seq_prefix, msgcontext

        remaining_pids = set(remaining_statements.keys())
        while(True):
            changed = False
            group = False
            for pid in remaining_pids.copy():
                if self.isGroupedProcess(pid):
                    remaining_pids.remove(pid)
                    changed = True
                    if self.debug:
                        print "[Debug:rewrite_statements] : Dont work on groups : ",  pid
                    break
                # If all statements from a pid is removed
                if not (pid in remaining_statements.keys()):
                    if self.debug:
                        print "[Debug:rewrite_statements] : completely sequentialized : ",
                        pid, remaining_statements
                    break
                if len(remaining_statements[pid]) == 0:
                    remaining_pids.remove(pid)
                    changed = True
                    if self.debug:
                        print "[Debug:rewrite_statements] : completely sequentialized 2 : ",
                        pid, remaining_statements
                    remaining_statements.pop(pid, None)
                    continue

                statement = remaining_statements[pid][0]
                # print "[Debug] ", statement.getText(), remaining_statements, seq_prefix, msgcontext
                output = self.rewriteOneStep(pid, statement,
                                             remaining_statements,
                                             msgcontext, seq_prefix)
                success, seq_prefix, msgcontext, remaining_statements = output

                if success:
                    # statements[pid].pop(0)
                    # if result != '':
                    #     seq_prefix.append(result)
                    changed = True

            # If no rewrite is possible
            if not changed:
                break
        if self.debug:
            print "[Debug:rewrite_statements:2] ", remaining_statements, seq_prefix, msgcontext
        return seq_prefix, msgcontext, remaining_statements

    def rewriteProgram(self, tree, outfile):
        # Build the statement lists
        self.visit(tree)

        msgcontext = {}
        temp = []
        for key in self.globaldecs:
            array_str = ', '.join([a.getText() for a in self.globaldecs[key]])
            temp.append("{}={{{}}};".format(key, array_str))
        for dec in self.declarations:
            temp.append(self.getDecString(dec))
        global_decs_str = '\n'.join(temp)

        statements = self.statement_lists.copy()
        rewritten = self.rewrite_statements([], msgcontext, statements)

        if not self.isEmptyMsgContext(rewritten[1]):
            print "Sequentializion failed"
            print "Remaining Messages: ", rewritten[1]
            exit(-1)

        seq_program = global_decs_str + "\n" + ";\n".join(rewritten[0])
        outfile.write(seq_program)
        if self.debug:
            print "Sequentialized Program:"
            print seq_program


class VariableRenamer(ParallelyListener):
    def __init__(self, stream):
        self.current_process = None
        self.rewriter = TokenStreamRewriter.TokenStreamRewriter(stream)
        self.done = set()

    def enterSingle(self, ctx):
        pid = ctx.processid()
        if isinstance(pid, ParallelyParser.NamedpContext):
            self.current_process = ctx.processid()
        elif isinstance(pid, ParallelyParser.VariablepContext):
            self.current_process = ctx.processid()
        else:
            self.current_process = ctx.processid().VAR()

    def enterVar(self, ctx):
        new_name = "_" + self.current_process.getText() 
        # self.rewriter.insertBeforeIndex(ctx.start.tokenIndex, new_name)
        # self.rewriter.insertBeforeToken(ctx.start, new_name)
        self.rewriter.insertAfterToken(ctx.stop, new_name)


class UnrollGroups(ParallelyListener):
    def __init__(self, stream):
        self.current_process = None
        self.rewriter = TokenStreamRewriter.TokenStreamRewriter(stream)
        self.done = set()

    def enterGroupedprogram(self, ctx):
        new_process_str = "{}:[{};{}]"
        processes = ctx.processset().processid()
        # print [p.getText() for p in processes]

        # removing the code for process groups
        self.rewriter.delete(self.rewriter.DEFAULT_PROGRAM_NAME,
                             ctx.start.tokenIndex, ctx.stop.tokenIndex)

        # Adding the unrolled text for the new procs
        new_procs = []
        for process in processes:
            processCode = new_process_str.format(process.getText(),
                                                 ctx.declaration().getText(),
                                                 ctx.statement().getText())
            # print ctx.parentCtx.getTokens(ctx)
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
        self.replacedone = False

    # def enterRepeat(self, ctx):
    #     cs = ctx.statement().start.getInputStream()
    #     statements = cs.getText(ctx.statement().start.start,
    #                             ctx.statement().stop.stop)
    #     rep_variable = int(ctx.INT().getText())
    #     edited = ''
    #     # removing the code for process groups
    #     self.rewriter.delete(self.rewriter.DEFAULT_PROGRAM_NAME,
    #                          ctx.start.tokenIndex, ctx.stop.tokenIndex)
    #     for var in range(rep_variable):
    #         edited += statements + ";\n"
    #     self.rewriter.insertAfter(ctx.stop.tokenIndex, edited)

    def enterRepeat(self, ctx):
        # Do only one replacement at a time
        if self.replacedone:
            return

        rep_variable = int(ctx.INT().getText())
        # TODO: Is there a way to avoid string manipulation?
        list_statements = ctx.statement()
        cs = list_statements[0].start.getInputStream()
        statements = cs.getText(list_statements[0].start.start,
                                list_statements[-1].stop.stop)
        print "------------------------------"
        print statements
        print "------------------------------"

        new_str = ''
        for var in range(rep_variable):
            new_str += "  " + statements + ";\n"
        self.rewriter.insertAfter(ctx.stop.tokenIndex + 1, new_str)
        self.rewriter.delete(self.rewriter.DEFAULT_PROGRAM_NAME,
                             ctx.start.tokenIndex,
                             ctx.stop.tokenIndex + 1)
        self.replacedone = True


def main(program_str, outfile, filename, debug, skiprename):
    input_stream = InputStream(program_str)
    # lexer = ParallelyLexer(input_stream)
    # stream = CommonTokenStream(lexer)
    # parser = ParallelyParser(stream)

    # tree = parser.parallelprogram()

    fullstart = time.time()

    if not skiprename:
        print "Unrolling Repeat statements"
        while(True):
            lexer = ParallelyLexer(input_stream)
            stream = CommonTokenStream(lexer)
            parser = ParallelyParser(stream)
            tree = parser.parallelprogram()
            unroller = unrollRepeat(stream)
            walker = ParseTreeWalker()
            walker.walk(unroller, tree)
            input_stream = InputStream(unroller.rewriter.getDefaultText())
            print unroller.replacedone
            if not unroller.replacedone:
                print unroller.replacedone
                input_stream = InputStream(unroller.rewriter.getDefaultText())
                break

        # if debug:
        debug_file = open("_DEBUG_UNROLLED_.txt", 'w')
        debug_file.write(input_stream.strdata)
        debug_file.close()

        lexer = ParallelyLexer(input_stream)
        stream = CommonTokenStream(lexer)
        parser = ParallelyParser(stream)
        tree = parser.parallelprogram()

        print "Renaming all variables"
        renamer = VariableRenamer(stream)
        walker = ParseTreeWalker()
        walker.walk(renamer, tree)

        start = time.time()

        # Run type checker on the renamed version
        input_stream = InputStream(renamer.rewriter.getDefaultText())
        # if debug:
        debug_file = open("_DEBUG_RENAMED_.txt", 'w')
        debug_file.write(input_stream.strdata)
        debug_file.close()

    print "Running type checker"
    lexer = ParallelyLexer(input_stream)
    stream = CommonTokenStream(lexer)
    parser = ParallelyParser(stream)

    tree = parser.parallelprogram()

    typechecker = parallelyTypeChecker(debug)
    start = time.time()
    type_checked = typechecker.visit(tree)
    end = time.time()

    if type_checked:
        print "{} passed type checker ({}s).".format(filename, end - start)
    else:
        print "{} failed type checker ({}s).".format(filename, end - start)
        exit(-1)

    # Sequentialization
    start2 = time.time()
    sequentializer = parallelySequentializer(debug)
    sequentializer.rewriteProgram(tree, outfile)
    end2 = time.time()
    print "Time for sequentialization :", end2 - start2

    print "Total time : ", (end2 - start2) + (end - start)


if __name__ == '__main__':
    parser = ArgumentParser()
    parser.add_argument("-f", dest="programfile",
                        help="File containing the code")
    parser.add_argument("-o", dest="outfile",
                        help="File to output the sequential code")
    parser.add_argument("-s", "--skip", action="store_true",
                        help="Skip renaming")
    parser.add_argument("-d", "--debug", action="store_true",
                        help="Print debug info")
    args = parser.parse_args()

    programfile = open(args.programfile, 'r')
    outfile = open(args.outfile, 'w')
    program_str = programfile.read()
    main(program_str, outfile, programfile.name, args.debug, args.skip)
