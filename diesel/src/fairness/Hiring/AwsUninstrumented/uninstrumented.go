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



func population_model() (int,float64,float64) {
	is_male := rand.Intn(2)
	col_rank := rand.NormFloat64() * 10 + 25
	var years_exp float64
	if (is_male > 0) {
		years_exp = rand.NormFloat64() * 5 + 15
	} else {
		years_exp = rand.NormFloat64() * 5 + 10
	}
	return is_male,col_rank,years_exp
}


func offer_job(gender int, col_rank float64, years_exp float64) int {

	if (col_rank <= 5) {
		return 1
	} else if (years_exp > -5) {
		return 1
	} else {
		return 0
	}
}

func getData(genders []int, college_rank []float64, years_exp []float64){
	var gender int
	var col float64
	var years float64

	for i:=0; i < len(genders); i++ {
		gender,col,years = population_model()
		genders[i] = gender
		college_rank[i] = col
		years_exp[i] = years
	}	
}	





var Q = []int {1};
const processors = 1
const datasize = 80000
const dataPerProcess = datasize/processors
const delta = 0.01


func func_Q(ind int){
  defer diesel.Wg.Done();
	//fmt.Println("Starting workers");
	var genders [dataPerProcess] int 
	var college_rank [dataPerProcess] float64 
	var years_exp [dataPerProcess] float64

	//what we stick into the Receive function has to have a fixed size
	diesel.ReceiveIntArray(genders[:],ind,0)
	diesel.ReceiveFloat64Array(college_rank[:],ind,0)
	diesel.ReceiveFloat64Array(years_exp[:],ind,0)

	var hire int	

	for i:=0; i < dataPerProcess; i++ {
		hire = offer_job(genders[i],college_rank[i],years_exp[i])
		_ = hire
		
	}

	
}

func main() {


	//fmt.Println("Starting main thread");

	var genders [datasize] int 
	var college_rank [datasize] float64 
	var years_exp [datasize] float64

	//creates the data by sampling the population model. Don't count this in the timing.
	getData(genders[:],college_rank[:],years_exp[:])
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
		diesel.SendFloat64Array(college_rank[start_ind:end_in],0,q)
		diesel.SendFloat64Array(years_exp[start_ind:end_in],0,q)
	}

	diesel.Wg.Done();
	//diesel.Wg.Wait()

	end := time.Now()
	elapsed := end.Sub(startTime)
	fmt.Println("Elapsed time :", elapsed.Nanoseconds())
}

