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
const processors = 1
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
	var Males diesel.BooleanTracker	= diesel.NewBooleanTracker()		//notice this
	Males.SetDelta(delta/2.)
	var Females diesel.BooleanTracker = diesel.NewBooleanTracker()		//notice this too
	Females.SetDelta(delta/2.)
	var DynMap [2] diesel.ProbInterval

	var probs [2] float64
	

	for i:=0; i < dataPerProcess; i++ {
		classification = offer_loan(genders[i],ages[i],education_nums[i])
		if (genders[i] == 1){
			Males.AddSample(classification)

		} else {
			Females.AddSample(classification)
		}
		
	}


	probs[0] = Males.GetMean()
	probs[1] = Females.GetMean()
	DynMap[0] = Males.GetInterval()
	DynMap[1] = Females.GetInterval()
	
	diesel.SendDynFloat64Array(probs[:],ind,0,DynMap[:],0)
	
}

func main() {


	//fmt.Println("Starting main thread");

	var genders [datasize] int 
	var ages [datasize] float64 
	var education_nums [datasize] float64

	//creates the data by sampling the population model. Don't count this in the timing.
	getData(genders[:],ages[:],education_nums[:])
	startTime := time.Now()

	var tmpDyn [2] diesel.ProbInterval

	var tmpFloats [2] float64
	
	var MaleHighIncomeProb float64
	var MaleHighIncomeProbs [processors]float64
	var FemaleHighIncomeProb float64
	var FemaleHighIncomeProbs [processors]float64
	var Ratio float64	

	var MaleHighIncomeFusedUI diesel.ProbInterval
	var MaleIncomeDynMap [processors]diesel.ProbInterval
	var FemaleHighIncomeFusedUI diesel.ProbInterval
	var FemaleIncomeDynMap [processors]diesel.ProbInterval		
	var RatioUI diesel.ProbInterval
	

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
		diesel.SendFloat64Array(education_nums[start_ind:end_in],0,q)
	}

	//get the dyn tracked vals from each processor
	for q := 1; q <= processors; q++ {
		
		diesel.ReceiveDynFloat64Array(tmpFloats[:],0,q,tmpDyn[:],0)

		MaleIncomeDynMap[q-1]=tmpDyn[0]
		MaleHighIncomeProbs[q-1]=tmpFloats[0]

		FemaleIncomeDynMap[q-1]=tmpDyn[1]
		FemaleHighIncomeProbs[q-1]=tmpFloats[1]

	}

	//I left this the same since we only care about the prob interval (not the total sum or number of samples seen anymore)
	//FUSE everything obtained from each processor
	MaleHighIncomeProb,MaleHighIncomeFusedUI = diesel.FuseFloat64(MaleHighIncomeProbs[:],MaleIncomeDynMap[:])
	FemaleHighIncomeProb,FemaleHighIncomeFusedUI = diesel.FuseFloat64(FemaleHighIncomeProbs[:],FemaleIncomeDynMap[:])

	//compute the ratio
	Ratio,RatioUI = diesel.DivProbInterval(MaleHighIncomeProb,FemaleHighIncomeProb,MaleHighIncomeFusedUI,FemaleHighIncomeFusedUI)
	//fmt.Println(RatioUI)
	//fmt.Println(Ratio)
	diesel.CheckFloat64(Ratio,RatioUI,float32(Ratio-0.8),delta)

	diesel.Wg.Done();
	//diesel.Wg.Wait()

	end := time.Now()
	elapsed := end.Sub(startTime)
	fmt.Println("Elapsed time :", elapsed.Nanoseconds())

}


