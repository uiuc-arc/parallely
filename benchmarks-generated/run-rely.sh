for benchmark in bfs-recover pagerank-recover sobel-recover sssp-recover bscholes-recover scale-recover sor-recover motion-recover; do
    echo $benchmark
    python ../parser/rely_recover.py -f ./src/$benchmark/$benchmark.par -s ./src/$benchmark/spec.rely -ifs | grep "Analysis time"
done
