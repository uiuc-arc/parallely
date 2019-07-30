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

str_probchoiceInt = "{} = parallely.Randchoice(float32({}), {}, {});\n"


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
        str_condsendInt = "parallely.Condsend({}, {}, {}, {});\n"
        str_condsendInt32 = "parallely.CondsendInt32({}, {}, {}, {});\n"
        str_condsendInt64 = "parallely.CondsendInt64({}, {}, {}, {});\n"
        str_condsendFloat32 = "parallely.CondsendFloat32({}, {}, {}, {});\n"
        str_condsendFloat64 = "parallely.CondsendFloat64({}, {}, {}, {});\n"

        str_condsendIntArray = "parallely.CondsendIntArray({}, {}[:], {}, {});\n"

        cond_var = ctx.var(0).getText()
        sent_var = ctx.var(1).getText()
        senttype = self.getType(ctx.fulltype())
        if senttype[1] == "int" and senttype[2] == 0:
            return str_condsendInt.format(cond_var, sent_var,
                                          self.pid, ctx.processid().getText())
        if senttype[1] == "int" and senttype[2] == 1:
            return str_condsendIntArray.format(cond_var, sent_var,
                                               self.pid, ctx.processid().getText())
        if senttype[1] == "int32" and senttype[2] == 0:
            return str_condsendInt32.format(cond_var, sent_var,
                                            self.pid, ctx.processid().getText())
        if senttype[1] == "int64" and senttype[2] == 0:
            return str_condsendInt64.format(cond_var, sent_var,
                                            self.pid, ctx.processid().getText())
        if senttype[1] == "float32" and senttype[2] == 0:
            return str_condsendFloat32.format(cond_var, sent_var,
                                              self.pid, ctx.processid().getText())
        if senttype[1] == "float64" and senttype[2] == 0:
            return str_condsendFloat64.format(cond_var, sent_var,
                                              self.pid, ctx.processid().getText())

    def visitSend(self, ctx):
        str_sendInt = "parallely.SendInt({}, {}, {});\n"
        str_sendInt32 = "parallely.SendInt32({}, {}, {});\n"
        str_sendInt64 = "parallely.SendInt64({}, {}, {});\n"
        str_sendIntArray = "parallely.SendIntArray({}, {}, {});\n"
        str_sendInt32Array = "parallely.SendInt32Array({}[:], {}, {});\n"
        str_sendInt64Array = "parallely.SendInt64Array({}[:], {}, {});\n"

        str_sendFloat64Array = "parallely.SendFloat64Array({}[:], {}, {});\n"
        str_sendFloat32Array = "parallely.SendFloat32Array({}[:], {}, {});\n"
        str_sendFloat64 = "parallely.SendFloat64({}, {}, {});\n"
        str_sendFloat32 = "parallely.SendFloat32({}, {}, {});\n"

        sent_var = ctx.var().getText()
        senttype = self.getType(ctx.fulltype())

        if senttype[1] == "int" and senttype[2] == 0:
            return str_sendInt.format(sent_var, self.pid, ctx.processid().getText())
        if senttype[1] == "int32" and senttype[2] == 0:
            return str_sendInt32.format(sent_var, self.pid, ctx.processid().getText())
        if senttype[1] == "int64" and senttype[2] == 0:
            return str_sendInt64.format(sent_var, self.pid, ctx.processid().getText())

        if (senttype[1] == "int" or senttype[1] == "int") and senttype[2] == 1:
            return str_sendIntArray.format(sent_var, self.pid, ctx.processid().getText())
        if (senttype[1] == "int32" or senttype[1] == "int") and senttype[2] == 1:
            return str_sendInt32Array.format(sent_var, self.pid, ctx.processid().getText())
        if (senttype[1] == "int64" or senttype[1] == "int") and senttype[2] == 1:
            return str_sendInt64Array.format(sent_var, self.pid, ctx.processid().getText())

        if senttype[1] == "float64" and senttype[2] == 1:
            return str_sendFloat64Array.format(sent_var, self.pid, ctx.processid().getText())
        if senttype[1] == "float32" and senttype[2] == 1:
            return str_sendFloat32Array.format(sent_var, self.pid, ctx.processid().getText())
        if senttype[1] == "float64" and senttype[2] == 0:
            return str_sendFloat64.format(sent_var, self.pid, ctx.processid().getText())
        if senttype[1] == "float32" and senttype[2] == 0:
            return str_sendFloat32.format(sent_var, self.pid, ctx.processid().getText())

    def visitReceive(self, ctx):
        str_RecInt = "parallely.ReceiveInt(&{}, {}, {});\n"
        str_RecInt32 = "parallely.ReceiveInt32(&{}, {}, {});\n"
        str_RecInt64 = "parallely.ReceiveInt64(&{}, {}, {});\n"

        str_RecIntArray = "parallely.ReceiveIntArray({}[:], {}, {});\n"
        str_RecInt32Array = "parallely.ReceiveInt32Array({}[:], {}, {});\n"
        str_RecInt64Array = "parallely.ReceiveInt64Array({}[:], {}, {});\n"

        str_RecFloat64 = "parallely.ReceiveFloat64(&{}, {}, {});\n"
        str_RecFloat32 = "parallely.ReceiveFloat32(&{}, {}, {});\n"

        str_RecFloat64Array = "parallely.ReceiveFloat64Array({}[:], {}, {});\n"
        str_RecFloat32Array = "parallely.ReceiveFloat32Array({}[:], {}, {});\n"

        # parallely.Condreceive(&b, &n, 0, 1);
        senttype = self.getType(ctx.fulltype())
        if senttype[1] == "int" and senttype[2] == 0:
            return str_RecInt.format(ctx.var().getText(),
                                     self.pid, ctx.processid().getText())
        if senttype[1] == "int32" and senttype[2] == 0:
            return str_RecInt32.format(ctx.var().getText(),
                                       self.pid, ctx.processid().getText())
        if senttype[1] == "int64" and senttype[2] == 0:
            return str_RecInt64.format(ctx.var().getText(),
                                       self.pid, ctx.processid().getText())
        if senttype[1] == "int" and senttype[2] == 1:
            return str_RecIntArray.format(ctx.var().getText(),
                                          self.pid, ctx.processid().getText())
        if senttype[1] == "int32" and senttype[2] == 1:
            return str_RecInt32Array.format(ctx.var().getText(),
                                            self.pid, ctx.processid().getText())
        if senttype[1] == "int64" and senttype[2] == 1:
            return str_RecInt64Array.format(ctx.var().getText(),
                                            self.pid, ctx.processid().getText())
        if senttype[1] == "float64" and senttype[2] == 0:
            return str_RecFloat64.format(ctx.var().getText(),
                                         self.pid, ctx.processid().getText())
        if senttype[1] == "float32" and senttype[2] == 0:
            return str_RecFloat32.format(ctx.var().getText(),
                                         self.pid, ctx.processid().getText())
        if senttype[1] == "float64" and senttype[2] == 1:
            return str_RecFloat64Array.format(ctx.var().getText(),
                                              self.pid, ctx.processid().getText())
        if senttype[1] == "float32" and senttype[2] == 1:
            return str_RecFloat32Array.format(ctx.var().getText(),
                                              self.pid, ctx.processid().getText())

    def visitCondreceive(self, ctx):
        str_condreceive = "parallely.Condreceive(&{}, &{}, {}, {});\n"
        str_condreceiveIntArray = "parallely.CondreceiveIntArray(&{}, {}[:], {}, {});\n"
        str_condreceiveInt32 = "parallely.CondreceiveInt32(&{}, &{}, {}, {});\n"
        str_condreceiveInt64 = "parallely.CondreceiveInt64(&{}, &{}, {}, {});\n"
        str_condreceiveFloat32 = "parallely.CondreceiveFloat32(&{}, &{}, {}, {});\n"
        str_condreceiveFloat64 = "parallely.CondreceiveFloat64(&{}, &{}, {}, {});\n"

        senttype = self.getType(ctx.fulltype())
        if senttype[1] == "int" and senttype[2] == 0:
            return str_condreceive.format(ctx.var(0).getText(), ctx.var(1).getText(),
                                          self.pid, ctx.processid().getText())
        if senttype[1] == "int" and senttype[2] == 1:
            return str_condreceiveIntArray.format(ctx.var(0).getText(), ctx.var(1).getText(),
                                                  self.pid, ctx.processid().getText())
        if senttype[1] == "int32" and senttype[2] == 0:
            return str_condreceiveInt32.format(ctx.var(0).getText(), ctx.var(1).getText(),
                                               self.pid, ctx.processid().getText())
        if senttype[1] == "int64" and senttype[2] == 0:
            return str_condreceiveInt64.format(ctx.var(0).getText(), ctx.var(1).getText(),
                                               self.pid, ctx.processid().getText())
        if senttype[1] == "float32" and senttype[2] == 0:
            return str_condreceiveFloat32.format(ctx.var(0).getText(), ctx.var(1).getText(),
                                                 self.pid, ctx.processid().getText())
        if senttype[1] == "float64" and senttype[2] == 0:
            return str_condreceiveFloat64.format(ctx.var(0).getText(), ctx.var(1).getText(),
                                                 self.pid, ctx.processid().getText())
        print "[Error] No condrec : ", senttype

    def visitProbassignment(self, ctx):
        assigned_var = ctx.var().getText()
        prob = ctx.probability().getText()
        return str_probchoiceInt.format(assigned_var, prob, ctx.expression(0).getText(),
                                        ctx.expression(1).getText())

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
        for proc in self.proc_groups[group_name]:
            self.pid = proc
            print "Translating process group: ", self.pid

            # Collect the declarations which should be at the top
            # Binding the pid int to the process name.
            dec_string = "{} := {};\n".format(element_name, proc)
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

            process_name = "func_" + self.pid
            self.process_list.append(process_name)

            process_code = dec_string + statement_string

            process_def_str = str_single_thread.format(process_name, process_code, element_name)
            self.process_defs.append(process_def_str)
            # print "--------------------"
            # print process_def_str
            # print "--------------------"

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
            translated = self.visit(statement)
            if translated:
                statement_string += translated
            else:
                print "[Error] Unable to transtate: ", statement.getText()

        process_name = "func_" + self.pid
        self.process_list.append(process_name)

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
        run_procs = ''.join(["go {}();\n".format(fname) for fname in self.process_list])
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
                        help="File containing the code")
    parser.add_argument("-o", dest="outfile",
                        help="File to output the sequential code")
    parser.add_argument("-d", "--debug", action="store_true",
                        help="Print debug info")
    parser.add_argument("-t", dest="template",
                        help="File containing the template")
    args = parser.parse_args()

    programfile = open(args.programfile, 'r')
    # outfile = open(args.outfile, 'w')
    program_str = programfile.read()
    main(program_str, args.outfile, programfile.name, args.template, args.debug)
