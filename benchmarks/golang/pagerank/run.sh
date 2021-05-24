cd ../../../gofrontend/

# Run translator
echo "**************************************"
echo "Running the translator Go -> Parallely"
echo "**************************************"
python -m newtranslator.translator.translator -f ../benchmarks/golang/pagerank/pagerank.go -o ../benchmarks/golang/pagerank/out.par

cd -
# Run sequentializer
echo "**************************************"
echo "Running the sequentializer"
echo "**************************************"
python ../../../parser/compiler.py -f out.par -o out.seq
