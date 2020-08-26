package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"dieseldist"
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

var Q = []int {1};


func func_0() {
  dieseldist.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldist.WaitForWorkers(Num_threads)
  var DynMap [9]dieseldist.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var cat_val int;
var cat_conf float32;
var cat int;
DynMap[0] = dieseldist.ProbInterval{1, 0};
var tempnn int;
DynMap[1] = dieseldist.ProbInterval{1, 0};
var tempnn1 int;
DynMap[2] = dieseldist.ProbInterval{1, 0};
var tempnn0 int;
DynMap[3] = dieseldist.ProbInterval{1, 0};
var tempr int;
DynMap[4] = dieseldist.ProbInterval{1, 0};
var tempcomb int;
DynMap[5] = dieseldist.ProbInterval{1, 0};
var distance float32;
DynMap[6] = dieseldist.ProbInterval{1, 0};
var slow int;
DynMap[7] = dieseldist.ProbInterval{1, 0};
var temp1 float32;
DynMap[8] = dieseldist.ProbInterval{1, 0};
var iter int;
 dieseldist.StartTiming() ;
iter = 0;
for __temp_0 := 0; __temp_0 < 10; __temp_0++ {
  fmt.Println(iter) ;
dieseldist.SendInt(iter, 0, 1);
 dieseldist.StartTimerPause() ;
cat_val,cat_conf=readCamera();
 dieseldist.StopTimerPause() ;
cat=cat_val;
DynMap[0] = dieseldist.ProbInterval{cat_conf, 0.0};
dieseldist.ReceiveFloat32(&distance, 0, 1);
my_chan_index = 1 * dieseldist.Numprocesses + 0;
__temp_rec_val_2 := dieseldist.GetDynValue(my_chan_index);
DynMap[6] = __temp_rec_val_2;
Distance = distance;
DynMap[1].Reliability = DynMap[0].Reliability;
tempnn = dieseldist.ConvBool(cat==1);
DynMap[2] = dieseldist.ProbInterval{1, 0};
tempnn1 = 1;
DynMap[3] = dieseldist.ProbInterval{1, 0};
tempnn0 = 0;
DynMap[8] = dieseldist.ProbInterval{1, 0};
temp1 = 100;
DynMap[4].Reliability = DynMap[6].Reliability + DynMap[8].Reliability - 1.0;
tempr = dieseldist.ConvBool(distance<temp1);
DynMap[5].Reliability = DynMap[1].Reliability + DynMap[4].Reliability - 1.0;
tempcomb = dieseldist.ConvBool(tempnn==1 && tempr==1);
temp_bool_1:= tempcomb; if temp_bool_1 != 0 { slow  = tempnn1 } else { slow = tempnn0 };
if temp_bool_1 != 0 {DynMap[7].Reliability  = DynMap[5].Reliability * DynMap[2].Reliability} else { DynMap[7].Reliability = DynMap[5].Reliability * DynMap[3].Reliability};
iter = iter+1;
 }
 dieseldist.EndTiming() ;
Category = slow;


  dieseldist.CleanupMain()
  fmt.Println("Ending thread : ", 0);
}
func func_Q(tid int) {
  dieseldist.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldist.PingMain(tid)
  var DynMap [4118]dieseldist.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  q := tid;
var radar_n float32;
var radar_fs float32;
var radar_alpha float32;
var radar_c float32;
var data64 [2048]float64;
dieseldist.InitDynArray(0, 2048, DynMap[:]);
var data [2048]float32;
dieseldist.InitDynArray(2048, 2048, DynMap[:]);
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
DynMap[4096] = dieseldist.ProbInterval{1, 0};
var z_real float32;
DynMap[4097] = dieseldist.ProbInterval{1, 0};
var z_imag float32;
DynMap[4098] = dieseldist.ProbInterval{1, 0};
var t_real float32;
DynMap[4099] = dieseldist.ProbInterval{1, 0};
var t_imag float32;
DynMap[4100] = dieseldist.ProbInterval{1, 0};
var w_real float32;
DynMap[4101] = dieseldist.ProbInterval{1, 0};
var w_imag float32;
DynMap[4102] = dieseldist.ProbInterval{1, 0};
var temp0 float64;
DynMap[4103] = dieseldist.ProbInterval{1, 0};
var temp1 float32;
DynMap[4104] = dieseldist.ProbInterval{1, 0};
var temp2 float32;
DynMap[4105] = dieseldist.ProbInterval{1, 0};
var temp3 float32;
DynMap[4106] = dieseldist.ProbInterval{1, 0};
var temp4 float32;
DynMap[4107] = dieseldist.ProbInterval{1, 0};
var temp5 float32;
DynMap[4108] = dieseldist.ProbInterval{1, 0};
var temp6 float32;
DynMap[4109] = dieseldist.ProbInterval{1, 0};
var temp7 float32;
DynMap[4110] = dieseldist.ProbInterval{1, 0};
var temp8 float32;
DynMap[4111] = dieseldist.ProbInterval{1, 0};
var temp9 float32;
DynMap[4112] = dieseldist.ProbInterval{1, 0};
var temp10 float32;
DynMap[4113] = dieseldist.ProbInterval{1, 0};
var temp11 float32;
DynMap[4114] = dieseldist.ProbInterval{1, 0};
var maxpsd float32;
DynMap[4115] = dieseldist.ProbInterval{1, 0};
var maxindex int;
DynMap[4116] = dieseldist.ProbInterval{1, 0};
var distance float32;
DynMap[4117] = dieseldist.ProbInterval{1, 0};
var iter int;
radar_n = 10.0;
radar_fs = 204800.0;
radar_alpha = 30000000000.0;
radar_c = 300000000.0;
iter = 0;
for __temp_2 := 0; __temp_2 < 10; __temp_2++ {
 dieseldist.ReceiveInt(&iter, tid, 0);
N = 1024;
logN = 10;
sign = -1.0;
transform_length = 1;
data64=bitReverse(Inputdata,N,logN);
i = 0;
for __temp_3 := 0; __temp_3 < 2048; __temp_3++ {
 _temp_index_1 := i;
temp0=data64[_temp_index_1];
DynMap[4103] = DynMap[0 + _temp_index_1];
temp2 = float32(temp0);
DynMap[4105].Reliability = DynMap[4103].Reliability;
 DynMap[4105].Delta = dieseldist.GetCastingError64to32(temp0, temp2);
_temp_index_2 := i;
data[_temp_index_2]=temp2;
DynMap[2048 + _temp_index_2] = DynMap[4105];
i = i+1;
 }
bit = 0;
for __temp_4 := 0; __temp_4 < 10; __temp_4++ {
 DynMap[4101] = dieseldist.ProbInterval{1, 0};
w_real = 1.0;
DynMap[4102] = dieseldist.ProbInterval{1, 0};
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
DynMap[4097] = DynMap[2048 + _temp_index_3];
_temp_index_4 := index2;
z_imag=data[_temp_index_4];
DynMap[4098] = DynMap[2048 + _temp_index_4];
DynMap[4109].Reliability = DynMap[4097].Reliability + DynMap[4101].Reliability - 1.0;
DynMap[4109].Delta = math.Abs(float64(w_real)) * DynMap[4101].Delta + math.Abs(float64(z_real)) * DynMap[4097].Delta + DynMap[4101].Delta*DynMap[4097].Delta;
temp6 = w_real*z_real;
DynMap[4110].Reliability = DynMap[4098].Reliability + DynMap[4102].Reliability - 1.0;
DynMap[4110].Delta = math.Abs(float64(w_imag)) * DynMap[4102].Delta + math.Abs(float64(z_imag)) * DynMap[4098].Delta + DynMap[4102].Delta*DynMap[4098].Delta;
temp7 = w_imag*z_imag;
DynMap[4111].Reliability = DynMap[4098].Reliability + DynMap[4101].Reliability - 1.0;
DynMap[4111].Delta = math.Abs(float64(w_real)) * DynMap[4101].Delta + math.Abs(float64(z_imag)) * DynMap[4098].Delta + DynMap[4101].Delta*DynMap[4098].Delta;
temp8 = w_real*z_imag;
DynMap[4112].Reliability = DynMap[4102].Reliability + DynMap[4097].Reliability - 1.0;
DynMap[4112].Delta = math.Abs(float64(w_imag)) * DynMap[4102].Delta + math.Abs(float64(z_real)) * DynMap[4097].Delta + DynMap[4102].Delta*DynMap[4097].Delta;
temp9 = w_imag*z_real;
DynMap[4099].Reliability = DynMap[4110].Reliability + DynMap[4109].Reliability - 1.0;
DynMap[4099].Delta = DynMap[4109].Delta + DynMap[4110].Delta;
t_real = temp6-temp7;
DynMap[4100].Reliability = DynMap[4112].Reliability + DynMap[4111].Reliability - 1.0;
DynMap[4100].Delta = DynMap[4111].Delta + DynMap[4112].Delta;
t_imag = temp8+temp9;
index3 = 2*i;
index4 = 2*i+1;
_temp_index_5 := index3;
temp1=data[_temp_index_5];
DynMap[4104] = DynMap[2048 + _temp_index_5];
_temp_index_6 := index4;
temp2=data[_temp_index_6];
DynMap[4105] = DynMap[2048 + _temp_index_6];
DynMap[4109].Reliability = DynMap[4104].Reliability + DynMap[4099].Reliability - 1.0;
DynMap[4109].Delta = DynMap[4104].Delta + DynMap[4099].Delta;
temp6 = temp1-t_real;
DynMap[4110].Reliability = DynMap[4105].Reliability + DynMap[4100].Reliability - 1.0;
DynMap[4110].Delta = DynMap[4105].Delta + DynMap[4100].Delta;
temp7 = temp2-t_imag;
_temp_index_7 := index1;
data[_temp_index_7]=temp6;
DynMap[2048 + _temp_index_7] = DynMap[4109];
_temp_index_8 := index2;
data[_temp_index_8]=temp7;
DynMap[2048 + _temp_index_8] = DynMap[4110];
DynMap[4106].Reliability = DynMap[4104].Reliability + DynMap[4099].Reliability - 1.0;
DynMap[4106].Delta = DynMap[4104].Delta + DynMap[4099].Delta;
temp3 = temp1+t_real;
_temp_index_9 := index3;
data[_temp_index_9]=temp3;
DynMap[2048 + _temp_index_9] = DynMap[4106];
DynMap[4107].Reliability = DynMap[4105].Reliability + DynMap[4100].Reliability - 1.0;
DynMap[4107].Delta = DynMap[4105].Delta + DynMap[4100].Delta;
temp4 = temp2+t_imag;
_temp_index_10 := index4;
data[_temp_index_10]=temp4;
DynMap[2048 + _temp_index_10] = DynMap[4107];
b = b+2*transform_length;
 }
DynMap[4109].Reliability = DynMap[4102].Reliability;
DynMap[4109].Delta = math.Abs(float64(s)) * DynMap[4102].Delta;
temp6 = s*w_imag;
DynMap[4110].Reliability = DynMap[4101].Reliability;
DynMap[4110].Delta = math.Abs(float64(s2)) * DynMap[4101].Delta;
temp7 = s2*w_real;
DynMap[4111].Reliability = DynMap[4110].Reliability + DynMap[4109].Reliability - 1.0;
DynMap[4111].Delta = DynMap[4109].Delta + DynMap[4110].Delta;
temp8 = temp6+temp7;
DynMap[4112].Reliability = DynMap[4101].Reliability;
DynMap[4112].Delta = math.Abs(float64(s)) * DynMap[4101].Delta;
temp9 = s*w_real;
DynMap[4113].Reliability = DynMap[4102].Reliability;
DynMap[4113].Delta = math.Abs(float64(s2)) * DynMap[4102].Delta;
temp10 = s2*w_imag;
DynMap[4114].Reliability = DynMap[4112].Reliability + DynMap[4113].Reliability - 1.0;
DynMap[4114].Delta = DynMap[4112].Delta + DynMap[4113].Delta;
temp11 = temp9-temp10;
DynMap[4099].Reliability = DynMap[4111].Reliability + DynMap[4101].Reliability - 1.0;
DynMap[4099].Delta = DynMap[4101].Delta + DynMap[4111].Delta;
t_real = w_real-temp8;
DynMap[4100].Reliability = DynMap[4102].Reliability + DynMap[4114].Reliability - 1.0;
DynMap[4100].Delta = DynMap[4102].Delta + DynMap[4114].Delta;
t_imag = w_imag+temp11;
DynMap[4101].Reliability = DynMap[4099].Reliability;
DynMap[4101].Delta = DynMap[4099].Delta;
w_real = t_real;
DynMap[4102].Reliability = DynMap[4100].Reliability;
DynMap[4102].Delta = DynMap[4100].Delta;
w_imag = t_imag;
a = a+1;
 }
bit = bit+1;
transform_length = transform_length*2;
 }
