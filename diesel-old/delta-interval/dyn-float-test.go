package main

import "fmt"
import "math/rand"
//import "time"

func GetIdx(row, col, cols int) int {
  return row*cols + col
}

const rows = 10
const cols = 10
const iterations = 100

func main() {
  randSource := rand.NewSource(123456789)//time.Now().UnixNano())
  randGen := rand.New(randSource)
  var array64 [rows*cols]DynFloat64
  var array32 [rows*cols]DynFloat32
  var array32_temp [rows*cols]DynFloat32
  for i:=0; i<rows*cols; i++ {
    array64[i] = MakeDynFloat64(randGen.Float64())
    array32[i] = DynFloat64To32(array64[i])
    array32_temp[i] = array32[i]
  }
  const02 := MakeDynFloat32(0.2)
  for iter:=0; iter<iterations; iter++ {
    //1 iteration of SOR
    for i:=1; i<rows-1; i++ {
      for j:=1; j<cols-1; j++ {
        sum1 := AddDynFloat32(array32[GetIdx(i-1,j,cols)], array32[GetIdx(i+1,j,cols)])
        sum2 := AddDynFloat32(array32[GetIdx(i,j-1,cols)], array32[GetIdx(i,j+1,cols)])
        sum1 = AddDynFloat32(sum1, sum2)
        sum1 = AddDynFloat32(sum1, array32[GetIdx(i,j,cols)])
        array32_temp[GetIdx(i,j,cols)] = MulDynFloat32(sum1, const02)
      }
    }
    /*
    maxInterval := 0.0
    maxIntervalIndex := -1
    for i:=0; i<rows*cols; i++ {
      array32[i] = array32_temp[i]
      low := array32[i].interval.low
      hig := array32[i].interval.hig
      num := float64(array32[i].num)
      if low > num || hig < num {
        fmt.Println("Error: out of interval!")
      }
      if hig-low >= maxInterval {
        maxInterval = hig-low
        maxIntervalIndex = i
      }
    }
    fmt.Println(iter+1, array32[maxIntervalIndex].interval.low, array32[maxIntervalIndex].num, array32[maxIntervalIndex].interval.hig)
    */
    
    maxInterval := 0.0
    maxIntervalIndex := -1
    for i:=0; i<rows*cols; i++ {
      array32[i] = array32_temp[i]
      if array32[i].delta >= maxInterval {
        maxInterval = array32[i].delta
        maxIntervalIndex = i
      }
    }
    fmt.Println(iter+1, array32[maxIntervalIndex].num, array32[maxIntervalIndex].delta)
    
  }
}
