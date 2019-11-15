package main

import "fmt"
import "math/rand"
import "time"

func sor(band int, channelin, channelout chan float32) {
  var array [rows*cols]float32
  var result [bandw*cols]float32
  for iter:=0; iter<iterations; iter++ {
    for i := range array {
      array[i] = <- channelin
    }
    bandStart := band*bandw
    for i := bandStart; i < bandStart+bandw; i++ {
      if i==0 || i==cols-1 {
        for j := 0; j < cols; j++ {
          result[GetIdx(i-bandStart,j,cols)] = array[GetIdx(i,j,cols)]
        }
      } else {
        result[GetIdx(i-bandStart,0,cols)] = array[GetIdx(i,0,cols)]
        for j := 1; j < cols-1; j++ {
          sum := array[GetIdx(i,j,cols)]+array[GetIdx(i-1,j,cols)]+array[GetIdx(i+1,j,cols)]+array[GetIdx(i,j-1,cols)]+array[GetIdx(i,j+1,cols)]
          result[GetIdx(i-bandStart,j,cols)] = sum*0.2
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
  randSource := rand.NewSource(seed)
  randGen := rand.New(randSource)
  var array64 [rows*cols]float64
  var array32 [rows*cols]float32
  var slice [bandw*cols]float32

  for i:=0; i<rows*cols; i++ {
    array64[i] = randGen.Float64()
    array32[i] = float32(array64[i])
  }

  var channels [bands*2]chan float32
  for i := range channels {
    channels[i] = make(chan float32, rows*cols)
  }

	startTime := time.Now()
	
  for i:=0; i<bands; i++ {
    go sor(i, channels[i], channels[i+bands])
  }

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

	elapsed := time.Since(startTime)
  fmt.Println(elapsed.Nanoseconds())
}
