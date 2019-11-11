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
	"time"
  "parallely"
)

var Src []float64
var Dest [262144]float64
var ImgSize int
var NumThreads int
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

func minf(x, y float64) float64 {
  if x > y {
    return y
  }
  return x
}

func maxf(x, y int) float64 {
  if x > y {
    return float64(x)
  }
  return float64(y)
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

func exp(input float64) float64 {
	return math.Exp(input);
}

var Q = []int {1,2,3,4,5,6,7,8};


func func_0() {
  defer parallely.Wg.Done()
  var DynMap = map[int] float64{};
  _ = DynMap;
  var s_height int;
var s_width int;
var t_height int;
var ts_height int;
var i int;
var j int;
var k int;
var myrows int;
var lastthread int;
var te_height int;
var inputImage []float64;
 inputImage=make([]float64, ImgSize);
var dest_slice [32768]float64;
parallely.InitDynArray(0, 32768, DynMap);
var outImage [262144]float64;
parallely.InitDynArray(32768, 262144, DynMap);
var temp float64;
DynMap[294912] = 1;
s_width = SWidth;
s_height = SHeight;
inputImage = Src;
t_height = s_height/(NumThreads-1);
i = 0;
for _, q := range(Q) {
 parallely.SendFloat64Array(inputImage[:], 0, q);
parallely.SendInt(s_height, 0, q);
parallely.SendInt(s_width, 0, q);
ts_height = i*t_height;
lastthread = parallely.ConvBool(i==(NumThreads-1));
if lastthread != 0 {
 te_height = s_height;
 } else {
 te_height = (i+1)*t_height;
 }
parallely.SendInt(ts_height, 0, q);
parallely.SendInt(te_height, 0, q);
i = i+1;
 }
i = 0;
for _, q := range(Q) {
 parallely.ReceiveDynFloat64ArrayO1(dest_slice[:], 0, q, DynMap, 0);
ts_height = i*t_height;
lastthread = parallely.ConvBool(i==(NumThreads-1));
if lastthread != 0 {
 te_height = s_height;
 } else {
 te_height = (i+1)*t_height;
 }
myrows = te_height-ts_height;
j = 0;
for __temp_0 := 0; __temp_0 < myrows; __temp_0++ {
 k = 0;
for __temp_1 := 0; __temp_1 < s_width; __temp_1++ {
 _temp_index_1 := j*s_width+k;
temp=dest_slice[_temp_index_1];
_temp_index_2 := (ts_height+j)*s_width+k;
outImage[_temp_index_2]=temp;
k = k+1;
// DynMap[294912] = DynMap[0 + _temp_index_1];
// DynMap[32768 + _temp_index_2] = DynMap[294912];
DynMap[294912] = DynMap[0+_temp_index_1];
DynMap[32768 + _temp_index_2] = DynMap[0+_temp_index_1];
 }
j = j+1;
 }
i = i+1;
 }
Dest = outImage;

fmt.Println("----------------------------");

fmt.Println("Spec checkarray(outImage, 0.99): ", parallely.CheckArray(32768, 0.99, 262144, DynMap));

fmt.Println("----------------------------");


  fmt.Println("Ending thread : ", 0);
}
func func_Q(tid int) {
  defer parallely.Wg.Done()
  var DynMap = map[int] float64{};
  _ = DynMap;
  q := tid;
var image []float64;
 image=make([]float64, ImgSize);
var dest [32768]float64;
parallely.InitDynArray(294913, 32768, DynMap);
var ts_height int;
var i int;
var j int;
var myrows int;
var s_height int;
var s_width int;
var te_height int;
var rs int;
var wght float64;
var wsum float64;
var val float64;
var iy int;
var ix int;
var tempy int;
var tempx int;
var x float64;
var y float64;
var dsq int;
var temp1 float64;
var temp2 float64;
var temp3 float64;
var tempd float64;
DynMap[327681] = 1;
parallely.ReceiveFloat64Array(image[:], tid, 0);
parallely.ReceiveInt(&s_height, tid, 0);
parallely.ReceiveInt(&s_width, tid, 0);
parallely.ReceiveInt(&ts_height, tid, 0);
parallely.ReceiveInt(&te_height, tid, 0);
myrows = te_height-ts_height;
i = 0;
for __temp_2 := 0; __temp_2 < myrows; __temp_2++ {
 j = 0;
for __temp_3 := 0; __temp_3 < s_width; __temp_3++ {
 rs = 10;
wght = 1;
wsum = 1;
val = 1;
iy = (i+ts_height)-rs;
for __temp_4 := 0; __temp_4 < 20; __temp_4++ {
 ix = j-rs;
for __temp_5 := 0; __temp_5 < 20; __temp_5++ {
 temp1=maxf(0,ix);
temp2=convertToFloat(s_width);
x=minf(temp2-1.0,temp1);
temp1=maxf(0,iy);
temp2=convertToFloat(s_height);
y=minf(temp2-1.0,temp1);
dsq = (ix-j)*(ix-j)+(iy-(i+ts_height))*(iy-(i+ts_height));
temp1=convertToFloat(dsq);
temp2 = (temp1*-1)/16;
temp3=exp(temp2);
wght = temp3/3.1416*2*4*4;
tempy=floorInt(y);
tempx=floorInt(x);
_temp_index_1 := tempy*s_width+tempx;
temp1=image[_temp_index_1];
val = temp1*wght+val;
wsum = wsum+wght;
ix = ix+1;
 }
iy = iy+1;
 }
tempd = parallely.RandchoiceFloat64(float32(0.9999), val/wsum, 0);
_temp_index_2 := i*s_width+j;
dest[_temp_index_2]=tempd;
j = j+1;
// DynMap[327681] = 0.9999;
// DynMap[294913 + _temp_index_2] = DynMap[327681];
DynMap[327681] = 0.9999;
DynMap[294913 + _temp_index_2] = 0.9999;
 }
i = i+1;
 }
parallely.SendDynFloat64ArrayO1(dest[:], tid, 0, DynMap, 294913);

  fmt.Println("Ending thread : ", q);
}

func main() {
	// rand.Seed(time.Now().UTC().UnixNano())

  fmt.Println("Starting main thread");
  NumThreads = 9;
	
	parallely.InitChannels(9);

  iFile := os.Args[1]
  oFile := os.Args[2]
  
  src_tmp, s_width, s_height, _ := ReadPpmFile(iFile)
	Src = make([]float64, len(src_tmp))
  SHeight = s_height
  SWidth = s_width
  DestSize = len(src_tmp)*4*4

  for i, _ := range src_tmp {
      Src[i] = float64(src_tmp[i])
  }

  ImgSize = len(src_tmp)

	startTime := time.Now()
	go func_0();
for _, index := range Q {
go func_Q(index);
}


	fmt.Println("Main thread waiting for others to finish");  
	parallely.Wg.Wait()

	end := time.Now()
	elapsed := end.Sub(startTime)
	fmt.Println("Elapsed time :", elapsed.Nanoseconds())

  tmp_dest := make([]int, len(Dest))
  for i, _ := range Dest {
      tmp_dest[i] = int(Dest[i])
  }

  WritePpmFile(tmp_dest, s_width, s_height, oFile)
}