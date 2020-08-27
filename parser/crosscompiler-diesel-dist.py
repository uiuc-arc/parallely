from antlr4 import CommonTokenStream
from antlr4 import InputStream
from ParallelyLexer import ParallelyLexer
from ParallelyParser import ParallelyParser
from ParallelyVisitor import ParallelyVisitor
from argparse import ArgumentParser
import collections
import time
import crosscompilerconstants as constants

LIBRARYNAME = "dieseldist"

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
            EXITWITHERROR("[Error] Cant handle process name variables")
            # return (True, pid.VAR().getText(),)
        else:
            return (True, pid.GLOBALVAR().getText(), pid.VAR().getText())

def EXITWITHERROR(msg):
    print(msg)
    exit(-1)    

class CountThreads(ParallelyVisitor):
    def __init__(self):
        print("Counting the number of processes")
        self.processes = {}
        self.processcount = 0

    def visitSingleglobaldec(self, ctx):
        global_var = ctx.GLOBALVAR().getText()
        members = [temp.getText() for temp in ctx.processid()]
        self.processes[global_var] = members

    # in theory pids are not int. Changing to simplify implementation
    def visitSingle(self, ctx):
        pid = isGroup(ctx.processid())
        if pid[0]:
            self.processcount += len(self.processes[pid[1]])
        else:
            self.processcount += 1


