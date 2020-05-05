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
  "diesel"
)

var Src [262144]float64
var Dest [4194304]float64
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
  defer diesel.Wg.Done();
  var DynMap [4194310]diesel.ProbInterval;
  _ = DynMap;
  var s_height int;
	var s_width int;
	var d_width int;
	var d_height int;
	var i int;
	var j int;
	var ul float64;
	var ll float64;
	var ur float64;
	var lr float64;
	var u_w float64;
	var l_w float64;
	var ul_w float64;
	DynMap[0] = diesel.ProbInterval{1, 0};
	var ll_w float64;
	DynMap[1] = diesel.ProbInterval{1, 0};
	var ur_w float64;
	DynMap[2] = diesel.ProbInterval{1, 0};
	var lr_w float64;
	DynMap[3] = diesel.ProbInterval{1, 0};
	var tempf float64;
	DynMap[4] = diesel.ProbInterval{1, 0};
	var tempf1 float64;
	DynMap[5] = diesel.ProbInterval{1, 0};
	var si float64;
	var sj float64;
	var delta float64;
	var cond int;
	var previ int;
	var prevj int;
	var nexti int;
	var nextj int;
	var image [262144]float64;
	var outImage [4194304]float64;
	// diesel.InitDynArray(6, 4194304, DynMap[:]);
	s_width = SWidth;
	s_height = SHeight;
	image = Src;
	d_height = 4*SHeight;
	d_width = 4*s_width;
	i = 0;
	delta = 1/4.0;
	si = 0;
	for __temp_0 := 0; __temp_0 < d_height; __temp_0++ {
		j = 0;
		sj = 0.0;
		for __temp_1 := 0; __temp_1 < d_width; __temp_1++ {
			previ=floorInt(si);
			nexti=ceilInt(si);
			prevj=floorInt(sj);
			nextj=ceilInt(sj);
			cond = diesel.ConvBool(SHeight<=nexti);
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
			// DynMap[0] = diesel.ProbInterval{0.9999, 0};
			ll_w = diesel.RandchoiceFloat64(float32(0.9999), (1.0-u_w)*l_w, 0);
			// DynMap[1] = diesel.ProbInterval{0.9999, 0};
			ur_w = diesel.RandchoiceFloat64(float32(0.9999), u_w*(1.0-l_w), 0);
			// DynMap[2] = diesel.ProbInterval{0.9999, 0};
			lr_w = diesel.RandchoiceFloat64(float32(0.9999), (1.0-u_w)*(1.0-l_w), 0);
			// DynMap[3] = diesel.ProbInterval{0.9999, 0};
			tempf1 = diesel.RandchoiceFloat64(float32(0.9999), ul*ul_w, 0);
			// DynMap[5].Reliability = DynMap[0].Reliability * 0.9999;
			tempf1 = diesel.RandchoiceFloat64(float32(0.9999), tempf1+ur*ur_w, 0);
			// DynMap[5].Reliability = diesel.Max(0.0, DynMap[2].Reliability + DynMap[5].Reliability - float32(1)) * 0.9999;
			tempf1 = diesel.RandchoiceFloat64(float32(0.9999), tempf1+ll*ll_w, 0);
			// DynMap[5].Reliability = diesel.Max(0.0, DynMap[1].Reliability + DynMap[5].Reliability - float32(1)) * 0.9999;
			tempf1 = diesel.RandchoiceFloat64(float32(0.9999), tempf1+lr*lr_w, 0);
			// DynMap[5].Reliability = diesel.Max(0.0, DynMap[3].Reliability + DynMap[5].Reliability - float32(1)) * 0.9999;
			_temp_index_5 := i*d_width+j;
			outImage[_temp_index_5]=tempf1;
			// DynMap[6 + _temp_index_5] = DynMap[5];
			sj = sj+delta;
			j = j+1;
		}
		si = si+delta;
		i = i+1;
	}

	for iii:=0; iii<d_height*d_width; iii++ {
		DynMap[6].Reliability = 0.9999 * 0.9999 * 0.9999 * 0.9999 * 0.9999 * 0.9999 * 0.9999 * 0.9999
	}

	// fmt.Println("----------------------------");
	// fmt.Println("Spec checkarray(outImage, 0.99): ", diesel.CheckArray(6, 0.99, 4194304, DynMap[:]));
	// fmt.Println("----------------------------");

	Dest = outImage;


  fmt.Println("Ending thread : ", 0);
}

func main() {
	// rand.Seed(time.Now().UTC().UnixNano())

  fmt.Println("Starting main thread");
  NumThreads = 1;
	
	diesel.InitChannels(1);

  iFile := os.Args[1]
  oFile := os.Args[2]
  
  src_tmp, s_width, s_height, _ := ReadPpmFile(iFile)
  SHeight = s_height
  SWidth = s_width
  DestSize = len(src_tmp)*4*4

  for i, _ := range src_tmp {
		Src[i] = float64(src_tmp[i])
  }

  ImgSize = len(src_tmp)

	startTime := time.Now()
	go func_0();


	fmt.Println("Main thread waiting for others to finish");  
	diesel.Wg.Wait()

	end := time.Now()
	elapsed := end.Sub(startTime)
	fmt.Println("Elapsed time :", elapsed.Nanoseconds())

  tmp_dest := make([]int, len(Dest))
  for i, _ := range Dest {
		tmp_dest[i] = int(Dest[i])
  }

  WritePpmFile(tmp_dest, s_width*4, s_height*4, oFile)
}
