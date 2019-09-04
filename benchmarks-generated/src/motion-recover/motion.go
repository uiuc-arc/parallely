package main

import (
  "os"
  "strconv"
  "math/rand"
  "time"
  "fmt"
  "parallely"
)

const BlockSize = 1600

func main() {
  rand.Seed(time.Now().UTC().UnixNano())
  
  if len(os.Args)!=2 {
    fmt.Println("Usage:\nmotion.go numThreads")
    os.Exit(0)
  }
  numThreads, _ := strconv.Atoi(os.Args[1])

  //randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

  blocks := make([][]uint8, numThreads+1)
  for i := range blocks {
    blocks[i] = make([]uint8, BlockSize)
    /*for j := range blocks[i] {
      blocks[i][j] = uint8(randGen.Intn(256))
    }*/
  }

  overallflag := false

  minSsd := 2147483647
  for i := 0; i < numThreads; i++ {
    ssd := 0
    for j := 0; j < BlockSize; j++ {
      var diff int
      var flag bool
      flag = false
      diff = int(blocks[0][j])-int(blocks[i+1][j])
      diff = parallely.RandchoiceFlag(0.999, diff, 0, &flag)
      if flag {
        flag = false
        diff = int(blocks[0][j])-int(blocks[i+1][j])
        diff = parallely.RandchoiceFlag(0.9999, diff, 0, &flag)
      }
      overallflag = overallflag || flag
      ssd = ssd + diff*diff
    }
    if ssd < minSsd {
      minSsd = ssd
    }
  }
  if overallflag {
    fmt.Println(1)
  } else {
    fmt.Println(0)
  }
}
