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

# str_single_thread = '''func {}() {{
#   defer parallely.Wg.Done()
#   _temp_index := -1;
#   _ = _temp_index;
#   var DynMap = map[parallely.DynKey] float64{{}};
#   _ = DynMap;
#   {}
#   fmt.Println("Ending thread : ", {});
# }}'''

# str_member_thread = '''func {}(tid int) {{
#   defer parallely.Wg.Done()
#   _temp_index := -1;
#   _ = _temp_index;
#   var DynMap = map[parallely.DynKey] float64{{}};
#   _ = DynMap;
#   {}
#   fmt.Println("Ending thread : ", {});
# }}'''

# str_probchoiceIntFlag = "{} = parallely.RandchoiceFlag(float32({}), {}, {}, &__flag_{});\n"
# str_probchoiceInt = "{} = parallely.Randchoice(float32({}), {}, {});\n"

# dyn_rec_str = '''my_chan_index := {} * parallely.Numprocesses + {};
# __temp_rec_val := <- parallely.DynamicChannelMap[my_chan_index];
# DynMap[parallely.DynKey{{Varname: \"{}\", Index: 0}}] = __temp_rec_val;
# '''

# ch_str = '''
# fmt.Println("----------------------------");\n
# fmt.Println("Spec checkarray({0}, {1}): ", parallely.CheckArray(\"{0}\", {1}, DynMap));\n
# fmt.Println("----------------------------");\n
# '''

# dyn_pchoice_str = '''
# DynMap[parallely.DynKey{{Varname: \"{}\", Index: 0}}] = parallely.Max(0.0, {} - float64({})) * {};
# '''

# dyn_assign_str = '''
# DynMap[parallely.DynKey{{Varname: \"{}\", Index: 0}}] = parallely.Max(0.0, {} - float64({}));
# '''

# dyn_precise = '''
# DynMap[parallely.DynKey{{Varname: \"{}\", Index: 0}}] = 1;
# '''

str_single_thread = '''func {}() {{
  defer parallely.Wg.Done()
  _temp_index := -1;
  _ = _temp_index;
  var DynMap = map[int] float64{{}};
  _ = DynMap;
  {}
  fmt.Println("Ending thread : ", {});
}}'''

str_member_thread = '''func {}(tid int) {{
  defer parallely.Wg.Done()
  _temp_index := -1;
  _ = _temp_index;
  var DynMap = map[int] float64{{}};
  _ = DynMap;
  {}
  fmt.Println("Ending thread : ", {});
}}'''

str_probchoiceIntFlag = "{} = parallely.RandchoiceFlag(float32({}), {}, {}, &__flag_{});\n"
str_probchoiceInt = "{} = parallely.Randchoice(float32({}), {}, {});\n"

dyn_rec_str = '''my_chan_index := {} * parallely.Numprocesses + {};
__temp_rec_val := <- parallely.DynamicChannelMap[my_chan_index];
DynMap[{}] = __temp_rec_val;
'''

ch_str = '''
fmt.Println("----------------------------");\n
fmt.Println("Spec checkarray({0}, {1}): ", parallely.CheckArray(\"{0}\", {1}, DynMap));\n
fmt.Println("----------------------------");\n
'''

dyn_pchoice_str = '''
DynMap[{}] = parallely.Max(0.0, {} - float64({})) * {};
'''

dyn_assign_str = '''
DynMap[{}] = parallely.Max(0.0, {} - float64({}));
'''

dyn_precise = '''
DynMap[{}] = 1;
'''

