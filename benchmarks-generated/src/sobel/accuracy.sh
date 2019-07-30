# First run original
rm -f sobel
rm -f temp.go
python ../../../parser/crosscompiler.py -f ./sobel.par -o temp.go -t __basic_go.txt &> _debug_.txt
go build

for i in {0..9}; do 
    ./sobel $i
    mv __output__.txt __output__precise_$i.txt
done


# Precision
rm -f sobel
rm -f temp.go
python ../../../parser/crosscompiler.py -f ./sobel-prec.par -o temp.go -t __basic_go.txt &> _debug_.txt
go build

for i in {0..9}; do 
    ./sobel $i
    mv __output__.txt __output__approx_$i.txt
done
