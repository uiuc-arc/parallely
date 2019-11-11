package main

import "fmt"
import "math/rand"
import "time"
import ."dynfloats"

func GetIdx(row, col, cols int) int {
  return row*cols + col
}

const Dim = 100
const Bands = 10
const BandW = Dim/Bands

func mm(band int, chin, chout chan []DynFloat32) {
  m1 := <- chin
  m2 := <- chin
  mr := make([]DynFloat32, BandW*Dim)
  for i:=band*BandW; i<(band+1)*BandW; i++ {
    for j:=0; j<Dim; j++ {
      sum := MakeDynFloat32(0)
      for k:=0; k<Dim; k++ {
        sum = AddDynFloat32(sum,MulDynFloat32(m1[GetIdx(i,k,Dim)],m2[GetIdx(k,j,Dim)]))
      }
      mr[GetIdx(i-band*BandW,j,Dim)] = sum
    }
  }
  chout <- mr
}

func main() {
  randSource := rand.NewSource(time.Now().UnixNano())
  randGen := rand.New(randSource)
  var m164 [Dim*Dim]DynFloat64
  var m132 [Dim*Dim]DynFloat32
  var m264 [Dim*Dim]DynFloat64
  var m232 [Dim*Dim]DynFloat32
  var mr32 [Dim*Dim]DynFloat32

  for i:=0; i<Dim*Dim; i++ {
    m164[i] = MakeDynFloat64(randGen.Float64())
    m132[i] = DynFloat64To32(m164[i])
    m264[i] = MakeDynFloat64(randGen.Float64())
    m232[i] = DynFloat64To32(m264[i])
  }

  channels := make([]chan []DynFloat32, Bands*2)
  for i := range channels {
    channels[i] = make(chan []DynFloat32, 2)
  }

  for i:=0; i<Bands; i++ {
    go mm(i, channels[i], channels[i+Bands])
  }

  startTime := time.Now()

  for band := 0; band < Bands; band++ {
    m1Copy := make([]DynFloat32, Dim*Dim)
    copy(m1Copy, m132[:])
    channels[band] <- m1Copy
    m2Copy := make([]DynFloat32, Dim*Dim)
    copy(m2Copy, m232[:])
    channels[band] <- m2Copy
  }
  for band := 0; band < Bands; band++ {
    data := <- channels[band+Bands]
    copy(mr32[GetIdx(band*BandW,0,Dim):GetIdx((band+1)*BandW,0,Dim)], data)
  }

  elapsed := time.Since(startTime)
  fmt.Println(elapsed)

  badCount := 0
  for i:=0; i<Dim*Dim; i++ {
    if mr32[i].Delta > 1e-3 {
      badCount += 1
    }
  }
  fmt.Println(badCount)
}
