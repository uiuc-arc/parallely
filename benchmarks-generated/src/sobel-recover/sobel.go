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

  array := make([]int, height*width)
  for i := 0; i < height*width; i++ {
    array[i] = i%23
  }
  result := make([]int, height*width)
  result_exact := make([]int, height*width)

  overallflag := false

  for i := 1; i < height-1; i++ {
    for j := 1; j < width-1; j++ {
      flag := false
      pix := array[Idx(i-1,j-1,width)] + 2*array[Idx(i-1,j,width)] + array[Idx(i-1,j+1,width)] - array[Idx(i+1,j-1,width)] - 2*array[Idx(i+1,j,width)] - array[Idx(i+1,j+1,width)]
      pix_exact := pix
      pix = parallely.RandchoiceFlag(0.99, pix, 0, &flag)
      if flag {
        flag = false
        pix := array[Idx(i-1,j-1,width)] + 2*array[Idx(i-1,j,width)] + array[Idx(i-1,j+1,width)] - array[Idx(i+1,j-1,width)] - 2*array[Idx(i+1,j,width)] - array[Idx(i+1,j+1,width)]
        pix = parallely.RandchoiceFlag(0.9999, pix, 0, &flag)
      }
      result[Idx(i,j,width)] = pix
      result_exact[Idx(i,j,width)] = pix_exact
      overallflag = overallflag || flag
    }
  }

  if overallflag {
    l2diff := 0.0
    l2a := 0.0
    l2b := 0.0
    for i := 0; i < height*width; i++ {
      diff := result[i] - result_exact[i]
      l2diff += float64(diff*diff)
      l2a += float64(result[i]*result[i])
      l2b += float64(result_exact[i]*result_exact[i])
    }
    fmt.Println(1,math.Sqrt(l2diff/(l2a*l2b)))
  } else {
    fmt.Println(0)
  }
}
