package main

import (
  "math"
  "errors"
  "io"
  "io/ioutil"
  "os"
  "regexp"
  "strconv"
  "fmt"
  "diesel"
)

var Src [262144]float64
var Dest [1048576]float64
var ImgSize int
var Num_threads int
var SWidth int
var SHeight int
var DestSize int

const (
  // single whitespace character
  ws = "[ \n\r\t\v\f]"
  // isolated comment
  cmt = "#[^\n\r]*"
  // comment sub expression
  cmts = "(" + ws + "*" + cmt + "[\n\r])"
  // number with leading comments
  num = "(" + cmts + "+" + ws + "*|" + ws + "+)([0-9]+)"
)

var rxHeader = regexp.MustCompile("^P6" + num + num + num + "(" + cmts + "*" + ")" + ws)
var rxComment = regexp.MustCompile(cmt)

// ReadFrom constructs a Bitmap object from an io.Reader.
func ReadPpmFrom(r io.Reader) (b []int, x int, y int, err error) {
  var all []byte
  all, err = ioutil.ReadAll(r)
  if err != nil {
    return
  }
  bss := rxHeader.FindSubmatch(all)
  if bss == nil {
    return nil, 0, 0, errors.New("unrecognized ppm header")
  }
  x, _ = strconv.Atoi(string(bss[3]))
  y, _ = strconv.Atoi(string(bss[6]))
  maxval, _ := strconv.Atoi(string(bss[9]))
  if maxval > 255 {
    return nil, 0, 0, errors.New("16 bit ppm not supported")
  }
  b = make([]int, x*y)
  b3 := all[len(bss[0]):]
  var n1 int
  for i := range b {
    b[i] = int(b3[n1]) * 255 / maxval
    n1 += 3
  }
  return
}

// ReadFile writes binary P6 format PPM from the specified filename.
func ReadPpmFile(fn string) (b []int, x int, y int, err error) {
  var f *os.File
  if f, err = os.Open(fn); err != nil {
    return
  }
  if b, x, y, err = ReadPpmFrom(f); err != nil {
    return
  }
  return b, x, y, f.Close()
}

// WriteTo outputs 8-bit P6 PPM format to an io.Writer.
func WritePpmTo(b []int, x int, y int, w io.Writer) (err error) {
  // magic number
  if _, err = fmt.Fprintln(w, "P6"); err != nil {
    return
  }
  // x, y, depth
  _, err = fmt.Fprintf(w, "%d %d\n255\n", x, y)
  if err != nil {
    return
  }
  // raster data in a single write
  b3 := make([]byte, 3*len(b))
  n1 := 0
  for _, px := range b {
    b3[n1] = byte(px)
    b3[n1+1] = byte(px)
    b3[n1+2] = byte(px)
    n1 += 3
  }
  if _, err = w.Write(b3); err != nil {
    return
  }
  return
}

// WriteFile writes to the specified filename.
func WritePpmFile(b []int, x int, y int, fn string) (err error) {
  var f *os.File
  if f, err = os.Create(fn); err != nil {
    return
  }
  if err = WritePpmTo(b,x,y,f); err != nil {
    return
  }
  return f.Close()
}

func Min(x, y int) int {
  if x < y {
    return x
  }
  return y
}

func Max(x, y int) int {
  if x > y {
    return x
  }
  return y
}

func Idx(i, j, width int) int {
  return i*width+j
}

func floorInt(input float64) int {
	return int(math.Floor(input));
}

func ceilInt(input float64) int {
	return int(math.Ceil(input));
}

func convertToFloat(x int) float64 {
	return float64(x)
}

var Q = []int {1,2,3,4,5,6,7,8};


