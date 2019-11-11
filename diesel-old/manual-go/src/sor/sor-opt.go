package main

import "fmt"
import "math/rand"
import "time"
import ."dynfloats"

func sor(band int, channelin, channelout chan float32, dcin, dcout chan float64) {
  var dmap = map[int] float64{}
  var array [rows*cols]float32
  var result [bandw*cols]float32
  for iter:=0; iter<iterations; iter++ {
    for i := range array {
      array[i] = <- channelin
    }
    dtemp := <- dcin
    for i := range array {
      dmap[i] = dtemp
    }
    bandStart := band*bandw
    for i := bandStart; i < bandStart+bandw; i++ {
      if i==0 || i==cols-1 {
        for j := 0; j < cols; j++ {
          result[GetIdx(i-bandStart,j,cols)] = array[GetIdx(i,j,cols)]
          dmap[rows*cols + GetIdx(i-bandStart,j,cols)] = dmap[GetIdx(i,j,cols)]
        }
      } else {
        result[GetIdx(i-bandStart,0,cols)] = array[GetIdx(i,0,cols)]
        dmap[rows*cols + GetIdx(i-bandStart,0,cols)] = dmap[GetIdx(i,0,cols)]
        for j := 1; j < cols-1; j++ {
          sum := array[GetIdx(i,j,cols)]+array[GetIdx(i-1,j,cols)]+array[GetIdx(i+1,j,cols)]+array[GetIdx(i,j-1,cols)]+array[GetIdx(i,j+1,cols)]
          dmap[(rows+bandw)*cols + 1] = dmap[GetIdx(i,j,cols)]+dmap[GetIdx(i-1,j,cols)]+dmap[GetIdx(i+1,j,cols)]+dmap[GetIdx(i,j-1,cols)]+dmap[GetIdx(i,j+1,cols)]
          result[GetIdx(i-bandStart,j,cols)] = sum*0.2
          dmap[rows*cols + GetIdx(i-bandStart,j,cols)] = dmap[(rows+bandw)*cols + 1]*0.2
        }
        result[GetIdx(i-bandStart,cols-1,cols)] = array[GetIdx(i,cols-1,cols)]
        dmap[rows*cols + GetIdx(i-bandStart,cols-1,cols)] = dmap[GetIdx(i,cols-1,cols)]
      }
    }
    dtemp = 0.0
    for i := range result {
      channelout <- result[i]
      if dmap[rows*cols + i] > dtemp {
        dtemp = dmap[rows*cols + i]
      }
    }
    dcout <- dtemp
  }
}

func main() {
  randSource := rand.NewSource(time.Now().UnixNano())
  randGen := rand.New(randSource)
  var dmap = map[int] float64{}
  //var array64 [rows*cols]float64
  var array32 [rows*cols]float32
  var slice [bandw*cols]float32

  for i:=0; i<rows*cols; i++ {
    temp64 := MakeDynFloat64(randGen.Float64())
    temp32 := DynFloat64To32(temp64)
    //array64[i] = randGen.Float64()
    array32[i] = temp32.Num
    dmap[i] = temp32.Delta
  }

  var channels [bands*2]chan float32
  var dchannels [bands*2]chan float64
  for i := range channels {
    channels[i] = make(chan float32, rows*cols)
    dchannels[i] = make(chan float64, rows*cols)
  }

  for i:=0; i<bands; i++ {
    go sor(i, channels[i], channels[i+bands], dchannels[i], dchannels[i+bands])
  }

  startTime := time.Now()

  for iter:=0; iter<iterations; iter++ {
    for band := 0; band < bands; band++ {
      dtemp := 0.0
      for i := range array32 {
        channels[band] <- array32[i]
        if dmap[i] > dtemp {
          dtemp = dmap[i]
        }
      }
      dchannels[band] <- dtemp
    }
    for band := 0; band < bands; band++ {
      for i := range slice {
        slice[i] = <- channels[band+bands]
      }
      dtemp := <- dchannels[band+bands]
      for i := range slice {
        dmap[rows*cols + i] = dtemp
      }
      copy(array32[GetIdx(band*bandw,0,cols):GetIdx((band+1)*bandw,0,cols)], slice[:])
      for i := range slice {
        dmap[GetIdx(band*bandw,0,cols) + i] = dmap[rows*cols + i]
      }
    }
  }

  badCount := 0
  minDelta := dmap[0]
  maxDelta := dmap[0]
  for i:=0; i<rows*cols; i++ {
    delta := dmap[i]
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
