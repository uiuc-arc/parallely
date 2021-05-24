# flake8: noqa E501

from antlr4 import CommonTokenStream
from antlr4 import InputStream
import TokenStreamRewriter
from antlr4 import ParseTreeWalker
from ParallelyLexer import ParallelyLexer
from ParallelyParser import ParallelyParser
from ParallelyListener import ParallelyListener
import time
from argparse import ArgumentParser

from sequentializer import parallelySequentializer
from renamer import VariableRenamer
from typechecker import parallelyTypeChecker
from unroller import unrollRepeat


# class unrollRepeat(ParallelyListener):
#     def __init__(self, stream):
#         self.rewriter = TokenStreamRewriter.TokenStreamRewriter(stream)
#         self.replacedone = False

#     # def enterRepeat(self, ctx):
#     #     cs = ctx.statement().start.getInputStream()
#     #     statements = cs.getText(ctx.statement().start.start,
#     #                             ctx.statement().stop.stop)
#     #     rep_variable = int(ctx.INT().getText())
#     #     edited = ''
#     #     # removing the code for process groups
#     #     self.rewriter.delete(self.rewriter.DEFAULT_PROGRAM_NAME,
#     #                          ctx.start.tokenIndex, ctx.stop.tokenIndex)
#     #     for var in range(rep_variable):
#     #         edited += statements + ";\n"
#     #     self.rewriter.insertAfter(ctx.stop.tokenIndex, edited)

#     def enterRepeat(self, ctx):
#         # Do only one replacement at a time
#         if self.replacedone:
#             return

#         rep_variable = int(ctx.INT().getText())
#         # TODO: Is there a way to avoid string manipulation?
#         list_statements = ctx.statement()
#         cs = list_statements[0].start.getInputStream()
#         statements = cs.getText(list_statements[0].start.start,
#                                 list_statements[-1].stop.stop)

#         # print "------------------------------"
#         # print "Unrolling {} times: ".format(rep_variable), statements
#         # print "------------------------------"
#         new_str = ''
#         for var in range(rep_variable):
#             new_str += "  " + statements + ";\n"
#         self.rewriter.insertAfter(ctx.stop.tokenIndex + 1, new_str)
#         self.rewriter.delete(self.rewriter.DEFAULT_PROGRAM_NAME,
#                              ctx.start.tokenIndex,
#                              ctx.stop.tokenIndex + 1)
#         self.replacedone = True


def main(program_str, outfile, filename, args):
    input_stream = InputStream(program_str)

    fullstart = time.time()
    
    if not args.skip:
        print("Unrolling Repeat statements")
        replacement = 0
        replacement_map = {}
        i = 0
        while(True):
            lexer = ParallelyLexer(input_stream)
            stream = CommonTokenStream(lexer)
            parser = ParallelyParser(stream)

            try:
                tree = parser.parallelprogram()
            except Exception as e:
                print("Parsing Error!!!")
                print(e)
                exit(-1)

            unroller = unrollRepeat(stream, replacement, replacement_map)
            walker = ParseTreeWalker()
            walker.walk(unroller, tree)
            input_stream = InputStream(unroller.rewriter.getDefaultText())
            replacement = unroller.replacement
            # print unroller.replacement, unroller.dummymap
            if not unroller.replacedone:
                input_stream = InputStream(unroller.rewriter.getDefaultText())
                break
                # if debug:
            i = i + 1

            debug_file = open("_DEBUG_UNROLLED_{}.txt".format(i), 'w')
            debug_file.write(input_stream.strdata)
            debug_file.close()
        unroller = unrollRepeat(stream, replacement - 1, replacement_map)
        new_program = unroller.replace_dummies(input_stream.strdata)
        debug_file = open("_DEBUG_UNROLLED_FINAL.txt", 'w')
        debug_file.write(new_program)
        debug_file.close()
    else:
        new_program = input_stream.strdata
    input_stream = InputStream(new_program)

    # if not args.skip:
    #     print "Recursively Unrolling Repeat statements"
    #     while(True):
    #         lexer = ParallelyLexer(input_stream)
    #         stream = CommonTokenStream(lexer)
    #         parser = ParallelyParser(stream)
    #         tree = parser.parallelprogram()

    #         unroller = unrollRepeat(stream)
    #         walker = ParseTreeWalker()
    #         walker.walk(unroller, tree)
    #         input_stream = InputStream(unroller.rewriter.getDefaultText())
    #         # print unroller.replacedone
    #         if not unroller.replacedone:
    #             # print unroller.replacedone
    #             input_stream = InputStream(unroller.rewriter.getDefaultText())
    #             break

    #     debug_file = open("_DEBUG_UNROLLED_.txt", 'w')
    #     debug_file.write(input_stream.strdata)
    #     debug_file.close()

    if not args.skipunroll:
        lexer = ParallelyLexer(input_stream)
        stream = CommonTokenStream(lexer)
        parser = ParallelyParser(stream)
        tree = parser.parallelprogram()

        print "Renaming all variables"
        renamer = VariableRenamer(stream)
        walker = ParseTreeWalker()
        walker.walk(renamer, tree)

        start = time.time()

        # Run type checker on the renamed version
        input_stream = InputStream(renamer.rewriter.getDefaultText())
        # if debug:
        debug_file = open("_DEBUG_RENAMED_.txt", 'w')
        debug_file.write(input_stream.strdata)
        debug_file.close()

    lexer = ParallelyLexer(input_stream)
    stream = CommonTokenStream(lexer)
    parser = ParallelyParser(stream)

    tree = parser.parallelprogram()

    start = time.time()
    if not args.skiptypes:
        print("----------------------------------------")
        print "Running type checker"
        typechecker = parallelyTypeChecker(args.debug)
        type_checked = typechecker.visit(tree)
        end = time.time()
        if type_checked:
            print "{} passed type checker ({}s).".format(filename, end - start)
            print("----------------------------------------")
        else:
            print "{} failed type checker ({}s).".format(filename, end - start)
            exit(-1)
    end = time.time()

    # Sequentialization
    print "Running sequentialization"
    start2 = time.time()
    sequentializer = parallelySequentializer(args.debug, args.annotate)
    sequentializer.rewriteProgram(tree, outfile)
    end2 = time.time()
    print "Time for sequentialization :", end2 - start2

    print "Total time : ", (end2 - start2) + (end - start)


if __name__ == '__main__':
    parser = ArgumentParser()
    parser.add_argument("-f", dest="programfile",
                        help="File containing the code")
    parser.add_argument("-o", dest="outfile",
                        help="File to output the sequential code")
    parser.add_argument("-s", "--skip", action="store_true",
                        help="Skip renaming")
    parser.add_argument("-t", "--skiptypes", action="store_true",
                        help="Skip type checking")
    parser.add_argument("-u", "--skipunroll", action="store_true",
                        help="Skip unrolling repeats")
    parser.add_argument("-d", "--debug", action="store_true",
                        help="Print debug info")
    parser.add_argument("-g", "--annotate", action="store_true",
                        help="annotate with debug info")
    args = parser.parse_args()

    programfile = open(args.programfile, 'r')
    outfile = open(args.outfile, 'w')
    program_str = programfile.read()
    main(program_str, outfile, programfile.name, args)
