package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"strings"
	"math"
	"time"
	"strconv"
)

type Option struct {
 K         float64 // strike price
 S0        float64 // strike at time 0
 r         float64 // risk free rate
 sigma     float64 // volatility
 eval_date string // current time
 exp_date  string // expiration date
 T         float64 // distance between exp and current
 right     string // ‘C’ = call, ‘P’ = put
 price     float64 // derived from info above
 delta     float64 // derived from info above
 theta     float64 // derived from info above
 gamma     float64 // derived from info above
}

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
    }	else {
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

func blackscholes(stockstrings [][]float64, channel chan float64) {
	var mystocks []float64
	for i := 0; i < len(stockstrings); i++ {
		c := stockstrings[i]
		optionPrice := BlkSchlsEqEuroNoDiv(c[0], c[1], c[2], c[4], c[5], c[6], c[8])
		mystocks = append(mystocks, optionPrice)
	}
	fmt.Println(mystocks[:5])
	for i := 0; i < len(mystocks); i++ {
		channel <- mystocks[i]
	}
}

func sum(s []int, c chan int) {
	fmt.Println(s)
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}

func main() {
	num_threads := 4
	
	data_bytes, _ := ioutil.ReadFile("in_4K.txt")
	data_string := string(data_bytes)
	data_str_array := strings.Split(data_string, "\n")

	var data_array []([] float64)

	for i := 1; i<len(data_str_array) ; i++ {
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

	workperthread := len(data_str_array)/num_threads
	fmt.Println("Work per threads :", workperthread)
	channels := make([]chan float64, num_threads)
	for i := range channels {
		channels[i] = make(chan float64, 100)
	}
		
	for i := 0; i < num_threads; i++ {
		go blackscholes(data_array[workperthread*i:workperthread*(i+1)], channels[i])
	}
	time.Sleep(5 * time.Second)

	var results []float64

	for i := 0; i < num_threads; i++ {
		for j:=0; j < workperthread; j++ {
			result := <-channels[i]
			results = append(results, result)
		}
	}

	f, _ := os.Create("output.txt")
	defer f.Close()

	
	for i := range results{
		f.WriteString(fmt.Sprintln(results[i]))	
	}
}
