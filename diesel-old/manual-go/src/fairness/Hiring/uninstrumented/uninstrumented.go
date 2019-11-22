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
	."dynfloat_fairness"
)


const workers = 10
const datasize = 20000
const global_delta = 0.1


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


func SendPersonArr(arr []Person, chout chan Person) {

  for i:=0; i<len(arr); i++ {
    chout <- arr[i]
  }
  
}

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

func fairness_func(i int, channelin chan Person, dynchannelout chan DynFairnessFloat){
     //var dynamic_fairness_map map[string]DynFairnessFloat 
     //_ = dynamic_fairness_map
     var data [] Person
     var decisions [] int
     /* var males int
     var females int
     var hired_males int
     var hired_females int
     var epsilon float64
     var delta float64
     delta = global_delta/(2*workers)    

     var hired_male_mean DynFairnessFloat = DynFairnessFloat{Val:0.,Epsilon:0.,Delta:delta}
     var hired_female_mean DynFairnessFloat = DynFairnessFloat{Val:0.,Epsilon:0.,Delta:delta}
     var ratio DynFairnessFloat = DynFairnessFloat{Val: 0.,Epsilon:0.,Delta:delta}
     */	
     //receive the Persons data array
     data = RecvPersonArr((datasize/workers),channelin)
     for i := range data { //data works

        
	//fmt.Println(epsilon)
	decisions = append(decisions,offer_job(data[i]))
        /*if data[i].gender > 0 {
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
	}*/


        //ratio = DivFloatFairness(hired_male_mean,hired_female_mean)

     }

     //dynchannelout <- ratio   
     
}

func main() {
     
     //fmt.Println("starting")   
     //var c = 0.8
     var data = get_input_data()
     //var means [workers] DynFairnessFloat
     //var LHS DynFairnessFloat

     //a send and recieve channel for each worker to the master
     var channels [workers]chan Person	//really an array of channels (for each worker) sends the persons to each worker
     var dynchannels [workers]chan DynFairnessFloat //channels used to send the dynamically tracked mean back to the master

     //make/initialize the channels
     for i := 0; i< workers; i++ {
     	 channels[i] = make(chan Person)
	 dynchannels[i] = make(chan DynFairnessFloat) //can change later
     }

     var startTime = time.Now()

     //start the goroutines
     for i:=0; i<workers; i++ {
     	 go fairness_func(i,channels[i],dynchannels[i])
     }


     //send the Persons data array
     for i:=0; i < workers; i++ {
	start_ind := (i*(datasize/workers))

	end_ind := start_ind + (datasize/workers)

	SendPersonArr(data[start_ind:end_ind], channels[i])
     }


    //get the local dynamically tracked eps-delta expressions
    /* for j:=0; j < workers; j++ {
	means[j] = <- dynchannels[j]
	//fmt.Println(means[j])
	if j == 0 {
		LHS = means[0]
	} else {
		LHS = AddFloatFairness(LHS,means[j])
	}
     }
     LHS = ConstMulFloatFairness((float64(1)/float64(workers)),LHS)
     //fmt.Println(LHS)		

     check := CheckIneq(LHS,(LHS.Val-c),global_delta)
     _ = check
     //fmt.Println(check)
     */
     var elapsed = time.Since(startTime)
     //fmt.Println("instrumented took ", elapsed)
     fmt.Println(elapsed.Nanoseconds())
}
