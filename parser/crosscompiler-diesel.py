import sys
import os
from antlr4 import CommonTokenStream
from antlr4 import InputStream
from antlr4 import TokenStreamRewriter
from antlr4 import ParseTreeWalker
from antlr4 import *
from ParallelyLexer import ParallelyLexer
from ParallelyParser import ParallelyParser
from antlr4.tree.Trees import Trees
from ParallelyVisitor import ParallelyVisitor
from ParallelyListener import ParallelyListener
from argparse import ArgumentParser
import random

key_error_msg = "Type error detected: Undeclared variable (probably : {})"

str_single_thread = '''func {}() {{
  defer parallely.Wg.Done()
  {}
  fmt.Println("Ending thread : ", {});
}}'''

str_member_thread = '''func {}(tid int) {{
  defer parallely.Wg.Done()
  {}
  fmt.Println("Ending thread : ", {});
}}'''

str_probchoiceInt = "{} = parallely.RandchoiceFlag(float32({}), {}, {}, &__flag_{});\n"


def isInt(s):
    try:
        int(s)
        return True
    except ValueError:
        return False


def isGroup(pid):
        if isinstance(pid, ParallelyParser.NamedpContext):
            return (False, pid.getText())
        elif isinstance(pid, ParallelyParser.VariablepContext):
            print "[Error] Cant handle process name variables"
            exit(-1)
            # return (True, pid.VAR().getText(),)
        else:
            return (True, pid.GLOBALVAR().getText(), pid.VAR().getText())


class CountThreads(ParallelyVisitor):
    def __init__(self):
        print "Counting the number of processes"
        self.processes = {}
        self.processcount = 0

    def visitSingleglobaldec(self, ctx):
        global_var = ctx.GLOBALVAR().getText()
        members = [temp.getText() for temp in ctx.processid()]
        self.processes[global_var] = members

    # in theory pids are not int. Changing to simplify implementation
    def visitSingle(self, ctx):
        pid = isGroup(ctx.processid())
        print pid
        if pid[0]:
            self.processcount += len(self.processes[pid[1]])
        else:
            self.processcount += 1


