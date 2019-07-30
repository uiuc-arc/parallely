# First run original
rm -f sor
rm -f temp.go
python ../../../parser/crosscompiler.py -f ./motion.par -o temp.go -t __basic_go.txt &> _debug_.txt
go build
time1="$(python getTime.py | tail -n1)"
# time1="$(echo $outstring | tail -n1)"
echo $time1

# Precision
rm -f sor
rm -f temp.go
python ../../../parser/crosscompiler.py -f ./motion-approxred.par -o temp.go -t __basic_go.txt &> _debug_.txt
go build
time2="$(python getTime.py | tail -n1)"
echo $time2

echo "scale=3; $time1/$time2" | bc 
