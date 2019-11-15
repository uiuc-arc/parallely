package main

import "fmt"
import "math/rand"
import "time"
import ."dynfloats"

func sor(band int, channelin, channelout chan float32, dcin, dcout chan float64) {
  var array [rows*cols]DynFloat32
  var result [bandw*cols]DynFloat32
  for iter:=0; iter<iterations; iter++ {
    for i := range array {
      array[i].Num = <- channelin
    }
    dtemp := <- dcin
    for i := range array {
      array[i].Delta = dtemp
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
          sum := AddDynFloat32(array[GetIdx(i,j,cols)], array[GetIdx(i-1,j,cols)])
          sum = AddDynFloat32(sum, array[GetIdx(i+1,j,cols)])
          sum = AddDynFloat32(sum, array[GetIdx(i,j-1,cols)])
          sum = AddDynFloat32(sum, array[GetIdx(i,j+1,cols)])
          result[GetIdx(i-bandStart,j,cols)] = MulDynFloat32(sum, const02)
        }
        result[GetIdx(i-bandStart,cols-1,cols)] = array[GetIdx(i,cols-1,cols)]
      }
    }
    dtemp = 0.0
    for i := range result {
      channelout <- result[i].Num
      if result[i].Delta > dtemp {
        dtemp = result[i].Delta
      }
    }
    dcout <- dtemp
  }
}

func main() {
  randSource := rand.NewSource(seed)
  randGen := rand.New(randSource)
  var array64 [rows*cols]DynFloat64
  var array32 [rows*cols]DynFloat32
  var slice [bandw*cols]DynFloat32

  for i:=0; i<rows*cols; i++ {
    array64[i] = MakeDynFloat64(randGen.Float64())
    array32[i] = DynFloat64To32(array64[i])
  }

  var channels [bands*2]chan float32
  var dchannels [bands*2]chan float64
  for i := range channels {
    channels[i] = make(chan float32)//, rows*cols)
    dchannels[i] = make(chan float64)//, 1)
  }

	startTime := time.Now()

  for i:=0; i<bands; i++ {
    go sor(i, channels[i], channels[i+bands], dchannels[i], dchannels[i+bands])
  }

  for iter:=0; iter<iterations; iter++ {
    for band := 0; band < bands; band++ {
      dtemp := 0.0
      for i := range array32 {
        channels[band] <- array32[i].Num
        if array32[i].Delta > dtemp {
          dtemp = array32[i].Delta
        }
      }
      dchannels[band] <- dtemp
    }
    for band := 0; band < bands; band++ {
      for i := range slice {
        slice[i].Num = <- channels[band+bands]
      }
      dtemp := <- dchannels[band+bands]
      for i := range slice {
        slice[i].Delta = dtemp
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
