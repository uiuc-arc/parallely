package main

/*import (
	"os"
	"fmt"
	"io/ioutil"
	"strings"
	"math"
	"time"
	"strconv"
	"math/rand"
	."dynfloats"
)*/

const workers = 10
const datasize = 10000

type DynFairness struct {
  epsilon float64
  delta float64
  
}

func fairness_func(i int, channelin, channelout chan float32,dynchannelin,dynchannelout chan DynFairness){
     var dynamic_fairness_map map[string]DynFairness 
}

func main() {

     startTime := time.Now()

     var dynamic_fairness_map map[string]DynFairness 
     var genders [datasize]float32
     var college_ranks [datasize]float32
     var years_exp [datasize]float32

     //a send and recieve channel for each worker to the master
     var channels [workers*2]chan float32
     var dynchannels [workers*2]chan DynFairness

     //make the channels
     for i := range channels {
     	 channels[i] = make(chan float32,3*datasize/workers)
	 dynchannels[i] = make(chan DynFairness) //can change later
     }

     //start the goroutines
     for i:=0; i<workers; i++ {
     	 go fairness_func(i,channels[i],channels[i+workers],dynchannels[i],dynchannels[i+workers])
     }

     

     

}