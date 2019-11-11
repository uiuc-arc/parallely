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

func blackscholes(chin chan [][]float32, chout chan float32) {
  stockstrings := <- chin
  var mystocks []float32
  for i := 0; i < len(stockstrings); i++ {
    c := stockstrings[i]
    optionPrice := BlkSchlsEqEuroNoDiv(c[0], c[1], c[2], c[4], c[5], c[6], c[8])
    mystocks = append(mystocks, optionPrice)
  }
  for i := 0; i < len(mystocks); i++ {
    chout <- float32(mystocks[i])
  }
}

func main() {
  num_threads := 8

  data_bytes, _ := ioutil.ReadFile("in_4K.txt")
  data_string := string(data_bytes)
  data_str_array := strings.Split(data_string, "\n")

  var data_array []([] float32)

  for i := 1; i<len(data_str_array) ; i++ {
    elements := strings.Split(data_str_array[i], " ")
    var floatelements []float32
    for j:=0; j<len(elements); j++ {
      if j==6 {
        if elements[j] == "P" {
          floatelements = append(floatelements, 1)
        } else {
          floatelements = append(floatelements, 0)
        }
        continue
      }
      s, _ := strconv.ParseFloat(elements[j], 64)
      floatelements = append(floatelements, float32(s))
    }
    data_array = append(data_array, floatelements)
  }

  workperthread := len(data_str_array)/num_threads
  coutput := make([]chan float32, num_threads)
  for i := range coutput {
    coutput[i] = make(chan float32, workperthread)
  }
  cinput := make([]chan [][]float32, num_threads)
  for i := range cinput {
    cinput[i] = make(chan [][]float32, 1)
  }

  for i := 0; i < num_threads; i++ {
    go blackscholes(cinput[i], coutput[i])
  }

  var results []float32

  startTime := time.Now()

  for i := 0; i < num_threads; i++ {
    cinput[i] <- data_array[workperthread*i:workperthread*(i+1)]
  }

  for i := 0; i < num_threads; i++ {
    for j:=0; j < workperthread; j++ {
      result := <- coutput[i]
      results = append(results, float32(result))
    }
  }

  elapsed := time.Since(startTime)
  fmt.Println(elapsed)
}
