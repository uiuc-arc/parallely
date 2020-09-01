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
  "dieseldistrel"
)

var Src [262144]float64
var Dest [262144]float64
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
  dieseldistrel.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldistrel.WaitForWorkers(Num_threads)
  var DynMap [294913]float32;
  var my_chan_index int;
  _ = my_chan_index;
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
	var dest_slice [32768]float64;
	dieseldistrel.InitDynArray(0, 32768, DynMap[:]);
	var outImage [262144]float64;
	dieseldistrel.InitDynArray(32768, 262144, DynMap[:]);
	var temp float64;
	DynMap[294912] = 1.0;
	dieseldist.StartTiming() ;
	s_width = SWidth;
	s_height = SHeight;
	t_height = s_height/(Num_threads-1);
	i = 0;
	for _, q := range(Q) {
		ts_height = i*t_height;
		lastthread = dieseldistrel.ConvBool(i==(Num_threads-1));
		if lastthread != 0 {
			te_height = s_height;
		} else {
			te_height = (i+1)*t_height;
		}
		dieseldistrel.SendInt(ts_height, 0, q);
		dieseldistrel.SendInt(te_height, 0, q);
		i = i+1;
	}
	i = 0;
	for _, q := range(Q) {
		dieseldistrel.ReceiveDynFloat64Array(dest_slice[:], 0, q, DynMap[:], 0);
		ts_height = i*t_height;
		lastthread = dieseldistrel.ConvBool(i==(Num_threads-1));
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
				DynMap[294912] = DynMap[0 + _temp_index_1];
				_temp_index_2 := (ts_height+j)*s_width+k;
				outImage[_temp_index_2]=temp;
				DynMap[32768 + _temp_index_2] = DynMap[294912];
				k = k+1;
			}
			j = j+1;
		}
		i = i+1;
	}
	dieseldist.EndTiming() ;
	Dest = outImage;


  dieseldistrel.CleanupMain()
  fmt.Println("Ending thread : ", 0);
}
func func_Q(tid int) {
  dieseldistrel.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldistrel.PingMain(tid)
  var DynMap [32769]float32;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  q := tid;
	var image [262144]float64;
	var dest [32768]float64;
	dieseldistrel.InitDynArray(0, 32768, DynMap[:]);
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
	var temp0 int;
	var temp1 float64;
	var temp2 float64;
	var temp3 float64;
	var temp4 float64;
	var tempd float64;
	DynMap[32768] = 1.0;
	image = Src;
	s_width = SWidth;
	s_height = SHeight;
	dieseldistrel.ReceiveInt(&ts_height, tid, 0);
	dieseldistrel.ReceiveInt(&te_height, tid, 0);
	myrows = te_height-ts_height;
	i = 0;
	for __temp_2 := 0; __temp_2 < myrows; __temp_2++ {
		j = 0;
		for __temp_3 := 0; __temp_3 < s_width; __temp_3++ {
			rs = 5;
			wght = 1;
			wsum = 1;
			val = 1;
			iy = (i+ts_height)-rs;
			for __temp_4 := 0; __temp_4 < 5; __temp_4++ {
				ix = j-rs;
				for __temp_5 := 0; __temp_5 < 10; __temp_5++ {
					temp0 = 0;
					temp1=maxf(temp0,ix);
					temp2=convertToFloat(s_width);
					temp4 = temp2-1.0;
					x=minf(temp4,temp1);
					temp1=maxf(temp0,iy);
					temp2=convertToFloat(s_height);
					temp4 = temp2-1.0;
					y=minf(temp4,temp1);
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
					wsum = dieseldistrel.RandchoiceFloat64(float32(0.9999), wsum+wght, -1);
					ix = ix+1;
				}
				iy = iy+1;
			}
			tempd = dieseldistrel.RandchoiceFloat64(float32(0.9999), val/wsum, 0);
			DynMap[32768] = 0.9999;
			_temp_index_2 := i*s_width+j;
			dest[_temp_index_2]=tempd;
			DynMap[0 + _temp_index_2] = DynMap[32768];
			j = j+1;
		}
		i = i+1;
	}
	dieseldistrel.SendDynFloat64Array(dest[:], tid, 0, DynMap[:], 0);

  fmt.Println("Ending thread : ", q);
}

func main() {
  fmt.Println("Starting main thread");
  Num_threads = 9;

  iFile := "temp-512.ppm"
  
  src_tmp, s_width, s_height, err := ReadPpmFile(iFile)
  if err != nil {
		fmt.Println("[ERROR] Input does not exist")
		os.Exit(-1)
  }

  SHeight = s_height
  SWidth = s_width
  DestSize = len(src_tmp)*4*4

  for i, _ := range src_tmp {
		Src[i] = float64(src_tmp[i])
  }

  ImgSize = len(src_tmp)

	func_0();


  // tmp_dest := make([]int, len(Dest))
  // for i, _ := range Dest {
  //     tmp_dest[i] = int(Dest[i])
  // }
  // WritePpmFile(tmp_dest, s_width, s_height, oFile)
}
