cd ../../../../gofrontend/

Run translator
echo "**************************************"
echo "Running the translator (Go -> Parallely)."
echo "**************************************"
python -m newtranslator.translator.translator -f ../benchmarks/golang/src/pagerank/pagerank.go -o ../benchmarks/golang/src/pagerank/out.par

# cd -
# # Run sequentializer
# echo "**************************************"
# echo "Running the sequentializer"
# echo "**************************************"
# python ../../../../parser/compiler.py -f out.par -o out.seq


# Run translator
echo "**************************************"
echo "Generating executable code (Go -> Go)."
echo "**************************************"
python -m newtranslator.translator.typedGoGenerator -f ../benchmarks/golang/src/pagerank/pagerank.go -o ../benchmarks/golang/src/pagerank/pagerank.exec.go