DynMap[4115] = dieseldist.ProbInterval{1, 0};
maxpsd = 0;
DynMap[4116] = dieseldist.ProbInterval{1, 0};
maxindex = 0;
i = 0;
DynMap[4096] = dieseldist.ProbInterval{1, 0};
di = 0;
for __temp_5 := 0; __temp_5 < N; __temp_5++ {
 index3 = 2*i;
index4 = 2*i+1;
_temp_index_11 := index3;
temp1=data[_temp_index_11];
DynMap[4104] = DynMap[2048 + _temp_index_11];
DynMap[4105].Reliability = DynMap[4104].Reliability;
DynMap[4105].Delta = math.Abs(float64(temp1)) * DynMap[4104].Delta + math.Abs(float64(temp1)) * DynMap[4104].Delta + DynMap[4104].Delta*DynMap[4104].Delta;
temp2 = temp1*temp1;
_temp_index_12 := index4;
temp3=data[_temp_index_12];
DynMap[4106] = DynMap[2048 + _temp_index_12];
DynMap[4107].Reliability = DynMap[4106].Reliability;
DynMap[4107].Delta = math.Abs(float64(temp3)) * DynMap[4106].Delta + math.Abs(float64(temp3)) * DynMap[4106].Delta + DynMap[4106].Delta*DynMap[4106].Delta;
temp4 = temp3*temp3;
DynMap[4108].Reliability = DynMap[4105].Reliability + DynMap[4107].Reliability - 1.0;
DynMap[4108].Delta = DynMap[4105].Delta + DynMap[4107].Delta;
temp5 = temp2+temp4;
DynMap[4109].Reliability = DynMap[4108].Reliability;
DynMap[4109].Delta =  DynMap[4108].Delta / math.Abs(float64(100.0));
temp6 = temp5/100.0;
maxpsd = dieseldist.DynCondFloat32GeqFloat32(temp6, maxpsd, DynMap[:], 4109, 4115, temp6, maxpsd, 4109, 4115, 4115);
maxindex = dieseldist.DynCondFloat32GeqInt(temp6, maxpsd, DynMap[:], 4109, 4115, di, maxindex, 4096, 4116, 4116);
i = i+1;
DynMap[4096].Reliability = DynMap[4096].Reliability;
DynMap[4096].Delta = DynMap[4096].Delta;
di = di+1;
 }
