package main

import (
  "os"
  "strconv"
  "fmt"
  "parallely"
)

func Idx(i, j, width int) int {
  return i*width+j
}

func main() {
  height, _ := strconv.Atoi(os.Args[1])
  width, _ := strconv.Atoi(os.Args[2])
  iterations, _ := strconv.Atoi(os.Args[3])

  array := make([]float64, height*width)
  result := make([]float64, height*width)

  overallflag := false

  for iter := 0; iter < iterations; iter++ {
    for i := 1; i < height-1; i++ {
      for j := 1; j < width-1; j++ {
        flag := false
        pix := 0.2*(array[Idx(i,j,width)] + array[Idx(i-1,j,width)] + array[Idx(i+1,j,width)] + array[Idx(i,j-1,width)] + array[Idx(i,j+1,width)])
        pix = parallely.RandchoiceFlagFloat64(0.999, pix, 0, &flag)
        if flag {
          flag = false
          pix := 0.2*(array[Idx(i,j,width)] + array[Idx(i-1,j,width)] + array[Idx(i+1,j,width)] + array[Idx(i,j-1,width)] + array[Idx(i,j+1,width)])
          pix = parallely.RandchoiceFlagFloat64(0.9999, pix, 0, &flag)
        }
        result[Idx(i,j,width)] = pix
        overallflag = overallflag || flag
      }
    }
  }
  if overallflag {
    fmt.Println(1)
  } else {
    fmt.Println(0)
  }
}
