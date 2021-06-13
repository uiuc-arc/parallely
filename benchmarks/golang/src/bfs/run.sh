cd ../../../../gofrontend/

# Run translator
echo "**************************************"
echo "Running the translator Go -> Parallely"
echo "**************************************"
python -m translator.translator -f ../benchmarks/golang/src/bfs/bfs.go -o ../benchmarks/golang/src/bfs/bfs.par

cd -
# Run sequentializer
echo "**************************************"
echo "Running the sequentializer"
echo "**************************************"
python ../../../../parser/compiler.py -f bfs.par -o bfs.seq

# Run translator
cd ../../../../gofrontend/
echo "**************************************"
echo "Generating executable code by renaming function calls with types (Go -> Go)."
python -m translator.typedGoGenerator -f ../benchmarks/golang/src/bfs/bfs.go -o ../benchmarks/golang/src/bfs/bfs.exec.go
echo "Use 'go run bfs.exec.go' to run the generated program" 
echo "**************************************"