temp6=getFloat32FromInt(maxindex);
DynMap[4104].Reliability = DynMap[4109].Reliability;
DynMap[4104].Delta = math.Abs(float64(radar_fs)) * DynMap[4109].Delta;
temp1 = temp6*radar_fs;
DynMap[4105].Reliability = DynMap[4104].Reliability;
DynMap[4105].Delta =  DynMap[4104].Delta / math.Abs(float64(radar_n));
temp2 = temp1/radar_n;
DynMap[4106] = dieseldist.ProbInterval{1, 0};
DynMap[4106] = dieseldist.ProbInterval{1, 0};
temp3 = 0.5*radar_c;
DynMap[4107].Reliability = DynMap[4106].Reliability + DynMap[4105].Reliability - 1.0;
DynMap[4107].Delta = math.Abs(float64(temp2)) * DynMap[4105].Delta + math.Abs(float64(temp3)) * DynMap[4106].Delta + DynMap[4105].Delta*DynMap[4106].Delta;
temp4 = temp2*temp3;
DynMap[4108].Reliability = DynMap[4107].Reliability;
DynMap[4108].Delta =  DynMap[4107].Delta / math.Abs(float64(radar_alpha));
temp5 = temp4/radar_alpha;
DynMap[4117].Reliability = DynMap[4108].Reliability;
DynMap[4117].Delta = DynMap[4108].Delta;
distance = temp5;
dieseldist.SendFloat32(distance, tid, 0);
dieseldist.SendDynVal(DynMap[4117], tid, 0);
iter = iter+1;
 }

  fmt.Println("Ending thread : ", q);
}

func main() {
	fmt.Println("Starting main thread")

	Num_threads = 2

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

	func_0();


	fmt.Println("Done!")
}
