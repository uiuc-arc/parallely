from ParallelyParser import ParallelyParser
from ParallelyVisitor import ParallelyVisitor
from ParallelyLexer import ParallelyLexer
from antlr4 import InputStream
from antlr4 import CommonTokenStream

from argparse import ArgumentParser
import time
import copy

class parallelySequentializer(ParallelyVisitor):
    def __init__(self, debug, annotate):
        self.group_map = {}
        self.statement_lists = {}
        self.msgcontext = {}
        self.globaldecs = {}
        self.declarations = []
        self.grouped_list = {}
        self.isProcessGroup = {}
        self.debug = debug
        self.annotate = annotate

    def debugMsg(self, msg):
        if self.debug:
            print("[Debug - TypeChecker] " + msg)

    def exitWithError(self, msg):
        print("[Error - TypeChecker]: " + msg)
        exit(-1)

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

        # print pid.getText(), is_group, len(statements)

        if not is_group[0]:
            self.statement_lists[pid.getText()] = statements
        else:
            self.group_map[pid.GLOBALVAR().getText()] = pid.VAR().getText()
            self.statement_lists[pid.GLOBALVAR().getText()] = statements

    def visitMultipleglobaldec(self, ctx):
        self.visit(ctx.getChild(0))
        self.visit(ctx.getChild(2))

    def visitSingleglobaldec(self, ctx):
        ind = ctx.GLOBALVAR().getText()
        members = ctx.processid()
        self.globaldecs[ind] = members

    def getType(self, fulltype):
        if isinstance(fulltype, ParallelyParser.SingletypeContext):
            return (fulltype.basictype().typequantifier().getText(),
                    fulltype.basictype().getChild(1).getText(), 0)
        elif isinstance(fulltype, ParallelyParser.ArraytypeContext):
            return (fulltype.basictype().typequantifier().getText(),
                    fulltype.basictype().getChild(1).getText(), 1)
        else:
            print "[Error] Unknown type : ", fulltype
            exit(-1)

    def annotateStr(self, str1, str2):
        if self.annotate:
            return str1 + " @ " + str(str2) + ";\n"
        else:
            return str1 + ";\n"

    def getDecString(self, dec):
        cs = dec.start.getInputStream()
        text = cs.getText(dec.start.start, dec.stop.stop)
        # dec_type_q, dec_type_t = self.getType(dec.fulltype())
        # # dec_type_t = statement.fulltype().getChild(1).getText()

        # dec_name = dec.var().getText()
        # newdec = "{} {} {};".format(dec_type_q, dec_type_t, dec_name)
        return self.annotateStr(text, str(dec.start.line))

    def appendIfExists(self, key, dict_in, val):
        if key in dict_in.keys():
            dict_in[key].append(val)
        else:
            dict_in[key] = [val]

    def handleSend(self, pid, statement, statement_list, msgcontext, seq_prefix):
        rec = statement.sender.getText()
        sent_type = self.getType(statement.fulltype())
        # = statement.fulltype().getChild(1).getText()
        sent_var = statement.var(1).getText()
        my_key = (rec, pid, sent_type)
        # print "1=======> ", sent_var, my_key, msgcontext
        self.appendIfExists(my_key, msgcontext, (sent_var, statement.start.line))
        # print "2=======> ", msgcontext
        new_statement_list = dict(statement_list)
        new_statement_list[pid].pop(0)
        return True, seq_prefix, dict(msgcontext), new_statement_list

    def handleCondSend(self, pid, statement, statement_list, msgcontext, seq_prefix):
        rec = statement.processid().getText()
        sent_type = self.getType(statement.fulltype())
        # sent_type_t = statement.fulltype().getChild(1).getText()
        guard_var = statement.var()[0].getText()
        sent_var = statement.var()[1].getText()
        my_key = (rec, pid, sent_type)
        self.appendIfExists(my_key, msgcontext, (sent_var, guard_var, statement.start.line))
        statement_list[pid].pop(0)
        return True, seq_prefix, msgcontext, statement_list

    def handleReceive(self, pid, statement, statement_list, msgcontext, seq_prefix):
        assigned_var = statement.var().getText()
        sender = statement.processid().getText()

        # if sender not in statement_list.keys():
        #     print ("[ERROR] Receiving from an unknown sender!")
        #     print ("[ERROR] ", sender, statement.getText(), statement_list)
        #     return False, seq_prefix, msgcontext, statement_list

        sent_type_t = self.getType(statement.fulltype())
        assign_symbol = statement.getChild(1).getText()
        my_key = (pid, sender, sent_type_t)
        # print(my_key, msgcontext.keys())

        # if the msgcontext is empty or top is a guarded expression exit
        if my_key in msgcontext.keys() and len(msgcontext[my_key]) > 0:
            rec_val = msgcontext[my_key].pop(0)
            rewrite = "{} {} {}".format(assigned_var, assign_symbol, rec_val[0])
            seq_prefix.append(self.annotateStr(rewrite,
                                               "{}, {}".format(rec_val[1],
                                                               statement.start.line)))
            statement_list[pid].pop(0)
            return True, seq_prefix, msgcontext, statement_list
        else:
            return False, seq_prefix, msgcontext, statement_list

    def handleCondReceive(self, pid, statement, statement_list, msgcontext, seq_prefix):
        guard_var = statement.var()[0].getText()
        assigned_var = statement.var()[1].getText()
        sender = statement.processid().getText()
        sent_type = self.getType(statement.fulltype())

        # print "====>", self.isProcessGroup, sender
        if sender in self.isProcessGroup and self.isProcessGroup[sender][0]:
            print ("[ERROR] Receiving from a group is not supported yet!")
            # self.rewrite_statements(seq_prefix, msgcontext, remaining_statements)
            exit(-1)
        else:
            my_key = (pid, sender, sent_type)
            if my_key in msgcontext.keys() and len(msgcontext[my_key]) > 0:
                rec_val, rec_guard, send_line = msgcontext[my_key].pop(0)
                out_format = "{} = 1 [{}] 0;\n{}={} [{}] {}"
                rewrite = out_format.format(guard_var, rec_guard,
                                            assigned_var, rec_val,
                                            rec_guard,
                                            assigned_var)
                seq_prefix.append(self.annotateStr(rewrite, "{}, {}".format(send_line,
                                                                            statement.start.line)))
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
        seq_prefix.append(self.annotateStr(result, str(statement.start.line)))
        return True, seq_prefix, msgcontext, statement_list

    def isEmptyMsgContext(self, msg_context):
        for key in msg_context.keys():
            if len(msg_context[key]) > 0:
                return False
        return True

    def handleFor(self, pid, statement, statement_list, msgcontext, seq_prefix):
        # For now assuming that the variable groups have the same iterator
        # this only handles for q in Q type loops
        out_template = "for {} in {} do {{\n{}\n}};\n"

        my_statements = statement.statement()
        target_group = statement.GLOBALVAR().getText()

        if not (target_group in statement_list.keys()):
            group_statements = []
        else:
            group_statements = list(statement_list[target_group])
        group_var = self.isProcessGroup[target_group][2]

        # print("#######################", statement.getText())
        # print("#######################", group_statements)
        # Need to check if pid itself is a group
        # Rewrite the smallest part from the group
        limit = 1
        found = False
        found_limit = -1
        found_output = {}
        while(True):
            if limit > len(group_statements):
                break

            tmp_statements = {}
            tmp_statements[pid] = list(my_statements)
            tmp_statements[group_var] = list(group_statements[:limit])

            # print("++++: ", tmp_statements)
            output = self.rewrite_statements([], {}, tmp_statements)
            # print("++++: ", output)
            
            # if (not output[0]) and found:
            #     print("****111111111111: ", output, found_limit)
            #     break
            if (len(output[2].keys()) == 0 and
                not (pid in output[3].keys()) and
                not (group_var in output[3].keys())):
                # print("****22222222222222: ", output, found_limit)
                found = True
                found_limit = limit
                found_output = tuple(output)
            elif found:
                # print("****111111111111: ", output, found_limit)                
                break
            limit += 1
        if found:
            new_statement_list = statement_list.copy()
            new_statement_list[pid].pop(0)
            # if group_var in output[3].keys():
            # print("@@@@@@@@@@@@@@@@@@@ New: ", group_statements[found_limit:], found_limit, len(group_statements))
            new_statement_list[target_group] = group_statements[found_limit:]
            # else:
            #     new_statement_list[target_group] = []
            new_seq = list(seq_prefix)
            translated_block = ''.join(found_output[1])
            new_seq.append(out_template.format(group_var, target_group, translated_block))
            # print("@@@@@@@@@@@@@@@@@@@ New: ", new_seq)
            # print("@@@@@@@@@@@@@@@@@@@ New: ", new_statement_list)
            return True, new_seq, msgcontext, new_statement_list
        else:
            print("===============================: ", found_limit, limit)
            return False, seq_prefix, msgcontext, statement_list
        # # limit = 0
        # # while True:
        #     # if limit > len(group_statements):
        #     #     # if self.debug:
        #     #     # print "Giving up : ", seq_prefix, msgcontext, statement_list
        #     #     return False, seq_prefix, msgcontext, statement_list

        #     tmp_statements = {}
        #     tmp_statements[pid] = list(my_statements)
        #     tmp_statements[group_var] = list(group_statements)
        #     # tmp_msgcontext = dict(msgcontext)

        #     output = self.rewrite_statements([], {}, tmp_statements)
        #     if self.isEmptyMsgContext(output[1]) and (pid not in output[2]):
        #         break
        #     if self.isEmptyMsgContext(output[1]) and (pid in output[2]) and len(output[2][pid]) == 0:
        #         break
        #     limit += 1

        # # Entire process was rewritten
        # if limit == 0:
        #     statement_list[pid].pop(0)
        #     # statement_list.pop(target_group, None)
        #     # print "--------------------------------------------"
        #     # print len(group_statements), limit, output
        #     # print "--------------------------------------------"
        #     if group_var in output[2]:
        #         statement_list[target_group] = output[2][group_var]
        #     else:
        #         statement_list[target_group] = []
        #         # print "Entire process was rewritten"
        #     # print "--------------------------------------------"

        #     rewrite = out_template.format(group_var, target_group, ''.join(output[0]) + ';')
        #     seq_prefix.append(self.annotateStr(rewrite, str(statement.start.line)))
        # # Only part of the process was rewritten
        # else:
        #     print statement_list.keys(), target_group, statement_list, seq_prefix, limit
        #     statement_list[pid].pop(0)
        #     statement_list[target_group] = group_statements[:limit]

        #     rewrite = out_template.format(group_var, target_group, ''.join(output[0]) + ';')
        #     seq_prefix.append(self.annotateStr(rewrite, str(statement.start.line)))

        # # print statement_list.keys(), target_group, statement_list
        # return True, seq_prefix, msgcontext, statement_list

    def handleRepeatVar(self, pid, statement, statement_list, msgcontext, seq_prefix):
        # TODO: *** Do the renaming step ***
        # For now assuming that the variable groups have the same iterator
        iter_number = statement.GLOBALVAR().getText()
        out_template = "repeat {} {{\n{}}}"
        my_statements = statement.statement()
        # target_group = statement.GLOBALVAR().getText()
        for threads in statement_list.keys():
            if threads != pid:
                group_var = self.isProcessGroup[threads][1]
                # print statement_list, group_var
                group_statements = statement_list[group_var]
                if isinstance(group_statements[0], ParallelyParser.RepeatvarContext):
                    tmp_statements = {}
                    tmp_statements[pid] = list(my_statements)
                    tmp_statements[group_var] = group_statements[0].statement()
                    tmp_msgcontext = dict(msgcontext)

                    output = self.rewrite_statements([], tmp_msgcontext, tmp_statements)
                    if (tmp_msgcontext != msgcontext):
                        print "[Error] rewriting repeatvar msgcontexts dont match"
                        exit(-1)
                    if (tmp_statements != {}):
                        print "[Error] rewriting repeatvar not both empty"
                        exit(-1)
                    seq_prefix.append(self.annotateStr(out_template.format(iter_number, ''.join(output[0])),
                                                       str(statement.start.line)))
                    statement_list[group_var].pop(0)
                    statement_list[pid].pop(0)
                    return True, seq_prefix, msgcontext, statement_list
        return False, seq_prefix, msgcontext, statement_list

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
        if isinstance(statement, ParallelyParser.WhileContext):
            print "while Implement"
            exit(-1)
        if isinstance(statement, ParallelyParser.RepeatvarContext):
            return self.handleRepeatVar(pid, statement, statement_list, msgcontext, seq_prefix)
            # out_template = "repeat {} {{{}}}"
            # bool_var = statement.VAR().getText()
            # # if_state = statement.statement().getText()
            # cs = statement.statement().start.getInputStream()
            # statements = cs.getText(statement.statement().start.start,
            #                         statement.statement().stop.stop)

            # result = out_template.format(bool_var, statements)
            # return True, result, msgcontext
        else:
            cs = statement.start.getInputStream()
            text = cs.getText(statement.start.start, statement.stop.stop)
            newstatementlist = statement_list.copy()
            newstatementlist[pid].pop(0)
            new_seq_prefix = list(seq_prefix)
            new_seq_prefix.append(self.annotateStr(text, statement.start.line))
            return True, new_seq_prefix, msgcontext, newstatementlist

    # def rewritePair(self, pid1, pid2, statement, msgcontext):
    #     temp_statement = {}
    #     temp_statement[pid1] = statement[pid1]
    #     for i in range(len(statement[pid2])):
    #         temp_statement[pid2] = statement[pid2][:i + 1]
    #         success = self.doRewriteProgram(pid1, temp_statement, msgcontext)
    #         if success[0]:
    #             return True, success[1], statement[pid2][i + 1:]
    #     return False, '', statement[pid2]

    def tryGroupedContextRule(self, statement_list, pid):
        out_template = "for {} in {} do {{\n{}\n}}"
        movables = []
        i = 0
        changed = False
        for statement in statement_list:
            if isinstance(statement, ParallelyParser.SendContext):
                break
            if isinstance(statement, ParallelyParser.ReceiveContext):
                break
            if isinstance(statement, ParallelyParser.CondsendContext):
                break
            if isinstance(statement, ParallelyParser.CondreceiveContext):
                break
            if isinstance(statement, ParallelyParser.ForloopContext):
                # might still be doable if there is no commuication?
                break
            if isinstance(statement, ParallelyParser.RepeatvarContext):
                break
            cs = statement.start.getInputStream()
            text = cs.getText(statement.start.start, statement.stop.stop + 2)
            movables.append(text)
            i += 1
        if len(movables) > 0:
            changed = True
        return changed, out_template.format(pid[-1], pid[-2], ''.join(movables)), statement_list[i:]

    def isGroupedProcess(self, pid):
        # pid created in the renaming process (I assume)
        if pid not in self.isProcessGroup:
            return False
        else:
            return self.isProcessGroup[pid][0]

    def rewrite_statements(self, seq_prefix, msgcontext, remaining_statements):
        while(True):
            for key in remaining_statements.keys():
                if len(remaining_statements[key]) == 0:
                    remaining_statements.pop(key, None)
            for key in msgcontext.keys():
                if len(msgcontext[key]) == 0:
                    msgcontext.pop(key, None)
            # print(seq_prefix, msgcontext, remaining_statements)

            if len(msgcontext.keys()) == 0 and len(remaining_statements.keys()) == 0:
                for key in remaining_statements.keys():
                    if len(remaining_statements[key]) == 0:
                        remaining_statements.pop(key, None)
                for key in msgcontext.keys():
                    if len(msgcontext[key]) == 0:
                        msgcontext.pop(key, None)
                return True, seq_prefix, msgcontext, remaining_statements
            # if len(remaining_statements.keys()==0):
            #     break
            changed = False
            for pid in remaining_statements.keys():
                if self.isGroupedProcess(pid):
                    continue
                statement = remaining_statements[pid][0]
                output = self.rewriteOneStep(pid, statement,
                                             remaining_statements,
                                             msgcontext, seq_prefix)
                if output[0]:
                    success, seq_prefix, msgcontext, remaining_statements = output
                    removed = False
                    changed = True
                    for key in remaining_statements.keys():
                        if len(remaining_statements[key]) == 0:
                            remaining_statements.pop(key, None)
                            removed = True
                    if removed:
                        break
            if not changed:
                for key in remaining_statements.keys():
                    if len(remaining_statements[key]) == 0:
                        remaining_statements.pop(key, None)
                for key in msgcontext.keys():
                    if len(msgcontext[key]) == 0:
                        msgcontext.pop(key, None)
                return False, seq_prefix, msgcontext, remaining_statements

        # self.debugMsg("[Debug - Seq] {} {} {}".format(remaining_statements,
        #                                               seq_prefix, msgcontext))

        # remaining_pids = set(remaining_statements.keys())
        # while(True):
        #     changed = False
        #     for pid in remaining_pids.copy():
        #         # print pid, changed, remaining_statements
        #         if self.isGroupedProcess(pid):
        #             if pid not in remaining_statements:
        #                 continue
        #             my_changed, moved, remaining = self.tryGroupedContextRule(remaining_statements[pid],
        #                                                                       self.isProcessGroup[pid])
        #             if my_changed:
        #                 changed = changed or my_changed
        #                 seq_prefix.append(moved)
        #                 remaining_statements[pid] = remaining
        #             # remaining_pids.remove(pid)
        #             # self.debugMsg("Sequential prefix: " + seq_prefix)

        #             if len(remaining_statements[pid]) == 0:
        #                 remaining_pids.remove(pid)
        #                 changed = True
        #                 if self.debug:
        #                     print "[Debug:rewrite_statements] : completely sequentialized 2 : ",
        #                     pid, remaining_statements
        #                 remaining_statements.pop(pid, None)
        #             continue
        #         # If all statements from a pid is removed
        #         if not (pid in remaining_statements.keys()):
        #             self.debugMsg("[Debug:rewrite_statements] : already completely sequentialized : {} {}".format(
        #                 pid, remaining_statements))
        #             break
        #         if len(remaining_statements[pid]) == 0:
        #             remaining_pids.remove(pid)
        #             changed = True
        #             self.debugMsg("[Debug:rewrite_statements] : completely sequentialized : {} {}".format(
        #                 pid, remaining_statements))
        #             remaining_statements.pop(pid, None)
        #             continue

        #         statement = remaining_statements[pid][0]
        #         # print "[Debug] ", seq_prefix
        #         output = self.rewriteOneStep(pid, statement,
        #                                      remaining_statements,
        #                                      msgcontext, seq_prefix)
        #         success, seq_prefix, msgcontext, remaining_statements = output

        #         if success:
        #             # statements[pid].pop(0)
        #             # if result != '':
        #             #     seq_prefix.append(result)
        #             changed = True

        #     # If no rewrite is possible
        #     if not changed:
        #         break
        # if self.debug:
        #     print "[Debug:rewrite_statements:2] ", remaining_statements, seq_prefix, msgcontext
        # return seq_prefix, msgcontext, remaining_statements

    def rewriteProgram(self, tree, outfile):
        # Build the statement lists
        self.visit(tree)

        msgcontext = {}

        # Move all declarations to the top of the program
        temp = []
        for key in self.globaldecs:
            array_str = ', '.join([a.getText() for a in self.globaldecs[key]])
            temp.append("{}={{{}}};\n".format(key, array_str))
        for dec in self.declarations:
            temp.append(self.getDecString(dec))
        global_decs_str = ''.join(temp)
        self.debugMsg("GLOBAL VARS: [" + global_decs_str + "\n]")

        statements = self.statement_lists.copy()
        rewritten = self.rewrite_statements([], msgcontext, statements)

        # print rewritten[1]
        if not rewritten[0]:
            print("Sequentializion failed")
            # print("Remaining statements: ", str(rewritten[3][:3])[:10])
            outfile.write(str(rewritten[1]) + str(rewritten[3]))
            print("Remaining Messages: ", rewritten[2])
            exit(-1)

        seq_program = global_decs_str + "\n" + "".join(rewritten[1])
        outfile.write(seq_program)
        if self.debug:
            print "Sequentialized Program:"
            print seq_program


if __name__ == '__main__':
    parser = ArgumentParser()
    parser.add_argument("-f", dest="programfile",
                        help="File containing the code", required=True)
    parser.add_argument("-o", dest="outfile",
                        help="File to output the sequential code", required=True)
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

    tree = parser.parallelprogram()

    # Sequentialization
    start2 = time.time()
    sequentializer = parallelySequentializer(args.debug, args.annotate)
    sequentializer.rewriteProgram(tree, outfile)
    end2 = time.time()
    print "Time for sequentialization :", end2 - start2
