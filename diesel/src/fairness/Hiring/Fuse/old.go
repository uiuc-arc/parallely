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





const datasize = 20000
const global_delta = 0.1


////////////////////////////////////////////////////////////////////////////////
type Person struct {
	gender int
	col_rank float64
	years_exp float64
}


func population_model() Person {
	is_male := rand.Intn(2)
	col_rank := rand.NormFloat64() * 10 + 25
	var years_exp float64
	if (is_male > 0) {
		years_exp = rand.NormFloat64() * 5 + 15
	} else {
		years_exp = rand.NormFloat64() * 5 + 10
	}

	p := Person{is_male, col_rank, years_exp}
	return p
}


func offer_job(p Person) int {

	if (p.col_rank <= 5) {
		return 1
	} else if (p.years_exp > -5) {
		return 1
	} else {
		return 0
	}
}



func hoeffding(n int, delta float64) (eps float64) {
	eps = math.Sqrt((0.6*math.Log((math.Log(float64(1.1*float64(n+1)))/math.Log(1.10)))+0.555*math.Log(24/delta))/float64(n+1))
	return
}

////////////////////////////////////////////////////////////////////////////////




//give me an array and a channel and I push each person thru the channel
func SendPersonArr(arr []Person, chout chan Person) {

  for i:=0; i<len(arr); i++ {
    chout <- arr[i]
  }
  
}

//give me a channel and a number of people to pull from the channel
func RecvPersonArr(num int, chin chan Person) (arr []Person) {
  for i:=0; i<(num); i++ {
    val := <- chin
    arr = append(arr,val)
  }
  return
}

func get_input_data() []Person {
	var data []Person
	for i:=0;i<datasize;i++{
		data = append(data, population_model())	//accessing slices as opposed to arrays is weird
	}
	return data
}


var Q = []int {1,2,3,4,5,6,7,8};


func func_Q(i int){
  defer diesel.Wg.Done();
	fmt.Println("Starting workers");
	var genders2 [] int 
	_ = genders2
	//var college_ranks [] float32
	//var years_exp [] float32

	 //receive the data sent from the master
	diesel.ReceiveIntArray(genders2,i,0)
	 //diesel.ReceiveFloat32Array(college_ranks,i,0)
	 //diesel.ReceiveFloat32Array(years_exp,i,0)

	fmt.Println(genders2)
	//fmt.Println(college_ranks)
	//fmt.Println(years_exp)
	//send the dynamically tracked values back to the master


}

func main() {

//  defer diesel.Wg.Done();
	var DynMapMaleHiredProbs []diesel.ProbInterval;     
	var DynMapFemaleHiredProbs []diesel.ProbInterval;     
       // var MaleHiredProbs []float32;
	//var FemaleHiredProbs []float32;

	 _ = DynMapMaleHiredProbs
	_ = DynMapFemaleHiredProbs

	var genders = [] int {2}
	_ = genders
	var college_ranks = [] float32 {1.1}
	_ = college_ranks
	var years_exp = [] float32 {2.2}
	_ = years_exp
	 


	//initialize the processors
	fmt.Println("Starting main thread");


	diesel.InitChannels(9);
	
//	for _, index := range Q {
	for q := 1; q < 9; q++ {
		go func_Q(q);
	
	}


	//send the data chunks to each processor
//	for _, q := range Q {
	for q := 1; q < 9; q++ {
		diesel.SendIntArray(genders,0,q)
			
		//diesel.SendFloat32Array(college_ranks,0,q)
		//diesel.SendFloat32Array(years_exp,0,q)
		_ = q
	}

	diesel.Wg.Done()

	diesel.Wg.Wait()
	/*
	//get the results back from each processor
	for _, index := range Q {
		//arguments are: (rec_var []float32, receiver, sender int, DynMap []ProbInterval, start int)
		diesel.ReceiveDynFloat32Array(MaleHiredProbs, 0, index, DynMapMaleHiredProbs, index)
		diesel.ReceiveDynFloat32Array(FemaleHiredProbs, 0, index, DynMapMaleHiredProbs, index)
	
	}*/

}

/*
     //var dynamic_fairness_map map[string]DynFairnessFloat 
     //_ = dynamic_fairness_map
     var data [] Person
     var decisions [] int
     var males int
     var females int
     var hired_males int
     var hired_females int
     var epsilon float64
     var m_epsilon float64
     var delta float64
     delta = global_delta/(2*workers)    


     var DynMapMaleHiredProbs [1]diesel.ProbInterval; 
     var DynMapFemaleHiredProbs [1]diesel.ProbInterval; 



	//Receive my chunk of the data to compute statistics on
	diesel.ReceiveIntArray(&s_height, tid, 0);
	diesel.ReceiveIntArray(&s_width, tid, 0);
	diesel.ReceiveIntArray(&d_width, tid, 0);
	diesel.ReceiveIntArray(&ts_height, tid, 0);

	
     //receive the Persons data array
     data = RecvPersonArr((datasize/workers),channelin)
     for i := range data { //data works

        
	//fmt.Println(epsilon)
	decisions = append(decisions,offer_job(data[i]))
        if data[i].gender > 0 {
		males = males + 1
		epsilon = math.Sqrt((0.6*math.Log((math.Log(float64(1.1*float64(males+1)))/math.Log(1.10)))+0.555*math.Log(24/delta))/float64(males+1))
		if decisions[i] > 0 {
			hired_males = hired_males + 1
		}
		hired_male_mean.Val = float64(hired_males)/float64(males)
		hired_male_mean.Epsilon = epsilon
		hired_male_mean.Delta = delta
        } else {
		females = females + 1
		epsilon = math.Sqrt((0.6*math.Log((math.Log(float64(1.1*float64(females+1)))/math.Log(1.10)))+0.555*math.Log(24/delta))/float64(females+1))
		if decisions[i] > 0 {
			hired_females = hired_females + 1
		}
		hired_female_mean.Val = float64(hired_females)/float64(females)
		hired_female_mean.Epsilon = epsilon
        	hired_female_mean.Delta = delta
	}


     }

   //  dynchannelout <- ratio   
   */ 

