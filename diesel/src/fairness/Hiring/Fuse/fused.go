package main

import (
	_ "os"
	"fmt"
	_ "io/ioutil"
	_ "strings"
	"math"
	_ "time"
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

func hoeffding(n int, delta float64) (eps float64) {
	eps = math.Sqrt((0.6*math.Log((math.Log(float64(1.1*float64(n+1)))/math.Log(1.10)))+0.555*math.Log(24/delta))/float64(n+1))
	return
}



var Q = []int {1,2,3,4,5,6,7,8};
const processors = 8
const datasize = 8000
const dataPerProcess = datasize/processors
const delta = 0.01


func func_Q(ind int){
  defer diesel.Wg.Done();
	fmt.Println("Starting workers");
	var genders [dataPerProcess] int 
	var college_rank [dataPerProcess] float64 
	var years_exp [dataPerProcess] float64

	//what we stick into the Receive function has to have a fixed size
	diesel.ReceiveIntArray(genders[:],ind,0)
	diesel.ReceiveFloat64Array(college_rank[:],ind,0)
	diesel.ReceiveFloat64Array(years_exp[:],ind,0)

	var hire int
	var males float64 =  0
	var females float64 = 0
	var hiredMales float64 = 0 
	var hiredFemales float64 = 0
	var maleHireProb float64 = 1
	var femaleHireProb float64 = 1
	var probs [2] float64
	_ = probs
	var eps float64 = 1
	var DynMap [2]diesel.ProbInterval;



	for i:=0; i < dataPerProcess; i++ {
		hire = offer_job(genders[i],college_rank[i],years_exp[i])
		if (genders[i] == 1){
			males = males + 1
			hiredMales = hiredMales + float64(hire)
			eps = hoeffding(int(males),delta)
			maleHireProb = hiredMales / males
			//This is what the explicit track statement does
			DynMap[0].Reliability = float32(eps) 
			DynMap[0].Delta =  delta / processors 

		} else {
			females = females + 1
			hiredFemales = hiredFemales + float64(hire)
			eps = hoeffding(int(females),delta)
			femaleHireProb = hiredFemales / females
			//This is what the explicit track statement does
			DynMap[1].Reliability = float32(eps) 
			DynMap[1].Delta = delta / processors
		}
		
	}

	diesel.SendFloat64(maleHireProb,ind,0)
	diesel.SendDynVal(DynMap[0], ind, 0) 
	diesel.SendFloat64(femaleHireProb,ind,0)
	diesel.SendDynVal(DynMap[1], ind, 0) 

	//probs[0] = maleHireProb
	//probs[1] = femaleHireProb
	//diesel.SendDynFloat64Array(probs[:],ind,0,DynMap[:],0)
	
}

func main() {

  // defer diesel.Wg.Done();


	fmt.Println("Starting main thread");

	var genders [datasize] int 
	var college_rank [datasize] float64 
	var years_exp [datasize] float64

	//creates the data by sampling the population model. Don't count this in the timing.
	getData(genders[:],college_rank[:],years_exp[:])
	
	var HireProbs [2*processors]float64
	_ = HireProbs
	var MaleHireProbs [processors]float64
	_ = MaleHireProbs
	var FemaleHireProbs [processors]float64
	_ = FemaleHireProbs


	var MaleHireDynMap [processors]diesel.ProbInterval
	var FemaleHireDynMap [processors]diesel.ProbInterval		
	
	//STARTS (Initializes) the processes
	diesel.InitChannels(9);
//	for _, index := range Q {
	for q := 1; q <= processors; q++ {
		go func_Q(q);
	
	}

	//send the data chunks to each processor
//	for _, q := range Q {
	for q := 1; q <= processors; q++ {

		var start_ind = (q-1)*(dataPerProcess)
		var end_in = q*dataPerProcess
		diesel.SendIntArray(genders[start_ind:end_in],0,q)
		diesel.SendFloat64Array(college_rank[start_ind:end_in],0,q)
		diesel.SendFloat64Array(years_exp[start_ind:end_in],0,q)
	}

	//get the dyn tracked vals from each processor
	for q := 1; q <= processors; q++ {
		//var ind = 2*(q-1)
		//diesel.ReceiveDynFloat64Array(HireProbs[ind:ind+2],0,q,DynMap[:],ind)
		diesel.ReceiveDynVal(&MaleHireDynMap[q-1],0,q)
		fmt.Println(MaleHireDynMap[q-1])
		diesel.ReceiveDynVal(&FemaleHireDynMap[q-1],0,q)
		fmt.Println(MaleHireDynMap[q-1])
	}
	

	//compute the ratio



	diesel.Wg.Done();
	diesel.Wg.Wait()


}


