package main

import (
       _ "os"
       _ "fmt"
       "time"
       _ "math/rand"
       //".dynfloat_fairness"
)

//declares a constant for the remainder of the program
const workers = 10
const datasize = 20000
const global_delta = 0.1

//THIS IS A PLACEHOLDER!
type DynFloat struct {
     d float64
}


//declare a salary type/class
type Salary struct {
     amount float64
}


func population_model() Salary {
     p := Salary{10.0}
     return p
}

func AvgSalary(arr []Salary) float64{
     var total float64 = 0.0
     for i:=0; i<len(arr);i++ {
     	 total := total + arr[i].amount
     }
     return total
}

func ReadFromChannel(num int; input_channel chan Salary) (arr []Salary){
  for i:=0; i<(num); i++ {
    val := <- input_channel
    arr = append(arr,val)
  }
  return
}

//each worker computes the average of their set of data
func worker(input_channel chan Salary,output_channel chan DynFloat){
     //read in everything from the channel
     var data []Salary = ReadFromChannel(input_channel)
     var avg float64 = AvgSalary(data)

     //how to get from float64 representing the average to a DynFloat??

}


func main(){

     var data = get_input_data()

     var input_channels [workers]chan Salary
     //IMPORTANT! Need to use new DynFloat isntead of old DynFairnessFloat!!
     var output_channels [workers]chan DynFloat

     //initialize the channels
     for i:=0; i<workers;i++{
     	 input_channels[i] = make(chan Salary)
	 output_channels[i] = make(chan DynFloat)
     }

     //timing info for experiments
     var startTime = time.now()

     //start go routines
     for j:=0;j<workers;j++{
     	 go worker(input_channels[j],output_channels[j])
     }



}