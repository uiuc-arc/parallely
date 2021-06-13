benchmark=scale

cd ../../../../gofrontend/

# Run translator
echo "**************************************"
echo "Running the translator Go -> Parallely"
echo "**************************************"
python -m translator.translator -f ../benchmarks/golang/src/$benchmark/$benchmark.go -o ../benchmarks/golang/src/$benchmark/$benchmark.par

cd - > /dev/null
# Run sequentializer
echo "**************************************"
echo "Running the sequentializer"
echo "**************************************"
python ../../../../parser/compiler.py -f $benchmark.par -o $benchmark.seq

# Run translator
cd ../../../../gofrontend/
echo "**************************************"
echo "Generating executable code by renaming function calls with types (Go -> Go)."
python -m translator.typedGoGenerator -f ../benchmarks/golang/src/$benchmark/$benchmark.go -o ../benchmarks/golang/src/$benchmark/$benchmark.exec.go
echo "Use 'go run $benchmark.exec.go scalefuncs.go ../../inputs/baboon.ppm out.ppm' to run the generated program" 
echo "**************************************"
