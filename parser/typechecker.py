# flake8: noqa E501

from ParallelyVisitor import ParallelyVisitor
from ParallelyParser import ParallelyParser

key_error_msg = "Type error detected: Undeclared variable (probably : {})"


class parallelyTypeChecker(ParallelyVisitor):
    def __init__(self, debug):
        self.typecontext = {}
        self.processgroups = {}
        self.debug = debug
        self.globaltypecontext = {}

    def baseTypesEqual(self, type1, type2, ctx):
        # Literals can take any type. for ease of use.
        if not (type1[1][:3] == type2[1][:3]):
            if(type1[2] == 2):
                return(type2[0], type1[1], 0)
            elif (type2[2] == 2):
                return(type1[0], type2[1], 0)
            else:
                print "Type error : ", ctx.getText(), type1[1], type2[1]
                exit(-1)
        else:
            if type1[0] == 'approx' or type2[0] == 'approx':
                return ('approx', type1[1])
            return type1

    def boolBaseTypesEqual(self, type1, type2, ctx):
        if(type1[2] == 2):
            return(type2[0], 'int32', 0)
        elif (type2[2] == 2):
            return(type1[0], 'int32', 0)
        elif type1[0] == 'approx' or type2[0] == 'approx':
            return ('approx', 'int32', 0)
        elif type1[1] == type2[1]:
            return (type1[0], type2[1], 0)
        else:
            print "Incompatible types in expressions : ", type1, type2, ctx.getText()
            exit(-1)

    def resultType(self, type1, type2, ctx):
        if(type1[2] == 1 or type2[2] == 1):
            print "Array types cannot occur in exporessions : ", type1, type2
            exit(-1)

        # Literals can take either type
        if(type1[2] == 2):
            return (type2[0], type2[1], 0)
        elif (type2[2] == 2):
            return (type1[0], type1[1], 0)
        elif (type1[0] == type2[0]):
            return (type1[0], type1[1], 0)
        elif (type1[0] == 'approx' or type2[0] == 'approx'):
            return ('approx', type1[1], 0)
        elif (type1[0] == 'dynamic' and type2[0] == 'dynamic'):
            return ('dynamic', type1[1], 0)
        else:
            print "Incompatible types in expressions : ", type1, type2, ctx.getText()
            exit(-1)

    ########################################
    # Expression type checking
    ########################################
    def visitLiteral(self, ctx):
        # print ("At expression Literal Int")
        # return (ParallelyLexer.PRECISETYPE, ParallelyLexer.INTTYPE)
        return ("precise", "int", 2)

    def visitFliteral(self, ctx):
        return ("precise", "float64", 2)

    def visitVariable(self, ctx):
        return self.typecontext[ctx.getText()]

    def visitVar(self, ctx):
        return self.typecontext[ctx.getText()]

    def visitGlobalvariable(self, ctx):
        return self.typecontext[ctx.getText()]

    def visitArrayvar(self, ctx):
        for expr in ctx.expression():
            expr_type = self.visit(expr)
            if expr_type[0] != 'precise':
                print ("[Error] Array indexes have to be precise :", ctx.getText())
                return False
            if expr_type[1] != 'int64' and expr_type[1] != 'int32':
                print ("[Error] Array indexes have to be int :", ctx.getText())
                return False

        array_qual, array_type, _ = self.typecontext[ctx.var().getText()]
        element_type = (array_qual, array_type, 0)
        return element_type

    # def visitVar(self, ctx):
    #     return self.typecontext[ctx.getText()]

    def visitSelect(self, ctx):
        return self.visit(ctx.expression())

    def visitMultiply(self, ctx):
        # Not checking if arraytype
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        res = self.resultType(type1, type2, ctx)
        # res = self.baseTypesEqual(type1, type2, ctx)
        return res

    def visitAdd(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        res = self.resultType(type1, type2, ctx)
        return res

    def visitMinus(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        # print ("[Debug] At expression minus", type1, type2)
        res = self.baseTypesEqual(type1, type2, ctx)
        return res

    def visitDivide(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        res = self.resultType(type1, type2, ctx)
        return res

    def visitProb(self, ctx):
        type1 = self.visit(ctx.expression(0))[0]
        type2 = self.visit(ctx.expression(1))[0]
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

    def visitGeq(self, ctx):
        type1 = self.visit(ctx.expression(0))
        type2 = self.visit(ctx.expression(1))
        return self.boolBaseTypesEqual(type1, type2, ctx)

    def visitLeq(self, ctx):
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
    # def getType(self, fulltype):
    #     if isinstance(fulltype, ParallelyParser.SingletypeContext):
    #         return (fulltype.basictype().typequantifier().getText(),
    #                 fulltype.basictype().getChild(1).getText(), 0)
    #     elif isinstance(fulltype, ParallelyParser.ArraytypeContext):
    #         return (fulltype.basictype().typequantifier().getText(),
    #                 fulltype.basictype().getChild(1).getText(), 1)
    #     else:
    #         print "[Error] Unknown type : ", fulltype
    #         exit(-1)
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
        # elif isinstance(fulltype, ParallelyParser.GlobalarrayContext):
        #     return (fulltype.basictype().typequantifier().getText(),
        #             fulltype.basictype().getChild(1).getText(), 1)
        else:
            print "[Error] Unknown type : ", fulltype.getText()
            exit(-1)

    def visitSingledeclaration(self, ctx):
        decl_type_q, decl_type_t, _ = self.getType(ctx.basictype())
        self.typecontext[ctx.var().getText()] = (decl_type_q, decl_type_t, 0)

    def visitArraydeclaration(self, ctx):
        decl_type_q, decl_type_t, _ = self.getType(ctx.basictype())
        self.typecontext[ctx.var().getText()] = (decl_type_q, decl_type_t, 1)

    def visitGlobalarray(self, ctx):
        decl_type_q, decl_type_t, _ = self.getType(ctx.basictype())
        self.globaltypecontext[ctx.GLOBALVAR().getText()] = (decl_type_q, decl_type_t, 1)

    def visitGlobalconst(self, ctx):
        decl_type_q, decl_type_t, _ = self.getType(ctx.basictype())
        self.globaltypecontext[ctx.GLOBALVAR().getText()] = (decl_type_q, decl_type_t, 0)

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

    def isAllowedFlow(self, type1, type2):
        # print "===========================>", type1, type2, type1 == 'dynamic'
        if (type1 == type2) or (type1 == 'approx'):
            return True
        elif (type1 == 'dynamic') and (type2 == 'precise'):
            return True
        else:
            print "Invalid flow {}<-{}".format(type1, type2)
            return False

    def visitArrayload(self, ctx):
        var_type = self.typecontext[ctx.var(0).getText()][0]
        array_type = self.typecontext[ctx.var(1).getText()][0]

        for expr in ctx.expression():
            expr_type = self.visit(expr)
            if expr_type[0] != 'precise':
                print ("Type error :", expr, expr_type[0])
                exit(-1)

        # Deadline day
        if self.isAllowedFlow(var_type, array_type):
            return True
        else:
            print ("Type error :", ctx.getText())
            exit(-1)

    def visitExpassignment(self, ctx):
        # print ctx.getText()
        # print ctx.expression().getText()
        var_type = self.typecontext[ctx.var().getText()]
        expr_type = self.visit(ctx.expression())
        if self.isAllowedFlow(var_type[0], expr_type[0]):
            if var_type[1] != expr_type[1]:
                print "[Error]: {}!={} ({})".format(var_type[1], expr_type[1], ctx.getText())
                return False

            # If literal allow to be trated as any type
            if expr_type[2] == 2:
                return True
            # If not literal types have to match
            elif var_type[2] == expr_type[2]:
                return True
            else:
                print "[Error]: possible array type ({})".format(ctx.getText())
                exit(-1)
        else:
            print "Type Error : {}, {}, {}".format(ctx.getText(),
                                                   var_type, expr_type)
            exit(-1)

    def visitBoolassignment(self, ctx):
        var_type = self.typecontext[ctx.var().getText()][0]
        expr_type = self.visit(ctx.boolexpression())
        if (var_type == expr_type):
            return True
        if (var_type[1] == expr_type[1]) and (var_type[0] == 'approx'):
            return True
        else:
            print "Type Error : {}, {}, {}".format(ctx.getText(),
                                                   var_type, expr_type)
            exit(-1)

    def visitIf(self, ctx):
        print "[DebugIf] ", ctx.getText(), ctx.ifs        
        guardtype = self.visit(ctx.getChild(1))
        if guardtype[0] != ('precise'):
            print "Type Error precise boolean expected. ", ctx.getText()
            exit(-1)
        for statement in ctx.ifs:
            typechecked = self.visit(statement)
            if not typechecked:
                print "[Error] failed to type check: {}".format(statement.getText())
                exit(-1)
            if self.debug:
                print "[Debug] ", statement.getText(), typechecked
        for statement in ctx.elses:
            typechecked = self.visit(statement)
            if not typechecked:
                print "[Error] failed to type check: {}".format(statement.getText())
                exit(-1)
            if self.debug:
                print "[Debug] ", statement.getText(), typechecked                
        return True
    
    def visitIfOnly(self, ctx):
        print "[Debug] ", ctx.getText(), ctx.ifs        
        guardtype = self.visit(ctx.getChild(1))
        if guardtype[0] != ('precise'):
            print "Type Error precise boolean expected. ", ctx.getText()
            exit(-1)
        for statement in ctx.ifs:
            typechecked = self.visit(statement)
            if not typechecked:
                print "[Error] failed to type check: {}".format(statement.getText())
                exit(-1)
            if self.debug:
                print "[Debug] ", statement.getText(), typechecked
        return True

    def visitSend(self, ctx):
        # At some point check if the first element is a pid
        var_type = self.typecontext[ctx.getChild(6).getText()]
        sent_type = self.getType(ctx.getChild(4))

        if var_type == sent_type:
            return True
        else:
            print "Type Error : {}".format(ctx.getText())
            exit(-1)

    def visitReceive(self, ctx):
        # At some point check if the first element is a pid
        var_type = self.typecontext[ctx.getChild(0).getText()]
        rec_type = self.getType(ctx.getChild(6))
        print var_type, rec_type
        # rec_type = ctx.getChild(6).getChild(1).getText()
        if var_type == rec_type:
            return True
        else:
            print "Type Error : {}".format(ctx.getText())
            exit(-1)

    def visitCondsend(self, ctx):
        variables = ctx.var()
        guard = self.typecontext[variables[0].getText()]
        var_type = self.typecontext[variables[1].getText()]

        if guard[0] != 'approx':
            err = "Type Error : {} has to be approx ({})".format(variables[0].getText(),
                                                                 ctx.getText())
            print err
            exit(-1)

        if var_type[0] != 'approx':
            err = "Type Error : {} has to be approx ({})".format(variables[1].getText(),
                                                                 ctx.getText())
            print err
            exit(-1)

        rec_type = self.getType(ctx.fulltype())
        if var_type == rec_type:
            return True
        else:
            print "[Error] Type Error (Conditional Receive): {}".format(ctx.getText())
            exit(-1)

    def visitCondreceive(self, ctx):
        # At some point check if the first element is a pid
        variables = ctx.var()
        signal = self.typecontext[variables[0].getText()]
        var_type = self.typecontext[variables[1].getText()]

        if signal[0] != 'approx':
            err = "[Error] Type Error : {} has to be approx. ({})".format(
                variables[0].getText(), ctx.getText())
            print err
            exit(-1)

        if var_type[0] != 'approx':
            err = "[Error] Type Error : {} has to be approx".format(
                variables[1].getText())
            print err
            exit(-1)

        rec_type = self.getType(ctx.fulltype())
        if var_type == rec_type:
            return True
        else:
            print "[Error] Type Error (Conditional Send): {}".format(ctx.getText())
            exit(-1)

    def visitDynsend(self, ctx):
        variables = ctx.var()
        var_type = self.typecontext[variables.getText()]

        if var_type[0] != 'dynamic':
            err = "Type Error : {} has to be dynamic ({})".format(variables.getText(), ctx.getText())
            print err
            return False

        sent_type = self.getType(ctx.fulltype())

        if var_type == sent_type:
            return True
        else:
            print "Type Error : {}!={} ({})".format(var_type, sent_type, ctx.getText())
            exit(-1)

    def visitDynreceive(self, ctx):
        # At some point check if the first element is a pid
        variables = ctx.var()
        var_type = self.typecontext[variables.getText()]

        if var_type[0] != 'dynamic':
            err = "[Error] Type Error : {} has to be dynamic".format(variables.getText())
            print err
            exit(-1)

        rec_type = self.getType(ctx.fulltype())
        if var_type == rec_type:
            return True
        else:
            print "[Error] Type Error (Conditional Send): {}".format(ctx.getText())
            exit(-1)

    def visitDyncondsend(self, ctx):
        variables = ctx.var()
        guard = self.visit(variables[0])
        var_type = self.visit(variables[1])

        if guard[0] != 'dynamic':
            err = "Type Error : {} has to be dynamic ({})".format(variables[0].getText(),
                                                                  ctx.getText())
            print err
            exit(-1)

        if var_type[0] != 'dynamic':
            err = "Type Error : {} has to be dynamic ({})".format(variables[1].getText(),
                                                                  ctx.getText())
            print err
            exit(-1)

        sent_type = self.getType(ctx.fulltype())
        if var_type == sent_type:
            return True
        else:
            print "Type Error {}!={} : ({})".format(var_type, sent_type, ctx.getText())
            exit(-1)

    def visitDyncondreceive(self, ctx):
        # At some point check if the first element is a pid
        variables = ctx.var()
        signal = self.typecontext[variables[0].getText()]
        var_type = self.typecontext[variables[1].getText()]

        if signal[0] != 'dynamic':
            err = "[Error] Type Error : {} has to be dynamic. ({})".format(
                variables[0].getText(), ctx.getText())
            print err
            exit(-1)

        if var_type[0] != 'dynamic':
            err = "[Error] Type Error : {} has to be dynamic".format(
                variables[1].getText())
            print err
            exit(-1)

        rec_type = self.getType(ctx.fulltype())
        if var_type == rec_type:
            return True
        else:
            print "[Error] Type Error (Conditional Send): {}".format(ctx.getText())
            exit(-1)
                
    def visitCast(self, ctx):
        type1 = self.visit(ctx.var(0))
        type2 = self.getType(ctx.fulltype())
        if type1 == type2:
            return True
        else:
            print "[Error] Type Error: {}".format(ctx.getText())            
            exit(-1)

    def visitTrack(self, ctx):
        type1 = self.visit(ctx.var(0))
        type2 = self.visit(ctx.var(1))

        if type1[0] != 'dynamic':
            print "[Error] Track must be assigned to a dynamic variable: {}".format(ctx.getText())
            exit(-1)
        elif type1[1] != type2[1]:
            print "[Error] Different types: {}".format(ctx.getText())
            exit(-1)
        else:
            return True

    def visitCheck(self, ctx):
        type1 = self.visit(ctx.var(0))
        type2 = self.visit(ctx.var(1))

        if type2[0] != 'dynamic':
            print "[Error] check must have a dynamic variable: {}".format(ctx.getText())
            exit(-1)
        if type1[0] != 'approx':
            print "[Error] check must assign to an approx variable: {}".format(ctx.getText())
            exit(-1)
        elif type1[1] != type2[1]:
            print "[Error] Different types: {}".format(ctx.getText())
            exit(-1)
        else:
            return True

    def visitForloop(self, ctx):
        temp_typecontext = dict(self.typecontext)

        # Binding process id to an int
        self.typecontext[ctx.VAR().getText()] = ("precise", "int32", 0)

        all_typechecked = True
        for statement in ctx.statement():
            print statement.getText()
            typechecked = self.visit(statement)
            if not typechecked:
                print "[Error] failed to type check: {}".format(statement.getText())
                exit(-1)
            if self.debug:
                print "[Debug] ", statement.getText(), typechecked
            all_typechecked = typechecked and all_typechecked

        self.typecontext = dict(temp_typecontext)
        return all_typechecked

    def visitRepeat(self, ctx):
        all_typechecked = True
        for statement in ctx.statement():
            typechecked = self.visit(statement)
            if self.debug:
                print "[Debug] ", statement.getText(), typechecked
            all_typechecked = typechecked and all_typechecked
        return all_typechecked

    def visitWhile(self, ctx):
        guardtype = self.visit(ctx.expression())
        if guardtype[0] != ('precise'):
            print "[Error] while loops need precise guards.", ctx.getText()
            exit(-1)
        all_typechecked = True
        for statement in ctx.statement():
            typechecked = self.visit(statement)
            if self.debug:
                print "[Debug] ", statement.getText(), typechecked
            all_typechecked = typechecked and all_typechecked
        return all_typechecked

    def visitInstrument(self, ctx):
        return True

    def visitRepeatvar(self, ctx):
        # var_type = self.typecontext[ctx.GLOBALVAR().getText()]
        # if not var_type[0] == 'precise':
        #     print "Type error: only precise int allowed in a repeat statement: ", ctx.getText()
        #     exit(-1)
        all_typechecked = True
        for statement in ctx.statement():
            typechecked = self.visit(statement)
            if self.debug:
                print "[Debug] ", statement.getText(), typechecked
            all_typechecked = typechecked and all_typechecked
        return all_typechecked

    def visitSingle(self, ctx):
        self.typecontext = dict(self.globaltypecontext)
        print self.typecontext
        pid = ctx.processid().getText()

        # Process Id is an int accesible by the program.
        # Doesnt have to be. Simplifies code for now.
        if isinstance(ctx.processid(), ParallelyParser.GroupedpContext):
            self.typecontext[ctx.processid().VAR().getText()] = ("precise", "int32", 0)
        else:
            self.typecontext[ctx.processid().getText()] = ("precise", "int32", 0)

        try:
            for declaration in ctx.declaration():
                self.visit(declaration)
        except Exception, e:
            print "Type error in declarations : ", pid
            print e
            exit(-1)

        all_typechecked = True
        temp_st = ''
        try:
            for statement in ctx.statement():
                temp_st = statement.getText()
                if self.debug:
                    print "[Debug - checking] ", statement.getText()
                typechecked = self.visit(statement)
                if not typechecked:
                    print "[ERROR] failed type checker: ", statement.getText()

                if self.debug:
                    print "[Debug - checked] ", statement.getText(), typechecked
                all_typechecked = typechecked and all_typechecked
        except KeyError, keyerror:
            print "[Error] Undeclared variable: ", temp_st[:20], keyerror
            exit(0)

        if not all_typechecked:
            print "Process {} failed typechecker".format(pid)
        self.typecontext = {}
        return all_typechecked

    

    def visitParcomposition(self, ctx):
        # Does nothing for now.
        # Only sets of procs allowed in these declarations
        if ctx.globaldec():
            for dec in ctx.globaldec():
                self.visit(dec)

        all_type_checked = True
        for current_program in ctx.program():
            type_checked = self.visit(current_program)
            all_type_checked = type_checked and all_type_checked

        return all_type_checked