class Translator(ParallelyVisitor):
    def __init__(self, dynamic, args):
        print("[crosscompiler] Starting translation")
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
            EXITWITHERROR("[Error] Unknown type : ", fulltype.getText())

    def visitCondsend(self, ctx):
        cond_send_str = {
            ("int", 0): "LIBRARYNAME.Condsend({}, {}, {}, {});\n",
            ("int32", 0): "LIBRARYNAME.CondsendInt32({}, {}, {}, {});\n",
            ("int64", 0): "LIBRARYNAME.CondsendInt64({}, {}, {}, {});\n",
            ("float32", 0): "LIBRARYNAME.CondsendFloat32({}, {}, {}, {});\n",
            ("float64", 0): "LIBRARYNAME.CondsendFloat64({}, {}, {}, {});\n",
            ("int", 1): "LIBRARYNAME.CondsendIntArray({}, {}[:], {}, {});\n",
            ("int", 1): "LIBRARYNAME.CondsendIntArray({}, {}[:], {}, {});\n",
        }

        cond_var = ctx.var(0).getText()
        sent_var = ctx.var(1).getText()
        senttype = self.getType(ctx.fulltype())
        return cond_send_str[(senttype[1], senttype[2])].format(cond_var, sent_var, self.pid,
                                                                ctx.processid().getText())

    def visitSend(self, ctx):
        send_str = {
            ("int", 0): "LIBRARYNAME.SendInt({}, {}, {});\n",
            ("int32", 0): "LIBRARYNAME.SendInt32({}, {}, {});\n",
            ("int64", 0): "LIBRARYNAME.SendInt64({}, {}, {});\n",
            ("float32", 0): "LIBRARYNAME.SendFloat32({}, {}, {});\n",
            ("float64", 0): "LIBRARYNAME.SendFloat64({}, {}, {});\n",
            ("int", 1): "LIBRARYNAME.SendIntArray({}[:], {}, {});\n",
            ("int32", 1): "LIBRARYNAME.SendInt32Array({}[:], {}, {});\n",
            ("int64", 1): "LIBRARYNAME.SendInt64Array({}[:], {}, {});\n",
            ("float32", 1): "LIBRARYNAME.SendFloat32Array({}[:], {}, {});\n",
            ("float64", 1): "LIBRARYNAME.SendFloat64Array({}[:], {}, {});\n"
        }

        dyn_send_str = {
            ("int", 0): "LIBRARYNAME.SendDynIntArray({}[:], {}, {}, DynMap[:], {});\n",
            ("float64", 0): "LIBRARYNAME.SendDynFloat64Array({}[:], {}, {}, DynMap[:], {});\n",
            ("float32", 0): "LIBRARYNAME.SendDynFloat32Array({}[:], {}, {}, DynMap[:], {});\n",
            ("int", 1): "LIBRARYNAME.SendDynIntArrayO1({}[:], {}, {}, DynMap[:], {});\n",
            ("float64", 1): "LIBRARYNAME.SendDynFloat64ArrayO1({}[:], {}, {}, DynMap[:], {});\n",
            ("float32", 1): "LIBRARYNAME.SendDynFloat32ArrayO1({}[:], {}, {}, DynMap[:], {});\n",
        }

        sent_var = ctx.var().getText()
        senttype = self.getType(ctx.fulltype())

        if (sent_var in self.arrays and self.enableDynamic and
                sent_var in self.primitiveTMap and self.primitiveTMap[sent_var] == 'dynamic'):
            t_str = dyn_send_str[senttype[1], self.args.arrayO1].format(sent_var, self.pid,
                                                                        ctx.processid().getText(),
                                                                        self.varMap[sent_var])
            return t_str

        s_str_0 = send_str[senttype[1], senttype[2]].format(sent_var, self.pid, ctx.processid().getText())
        d_str = ""
        if self.enableDynamic and sent_var in self.primitiveTMap and self.primitiveTMap[sent_var] == 'dynamic':
            v_str = "DynMap[{}]".format(self.varMap[sent_var])
            d_str = "LIBRARYNAME.SendDynVal({}, {}, {});\n".format(v_str, self.pid, ctx.processid().getText())
        return s_str_0 + d_str

    def visitReceive(self, ctx):
        rec_str = {
            ("int", 0): "LIBRARYNAME.ReceiveInt(&{}, {}, {});\n",
            ("int32", 0): "LIBRARYNAME.ReceiveInt32(&{}, {}, {});\n",
            ("int64", 0): "LIBRARYNAME.ReceiveInt64(&{}, {}, {});\n",
            ("float32", 0): "LIBRARYNAME.ReceiveFloat32(&{}, {}, {});\n",
            ("float64", 0): "LIBRARYNAME.ReceiveFloat64(&{}, {}, {});\n",
            ("int", 1): "LIBRARYNAME.ReceiveIntArray({}[:], {}, {});\n",
            ("int32", 1): "LIBRARYNAME.ReceiveInt32Array({}[:], {}, {});\n",
            ("int64", 1): "LIBRARYNAME.ReceiveInt64Array({}[:], {}, {});\n",
            ("float32", 1): "LIBRARYNAME.ReceiveFloat32Array({}[:], {}, {});\n",
            ("float64", 1): "LIBRARYNAME.ReceiveFloat64Array({}[:], {}, {});\n"
        }

        dyn_rec_dict = {
            ("int", 0): "LIBRARYNAME.ReceiveDynIntArray({}[:], {}, {}, DynMap[:], {});\n",
            ("float64", 0): "LIBRARYNAME.ReceiveDynFloat64Array({}[:], {}, {}, DynMap[:], {});\n",
            ("float32", 0): "LIBRARYNAME.ReceiveDynFloat32Array({}[:], {}, {}, DynMap[:], {});\n",
            ("int", 1): "LIBRARYNAME.ReceiveDynIntArrayO1({}[:], {}, {}, DynMap[:], {});\n",
            ("float64", 1): "LIBRARYNAME.ReceiveDynFloat64ArrayO1({}[:], {}, {}, DynMap[:], {});\n",
            ("float32", 1): "LIBRARYNAME.ReceiveDynFloat32ArrayO1({}[:], {}, {}, DynMap[:], {});\n",
        }

        noisy_dyn_rec_dict = {
            ("int", 0): "LIBRARYNAME.NoisyReceiveDynIntArray({}[:], {}, {}, DynMap[:], {});\n",
            ("float64", 0): "LIBRARYNAME.NoisyReceiveDynFloat64Array({}[:], {}, {}, DynMap[:], {});\n",
            # ("float32", 0): "LIBRARYNAME.NoisyReceiveDynFloat32Array({}[:], {}, {}, DynMap[:], {});\n",
            ("int", 1): "LIBRARYNAME.NoisyReceiveDynIntArrayO1({}[:], {}, {}, DynMap[:], {});\n",
            ("float64", 1): "LIBRARYNAME.NoisyReceiveDynFloat64ArrayO1({}[:], {}, {}, DynMap[:], {});\n",
            # ("float32", 1): "LIBRARYNAME.NoisyReceiveDynFloat32ArrayO1({}[:], {}, {}, DynMap[:], {});\n",
        }

        senttype = self.getType(ctx.fulltype())
        rec_var = ctx.var().getText()

        if (rec_var in self.arrays and self.enableDynamic and
                rec_var in self.primitiveTMap and self.primitiveTMap[rec_var] == 'dynamic'):
            if not self.args.noisy:
                return dyn_rec_dict[senttype[1], self.args.arrayO1].format(ctx.var().getText(),
                                                                           self.pid,
                                                                           ctx.processid().getText(),
                                                                           self.varMap[rec_var])
            else:
                return noisy_dyn_rec_dict[senttype[1], self.args.arrayO1].format(ctx.var().getText(),
                                                                           self.pid,
                                                                           ctx.processid().getText(),
                                                                           self.varMap[rec_var])

        rec_str_0 = rec_str[senttype[1], senttype[2]].format(ctx.var().getText(),
                                                             self.pid, ctx.processid().getText())
        d_str = ""
        if self.enableDynamic and rec_var in self.primitiveTMap and self.primitiveTMap[rec_var] == 'dynamic':
            self.tempvarnum += 1
            d_str = constants.dyn_rec_str.format(ctx.processid().getText(),
                                                 self.pid, self.varMap[ctx.var().getText()], self.tempvarnum)

        return rec_str_0 + d_str

    def visitCondreceive(self, ctx):
        rec_str = {
            ("int", 0): "LIBRARYNAME.Condreceive(&{}, &{}, {}, {});\n",
            ("int32", 0): "LIBRARYNAME.CondreceiveInt32(&{}, &{}, {}, {});\n",
            ("int64", 0): "LIBRARYNAME.CondreceiveInt64(&{}, &{}, {}, {});\n",
            # ("float32", 0):
            # ("float64", 0):
            ("int", 1): "LIBRARYNAME.CondreceiveIntArray(&{}, {}[:], {}, {});\n",
            # ("int32", 1):
            # ("int64", 1):
            ("float32", 1): "LIBRARYNAME.CondreceiveFloat32(&{}, &{}, {}, {});\n",
            ("float64", 1): "LIBRARYNAME.CondreceiveFloat64(&{}, &{}, {}, {});\n"
        }
        senttype = self.getType(ctx.fulltype())
        return rec_str[senttype[1], senttype[2]].format(ctx.var(0).getText(), ctx.var(1).getText(),
                                                        self.pid, ctx.processid().getText())

    def visitProbassignment(self, ctx):
        rand_str = {
            ("int"): "{} = LIBRARYNAME.Randchoice(float32({}), {}, {});\n",
            ("float64"): "{} = LIBRARYNAME.RandchoiceFloat64(float32({}), {}, {});\n",
        }

        a_var = ctx.var().getText()
        prob = ctx.probability().getText()

        if (self.recovernum == 0):
            p_str = rand_str[self.typeMap[a_var][0]].format(a_var, prob, ctx.expression(0).getText(),
                                                            ctx.expression(1).getText())
        else:
            p_str = constants.str_probchoiceIntFlag.format(a_var, prob, ctx.expression(0).getText(),
                                                           ctx.expression(1).getText(), self.recovernum)
        if self.enableDynamic and a_var in self.primitiveTMap and self.primitiveTMap[a_var] == 'dynamic':
            if self.args.accuracy and (not self.args.reliability):
                dyn_str = constants.prob_assignment_str_const[LIBRARYNAME].format(self.varMap[a_var],
                                                                                  ctx.probability().getText())
                print("[crosscompiler] Warning - Accuracy analysis for probchoice resolves to infinity")
                return p_str + dyn_str
                
            var_list = list(set(self.getVarList(ctx.precise)))
            if len(var_list) == 0:
                dyn_str = constants.prob_assignment_str_const[LIBRARYNAME].format(self.varMap[a_var],
                                                                                  ctx.probability().getText())
            elif len(var_list) == 1:
                dyn_str = constants.prob_assignment_str_single[LIBRARYNAME].format(self.varMap[a_var],
                                                                                   self.varMap[var_list[0]],
                                                                                   ctx.probability().getText())
            else:
                sum_str = []
                for var in var_list:
                    sum_str.append(constants.access_reliability[LIBRARYNAME].format(self.varMap[var]))
                dyn_str = constants.dyn_pchoice_str[LIBRARYNAME].format(self.varMap[a_var], " + ".join(sum_str),
                                                                        len(var_list) - 1,
                                                                        ctx.probability().getText())
            if self.args.accuracy:
                dyn_str += constants.dyn_pchoice_accuracy[LIBRARYNAME].format(self.varMap[a_var])
            p_str = p_str + dyn_str
        return p_str

    # def getDynUpdate(self, expression, probability, updated_var, isarray):
    #     var_list = list(set(self.getVarList(expression)))
    #     temp_dyn = DynUpdate(self.varMap[updated_var], [(self.varMap[var], 0) for var in var_list],
    #                          float(probability), isarray)
    #     return temp_dyn

    def handleExpression(self, ctx):
        # print ctx.getText()
        convert_str = "LIBRARYNAME.ConvBool({})"
        if isinstance(ctx, ParallelyParser.SelectContext):
            return self.handleExpression(ctx.expression())
        if (isinstance(ctx, ParallelyParser.EqContext) or
                isinstance(ctx, ParallelyParser.GeqContext) or
                isinstance(ctx, ParallelyParser.LeqContext) or
                isinstance(ctx, ParallelyParser.LessContext) or
                isinstance(ctx, ParallelyParser.GreaterContext)):           
            # If it is a boolean statement
            return convert_str.format(ctx.getText())
        if  isinstance(ctx, ParallelyParser.AndexpContext):
            convert_str = "LIBRARYNAME.ConvBool({}==1 && {}==1)".format(ctx.expression(0).getText(),
                                                                       ctx.expression(1).getText())
            return convert_str
        if  isinstance(ctx, ParallelyParser.OrexpContext):
            convert_str = "LIBRARYNAME.ConvBool({}==1 || {}==1)".format(ctx.expression(0).getText(),
                                                                       ctx.expression(1).getText())
            return convert_str
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

    def visitDyncondassignmentgeq(self, ctx):
        a_var = ctx.assigned.getText()
        l_var = ctx.lvar.getText()
        r_var = ctx.rvar.getText()
        o1_var = ctx.ifvar.getText()
        o2_var = ctx.elsevar.getText()
        
        if self.enableDynamic:
            a_var_index = self.varMap[a_var]
            l_var_index = self.varMap[l_var]
            r_var_index = self.varMap[r_var]
            o1_var_index = self.varMap[o1_var]
            o2_var_index = self.varMap[o2_var]
            convert_cond_str = constants.convert_cond_str_list[(self.typeMap[l_var][0],
                                                                self.typeMap[a_var][0])]
            out_str = convert_cond_str.format(a_var, l_var, r_var,
                                              l_var_index, r_var_index,
                                              o1_var, o2_var,
                                              o1_var_index, o2_var_index, a_var_index)
        else:
            cond_str = "if {}>={} {{ {}={} }} else {{ {}={} }};\n"
            out_str = cond_str.format(l_var, r_var, a_var, o1_var, a_var, o2_var)
        return out_str

    def visitCondassignment(self, ctx):
        assign_str = "temp_bool_{4}:= {0}; if temp_bool_{4} != 0 {{ {1}  = {2} }} else {{ {1} = {3} }};\n"
        a_var = ctx.var()[0].getText()
        b_var = ctx.condition.getText()
        o1_var = ctx.ifvar.getText()
        o2_var = ctx.elsevar.getText()

        self.tempindexnum += 1
        out_str = assign_str.format(b_var, a_var, o1_var, o2_var, self.tempindexnum)
        d_str = ""
        if self.enableDynamic and a_var in self.primitiveTMap and self.primitiveTMap[a_var] == 'dynamic':
            d_str = constants.condassignment_dyn_str[LIBRARYNAME].format(self.tempindexnum,
                                                                         self.varMap[a_var],
                                                                         self.varMap[b_var],
                                   self.varMap[o1_var], self.varMap[o2_var])
        return out_str + d_str

    def visitExpassignment(self, ctx):
        assign_str = "{} = {};\n"
        expr_str = self.handleExpression(ctx.expression())
        var_str = ctx.var().getText()

        # hack to handle array copy
        # print var_str, expr_str, self.arrays
        if (self.enableDynamic and var_str in self.arrays and
                expr_str in self.arrays and self.primitiveTMap[var_str] == 'dynamic'):
            if self.primitiveTMap[expr_str] == 'precise':
                dyn_c_str = "LIBRARYNAME.InitDynArray({}, {}, DynMap[:]);\n".format(self.varMap[var_str],
                                                                                    self.arraySize[var_str])
            else:
                dyn_c_str = "LIBRARYNAME.CopyDynArray({},{},{},DynMap[:]);\n".format(self.varMap[var_str],
                                                                                     self.varMap[expr_str],
                                                                                     self.arraySize[var_str])
            return ctx.getText() + ";\n" + dyn_c_str

        dyn_str = ""
        acc_str = ""
        
        # Not global and dynamic
        if self.enableDynamic and var_str in self.primitiveTMap and self.primitiveTMap[var_str] == 'dynamic':
            var_list = list(set(self.getVarList(ctx.expression())))

            if len(var_list) == 0:
                dyn_str = constants.dyn_expression_str_precise[LIBRARYNAME].format(self.varMap[var_str])
            elif len(var_list) == 1:
                dyn_str = constants.dyn_expression_str_single[LIBRARYNAME].format(self.varMap[var_str],
                                                                                  self.varMap[var_list[0]])
                if self.args.accuracy:
                    acc_str = self.getAccuracyStr(ctx.expression(), var_str)
            else:
                sum_str = []
                for var in var_list:
                    sum_str.append(constants.access_reliability[LIBRARYNAME].format(self.varMap[var]))
                dyn_str = constants.dyn_assign_str[LIBRARYNAME].format(self.varMap[var_str],
                                                                       " + ".join(sum_str),
                                                                       len(var_list) - 1)
                if self.args.accuracy:
                    acc_str = self.getAccuracyStr(ctx.expression(), var_str)
        return  dyn_str + acc_str + assign_str.format(var_str, expr_str)

    def getAccuracyStr(self, ctx, var_str):
        if (isinstance(ctx, ParallelyParser.FliteralContext)
            or isinstance(ctx, ParallelyParser.LiteralContext)):
            return "" # "DynMap[{}].Delta = 0;\n".format(self.varMap[var_str])
        elif isinstance(ctx, ParallelyParser.EqContext):
            return ""
        # We need to fix this
        elif isinstance(ctx, ParallelyParser.GreaterContext):
            return ""
        elif isinstance(ctx, ParallelyParser.SelectContext):
            return self.getAccuracyStr(ctx.expression(), var_str)
        elif isinstance(ctx, ParallelyParser.VariableContext):
            return constants.dyn_accuracy_update[LIBRARYNAME].format(self.varMap[var_str],
                                                                     self.varMap[ctx.getText()])
        elif isinstance(ctx, ParallelyParser.VarContext):
            return constants.dyn_accuracy_update[LIBRARYNAME].format(self.varMap[var_str],
                                                                     self.varMap[ctx.getText()])
        elif isinstance(ctx, ParallelyParser.AddContext) or isinstance(ctx, ParallelyParser.MinusContext):
            var_list = self.getVarList(ctx)
            if len(var_list) == 0:
                return ""
            if len(var_list) == 1:
                return constants.dyn_accuracy_update[LIBRARYNAME].format(self.varMap[var_str],
                                                                         self.varMap[var_list[0]])
            elif len(var_list) == 2:
                d_str = constants.dyn_accuracy_update_double[LIBRARYNAME]
                return d_str.format(self.varMap[var_str],
                                    self.varMap[var_list[0]],
                                    self.varMap[var_list[1]])
            else:
                EXITWITHERROR("[ERROR]: Only support simple expressions: ", ctx.getText(), var_list)
                
        elif isinstance(ctx, ParallelyParser.MultiplyContext):
            var_list = self.getVarList(ctx)
            if len(var_list) == 0:
                EXITWITHERROR("[ERROR]: should not have a separate accuracy string: ",
                              ctx.getText(), var_list)
            if len(var_list) == 1:
                dyn_str = constants.dyn_accuracy_mult_single_str[LIBRARYNAME]
                if (isinstance(ctx.expression(0), ParallelyParser.FliteralContext) or
                    isinstance(ctx.expression(0), ParallelyParser.LiteralContext) or
                    (ctx.expression(0).getText() in self.primitiveTMap and
                     self.primitiveTMap[ctx.expression(0).getText()] == 'dynamic')):
                    return dyn_str.format(self.varMap[var_str], ctx.expression(1).getText(),
                                          self.varMap[var_list[0]])
                elif (isinstance(ctx.expression(1), ParallelyParser.FliteralContext) or
                      isinstance(ctx.expression(1), ParallelyParser.LiteralContext) or
                      (ctx.expression(1).getText() in self.primitiveTMap and
                       self.primitiveTMap[ctx.expression(1).getText()] == 'dynamic')):
                    dyn_str = constants.dyn_accuracy_mult_single_str[LIBRARYNAME]
                    return dyn_str.format(self.varMap[var_str],
                                          ctx.expression(0).getText(), self.varMap[var_list[0]])
                else:
                    return dyn_str.format(var_str, var_list[0], self.varMap[var_list[0]],
                                      var_list[1], self.varMap[var_list[1]])
            elif len(var_list) == 2:
                upd_str = constants.dyn_accuracy_mult_double_str[LIBRARYNAME]
                return upd_str.format(self.varMap[var_str], var_list[0], self.varMap[var_list[0]],
                                      var_list[1], self.varMap[var_list[1]])
        elif isinstance(ctx, ParallelyParser.DivideContext):
            var_list = self.getVarList(ctx)
            #implement the zero check at some point
            if len(var_list) == 1:
                if (isinstance(ctx.expression(0), ParallelyParser.FliteralContext) or
                    isinstance(ctx.expression(0), ParallelyParser.LiteralContext) or
                    (ctx.expression(0).getText() in self.primitiveTMap and
                     self.primitiveTMap[ctx.expression(0).getText()] == 'dynamic')):
                    upd_str = constants.dyn_accuracy_div_single_str_0[LIBRARYNAME]
                    return upd_str.format(self.varMap[var_str],
                                          ctx.expression(1).getText(), self.varMap[var_list[0]])

                elif (isinstance(ctx.expression(1), ParallelyParser.FliteralContext) or
                      isinstance(ctx.expression(1), ParallelyParser.LiteralContext) or
                (ctx.expression(1).getText() in self.primitiveTMap and self.primitiveTMap[ctx.expression(1).getText()] == 'dynamic')):
                    upd_str = constants.dyn_accuracy_div_single_str_1[LIBRARYNAME]
                    return upd_str.format(self.varMap[var_str],
                                          ctx.expression(0).getText(), self.varMap[var_list[0]])
                # return upd_str.format(var_str, var_list[0], self.varMap[var_list[0]],
                #                       var_list[1], self.varMap[var_list[1]])
            elif len(var_list) == 2:
                upd_str = constants.dyn_accuracy_div_double_str[LIBRARYNAME]
                return upd_str.format(self.varMap[var_str], var_list[0], self.varMap[var_list[0]],
                                      var_list[1], self.varMap[var_list[1]])
        print("[WARNING!!!!!] Not generating accuracy tracking for: ", ctx.getText(), type(ctx))
        return ""

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
            return go_str.format(index_expr, assigned_var, array_var, self.tempindexnum) + d_str
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
            var_list = list(set(self.getVarList(ctx.expression()[1])))
            if len(var_list) == 0:
                dyn_str = constants.dyn_array_store_precise[LIBRARYNAME].format(self.varMap[a_var],
                                                                                self.tempindexnum)
            elif len(var_list) == 1:
                dyn_str = constants.dyn_array_store_single_dyn.format(self.varMap[a_var],
                                                                      self.varMap[var_list[0]],
                                                                      self.tempindexnum)
            else:
                EXITWITHERROR("Array Store Complex Expression: " + ctx.getString()) 
                # dyn_upd_map = "DynMap[{0} + _temp_index_{3}] = dieseldist.Max(0.0, {1} - float64({2}));\n"
                # sum_str = []
                # for var in var_list:
                #     sum_str.append("DynMap[{}]".format(self.varMap[var]))
                # dyn_str = dyn_upd_map.format(self.varMap[a_var], " + ".join(sum_str),
                #                              len(var_list) - 1, self.tempindexnum)
        return r_str + dyn_str

    def visitCast(self, ctx):
        resultType = self.getType(ctx.fulltype())
        assignedvar = ctx.var(0).getText()
        castedvar = ctx.var(1).getText()
        # Array type
        if resultType[1] == "float64" and resultType[2] == 1:
            return "LIBRARYNAME.Cast32to64Array({}[:], {}[:]);\n".format(assignedvar,
                                                                         castedvar)
        if resultType[1] == "float32" and resultType[2] == 1:
            return "LIBRARYNAME.Cast64to32Array({}[:], {}[:]);\n".format(assignedvar,
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
                if castedvar in self.varMap:
                    d_str = constants.dyn_cast_dyn_str[LIBRARYNAME]
                    d_str = d_str.format(self.varMap[assignedvar], castedvar,
                                         assignedvar, self.varMap[castedvar])
                else:
                    d_str = constants.dyn_cast_precise_str[LIBRARYNAME]
                    d_str = d_str.format(self.varMap[assignedvar], castedvar, assignedvar)
            return "{} = float32({});\n".format(assignedvar, castedvar) + d_str

    def visitTrack(self, ctx):
        statement_string="{}={};\n".format(ctx.var(0).getText(), ctx.var(1).getText())
        if self.enableDynamic:
            updstr = constants.dyn_track_str[LIBRARYNAME].format(self.varMap[ctx.var(0).getText()],
                                                                 ctx.probability().getText(),
                                                                 ctx.FLOAT().getText())
            return statement_string + updstr
            
        else:
            updstr = "_ = {};_ = {};\n".format(ctx.probability().getText(),
                                        ctx.FLOAT().getText())                        
            return statement_string + updstr        

    def visitIfonly(self, ctx):
        str_if_only = "if {} != 0 {{\n {} }}\n"
        cond_var = ctx.var().getText()
        statement_string = ''
        for statement in ctx.ifs:
            translated = self.visit(statement)
            if translated is not None:
                statement_string += translated
            else:
                EXITWITHERROR("[Error] Unable to translate: ", statement.getText())
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
                EXITWITHERROR("[Error] Unable to translate: ", statement.getText())

        else_statement_string = ''
        for statement in ctx.elses:
            translated = self.visit(statement)
            if translated is not None:
                else_statement_string += translated
            else:
                EXITWITHERROR("[Error] Unable to translate: ", statement.getText())

        return str_if.format(cond_var, statement_string, else_statement_string)

    def getReadDynString(self, member):
        if member[1] == 0:
            return "DynMap[{}]".format(member[0])
        else:
            return "DynMap[{}+{}]".format(member[0], member[1])

    def visitRepeatlvar(self, ctx):
        repeatVar = ctx.var().getText()
        temp_var_name = "__temp_{}".format(self.tempvarnum)
        self.tempvarnum += 1

        statement_string = ''
        for statement in ctx.statement():
            translated = self.visit(statement)
            if translated is not None:
                statement_string += translated
            else:
                EXITWITHERROR("[Error] Unable to translate: ", statement.getText())

        str_for_loop = "for {} := 0; {} < {}; {}++ {{\n {} }}\n"
        return str_for_loop.format(temp_var_name, temp_var_name, repeatVar,
                                   temp_var_name, statement_string)

    def visitWhile(self, ctx):
        while_str = "for {} {{\n {} }}\n"
        condition = ctx.cond.getText()
        
        statement_string = ''
        for statement in ctx.statement():
            translated = self.visit(statement)
            if translated is not None:
                statement_string += translated
            else:
                EXITWITHERROR("[Error] Unable to translate: ", statement.getText())
        return while_str.format(condition, statement_string)            

    def visitRepeat(self, ctx):
        pre_string = ''
        repeatNum = ctx.INT().getText()
        temp_var_name = "__temp_{}".format(self.tempvarnum)
        self.tempvarnum += 1

        statement_string = ''
        for statement in ctx.statement():
            translated = self.visit(statement)
            if translated is not None:
                statement_string += translated
            else:
                EXITWITHERROR("[Error] Unable to translate: ", statement.getText())
                
        str_for_loop = pre_string + "for {} := 0; {} < {}; {}++ {{\n {} }}\n"
        return str_for_loop.format(temp_var_name, temp_var_name, repeatNum,
                                   temp_var_name, statement_string)

    def visitForloop(self, ctx):
        pre_string = ''

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
                EXITWITHERROR("[Error] Unable to translate: ", statement.getText())

        str_for_loop = pre_string + "for _, {} := range({}) {{\n {} }}\n"
        return str_for_loop.format(var_name, group_name, statement_string)

    def visitFunc(self, ctx):
        return ctx.getText() + ";\n"

    def visitInstrument(self, ctx):
        if self.args.instrument:
            return ctx.code.text[2:-2] + ";\n"
        else:
            return ""

    def visitSpeccheckarray(self, ctx):
        if not self.enableDynamic:
            return ""
        statement_string = ''
        checked_var = ctx.var().getText()
        checked_val = ctx.probability().getText()

        return statement_string + constants.ch_str.format(self.varMap[checked_var], checked_val,
                                                          self.arraySize[checked_var], checked_var)

    def isGroup(self, pid):
        if isinstance(pid, ParallelyParser.NamedpContext):
            return (False, pid.getText())
        elif isinstance(pid, ParallelyParser.VariablepContext):
            EXITWITHERROR("[Error] Cant handle process name variables")
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

                # only works for 1 dimentional so far
                d_str = ""
                if self.enableDynamic:
                    d_str = "LIBRARYNAME.InitDynArray({}, {}, DynMap[:]);\n".format(self.varMap[varname],
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
                EXITWITHERROR("Not supporting dynamic unbounded arrays")

            d_str = "LIBRARYNAME.InitDynArray(\"{}\", {}, DynMap);\n".format(varname,
                                                                             decl.GLOBALVAR()[0])

            if len(dim) == 1:
                return dyn_array_dec.format(varname, "[]", dectype[1], dim[0]) + d_str
            if len(dim) > 1:
                return dyn_array_dec.format(varname, "[]", dectype[1], dim[0]) + d_str
            else:
                EXITWITHERROR("[Error] Unable to translate: ", decl.getText())
        else:
            varname = decl.var().getText()
            self.primitiveTMap[varname] = dectype[0]
            self.typeMap[varname] = (dectype[1], dectype[2])

            if self.enableDynamic and dectype[0] == "dynamic":
                self.varMap[varname] = self.varNum
                self.varNum += 1
                self.dynsize += 1
                d_init_str = "var {0} {1};\n".format(varname, dectype[1])
                d_track_str = constants.dyn_init_str[LIBRARYNAME].format(self.varMap[varname])
                return d_init_str + d_track_str
            else:
                return str_single_dec.format(varname, dectype[1])

    def handleGroup(self, group_name, element_name, ctx):
        self.pid = "tid"
        print("Translating process group: ", group_name)

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
                EXITWITHERROR("[Error] Unable to translate: ", statement.getText())
        process_name = "func_" + group_name
        self.process_list.append(((process_name, group_name), 1))
        process_code = dec_string + statement_string
        process_def_str = constants.multiple_process_thread[LIBRARYNAME].format(process_name, self.dynsize,
                                                                                process_code, element_name)
        self.process_defs.append(process_def_str)

    def visitSingle(self, ctx):
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

        print("Translating process: ", self.pid)

        # Collect the declarations which should be at the top
        dec_string = ""
        for decl in ctx.declaration():
            dec_string += self.handleDec(decl)
            
        statement_string = ""
        for statement in ctx.statement():
            self.recovernum = 0
            translated = self.visit(statement)
            # print statement, translated
            if translated is not None:
                statement_string += translated
            else:
                EXITWITHERROR("[Error] Unable to translate: ", statement.getText())

        process_name = "func_" + self.pid
        self.process_list.append((process_name, 0))
        process_code = dec_string + statement_string
        process_def_str = constants.single_process_thread[LIBRARYNAME].format(process_name,
                                                                    self.dynsize, process_code, self.pid)
        # print "--------------------"
        # print process_def_str
        # print "--------------------"
        self.process_defs.append(process_def_str)

    def translate(self, tree, numthreads, proc_groups_in, fout_name, maintemplate, worktemplate):
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

        all_global_decs = ''.join(self.globaldecs)
        all_process_defs = '\n'.join(self.process_defs).replace('LIBRARYNAME', LIBRARYNAME)        

        run_procs = ''
        for fname, is_group in self.process_list:
            if is_group:
                run_procs = fname[0]
                template_str = open(worktemplate, 'r').readlines()
                with open("worker.go", "w") as fout:
                    for line in template_str:
                        newline = line.replace('__NUM_THREADS__', str(numthreads))
                        newline = newline.replace('__GLOBAL_DECS__', all_global_decs)
                        newline = newline.replace('__FUNC_DECS__', all_process_defs)
                        newline = newline.replace('__START__THREADS__', run_procs)
                        newline = newline.replace('LIBRARYNAME', LIBRARYNAME)
                        fout.write(newline)                
            else:
                run_procs = "{}();\n".format(fname)
                template_str = open(maintemplate, 'r').readlines()
                with open("main.go", "w") as fout:
                    for line in template_str:
                        newline = line.replace('__NUM_THREADS__', str(numthreads))
                        newline = newline.replace('__GLOBAL_DECS__', all_global_decs)
                        newline = newline.replace('__FUNC_DECS__', all_process_defs)
                        newline = newline.replace('__START__THREADS__', run_procs)
                        newline = newline.replace('LIBRARYNAME', LIBRARYNAME)
                        fout.write(newline)                

def main(program_str, outfile, filename, maintemplate, worktemplate, debug, dynamic, args):
    print("Starting the cross compilation")
    input_stream = InputStream(program_str)
    lexer = ParallelyLexer(input_stream)
    stream = CommonTokenStream(lexer)
    parser = ParallelyParser(stream)

    tree = parser.parallelprogram()

    threadcounter = CountThreads()
    threadcounter.visit(tree)

    print("Number of processes found: {}".format(threadcounter.processcount))

    translator = Translator(dynamic, args)
    translator.translate(tree, threadcounter.processcount, threadcounter.processes,
                         outfile, maintemplate, worktemplate)


if __name__ == '__main__':
    parser = ArgumentParser()
    parser.add_argument("-f", dest="programfile",
                        help="File containing the code", required=True)
    parser.add_argument("-o", dest="outfile",
                        help="File to output the sequential code")
    parser.add_argument("-d", "--debug", action="store_true",
                        help="Print debug info")
    parser.add_argument("-tm", dest="maintemplate",
                        help="File containing the template for main")
    parser.add_argument("-tw", dest="worktemplate",
                        help="File containing the template for worker")
    # parser.add_argument("-t", dest="template",
    #                     help="File containing the template")    
    parser.add_argument("-dyn", "--dynamic", action="store_true",
                        help="Enable dynamic tracking")
    parser.add_argument("-a", "--arrayO1", action="store_true",
                        help="Inline tracking")
    parser.add_argument("-i", "--instrument", action="store_true",
                        help="Add instrumentation")
    parser.add_argument("-n", "--noisy", action="store_true",
                        help="Use the noisy channel function")
    parser.add_argument("-acc", "--accuracy", action="store_true", default=False,
                        help="Enable dynamic tracking of accuracy")
    parser.add_argument("-rel", "--reliability", action="store_true", default=False,
                        help="Enable dynamic tracking of accuracy")        
    args = parser.parse_args()


    # Plan is to combine the other version with this later
    # if args.distributed:
    #     print("[crosscompiler] Using the distributed runtime")
    #     if args.maintemplate is None or args.worktemplate is None:
    #         print("[crosscompiler] Please provide two templates -tw and -tm")
    #         exit(-1)
    # else:
    #     print("[crosscompiler] Using the shared memory runtime")
    #     if args.template is None:
    #         print("[crosscompiler] Please provide the template file name with -t")
    #         exit(-1)

    if args.reliability and args.accuracy:
        LIBRARYNAME = "dieseldist"
        print("[crosscompiler] Tracking both accuracy and reliability")
    elif args.reliability:
        LIBRARYNAME = "dieseldistrel"
        print("[crosscompiler] Tracking only reliability")
    elif args.accuracy:
        LIBRARYNAME = "dieseldistacc"
        print("[crosscompiler] Tracking only accuracy")
    else:
        print("[crosscompiler] Defaulting to tracking both accuracy and reliability")

    if args.dynamic:
        print("[crosscompiler] Enabling dynamic tracking")
    if args.arrayO1:
        print("[crosscompiler] Enabling array optimization: Send one value")
    if args.instrument:
        print("[crosscompiler] Enabling instrumentation")

    programfile = open(args.programfile, 'r')
    # outfile = open(args.outfile, 'w')
    program_str = programfile.read()

    startTime = time.time()
    main(program_str, args.outfile, programfile.name, args.maintemplate, args.worktemplate,
         args.debug, args.dynamic, args)
    print("[crosscompiler] Done!")
    print("[crosscompiler] Elapsed time : ", time.time()-startTime)
