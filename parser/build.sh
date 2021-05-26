ANTLRJAR=/home/vimuth/Downloads/antlr4-4.7.1-complete.jar

java -Xmx500M -cp "$ANTLRJAR:$CLASSPATH" org.antlr.v4.Tool -Dlanguage=Python2 -visitor Parallely.g4
