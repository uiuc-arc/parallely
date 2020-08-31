package main



import (
	_ "os"
	"fmt"
	_ "io/ioutil"
	_ "strings"
	_ "math"
	"time"
	_ "strconv"
	"math/rand"
	"../../../diesel"
)



func bernoulli(p float64) int {
     var x = rand.Float64()>p
     if x {
     	return 1
     } else {
       return 0
     }
}

func gaussian(mu, sigma float64) float64 {
     var x = rand.NormFloat64() * sigma + mu
     return x
}

func population_model() (gender int, age, capital_gain float64) {
        //var age float64
	//var capital_gain float64
	gender = bernoulli(0.667)
	if (gender < 1) {
	   capital_gain = gaussian(568.4105, 24248365.5428)
	   if (capital_gain<7298.00){
	      age = gaussian(38.4208, 184.9151)
	   } else {
	     age = gaussian(38.8125, 193.4918)
	   }
	} else {
	   capital_gain = gaussian(1329.3700, 69327473.1006)
	   if (capital_gain<5178.00){
	      age  = gaussian(38.6361, 187.2435)
	   } else {
	     age = gaussian(38.2668, 187.2747)
	   }
	}

	return

}

func offer_loan(gender int, age float64, capital_gain float64) int {

     	N_age := (age-17.0) / 62.0
	N_capital_gain := (capital_gain)/22040.0
	t := -0.0008 * N_age + -5.7337 * N_capital_gain + 1.0003
	if (gender>0) {
		t = t + 0.0001
	}

	if (t < 0) {
	   return 1
	} else {
	  return 0
	}

}


func getData(genders []int, ages []float64, capital_gains []float64){
	var gender int
	var age float64
	var capital_gain float64

	for i:=0; i < len(genders); i++ {
		gender,age,capital_gain = population_model() //return a random sample
		genders[i] = gender
		ages[i] = age
		capital_gains[i] = capital_gain
	}	
}	



var Q = []int {1,2,3,4,5,6,7,8};
const processors = 1
const datasize = 80000
const dataPerProcess = datasize/processors
const delta = 0.01


func func_Q(ind int){
  defer diesel.Wg.Done();
	//fmt.Println("Starting workers");
	var genders [dataPerProcess] int 
	var ages [dataPerProcess] float64 
	var capital_gains [dataPerProcess] float64

	//what we stick into the Receive function has to have a fixed size
	diesel.ReceiveIntArray(genders[:],ind,0)
	diesel.ReceiveFloat64Array(ages[:],ind,0)
	diesel.ReceiveFloat64Array(capital_gains[:],ind,0)

	var classification int	
	

	for i:=0; i < dataPerProcess; i++ {
		classification = offer_loan(genders[i],ages[i],capital_gains[i])
		_ = classification
	}


	
}

func main() {


	//fmt.Println("Starting main thread");

	var genders [datasize] int 
	var ages [datasize] float64 
	var capital_gains [datasize] float64

	//creates the data by sampling the population model. Don't count this in the timing.
	getData(genders[:],ages[:],capital_gains[:])
	startTime := time.Now()

	

	//STARTS (Initializes) the processes
	diesel.InitChannels(9);
	for q := 1; q <= processors; q++ {
		go func_Q(q);
	
	}

	//send the data chunks to each processor
	for q := 1; q <= processors; q++ {

		var start_ind = (q-1)*(dataPerProcess)
		var end_in = q*dataPerProcess
		diesel.SendIntArray(genders[start_ind:end_in],0,q)
		diesel.SendFloat64Array(ages[start_ind:end_in],0,q)
		diesel.SendFloat64Array(capital_gains[start_ind:end_in],0,q)
	}


	diesel.Wg.Done();
	//diesel.Wg.Wait()

	end := time.Now()
	elapsed := end.Sub(startTime)
	fmt.Println("Elapsed time :", elapsed.Nanoseconds())

}