class Translator(ParallelyVisitor):
    def __init__(self):
        print "Starting translation"
        self.pid = None
        self.process_defs = []
        self.process_list = []
        self.globaldecs = []
        self.tempvarnum = 0
        self.recovernum = 0

    def visitSingleglobaldec(self, ctx):
        str_global_dec = "var {} = []int {{{}}};\n"
        varname = ctx.GLOBALVAR().getText()
        members = [t.getText() for t in ctx.processid()]
        # Q = {2,3,4,5};
        global_str = str_global_dec.format(varname, ','.join(members))
        self.globaldecs.append(global_str)

    def visitGlobalconst(self, ctx):
        str_global_dec = "var {} {};\n"
        mytype = self.getType(ctx.basictype())
        varname = ctx.GLOBALVAR().getText()
        # Q = {2,3,4,5};
        global_str = str_global_dec.format(varname, mytype[1])
        self.globaldecs.append(global_str)

    def visitGlobalarray(self, ctx):
        str_global_dec = "var {} []{};\n"
        mytype = self.getType(ctx.basictype())
        varname = ctx.GLOBALVAR().getText()
        # Q = {2,3,4,5};
        global_str = str_global_dec.format(varname, mytype[1])
        self.globaldecs.append(global_str)

    def getType(self, fulltype):
        if isinstance(fulltype, ParallelyParser.BasictypeContext):
            return (fulltype.typequantifier().getText(),
                    fulltype.getChild(1).getText(), 0)
        if isinstance(fulltype, ParallelyParser.SingletypeContext):
            return (fulltype.basictype().typequantifier().getText(),
                    fulltype.basictype().getChild(1).getText(), 0)
        elif isinstance(fulltype, ParallelyParser.ArraytypeContext):
            return (fulltype.basictype().typequantifier().getText(),
                    fulltype.basictype().getChild(1).getText(), 1)
        else:
            print "[Error] Unknown type : ", fulltype.getText()
            exit(-1)

    def visitCondsend(self, ctx):
        cond_send_str = {
            ("int", 0): "parallely.Condsend({}, {}, {}, {});\n",
            ("int32", 0): "parallely.CondsendInt32({}, {}, {}, {});\n",
            ("int64", 0): "parallely.CondsendInt64({}, {}, {}, {});\n",
            ("float32", 0): "parallely.CondsendFloat32({}, {}, {}, {});\n",
            ("float64", 0): "parallely.CondsendFloat64({}, {}, {}, {});\n",
            ("int", 1): "parallely.CondsendIntArray({}, {}[:], {}, {});\n",
        }

        cond_var = ctx.var(0).getText()
        sent_var = ctx.var(1).getText()
        senttype = self.getType(ctx.fulltype())

        return cond_send_str[(senttype[1], senttype[2])].format(cond_var, sent_var,
                                                                self.pid, ctx.processid().getText())

    def visitSend(self, ctx):
        send_str = {
            ("int", 0): "parallely.SendInt({}, {}, {});\n",
            ("int32", 0): "parallely.SendInt32({}, {}, {});\n",
            ("int64", 0): "parallely.SendInt64({}, {}, {});\n",
            ("float32", 0): "parallely.SendFloat32({}, {}, {});\n",
            ("float64", 0): "parallely.SendFloat64({}, {}, {});\n",
            ("int", 1): "parallely.SendIntArray({}, {}, {});\n",
            ("int32", 1): "parallely.SendInt32Array({}[:], {}, {});\n",
            ("int64", 1): "parallely.SendInt64Array({}[:], {}, {});\n",
            ("float32", 1): "parallely.SendFloat32Array({}[:], {}, {});\n",
            ("float64", 1): "parallely.SendFloat64Array({}[:], {}, {});\n"
        }
        sent_var = ctx.var().getText()
        senttype = self.getType(ctx.fulltype())
        return send_str[senttype[1], senttype[2]].format(sent_var, self.pid, ctx.processid().getText())

    def visitReceive(self, ctx):
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

        senttype = self.getType(ctx.fulltype())
        return rec_str[senttype[1], senttype[2]].format(ctx.var().getText(),
                                                        self.pid, ctx.processid().getText())

    def visitCondreceive(self, ctx):
        rec_str = {
            ("int", 0): "parallely.Condreceive(&{}, &{}, {}, {});\n",
            ("int32", 0): "parallely.CondreceiveInt32(&{}, &{}, {}, {});\n",
            ("int64", 0): "parallely.CondreceiveInt64(&{}, &{}, {}, {});\n",
            # ("float32", 0):
            # ("float64", 0):
            ("int", 1): "parallely.CondreceiveIntArray(&{}, {}[:], {}, {});\n",
            # ("int32", 1):
            # ("int64", 1):
            ("float32", 1): "parallely.CondreceiveFloat32(&{}, &{}, {}, {});\n",
            ("float64", 1): "parallely.CondreceiveFloat64(&{}, &{}, {}, {});\n"
        }
        senttype = self.getType(ctx.fulltype())
        return senttype[senttype[1], senttype[2]].format(ctx.var(0).getText(), ctx.var(1).getText(),
                                                         self.pid, ctx.processid().getText())

    def visitProbassignment(self, ctx):
        assigned_var = ctx.var().getText()
        prob = ctx.probability().getText()
        return str_probchoiceInt.format(assigned_var, prob, ctx.expression(0).getText(),
                                        ctx.expression(1).getText(), self.recovernum)

    def handleExpression(self, ctx):
        convert_str = "parallely.ConvBool({})"
        if isinstance(ctx, ParallelyParser.SelectContext):
            return self.handleExpression(ctx.expression())
        if (isinstance(ctx, ParallelyParser.EqContext) or isinstance(ctx, ParallelyParser.GeqContext) or isinstance(ctx, ParallelyParser.LeqContext) or isinstance(ctx, ParallelyParser.LessContext) or isinstance(ctx, ParallelyParser.GreaterContext) or isinstance(ctx, ParallelyParser.AndexpContext)):
            return convert_str.format(ctx.getText())
        else:
            return ctx.getText()

    def visitExpassignment(self, ctx):
        assign_str = "{} = {};\n"
        expr_str = self.handleExpression(ctx.expression())
        var_str = ctx.var().getText()
        return assign_str.format(var_str, expr_str)

    def visitGexpassignment(self, ctx):
        assign_str = "{} = {};\n"
        expr_str = self.handleExpression(ctx.expression())
        var_str = ctx.GLOBALVAR().getText()
        return assign_str.format(var_str, expr_str)

    def visitArrayload(self, ctx):
        # print ctx.getText()
        return ctx.getText() + ";\n"

    def visitArrayassignment(self, ctx):
        # print ctx.getText()
        return ctx.getText() + ";\n"

    def visitCast(self, ctx):
        resultType = self.getType(ctx.fulltype())
        assignedvar = ctx.var(0).getText()
        castedvar = ctx.var(1).getText()
        # Array type
        if resultType[1] == "float64" and resultType[2] == 1:
            return "parallely.Cast32to64Array({}[:], {}[:]);\n".format(assignedvar,
                                                                       castedvar)
        if resultType[1] == "float32" and resultType[2] == 1:
            return "parallely.Cast64to32Array({}[:], {}[:]);\n".format(assignedvar,
                                                                       castedvar)
        # Regular cast
        if resultType[1] == "float64" and resultType[2] == 0:
            return "{} = float64({});\n".format(assignedvar, castedvar)
        if resultType[1] == "float32" and resultType[2] == 0:
            return "{} = float32({});\n".format(assignedvar, castedvar)

    def visitIfonly(self, ctx):
        str_if_only = "if {} != 0 {{\n {} }}\n"
        cond_var = ctx.var().getText()

        statement_string = ''
        for statement in ctx.ifs:
            translated = self.visit(statement)
            if translated:
                statement_string += translated
            else:
                print "[Error] Unable to transtate: ", statement.getText()
                exit(-1)
        # print str_if_only.format(cond_var, statement_string)
        return str_if_only.format(cond_var, statement_string)

    def visitIf(self, ctx):
        str_if = "if {} != 0 {{\n {} }} else {{\n {} }}\n"
        cond_var = ctx.var().getText()

        statement_string = ''
        for statement in ctx.ifs:
            translated = self.visit(statement)
            if translated:
                statement_string += translated
            else:
                print "[Error] Unable to transtate: ", statement.getText()
                exit(-1)

        else_statement_string = ''
        for statement in ctx.elses:
            translated = self.visit(statement)
            if translated:
                else_statement_string += translated
            else:
                print "[Error] Unable to transtate: ", statement.getText()
                exit(-1)
        # print str_if_only.format(cond_var, statement_string)
        return str_if.format(cond_var, statement_string, else_statement_string)

    def visitRepeatlvar(self, ctx):
        repeatVar = ctx.var().getText()
        temp_var_name = "__temp_{}".format(self.tempvarnum)
        self.tempvarnum += 1

        statement_string = ''
        for statement in ctx.statement():
            translated = self.visit(statement)
            if translated:
                statement_string += translated
            else:
                print "[Error] Unable to transtate: ", statement.getText()
                exit(-1)
        str_for_loop = "for {} := 0; {} < {}; {}++ {{\n {} }}\n"
        return str_for_loop.format(temp_var_name, temp_var_name, repeatVar, temp_var_name, statement_string)

    def visitRecover(self, ctx):
        self.recovernum += 1

        starting_level = self.recovernum
        temp_flag_name = "__flag_{}".format(self.recovernum)

        recover_str = "{} := false;\n {}\n if {} {{\n {} = false;\n {}\n }}\n {}\n"

        try_statement_string = ''
        for statement in ctx.trys:
            translated = self.visit(statement)
            if translated:
                try_statement_string += translated
            else:
                print "[Error] Unable to transtate: ", statement.getText()
                exit(-1)
        recovers_statement_string = ''
        for statement in ctx.trys:
            translated = self.visit(statement)
            if translated:
                recovers_statement_string += translated
            else:
                print "[Error] Unable to transtate: ", statement.getText()
                exit(-1)

        how_deep = self.recovernum
        combine_str = ""

        self.recovernum -= 1

        # For each outer flag set it to the value of the inner flags
        for f in range(starting_level - 1):
            temp_str = ""
            oneflag = "__flag_{} = __flag_{} {};\n"
            for f in range(starting_level - 1):
               temp_str += "|| __flag_{}".format(starting_level + f)
            combine_str += oneflag.format(f+1, f+1, temp_str)

        final_str = combine_str

        # final_str = final_str.format(temp_flag_name, temp_flag_name + combine_str)

        program_str = recover_str.format(temp_flag_name,
                                         try_statement_string,
                                         temp_flag_name,
                                         temp_flag_name,
                                         recovers_statement_string,
                                         final_str)

        return program_str

    def visitRepeat(self, ctx):
        repeatNum = ctx.INT().getText()
        temp_var_name = "__temp_{}".format(self.tempvarnum)
        self.tempvarnum += 1

        statement_string = ''
        for statement in ctx.statement():
            translated = self.visit(statement)
            if translated:
                statement_string += translated
            else:
                print "[Error] Unable to transtate: ", statement.getText()
                exit(-1)
        str_for_loop = "for {} := 0; {} < {}; {}++ {{\n {} }}\n"
        return str_for_loop.format(temp_var_name, temp_var_name, repeatNum, temp_var_name, statement_string)

    def visitForloop(self, ctx):
        group_name = ctx.GLOBALVAR().getText()
        var_name = ctx.VAR().getText()
        # for proc in self.proc_groups[group_name]:
        statement_string = ''
        for statement in ctx.statement():
            translated = self.visit(statement)
            if translated:
                statement_string += translated
            else:
                print "[Error] Unable to transtate: ", statement.getText()
                exit(-1)
        str_for_loop = "for _, {} := range({}) {{\n {} }}\n"
        return str_for_loop.format(var_name, group_name, statement_string)

    def visitFunc(self, ctx):
        return ctx.getText() + ";\n"

    def isGroup(self, pid):
        if isinstance(pid, ParallelyParser.NamedpContext):
            return (False, pid.getText())
        elif isinstance(pid, ParallelyParser.VariablepContext):
            print "[Error] Cant handle process name variables"
            exit(-1)
            # return (True, pid.VAR().getText(),)
        else:
            return (True, pid.GLOBALVAR().getText(), pid.VAR().getText())

    def handleDec(self, decl):
        str_single_dec = "var {} {};\n"
        str_array_dec = "var {} {}{};\n"
        dectype = self.getType(decl.basictype())

        # Array declaration
        if isinstance(decl, ParallelyParser.ArraydeclarationContext):
            varname = decl.var().getText()
            dim = ""
            if decl.INT():
                for dimention in decl.INT():
                    dim += "[{}]".format(dimention)
            else:
                dim += "[]"
            return str_array_dec.format(varname, dim, dectype[1])
        else:
            varname = decl.var().getText()
            return str_single_dec.format(varname, dectype[1])

    def handleGroup(self, group_name, element_name, ctx):
        print group_name, element_name, self.proc_groups
        # for proc in self.proc_groups[group_name]:
        self.pid = "tid"
        print "Translating process group: ", group_name

        # Collect the declarations which should be at the top
        # Binding the pid int to the process name.
        dec_string = "{} := tid;\n".format(element_name)
        for decl in ctx.declaration():
            dec_string += self.handleDec(decl)

        statement_string = ""
        for statement in ctx.statement():
            translated = self.visit(statement)
            if translated:
                statement_string += translated
            else:
                print "[Error] Unable to transtate: ", statement.getText()
                exit(-1)

        process_name = "func_" + group_name
        self.process_list.append(((process_name, group_name), 1))

        process_code = dec_string + statement_string

        process_def_str = str_member_thread.format(process_name, process_code, element_name)
        self.process_defs.append(process_def_str)

    def visitSingle(self, ctx):
        if self.isGroup(ctx.processid())[0]:
            self.handleGroup(self.isGroup(ctx.processid())[1], self.isGroup(ctx.processid())[2], ctx)
            return

        self.pid = ctx.processid().getText()

        print "Translating process: ", self.pid

        # Collect the declarations which should be at the top
        dec_string = ""
        for decl in ctx.declaration():
            dec_string += self.handleDec(decl)
            # dectype = self.getType(decl.fulltype())[1]
            # varname = decl.var().getText()
            # dec_string += str_single_dec.format(varname, dectype)

        statement_string = ""
        for statement in ctx.statement():
            self.recovernum = 0
            translated = self.visit(statement)
            if translated:
                statement_string += translated
            else:
                print "[Error] Unable to transtate: ", statement.getText()

        process_name = "func_" + self.pid
        self.process_list.append((process_name, 0))

        process_code = dec_string + statement_string

        process_def_str = str_single_thread.format(process_name, process_code, self.pid)
        # print "--------------------"
        # print process_def_str
        # print "--------------------"
        self.process_defs.append(process_def_str)

    def translate(self, tree, numthreads, proc_groups_in, fout_name, template):
        self.proc_groups = proc_groups_in
        self.visit(tree)

        # print self.process_defs

        all_global_decs = ''.join(self.globaldecs)
        all_process_defs = '\n'.join(self.process_defs)

        run_procs = ''
        for fname, is_group in self.process_list:
            if is_group:
                run_procs += "for _, index := range {} {{\ngo {}(index);\n}}\n".format(fname[1], fname[0])
            else:
                run_procs += "go {}();\n".format(fname)

        # print "--------------------"
        # print all_process_defs
        # print run_procs
        # print "--------------------"

        # There has to be a better way to read in the template
        # __location__ = os.path.realpath(os.path.join(os.getcwd(),
        #                                              os.path.dirname(__file__)))

        # template_str = open(os.path.join(__location__, '__basic_go.txt'), 'r').readlines()
        template_str = open(template, 'r').readlines()
        with open(fout_name, "w") as fout:
            for line in template_str:
                newline = line.replace('__NUM_THREADS__', str(numthreads))
                newline = newline.replace('__GLOBAL_DECS__', all_global_decs)
                newline = newline.replace('__FUNC_DECS__', all_process_defs)
                newline = newline.replace('__START__THREADS__', run_procs)
                fout.write(newline)


