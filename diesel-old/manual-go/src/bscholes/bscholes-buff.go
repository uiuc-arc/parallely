package main

import (
  "fmt"
  "io/ioutil"
  "strings"
  "math"
  "strconv"
  "time"
  ."dynfloats"
)

const optimize = true

func CNDF(InputX float32, dInputX float64) (float32, float64) {
    var D [14]float64
    var sign int // 0
    var OutputX float32 // 1
    var xInput float32 // 2
    var xNPrimeofX float32 // 3
    var expValues float32 // 4
    var xK2 float32 // 5
    var xK2_2, xK2_3, xK2_4, xK2_5, xLocal, xLocal_1, xLocal_2, xLocal_3 float32
    //    6      7      8      9      10       11        12        13

    // Check for negative value of InputX
    if InputX < 0.0 {
        InputX = -InputX
        sign = 1
    }  else {
      sign = 0
    }
    D[0] = 0.0

    xInput = InputX
    D[2] = dInputX

    // Compute NPrimeX term common to both four & six decimal accuracy calcs
    prod := float32(-0.5) * InputX * InputX
    Dprod := math.Abs(float64(InputX))*dInputX + 0.5*dInputX*dInputX
    expValues = float32(math.Exp(float64(prod)))
    D[4] = math.Exp(Dprod+float64(prod))-math.Exp(float64(prod))
    xNPrimeofX = expValues
    D[3] = D[4]
    xNPrimeofX = xNPrimeofX * 0.39894228040143270286
    D[3] = D[3] * 0.39894228040143270286

    xK2 = 0.2316419 * xInput
    D[5] = 0.2316419 * D[2]
    xK2 = 1.0 + xK2
    xK2 = 1.0 / xK2
    D[5] = D[5]/math.Abs(float64(xK2))/(math.Abs(float64(xK2))-D[5])
    xK2_2 = xK2 * xK2
    D[6] = 2*math.Abs(float64(xK2))*D[5] + D[5]*D[5]
    xK2_3 = xK2_2 * xK2
    D[7] = math.Abs(float64(xK2_2))*D[5] + math.Abs(float64(xK2))*D[6] + D[5]*D[6]
    xK2_4 = xK2_3 * xK2
    D[8] = math.Abs(float64(xK2_3))*D[5] + math.Abs(float64(xK2))*D[7] + D[5]*D[7]
    xK2_5 = xK2_4 * xK2
    D[9] = math.Abs(float64(xK2_4))*D[5] + math.Abs(float64(xK2))*D[8] + D[5]*D[8]

    xLocal_1 = xK2 * 0.319381530
    D[11] = D[5] * 0.319381530
    xLocal_2 = xK2_2 * (-0.356563782)
    D[12] = D[6] * (-0.356563782)
    xLocal_3 = xK2_3 * 1.781477937
    D[13] = D[7] * 1.781477937
    xLocal_2 = xLocal_2 + xLocal_3
    D[12] = D[12] + D[13]
    xLocal_3 = xK2_4 * (-1.821255978)
    D[13] = D[8] * (-1.821255978)
    xLocal_2 = xLocal_2 + xLocal_3
    D[12] = D[12] + D[13]
    xLocal_3 = xK2_5 * 1.330274429
    D[13] = D[9] * 1.330274429
    xLocal_2 = xLocal_2 + xLocal_3
    D[12] = D[12] + D[13]

    xLocal_1 = xLocal_2 + xLocal_1
    D[11] = D[12] + D[11]
    xLocal   = xLocal_1 * xNPrimeofX
    D[10] = math.Abs(float64(xLocal_1))*D[3] + math.Abs(float64(xNPrimeofX))*D[11] + D[3]*D[11]
    xLocal   = 1.0 - xLocal

    OutputX  = xLocal
    D[1] = D[10]

    if sign==1 {
        OutputX = 1.0 - OutputX
    }

    return OutputX, D[1]
}

