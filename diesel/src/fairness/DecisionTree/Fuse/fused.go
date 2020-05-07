package main
//taken from: https://github.com/sedrews/fairsquare/blob/master/oopsla/noqual/M_BN_F_DT_V2_D2_N4.fr (Decision Tree V2, uses same population model)
import (
	_ "os"
	"fmt"
	_ "io/ioutil"
	_ "strings"
	_ "math"
	_ "time"
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

	if (capital_gain >= 7073.5){
		if (age < 20){
		    t = 1
		}
		else{
		    t = 0
		}
	}
	else{
		t = 1
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
const processors = 8
const datasize = 80000
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
	
	var eps float64 = 1
	var DynMap [2]diesel.ProbInterval;



	for i:=0; i < dataPerProcess; i++ {
		hire = offer_loan(genders[i],college_rank[i],years_exp[i])
		if (genders[i] == 1){
			males = males + 1
			hiredMales = hiredMales + float64(hire)
			eps = diesel.Hoeffding(int(males),1-delta/2)
			maleHireProb = hiredMales / males
			//This is what the explicit track statement does
			DynMap[0].Reliability = float32(eps) 
			DynMap[0].Delta =  delta / processors 

		} else {
			females = females + 1
			hiredFemales = hiredFemales + float64(hire)
			eps = diesel.Hoeffding(int(females),1-delta/2)
			femaleHireProb = hiredFemales / females
			//This is what the explicit track statement does
			DynMap[1].Reliability = float32(eps) 
			DynMap[1].Delta = delta / processors
		}
		
	}


	probs[0] = maleHireProb
	probs[1] = femaleHireProb
	diesel.SendDynFloat64Array(probs[:],ind,0,DynMap[:],0)
	
}

func main() {

  // defer diesel.Wg.Done();


	fmt.Println("Starting main thread");

	var genders [datasize] int 
	var college_rank [datasize] float64 
	var years_exp [datasize] float64

	//creates the data by sampling the population model. Don't count this in the timing.
	getData(genders[:],college_rank[:],years_exp[:])
	

	var tmpDyn [2] diesel.ProbInterval

	var tmpFloats [2] float64
	
	var MaleHireProb float64
	var MaleHireProbs [processors]float64
	var FemaleHireProb float64
	var FemaleHireProbs [processors]float64
	var Ratio float64	

	var MaleHireUI diesel.ProbInterval
	var MaleHireDynMap [processors]diesel.ProbInterval
	var FemaleHireUI diesel.ProbInterval
	var FemaleHireDynMap [processors]diesel.ProbInterval		
	var RatioUI diesel.ProbInterval
	

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
		
		diesel.ReceiveDynFloat64Array(tmpFloats[:],0,q,tmpDyn[:],0)

		MaleHireDynMap[q-1]=tmpDyn[0]
		MaleHireProbs[q-1]=tmpFloats[0]

		FemaleHireDynMap[q-1]=tmpDyn[1]
		FemaleHireProbs[q-1]=tmpFloats[1]


	}


	//FUSE everything obtained from each processor
	MaleHireProb,MaleHireUI = diesel.FuseFloat64(MaleHireProbs[:],MaleHireDynMap[:])
	FemaleHireProb,FemaleHireUI = diesel.FuseFloat64(FemaleHireProbs[:],FemaleHireDynMap[:])

	//compute the ratio
	Ratio,RatioUI = diesel.DivProbInterval(MaleHireProb,FemaleHireProb,MaleHireUI,FemaleHireUI)
	fmt.Println(RatioUI)
	fmt.Println(Ratio)
	diesel.CheckFloat64(Ratio,RatioUI,float32(Ratio-0.8),delta)

	diesel.Wg.Done();
	diesel.Wg.Wait()


}


