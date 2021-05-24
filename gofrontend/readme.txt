Need to run from home folder due to antlr generated code using
relative imports. Once the language is fully finalized we should be
able to change this.

python -m newtranslator.translator.translator -f ./src/kmeans/kmeans.go 
