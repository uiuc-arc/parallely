package main

import "fmt"
import "math"
import "math/rand"
import "time"
import ."dynfloats"

const optimize = true

func mm(band int, chin, chout chan float32, dchin, dchout chan float64) {
  var D[(2*Dim+BandW)*Dim+1][2]float64
  const sumOffset = (2*Dim+BandW)*Dim
  var m1 [Dim*Dim]float32
  var m2 [Dim*Dim]float32
  var mr [BandW*Dim]float32
  for i := range m1 {
    D[i][0] = 1.0
    D[Dim*Dim + i][0] = 1.0
  }
  for i := range mr {
    D[2*Dim*Dim + i][0] = 1.0
  }
  D[sumOffset][0] = 1.0
  RecvF32ArrBoth(m1[:],D[:],Dim*Dim,0,chin,dchin,optimize)
  RecvF32ArrBoth(m2[:],D[Dim*Dim:],Dim*Dim,0,chin,dchin,optimize)
  for i:=band*BandW; i<(band+1)*BandW; i++ {
    for j:=0; j<Dim; j++ {
      sum := float32(0.0)
      D[sumOffset][1] = 0.0
      for k:=0; k<Dim; k++ {
        m1idx := GetIdx(i,k,Dim)
        m2idx := GetIdx(k,j,Dim)
        sum += m1[m1idx] * m2[m2idx]
        D[sumOffset][0] = (D[sumOffset][0] + D[m1idx][0] + D[Dim*Dim + m2idx][0] - 2)*0.99999
        D[sumOffset][1] += math.Abs(float64(m1[m1idx]))*D[Dim*Dim + m2idx][1] + math.Abs(float64(m2[m2idx]))*D[m1idx][1] + D[m1idx][1]*D[Dim*Dim + m2idx][1]
      }
      mr[GetIdx(i-band*BandW,j,Dim)] = sum
      D[2*Dim*Dim + GetIdx(i-band*BandW,j,Dim)][0] = D[sumOffset][0]
      D[2*Dim*Dim + GetIdx(i-band*BandW,j,Dim)][1] = D[sumOffset][1]
    }
  }
  SendF32ArrBoth(mr[:],D[2*Dim*Dim:],BandW*Dim,0,chout,dchout,optimize)
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
  var D[(3*Dim+BandW)*Dim][2]float64

  for i:=0; i<Dim*Dim; i++ {
    m164[i] = randGen.Float64()
    m132[i] = float32(m164[i])
    D[i][1] = math.Abs(m164[i]-float64(m132[i]))
    D[i][0] = 1.0
    m264[i] = randGen.Float64()
    m232[i] = float32(m264[i])
    D[Dim*Dim + i][1] = math.Abs(m264[i]-float64(m232[i]))
    D[Dim*Dim + i][0] = 1.0
    D[2*Dim*Dim + i][0] = 1.0
  }
  for i := range slice {
    D[3*Dim*Dim + i][0] = 1.0
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
    SendF32ArrBoth(m132[:],D[:],Dim*Dim,0,channels[band],dchannels[band],optimize)
    SendF32ArrBoth(m232[:],D[Dim*Dim:],Dim*Dim,0,channels[band],dchannels[band],optimize)
  }
  for band := 0; band < Bands; band++ {
    RecvF32ArrBoth(slice[:],D[3*Dim*Dim:],BandW*Dim,0,channels[band+Bands],dchannels[band+Bands],optimize)
    copy(mr32[GetIdx(band*BandW,0,Dim):GetIdx((band+1)*BandW,0,Dim)], slice[:])
    copy(D[2*Dim*Dim+GetIdx(band*BandW,0,Dim):2*Dim*Dim+GetIdx((band+1)*BandW,0,Dim)], D[3*Dim*Dim:])
  }

  badCount := 0
  //minDelta := dmap[0]
  //maxDelta := dmap[0]
  for i:=0; i<Dim*Dim; i++ {
    delta := D[2*Dim*Dim+i][1]
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
