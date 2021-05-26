ANTLRJAR=/home/vimuth/Downloads/antlr4-4.7.1-complete.jar

cd gofrontend/newtranslator/antlrgenerated/
java -Xmx500M -cp "/usr/local/lib/antlr-4.7.1-complete.jar:$CLASSPATH" org.antlr.v4.Tool -Dlanguage=Python2 -visitor GoLexer.g4 
java -Xmx500M -cp "/usr/local/lib/antlr-4.7.1-complete.jar:$CLASSPATH" org.antlr.v4.Tool -Dlanguage=Python2 -visitor GoParser.g4

cd -
cd parser/
java -Xmx500M -cp "$ANTLRJAR:$CLASSPATH" org.antlr.v4.Tool -Dlanguage=Python2 -visitor Parallely.g4
