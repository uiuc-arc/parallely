for benchmark in bfs-recover pagerank-recover sobel-recover sssp-recover bscholes-recover scale-recover sor-recover motion-recover; do
    echo python ../parser/rely_recover.py -f ./src/$benchmark/$benchmark-kernel.par -s ./src/$benchmark/kernel.spec -func ../parser/unreliable-chcker.json
    python ../parser/rely_recover.py -f ./src/$benchmark/$benchmark-kernel.par -s ./src/$benchmark/kernel.spec -func ../parser/unreliable-chcker.json
    echo "#############################################"
done
