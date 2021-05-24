from antlr4 import Parser
from antlr4 import Lexer
from newtranslator.antlrgenerated.GoLexer import GoLexer


class GoParserBase(Parser):
    def lineTerminatorAhead(self):
        possibleIndexEosToken = self.getCurrentToken().tokenIndex - 1

        if possibleIndexEosToken == -1:
            return True

        ahead = self._input.get(possibleIndexEosToken)
        if ahead.channel != Lexer.HIDDEN:
            return False

        if ahead.type == GoLexer.TERMINATOR:
            return True

        if ahead.type == GoLexer.WS:
            possibleIndexEosToken = self.getCurrentToken().tokenIndex - 2
            ahead = self._input.get(possibleIndexEosToken)

        text = ahead.text
        token_type = ahead.type

        return ((token_type == GoLexer.COMMENT) and
                ((text.find("\r") > -1) or (text.find("\n") > -1))) or (token_type == GoLexer.TERMINATOR)

    def noTerminatorBetween(self, i):
        stream = self._input
        tokens = stream.getHiddenTokensToLeft(stream.LT(i).tokenIndex)

        if not tokens:
            return True

        for token in tokens:
            if token.text.find("\n") > -1:
                return False
        return True

    def checkPreviousTokenText(self, i):
        stream = self._input
        prev_token = stream.LT(1).text
        if not prev_token:
            return False
        return prev_token == i

    def noTerminatorAfterParams(self, tokenOffset):
        leftParams = 1
        rightParams = 0
        stream = self._input

        if (stream.LT(tokenOffset).type == GoLexer.L_PAREN):
            while leftParams != rightParams:
                tokenOffset += 1
                val_type = stream.LT(tokenOffset).type

                if val_type == GoLexer.L_PAREN:
                    leftParams += 1
                if val_type == GoLexer.R_PAREN:
                    rightParams += 1
            tokenOffset += 1
            return self.noTerminatorBetween(tokenOffset)
        return True
