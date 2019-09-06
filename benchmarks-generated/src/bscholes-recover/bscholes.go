package main

import (
  "os"
  "fmt"
  "io/ioutil"
  "strings"
  "math"
  "time"
  "math/rand"
  "strconv"
  "parallely"
)

const inv_sqrt_2xPI = 0.39894228040143270286

func CNDF(InputX float64) float64 {
  var sign int;

  var OutputX float64;
  var xInput float64;
  var xNPrimeofX float64;
  var expValues float64;
  var xK2 float64;
  var xK2_2, xK2_3 float64;
  var xK2_4, xK2_5 float64;
  var xLocal, xLocal_1 float64;
  var xLocal_2, xLocal_3 float64;

  // Check for negative value of InputX
  if InputX < 0.0 {
      InputX = -InputX;
      sign = 1;
  }  else {
    sign = 0;
  }

  xInput = InputX;

  // Compute NPrimeX term common to both four & six decimal accuracy calcs
  expValues = math.Exp(-0.5 * InputX * InputX);
  xNPrimeofX = expValues;
  xNPrimeofX = xNPrimeofX * inv_sqrt_2xPI;

  xK2 = 0.2316419 * xInput;
  xK2 = 1.0 + xK2;
  xK2 = 1.0 / xK2;
  xK2_2 = xK2 * xK2;
  xK2_3 = xK2_2 * xK2;
  xK2_4 = xK2_3 * xK2;
  xK2_5 = xK2_4 * xK2;

  xLocal_1 = xK2 * 0.319381530;
  xLocal_2 = xK2_2 * (-0.356563782);
  xLocal_3 = xK2_3 * 1.781477937;
  xLocal_2 = xLocal_2 + xLocal_3;
  xLocal_3 = xK2_4 * (-1.821255978);
  xLocal_2 = xLocal_2 + xLocal_3;
  xLocal_3 = xK2_5 * 1.330274429;
  xLocal_2 = xLocal_2 + xLocal_3;

  xLocal_1 = xLocal_2 + xLocal_1;
  xLocal   = xLocal_1 * xNPrimeofX;
  xLocal   = 1.0 - xLocal;

  OutputX  = xLocal;

  if sign==1 {
      OutputX = 1.0 - OutputX;
  }

  return OutputX;
}

func BlkSchlsEqEuroNoDiv(sptprice, strike, rate, volatility, time float64, otype float64, timet float64) float64 {
  var OptionPrice float64

  // local private working variables for the calculation
  var xStockPrice, xStrikePrice, xRiskFreeRate, xVolatility, xTime, xSqrtTime float64
  var logValues, xLogTerm, xD1 float64
  var xD2, xPowerTerm, xDen, d1, d2, FutureValueX, NofXd1, NofXd2, NegNofXd1, NegNofXd2 float64

  xStockPrice = sptprice;
  xStrikePrice = strike;
  xRiskFreeRate = rate;
  xVolatility = volatility;

  xTime = time;
  xSqrtTime = math.Sqrt(xTime);

  logValues = math.Log( xStockPrice / xStrikePrice );

  xLogTerm = logValues;


  xPowerTerm = xVolatility * xVolatility;
  xPowerTerm = xPowerTerm * 0.5;

  xD1 = xRiskFreeRate + xPowerTerm;
  xD1 = xD1 * xTime;
  xD1 = xD1 + xLogTerm;

  xDen = xVolatility * xSqrtTime;
  xD1 = xD1 / xDen;
  xD2 = xD1 -  xDen;

  d1 = xD1;
  d2 = xD2;

  NofXd1 = CNDF( d1 );
  NofXd2 = CNDF( d2 );

  FutureValueX = xStrikePrice * ( math.Exp( -(rate)*(time) ) );
  if (otype == 0) {
    OptionPrice = (xStockPrice * NofXd1) - (FutureValueX * NofXd2);
  } else {
    NegNofXd1 = (1.0 - NofXd1);
    NegNofXd2 = (1.0 - NofXd2);
    OptionPrice = (FutureValueX * NegNofXd2) - (xStockPrice * NegNofXd1);
  }

  return OptionPrice;
}

func main() {
  rand.Seed(time.Now().UTC().UnixNano())

  data_bytes, _ := ioutil.ReadFile("in_4K.txt")
  data_string := string(data_bytes)
  data_str_array := strings.Split(data_string, "\n")

  var data_array []([] float64)

  for i := 1; i<len(data_str_array)-1 ; i++ {
    elements := strings.Split(data_str_array[i], " ")
    var floatelements []float64
    for j:=0; j<len(elements); j++ {
      if j==6 {
        if elements[j] == "P" {
          floatelements = append(floatelements, 1);
        } else {
          floatelements = append(floatelements, 0);
        }
        continue
      }
      s, _ := strconv.ParseFloat(elements[j], 64)
      floatelements = append(floatelements, s)
    }
    data_array = append(data_array, floatelements)
  }

  var results []float64

  overallflag := false

  for i := 0; i < len(data_array); i++ {
    c := data_array[i]
    var optionPrice float64
    flag := false
    optionPrice = BlkSchlsEqEuroNoDiv(c[0], c[1], c[2], c[4], c[5], c[6], c[8])
    optionPrice = parallely.RandchoiceFlagFloat64(0.999, optionPrice, 0, &flag)
    if flag {
      flag = false
      optionPrice = BlkSchlsEqEuroNoDiv(c[0], c[1], c[2], c[4], c[5], c[6], c[8])
      optionPrice = parallely.RandchoiceFlagFloat64(0.9999, optionPrice, 0, &flag)
    }
    overallflag = overallflag || flag
    results = append(results,optionPrice)
  }

  if overallflag {
    fmt.Print(1," ")
    exact_result_bytes, _ := ioutil.ReadFile("output-exact.txt")
    exact_result_strs := strings.Split(string(exact_result_bytes), "\n")
    l2diff := 0.0
    l2a := 0.0
    l2b := 0.0
    for i := 0; i < len(data_array); i++ {
      exact, _ := strconv.ParseFloat(exact_result_strs[i], 64)
      diff := results[i] - exact
      l2diff += diff*diff
      l2a += exact*exact
      l2b += results[i]*results[i]
    }
    fmt.Println(math.Sqrt(l2diff/(l2a*l2b)))
  } else {
    fmt.Println(0)
  }

  /*f, _ := os.Create("output.txt")
  defer f.Close()

  for i := range results{
    f.WriteString(fmt.Sprintln(results[i]))
  }*/
}
