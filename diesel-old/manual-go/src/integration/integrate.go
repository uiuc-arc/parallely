package main

import "fmt"
import "time"
import ."dynfloats"

const Max = float32(10)
const Min = float32(0)
const Delta = float32(0.00390625)
const Threads = 10
const IntervalPerThread = (Max-Min)/float32(Threads)
const DivisionsPerThread = int(IntervalPerThread/Delta)

func calc(x DynFloat32) DynFloat32 {
  return MulDynFloat32(x,x)
}

func integrate(input chan int, output chan DynFloat32) {
  sum := MakeDynFloat32(0)
  dynDelta := MakeDynFloat32(Delta)
  threadId := <- input
  x := MakeDynFloat32(Min + IntervalPerThread*float32(threadId))
  for i:=0; i<DivisionsPerThread; i++ {
    sum = AddDynFloat32(sum,MulDynFloat32(calc(x),dynDelta))
    x = AddDynFloat32(x,dynDelta)
  }
  output <- sum
}

func main() {
  inputChans := make([]chan int, Threads)
  for i := range inputChans {
    inputChans[i] = make(chan int, 1)
  }
  outputChans := make([]chan DynFloat32, Threads)
  for i := range inputChans {
    outputChans[i] = make(chan DynFloat32, 1)
  }

  for i:=0; i<Threads; i++ {
    go integrate(inputChans[i], outputChans[i])
  }

  startTime := time.Now()

  for i:=0; i<Threads; i++ {
    inputChans[i] <- i
  }
  sum := MakeDynFloat32(0)
  for i:=0; i<Threads; i++ {
    temp := <- outputChans[i]
    sum = AddDynFloat32(sum,temp)
  }

  elapsed := time.Since(startTime)
  fmt.Println(elapsed)

  fmt.Println(sum.Delta)
}
