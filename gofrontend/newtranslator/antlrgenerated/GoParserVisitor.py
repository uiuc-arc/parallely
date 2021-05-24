# Generated from GoParser.g4 by ANTLR 4.7.1
from antlr4 import *

# This class defines a complete generic visitor for a parse tree produced by GoParser.

class GoParserVisitor(ParseTreeVisitor):

    # Visit a parse tree produced by GoParser#sourceFile.
    def visitSourceFile(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#packageClause.
    def visitPackageClause(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#importDecl.
    def visitImportDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#importSpec.
    def visitImportSpec(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#importPath.
    def visitImportPath(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#declaration.
    def visitDeclaration(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#constDecl.
    def visitConstDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#constSpec.
    def visitConstSpec(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#identifierList.
    def visitIdentifierList(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#expressionList.
    def visitExpressionList(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#typeDecl.
    def visitTypeDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#typeSpec.
    def visitTypeSpec(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#functionDecl.
    def visitFunctionDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#methodDecl.
    def visitMethodDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#receiver.
    def visitReceiver(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#varDecl.
    def visitVarDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#approxVarDecl.
    def visitApproxVarDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#varSpec.
    def visitVarSpec(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#block.
    def visitBlock(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#statementList.
    def visitStatementList(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#stmtdec.
    def visitStmtdec(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#stmtfunction.
    def visitStmtfunction(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#stmtlabeled.
    def visitStmtlabeled(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#stmtsimple.
    def visitStmtsimple(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#stmtgo.
    def visitStmtgo(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#stmtret.
    def visitStmtret(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#stmtblock.
    def visitStmtblock(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#stmtif.
    def visitStmtif(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#stmtfor.
    def visitStmtfor(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#stmtdefer.
    def visitStmtdefer(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#simpleStmt.
    def visitSimpleStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#expressionStmt.
    def visitExpressionStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#probc.
    def visitProbc(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#rec.
    def visitRec(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#condrec.
    def visitCondrec(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#origsend.
    def visitOrigsend(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#newsend.
    def visitNewsend(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#noisysend.
    def visitNoisysend(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#condsend.
    def visitCondsend(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#incDecStmt.
    def visitIncDecStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#assignment.
    def visitAssignment(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#assign_op.
    def visitAssign_op(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#shortVarDecl.
    def visitShortVarDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#emptyStmt.
    def visitEmptyStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#labeledStmt.
    def visitLabeledStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#returnStmt.
    def visitReturnStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#breakStmt.
    def visitBreakStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#continueStmt.
    def visitContinueStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#gotoStmt.
    def visitGotoStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#fallthroughStmt.
    def visitFallthroughStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#deferStmt.
    def visitDeferStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#simpleif.
    def visitSimpleif(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#switchStmt.
    def visitSwitchStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#exprSwitchStmt.
    def visitExprSwitchStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#exprCaseClause.
    def visitExprCaseClause(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#exprSwitchCase.
    def visitExprSwitchCase(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#typeSwitchStmt.
    def visitTypeSwitchStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#typeSwitchGuard.
    def visitTypeSwitchGuard(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#typeCaseClause.
    def visitTypeCaseClause(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#typeSwitchCase.
    def visitTypeSwitchCase(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#typeList.
    def visitTypeList(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#selectStmt.
    def visitSelectStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#commClause.
    def visitCommClause(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#commCase.
    def visitCommCase(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#recvStmt.
    def visitRecvStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#forStmt.
    def visitForStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#forClause.
    def visitForClause(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#rangeClause.
    def visitRangeClause(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#goStmt.
    def visitGoStmt(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#type_.
    def visitType_(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#typeName.
    def visitTypeName(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#typeLit.
    def visitTypeLit(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#arrayType.
    def visitArrayType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#arrayLength.
    def visitArrayLength(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#elementType.
    def visitElementType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#pointerType.
    def visitPointerType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#interfaceType.
    def visitInterfaceType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#sliceType.
    def visitSliceType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#mapType.
    def visitMapType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#channelType.
    def visitChannelType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#methodSpec.
    def visitMethodSpec(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#functionType.
    def visitFunctionType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#signature.
    def visitSignature(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#result.
    def visitResult(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#parameters.
    def visitParameters(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#parameterDecl.
    def visitParameterDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#expression.
    def visitExpression(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#primaryExpr.
    def visitPrimaryExpr(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#unaryExpr.
    def visitUnaryExpr(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#conversion.
    def visitConversion(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#operand.
    def visitOperand(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#literal.
    def visitLiteral(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#basicLit.
    def visitBasicLit(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#integer.
    def visitInteger(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#operandName.
    def visitOperandName(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#qualifiedIdent.
    def visitQualifiedIdent(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#compositeLit.
    def visitCompositeLit(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#literalType.
    def visitLiteralType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#literalValue.
    def visitLiteralValue(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#elementList.
    def visitElementList(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#keyedElement.
    def visitKeyedElement(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#key.
    def visitKey(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#element.
    def visitElement(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#structType.
    def visitStructType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#fieldDecl.
    def visitFieldDecl(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#string_.
    def visitString_(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#anonymousField.
    def visitAnonymousField(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#functionLit.
    def visitFunctionLit(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#index.
    def visitIndex(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#goslice.
    def visitGoslice(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#typeAssertion.
    def visitTypeAssertion(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#arguments.
    def visitArguments(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#methodExpr.
    def visitMethodExpr(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#receiverType.
    def visitReceiverType(self, ctx):
        return self.visitChildren(ctx)


    # Visit a parse tree produced by GoParser#eos.
    def visitEos(self, ctx):
        return self.visitChildren(ctx)


