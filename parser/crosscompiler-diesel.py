from antlr4 import CommonTokenStream
from antlr4 import InputStream
from ParallelyLexer import ParallelyLexer
from ParallelyParser import ParallelyParser
from ParallelyVisitor import ParallelyVisitor
from argparse import ArgumentParser
import collections
import time

key_error_msg = "Type error detected: Undeclared variable (probably : {})"

str_single_thread = '''func {}() {{
  defer diesel.Wg.Done();
  var DynMap [{}]diesel.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  {}

  fmt.Println("Ending thread : ", {});
}}'''

str_member_thread = '''func {}(tid int) {{
  defer diesel.Wg.Done();
  var DynMap [{}]diesel.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  {}
  fmt.Println("Ending thread : ", {});
}}'''

str_probchoiceIntFlag = "{} = diesel.RandchoiceFlag(float32({}), {}, {}, &__flag_{});\n"
str_probchoiceInt = "{} = diesel.Randchoice(float32({}), {}, {});\n"

dyn_rec_str = '''my_chan_index = {0} * diesel.Numprocesses + {1};
__temp_rec_val_{3} := <- diesel.DynamicChannelMap[my_chan_index];
DynMap[{2}] = __temp_rec_val_{3};
'''

ch_str = '''
fmt.Println("----------------------------");\n
fmt.Println("Spec checkarray({3}, {1}): ", diesel.CheckArray({0}, {1}, {2}, DynMap[:]));\n
fmt.Println("----------------------------");\n
'''

dyn_pchoice_str = '''DynMap[{}].Reliability = diesel.Max(0.0, {} - float32({})) * {};
'''

dyn_assign_str = '''DynMap[{}].Reliability = {} - {}.0;
'''

dyn_precise = '''DynMap[{}] = diesel.ProbInterval{{1, 0}};\n'''

t_d_str = '''if temp_bool_{0} != 0 {{DynMap[{1}]  = DynMap[{2}] + DynMap[{4}] - 1.0}} else {{ DynMap[{1}] = DynMap[{3} ] + DynMap[{4}] - 1.0}};\n'''

t_d_str2 = '''if temp_bool_{0} != 0 {{
    {2}}} else {{
    {3}}};\n'''

