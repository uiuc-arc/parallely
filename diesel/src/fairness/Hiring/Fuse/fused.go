package main

import (
	_ "os"
	"fmt"
	_ "io/ioutil"
	_ "strings"
	_ "math"
	_ "time"
	_ "strconv"
	_ "math/rand"
	"../../../diesel"
)


var Q = []int {1,2,3,4,5,6,7,8};


func func_Q(i int){
  defer diesel.Wg.Done();
	fmt.Println("Starting workers");
	var genders [] int 

	diesel.ReceiveIntArray(genders,i,0)


	fmt.Println(genders)		//why is this not printing [2]??



}

func main() {

  defer diesel.Wg.Done();

	var genders = [] int {2}	//a trivial array with a single element

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
			
	}

	//diesel.Wg.Done()

	diesel.Wg.Wait()


}


