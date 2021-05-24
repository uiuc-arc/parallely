cd ../../../gofrontend/

# Run translator
echo "**************************************"
echo "Running the translator Go -> Parallely"
echo "**************************************"
python -m newtranslator.translator.translator -f ../benchmarks/golang/bfs/bfs.go -o ../benchmarks/golang/bfs/bfs.par

# cd -
# # Run sequentializer
# echo "**************************************"
# echo "Running the sequentializer"
# echo "**************************************"
# python ../../../parser/compiler.py -f bfs.par -o bfs.seq

