# First run original
rm -f sor
rm -f temp.go
python ../../../parser/crosscompiler.py -f ./sor.par -o temp.go -t __basic_go.txt &> _debug_.txt
go build
python getTime.py

# Precision
rm -f sor
rm -f temp.go
python ../../../parser/crosscompiler.py -f ./sor-prec.par -o temp.go -t __basic_go.txt &> _debug_.txt
go build
python getTime.py 
