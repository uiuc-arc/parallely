# Generated from GoLang.g4 by ANTLR 4.7.1
from antlr4 import *

# This class defines a complete generic visitor for a parse tree produced by GoLangParser.

class GoLangVisitor(ParseTreeVisitor):

    # Visit a parse tree produced by GoLangParser#sourceFile.
    def visitSourceFile(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#packageClause.
    def visitPackageClause(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#importDecl.
    def visitImportDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#importSpec.
    def visitImportSpec(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#importPath.
    def visitImportPath(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#declaration.
    def visitDeclaration(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#constDecl.
    def visitConstDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#constSpec.
    def visitConstSpec(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#identifierList.
    def visitIdentifierList(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#expressionList.
    def visitExpressionList(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#typeDecl.
    def visitTypeDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#typeSpec.
    def visitTypeSpec(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#functionDecl.
    def visitFunctionDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#methodDecl.
    def visitMethodDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#receiver.
    def visitReceiver(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#varDecl.
    def visitVarDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#varSpec.
    def visitVarSpec(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#block.
    def visitBlock(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#statementList.
    def visitStatementList(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#smtdec.
    def visitSmtdec(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#smtfunction.
    def visitSmtfunction(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#smtfunctionassign.
    def visitSmtfunctionassign(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#smtlabeled.
    def visitSmtlabeled(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#smtsimple.
    def visitSmtsimple(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#smtgo.
    def visitSmtgo(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#smtret.
    def visitSmtret(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#smtblock.
    def visitSmtblock(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#smtif.
    def visitSmtif(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#smtfor.
    def visitSmtfor(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#smtdefer.
    def visitSmtdefer(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#smtsend.
    def visitSmtsend(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#smtinc.
    def visitSmtinc(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#smtass.
    def visitSmtass(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#smtshortdec.
    def visitSmtshortdec(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#smtempty.
    def visitSmtempty(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#expressionStmt.
    def visitExpressionStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#sendStmt.
    def visitSendStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#incDecStmt.
    def visitIncDecStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#assignment.
    def visitAssignment(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#assign_op.
    def visitAssign_op(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#shortVarDecl.
    def visitShortVarDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#emptyStmt.
    def visitEmptyStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#labeledStmt.
    def visitLabeledStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#returnStmt.
    def visitReturnStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#breakStmt.
    def visitBreakStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#continueStmt.
    def visitContinueStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#gotoStmt.
    def visitGotoStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#fallthroughStmt.
    def visitFallthroughStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#deferStmt.
    def visitDeferStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#ifStmt.
    def visitIfStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#switchStmt.
    def visitSwitchStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#exprSwitchStmt.
    def visitExprSwitchStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#exprCaseClause.
    def visitExprCaseClause(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#exprSwitchCase.
    def visitExprSwitchCase(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#typeSwitchStmt.
    def visitTypeSwitchStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#typeSwitchGuard.
    def visitTypeSwitchGuard(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#typeCaseClause.
    def visitTypeCaseClause(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#typeSwitchCase.
    def visitTypeSwitchCase(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#typeList.
    def visitTypeList(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#selectStmt.
    def visitSelectStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#commClause.
    def visitCommClause(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#commCase.
    def visitCommCase(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#recvStmt.
    def visitRecvStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#forStmt.
    def visitForStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#forClause.
    def visitForClause(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#rangeClause.
    def visitRangeClause(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#goStmt.
    def visitGoStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#type_.
    def visitType_(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#typeName.
    def visitTypeName(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#typeLit.
    def visitTypeLit(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#arrayType.
    def visitArrayType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#arrayLength.
    def visitArrayLength(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#elementType.
    def visitElementType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#pointerType.
    def visitPointerType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#interfaceType.
    def visitInterfaceType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#sliceType.
    def visitSliceType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#mapType.
    def visitMapType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#channelType.
    def visitChannelType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#methodSpec.
    def visitMethodSpec(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#functionType.
    def visitFunctionType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#signature.
    def visitSignature(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#result.
    def visitResult(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#parameters.
    def visitParameters(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#parameterDecl.
    def visitParameterDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#expression.
    def visitExpression(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#primaryExpr.
    def visitPrimaryExpr(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#unaryExpr.
    def visitUnaryExpr(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#conversion.
    def visitConversion(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#operand.
    def visitOperand(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#literal.
    def visitLiteral(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#basicLit.
    def visitBasicLit(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#integer.
    def visitInteger(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#operandName.
    def visitOperandName(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#qualifiedIdent.
    def visitQualifiedIdent(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#compositeLit.
    def visitCompositeLit(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#literalType.
    def visitLiteralType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#literalValue.
    def visitLiteralValue(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#elementList.
    def visitElementList(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#keyedElement.
    def visitKeyedElement(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#key.
    def visitKey(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#element.
    def visitElement(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#structType.
    def visitStructType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#fieldDecl.
    def visitFieldDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#string_.
    def visitString_(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#anonymousField.
    def visitAnonymousField(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#functionLit.
    def visitFunctionLit(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#index.
    def visitIndex(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#goslice.
    def visitGoslice(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#typeAssertion.
    def visitTypeAssertion(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#arguments.
    def visitArguments(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#methodExpr.
    def visitMethodExpr(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#receiverType.
    def visitReceiverType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoLangParser#eos.
    def visitEos(self, ctx):
        return self.visitChildren(ctx)


