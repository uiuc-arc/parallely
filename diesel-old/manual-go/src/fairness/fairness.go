package main

import (
	_ "os"
	"fmt"
	_ "io/ioutil"
	_ "strings"
	_ "math"
	_ "time"
	_ "strconv"
	"math/rand"
	."dynfloat_fairness"
)


const workers = 10
const datasize = 10000


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
     var dynamic_fairness_map map[string]DynFairnessFloat 
     _ = dynamic_fairness_map
     var data [] Person

     //receive the Persons data array
     data = RecvPersonArr((datasize/workers),channelin)
     
}

func main() {

     fmt.Println("starting")   

     var data = get_input_data()

     var dynamic_fairness_map map[string]DynFairnessFloat
     _ = dynamic_fairness_map

     // startTime := time.Now()

     //a send and recieve channel for each worker to the master
     var channels [workers]chan Person	//really an array of channels (for each worker) sends the persons to each worker
     var dynchannels [workers]chan DynFairnessFloat //channels used to send the dynamically tracked mean back to the master

     //make/initialize the channels
     for i := 0; i< workers; i++ {
     	 channels[i] = make(chan Person)
	 dynchannels[i] = make(chan DynFairnessFloat) //can change later
     }


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

}
