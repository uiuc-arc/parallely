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

  array := make([]int, height*width)
  result := make([]int, height*width)

  overallflag := false

  for i := 1; i < height-1; i++ {
    for j := 1; j < width-1; j++ {
      flag := false
      pix := array[Idx(i-1,j-1,width)] + 2*array[Idx(i-1,j,width)] + array[Idx(i-1,j+1,width)] - array[Idx(i+1,j-1,width)] - 2*array[Idx(i+1,j,width)] - array[Idx(i+1,j+1,width)]
      pix = parallely.RandchoiceFlag(0.99, pix, 0, &flag)
      if flag {
        flag = false
        pix := array[Idx(i-1,j-1,width)] + 2*array[Idx(i-1,j,width)] + array[Idx(i-1,j+1,width)] - array[Idx(i+1,j-1,width)] - 2*array[Idx(i+1,j,width)] - array[Idx(i+1,j+1,width)]
        pix = parallely.RandchoiceFlag(0.9999, pix, 0, &flag)
      }
      result[Idx(i,j,width)] = pix
      overallflag = overallflag || flag
    }
  }
  if overallflag {
    fmt.Println(1)
  } else {
    fmt.Println(0)
  }
}
