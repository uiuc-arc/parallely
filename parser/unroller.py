from ParallelyListener import ParallelyListener
import TokenStreamRewriter
from ParallelyLexer import ParallelyLexer
from ParallelyParser import ParallelyParser

class unrollRepeat(ParallelyListener):
    def __init__(self, stream, replacement, replacement_map):
        self.rewriter = TokenStreamRewriter.TokenStreamRewriter(stream)
        self.replacedone = False
        self.replacement = replacement
        self.dummymap = replacement_map

    # def exitRepeat(self, ctx):
    def replace_dummies(self, program_str):
        while(self.replacement >= 0):
            # print program_str, self.replacement
            dummy_syntax = "<dummy {}>".format(self.replacement)
            program_str = program_str.replace(dummy_syntax, self.dummymap[self.replacement])
            self.replacement -= 1            
        return program_str

    def exitRepeat(self, ctx):
        if self.replacedone:
            return

        dummy_syntax = "<dummy {}>".format(self.replacement)

        rep_variable = int(ctx.INT().getText())
        # TODO: Is there a way to avoid string manipulation?
        list_statements = ctx.statement()
        cs = list_statements[0].start.getInputStream()
        statements = cs.getText(list_statements[0].start.start,
                                list_statements[-1].stop.stop)
        self.dummymap[self.replacement] = statements

        new_str = ''
        # for var in range(rep_variable):
        new_str = ("  " + dummy_syntax + ";\n") * rep_variable
        self.rewriter.insertAfter(ctx.stop.tokenIndex + 1, new_str)
        self.rewriter.delete(self.rewriter.DEFAULT_PROGRAM_NAME,
                             ctx.start.tokenIndex,
                             ctx.stop.tokenIndex + 1)
        self.replacedone = True
        self.replacement += 1
        return
