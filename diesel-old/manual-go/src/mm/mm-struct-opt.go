package main

import "fmt"
import "math/rand"
import "time"
import ."dynfloats"

func mm(band int, chin, chout chan float32, dchin, dchout chan float64) {
  var m1 [Dim*Dim]DynFloat32
  var m2 [Dim*Dim]DynFloat32
  var mr [BandW*Dim]DynFloat32
  for i := range m1 {
    m1[i].Num = <- chin
  }
  dtemp := <- dchin
  for i := range m1 {
    m1[i].Delta = dtemp
  }
  for i := range m2 {
    m2[i].Num = <- chin
  }
  dtemp = <- dchin
  for i := range m2 {
    m2[i].Delta = dtemp
  }
  for i:=band*BandW; i<(band+1)*BandW; i++ {
    for j:=0; j<Dim; j++ {
      sum := MakeDynFloat32(0)
      for k:=0; k<Dim; k++ {
        sum = AddDynFloat32(sum,MulDynFloat32(m1[GetIdx(i,k,Dim)],m2[GetIdx(k,j,Dim)]))
      }
      mr[GetIdx(i-band*BandW,j,Dim)] = sum
    }
  }
  dtemp = 0.0
  for i := range mr {
    chout <- mr[i].Num
    if mr[i].Delta > dtemp {
      dtemp = mr[i].Delta
    }
  }
  dchout <- dtemp
}

func main() {
  randSource := rand.NewSource(time.Now().UnixNano())
  randGen := rand.New(randSource)
  var m164 [Dim*Dim]DynFloat64
  var m132 [Dim*Dim]DynFloat32
  var m264 [Dim*Dim]DynFloat64
  var m232 [Dim*Dim]DynFloat32
  var mr32 [Dim*Dim]DynFloat32
  var slice [BandW*Dim]DynFloat32

  for i:=0; i<Dim*Dim; i++ {
    m164[i] = MakeDynFloat64(randGen.Float64())
    m132[i] = DynFloat64To32(m164[i])
    m264[i] = MakeDynFloat64(randGen.Float64())
    m232[i] = DynFloat64To32(m264[i])
  }

  var channels [Bands*2]chan float32
  var dchannels [Bands*2]chan float64
  for i := range channels {
    channels[i] = make(chan float32)
    dchannels[i] = make(chan float64)
  }

  for i:=0; i<Bands; i++ {
    go mm(i, channels[i], channels[i+Bands], dchannels[i], dchannels[i+Bands])
  }

  startTime := time.Now()

  for band := 0; band < Bands; band++ {
    dtemp := 0.0
    for i := range m132 {
      channels[band] <- m132[i].Num
      if m132[i].Delta > dtemp {
        dtemp = m132[i].Delta
      }
    }
    dchannels[band] <- dtemp
    for i := range m232 {
      channels[band] <- m232[i].Num
      if m232[i].Delta > dtemp {
        dtemp = m232[i].Delta
      }
    }
    dchannels[band] <- dtemp
  }
  for band := 0; band < Bands; band++ {
    for i := range slice {
      slice[i].Num = <- channels[band+Bands]
    }
    dtemp := <- dchannels[band+Bands]
    for i := range slice {
      slice[i].Delta = dtemp
    }
    copy(mr32[GetIdx(band*BandW,0,Dim):GetIdx((band+1)*BandW,0,Dim)], slice[:])
  }

  badCount := 0
  for i:=0; i<Dim*Dim; i++ {
    if mr32[i].Delta > 1e-3 {
      badCount += 1
    }
  }

  elapsed := time.Since(startTime)
  fmt.Println(elapsed.Nanoseconds(), badCount)
}
