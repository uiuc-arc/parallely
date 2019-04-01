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

func sor(height, width, tsHeight, teHeight int, omega float64, channel chan []float64) {
  array := <- channel
  result := make([]float64, (teHeight-tsHeight)*width)
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
      result[Idx(i-tsHeight,j,width)] = omega/4.0*(up+down+left+right) + (1.0-omega)*center
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

  channels := make([]chan []float64, numThreads)
  for i := range channels {
    channels[i] = make(chan []float64, 1)
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
  for iter := 0; iter < iterations; iter++ {
    for thr := 0; thr < numThreads; thr++ {
      arrayCopy := make([]float64, height*width)
      copy(arrayCopy,array)
      channels[thr] <- arrayCopy
      //channels[thr] <- array
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
      copy(array[Idx(tsHeight,0,width):Idx(teHeight,0,width)], tile)
    }
  }

  f, _ := os.Create("/tmp/original.dat")
  defer f.Close()
  encoder := gob.NewEncoder(f)
  encoder.Encode(array)

  elapsed := time.Since(startTime)
  fmt.Println(elapsed)
}