def main(program_str, outfile, filename, template, debug):
    print "Starting the cross compilation"
    input_stream = InputStream(program_str)
    lexer = ParallelyLexer(input_stream)
    stream = CommonTokenStream(lexer)
    parser = ParallelyParser(stream)

    tree = parser.parallelprogram()

    threadcounter = CountThreads()
    threadcounter.visit(tree)

    print "Number of processes found: {}".format(threadcounter.processcount)

    translator = Translator()
    translator.translate(tree, threadcounter.processcount, threadcounter.processes, outfile, template)

    # tree = parser.parallelprogram()
    # renamer = IdentifyChannels()
    # walker = ParseTreeWalker()
    # walker.walk(renamer, tree)

    # channels, ch_decs = renamer.getChannelSet()
    # print ch_decs

    # var_definition = ''
    # for channel in channels:
    #     ch_name = channels[channel]
    #     ch_type = channel[3]
    #     var_definition += "var {} chan {}\n".format(ch_name, ch_type)

    # gotranslator = Translator(channels, debug)
    # generated_program, evocation_str = gotranslator.visit(tree)

    # print "Generating the program file"
    # template = open("__golang_template.txt", "rt")
    # for line in template.readlines():
    #     newline = line.replace('__VAR_DEFS__', var_definition)
    #     newline = newline.replace('__CHANEL_MAKES__', '\n'.join(ch_decs))
    #     newline = newline.replace('__FUNC_DEFS__', generated_program)
    #     newline = newline.replace('__FUNCTION_CALLS__', evocation_str)
    #     outfile.write(newline)

    # outfile.close()


if __name__ == '__main__':
    parser = ArgumentParser()
    parser.add_argument("-f", dest="programfile",
                        help="File containing the code", required=True)
    parser.add_argument("-o", dest="outfile",
                        help="File to output the sequential code", required=True)
    parser.add_argument("-d", "--debug", action="store_true",
                        help="Print debug info")
    parser.add_argument("-t", dest="template",
                        help="File containing the template", required=True)
    args = parser.parse_args()

    programfile = open(args.programfile, 'r')
    # outfile = open(args.outfile, 'w')
    program_str = programfile.read()
    main(program_str, args.outfile, programfile.name, args.template, args.debug)
