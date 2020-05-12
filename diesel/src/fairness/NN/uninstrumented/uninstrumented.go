package main
//taken from: https://github.com/sedrews/fairsquare/blob/master/oopsla/noqual/M_BNc_F_NN_V2_H1.fr
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

func population_model() (gender int, age, education_num float64) {
        //var age float64
	//var education_num float64
	gender = bernoulli(0.667)
	if (gender < 1) {
	   education_num = gaussian(568.4105, 24248365.5428)
	   if (education_num<7298.00){
	      age = gaussian(38.4208, 184.9151)
	      education_num = gaussian(10.0827, 6.5096)	
	   } else {
	     age = gaussian(38.8125, 193.4918)
             education_num = gaussian(10.1041, 6.1522)
	   }
	} else {
	   education_num = gaussian(1329.3700, 69327473.1006)
	   if (education_num<5178.00){
	      age  = gaussian(38.6361, 187.2435)
            education_num = gaussian(10.0817, 6.4841)
	   } else {
	     age = gaussian(38.2668, 187.2747)
             education_num = gaussian(10.0974, 7.1793)
	   }
	}
    if (education_num > age){
        age = education_num 
    }
	return
}



func offer_loan(gender int, age float64, education_num float64) (t int) {
	
    N_age := ((age - 17.0) / 73.0  - 0.5) * 10 + 0.5
    N_education_num := ((education_num - 3.0) / 13.0  - 0.5) * 10 + 0.5
    h1 :=  0.1718 * N_age +  1.1416 * N_education_num +  0.4754
    if (h1 < 0){
        h1 = 0
    }
    o1 :=  0.4778 * h1 +  1.2091
    if (o1 < 0){
        o1 = 0
    }
    o2 :=  1.9717 * h1 + -0.3104
    if (o2 < 0){
        o2 = 0
    }


    if (o1 < o2){
	t = 1
    } else {
	t = 0
    }
    return 

}

func getData(genders []int, ages []float64, education_nums []float64){
	var gender int
	var age float64
	var education_num float64

	for i:=0; i < len(genders); i++ {
		gender,age,education_num = population_model() //return a random sample
		genders[i] = gender
		ages[i] = age
		education_nums[i] = education_num
	}	
}	





var Q = []int {1,2,3,4,5,6,7,8};
const processors = 8
const datasize = 80000
const dataPerProcess = datasize/processors
const delta = 0.01


func func_Q(ind int){
  defer diesel.Wg.Done();
	//fmt.Println("Starting workers");
	var genders [dataPerProcess] int 
	var ages [dataPerProcess] float64 
	var education_nums [dataPerProcess] float64

	//what we stick into the Receive function has to have a fixed size
	diesel.ReceiveIntArray(genders[:],ind,0)
	diesel.ReceiveFloat64Array(ages[:],ind,0)
	diesel.ReceiveFloat64Array(education_nums[:],ind,0)

	var classification int


	for i:=0; i < dataPerProcess; i++ {
		classification = offer_loan(genders[i],ages[i],education_nums[i])
		_ = classification
	}



	
}

func main() {


	//fmt.Println("Starting main thread");

	var genders [datasize] int 
	var ages [datasize] float64 
	var education_nums [datasize] float64

	//creates the data by sampling the population model. Don't count this in the timing.
	getData(genders[:],ages[:],education_nums[:])
	
	startTime := time.Now()


	//STARTS (Initializes) the processes
	diesel.InitChannels(9);

	for q := 1; q <= processors; q++ {
		go func_Q(q);
	
	}

	for q := 1; q <= processors; q++ {

		var start_ind = (q-1)*(dataPerProcess)
		var end_in = q*dataPerProcess
		diesel.SendIntArray(genders[start_ind:end_in],0,q)
		diesel.SendFloat64Array(ages[start_ind:end_in],0,q)
		diesel.SendFloat64Array(education_nums[start_ind:end_in],0,q)
	}


	diesel.Wg.Done();
	diesel.Wg.Wait()


	end := time.Now()
	elapsed := end.Sub(startTime)
	fmt.Println("Elapsed time :", elapsed.Nanoseconds())

}


