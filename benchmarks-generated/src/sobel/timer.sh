# First run original
rm -f sobel
rm -f temp.go
python ../../../parser/crosscompiler.py -f ./sobel.par -o temp.go -t __basic_go.txt &> _debug_.txt
go build
python getTime.py

# Precision
rm -f sobel
rm -f temp.go
python ../../../parser/crosscompiler.py -f ./sobel-prec.par -o temp.go -t __basic_go.txt &> _debug_.txt
go build
python getTime.py 
