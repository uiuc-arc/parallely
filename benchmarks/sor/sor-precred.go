package main

import (
  "os"
  "strconv"
  "math/rand"
  "time"
  "fmt"
  "encoding/gob"
)

func Idx(i, j, width int) int {
  return i*width+j
}

func sliceTof32(slice []float64) []float32 {
  length := len(slice)
  result := make([]float32, length)
  for i := 0; i < length; i++ {
    result[i] = float32(slice[i])
  }
  return result
}

func sliceTof64(slice []float32) []float64 {
  length := len(slice)
  result := make([]float64, length)
  for i := 0; i < length; i++ {
    result[i] = float64(slice[i])
  }
  return result
}

func sor(height, width, tsHeight, teHeight int, omega float64, channel chan []float32) {
  array := <- channel
  result := make([]float32, (teHeight-tsHeight)*width)
  for i := tsHeight; i < teHeight; i++ {
    if i==0 || i==height-1 {
      continue
    }
    for j := 1; j < width-1; j++ {
      up := array[Idx(i-1,j,width)]
      down := array[Idx(i+1,j,width)]
      left := array[Idx(i,j-1,width)]
      right := array[Idx(i,j+1,width)]
      center := array[Idx(i,j,width)]
      result[Idx(i-tsHeight,j,width)] = float32(omega)/4.0*(up+down+left+right) + (1.0-float32(omega))*center
    }
  }
  channel <- result
}

func main() {
  if len(os.Args)!=7 {
    fmt.Println("Usage:\nsor-precred.go height width iterations omega numThreads seed")
    os.Exit(0)
  }
  height, _ := strconv.Atoi(os.Args[1])
  width, _ := strconv.Atoi(os.Args[2])
  iterations, _ := strconv.Atoi(os.Args[3])
  omega, _ := strconv.ParseFloat(os.Args[4],64)
  numThreads, _ := strconv.Atoi(os.Args[5])
  seed, _ := strconv.Atoi(os.Args[6])

  randGen := rand.New(rand.NewSource(int64(seed)))

  array := make([]float64, height*width)
  for i := 0; i < height*width; i++ {
    array[i] = randGen.Float64()
  }

  channels := make([]chan []float32, numThreads)
  for i := range channels {
    channels[i] = make(chan []float32, 1)
  }

  startTime := time.Now()

  tHeight := height/numThreads
  for thr := 0; thr < numThreads; thr++ {
    tsHeight := tHeight*thr
    var teHeight int
    if thr==numThreads-1 {
      teHeight = height-1
    } else {
      teHeight = tHeight*(thr+1)
    }
    go sor(height, width, tsHeight, teHeight, omega, channels[thr])
  }
  array32 := sliceTof32(array)
  for iter := 0; iter < iterations; iter++ {
    for thr := 0; thr < numThreads; thr++ {
      array32Copy := make([]float32, height*width)
      copy(array32Copy,array32)
      channels[thr] <- array32Copy
      //channels[thr] <- array32
    }
    for thr := 0; thr < numThreads; thr++ {
      tsHeight := tHeight*thr
      var teHeight int
      if thr==numThreads-1 {
        teHeight = height-1
      } else {
        teHeight = tHeight*(thr+1)
      }
      tile := <- channels[thr]
      copy(array32[Idx(tsHeight,0,width):Idx(teHeight,0,width)], tile)
    }
  }
  array = sliceTof64(array32)

  f, _ := os.Create("/tmp/precred.dat")
  defer f.Close()
  encoder := gob.NewEncoder(f)
  encoder.Encode(array)

  elapsed := time.Since(startTime)
  fmt.Println(elapsed)
}
