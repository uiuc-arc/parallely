ANTLRJAR=<Update with location of antlr jar>

java -Xmx500M -cp "$ANTLRJAR:$CLASSPATH" org.antlr.v4.Tool -Dlanguage=Python2 -visitor Parallely.g4
