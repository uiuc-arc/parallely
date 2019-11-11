package main

import "fmt"
import "math/rand"
import "time"
import ."dynfloats"

func sor(band int, channelin, channelout chan DynFloat32) {
  var array [rows*cols]DynFloat32
  var result [bandw*cols]DynFloat32
  for iter:=0; iter<iterations; iter++ {
    for i := range array {
      array[i] = <- channelin
    }
    bandStart := band*bandw
    const02 := MakeDynFloat32(0.2)
    for i := bandStart; i < bandStart+bandw; i++ {
      if i==0 || i==cols-1 {
        for j := 0; j < cols; j++ {
          result[GetIdx(i-bandStart,j,cols)] = array[GetIdx(i,j,cols)]
        }
      } else {
        result[GetIdx(i-bandStart,0,cols)] = array[GetIdx(i,0,cols)]
        for j := 1; j < cols-1; j++ {
          sum1 := AddDynFloat32(array[GetIdx(i-1,j,cols)], array[GetIdx(i+1,j,cols)])
          sum2 := AddDynFloat32(array[GetIdx(i,j-1,cols)], array[GetIdx(i,j+1,cols)])
          sum1 = AddDynFloat32(sum1, sum2)
          sum1 = AddDynFloat32(sum1, array[GetIdx(i,j,cols)])
          result[GetIdx(i-bandStart,j,cols)] = MulDynFloat32(sum1, const02)
        }
        result[GetIdx(i-bandStart,cols-1,cols)] = array[GetIdx(i,cols-1,cols)]
      }
    }
    for i := range result {
      channelout <- result[i]
    }
  }
}

func main() {
  randSource := rand.NewSource(time.Now().UnixNano())
  randGen := rand.New(randSource)
  var array64 [rows*cols]DynFloat64
  var array32 [rows*cols]DynFloat32
  var slice [bandw*cols]DynFloat32

  for i:=0; i<rows*cols; i++ {
    array64[i] = MakeDynFloat64(randGen.Float64())
    array32[i] = DynFloat64To32(array64[i])
  }

  var channels [bands*2]chan DynFloat32
  for i := range channels {
    channels[i] = make(chan DynFloat32, rows*cols)
  }

  for i:=0; i<bands; i++ {
    go sor(i, channels[i], channels[i+bands])
  }

  startTime := time.Now()

  for iter:=0; iter<iterations; iter++ {
    for band := 0; band < bands; band++ {
      for i := range array32 {
        channels[band] <- array32[i]
      }
    }
    for band := 0; band < bands; band++ {
      for i := range slice {
        slice[i] = <- channels[band+bands]
      }
      copy(array32[GetIdx(band*bandw,0,cols):GetIdx((band+1)*bandw,0,cols)], slice[:])
    }
  }

  badCount := 0
  minDelta := array32[0].Delta
  maxDelta := array32[0].Delta
  for i:=0; i<rows*cols; i++ {
    delta := array32[i].Delta
    if delta > 1e-4 {
      badCount += 1
    }
    if maxDelta < delta {
      maxDelta = delta
    }
    if minDelta > delta {
      minDelta = delta
    }
  }

  elapsed := time.Since(startTime)
  fmt.Println(elapsed.Nanoseconds(), badCount, minDelta, maxDelta)
}
