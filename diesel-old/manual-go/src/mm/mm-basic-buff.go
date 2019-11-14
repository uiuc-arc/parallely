package main

import "fmt"
import "math/rand"
import "time"

func mm(band int, chin, chout chan float32) {
  var m1 [Dim*Dim]float32
  var m2 [Dim*Dim]float32
  var mr [BandW*Dim]float32
  for i := range m1 {
    m1[i] = <- chin
  }
  for i := range m2 {
    m2[i] = <- chin
  }
  for i:=band*BandW; i<(band+1)*BandW; i++ {
    for j:=0; j<Dim; j++ {
      sum := float32(0.0)
      for k:=0; k<Dim; k++ {
        sum += m1[GetIdx(i,k,Dim)] * m2[GetIdx(k,j,Dim)]
      }
      mr[GetIdx(i-band*BandW,j,Dim)] = sum
    }
  }
  for i := range mr {
    chout <- mr[i]
  }
}

func main() {
  randSource := rand.NewSource(time.Now().UnixNano())
  randGen := rand.New(randSource)
  var m164 [Dim*Dim]float64
  var m132 [Dim*Dim]float32
  var m264 [Dim*Dim]float64
  var m232 [Dim*Dim]float32
  var mr32 [Dim*Dim]float32
  var slice [BandW*Dim]float32

  for i:=0; i<Dim*Dim; i++ {
    m164[i] = randGen.Float64()
    m132[i] = float32(m164[i])
    m264[i] = randGen.Float64()
    m232[i] = float32(m264[i])
  }

  var channels [Bands*2]chan float32
  for i := range channels {
    channels[i] = make(chan float32, 2*Dim*Dim)
  }

  for i:=0; i<Bands; i++ {
    go mm(i, channels[i], channels[i+Bands])
  }

  startTime := time.Now()

  for band := 0; band < Bands; band++ {
    for i := range m132 {
      channels[band] <- m132[i]
    }
    for i := range m232 {
      channels[band] <- m232[i]
    }
  }
  for band := 0; band < Bands; band++ {
    for i := range slice {
      slice[i] = <- channels[band+Bands]
    }
    copy(mr32[GetIdx(band*BandW,0,Dim):GetIdx((band+1)*BandW,0,Dim)], slice[:])
  }

  elapsed := time.Since(startTime)
  fmt.Println(elapsed.Nanoseconds())
}
