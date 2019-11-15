package main

import (
  "fmt"
  "io/ioutil"
  "strings"
  "math"
  "strconv"
  "time"
)

func CNDF(InputX float32) float32 {
    var sign int

    var OutputX float32
    var xInput float32
    var xNPrimeofX float32
    var expValues float32
    var xK2 float32
    var xK2_2, xK2_3 float32
    var xK2_4, xK2_5 float32
    var xLocal, xLocal_1 float32
    var xLocal_2, xLocal_3 float32

    // Check for negative value of InputX
    if InputX < 0.0 {
        InputX = -InputX
        sign = 1
    }  else {
      sign = 0
    }

    xInput = InputX

    // Compute NPrimeX term common to both four & six decimal accuracy calcs
    expValues = float32(math.Exp(float64(-0.5 * InputX * InputX)))
    xNPrimeofX = expValues
    xNPrimeofX = xNPrimeofX * 0.39894228040143270286

    xK2 = 0.2316419 * xInput
    xK2 = 1.0 + xK2
    xK2 = 1.0 / xK2
    xK2_2 = xK2 * xK2
    xK2_3 = xK2_2 * xK2
    xK2_4 = xK2_3 * xK2
    xK2_5 = xK2_4 * xK2

    xLocal_1 = xK2 * 0.319381530
    xLocal_2 = xK2_2 * (-0.356563782)
    xLocal_3 = xK2_3 * 1.781477937
    xLocal_2 = xLocal_2 + xLocal_3
    xLocal_3 = xK2_4 * (-1.821255978)
    xLocal_2 = xLocal_2 + xLocal_3
    xLocal_3 = xK2_5 * 1.330274429
    xLocal_2 = xLocal_2 + xLocal_3

    xLocal_1 = xLocal_2 + xLocal_1
    xLocal   = xLocal_1 * xNPrimeofX
    xLocal   = 1.0 - xLocal

    OutputX  = xLocal

    if sign==1 {
        OutputX = 1.0 - OutputX
    }

    return OutputX
}

func BlkSchlsEqEuroNoDiv(sptprice, strike, rate, volatility, time float32, otype float32, timet float32) float32 {
  var OptionPrice float32

  // local private working variables for the calculation
  var xStockPrice, xStrikePrice, xRiskFreeRate, xVolatility, xTime, xSqrtTime float32
  var logValues, xLogTerm, xD1 float32
  var xD2, xPowerTerm, xDen, d1, d2, FutureValueX, NofXd1, NofXd2, NegNofXd1, NegNofXd2 float32

  xStockPrice = sptprice
  xStrikePrice = strike
  xRiskFreeRate = rate
  xVolatility = volatility

  xTime = time
  xSqrtTime = float32(math.Sqrt(float64(xTime)))

  logValues = float32(math.Log( float64(xStockPrice / xStrikePrice) ))

  xLogTerm = logValues


  xPowerTerm = xVolatility * xVolatility
  xPowerTerm = xPowerTerm * 0.5

  xD1 = xRiskFreeRate + xPowerTerm
  xD1 = xD1 * xTime
  xD1 = xD1 + xLogTerm

  xDen = xVolatility * xSqrtTime
  xD1 = xD1 / xDen
  xD2 = xD1 -  xDen

  d1 = xD1
  d2 = xD2

  NofXd1 = CNDF( d1 )
  NofXd2 = CNDF( d2 )

  FutureValueX = xStrikePrice * ( float32(math.Exp( float64(-(rate)*(time)) )) )
  if (otype == 0) {
    OptionPrice = (xStockPrice * NofXd1) - (FutureValueX * NofXd2)
  } else {
    NegNofXd1 = (1.0 - NofXd1)
    NegNofXd2 = (1.0 - NofXd2)
    OptionPrice = (FutureValueX * NegNofXd2) - (xStockPrice * NegNofXd1)
  }

  return OptionPrice
}

func blackscholes(chin chan float32, chout chan float32) {
  var data_array [workPerThread*9]float32
  var results [workPerThread]float32
  for i := 0; i < workPerThread*9; i++ {
    data_array[i] = <- chin
  }
  for i := 0; i < workPerThread; i++ {
    c := data_array[i*9:i*9+9]
    results[i] = BlkSchlsEqEuroNoDiv(c[0], c[1], c[2], c[4], c[5], c[6], c[8])
  }
  for i := 0; i < workPerThread; i++ {
    chout <- results[i]
  }
}

func main() {
  data_bytes, _ := ioutil.ReadFile(inputFileName)
  data_string := string(data_bytes)
  data_str_array := strings.Split(data_string, "\n")

  var data_array [totalWork*9]float32
  var results [totalWork]float32

  for i := 1; i<totalWork ; i++ {
    elements := strings.Split(data_str_array[i], " ")
    for j:=0; j<9; j++ {
      if j==6 {
        if elements[j] == "P" {
          data_array[(i-1)*9+j] = float32(1.0)
        } else {
          data_array[(i-1)*9+j] = float32(0.0)
        }
        continue
      }
      s, _ := strconv.ParseFloat(elements[j], 64)
      data_array[(i-1)*9+j] = float32(s)
    }
  }

  var chin [numThreads]chan float32
  var chout [numThreads]chan float32
  for i := range chout {
    chin[i] = make(chan float32, workPerThread*9)
    chout[i] = make(chan float32, workPerThread)
  }

  for i := 0; i < numThreads; i++ {
    go blackscholes(chin[i], chout[i])
  }

  startTime := time.Now()

  for i := 0; i < numThreads; i++ {
    for j := 0; j < workPerThread*9; j++ {
      chin[i] <- data_array[i*workPerThread*9+j]
    }
  }

  for i := 0; i < numThreads; i++ {
    for j:=0; j < workPerThread; j++ {
      result := <- chout[i]
      results[i*workPerThread+j] = result
    }
  }

  elapsed := time.Since(startTime)
  fmt.Println(elapsed.Nanoseconds())
}
