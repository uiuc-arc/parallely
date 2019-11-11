package main

import "fmt"
import "math/rand"
import "time"
import ."dynfloats"

func GetIdx(row, col, cols int) int {
  return row*cols + col
}

const rows = 100
const cols = 100
const bands = 10
const bandw = 10

func sor(band int, channelin, channelout chan []DynFloat32) {
  Const2 := MakeDynFloat32(2)
  ConstM1 := MakeDynFloat32(-1)
  ConstM2 := MakeDynFloat32(-2)
  array := <- channelin
  result := make([]DynFloat32, bandw*cols)
  bandStart := band*bandw
  for i := bandStart; i < bandStart+bandw; i++ {
    if i==0 || i==cols-1 {
      for j := 0; j < cols; j++ {
        result[GetIdx(i-bandStart,j,cols)] = array[GetIdx(i,j,cols)]
      }
    } else {
      result[GetIdx(i-bandStart,0,cols)] = array[GetIdx(i,0,cols)]
      for j := 1; j < cols-1; j++ {
        sum := MakeDynFloat32(0)
        sum = AddDynFloat32(sum,array[GetIdx(i-1,j-1,cols)])
        sum = AddDynFloat32(sum,MulDynFloat32(array[GetIdx(i-1,j  ,cols)],Const2 ))
        sum = AddDynFloat32(sum,array[GetIdx(i-1,j+1,cols)])
        sum = AddDynFloat32(sum,MulDynFloat32(array[GetIdx(i+1,j-1,cols)],ConstM1))
        sum = AddDynFloat32(sum,MulDynFloat32(array[GetIdx(i+1,j  ,cols)],ConstM2))
        sum = AddDynFloat32(sum,MulDynFloat32(array[GetIdx(i+1,j+1,cols)],ConstM1))
        result[GetIdx(i-bandStart,j,cols)] = sum
      }
      result[GetIdx(i-bandStart,cols-1,cols)] = array[GetIdx(i,cols-1,cols)]
    }
  }
  channelout <- result
}

func main() {
  randSource := rand.NewSource(time.Now().UnixNano())
  randGen := rand.New(randSource)
  var array64 [rows*cols]DynFloat64
  var array32 [rows*cols]DynFloat32

  for i:=0; i<rows*cols; i++ {
    array64[i] = MakeDynFloat64(randGen.Float64())
    array32[i] = DynFloat64To32(array64[i])
  }

  channels := make([]chan []DynFloat32, bands*2)
  for i := range channels {
    channels[i] = make(chan []DynFloat32, 1)
  }

  for i:=0; i<bands; i++ {
    go sor(i, channels[i], channels[i+bands])
  }

  startTime := time.Now()

  for band := 0; band < bands; band++ {
    array32Copy := make([]DynFloat32, rows*cols)
    copy(array32Copy, array32[:])
    channels[band] <- array32Copy
  }
  for band := 0; band < bands; band++ {
    data := <- channels[band+bands]
    copy(array32[GetIdx(band*bandw,0,cols):GetIdx((band+1)*bandw,0,cols)], data)
  }

  elapsed := time.Since(startTime)
  fmt.Println(elapsed)

  badCount := 0
  for i:=0; i<rows*cols; i++ {
    if array32[i].Delta > 1e-5 {
      badCount += 1
    }
  }
  fmt.Println(badCount)
}
