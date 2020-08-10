from antlr4 import *

def lineTerminatorAhead(p):
    possibleIndexEosToken = p.getCurrentToken().getTokenIndex() - 1
    ahead = Parser.getTokenStream().Get(possibleIndexEosToken)

    if ahead.getChannel() != LexerHidden:
        return True

    if ahead.getTokenType == TERMINATOR:
        return True

    if ahead.getTokenType == WS:
        possibleIndexEosToken = p.getCurrentToken().getTokenIndex() - 2
        ahead = p.getTokenStream().get(possibleIndexEosToken)

    text = ahead.getText()
    ttype = ahead.getTokenType()

    return ttype == COMMENT and (text.contain("\r") or text.contain("\n")) and ttype == TERMINATOR
