package main

import "fmt"
import "math/rand"
import ."dynfloats"

const optimize = true

func sor(band int, channelin, channelout chan float32, dcin, dcout chan float64) {
  var dmap [rows*cols+bandw*cols+1]float64
  var array [rows*cols]float32
  var result [bandw*cols]float32
  for iter:=0; iter<iterations; iter++ {
    RecvF32ArrAcc(array[:], dmap[:], rows*cols, 0, channelin, dcin, optimize)
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
          dmap[(rows+bandw)*cols] = dmap[GetIdx(i,j,cols)]+dmap[GetIdx(i-1,j,cols)]+dmap[GetIdx(i+1,j,cols)]+dmap[GetIdx(i,j-1,cols)]+dmap[GetIdx(i,j+1,cols)]
          result[GetIdx(i-bandStart,j,cols)] = sum*0.2
          dmap[rows*cols + GetIdx(i-bandStart,j,cols)] = dmap[(rows+bandw)*cols]*0.2
        }
        result[GetIdx(i-bandStart,cols-1,cols)] = array[GetIdx(i,cols-1,cols)]
        dmap[rows*cols + GetIdx(i-bandStart,cols-1,cols)] = dmap[GetIdx(i,cols-1,cols)]
      }
    }
    SendF32ArrAcc(result[:], dmap[rows*cols:], bandw*cols, 0, channelout, dcout, optimize)
  }
}

func printMaxDelta(dmap [(rows+bandw)*cols]float64, array32 [rows*cols]float32) {
  const boundary = 11
  maxDelta := 0.0
  //maxI := -1
  //maxJ := -1
  for i:=boundary; i<rows-boundary; i++ {
    for j:=boundary; j<cols-boundary; j++ {
      idx := GetIdx(i,j,cols)
      delta := dmap[idx]
      if maxDelta < delta {
        maxDelta = delta
        //maxI = i
        //maxJ = j
      }
    }
  }
  fmt.Println(maxDelta/*, maxI, maxJ*/)
}

func main() {
  randSource := rand.NewSource(seed)
  randGen := rand.New(randSource)
  var dmap [(rows+bandw)*cols]float64
  var array32 [rows*cols]float32
  var slice [bandw*cols]float32

  for i:=0; i<rows*cols; i++ {
    temp64 := MakeDynFloat64(randGen.Float64()*10.0-5.0)
    temp32 := DynFloat64To32(temp64)
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

  printMaxDelta(dmap, array32)

  for iter:=0; iter<iterations; iter++ {
    for band := 0; band < bands; band++ {
      SendF32ArrAcc(array32[:], dmap[:], rows*cols, 0, channels[band], dchannels[band], optimize)
    }
    for band := 0; band < bands; band++ {
      RecvF32ArrAcc(slice[:], dmap[rows*cols:], bandw*cols, 0, channels[band+bands], dchannels[band+bands], optimize)
      copy(array32[GetIdx(band*bandw,0,cols):GetIdx((band+1)*bandw,0,cols)], slice[:])
      copy(dmap[GetIdx(band*bandw,0,cols):GetIdx((band+1)*bandw,0,cols)], dmap[rows*cols:(rows+bandw)*cols])
    }
    //fmt.Println(array32[GetIdx(rows/2,cols/2,cols)], dmap[GetIdx(rows/2,cols/2,cols)])
    printMaxDelta(dmap, array32)
  }
}
