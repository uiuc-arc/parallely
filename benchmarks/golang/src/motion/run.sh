benchmark=motion

cd ../../../../gofrontend/

# Run translator
echo "**************************************"
echo "Running the translator Go -> Parallely"
echo "**************************************"
python -m newtranslator.translator.translator -f ../benchmarks/golang/src/$benchmark/$benchmark.go -o ../benchmarks/golang/src/$benchmark/$benchmark.par

cd -
# # Run sequentializer
# echo "**************************************"
# echo "Running the sequentializer"
# echo "**************************************"
# python ../../../../parser/compiler.py -f $benchmark.par -o $benchmark.seq

# Run translator
cd ../../../../gofrontend/
echo "**************************************"
echo "Generating executable code by renaming function calls with types (Go -> Go)."
python -m newtranslator.translator.typedGoGenerator -f ../benchmarks/golang/src/$benchmark/$benchmark.go -o ../benchmarks/golang/src/$benchmark/$benchmark.exec.go
echo "Use 'go run $benchmark.exec.go' to run the generated program" 
echo "**************************************"
