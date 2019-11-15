package main

import "fmt"
import "math"
import "math/rand"
import "time"
import ."dynfloats"

const optimize = true

func mm(band int, chin, chout chan float32, dchin, dchout chan float64) {
  var D[(2*Dim+BandW)*Dim+1]float64
  var m1 [Dim*Dim]float32
  var m2 [Dim*Dim]float32
  var mr [BandW*Dim]float32
  RecvF32ArrAcc(m1[:],D[:],Dim*Dim,0,chin,dchin,optimize)
  RecvF32ArrAcc(m2[:],D[Dim*Dim:],Dim*Dim,0,chin,dchin,optimize)
  for i:=band*BandW; i<(band+1)*BandW; i++ {
    for j:=0; j<Dim; j++ {
      sum := float32(0.0)
      const sumOffset = (2*Dim+BandW)*Dim
      D[sumOffset] = 0.0
      for k:=0; k<Dim; k++ {
        m1idx := GetIdx(i,k,Dim)
        m2idx := GetIdx(k,j,Dim)
        sum += m1[m1idx] * m2[m2idx]
        D[sumOffset] += math.Abs(float64(m1[m1idx]))*D[Dim*Dim + m2idx] + math.Abs(float64(m2[m2idx]))*D[m1idx] + D[m1idx]*D[Dim*Dim + m2idx]
      }
      mr[GetIdx(i-band*BandW,j,Dim)] = sum
      D[2*Dim*Dim + GetIdx(i-band*BandW,j,Dim)] = D[sumOffset]
    }
  }
  SendF32ArrAcc(mr[:],D[2*Dim*Dim:],BandW*Dim,0,chout,dchout,optimize)
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
  var D[(3*Dim+BandW)*Dim]float64

  for i:=0; i<Dim*Dim; i++ {
    m164[i] = randGen.Float64()
    m132[i] = float32(m164[i])
    D[i] = math.Abs(m164[i]-float64(m132[i]))
    m264[i] = randGen.Float64()
    m232[i] = float32(m264[i])
    D[Dim*Dim + i] = math.Abs(m264[i]-float64(m232[i]))
  }

  var channels [Bands*2]chan float32
  var dchannels [Bands*2]chan float64
  for i := range channels {
    channels[i] = make(chan float32, 2*Dim*Dim)
    dchannels[i] = make(chan float64, 2*Dim*Dim)
  }

  for i:=0; i<Bands; i++ {
    go mm(i, channels[i], channels[i+Bands], dchannels[i], dchannels[i+Bands])
  }

  startTime := time.Now()

  for band := 0; band < Bands; band++ {
    SendF32ArrAcc(m132[:],D[:],Dim*Dim,0,channels[band],dchannels[band],optimize)
    SendF32ArrAcc(m232[:],D[Dim*Dim:],Dim*Dim,0,channels[band],dchannels[band],optimize)
  }
  for band := 0; band < Bands; band++ {
    RecvF32ArrAcc(slice[:],D[3*Dim*Dim:],BandW*Dim,0,channels[band+Bands],dchannels[band+Bands],optimize)
    copy(mr32[GetIdx(band*BandW,0,Dim):GetIdx((band+1)*BandW,0,Dim)], slice[:])
    copy(D[2*Dim*Dim+GetIdx(band*BandW,0,Dim):2*Dim*Dim+GetIdx((band+1)*BandW,0,Dim)], D[3*Dim*Dim:])
  }

  badCount := 0
  //minDelta := dmap[0]
  //maxDelta := dmap[0]
  for i:=0; i<Dim*Dim; i++ {
    delta := D[2*Dim*Dim+i]
    if delta > 1e-4 {
      badCount += 1
    }
    //if maxDelta < delta {
    //  maxDelta = delta
    //}
    //if minDelta > delta {
    //  minDelta = delta
    //}
  }

  elapsed := time.Since(startTime)
  fmt.Println(elapsed.Nanoseconds(), badCount)
}
