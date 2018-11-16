package main

import (
  "os"
  "strconv"
  "math/rand"
  "time"
  "fmt"
  "math"
)

func Idx(i, j, width int) int {
  return i*width+j
}

func main() {
  if len(os.Args)!=5 {
    fmt.Println("Usage:\nsor-sequential.go height width iterations omega")
    os.Exit(0)
  }
  height, _ := strconv.Atoi(os.Args[1])
  width, _ := strconv.Atoi(os.Args[2])
  iterations, _ := strconv.Atoi(os.Args[3])
  omega, _ := strconv.ParseFloat(os.Args[4],64)

  randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

  array := make([]float32, height*width)
  array_err := make([]float64, height*width)
  conversionError := math.Pow(2.0,-24.0)
  for i := 0; i < height*width; i++ {
    array[i] = randGen.Float32()
    array_err[i] = conversionError
  }

  startTime := time.Now()

  result := make([]float32, height*width)
  result_err := make([]float64, height*width)
  for iter := 0; iter < iterations; iter++ {
    for i := 1; i < height-1; i++ {
      for j := 1; j < width-1; j++ {
        up := array[Idx(i-1,j,width)]
        up_err := array_err[Idx(i-1,j,width)]
        down := array[Idx(i+1,j,width)]
        down_err := array_err[Idx(i+1,j,width)]
        left := array[Idx(i,j-1,width)]
        left_err := array_err[Idx(i,j-1,width)]
        right := array[Idx(i,j+1,width)]
        right_err := array_err[Idx(i,j+1,width)]
        center := array[Idx(i,j,width)]
        center_err := array_err[Idx(i,j,width)]
        result[Idx(i,j,width)] = float32(omega)/4.0*(up+down+left+right) + (1.0-float32(omega))*center
        result_err[Idx(i,j,width)] = omega/4.0*(up_err+down_err+left_err+right_err) + (1.0-omega)*center_err
        //result_err[Idx(i,j,width)] = math.Max(math.Max(math.Max(up_err,down_err),math.Max(left_err,right_err)),center_err)
      }
    }
    array = result
    array_err = result_err
  }

  maxError := 0.0
  for i := 0; i < height*width; i++ {
    if array_err[i] > maxError {
      maxError = array_err[i]
    }
  }
  fmt.Println("Max additive error",maxError)
  //fmt.Println("Max multiplicative error",maxError)

  elapsed := time.Since(startTime)
  fmt.Println(elapsed)
}