t_d_str = '''
if {0} != 0 {{
    DynMap[{1}]  = DynMap[{2}] * DynMap[{4}]
}} else {{
    DynMap[{1}] = DynMap[{3}] * DynMap[{4}]
}};\n'''


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
    def __init__(self, dynamic, args):
        print "Starting translation"
        self.pid = None
        self.process_defs = []
        self.process_list = []
        self.globaldecs = []
        self.tempvarnum = 0
        self.recovernum = 0
        self.primitiveTMap = {}
        self.typeMap = {}
        self.arrays = []
        self.enableDynamic = dynamic
        self.args = args
        self.varMap = {}
        self.varNum = 0
        # Only support 1D arrays
        self.arraySize = {}

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

    # # We will ignore global variables. Assume that they are in the scaffolding
    # def visitGlobalarray(self, ctx):
    #     str_global_dec = "var {} []{};\n"
    #     mytype = self.getType(ctx.basictype())
    #     varname = ctx.GLOBALVAR().getText()
    #     # Q = {2,3,4,5};
    #     global_str = str_global_dec.format(varname, mytype[1])
    #     self.globaldecs.append(global_str)

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
            ("int", 1): "parallely.SendIntArray({}[:], {}, {});\n",
            ("int32", 1): "parallely.SendInt32Array({}[:], {}, {});\n",
            ("int64", 1): "parallely.SendInt64Array({}[:], {}, {});\n",
            ("float32", 1): "parallely.SendFloat32Array({}[:], {}, {});\n",
            ("float64", 1): "parallely.SendFloat64Array({}[:], {}, {});\n"
        }

        dyn_send_str = {
            ("int", 0): "parallely.SendDynIntArray({}[:], {}, {}, DynMap, {});\n",
            ("float64", 0): "parallely.SendDynFloat64Array({}[:], {}, {}, DynMap, {});\n",
            ("int", 1): "parallely.SendDynIntArrayO1({}[:], {}, {}, DynMap, {});\n",
            ("float64", 1): "parallely.SendDynFloat64ArrayO1({}[:], {}, {}, DynMap, {});\n",
        }

        sent_var = ctx.var().getText()
        senttype = self.getType(ctx.fulltype())

        # oplevel = 0
        # if self.args.arrayO1:
        #     oplevel = 1

        if (sent_var in self.arrays and self.enableDynamic and
                sent_var in self.primitiveTMap and self.primitiveTMap[sent_var] == 'dynamic'):
            t_str = dyn_send_str[senttype[1], self.args.arrayO1].format(sent_var, self.pid,
                                                                        ctx.processid().getText(),
                                                                        self.arraySize[sent_var])
            print t_str
            return t_str

        s_str_0 = send_str[senttype[1], senttype[2]].format(sent_var, self.pid, ctx.processid().getText())
        d_str = ""
        if self.enableDynamic and sent_var in self.primitiveTMap and self.primitiveTMap[sent_var] == 'dynamic':
            v_str = "DynMap[{}]".format(self.varMap[sent_var])
            d_str = "parallely.SendDynVal({}, {}, {});\n".format(v_str, self.pid, ctx.processid().getText())

        return s_str_0 + d_str

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

        dyn_rec_dict = {
            ("int", 0): "parallely.ReceiveDynIntArray({}[:], {}, {}, DynMap, {});\n",
            ("float64", 0): "parallely.ReceiveDynFloat64Array({}[:], {}, {}, DynMap, {});\n",
            ("int", 1): "parallely.ReceiveDynIntArrayO1({}[:], {}, {}, DynMap, {});\n",
            ("float64", 1): "parallely.ReceiveDynFloat64ArrayO1({}[:], {}, {}, DynMap, {});\n",
        }

        senttype = self.getType(ctx.fulltype())
        rec_var = ctx.var().getText()

        if (rec_var in self.arrays and self.enableDynamic and
                rec_var in self.primitiveTMap and self.primitiveTMap[rec_var] == 'dynamic'):
            return dyn_rec_dict[senttype[1], self.args.arrayO1].format(ctx.var().getText(),
                                                                       self.pid,
                                                                       ctx.processid().getText(),
                                                                       self.varMap[rec_var])

        rec_str_0 = rec_str[senttype[1], senttype[2]].format(ctx.var().getText(),
                                                             self.pid, ctx.processid().getText())
        d_str = ""
        if self.enableDynamic and rec_var in self.primitiveTMap and self.primitiveTMap[rec_var] == 'dynamic':
            d_str = dyn_rec_str.format(ctx.processid().getText(), self.pid, self.varMap[ctx.var().getText()])

        return rec_str_0 + d_str

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
        return rec_str[senttype[1], senttype[2]].format(ctx.var(0).getText(), ctx.var(1).getText(),
                                                        self.pid, ctx.processid().getText())

    def visitProbassignment(self, ctx):
        rand_str = {
            ("int"): "{} = parallely.Randchoice(float32({}), {}, {});\n",
            ("float64"): "{} = parallely.RandchoiceFloat64(float32({}), {}, {});\n",
            # ("int32", 0): "parallely.CondreceiveInt32(&{}, &{}, {}, {});\n",
            # ("int64", 0): "parallely.CondreceiveInt64(&{}, &{}, {}, {});\n",
            # ("float32", 0):
            # ("float64", 0):
            # ("int", 1): "parallely.CondreceiveIntArray(&{}, {}[:], {}, {});\n",
            # # ("int32", 1):
            # # ("int64", 1):
            # ("float32", 1): "parallely.CondreceiveFloat32(&{}, &{}, {}, {});\n",
            # ("float64", 1): "parallely.CondreceiveFloat64(&{}, &{}, {}, {});\n"
        }

        a_var = ctx.var().getText()
        prob = ctx.probability().getText()

        if (self.recovernum == 0):
            p_str = rand_str[self.typeMap[a_var][0]].format(a_var, prob, ctx.expression(0).getText(),
                                                            ctx.expression(1).getText())
        else:
            p_str = str_probchoiceIntFlag.format(a_var, prob, ctx.expression(0).getText(),
                                                 ctx.expression(1).getText(), self.recovernum)

        if self.enableDynamic and a_var in self.primitiveTMap and self.primitiveTMap[a_var] == 'dynamic':
            var_list = self.getVarList(ctx.precise)
            if len(var_list) == 0:
                dyn_str = "DynMap[{}] = {};\n".format(a_var, ctx.probability().getText())
            elif len(var_list) == 1:
                dyn_str = "DynMap[{}] = DynMap[{}] * {};\n".format(self.varMap[a_var],
                                                                   self.varMap[var_list[0]],
                                                                   ctx.probability().getText())
            else:
                sum_str = []
                for var in var_list:
                    sum_str.append("DynMap[{}]".format(self.varMap[var]))
                dyn_str = dyn_pchoice_str.format(self.varMap[a_var], " + ".join(sum_str),
                                                 len(var_list) - 1, ctx.probability().getText())
            p_str = p_str + dyn_str
        return p_str

    def handleExpression(self, ctx):
        convert_str = "parallely.ConvBool({})"
        if isinstance(ctx, ParallelyParser.SelectContext):
            return self.handleExpression(ctx.expression())
        if (isinstance(ctx, ParallelyParser.EqContext) or
                isinstance(ctx, ParallelyParser.GeqContext) or
                isinstance(ctx, ParallelyParser.LeqContext) or
                isinstance(ctx, ParallelyParser.LessContext) or
                isinstance(ctx, ParallelyParser.GreaterContext) or
                isinstance(ctx, ParallelyParser.AndexpContext)):
            # If it is a boolean statement
            return convert_str.format(ctx.getText())
        else:
            return ctx.getText()

    def getVarList(self, expression):
        if isinstance(expression, ParallelyParser.SelectContext):
            return self.getVarList(expression.expression())

        if isinstance(expression, ParallelyParser.VariableContext):
            # We will drop precise variables as their reliability is guarateed to be 1
            if self.primitiveTMap[expression.getText()] == 'precise':
                return []
            return [expression.getText().encode("ascii")]

        if (isinstance(expression, ParallelyParser.LiteralContext) or
                isinstance(expression, ParallelyParser.FliteralContext)):
            return []

        dyn_list = []
        for expression_part in expression.expression():
            partial_list = self.getVarList(expression_part)
            dyn_list = dyn_list + partial_list

        return dyn_list

    def visitCondassignment(self, ctx):
        assign_str = "if {0} != 0 {{ {1}  = {2} }} else {{ {1} = {3} }};\n"
        a_var = ctx.var()[0].getText()
        b_var = ctx.condition.getText()
        o1_var = ctx.ifvar.getText()
        o2_var = ctx.elsevar.getText()

        out_str = assign_str.format(b_var, a_var, o1_var, o2_var)
        d_str = ""
        if self.enableDynamic and a_var in self.primitiveTMap and self.primitiveTMap[a_var] == 'dynamic':
            # t_d_str = '''
            # if {0} != 0 {{
            #     DynMap[{1}]  = DynMap[{2}] * DynMap[{4}]
            # }} else {{
            #     DynMap[{1}] = DynMap[{3}] * DynMap[{4}]
            # }};\n'''
            # if
            d_str = t_d_str.format(b_var, self.varMap[b_var], self.varMap[a_var],
                                   self.varMap[o1_var], self.varMap[o2_var])

        return out_str + d_str

    def visitExpassignment(self, ctx):
        assign_str = "{} = {};\n"
        expr_str = self.handleExpression(ctx.expression())
        var_str = ctx.var().getText()

        # hack to handle array copy
        # print var_str, expr_str, self.arrays
        if self.enableDynamic and var_str in self.arrays and expr_str in self.arrays:
            if self.primitiveTMap[expr_str] == 'precise':
                dyn_c_str = "parallely.InitDynArray({}, {}, DynMap);\n".format(self.varMap[var_str],
                                                                               self.arraySize[var_str])
            else:
                dyn_c_str = "parallely.CopyDynArray({}, {}, {}, DynMap);\n".format(self.varMap[var_str],
                                                                                   self.varMap[expr_str],
                                                                                   self.arraySize[var_str])
            return ctx.getText() + ";\n" + dyn_c_str

        dyn_str = ""
        # Not global and dynamic
        if self.enableDynamic and var_str in self.primitiveTMap and self.primitiveTMap[var_str] == 'dynamic':
            var_list = self.getVarList(ctx.expression())
            # print ctx.getText(), var_list
            if len(var_list) == 0:
                dyn_str = dyn_precise.format(var_str)
            elif len(var_list) == 1:
                dyn_str = "DynMap[{} + _temp_index] = DynMap[{}];\n".format(self.varMap[var_str],
                                                                            self.varMap[var_list[0]])
            else:
                sum_str = []
                for var in var_list:
                    sum_str.append("DynMap[{}]".format(self.varMap[var]))
                dyn_str = dyn_assign_str.format(self.varMap[var_str], " + ".join(sum_str), len(var_list) - 1)
        return assign_str.format(var_str, expr_str) + dyn_str

    def visitGexpassignment(self, ctx):
        assign_str = "{} = {};\n"
        expr_str = self.handleExpression(ctx.expression())
        var_str = ctx.GLOBALVAR().getText()
        return assign_str.format(var_str, expr_str)

    # For now only supports 1d arrays
    def visitArrayload(self, ctx):
        assigned_var = ctx.var()[0].getText()
        go_str = "_temp_index = {};\n{}={}[_temp_index];\n"
        dyn_upd_map = "DynMap[{}] = DynMap[{} + _temp_index];\n"
        # print self.primitiveTMap, assigned_var
        if (self.enableDynamic and assigned_var in self.primitiveTMap and
                self.primitiveTMap[assigned_var] == 'dynamic'):
            index_expr = ctx.expression()[0].getText()
            array_var = ctx.var()[1].getText()
            # assigned_var = ctx.var()[0]
            return go_str.format(index_expr, assigned_var, array_var) + dyn_upd_map.format(self.varMap[assigned_var], self.varMap[array_var])
        return ctx.getText() + ";\n"

    def visitArrayassignment(self, ctx):
        # print ctx.getText()
        a_var = ctx.var().getText()
        go_str = "_temp_index = {};\n{}[_temp_index]={};\n"
        index_expr = ctx.expression()[0].getText()
        a_expr = ctx.expression()[1].getText()

        r_str = go_str.format(index_expr, a_var, a_expr)

        dyn_str = ""
        if self.enableDynamic and a_var in self.primitiveTMap and self.primitiveTMap[a_var] == 'dynamic':
            # DynMap[{}] = parallely.Max(0.0, {} - float64({}));
            dyn_upd_map = "DynMap[{} + _temp_index] = parallely.Max(0.0, {} - float64({}));\n"
            var_list = self.getVarList(ctx.expression()[1])
            if len(var_list) == 0:
                dyn_str = "DynMap[{} + _temp_index] = 1;\n".format(self.varMap[a_var])
            elif len(var_list) == 1:
                dyn_str = "DynMap[{} + _temp_index] = DynMap[{}];\n".format(self.varMap[a_var],
                                                                            self.varMap[var_list[0]])
            else:
                sum_str = []
                for var in var_list:
                    sum_str.append("DynMap[{}]".format(self.varMap[var]))
                dyn_str = dyn_upd_map.format(self.varMap[a_var], " + ".join(sum_str), len(var_list) - 1)
            # array_str = "[]string{" + ", ".join(var_list) + "}"
            # r_str += dyn_upd_map.format(a_var, var_list)

        # + dyn_upd_map.format(assigned_var, array_var)
        # go_str = "_temp_index := {};\n{}={}[_temp_index];\n"
        return r_str + dyn_str

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

        # how_deep = self.recovernum
        combine_str = ""
        self.recovernum -= 1

        # For each outer flag set it to the value of the inner flags
        for f in range(starting_level - 1):
            temp_str = ""
            oneflag = "__flag_{} = __flag_{} {};\n"
            for f in range(starting_level - 1):
               temp_str += "|| __flag_{}".format(starting_level + f)
            combine_str += oneflag.format(f + 1, f + 1, temp_str)

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

    def visitSpeccheckarray(self, ctx):
        checked_var = ctx.var().getText()
        checked_val = ctx.probability().getText()
        return ch_str.format(checked_var, checked_val)

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
            self.arrays.append(varname)
            self.primitiveTMap[varname] = dectype[0]
            self.typeMap[varname] = (dectype[1], dectype[2])

            # self.varNum += 1
            dim = ""
            if decl.INT():
                for dimention in decl.INT():
                    dim += "[{}]".format(dimention)
            else:
                dim += "[]"

            if self.enableDynamic and dectype[0] == "dynamic":
                self.varMap[varname] = self.varNum
                self.arraySize[varname] = decl.INT()[0]
                self.varNum += int(decl.INT()[0].getText())

                # Only works for 1 dimentional so far
                d_str = ""
                if self.enableDynamic:
                    d_str = "parallely.InitDynArray({}, {}, DynMap);\n".format(self.varMap[varname],
                                                                               decl.INT()[0])
                return str_array_dec.format(varname, dim, dectype[1]) + d_str
            else:
                return str_array_dec.format(varname, dim, dectype[1])
        if isinstance(decl, ParallelyParser.DynarraydeclarationContext):
            dyn_array_dec = "var {0} {1}{2};\n {0}=make({1}{2}, {3});\n"
            varname = decl.var().getText()
            self.arrays.append(varname)
            self.primitiveTMap[varname] = dectype[0]
            self.typeMap[varname] = (dectype[1], dectype[2])
            # self.varMap[varname] = self.varNum
            # self.varNum += 1
            # self.arraySize[varname] = decl.INT()[0]

            dim = []
            for dimention in decl.GLOBALVAR():
                dim.append(dimention)

            d_str = ""
            if self.enableDynamic and dectype[0] == "dynamic":
                # Only works for 1 dimentional so far
                print "Not supporting dynamic unbounded arrays"
                exit(-1)
                d_str = "parallely.InitDynArray(\"{}\", {}, DynMap);\n".format(varname,
                                                                               decl.GLOBALVAR()[0])

            if len(dim) == 1:
                return dyn_array_dec.format(varname, "[]", dectype[1], dim[0]) + d_str
            if len(dim) > 1:
                return dyn_array_dec.format(varname, "[]", dectype[1], dim[0]) + d_str
            else:
                print "[Error] Unable to transtate: ", decl.getText()
                exit(-1)
        else:
            varname = decl.var().getText()
            self.primitiveTMap[varname] = dectype[0]
            self.typeMap[varname] = (dectype[1], dectype[2])
            if self.enableDynamic and dectype[0] == "dynamic":
                self.varMap[varname] = self.varNum
                self.varNum += 1
                d_init_str = "var {0} {1};\nDynMap[{2}] = 1;\n"
                return d_init_str.format(varname, dectype[1], self.varMap[varname])
            else:
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
        # self.primitiveTMap = {}
        # self.typeMap = {}
        # self.varMap = {}
        # self.varNum = 0

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
                exit(-1)

        process_name = "func_" + self.pid
        self.process_list.append((process_name, 0))

        process_code = dec_string + statement_string

        process_def_str = str_single_thread.format(process_name, process_code, self.pid)
        # print "--------------------"
        # print process_def_str
        # print "--------------------"
        self.process_defs.append(process_def_str)

    def translate(self, tree, numthreads, proc_groups_in, fout_name, template):
        for gdec in tree.globaldec():
            if isinstance(gdec, ParallelyParser.GlobalarrayContext):
                self.varMap[gdec.GLOBALVAR().getText()] = self.varNum
                self.arraySize[gdec.GLOBALVAR().getText()] = gdec.INT()
                self.varNum += int(gdec.INT().getText())
                self.arrays.append(gdec.GLOBALVAR().getText())
                dectype = self.getType(gdec.basictype())
                self.primitiveTMap[gdec.GLOBALVAR().getText()] = dectype[0]
                # print "#########", gdec.GLOBALVAR().getText()

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


def main(program_str, outfile, filename, template, debug, dynamic, args):
    print "Starting the cross compilation"
    input_stream = InputStream(program_str)
    lexer = ParallelyLexer(input_stream)
    stream = CommonTokenStream(lexer)
    parser = ParallelyParser(stream)

    tree = parser.parallelprogram()

    threadcounter = CountThreads()
    threadcounter.visit(tree)

    print "Number of processes found: {}".format(threadcounter.processcount)

    translator = Translator(dynamic, args)
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
    parser.add_argument("-dyn", "--dynamic", action="store_true",
                        help="Enable dynamic tracking")
    parser.add_argument("-a", "--arrayO1", action="store_true",
                        help="Inline tracking")
    args = parser.parse_args()

    if args.dynamic:
        print "Enabling dynamic tracking"
    if args.arrayO1:
        print "Enabling array optimization: Send one value"

    programfile = open(args.programfile, 'r')
    # outfile = open(args.outfile, 'w')
    program_str = programfile.read()
    main(program_str, args.outfile, programfile.name, args.template,
         args.debug, args.dynamic, args)
