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
	age float64
	capital_gain float64
}
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
func population_model() Person {
        var age float64
	var capital_gain float64
	gender := bernoulli(0.667)
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

	p := Person{gender,age, capital_gain}
	return p
}
func F(p Person) int {
     	N_age := (p.age-17.0) / 62.0
	N_capital_gain := (p.capital_gain)/22040.0
	t := -0.0008 * N_age + -5.7337 * N_capital_gain + 1.0003
	if (p.gender>0) {
		t = t + 0.0001
	}
	if t < 0 {
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
     var data [] Person
     var decisions [] int
     data = RecvPersonArr((datasize/workers),channelin)
     for i := range data { //data works
	decisions = append(decisions,F(data[i]))
     }
}
func main() {
     var data = get_input_data()
     var channels [workers]chan Person	//really an array of channels (for each worker) sends the persons to each worker
     var dynchannels [workers]chan DynFairnessFloat //channels used to send the dynamically tracked mean back to the master
     for i := 0; i< workers; i++ {
     	 channels[i] = make(chan Person)
	 dynchannels[i] = make(chan DynFairnessFloat) //can change later
     }
     var startTime = time.Now()
     for i:=0; i<workers; i++ {
     	 go fairness_func(i,channels[i],dynchannels[i])
     }
     for i:=0; i < workers; i++ {
	start_ind := (i*(datasize/workers))
	end_ind := start_ind + (datasize/workers)
	SendPersonArr(data[start_ind:end_ind], channels[i])
     }
     var elapsed = time.Since(startTime)
     fmt.Println(elapsed.Nanoseconds())
}
