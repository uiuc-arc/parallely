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
  "parallely"
  "math/rand"
  "time"
)

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

func scale_kernel(i float64, j float64, src []int, s_height int, s_width int) int {
  previ := int(math.Floor(i))
  nexti := int(math.Ceil(i))
  prevj := int(math.Floor(j))
  nextj := int(math.Ceil(j))

  if s_height <= nexti {
    previ = Max(s_height-2,0)
    nexti = Min(previ+1,s_height-1)
  }
  if s_width <= nextj {
    prevj = Max(s_width-2,0)
    nextj = Min(prevj+1,s_width-1)
  }

  ul := src[Idx(previ,prevj,s_width)]
  ll := src[Idx(nexti,prevj,s_width)]
  ur := src[Idx(previ,nextj,s_width)]
  lr := src[Idx(nexti,nextj,s_width)]

  //BUG! weights can be 0 if i or j are integers
  //ul_w := (float64(nextj)-j)*(float64(nexti)-i)
  //ll_w := (float64(nextj)-j)*(i-float64(previ))
  //ur_w := (j-float64(prevj))*(float64(nexti)-i)
  //lr_w := (j-float64(prevj))*(i-float64(previ))

  u_w := float64(nextj)-j
  l_w := float64(nexti)-i
  ul_w := u_w*l_w
  ll_w := (1.0-u_w)*l_w
  ur_w := u_w*(1.0-l_w)
  lr_w := (1.0-u_w)*(1.0-l_w)
  

  return int(float64(ul)*ul_w+float64(ur)*ur_w+float64(ll)*ll_w+float64(lr)*lr_w)
}

func scale(f float64, src []int, s_width int, s_height int, dest []int, d_height int, d_width int) {
  overallflag := false

  si := 0.0
  delta := 1.0/f

  for i := 0; i < d_height; i++ {
    sj := 0.0
    for j := 0; j < d_width; j++ {
      //BEGIN try block
      var dest_pix int
      flag := false
      dest_pix = scale_kernel(si,sj,src,s_height,s_width)
      dest[Idx(i,j,d_width)] = parallely.RandchoiceFlag(0.999, dest_pix, 0, &flag)
      //CHECK
      if flag {
        //BEGIN redo block
        flag = false
        dest_pix = scale_kernel(si,sj,src,s_height,s_width)
        dest[Idx(i,j,d_width)] = parallely.RandchoiceFlag(0.9999, dest_pix, 0, &flag)
      }
      //END tcr block
      overallflag = overallflag || flag
      sj += delta
    }
    si += delta
  }
  if overallflag {
    fmt.Println(1)
  } else {
    fmt.Println(0)
  }
}

func main() {
  rand.Seed(time.Now().UTC().UnixNano())
  iFile := os.Args[1]
  //oFile := os.Args[2]
  f, _ := strconv.ParseFloat(os.Args[3],64)
  src, s_width, s_height, _ := ReadPpmFile(iFile)
  d_height := int(f*float64(s_height))
  d_width := int(f*float64(s_width))
  dest := make([]int, d_height*d_width)
  scale(f,src,s_width,s_height,dest,d_height,d_width)
  //WritePpmFile(dest,d_width,d_height,oFile)
}
