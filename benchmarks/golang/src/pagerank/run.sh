cd ../../../../gofrontend/

# Run translator
echo "**************************************"
echo "Running the translator (Go -> Parallely)."
echo "**************************************"
python -m translator.translator -f ../benchmarks/golang/src/pagerank/pagerank.go -o ../benchmarks/golang/src/pagerank/out.par

cd - > /dev/null
# Run sequentializer
echo "**************************************"
echo "Running the sequentializer"
echo "**************************************"
python ../../../../parser/compiler.py -f out.par -o out.seq


# Run translator
cd ../../../../gofrontend/
echo "**************************************"
echo "Generating executable code by renaming function calls with types (Go -> Go)."
python -m translator.typedGoGenerator -f ../benchmarks/golang/src/pagerank/pagerank.go -o ../benchmarks/golang/src/pagerank/pagerank.exec.go
echo "Use 'go run pagerank.exec.go' to run the generated program" 
echo "**************************************"
