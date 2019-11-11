package main

import "fmt"
import "math/rand"
import "time"

func GetIdx(row, col, cols int) int {
  return row*cols + col
}

const Dim = 100
const Bands = 10
const BandW = Dim/Bands

func mm(band int, chin, chout chan []float32) {
  m1 := <- chin
  m2 := <- chin
  mr := make([]float32, BandW*Dim)
  for i:=band*BandW; i<(band+1)*BandW; i++ {
    for j:=0; j<Dim; j++ {
      sum := float32(0)
      for k:=0; k<Dim; k++ {
        sum = sum + m1[GetIdx(i,k,Dim)]*m2[GetIdx(k,j,Dim)]
      }
      mr[GetIdx(i-band*BandW,j,Dim)] = sum
    }
  }
  chout <- mr
}

func main() {
  randSource := rand.NewSource(time.Now().UnixNano())
  randGen := rand.New(randSource)
  var m164 [Dim*Dim]float64
  var m132 [Dim*Dim]float32
  var m264 [Dim*Dim]float64
  var m232 [Dim*Dim]float32
  var mr32 [Dim*Dim]float32

  for i:=0; i<Dim*Dim; i++ {
    m164[i] = randGen.Float64()
    m132[i] = float32(m164[i])
    m264[i] = randGen.Float64()
    m232[i] = float32(m264[i])
  }

  channels := make([]chan []float32, Bands*2)
  for i := range channels {
    channels[i] = make(chan []float32, 2)
  }

  for i:=0; i<Bands; i++ {
    go mm(i, channels[i], channels[i+Bands])
  }

  startTime := time.Now()

  for band := 0; band < Bands; band++ {
    m1Copy := make([]float32, Dim*Dim)
    copy(m1Copy, m132[:])
    channels[band] <- m1Copy
    m2Copy := make([]float32, Dim*Dim)
    copy(m2Copy, m232[:])
    channels[band] <- m2Copy
  }
  for band := 0; band < Bands; band++ {
    data := <- channels[band+Bands]
    copy(mr32[GetIdx(band*BandW,0,Dim):GetIdx((band+1)*BandW,0,Dim)], data)
  }

  elapsed := time.Since(startTime)
  fmt.Println(elapsed)
}
