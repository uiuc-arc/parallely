from ParallelyParser import ParallelyParser
from ParallelyLexer import ParallelyLexer
from ParallelyListener import ParallelyListener
import TokenStreamRewriter


class VariableRenamer(ParallelyListener):
    def __init__(self, stream):
        self.current_process = None
        self.rewriter = TokenStreamRewriter.TokenStreamRewriter(stream)
        self.done = set()
        self.skiplist = []

    def enterSingle(self, ctx):
        pid = ctx.processid()
        if isinstance(pid, ParallelyParser.NamedpContext):
            self.current_process = ctx.processid()
        elif isinstance(pid, ParallelyParser.VariablepContext):
            self.current_process = ctx.processid()
        else:
            self.current_process = ctx.processid().VAR()
            self.skiplist.append(ctx.processid().VAR().getText())

    def exitSingle(self, ctx):
        if isinstance(ctx.processid(), ParallelyParser.GroupedpContext):
            top = self.skiplist.pop(0)
            if top != ctx.processid().VAR().getText():
                print "[Error] Does not match: ", top, ctx.processid()

    def enterForloop(self, ctx):
        iter_var = ctx.VAR().getText()
        self.skiplist.append(iter_var)

    def exitForloop(self, ctx):
        iter_var = ctx.VAR().getText()
        top = self.skiplist.pop(0)
        if top != iter_var:
            print "[Error] Does not match"

    def enterLocalvariable(self, ctx):
        if not ctx.getText() in self.skiplist:
            new_name = "$" + self.current_process.getText()
            self.rewriter.insertAfterToken(ctx.stop, new_name)