func func_0() {
  diesel.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  diesel.WaitForWorkers(Num_threads)
  var DynMap [1179649]diesel.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var s_height int;
var s_width int;
var d_height int;
var d_width int;
var t_height int;
var ts_height int;
var i int;
var j int;
var k int;
var myrows int;
var lastthread int;
var te_height int;
var dest_slice [131072]float64;
diesel.InitDynArray(0, 131072, DynMap[:]);
var outImage [1048576]float64;
diesel.InitDynArray(131072, 1048576, DynMap[:]);
var temp float64;
DynMap[1179648] = diesel.ProbInterval{1, 0};
 diesel.StartTiming() ;
s_width = SWidth;
s_height = SHeight;
d_width = 2*s_width;
d_height = 2*s_height;
t_height = d_height/(Num_threads-1);
i = 0;
for _, q := range(Q) {
 diesel.SendInt(s_height, 0, q);
diesel.SendInt(s_width, 0, q);
diesel.SendInt(d_width, 0, q);
ts_height = i*t_height;
lastthread = diesel.ConvBool(i==(Num_threads-1));
if lastthread != 0 {
 te_height = d_height;
 } else {
 te_height = (i+1)*t_height;
 }
diesel.SendInt(ts_height, 0, q);
diesel.SendInt(te_height, 0, q);
i = i+1;
 }
i = 0;
for _, q := range(Q) {
 diesel.ReceiveDynFloat64ArrayO1(dest_slice[:], 0, q, DynMap[:], 0);
ts_height = i*t_height;
lastthread = diesel.ConvBool(i==(Num_threads-1));
if lastthread != 0 {
 te_height = d_height;
 } else {
 te_height = (i+1)*t_height;
 }
myrows = te_height-ts_height;
j = 0;
for __temp_0 := 0; __temp_0 < myrows; __temp_0++ {
 k = 0;
for __temp_1 := 0; __temp_1 < d_width; __temp_1++ {
 _temp_index_1 := j*d_width+k;
temp=dest_slice[_temp_index_1];
DynMap[1179648] = DynMap[0 + _temp_index_1];
_temp_index_2 := (ts_height+j)*d_width+k;
outImage[_temp_index_2]=temp;
DynMap[131072 + _temp_index_2] = DynMap[1179648];
k = k+1;
 }
j = j+1;
 }
i = i+1;
 }
 diesel.EndTiming() ;
Dest = outImage;


  diesel.CleanupMain()
  fmt.Println("Ending thread : ", 0);
}
func func_Q(tid int) {
  diesel.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  diesel.PingMain()
  var DynMap [131073]diesel.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  q := tid;
var image [262144]float64;
var dest [131072]float64;
diesel.InitDynArray(0, 131072, DynMap[:]);
var ts_height int;
var i int;
var j int;
var myrows int;
var si float64;
var sj float64;
var delta float64;
var s_height int;
var s_width int;
var te_height int;
var d_width int;
var cond int;
var previ int;
var prevj int;
var nexti int;
var nextj int;
var ul float64;
var ll float64;
var ur float64;
var lr float64;
var u_w float64;
var l_w float64;
var ul_w float64;
var ll_w float64;
var ur_w float64;
var lr_w float64;
var tempf float64;
var tempf1 float64;
DynMap[131072] = diesel.ProbInterval{1, 0};
image = Src;
diesel.ReceiveInt(&s_height, tid, 0);
diesel.ReceiveInt(&s_width, tid, 0);
diesel.ReceiveInt(&d_width, tid, 0);
diesel.ReceiveInt(&ts_height, tid, 0);
diesel.ReceiveInt(&te_height, tid, 0);
myrows = te_height-ts_height;
i = 0;
delta = 1/2.0;
tempf=convertToFloat(ts_height);
si = tempf*delta;
for __temp_2 := 0; __temp_2 < myrows; __temp_2++ {
 j = 0;
sj = 0.0;
for __temp_3 := 0; __temp_3 < d_width; __temp_3++ {
 previ=floorInt(si);
nexti=ceilInt(si);
prevj=floorInt(sj);
nextj=ceilInt(sj);
cond = diesel.ConvBool(s_height<=nexti);
if cond != 0 {
 previ = s_height-2;
nexti = s_height-1;
 }
cond = diesel.ConvBool(s_width<=nextj);
if cond != 0 {
 prevj = s_width-2;
nextj = s_width-1;
 }
_temp_index_1 := previ*s_width+prevj;
ul=image[_temp_index_1];
_temp_index_2 := nexti*s_width+prevj;
ll=image[_temp_index_2];
_temp_index_3 := previ*s_width+nextj;
ur=image[_temp_index_3];
_temp_index_4 := nexti*s_width+nextj;
lr=image[_temp_index_4];
tempf=convertToFloat(nexti);
u_w = tempf-si;
tempf=convertToFloat(nextj);
l_w = tempf-sj;
ul_w = diesel.RandchoiceFloat64(float32(0.9999), u_w*l_w, 0);
ll_w = (1.0-u_w)*l_w;
ur_w = u_w*(1.0-l_w);
lr_w = (1.0-u_w)*(1.0-l_w);
tempf1 = diesel.RandchoiceFloat64(float32(0.9999), ul*ul_w, 0);
DynMap[131072] = diesel.ProbInterval{0.9999, 0};
DynMap[131072].Reliability = DynMap[131072].Reliability;
DynMap[131072].Delta = DynMap[131072].Delta;
tempf1 = tempf1+ur*ur_w;
DynMap[131072].Reliability = DynMap[131072].Reliability;
DynMap[131072].Delta = DynMap[131072].Delta;
tempf1 = tempf1+ll*ll_w;
DynMap[131072].Reliability = DynMap[131072].Reliability;
DynMap[131072].Delta = DynMap[131072].Delta;
tempf1 = tempf1+lr*lr_w;
tempf1 = diesel.RandchoiceFloat64(float32(0.99), tempf1, 0);
DynMap[131072].Reliability = DynMap[131072].Reliability * 0.99;
_temp_index_5 := i*d_width+j;
dest[_temp_index_5]=tempf1;
DynMap[0 + _temp_index_5] = DynMap[131072];
sj = sj+delta;
j = j+1;
 }
si = si+delta;
i = i+1;
 }
diesel.SendDynFloat64ArrayO1(dest[:], tid, 0, DynMap[:], 0);

  diesel.CleanupMain()
  fmt.Println("Ending thread : ", q);
}

func main() {
  Num_threads = 9;	

  iFile := "./baboon.ppm"

	tid, _ := strconv.Atoi(os.Args[1])	
	fmt.Println("Starting worker thread: ", tid)
  
  src_tmp, s_width, s_height, _ := ReadPpmFile(iFile)
  SHeight = s_height
  SWidth = s_width
  DestSize = len(src_tmp)*2*2

  for i, _ := range src_tmp {
      Src[i] = float64(src_tmp[i])
  }

  ImgSize = len(src_tmp)

	// startTime := time.Now()
	func_Q(tid)

	// fmt.Println("Main thread waiting for others to finish");  
	// diesel.Wg.Wait()

	// end := time.Now()
	// elapsed := end.Sub(startTime)
	// fmt.Println("Elapsed time :", elapsed.Nanoseconds())

  // tmp_dest := make([]int, len(Dest))
  // for i, _ := range Dest {
  //     tmp_dest[i] = int(Dest[i])
  // }

  // WritePpmFile(tmp_dest, s_width*4, s_height*4, oFile)
}
