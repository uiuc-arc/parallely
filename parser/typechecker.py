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

    def debugMsg(self, msg):
        if self.debug:
            print("[Debug - TypeChecker] " + msg)

    def exitWithError(self, msg):
        print("[Error - TypeChecker]: " + msg)
        exit(-1)

    def baseTypesEqual(self, type1, type2, ctx):
        # Literals can take any type. for ease of use.
        if not (type1[1][:3] == type2[1][:3]):
            if(type1[2] == 2):
                return(type2[0], type1[1], 0)
            elif (type2[2] == 2):
                return(type1[0], type2[1], 0)
            else:
                self.exitWithError("{} != {} ({})".format(type1[1], type2[1], ctx.getText()))
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
            self.exitWithError("{} != {} ({})".format(type1[1], type2[1], ctx.getText()))

    def resultType(self, type1, type2, ctx):
        if(type1[2] == 1 or type2[2] == 1):
            self.exitWithError(
                "Array types in exporessions ({} or {})".format(type1, type2, ctx.getText()))

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
        # TODO: ugly
        elif ((type1[0] == 'dynamic' and type2[0] == 'precise') or
              (type1[0] == 'precise' and type2[0] == 'dynamic')):
            return ('dynamic', type1[1], 0)
        else:
            self.exitWithError("{} != {} ({})".format(type1[1], type2[1], ctx.getText()))

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
        # print "--------------", self.typecontext[ctx.getText()]
        return self.typecontext[ctx.getText()]

    def visitVar(self, ctx):
        # print "--------------", self.typecontext[ctx.getText()]
        return self.typecontext[ctx.getText()]

    def visitLocalvariable(self, ctx):
        return self.typecontext[ctx.getText()]

    def visitGlobalvariable(self, ctx):
        return self.typecontext[ctx.getText()]

    def visitArrayvar(self, ctx):
        for expr in ctx.expression():
            expr_type = self.visit(expr)
            if expr_type[0] != 'precise':
                self.exitWithError("Array indexes have to be precise: ({})".format(ctx.getText()))
            if expr_type[1] != 'int64' and expr_type[1] != 'int32':
                self.exitWithError("Array indexes have to be int: ({})".format(ctx.getText()))

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
            self.exitWithError("{} != {} ({})".format(type1[1], type2[1], ctx.getText()))
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
            self.exitWithError("Unknown type - {}".format(fulltype.getText()))

    def visitSingledeclaration(self, ctx):
        decl_type_q, decl_type_t, _ = self.getType(ctx.basictype())
        self.typecontext[ctx.var().getText()] = (decl_type_q, decl_type_t, 0)

    def visitGlobalexternal(self, ctx):
        decl_type_q, decl_type_t, _ = self.getType(ctx.basictype())
        self.globaltypecontext[ctx.GLOBALVAR().getText()] = (decl_type_q, decl_type_t, 0)

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
        type1 = self.visit(ctx.getChild(0))
        type2 = self.visit(ctx.getChild(2))
        return (type1 and type2)

    # Removed blocks from the grammer. Still keeping this here
    def visitBlock(self, ctx):
        return self.visit(ctx.getChild(1))

    def isAllowedFlow(self, type1, type2):
        if (type1 == type2) or (type1 == 'approx'):
            return True
        elif (type1 == 'dynamic') and (type2 == 'precise'):
            return True
        else:
            return False

    def visitArrayload(self, ctx):
        var_type = self.typecontext[ctx.var(0).getText()][0]
        array_type = self.typecontext[ctx.var(1).getText()][0]

        for expr in ctx.expression():
            expr_type = self.visit(expr)
            if expr_type[0] != 'precise':
                self.exitWithError("Array indexes need to be precise: {}".format(expr.getText()))

        # Deadline day
        if self.isAllowedFlow(var_type, array_type):
            return True
        else:
            self.exitWithError("Invalid flow {}<-{}".format(var_type, array_type))

    def visitExpassignment(self, ctx):
        # print ctx.getText()
        # print "**************", ctx.expression().getText()
        var_type = self.typecontext[ctx.var().getText()]
        expr_type = self.visit(ctx.expression())
        if self.isAllowedFlow(var_type[0], expr_type[0]):
            if var_type[1] != expr_type[1]:
                self.exitWithError("{} != {} ({})".format(var_type[1], expr_type[1], ctx.getText()))

            # If literal allow to be trated as any type
            if expr_type[2] == 2:
                return True
            # If not literal types have to match
            elif var_type[2] == expr_type[2]:
                return True
            else:
                self.exitWithError("{} != {} ({})".format(var_type, expr_type, ctx.getText()))
                # print "[Error]: possible array type ({})".format(ctx.getText())
        else:
            self.exitWithError("Invalid flow {}<-{} ({})".format(var_type, expr_type, ctx.getText()))

    def visitIf(self, ctx):
        self.debugMsg("If: {}" + ctx.getText())
        guardtype = self.visit(ctx.var())
        # print self.typecontext, guardtype, ctx.var().getText(), type(ctx.var())
        if guardtype[0] != ('precise'):
            self.exitWithError("Boolean guards have to be precise: {}".format(ctx.getText()))
        for statement in ctx.ifs:
            typechecked = self.visit(statement)
            if not typechecked:
                self.exitWithError("Failed to typecheck: {}".format(statement.getText()))
            self.debugMsg("If branch - statement - {}".format(statement.getText()))
        for statement in ctx.elses:
            typechecked = self.visit(statement)
            if not typechecked:
                self.exitWithError("Failed to typecheck: {}".format(statement.getText()))
            self.debugMsg("Else branch - statement - {}".format(statement.getText()))
        return True

    def visitIfonly(self, ctx):
        guardtype = self.visit(ctx.getChild(1))
        if guardtype[0] != ('precise'):
            self.exitWithError("Boolean guards have to be precise: {}".format(ctx.getText()))
        for statement in ctx.ifs:
            typechecked = self.visit(statement)
            if not typechecked:
                self.exitWithError("Failed to typecheck: {}".format(statement.getText()))
            self.debugMsg("If branch - statement - {}".format(statement.getText()))
        return True

    def visitSend(self, ctx):
        # At some point check if the first element is a pid
        var_type = self.typecontext[ctx.getChild(6).getText()]
        sent_type = self.getType(ctx.getChild(4))

        if var_type == sent_type:
            return True
        else:
            self.exitWithError("{} != {} ({})".format(var_type, sent_type, ctx.getText()))

    def visitReceive(self, ctx):
        # At some point check if the first element is a pid
        var_type = self.typecontext[ctx.getChild(0).getText()]
        rec_type = self.getType(ctx.getChild(6))
        # rec_type = ctx.getChild(6).getChild(1).getText()
        if var_type == rec_type:
            return True
        else:
            self.exitWithError("{} != {} ({})".format(rec_type, var_type, ctx.getText()))            

    def visitCondsend(self, ctx):
        variables = ctx.var()
        guard = self.typecontext[variables[0].getText()]
        var_type = self.typecontext[variables[1].getText()]

        if guard[0] != 'approx':
            self.exitWithError("Condsend guard has to be approx {} ({})".format(guard, ctx.getText()))
        if var_type[0] != 'approx':
            self.exitWithError("Condsend data has to be approx {} ({})".format(var_type, ctx.getText()))

        rec_type = self.getType(ctx.fulltype())
        if var_type == rec_type:
            return True
        else:
            self.exitWithError("{} != {} ({})".format(rec_type, var_type, ctx.getText()))

    def visitCondreceive(self, ctx):
        # At some point check if the first element is a pid
        variables = ctx.var()
        signal = self.typecontext[variables[0].getText()]
        var_type = self.typecontext[variables[1].getText()]

        if signal[0] != 'approx':
            self.exitWithError("Condrec signal has to be approx {} ({})".format(signal, ctx.getText()))            
        if var_type[0] != 'approx':
            self.exitWithError("Condsend data has to be approx {} ({})".format(var_type, ctx.getText()))

        rec_type = self.getType(ctx.fulltype())
        if var_type == rec_type:
            return True
        else:
            self.exitWithError("{} != {} ({})".format(rec_type, var_type, ctx.getText()))

    def visitDynsend(self, ctx):
        variables = ctx.var()
        var_type = self.typecontext[variables.getText()]

        if var_type[0] != 'dynamic':
            self.exitWithError("Dynsend type not dynamic {} ({})".format(var_type, ctx.getText()))               

        sent_type = self.getType(ctx.fulltype())
        if var_type == sent_type:
            return True
        else:
            self.exitWithError("{} != {} ({})".format(rec_type, var_type, ctx.getText()))

    def visitDynreceive(self, ctx):
        # At some point check if the first element is a pid
        variables = ctx.var()
        var_type = self.typecontext[variables.getText()]

        if var_type[0] != 'dynamic':
            self.exitWithError("Dynrec type not dynamic {} ({})".format(var_type, ctx.getText()))               

        rec_type = self.getType(ctx.fulltype())
        if var_type == rec_type:
            return True
        else:
            self.exitWithError("{} != {} ({})".format(rec_type, var_type, ctx.getText()))

    def visitDyncondsend(self, ctx):
        variables = ctx.var()
        guard = self.visit(variables[0])
        var_type = self.visit(variables[1])
        
        if guard[0] != 'dynamic':
            self.exitWithError("Dynsend guard has to be dynamic {} ({})".format(guard, ctx.getText()))
        if var_type[0] != 'dynamic':
            self.exitWithError("Dynrec type not dynamic {} ({})".format(var_type, ctx.getText()))     

        sent_type = self.getType(ctx.fulltype())
        if var_type == sent_type:
            return True
        else:
            self.exitWithError("{} != {} ({})".format(rec_type, var_type, ctx.getText()))

    def visitDyncondreceive(self, ctx):
        # At some point check if the first element is a pid
        variables = ctx.var()
        signal = self.typecontext[variables[0].getText()]
        var_type = self.typecontext[variables[1].getText()]

        if signal[0] != 'dynamic':
            self.exitWithError("Dynrec signal has to be dynamic {} ({})".format(signal, ctx.getText()))
        if var_type[0] != 'dynamic':
            self.exitWithError("Dynrec type not dynamic {} ({})".format(var_type, ctx.getText()))     

        rec_type = self.getType(ctx.fulltype())
        if var_type == rec_type:
            return True
        else:
            self.exitWithError("{} != {} ({})".format(rec_type, var_type, ctx.getText()))

    def visitCast(self, ctx):
        type1 = self.visit(ctx.var(0))
        type2 = self.getType(ctx.fulltype())
        if type1 == type2:
            return True
        else:
            self.exitWithError("{} != {} ({})".format(type1, type2, ctx.getText()))

    def visitTrack(self, ctx):
        type1 = self.visit(ctx.var(0))
        type2 = self.visit(ctx.var(1))

        if type1[0] != 'dynamic':
            self.exitWithError("Has to be dynamic {} ({})".format(type1, ctx.getText()))            
        elif type1[1] != type2[1]:
            self.exitWithError("{} != {} ({})".format(type1, type2, ctx.getText()))
        else:
            return True

    def visitCheck(self, ctx):
        type1 = self.visit(ctx.var(0))
        type2 = self.visit(ctx.var(1))

        if type2[0] != 'dynamic':
            self.exitWithError("Has to be dynamic {} ({})".format(type2, ctx.getText()))
        if type1[0] != 'approx':
            self.exitWithError("check must assign to an approx variable {} ({})".format(type2, ctx.getText()))
        elif type1[1] != type2[1]:
            self.exitWithError("{} != {} ({})".format(type1, type2, ctx.getText()))
        else:
            return True

    def visitForloop(self, ctx):
        temp_typecontext = dict(self.typecontext)

        # Binding process id to an int
        self.typecontext[ctx.VAR().getText()] = ("precise", "int32", 0)

        all_typechecked = True
        for statement in ctx.statement():
            typechecked = self.visit(statement)
            if not typechecked:
                self.exitWithError("Failed to type check: ({})".format(statement.getText()))
            self.debugMsg("Inside for loop : " + statement.getText())
            all_typechecked = typechecked and all_typechecked

        self.typecontext = dict(temp_typecontext)
        return all_typechecked

    def visitRepeat(self, ctx):
        all_typechecked = True
        for statement in ctx.statement():
            typechecked = self.visit(statement)
            self.debugMsg("Inside repeat = ({}) : {}".format(statement.getText(), typechecked))             
            if not typechecked:
                self.exitWithError("Failed to type check: ({})".format(statement.getText()))
            all_typechecked = typechecked and all_typechecked
        return all_typechecked

    def visitRepeatlvar(self, ctx):
        # Need to check if the lvar is precise
        guardtype = self.visit(ctx.var())
        if guardtype[0] != ('precise'):
            self.exitWithError("repeat guard has to be precise {} ({})".format(guardtype, ctx.getText()))

        all_typechecked = True
        for statement in ctx.statement():
            typechecked = self.visit(statement)
            if not typechecked:
                self.exitWithError("Failed to type check inside repeat: ({})".format(statement.getText()))
            self.debugMsg("Inside repeatlvar = ({}) : {}".format(statement.getText(), typechecked))             
            all_typechecked = typechecked and all_typechecked
        return all_typechecked

    def visitWhile(self, ctx):
        guardtype = self.visit(ctx.expression())
        if guardtype[0] != ('precise'):
            self.exitWithError("while guard has to be precise {} ({})".format(guardtype, ctx.getText()))
        all_typechecked = True
        for statement in ctx.statement():
            typechecked = self.visit(statement)
            if not typechecked:
                self.exitWithError("Failed to type check: ({})".format(statement.getText()))            
            self.debugMsg("Inside while = ({}) : {}".format(statement.getText(), typechecked))
            all_typechecked = typechecked and all_typechecked
        return all_typechecked

    # At some point may need to check against function spec
    def visitFunc(self, ctx):
        return True

    def visitInstrument(self, ctx):
        return True

    # Need to check that only dynamic type go here
    def visitSpeccheck(self, ctx):
        for var in ctx.rel_factor:
            vartype = self.visit(var)
            if vartype[0] != ('dynamic'):
                self.exitWithError("Has to be dynamic {} ({})".format(vartype, ctx.getText()))
        return True

    # Need to check that only dynamic type go here
    def visitSpeccheckarray(self, ctx):
        vartype = self.visit(ctx.var())
        if vartype[0] != ('dynamic'):
            self.exitWithError("Has to be dynamic {} ({})".format(vartype, ctx.getText()))            
        return True

    def visitRepeatvar(self, ctx):
        # var_type = self.typecontext[ctx.GLOBALVAR().getText()]
        # if not var_type[0] == 'precise':
        #     print "Type error: only precise int allowed in a repeat statement: ", ctx.getText()
        #     exit(-1)
        all_typechecked = True
        for statement in ctx.statement():
            typechecked = self.visit(statement)
            if not typechecked:
                self.exitWithError("Failed to type check: ({})".format(statement.getText()))            
            self.debugMsg("Inside while = ({}) : {}".format(statement.getText(), typechecked))
            all_typechecked = typechecked and all_typechecked
        return all_typechecked

    def visitSingle(self, ctx):
        self.typecontext = dict(self.globaltypecontext)
        # print self.typecontext
        pid = ctx.processid().getText()

        print("Type checking: process {}".format(pid))

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
            self.exitWithError("Type error in declarations ({}): {}".format(pid, ctx))

        all_typechecked = True
        temp_st = ''
        try:
            for statement in ctx.statement():
                temp_st = statement.getText()
                self.debugMsg("[Debug - checking]: "+ statement.getText())
                typechecked = self.visit(statement)
                if not typechecked:
                    self.exitWithError("Failed to type check: ({})".format(statement.getText()))
                self.debugMsg("[Debug - checked]: {}: {}".format(statement.getText(), typechecked))
                all_typechecked = typechecked and all_typechecked
        except KeyError, keyerror:
            print "[TypeError] Undeclared variable: ", temp_st[:20], keyerror
            exit(0)

        if not all_typechecked:
            print "Process {} failed typechecker".format(pid)
        self.typecontext = {}
        return all_typechecked

    def visitParcomposition(self, ctx):
        print("Starting the type checker...")
        
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
