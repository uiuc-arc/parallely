Parallely
======

Parallely, is a programming language and a system for verification of
approximations in parallel message-passing programs. Parallely’s
language can express various software and hardware level
approximations that reduce the computation and communication
overheads at the cost of result accuracy.

Parallely takes as input programs written in Go programming language
(with some extensions) and generates a equivalent sequential program
that maintain the program semantics. This sequential program can be
used to verify important safety and accuracy properties of the
original program.

---

Directory Structure
-------------------
* `gofrontend` contains the parser for the programs in Go language and
the translator that converts the program to Parallely intermediate language.
* `benchmarks` contains example programs.
* `parser` contains the code for analyzing parallely programs. The
compiler unrolls bounded loops, performs type checking and generates a
equivalent sequential program.

Setup
-------------------

* Setup ANTLR v4 following the instructions in
  https://github.com/antlr/antlr4/blob/master/doc/getting-started.md
* Update the location of the downloaded jar file in the `build.sh`
  file
* Run `build.sh`: This generates the Lexer and Parser code for both
  the GoLang frontend and the Parallely intermediate language

Example
-------------------

`benchmarks/golang/pagerank/pagerank.go` contains example
code for the Pagerank computation using a set of parallel threads.

Run the script in 


Need to run from home folder due to antlr generated code using
relative imports. Once the language is fully finalized we should be
able to change this.

python -m newtranslator.translator.translator -f ./src/kmeans/kmeans.go
