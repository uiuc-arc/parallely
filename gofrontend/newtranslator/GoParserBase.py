from antlr4 import *

def lineTerminatorAhead():
    return True

class GoParserBase(Parser):    
    def lineTerminatorAhead(self):
        return True
    
    def noTerminatorBetween(i):
        return True


