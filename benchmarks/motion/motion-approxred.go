package main

import (
  "os"
  "strconv"
  "math/rand"
  "time"
  "fmt"
)

const BlockSize = 64
const SampleRate = 0.01

func calcSsd(blockChan <-chan []uint8, ssdChan chan<- int) {
  myBlock := <- blockChan
  compBlock := <- blockChan
  if rand.Float32() < SampleRate {
    ssd := 0
    for i := 0; i < BlockSize; i++ {
      diff := int(myBlock[i])-int(compBlock[i])
      ssd += diff*diff
    }
    ssdChan <- ssd
  } else {
    ssdChan <- -1
  }
}

func main() {
  if len(os.Args)!=2 {
    fmt.Println("Usage:\nmotion.go numThreads")
    os.Exit(0)
  }
  numThreads, _ := strconv.Atoi(os.Args[1])

  randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

  blocks := make([][]uint8, numThreads+1)
  for i := range blocks {
    blocks[i] = make([]uint8, BlockSize)
    for j := range blocks[i] {
      blocks[i][j] = uint8(randGen.Intn(256))
    }
  }

  startTime := time.Now()

  blockChans := make([]chan []uint8, numThreads)
  for i := range blockChans {
    blockChans[i] = make(chan []uint8, 2)
  }
  ssdChans := make([]chan int, numThreads)
  for i := range ssdChans {
    ssdChans[i] = make(chan int, 1)
  }

  for i := 0; i < numThreads; i++ {
    go calcSsd(blockChans[i],ssdChans[i])
  }

  for i := 0; i < numThreads; i++ {
    myBlockCopy := make([]uint8, BlockSize)
    copy(myBlockCopy,blocks[i])
    blockChans[i] <- myBlockCopy
    refBlockCopy := make([]uint8, BlockSize)
    copy(refBlockCopy,blocks[numThreads])
    blockChans[i] <- refBlockCopy
  }
  minSsd := 2147483647
  minBlock := -1
  skippedBlocks := 0
  for i := 0; i < numThreads; i++ {
    ssd := <- ssdChans[i]
    if ssd>=0 {
      if ssd<minSsd {
        minSsd = ssd
        minBlock = i
      }
    } else {
      skippedBlocks++
    }
  }

  elapsed := time.Since(startTime)

  fmt.Println(elapsed,minBlock,minSsd,skippedBlocks)
}