func BlkSchlsEqEuroNoDiv(sptprice, strike, rate, volatility, time, otype, timet float32, dsptprice, dstrike, drate, dvolatility, dtime, dotype, dtimet float64) (float32, float64) {

  var D [20]float64

  // local private working variables for the calculation
  var xStockPrice, xStrikePrice, xRiskFreeRate, xVolatility, xTime, xSqrtTime, logValues, xLogTerm, xD1, xD2, xPowerTerm, xDen, d1, d2, FutureValueX, NofXd1, NofXd2, NegNofXd1, NegNofXd2 float32
  //       0             1              2            3         4        5           6        7       8    9       10       11   12  13       14         15      16        17         18
  var OptionPrice float32 // 19

  xStockPrice = sptprice
  D[0] = dsptprice
  xStrikePrice = strike
  D[1] = dstrike
  xRiskFreeRate = rate
  D[2] = drate
  xVolatility = volatility
  D[3] = dvolatility

  xTime = time
  D[4] = dtime
  xSqrtTime = float32(math.Sqrt(float64(xTime)))
  D[5] = float64(xSqrtTime) - math.Sqrt(float64(xTime)-D[4])

  numerator := float32(math.Log(float64(xStockPrice)))
  Dnumerator := float64(numerator) - math.Log(float64(xStockPrice)-D[0])
  denominator := float32(math.Log(float64(xStrikePrice)))
  Ddenominator := float64(denominator) - math.Log(float64(xStrikePrice)-D[1])
  logValues = numerator - denominator
  D[6] = Dnumerator + Ddenominator

  xLogTerm = logValues
  D[7] = D[6]

  xPowerTerm = xVolatility * xVolatility
  D[10] = 2*math.Abs(float64(xVolatility))*D[3] + D[3]*D[3]
  xPowerTerm = xPowerTerm * 0.5
  D[10] *= 0.5

  xD1 = xRiskFreeRate + xPowerTerm
  D[8] = D[2]+D[10]
  xD1 = xD1 * xTime
  D[8] = math.Abs(float64(xD1))*D[4] + math.Abs(float64(xTime))*D[8] + D[4]*D[8]
  xD1 = xD1 + xLogTerm
  D[8] = D[8]+D[7]

  xDen = xVolatility * xSqrtTime
  D[11] = math.Abs(float64(xVolatility))*D[5] + math.Abs(float64(xSqrtTime))*D[3] + D[5]*D[3]
  xD1 = xD1 / xDen
  D[8] = (math.Abs(float64(xDen))*D[8] + math.Abs(float64(xD1))*D[11])/math.Abs(float64(xDen))/(math.Abs(float64(xDen))-D[11])
  xD2 = xD1 -  xDen
  D[9] = D[8]+D[11]

  d1 = xD1
  D[12] = D[8]
  d2 = xD2
  D[13] = D[9]

  NofXd1, D[15] = CNDF( d1, D[12] )
  NofXd2, D[16] = CNDF( d2, D[13] )

  exp := float32(math.Exp( float64(-(rate)*(time)) ))
  FutureValueX = xStrikePrice * exp
  Dprod := math.Abs(float64(rate))*dtime + math.Abs(float64(time))*drate + drate*dtime
  Dexp := math.Exp(Dprod+float64(-(rate)*(time)))-float64(exp)
  D[14] = math.Abs(float64(xStrikePrice))*Dexp + math.Abs(float64(exp))*D[1] + D[1]*Dexp
  if (otype == 0) {
    OptionPrice = (xStockPrice * NofXd1) - (FutureValueX * NofXd2)
    D[19] = math.Abs(float64(xStockPrice))*D[15] + math.Abs(float64(NofXd1))*D[0] + D[0]*D[15] + math.Abs(float64(FutureValueX))*D[16] + math.Abs(float64(NofXd2))*D[14] + D[14]*D[16]
  } else {
    NegNofXd1 = (1.0 - NofXd1)
    D[17] = D[15]
    NegNofXd2 = (1.0 - NofXd2)
    D[18] = D[16]
    OptionPrice = (FutureValueX * NegNofXd2) - (xStockPrice * NegNofXd1)
    D[19] = math.Abs(float64(xStockPrice))*D[17] + math.Abs(float64(NegNofXd1))*D[0] + D[0]*D[17] + math.Abs(float64(FutureValueX))*D[18] + math.Abs(float64(NegNofXd2))*D[14] + D[14]*D[18]
  }

  return OptionPrice, D[19]
}

func blackscholes(chin, chout chan float32, dchin, dchout chan float64) {
  var data_array [workPerThread*9]float32
  var results [workPerThread]float32
  var D [workPerThread*10]float64
  RecvF32ArrAcc(data_array[:], D[:], workPerThread*9, 0, chin, dchin, optimize)
  for i := 0; i < workPerThread; i++ {
    c := data_array[i*9:i*9+9]
    d := D[i*9:i*9+9]
    results[i], D[workPerThread*9 + i] = BlkSchlsEqEuroNoDiv(c[0], c[1], c[2], c[4], c[5], c[6], c[8], d[0], d[1], d[2], d[4], d[5], d[6], d[8])
  }
  SendF32ArrAcc(results[:], D[workPerThread*9:], workPerThread, 0, chout, dchout, optimize)
}

func main() {
  

  data_bytes, _ := ioutil.ReadFile(inputFileName)
  data_string := string(data_bytes)
  data_str_array := strings.Split(data_string, "\n")

  var data_array [totalWork*9]float32
  var results [totalWork]float32
  var D [totalWork*10]float64

  for i := 1; i<totalWork ; i++ {
    elements := strings.Split(data_str_array[i], " ")
    for j:=0; j<9; j++ {
      if j==6 {
        if elements[j] == "P" {
          data_array[(i-1)*9+j] = float32(1.0)
        } else {
          data_array[(i-1)*9+j] = float32(0.0)
        }
        D[(i-1)*9+j] = 0.0
        continue
      }
      s, _ := strconv.ParseFloat(elements[j], 64)
      data_array[(i-1)*9+j] = float32(s)
      D[(i-1)*9+j] = math.Abs(s-float64(data_array[(i-1)*9+j]))
    }
  }

  var chin [numThreads]chan float32
  var dchin [numThreads]chan float64
  var chout [numThreads]chan float32
  var dchout [numThreads]chan float64
  for i := range chout {
    chin[i] = make(chan float32, workPerThread*9)
    dchin[i] = make(chan float64, workPerThread*9)
    chout[i] = make(chan float32, workPerThread)
    dchout[i] = make(chan float64, workPerThread)
  }

  for i := 0; i < numThreads; i++ {
    go blackscholes(chin[i], chout[i], dchin[i], dchout[i])
  }

  startTime := time.Now()

  for i := 0; i < numThreads; i++ {
    SendF32ArrAcc(data_array[i*workPerThread*9:(i+1)*workPerThread*9], D[i*workPerThread*9:(i+1)*workPerThread*9], workPerThread*9, 0, chin[i], dchin[i], optimize)
  }

  for i := 0; i < numThreads; i++ {
    RecvF32ArrAcc(results[i*workPerThread:(i+1)*workPerThread], D[totalWork*9:], workPerThread, i*workPerThread, chout[i], dchout[i], optimize)
  }

  badCount := 0
  //minDelta := dmap[0]
  //maxDelta := dmap[0]
  for i:=0; i<totalWork; i++ {
    delta := D[totalWork*9+i]
    if delta > 1e-2 {
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

