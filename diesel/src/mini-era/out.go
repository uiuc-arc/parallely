package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"diesel"
  "time"
  "bufio"
	"math"
  "unsafe"
	"math/rand"
  "os"
  "math/bits"
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

func getFloat32FromInt(i int) float32 {
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

func max(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func convertToFloat(x int) float64 {
	return float64(x)
}

func parseOutput(outstr string) (result int, confidence float32) {
	r := regexp.MustCompile(`Prediction: .*`)
	matches := r.FindAllString(outstr, 1)
	if len(matches) == 0 {
		fmt.Println("could not read the output")
		return -1, -1.0
	}
	outparts := strings.Fields(matches[0])
	cat, err1 := strconv.Atoi(outparts[1])
	conf, err2 := strconv.ParseFloat(outparts[2], 32)
	if err1 != nil || err2 != nil {
		fmt.Println("could not read the output")
		return -1, -1.0
	}
	return cat, float32(conf)
}

func readCamera() (result int, confidence float32) {
	cmd := exec.Command("python3", "mio_inference_single.py")
	cmd.Dir = "./CNN_MIO_KERAS/"

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running the python script")
	}

	return parseOutput(string(out))
}

var Num_threads int
var Category int
var Inputdata [2048]float64
var Outdata [2048]float32
var Pi = float32(3.141592653589)
var Distance float32



func func_0() {
  defer diesel.Wg.Done();
  var DynMap [4125]diesel.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var cat_val int;
var cat_conf float32;
var cat int;
DynMap[0] = diesel.ProbInterval{1, 0};
var tempnn int;
DynMap[1] = diesel.ProbInterval{1, 0};
var tempnn1 int;
DynMap[2] = diesel.ProbInterval{1, 0};
var tempnn0 int;
DynMap[3] = diesel.ProbInterval{1, 0};
var tempr int;
DynMap[4] = diesel.ProbInterval{1, 0};
var tempcomb int;
DynMap[5] = diesel.ProbInterval{1, 0};
var slow int;
DynMap[6] = diesel.ProbInterval{1, 0};
var radar_n float32;
var radar_fs float32;
var radar_alpha float32;
var radar_c float32;
var data64 [2048]float64;
diesel.InitDynArray(7, 2048, DynMap[:]);
var data [2048]float32;
diesel.InitDynArray(2055, 2048, DynMap[:]);
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
var ptemp0 float32;
var ptemp1 float32;
var di int;
DynMap[4103] = diesel.ProbInterval{1, 0};
var z_real float32;
DynMap[4104] = diesel.ProbInterval{1, 0};
var z_imag float32;
DynMap[4105] = diesel.ProbInterval{1, 0};
var t_real float32;
DynMap[4106] = diesel.ProbInterval{1, 0};
var t_imag float32;
DynMap[4107] = diesel.ProbInterval{1, 0};
var w_real float32;
DynMap[4108] = diesel.ProbInterval{1, 0};
var w_imag float32;
DynMap[4109] = diesel.ProbInterval{1, 0};
var temp0 float64;
DynMap[4110] = diesel.ProbInterval{1, 0};
var temp1 float32;
DynMap[4111] = diesel.ProbInterval{1, 0};
var temp2 float32;
DynMap[4112] = diesel.ProbInterval{1, 0};
var temp3 float32;
DynMap[4113] = diesel.ProbInterval{1, 0};
var temp4 float32;
DynMap[4114] = diesel.ProbInterval{1, 0};
var temp5 float32;
DynMap[4115] = diesel.ProbInterval{1, 0};
var temp6 float32;
DynMap[4116] = diesel.ProbInterval{1, 0};
var temp7 float32;
DynMap[4117] = diesel.ProbInterval{1, 0};
var temp8 float32;
DynMap[4118] = diesel.ProbInterval{1, 0};
var temp9 float32;
DynMap[4119] = diesel.ProbInterval{1, 0};
var temp10 float32;
DynMap[4120] = diesel.ProbInterval{1, 0};
var temp11 float32;
DynMap[4121] = diesel.ProbInterval{1, 0};
var maxpsd float32;
DynMap[4122] = diesel.ProbInterval{1, 0};
var maxindex int;
DynMap[4123] = diesel.ProbInterval{1, 0};
var distance float32;
DynMap[4124] = diesel.ProbInterval{1, 0};
var iter int;
radar_n = 10.0;
radar_fs = 204800.0;
radar_alpha = 30000000000.0;
radar_c = 300000000.0;
 var totaltime int64 ;
 totaltime = 0 ;
 var starttime time.Time ;
 var elapsed time.Duration ;
 starttime = time.Now() ;
iter = 0;
for __temp_0 := 0; __temp_0 < 10; __temp_0++ {
  fmt.Println("Running Iteration: " + strconv.Itoa(iter)) ;
 elapsed = time.Since(starttime) ;
 totaltime += elapsed.Nanoseconds() ;
cat_val,cat_conf=readCamera();
 starttime = time.Now() ;
cat=cat_val;
DynMap[0] = diesel.ProbInterval{cat_conf, 0.0};
N = 1024;
logN = 10;
sign = -1.0;
transform_length = 1;
data64=bitReverse(Inputdata,N,logN);
DynMap[7] = diesel.ProbInterval{1, 0};
i = 0;
for __temp_1 := 0; __temp_1 < 2048; __temp_1++ {
 _temp_index_1 := i;
temp0=data64[_temp_index_1];
DynMap[4110] = DynMap[7 + _temp_index_1];
temp2 = float32(temp0);
DynMap[4112].Reliability = DynMap[4110].Reliability;
 DynMap[4112].Delta = diesel.GetCastingError64to32(temp0, temp2);
_temp_index_2 := i;
data[_temp_index_2]=temp2;
DynMap[2055 + _temp_index_2] = DynMap[4112];
i = i+1;
 }
bit = 0;
for __temp_2 := 0; __temp_2 < 10; __temp_2++ {
 DynMap[4108] = diesel.ProbInterval{1, 0};
w_real = 1.0;
DynMap[4109] = diesel.ProbInterval{1, 0};
w_imag = 0.0;
ptemp0=getFloat(transform_length);
ptemp1 = sign*Pi;
theta = ptemp1/ptemp0;
s=getSin32(theta);
ptemp0 = 0.5;
ptemp1 = ptemp0*theta;
t=getSin32(ptemp1);
ptemp0 = 2.0;
ptemp1=getFloat32(temp0);
s2 = ptemp1*t*t;
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
DynMap[4104] = DynMap[2055 + _temp_index_3];
_temp_index_4 := index2;
z_imag=data[_temp_index_4];
DynMap[4105] = DynMap[2055 + _temp_index_4];
DynMap[4116].Reliability = DynMap[4104].Reliability + DynMap[4108].Reliability - 1.0;
DynMap[4116].Delta = math.Abs(float64(w_real)) * DynMap[4108].Delta + math.Abs(float64(z_real)) * DynMap[4104].Delta + DynMap[4108].Delta*DynMap[4104].Delta;
temp6 = w_real*z_real;
DynMap[4117].Reliability = DynMap[4105].Reliability + DynMap[4109].Reliability - 1.0;
DynMap[4117].Delta = math.Abs(float64(w_imag)) * DynMap[4109].Delta + math.Abs(float64(z_imag)) * DynMap[4105].Delta + DynMap[4109].Delta*DynMap[4105].Delta;
temp7 = w_imag*z_imag;
DynMap[4118].Reliability = DynMap[4105].Reliability + DynMap[4108].Reliability - 1.0;
DynMap[4118].Delta = math.Abs(float64(w_real)) * DynMap[4108].Delta + math.Abs(float64(z_imag)) * DynMap[4105].Delta + DynMap[4108].Delta*DynMap[4105].Delta;
temp8 = w_real*z_imag;
DynMap[4119].Reliability = DynMap[4109].Reliability + DynMap[4104].Reliability - 1.0;
DynMap[4119].Delta = math.Abs(float64(w_imag)) * DynMap[4109].Delta + math.Abs(float64(z_real)) * DynMap[4104].Delta + DynMap[4109].Delta*DynMap[4104].Delta;
temp9 = w_imag*z_real;
DynMap[4106].Reliability = DynMap[4117].Reliability + DynMap[4116].Reliability - 1.0;
DynMap[4106].Delta = DynMap[4116].Delta + DynMap[4117].Delta;
t_real = temp6-temp7;
DynMap[4107].Reliability = DynMap[4119].Reliability + DynMap[4118].Reliability - 1.0;
DynMap[4107].Delta = DynMap[4118].Delta + DynMap[4119].Delta;
t_imag = temp8+temp9;
index3 = 2*i;
index4 = 2*i+1;
_temp_index_5 := index3;
temp1=data[_temp_index_5];
DynMap[4111] = DynMap[2055 + _temp_index_5];
_temp_index_6 := index4;
temp2=data[_temp_index_6];
DynMap[4112] = DynMap[2055 + _temp_index_6];
DynMap[4116].Reliability = DynMap[4111].Reliability + DynMap[4106].Reliability - 1.0;
DynMap[4116].Delta = DynMap[4111].Delta + DynMap[4106].Delta;
temp6 = temp1-t_real;
DynMap[4117].Reliability = DynMap[4112].Reliability + DynMap[4107].Reliability - 1.0;
DynMap[4117].Delta = DynMap[4112].Delta + DynMap[4107].Delta;
temp7 = temp2-t_imag;
_temp_index_7 := index1;
data[_temp_index_7]=temp6;
DynMap[2055 + _temp_index_7] = DynMap[4116];
_temp_index_8 := index2;
data[_temp_index_8]=temp7;
DynMap[2055 + _temp_index_8] = DynMap[4117];
DynMap[4113].Reliability = DynMap[4111].Reliability + DynMap[4106].Reliability - 1.0;
DynMap[4113].Delta = DynMap[4111].Delta + DynMap[4106].Delta;
temp3 = temp1+t_real;
_temp_index_9 := index3;
data[_temp_index_9]=temp3;
DynMap[2055 + _temp_index_9] = DynMap[4113];
DynMap[4114].Reliability = DynMap[4112].Reliability + DynMap[4107].Reliability - 1.0;
DynMap[4114].Delta = DynMap[4112].Delta + DynMap[4107].Delta;
temp4 = temp2+t_imag;
_temp_index_10 := index4;
data[_temp_index_10]=temp4;
DynMap[2055 + _temp_index_10] = DynMap[4114];
b = b+2*transform_length;
 }
DynMap[4116].Reliability = DynMap[4109].Reliability;
DynMap[4116].Delta = math.Abs(float64(s)) * DynMap[4109].Delta;
temp6 = s*w_imag;
DynMap[4117].Reliability = DynMap[4108].Reliability;
DynMap[4117].Delta = math.Abs(float64(s2)) * DynMap[4108].Delta;
temp7 = s2*w_real;
DynMap[4118].Reliability = DynMap[4117].Reliability + DynMap[4116].Reliability - 1.0;
DynMap[4118].Delta = DynMap[4116].Delta + DynMap[4117].Delta;
temp8 = temp6+temp7;
DynMap[4119].Reliability = DynMap[4108].Reliability;
DynMap[4119].Delta = math.Abs(float64(s)) * DynMap[4108].Delta;
temp9 = s*w_real;
DynMap[4120].Reliability = DynMap[4109].Reliability;
DynMap[4120].Delta = math.Abs(float64(s2)) * DynMap[4109].Delta;
temp10 = s2*w_imag;
DynMap[4121].Reliability = DynMap[4119].Reliability + DynMap[4120].Reliability - 1.0;
DynMap[4121].Delta = DynMap[4119].Delta + DynMap[4120].Delta;
temp11 = temp9-temp10;
DynMap[4106].Reliability = DynMap[4118].Reliability + DynMap[4108].Reliability - 1.0;
DynMap[4106].Delta = DynMap[4108].Delta + DynMap[4118].Delta;
t_real = w_real-temp8;
DynMap[4107].Reliability = DynMap[4109].Reliability + DynMap[4121].Reliability - 1.0;
DynMap[4107].Delta = DynMap[4109].Delta + DynMap[4121].Delta;
t_imag = w_imag+temp11;
DynMap[4108].Reliability = DynMap[4106].Reliability;
DynMap[4108].Delta = DynMap[4106].Delta;
w_real = t_real;
DynMap[4109].Reliability = DynMap[4107].Reliability;
DynMap[4109].Delta = DynMap[4107].Delta;
w_imag = t_imag;
a = a+1;
 }
bit = bit+1;
transform_length = transform_length*2;
 }
DynMap[4122] = diesel.ProbInterval{1, 0};
maxpsd = 0;
DynMap[4123] = diesel.ProbInterval{1, 0};
maxindex = 0;
i = 0;
DynMap[4103] = diesel.ProbInterval{1, 0};
di = 0;
for __temp_3 := 0; __temp_3 < N; __temp_3++ {
 index3 = 2*i;
index4 = 2*i+1;
_temp_index_11 := index3;
temp1=data[_temp_index_11];
DynMap[4111] = DynMap[2055 + _temp_index_11];
DynMap[4112].Reliability = DynMap[4111].Reliability;
DynMap[4112].Delta = math.Abs(float64(temp1)) * DynMap[4111].Delta + math.Abs(float64(temp1)) * DynMap[4111].Delta + DynMap[4111].Delta*DynMap[4111].Delta;
temp2 = temp1*temp1;
_temp_index_12 := index4;
temp3=data[_temp_index_12];
DynMap[4113] = DynMap[2055 + _temp_index_12];
DynMap[4114].Reliability = DynMap[4113].Reliability;
DynMap[4114].Delta = math.Abs(float64(temp3)) * DynMap[4113].Delta + math.Abs(float64(temp3)) * DynMap[4113].Delta + DynMap[4113].Delta*DynMap[4113].Delta;
temp4 = temp3*temp3;
DynMap[4115].Reliability = DynMap[4112].Reliability + DynMap[4114].Reliability - 1.0;
DynMap[4115].Delta = DynMap[4112].Delta + DynMap[4114].Delta;
temp5 = temp2+temp4;
DynMap[4116].Reliability = DynMap[4115].Reliability;
DynMap[4116].Delta =  DynMap[4115].Delta / math.Abs(float64(100.0));
temp6 = temp5/100.0;
maxpsd = diesel.DynCondFloat32GeqFloat32(temp6, maxpsd, DynMap[:], 4116, 4122, temp6, maxpsd, 4116, 4122, 4122);
maxindex = diesel.DynCondFloat32GeqInt(temp6, maxpsd, DynMap[:], 4116, 4122, di, maxindex, 4103, 4123, 4123);
i = i+1;
DynMap[4103].Reliability = DynMap[4103].Reliability;
DynMap[4103].Delta = DynMap[4103].Delta;
di = di+1;
 }
temp6=getFloat32FromInt(maxindex);
DynMap[4116] = diesel.ProbInterval{1, 0};
DynMap[4111].Reliability = DynMap[4116].Reliability;
DynMap[4111].Delta = math.Abs(float64(radar_fs)) * DynMap[4116].Delta;
temp1 = temp6*radar_fs;
DynMap[4112].Reliability = DynMap[4111].Reliability;
DynMap[4112].Delta =  DynMap[4111].Delta / math.Abs(float64(radar_n));
temp2 = temp1/radar_n;
DynMap[4113] = diesel.ProbInterval{1, 0};
DynMap[4113] = diesel.ProbInterval{1, 0};
temp3 = 0.5*radar_c;
DynMap[4114].Reliability = DynMap[4113].Reliability + DynMap[4112].Reliability - 1.0;
DynMap[4114].Delta = math.Abs(float64(temp2)) * DynMap[4112].Delta + math.Abs(float64(temp3)) * DynMap[4113].Delta + DynMap[4112].Delta*DynMap[4113].Delta;
temp4 = temp2*temp3;
DynMap[4115].Reliability = DynMap[4114].Reliability;
DynMap[4115].Delta =  DynMap[4114].Delta / math.Abs(float64(radar_alpha));
temp5 = temp4/radar_alpha;
DynMap[4124].Reliability = DynMap[4115].Reliability;
DynMap[4124].Delta = DynMap[4115].Delta;
distance = temp5;
Outdata = data;
Distance = distance;
DynMap[1].Reliability = DynMap[0].Reliability;
tempnn = diesel.ConvBool(cat==1);
DynMap[2] = diesel.ProbInterval{1, 0};
tempnn1 = 1;
DynMap[3] = diesel.ProbInterval{1, 0};
tempnn0 = 0;
DynMap[4111] = diesel.ProbInterval{1, 0};
temp1 = 100;
DynMap[4].Reliability = DynMap[4124].Reliability + DynMap[4111].Reliability - 1.0;
tempr = diesel.ConvBool(distance<temp1);
DynMap[5].Reliability = DynMap[1].Reliability + DynMap[4].Reliability - 1.0;
tempcomb = diesel.ConvBool(tempnn==1 && tempr==1);
temp_bool_13:= tempcomb; if temp_bool_13 != 0 { slow  = tempnn1 } else { slow = tempnn0 };
if temp_bool_13 != 0 {DynMap[6].Reliability  = DynMap[5].Reliability * DynMap[2].Reliability} else { DynMap[6].Reliability = DynMap[5].Reliability * DynMap[3].Reliability};
iter = iter+1;
 }
 fmt.Println("Elapsed time : ", totaltime) ;
Category = slow;


  fmt.Println("Ending thread : ", 0);
}

func main() {
	fmt.Println("Starting main thread")

	Num_threads = 1

	diesel.InitChannels(1)

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

	fmt.Println("Starting the iterations")

	go func_0();


	// cmd := exec.Command("python3", "mio_inference_single.py")
	// cmd.Dir = "./CNN_MIO_KERAS/"
	// out, err := cmd.CombinedOutput()
	// if err != nil {
	// 	fmt.Println("Error running the python script")
	// }
	// fmt.Println(parseOutput(string(out)))

	fmt.Println("Main thread waiting for others to finish")
	diesel.Wg.Wait()

	fmt.Println("Done!")
}
