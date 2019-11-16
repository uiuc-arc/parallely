package main

import (
	_ "os"
	_ "fmt"
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
	is_male := rand.Intn(1)
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

func SendPersonArr(arr []Person, ind, num int, chout chan Person) {
  for i:=0; i<num; i++ {
    chout <- arr[ind+i]
  }
}

func RecvPersonArr(arr []Person, num int, chin chan Person) {
  for i:=0; i<num; i++ {
    arr[i] = <- chin
  }
}



func fairness_func(i int, channelin chan Person,dynchannelout chan DynFairnessFloat){
     var dynamic_fairness_map map[string]DynFairnessFloat 
     _ = dynamic_fairness_map
}

func main() {

    

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
  

}
