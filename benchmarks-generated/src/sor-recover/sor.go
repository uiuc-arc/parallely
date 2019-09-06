package main

import (
  "os"
  "strconv"
  "fmt"
  "parallely"
  "math/rand"
  "time"
  "math"
)

func Idx(i, j, width int) int {
  return i*width+j
}

func main() {
  rand.Seed(time.Now().UTC().UnixNano())
  height, _ := strconv.Atoi(os.Args[1])
  width, _ := strconv.Atoi(os.Args[2])
  iterations, _ := strconv.Atoi(os.Args[3])

  array := make([]float64, height*width)
  for i := 0; i < height*width; i++ {
    array[i] = float64(i%23)
  }
  result := make([]float64, height*width)
  result_exact := make([]float64, height*width)

  overallflag := false

  for iter := 0; iter < iterations; iter++ {
    for i := 1; i < height-1; i++ {
      for j := 1; j < width-1; j++ {
        flag := false
        var pix float64
        pix = 0.2*(array[Idx(i,j,width)] + array[Idx(i-1,j,width)] + array[Idx(i+1,j,width)] + array[Idx(i,j-1,width)] + array[Idx(i,j+1,width)])
        pix = parallely.RandchoiceFlagFloat64(0.999, pix, 0, &flag)
        if flag {
          flag = false
          pix = 0.2*(array[Idx(i,j,width)] + array[Idx(i-1,j,width)] + array[Idx(i+1,j,width)] + array[Idx(i,j-1,width)] + array[Idx(i,j+1,width)])
          pix = parallely.RandchoiceFlagFloat64(0.9999, pix, 0, &flag)
        }
        result[Idx(i,j,width)] = pix
        result_exact[Idx(i,j,width)] = 0.2*(array[Idx(i,j,width)] + array[Idx(i-1,j,width)] + array[Idx(i+1,j,width)] + array[Idx(i,j-1,width)] + array[Idx(i,j+1,width)])
        overallflag = overallflag || flag
      }
    }
  }

  if overallflag {
    l2diff := 0.0
    l2a := 0.0
    l2b := 0.0
    for i := 0; i < height*width; i++ {
      diff := result[i] - result_exact[i]
      l2diff += diff*diff
      l2a += result[i]*result[i]
      l2b += result_exact[i]*result_exact[i]
    }
    fmt.Println(1,math.Sqrt(l2diff/(l2a*l2b)))
  } else {
    fmt.Println(0)
  }
}