DynUpdate = collections.namedtuple('DynUpdate', ['updated', 'sum', 'multiplicative', 'isarray'])
ConditionalDynUpdate = collections.namedtuple('ConditionalDynUpdate',
                                              ['condition', 'ifop', 'elseop', 'updated'])


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
        self.allTracking = []
        self.tracking = []
        self.tempindexnum = 0
        self.dynsize = 0

    def visitSingleglobaldec(self, ctx):
        str_global_dec = "var {} = []int {{{}}};\n"
        varname = ctx.GLOBALVAR().getText()
        members = [t.getText() for t in ctx.processid()]
        # Q = {2,3,4,5};
        global_str = str_global_dec.format(varname, ','.join(members))
        self.globaldecs.append(global_str)

    # def visitGlobalconst(self, ctx):
    #     str_global_dec = "var {} {};\n"
    #     mytype = self.getType(ctx.basictype())
    #     varname = ctx.GLOBALVAR().getText()
    #     # Q = {2,3,4,5};
    #     global_str = str_global_dec.format(varname, mytype[1])
    #     self.globaldecs.append(global_str)

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
            ("int", 0): "diesel.Condsend({}, {}, {}, {});\n",
            ("int32", 0): "diesel.CondsendInt32({}, {}, {}, {});\n",
            ("int64", 0): "diesel.CondsendInt64({}, {}, {}, {});\n",
            ("float32", 0): "diesel.CondsendFloat32({}, {}, {}, {});\n",
            ("float64", 0): "diesel.CondsendFloat64({}, {}, {}, {});\n",
            ("int", 1): "diesel.CondsendIntArray({}, {}[:], {}, {});\n",
            ("int", 1): "diesel.CondsendIntArray({}, {}[:], {}, {});\n",
        }

        cond_var = ctx.var(0).getText()
        sent_var = ctx.var(1).getText()
        senttype = self.getType(ctx.fulltype())
        return cond_send_str[(senttype[1], senttype[2])].format(cond_var, sent_var,
                                                                self.pid, ctx.processid().getText())

    def visitSend(self, ctx):
        send_str = {
            ("int", 0): "diesel.SendInt({}, {}, {});\n",
            ("int32", 0): "diesel.SendInt32({}, {}, {});\n",
            ("int64", 0): "diesel.SendInt64({}, {}, {});\n",
            ("float32", 0): "diesel.SendFloat32({}, {}, {});\n",
            ("float64", 0): "diesel.SendFloat64({}, {}, {});\n",
            ("int", 1): "diesel.SendIntArray({}[:], {}, {});\n",
            ("int32", 1): "diesel.SendInt32Array({}[:], {}, {});\n",
            ("int64", 1): "diesel.SendInt64Array({}[:], {}, {});\n",
            ("float32", 1): "diesel.SendFloat32Array({}[:], {}, {});\n",
            ("float64", 1): "diesel.SendFloat64Array({}[:], {}, {});\n"
        }

        dyn_send_str = {
            ("int", 0): "diesel.SendDynIntArray({}[:], {}, {}, DynMap[:], {});\n",
            ("float64", 0): "diesel.SendDynFloat64Array({}[:], {}, {}, DynMap[:], {});\n",
            ("float32", 0): "diesel.SendDynFloat32Array({}[:], {}, {}, DynMap[:], {});\n",
            ("int", 1): "diesel.SendDynIntArrayO1({}[:], {}, {}, DynMap[:], {});\n",
            ("float64", 1): "diesel.SendDynFloat64ArrayO1({}[:], {}, {}, DynMap[:], {});\n",
            ("float32", 1): "diesel.SendDynFloat32ArrayO1({}[:], {}, {}, DynMap[:], {});\n",
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
                                                                        self.varMap[sent_var])
            # print t_str
            return t_str

        s_str_0 = send_str[senttype[1], senttype[2]].format(sent_var, self.pid, ctx.processid().getText())
        d_str = ""
        if self.enableDynamic and sent_var in self.primitiveTMap and self.primitiveTMap[sent_var] == 'dynamic':
            v_str = "DynMap[{}]".format(self.varMap[sent_var])
            d_str = "diesel.SendDynVal({}, {}, {});\n".format(v_str, self.pid, ctx.processid().getText())

        if self.args.gather:
            dyn_str = ''.join(self.trackingStatements)
            d_str = dyn_str + d_str
            d_str_opt = self.getDynString()
            # print d_str_opt
            d_str = d_str_opt + d_str

        self.trackingStatements = []
        self.tracking = []
        return s_str_0 + d_str

    def visitReceive(self, ctx):
        rec_str = {
            ("int", 0): "diesel.ReceiveInt(&{}, {}, {});\n",
            ("int32", 0): "diesel.ReceiveInt32(&{}, {}, {});\n",
            ("int64", 0): "diesel.ReceiveInt64(&{}, {}, {});\n",
            ("float32", 0): "diesel.ReceiveFloat32(&{}, {}, {});\n",
            ("float64", 0): "diesel.ReceiveFloat64(&{}, {}, {});\n",
            ("int", 1): "diesel.ReceiveIntArray({}[:], {}, {});\n",
            ("int32", 1): "diesel.ReceiveInt32Array({}[:], {}, {});\n",
            ("int64", 1): "diesel.ReceiveInt64Array({}[:], {}, {});\n",
            ("float32", 1): "diesel.ReceiveFloat32Array({}[:], {}, {});\n",
            ("float64", 1): "diesel.ReceiveFloat64Array({}[:], {}, {});\n"
        }

        dyn_rec_dict = {
            ("int", 0): "diesel.ReceiveDynIntArray({}[:], {}, {}, DynMap[:], {});\n",
            ("float64", 0): "diesel.ReceiveDynFloat64Array({}[:], {}, {}, DynMap[:], {});\n",
            ("float32", 0): "diesel.ReceiveDynFloat32Array({}[:], {}, {}, DynMap[:], {});\n",
            ("int", 1): "diesel.ReceiveDynIntArrayO1({}[:], {}, {}, DynMap[:], {});\n",
            ("float64", 1): "diesel.ReceiveDynFloat64ArrayO1({}[:], {}, {}, DynMap[:], {});\n",
            ("float32", 1): "diesel.ReceiveDynFloat32ArrayO1({}[:], {}, {}, DynMap[:], {});\n",
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
            self.tempvarnum += 1
            d_str = dyn_rec_str.format(ctx.processid().getText(),
                                       self.pid, self.varMap[ctx.var().getText()], self.tempvarnum)

        return rec_str_0 + d_str

    def visitCondreceive(self, ctx):
        rec_str = {
            ("int", 0): "diesel.Condreceive(&{}, &{}, {}, {});\n",
            ("int32", 0): "diesel.CondreceiveInt32(&{}, &{}, {}, {});\n",
            ("int64", 0): "diesel.CondreceiveInt64(&{}, &{}, {}, {});\n",
            # ("float32", 0):
            # ("float64", 0):
            ("int", 1): "diesel.CondreceiveIntArray(&{}, {}[:], {}, {});\n",
            # ("int32", 1):
            # ("int64", 1):
            ("float32", 1): "diesel.CondreceiveFloat32(&{}, &{}, {}, {});\n",
            ("float64", 1): "diesel.CondreceiveFloat64(&{}, &{}, {}, {});\n"
        }
        senttype = self.getType(ctx.fulltype())
        return rec_str[senttype[1], senttype[2]].format(ctx.var(0).getText(), ctx.var(1).getText(),
                                                        self.pid, ctx.processid().getText())

    def visitProbassignment(self, ctx):
        rand_str = {
            ("int"): "{} = diesel.Randchoice(float32({}), {}, {});\n",
            ("float64"): "{} = diesel.RandchoiceFloat64(float32({}), {}, {});\n",
            # ("int32", 0): "diesel.CondreceiveInt32(&{}, &{}, {}, {});\n",
            # ("int64", 0): "diesel.CondreceiveInt64(&{}, &{}, {}, {});\n",
            # ("float32", 0):
            # ("float64", 0):
            # ("int", 1): "diesel.CondreceiveIntArray(&{}, {}[:], {}, {});\n",
            # # ("int32", 1):
            # # ("int64", 1):
            # ("float32", 1): "diesel.CondreceiveFloat32(&{}, &{}, {}, {});\n",
            # ("float64", 1): "diesel.CondreceiveFloat64(&{}, &{}, {}, {});\n"
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
            var_list = list(set(self.getVarList(ctx.precise)))
            if len(var_list) == 0:
                dyn_str = "DynMap[{}] = diesel.ProbInterval{{{}, 0}};\n".format(self.varMap[a_var], ctx.probability().getText())
            elif len(var_list) == 1:
                dyn_str = "DynMap[{}].Reliability = DynMap[{}].Reliability * {};\n".format(self.varMap[a_var],
                                                                                           self.varMap[var_list[0]],
                                                                                           ctx.probability().getText())
            else:
                sum_str = []
                for var in var_list:
                    sum_str.append("DynMap[{}].Reliability".format(self.varMap[var]))
                dyn_str = dyn_pchoice_str.format(self.varMap[a_var], " + ".join(sum_str),
                                                 len(var_list) - 1, ctx.probability().getText())
            self.trackingStatements.append("// " + dyn_str)
            self.allTracking.append(self.getDynUpdate(ctx.precise, ctx.probability().getText(), a_var, 0))
            self.tracking.append(self.getDynUpdate(ctx.precise, ctx.probability().getText(), a_var, 0))
            if not self.args.gather:
                p_str = p_str + dyn_str
        return p_str

    def getDynUpdate(self, expression, probability, updated_var, isarray):
        var_list = list(set(self.getVarList(expression)))
        temp_dyn = DynUpdate(self.varMap[updated_var], [(self.varMap[var], 0) for var in var_list],
                             float(probability), isarray)
        return temp_dyn

    def handleExpression(self, ctx):
        convert_str = "diesel.ConvBool({})"
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

        return list(dyn_list)

    def visitCondassignment(self, ctx):
        assign_str = "temp_bool_{4}:= {0}; if temp_bool_{4} != 0 {{ {1}  = {2} }} else {{ {1} = {3} }};\n"
        a_var = ctx.var()[0].getText()
        b_var = ctx.condition.getText()
        o1_var = ctx.ifvar.getText()
        o2_var = ctx.elsevar.getText()

        self.tempindexnum += 1
        out_str = assign_str.format(b_var, a_var, o1_var, o2_var, self.tempindexnum)
        d_str = ""

        t_d_str = '''if temp_bool_{0} != 0 {{DynMap[{1}]  = DynMap[{2}] + DynMap[{3}] - 1.0}} else {{ DynMap[{1}] = DynMap[{2}] + DynMap[{4}] - 1.0}};\n'''

        if self.enableDynamic and a_var in self.primitiveTMap and self.primitiveTMap[a_var] == 'dynamic':
            d_str = t_d_str.format(self.tempindexnum, self.varMap[a_var], self.varMap[b_var],
                                   self.varMap[o1_var], self.varMap[o2_var])
            temp_cupd = ConditionalDynUpdate(self.tempindexnum,
                                             DynUpdate(self.varMap[a_var],
                                                       [(self.varMap[o1_var], 0),
                                                        (self.varMap[b_var], 0)], 1, 0),
                                             DynUpdate(self.varMap[a_var],
                                                       [(self.varMap[o2_var], 0),
                                                        (self.varMap[b_var], 0)], 1, 0),
                                             self.varMap[a_var])
            # print "*******: ", temp_cupd.updated, out_str, d_str
            self.trackingStatements.append("// " + d_str)
            self.tracking.append(temp_cupd)
            self.allTracking.append(temp_cupd)

        if not self.args.gather:
            return out_str + d_str
        return out_str

    def visitExpassignment(self, ctx):
        assign_str = "{} = {};\n"
        expr_str = self.handleExpression(ctx.expression())
        var_str = ctx.var().getText()

        # hack to handle array copy
        # print var_str, expr_str, self.arrays
        if (self.enableDynamic and var_str in self.arrays and
                expr_str in self.arrays and self.primitiveTMap[var_str] == 'dynamic'):
            if self.primitiveTMap[expr_str] == 'precise':
                dyn_c_str = "diesel.InitDynArray({}, {}, DynMap[:]);\n".format(self.varMap[var_str],
                                                                               self.arraySize[var_str])
            else:
                dyn_c_str = "diesel.CopyDynArray({}, {}, {}, DynMap[:]);\n".format(self.varMap[var_str],
                                                                                   self.varMap[expr_str],
                                                                                   self.arraySize[var_str])
            return ctx.getText() + ";\n" + dyn_c_str

        dyn_str = ""
        acc_str = ""
        # Not global and dynamic
        if self.enableDynamic and var_str in self.primitiveTMap and self.primitiveTMap[var_str] == 'dynamic':
            var_list = list(set(self.getVarList(ctx.expression())))

            if len(var_list) == 0:
                dyn_str = dyn_precise.format(self.varMap[var_str])
            elif len(var_list) == 1:
                dyn_str = "DynMap[{}].Reliability = DynMap[{}].Reliability;\n".format(self.varMap[var_str],
                                                                                      self.varMap[var_list[0]])
            else:
                sum_str = []
                for var in var_list:
                    sum_str.append("DynMap[{}].Reliability".format(self.varMap[var]))
                dyn_str = dyn_assign_str.format(self.varMap[var_str], " + ".join(sum_str), len(var_list) - 1)

            acc_str = self.getAccuracyStr(ctx.expression(), var_str)
            print acc_str

            self.trackingStatements.append("// " + dyn_str)
            self.allTracking.append(self.getDynUpdate(ctx.expression(), 1, var_str, 0))
            self.tracking.append(self.getDynUpdate(ctx.expression(), 1, var_str, 0))
        if not self.args.gather:
            return  dyn_str + acc_str + assign_str.format(var_str, expr_str)
        return assign_str.format(var_str, expr_str)

    def getAccuracyStr(self, ctx, var_str):
        if isinstance(ctx, ParallelyParser.FliteralContext) or isinstance(ctx, ParallelyParser.LiteralContext):
            return "" # "DynMap[{}].Delta = 0;\n".format(self.varMap[var_str])
        if isinstance(ctx, ParallelyParser.EqContext):
            return ""
        # We need to fix this
        if isinstance(ctx, ParallelyParser.GreaterContext):
            return ""
        if isinstance(ctx, ParallelyParser.SelectContext):
            return self.getAccuracyStr(ctx.expression(), var_str)
        if isinstance(ctx, ParallelyParser.VariableContext):
            return "DynMap[{}].Delta = DynMap[{}].Delta;\n".format(self.varMap[var_str],
                                                                   self.varMap[ctx.getText()])
        if isinstance(ctx, ParallelyParser.VarContext):
            return "DynMap[{}].Delta = DynMap[{}].Delta;\n".format(self.varMap[var_str],
                                                                   self.varMap[ctx.getText()])
        if isinstance(ctx, ParallelyParser.AddContext) or isinstance(ctx, ParallelyParser.MinusContext):
            var_list = self.getVarList(ctx)
            if len(var_list) == 0:
                return ""
            if len(var_list) == 1:
                return "DynMap[{}].Delta = DynMap[{}].Delta;\n".format(self.varMap[var_str],
                                                                       self.varMap[var_list[0]])
            elif len(var_list) == 2:
                d_str = "DynMap[{}].Delta = DynMap[{}].Delta + DynMap[{}].Delta;\n"
                return d_str.format(self.varMap[var_str],
                                    self.varMap[var_list[0]],
                                    self.varMap[var_list[1]])
            else:
                print("[ERROR]: Only support simple expressions: ", ctx.getText(), var_list)
                exit(-1)
        if isinstance(ctx, ParallelyParser.MultiplyContext):
            var_list = self.getVarList(ctx)
            upd_str = "DynMap[{0}].Delta = math.Abs({1}) * DynMap[{2}].Delta;\n"
            if len(var_list) == 1:
                # print ctx.getText(), var_list,ctx.expression(0).getText(), self.primitiveTMap
                if (isinstance(ctx.expression(0), ParallelyParser.FliteralContext) or
                    isinstance(ctx.expression(0), ParallelyParser.LiteralContext) or
                    (ctx.expression(0).getText() in self.primitiveTMap and self.primitiveTMap[ctx.expression(0).getText()] == 'dynamic')):
                    return upd_str.format(self.varMap[var_str],
                                          ctx.expression(1).getText(), self.varMap[var_list[0]])

                elif (isinstance(ctx.expression(1), ParallelyParser.FliteralContext) or
                      isinstance(ctx.expression(1), ParallelyParser.LiteralContext) or
                      (ctx.expression(1).getText() in self.primitiveTMap and self.primitiveTMap[ctx.expression(1).getText()] == 'dynamic')):
                    return upd_str.format(self.varMap[var_str],
                                          ctx.expression(0).getText(), self.varMap[var_list[0]])
                else:
                    return upd_str.format(var_str, var_list[0], self.varMap[var_list[0]],
                                      var_list[1], self.varMap[var_list[1]])
            elif len(var_list) == 2:
                upd_str = "DynMap[{0}].Delta = math.Abs({1}) * DynMap[{2}].Delta + math.Abs({3}) * DynMap[{4}].Delta + DynMap[{2}].Delta*DynMap[{4}].Delta;\n"
                return upd_str.format(self.varMap[var_str], var_list[0], self.varMap[var_list[0]],
                                      var_list[1], self.varMap[var_list[1]])
        if isinstance(ctx, ParallelyParser.DivideContext):
            var_list = self.getVarList(ctx)
            #implement the zero check at some point
            if len(var_list) == 1:
                if (isinstance(ctx.expression(0), ParallelyParser.FliteralContext) or
                    isinstance(ctx.expression(0), ParallelyParser.LiteralContext) or
                    (ctx.expression(0).getText() in self.primitiveTMap and self.primitiveTMap[ctx.expression(0).getText()] == 'dynamic')):
                    # print ":::::::::::", ctx.getText(), var_list
                    upd_str = "DynMap[{0}].Delta =  DynMap[{2}].Delta / math.Abs({1});\n"
                    return upd_str.format(self.varMap[var_str],
                                          ctx.expression(1).getText(), self.varMap[var_list[0]])

                elif (isinstance(ctx.expression(1), ParallelyParser.FliteralContext) or
                      isinstance(ctx.expression(1), ParallelyParser.LiteralContext) or
                (ctx.expression(1).getText() in self.primitiveTMap and self.primitiveTMap[ctx.expression(1).getText()] == 'dynamic')):
                    upd_str = "DynMap[{0}].Delta =  DynMap[{2}].Delta * math.Abs({1});\n"
                    return upd_str.format(self.varMap[var_str],
                                          ctx.expression(0).getText(), self.varMap[var_list[0]])

                return upd_str.format(var_str, var_list[0], self.varMap[var_list[0]],
                                      var_list[1], self.varMap[var_list[1]])
            elif len(var_list) == 2:
                upd_str = "DynMap[{0}].Delta = math.Abs({1}) * DynMap[{2}].Delta + math.Abs({3}) * DynMap[{4}].Delta / (math.Abs({3}) * (math.Abs({1})-DynMap[{4}].Delta));\n"
                return upd_str.format(self.varMap[var_str], var_list[0], self.varMap[var_list[0]],
                                      var_list[1], self.varMap[var_list[1]])

        print ctx.getText(), type(ctx)
        exit(-1)

    def visitGexpassignment(self, ctx):
        assign_str = "{} = {};\n"
        expr_str = self.handleExpression(ctx.expression())
        var_str = ctx.GLOBALVAR().getText()
        return assign_str.format(var_str, expr_str)

    # For now only supports 1d arrays
    def visitArrayload(self, ctx):
        self.tempindexnum += 1
        assigned_var = ctx.var()[0].getText()
        go_str = "_temp_index_{3} := {0};\n{1}={2}[_temp_index_{3}];\n"
        dyn_upd_map = "DynMap[{}] = DynMap[{} + _temp_index_{}];\n"
        index_expr = ctx.expression()[0].getText()
        array_var = ctx.var()[1].getText()
        # print self.primitiveTMap, assigned_var
        if (self.enableDynamic and assigned_var in self.primitiveTMap and
                self.primitiveTMap[assigned_var] == 'dynamic'):
            # assigned_var = ctx.var()[0]
            d_str = dyn_upd_map.format(self.varMap[assigned_var], self.varMap[array_var], self.tempindexnum)
            self.trackingStatements.append("// " + d_str)

            temp_dyn = DynUpdate(self.varMap[assigned_var],
                                 [(self.varMap[array_var], '_temp_index_{}'.format(self.tempindexnum))],
                                 1, 0)
            self.tracking.append(temp_dyn)
            self.allTracking.append(temp_dyn)
            if not self.args.gather:
                return go_str.format(index_expr, assigned_var, array_var, self.tempindexnum) + d_str
            else:
                return go_str.format(index_expr, assigned_var, array_var, self.tempindexnum)
        return go_str.format(index_expr, assigned_var, array_var, self.tempindexnum)

    def visitArrayassignment(self, ctx):
        # print ctx.getText()
        self.tempindexnum += 1
        a_var = ctx.var().getText()
        go_str = "_temp_index_{3} := {0};\n{1}[_temp_index_{3}]={2};\n"
        index_expr = ctx.expression()[0].getText()
        a_expr = ctx.expression()[1].getText()

        r_str = go_str.format(index_expr, a_var, a_expr, self.tempindexnum)

        dyn_str = ""
        if self.enableDynamic and a_var in self.primitiveTMap and self.primitiveTMap[a_var] == 'dynamic':
            # DynMap[{}] = parallely.Max(0.0, {} - float64({}));
            dyn_upd_map = "DynMap[{0} + _temp_index_{3}] = diesel.Max(0.0, {1} - float64({2}));\n"
            var_list = list(set(self.getVarList(ctx.expression()[1])))
            if len(var_list) == 0:
                dyn_str = "DynMap[{} + _temp_index_{}] = diesel.ProbInterval{{1, 0}};\n".format(self.varMap[a_var], self.tempindexnum)
            elif len(var_list) == 1:
                dyn_str = "DynMap[{0} + _temp_index_{2}] = DynMap[{1}];\n".format(self.varMap[a_var],
                                                                                  self.varMap[var_list[0]],
                                                                                  self.tempindexnum)
            else:
                sum_str = []
                for var in var_list:
                    sum_str.append("DynMap[{}]".format(self.varMap[var]))
                dyn_str = dyn_upd_map.format(self.varMap[a_var], " + ".join(sum_str),
                                             len(var_list) - 1, self.tempindexnum)
            self.trackingStatements.append("// " + dyn_str)

            temp_dyn = self.getDynUpdate(ctx.expression()[1], 1, a_var, self.tempindexnum)
            self.allTracking.append(temp_dyn)
            self.tracking.append(temp_dyn)

        if not self.args.gather:
            return r_str + dyn_str
        else:
            return r_str

    def visitCast(self, ctx):
        resultType = self.getType(ctx.fulltype())
        assignedvar = ctx.var(0).getText()
        castedvar = ctx.var(1).getText()
        # Array type
        if resultType[1] == "float64" and resultType[2] == 1:
            return "diesel.Cast32to64Array({}[:], {}[:]);\n".format(assignedvar,
                                                                    castedvar)
        if resultType[1] == "float32" and resultType[2] == 1:
            return "diesel.Cast64to32Array({}[:], {}[:]);\n".format(assignedvar,
                                                                    castedvar)
        # Regular cast
        if resultType[1] == "float64" and resultType[2] == 0:
            r_str = "{} = float64({});\n".format(assignedvar, castedvar)

            # Should reliability go to zero???
            d_str = ""
            if self.enableDynamic:
                d_str = "DynMap[{}] = DynMap[{}];\n".format(self.varMap[assignedvar],
                                                            self.varMap[castedvar])
            return r_str + d_str
        if resultType[1] == "float32" and resultType[2] == 0:
            d_str = ""
            if self.enableDynamic:
                d_str = "DynMap[{0}].Reliability = 0;\n DynMap[{0}].Delta = diesel.GetCastingError64to32({1}, {2});\n"
                d_str = d_str.format(self.varMap[assignedvar], castedvar, assignedvar)
            return "{} = float32({});\n".format(assignedvar, castedvar) + d_str

    def visitTrack(self, ctx):
        statement_string="{}={};\n".format(ctx.var(0).getText(), ctx.var(1).getText())
        if self.enableDynamic:
            updstr = "DynMap[{}] = diesel.ProbInterval{{{}, {}}};\n".format(self.varMap[ctx.var(0).getText()],
                                                                          ctx.probability().getText(),
                                                                          ctx.FLOAT().getText())
            self.tracking.append(updstr)
            self.allTracking.append(updstr)
            if not self.args.gather:
                return statement_string + updstr
            else:
                return statement_string
            
            
        return statement_string
        

    def visitIfonly(self, ctx):
        str_if_only = "if {} != 0 {{\n {} }}\n"
        cond_var = ctx.var().getText()

        statement_string = ''
        for statement in ctx.ifs:
            translated = self.visit(statement)
            if translated is not None:
                statement_string += translated
            else:
                print "[Error] Unable to translate: ", statement.getText()
                exit(-1)
        # print str_if_only.format(cond_var, statement_string)
        return str_if_only.format(cond_var, statement_string)

    def visitIf(self, ctx):
        str_if = "if {} != 0 {{\n {} }} else {{\n {} }}\n"
        cond_var = ctx.var().getText()

        temp_track_strings = list(self.trackingStatements)
        temp_track = list(self.tracking)

        statement_string = ''
        for statement in ctx.ifs:
            translated = self.visit(statement)
            if translated is not None:
                statement_string += translated
            else:
                print "[Error] Unable to translate: ", statement.getText()
                exit(-1)

        if self.args.gather:
            dyn_str = ''.join(self.trackingStatements)
            d_str_opt = self.getDynString()
            statement_string += (dyn_str + d_str_opt)

        self.trackingStatements = list(temp_track_strings)
        self.tracking = list(temp_track)

        else_statement_string = ''
        for statement in ctx.elses:
            translated = self.visit(statement)
            if translated is not None:
                else_statement_string += translated
            else:
                print "[Error] Unable to translate: ", statement.getText()
                exit(-1)

        if self.args.gather:
            dyn_str = ''.join(self.trackingStatements)
            d_str_opt = self.getDynString()
            else_statement_string += (dyn_str + d_str_opt)

        self.trackingStatements = []
        self.tracking = []
        # print str_if_only.format(cond_var, statement_string)
        return str_if.format(cond_var, statement_string, else_statement_string)

    def getReadDynString(self, member):
        if member[1] == 0:
            return "DynMap[{}]".format(member[0])
        else:
            return "DynMap[{}+{}]".format(member[0], member[1])

    def getUpdateString(self, var_list, probability):
        if len(var_list) == 0:
            return str(probability)
        elif len(var_list) == 1:
            if probability == 1:
                return self.getReadDynString(var_list[0])
            else:
                if var_list[0][1] == 0:
                    return "DynMap[{}] * {}".format(var_list[0][0], probability)
                else:
                    return "DynMap[{}+{}] * {}".format(var_list[0][0], var_list[0][1], probability)
        else:
            sum_str = []
            for var in var_list:
                    sum_str.append(self.getReadDynString(var))
            if len(var_list) == 2:
                upd_str = " + ".join(sum_str) + " - 1.0"
            else:
                upd_str = " + ".join(sum_str) + " - float64({})".format(len(var_list) - 1)
            if probability == 1:
                return upd_str
            else:
                return "({}) * {}".format(upd_str, probability)

    def updateCondDyn(self, val, condDyn, upd):
        # print "-----------: ", condDyn.updated
        # collections.namedtuple('ConditionalDynUpdate', ['condition', 'ifop', 'elseop'])
        if (val, 0) in condDyn.ifop.sum:
            newDyn1 = DynUpdate(condDyn.ifop.updated, condDyn.ifop.sum + upd.sum,
                                condDyn.ifop.multiplicative * upd.multiplicative,
                                condDyn.ifop.isarray)
        else:
            newDyn1 = condDyn.ifop

        if (val, 0) in condDyn.elseop.sum:
            newDyn2 = DynUpdate(condDyn.elseop.updated, condDyn.elseop.sum + upd.sum,
                                condDyn.elseop.multiplicative * upd.multiplicative,
                                condDyn.elseop.isarray)
        else:
            newDyn2 = condDyn.elseop
        return ConditionalDynUpdate(condDyn.condition, newDyn1, newDyn2, condDyn.updated)

    def simpleDCE(self):
        remove_list = []
        for i, upd in enumerate(self.tracking):
            # For now keep the last
            # TODO: optimize?
            if i == len(self.tracking) - 1:
                break
            if not isinstance(upd, ConditionalDynUpdate):
                for j in range(i + 1, len(self.tracking)):
                    # check if it apprears in any statement later
                    if isinstance(self.tracking[j], ConditionalDynUpdate):
                        if ((upd.updated, upd.isarray) in self.tracking[j].ifop.sum or
                                (upd.updated, upd.isarray) in self.tracking[j].elseop.sum):
                            break
                        if upd.updated == self.tracking[j].updated:
                            remove_list.append(upd)
                            break
                    else:
                        if (upd.updated, upd.isarray) in self.tracking[j].sum:
                            break
                        if (upd.updated, upd.isarray) == (self.tracking[j].updated, self.tracking[j].isarray):
                            remove_list.append(upd)
                            break
            else:
                for j in range(i + 1, len(self.tracking)):
                    # check if it apprears in any statement later
                    if isinstance(self.tracking[j], ConditionalDynUpdate):
                        if ((upd.updated, 0) in self.tracking[j].ifop.sum or
                                (upd.updated, 0) in self.tracking[j].elseop.sum):
                            break
                        if upd.updated == self.tracking[j].updated:
                            remove_list.append(upd)
                            break
                    else:
                        if (upd.updated, 0) in self.tracking[j].sum:
                            break
                        if (upd.updated, 0) == (self.tracking[j].updated, self.tracking[j].isarray):
                            remove_list.append(upd)
                            break

        for item in remove_list:
            self.tracking.remove(item)

    def simpleOptimizeConsts(self):
        # First perform a simple constant replacement
        for i, upd in enumerate(self.tracking):
            if isinstance(upd, ConditionalDynUpdate):
                continue
            # Last one
            if i == len(self.tracking) - 1:
                break
            if upd.isarray:
                # print "Skipping: ", upd
                continue
            # print "Trying to substitute: ", upd, i, range(i + 1, len(self.tracking))
            val = upd.updated
            if (val,0) in upd.sum:
                continue
            for j in range(i + 1, len(self.tracking)):
                if isinstance(self.tracking[j], ConditionalDynUpdate):
                    self.tracking[j] = self.updateCondDyn((val, 0), self.tracking[j], upd)
                    if (self.tracking[j].updated, 0) in upd.sum:
                        break
                else:
                    if (val, 0) in self.tracking[j].sum:
                        # print "==>:", self.tracking[j]
                        self.tracking[j].sum.remove((val, 0))
                        newDyn = DynUpdate(self.tracking[j].updated, self.tracking[j].sum + upd.sum,
                                           self.tracking[j].multiplicative * upd.multiplicative,
                                           self.tracking[j].isarray)
                        self.tracking[j] = newDyn
                        # print "====>:", self.tracking[j]
                    if (self.tracking[j].updated, self.tracking[j].isarray) in upd.sum:
                        break
                if val == self.tracking[j].updated:
                    break

        # Next perform a very simple DCE
        self.simpleDCE()
        return self.tracking

    def getDynString(self):
        dyn_str = []
        dyn_upd_map = "DynMap[{}] = {};\n"
        dyn_upd_map_array = "DynMap[{} + _temp_index_{}] = {};\n"

        tracking_list = self.simpleOptimizeConsts()

        if len(tracking_list) == 0:
            return ''

        for update in tracking_list:
            if isinstance(update, ConditionalDynUpdate):
                # print "===========: ", update.updated
                # t_d_str.format(self.tempindexnum, self.varMap[b_var], self.varMap[a_var],
                #                    self.varMap[o1_var], self.varMap[o2_var])
                option1 = dyn_upd_map.format(update.updated,
                                             self.getUpdateString(update.ifop.sum,
                                                                  update.ifop.multiplicative))
                option2 = dyn_upd_map.format(update.updated,
                                             self.getUpdateString(update.elseop.sum,
                                                                  update.elseop.multiplicative))
                dyn_str.append(t_d_str2.format(update.condition, update.updated, option1, option2))
            else:
                if update.isarray != 0:
                    dyn_str.append(dyn_upd_map_array.format(update.updated,
                                                            update.isarray,
                                                            self.getUpdateString(update.sum,
                                                                                 update.multiplicative)))
                else:
                    dyn_str.append(dyn_upd_map.format(update.updated,
                                                      self.getUpdateString(update.sum,
                                                                           update.multiplicative)))
        return ''.join(dyn_str)

    def visitRepeatlvar(self, ctx):
        pre_string = ''
        if self.args.gather:
            dyn_str = ''.join(self.trackingStatements)
            pre_string = dyn_str
            d_str_opt = self.getDynString()
            pre_string += d_str_opt

        self.trackingStatements = []
        self.tracking = []

        repeatVar = ctx.var().getText()
        temp_var_name = "__temp_{}".format(self.tempvarnum)
        self.tempvarnum += 1

        statement_string = ''
        for statement in ctx.statement():
            translated = self.visit(statement)
            if translated is not None:
                statement_string += translated
            else:
                print "[Error] Unable to translate: ", statement.getText()
                exit(-1)

        if self.args.gather:
            dyn_str = ''.join(self.trackingStatements)
            statement_string += dyn_str
            d_str_opt = self.getDynString()
            # print d_str_opt
            statement_string += d_str_opt
        self.trackingStatements = []
        self.tracking = []

        str_for_loop = pre_string + "for {} := 0; {} < {}; {}++ {{\n {} }}\n"
        return str_for_loop.format(temp_var_name, temp_var_name, repeatVar, temp_var_name, statement_string)


    def visitRepeat(self, ctx):
        pre_string = ''
        if self.args.gather:
            dyn_str = ''.join(self.trackingStatements)
            pre_string = dyn_str
            d_str_opt = self.getDynString()
            pre_string += d_str_opt

        self.trackingStatements = []
        self.tracking = []

        repeatNum = ctx.INT().getText()
        temp_var_name = "__temp_{}".format(self.tempvarnum)
        self.tempvarnum += 1

        statement_string = ''
        for statement in ctx.statement():
            translated = self.visit(statement)
            if translated is not None:
                statement_string += translated
            else:
                print "[Error] Unable to translate: ", statement.getText()
                exit(-1)

        if self.args.gather:
            dyn_str = ''.join(self.trackingStatements)
            statement_string += dyn_str
            d_str_opt = self.getDynString()
            # print d_str_opt
            statement_string += d_str_opt

        self.tracking = []
        self.trackingStatements = []
        str_for_loop = pre_string + "for {} := 0; {} < {}; {}++ {{\n {} }}\n"
        return str_for_loop.format(temp_var_name, temp_var_name, repeatNum, temp_var_name, statement_string)

    def visitForloop(self, ctx):
        pre_string = ''
        if self.args.gather:
            dyn_str = ''.join(self.trackingStatements)
            pre_string = dyn_str
            d_str_opt = self.getDynString()
            pre_string += d_str_opt

        self.trackingStatements = []
        self.tracking = []

        group_name = ctx.GLOBALVAR().getText()
        var_name = ctx.VAR().getText()
        # for proc in self.proc_groups[group_name]:
        statement_string = ''
        for statement in ctx.statement():
            translated = self.visit(statement)
            if translated is not None:
                statement_string += translated
            else:
                print "[Error] Unable to translate: ", statement.getText()
                exit(-1)

        if self.args.gather:
            dyn_str = ''.join(self.trackingStatements)
            statement_string += dyn_str
            d_str_opt = self.getDynString()
            # print d_str_opt
            statement_string += d_str_opt

        self.tracking = []
        self.trackingStatements = []
        str_for_loop = pre_string + "for _, {} := range({}) {{\n {} }}\n"
        return str_for_loop.format(var_name, group_name, statement_string)

    def visitFunc(self, ctx):
        return ctx.getText() + ";\n"

    def visitInstrument(self, ctx):
        # print ctx.getText()
        if self.args.instrument:
            return ctx.code.text[2:] + ";\n"
        else:
            return ""

    def visitSpeccheckarray(self, ctx):
        if not self.enableDynamic:
            return ""
        statement_string = ''
        if self.args.gather:
            dyn_str = ''.join(self.trackingStatements)
            statement_string += dyn_str
            d_str_opt = self.getDynString()
            # print d_str_opt
            statement_string += d_str_opt

        self.tracking = []
        self.trackingStatements = []
        checked_var = ctx.var().getText()
        checked_val = ctx.probability().getText()

        return statement_string + ch_str.format(self.varMap[checked_var], checked_val,
                                                self.arraySize[checked_var], checked_var)

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
        # print dectype

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
                self.dynsize += int(decl.INT()[0].getText())

                # print "Increasing dynamic size: ", self.dynsize

                # only works for 1 dimentional so far
                d_str = ""
                if self.enableDynamic:
                    d_str = "diesel.InitDynArray({}, {}, DynMap[:]);\n".format(self.varMap[varname],
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

            dim = []
            for dimention in decl.GLOBALVAR():
                dim.append(dimention)

            d_str = ""
            if self.enableDynamic and dectype[0] == "dynamic":
                # Only works for 1 dimentional so far
                print "Not supporting dynamic unbounded arrays"
                exit(-1)
                d_str = "diesel.InitDynArray(\"{}\", {}, DynMap);\n".format(varname,
                                                                               decl.GLOBALVAR()[0])

            if len(dim) == 1:
                return dyn_array_dec.format(varname, "[]", dectype[1], dim[0]) + d_str
            if len(dim) > 1:
                return dyn_array_dec.format(varname, "[]", dectype[1], dim[0]) + d_str
            else:
                print "[Error] Unable to translate: ", decl.getText()
                exit(-1)
        else:
            varname = decl.var().getText()
            self.primitiveTMap[varname] = dectype[0]
            self.typeMap[varname] = (dectype[1], dectype[2])
            # print "****:" ,varname
            if self.enableDynamic and dectype[0] == "dynamic":
                self.varMap[varname] = self.varNum
                self.varNum += 1
                self.dynsize += 1
                print "Increasing dynamic size: ", self.dynsize
                d_init_str = "var {0} {1};\nDynMap[{2}] = diesel.ProbInterval{{1, 0}};\n"
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

        # print dec_string
        
        statement_string = ""
        for statement in ctx.statement():
            translated = self.visit(statement)
            if translated is not None:
                statement_string += translated
            else:
                print "[Error] Unable to translate: ", statement.getText()
                exit(-1)

        if self.args.gather:
            dyn_str = ''.join(self.trackingStatements)
            statement_string += dyn_str
            d_str_opt = self.getDynString()
            # print d_str_opt
            statement_string += d_str_opt

        self.tracking = []
        self.trackingStatements = []
        process_name = "func_" + group_name
        self.process_list.append(((process_name, group_name), 1))

        process_code = dec_string + statement_string

        process_def_str = str_member_thread.format(process_name, self.dynsize, process_code, element_name)
        self.process_defs.append(process_def_str)

    def visitSingle(self, ctx):
        # self.primitiveTMap = {}
        # self.typeMap = {}
        # self.varMap = {}
        self.varNum = 0
        self.trackingStatements = []
        self.tracking = []
        self.allTracking = []
        self.tempindexnum = 0
        self.dynsize = 0
        
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
            if translated is not None:
                statement_string += translated
            else:
                print "[Error] Unable to translate: ", statement.getText()
                exit(-1)

        # print self.VarMap
        if self.args.gather:
            dyn_str = ''.join(self.trackingStatements)
            self.trackingStatements = []
            statement_string += dyn_str
            d_str_opt = self.getDynString()
            # print d_str_opt
            statement_string += d_str_opt

        process_name = "func_" + self.pid
        self.process_list.append((process_name, 0))

        process_code = dec_string + statement_string

        process_def_str = str_single_thread.format(process_name, self.dynsize, process_code, self.pid)
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
    parser.add_argument("-i", "--instrument", action="store_true",
                        help="Add instrumentation")
    parser.add_argument("-g", "--gather", action="store_true",
                        help="Collect tracking together")
    args = parser.parse_args()

    if args.dynamic:
        print "Enabling dynamic tracking"
    if args.arrayO1:
        print "Enabling array optimization: Send one value"
    if args.gather:
        print "Enabling gather optimization + Simple DCE"

    programfile = open(args.programfile, 'r')
    # outfile = open(args.outfile, 'w')
    program_str = programfile.read()

    startTime = time.time()
    main(program_str, args.outfile, programfile.name, args.template,
         args.debug, args.dynamic, args)
    print "Done!";
    print "Elapsed time : ", time.time()-startTime;
