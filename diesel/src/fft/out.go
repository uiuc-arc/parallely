package main

import (
	"diesel"
  "bufio"
	"math"
	"fmt"
	"math/bits"
	"unsafe"
	"math/rand"
	"time"
  "strconv"
  "os"  
)

func getDistance(max_index int) float32 {
	RADAR_N := 10
	RADAR_fs := 204800.0
	RADAR_alpha := 30000000000.0
	RADAR_c := 300000000.0
	return (float32(max_index) * float32(RADAR_fs) / float32(RADAR_N)) * float32(0.5*RADAR_c) / float32(RADAR_alpha)
}

func getFloat(i int) float32 {
 return float32(i)
}

func getFloat32(i float64) float32 {
 return float32(i)
}

func getSin32(i float32) float32 {
 return float32(math.Sin(float64(i)))
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func _rev(v uint) uint {
	r := v
	s := unsafe.Sizeof(v)*8 - 1

	for v >>= 1; v != 0; v >>= 1 {
		r <<= 1
		r |= v & 1
		s--
	}
	r <<= s

	return r
}

func bitReverse(w [2048]float64, N int, bitsin int) [2048]float64 {
	var i, r, s, shift uint

	var t_real, t_imag float64

	s = uint(unsafe.Sizeof(i)*8 - 1)
	shift = s - uint(bitsin) + 1

	for i = 0; i < uint(N); i++ {
		r = _rev(i)
		r2 := bits.Reverse(i)
		r >>= shift
		r2 >>= shift

		if i < r {
			t_real = w[2*i]
			t_imag = w[2*i+1]
			w[2*i] = w[2*r]
			w[2*i+1] = w[2*r+1]
			w[2*r] = t_real
			w[2*r+1] = t_imag
		}
	}

	return w
}

var Num_threads int
var Inputdata [2048]float64
var Outdata [2048]float32
var Pi = float32(3.141592653589)
var Distance float32



func func_0() {
  defer diesel.Wg.Done();
  var DynMap [4119]diesel.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var data64 [2048]float64;
diesel.InitDynArray(0, 2048, DynMap[:]);
var data [2048]float32;
diesel.InitDynArray(2048, 2048, DynMap[:]);
var N int;
var logN int;
var bit int;
var sign float32;
var a int;
var b int;
var index1 int;
var index2 int;
var index3 int;
var index4 int;
var i int;
var j int;
var transform_length int;
var s float32;
var t float32;
var s2 float32;
var theta float32;
DynMap[4096] = diesel.ProbInterval{1, 0};
var di int;
DynMap[4097] = diesel.ProbInterval{1, 0};
var z_real float32;
DynMap[4098] = diesel.ProbInterval{1, 0};
var z_imag float32;
DynMap[4099] = diesel.ProbInterval{1, 0};
var t_real float32;
DynMap[4100] = diesel.ProbInterval{1, 0};
var t_imag float32;
DynMap[4101] = diesel.ProbInterval{1, 0};
var w_real float32;
DynMap[4102] = diesel.ProbInterval{1, 0};
var w_imag float32;
DynMap[4103] = diesel.ProbInterval{1, 0};
var temp0 float64;
DynMap[4104] = diesel.ProbInterval{1, 0};
var temp1 float32;
DynMap[4105] = diesel.ProbInterval{1, 0};
var temp2 float32;
DynMap[4106] = diesel.ProbInterval{1, 0};
var temp3 float32;
DynMap[4107] = diesel.ProbInterval{1, 0};
var temp4 float32;
DynMap[4108] = diesel.ProbInterval{1, 0};
var temp5 float32;
DynMap[4109] = diesel.ProbInterval{1, 0};
var temp6 float32;
DynMap[4110] = diesel.ProbInterval{1, 0};
var temp7 float32;
DynMap[4111] = diesel.ProbInterval{1, 0};
var temp8 float32;
DynMap[4112] = diesel.ProbInterval{1, 0};
var temp9 float32;
DynMap[4113] = diesel.ProbInterval{1, 0};
var temp10 float32;
DynMap[4114] = diesel.ProbInterval{1, 0};
var temp11 float32;
DynMap[4115] = diesel.ProbInterval{1, 0};
var maxpsd float32;
DynMap[4116] = diesel.ProbInterval{1, 0};
var maxindex int;
DynMap[4117] = diesel.ProbInterval{1, 0};
var distance float32;
DynMap[4118] = diesel.ProbInterval{1, 0};
N = 1024;
logN = 10;
sign = -1.0;
transform_length = 1;
data64=bitReverse(Inputdata,N,logN);
i = 0;
for __temp_0 := 0; __temp_0 < 2048; __temp_0++ {
 _temp_index_1 := i;
temp0=data64[_temp_index_1];
DynMap[4104] = DynMap[0 + _temp_index_1];
temp2 = float32(temp0);
DynMap[4106].Reliability = DynMap[4104].Reliability;
 DynMap[4106].Delta = diesel.GetCastingError64to32(temp0, temp2);
_temp_index_2 := i;
data[_temp_index_2]=temp2;
DynMap[2048 + _temp_index_2] = DynMap[4106];
i = i+1;
 }
bit = 0;
for __temp_1 := 0; __temp_1 < 10; __temp_1++ {
 DynMap[4102] = diesel.ProbInterval{1, 0};
w_real = 1.0;
DynMap[4103] = diesel.ProbInterval{1, 0};
w_imag = 0.0;
temp5=getFloat(transform_length);
temp2=getFloat32(1.0);
DynMap[4107].Reliability = DynMap[4106].Reliability;
DynMap[4107].Delta = math.Abs(float64(sign)) * DynMap[4106].Delta;
temp3 = temp2*sign;
DynMap[4108].Reliability = DynMap[4107].Reliability;
DynMap[4108].Delta = math.Abs(float64(Pi)) * DynMap[4107].Delta;
temp4 = temp3*Pi;
DynMap[4096].Reliability = DynMap[4109].Reliability + DynMap[4108].Reliability - 1.0;
DynMap[4096].Delta = math.Abs(float64(temp4)) * DynMap[4108].Delta + math.Abs(float64(temp5)) * DynMap[4109].Delta / (math.Abs(float64(temp5)) * (math.Abs(float64(temp4))-DynMap[4109].Delta));
theta = temp4/temp5;
s=getSin32(theta);
t=getSin32(0.5*theta);
temp2=getFloat32(2.0);
s2 = temp2*t*t;
a = 0;
for a<transform_length {
 b = 0;
for b<N {
 i = b+a;
j = b+a+transform_length;
index1 = 2*j;
index2 = index1+1;
_temp_index_3 := index1;
z_real=data[_temp_index_3];
DynMap[4098] = DynMap[2048 + _temp_index_3];
_temp_index_4 := index2;
z_imag=data[_temp_index_4];
DynMap[4099] = DynMap[2048 + _temp_index_4];
DynMap[4110].Reliability = DynMap[4098].Reliability + DynMap[4102].Reliability - 1.0;
DynMap[4110].Delta = math.Abs(float64(w_real)) * DynMap[4102].Delta + math.Abs(float64(z_real)) * DynMap[4098].Delta + DynMap[4102].Delta*DynMap[4098].Delta;
temp6 = w_real*z_real;
DynMap[4111].Reliability = DynMap[4099].Reliability + DynMap[4103].Reliability - 1.0;
DynMap[4111].Delta = math.Abs(float64(w_imag)) * DynMap[4103].Delta + math.Abs(float64(z_imag)) * DynMap[4099].Delta + DynMap[4103].Delta*DynMap[4099].Delta;
temp7 = w_imag*z_imag;
DynMap[4112].Reliability = DynMap[4099].Reliability + DynMap[4102].Reliability - 1.0;
DynMap[4112].Delta = math.Abs(float64(w_real)) * DynMap[4102].Delta + math.Abs(float64(z_imag)) * DynMap[4099].Delta + DynMap[4102].Delta*DynMap[4099].Delta;
temp8 = w_real*z_imag;
DynMap[4113].Reliability = DynMap[4103].Reliability + DynMap[4098].Reliability - 1.0;
DynMap[4113].Delta = math.Abs(float64(w_imag)) * DynMap[4103].Delta + math.Abs(float64(z_real)) * DynMap[4098].Delta + DynMap[4103].Delta*DynMap[4098].Delta;
temp9 = w_imag*z_real;
DynMap[4100].Reliability = DynMap[4111].Reliability + DynMap[4110].Reliability - 1.0;
DynMap[4100].Delta = DynMap[4110].Delta + DynMap[4111].Delta;
t_real = temp6-temp7;
DynMap[4101].Reliability = DynMap[4113].Reliability + DynMap[4112].Reliability - 1.0;
DynMap[4101].Delta = DynMap[4112].Delta + DynMap[4113].Delta;
t_imag = temp8+temp9;
index3 = 2*i;
index4 = 2*i+1;
_temp_index_5 := index3;
temp1=data[_temp_index_5];
DynMap[4105] = DynMap[2048 + _temp_index_5];
_temp_index_6 := index4;
temp2=data[_temp_index_6];
DynMap[4106] = DynMap[2048 + _temp_index_6];
DynMap[4110].Reliability = DynMap[4105].Reliability + DynMap[4100].Reliability - 1.0;
DynMap[4110].Delta = DynMap[4105].Delta + DynMap[4100].Delta;
temp6 = temp1-t_real;
DynMap[4111].Reliability = DynMap[4106].Reliability + DynMap[4101].Reliability - 1.0;
DynMap[4111].Delta = DynMap[4106].Delta + DynMap[4101].Delta;
temp7 = temp2-t_imag;
_temp_index_7 := index1;
data[_temp_index_7]=temp6;
DynMap[2048 + _temp_index_7] = DynMap[4110];
_temp_index_8 := index2;
data[_temp_index_8]=temp7;
DynMap[2048 + _temp_index_8] = DynMap[4111];
DynMap[4107].Reliability = DynMap[4105].Reliability + DynMap[4100].Reliability - 1.0;
DynMap[4107].Delta = DynMap[4105].Delta + DynMap[4100].Delta;
temp3 = temp1+t_real;
_temp_index_9 := index3;
data[_temp_index_9]=temp3;
DynMap[2048 + _temp_index_9] = DynMap[4107];
DynMap[4108].Reliability = DynMap[4106].Reliability + DynMap[4101].Reliability - 1.0;
DynMap[4108].Delta = DynMap[4106].Delta + DynMap[4101].Delta;
temp4 = temp2+t_imag;
_temp_index_10 := index4;
data[_temp_index_10]=temp4;
DynMap[2048 + _temp_index_10] = DynMap[4108];
b = b+2*transform_length;
 }
DynMap[4110].Reliability = DynMap[4103].Reliability;
DynMap[4110].Delta = math.Abs(float64(s)) * DynMap[4103].Delta;
temp6 = s*w_imag;
DynMap[4111].Reliability = DynMap[4102].Reliability;
DynMap[4111].Delta = math.Abs(float64(s2)) * DynMap[4102].Delta;
temp7 = s2*w_real;
DynMap[4112].Reliability = DynMap[4111].Reliability + DynMap[4110].Reliability - 1.0;
DynMap[4112].Delta = DynMap[4110].Delta + DynMap[4111].Delta;
temp8 = temp6+temp7;
DynMap[4113].Reliability = DynMap[4102].Reliability;
DynMap[4113].Delta = math.Abs(float64(s)) * DynMap[4102].Delta;
temp9 = s*w_real;
DynMap[4114].Reliability = DynMap[4103].Reliability;
DynMap[4114].Delta = math.Abs(float64(s2)) * DynMap[4103].Delta;
temp10 = s2*w_imag;
DynMap[4115].Reliability = DynMap[4113].Reliability + DynMap[4114].Reliability - 1.0;
DynMap[4115].Delta = DynMap[4113].Delta + DynMap[4114].Delta;
temp11 = temp9-temp10;
DynMap[4100].Reliability = DynMap[4112].Reliability + DynMap[4102].Reliability - 1.0;
DynMap[4100].Delta = DynMap[4102].Delta + DynMap[4112].Delta;
t_real = w_real-temp8;
DynMap[4101].Reliability = DynMap[4103].Reliability + DynMap[4115].Reliability - 1.0;
DynMap[4101].Delta = DynMap[4103].Delta + DynMap[4115].Delta;
t_imag = w_imag+temp11;
DynMap[4102].Reliability = DynMap[4100].Reliability;
DynMap[4102].Delta = DynMap[4100].Delta;
w_real = t_real;
DynMap[4103].Reliability = DynMap[4101].Reliability;
DynMap[4103].Delta = DynMap[4101].Delta;
w_imag = t_imag;
a = a+1;
 }
bit = bit+1;
transform_length = transform_length*2;
 }
DynMap[4116] = diesel.ProbInterval{1, 0};
maxpsd = 0;
DynMap[4117] = diesel.ProbInterval{1, 0};
maxindex = 0;
i = 0;
DynMap[4097] = diesel.ProbInterval{1, 0};
di = 0;
for __temp_2 := 0; __temp_2 < N; __temp_2++ {
 index3 = 2*i;
index4 = 2*i+1;
_temp_index_11 := index3;
temp1=data[_temp_index_11];
DynMap[4105] = DynMap[2048 + _temp_index_11];
DynMap[4106].Reliability = DynMap[4105].Reliability;
DynMap[4106].Delta = math.Abs(float64(temp1)) * DynMap[4105].Delta + math.Abs(float64(temp1)) * DynMap[4105].Delta + DynMap[4105].Delta*DynMap[4105].Delta;
temp2 = temp1*temp1;
_temp_index_12 := index4;
temp3=data[_temp_index_12];
DynMap[4107] = DynMap[2048 + _temp_index_12];
DynMap[4108].Reliability = DynMap[4107].Reliability;
DynMap[4108].Delta = math.Abs(float64(temp3)) * DynMap[4107].Delta + math.Abs(float64(temp3)) * DynMap[4107].Delta + DynMap[4107].Delta*DynMap[4107].Delta;
temp4 = temp3*temp3;
DynMap[4109].Reliability = DynMap[4106].Reliability + DynMap[4108].Reliability - 1.0;
DynMap[4109].Delta = DynMap[4106].Delta + DynMap[4108].Delta;
temp5 = temp2+temp4;
DynMap[4110].Reliability = DynMap[4109].Reliability;
DynMap[4110].Delta =  DynMap[4109].Delta / math.Abs(float64(100.0));
temp6 = temp5/100.0;
maxpsd = diesel.DynCondFloat32GeqFloat32(temp6, maxpsd, DynMap[:], 4110, 4116, temp6, maxpsd, 4110, 4116, 4116);
maxindex = diesel.DynCondFloat32GeqInt(temp6, maxpsd, DynMap[:], 4110, 4116, di, maxindex, 4097, 4117, 4117);
i = i+1;
DynMap[4097].Reliability = DynMap[4097].Reliability;
DynMap[4097].Delta = DynMap[4097].Delta;
di = di+1;
 }
distance=getDistance(maxindex);
Outdata = data;
Distance = distance;


  fmt.Println("Ending thread : ", 0);
}

func main() {
  Num_threads = 1;
	
	diesel.InitChannels(1);

	rand.Seed(time.Now().UnixNano())

  dat, err := os.Open("signal.txt")
  check(err)
  defer dat.Close()

  scanner := bufio.NewScanner(dat)
	scanner.Split(bufio.ScanLines)

  i := 0
	for scanner.Scan() {
    Inputdata[i], err = strconv.ParseFloat(scanner.Text(), 64)
    check(err)
    i = i + 1
	}

	// for i,_ := range(Inputdata) {
	// 	Inputdata[i] = rand.Float64()
	// }

  fmt.Println("Starting the iterations")
  fmt.Println(Inputdata[:20])

  startTime := time.Now()

	go func_0();


	fmt.Println("Main thread waiting for others to finish");  
	diesel.Wg.Wait()
  elapsed := time.Since(startTime)

	fmt.Println("Done!");
  fmt.Println("Elapsed time : ", elapsed.Nanoseconds());

  fmt.Println(Outdata[:20])
}
