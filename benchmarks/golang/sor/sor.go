package main

import "fmt"
import "parallely"
import "math/rand"
import "time"
import "os"

var input_array [1000000]float64;

var Q = []int {1,2,3,4,5,6,7,8,9,10};


func func_0() {
  // defer parallely.Wg.Done()
  var slice [100000]float32;
	var i int;
	var iter int;	
	var idx int;
	var elem float32;
	var tempelem float64;
	var lowarray [1000000]float32;
	
	lowarray = Parallely.CastF32(input_array)
	
	for iter := 0; iter < 10; iter++ /*maxiterations=10*/ {
		for _, q := range(Q) {
			send(q, lowarray);
		}
		for _, q := range(Q) {
			slice = receive(0)
			for i := 0; i < 100000; i++ /*maxiterations=100000*/ {
				idx = (q-1)*100000+i;
				elem=slice[i];
				tempelem = float64(elem);
				input_array[idx]=tempelem;
			}
		}
	}

}

func func_Q(q int) {
  // defer parallely.Wg.Done()
	var array [1000000]float32;
	var slice [100000]float32;
	var i int;
	var tempi int;	
	var iter int;
	var iter2 int;
	var iter3 int;	
	var j int;
	var k int;
	var conditional int;
	var point float32;
	var temp float32;
	var temp1 float32;
	var temp2 float32;
	var temp3 float32;
	var temp4 float32;
	var temp5 float32;

	for iter := 0; iter < 10; iter++ /*maxiterations=10*/ {
		array = receive(0)
		i = (q-1)*10
		k = 0
		for iter2 := 0; iter2 < 10; iter2++ /*maxiterations=10*/ {
			conditional = ((i<999)&&(i>0))
			if conditional != 0 {
				j = 1;
				for iter3 := 0; iter3 < 998; iter3++ /*maxiterations=998*/ {
					tempi = i*1000+j-1
					temp1 = array[tempi]
					tempi = i*1000+j+1
					temp2 = array[tempi]
					tempi = (i-1)*1000+j
					temp3 = array[tempi]
					tempi = (i-1)*1000+j
					temp4 = array[tempi]
					tempi = (i+1)*1000+j 
					temp5 = array[tempi]
					temp = temp1 + temp2 + temp3 + temp4 + temp5
					point = 0.2*temp;
					slice[k]=point;
					j = j+1;
					k = k+1;
				}
			} else {
				k = k+1000;
			}
			i = i+1;
		}
		send(0, slice)
	}
}


func initArray() {
  for i := 0; i < 1000000; i++ {
		input_array[i] = rand.Float64()
  }
}

func main() {
	fmt.Println("Starting main thread");
	
	parallely.InitChannels(11);

  initArray();

  startTime := time.Now()
  
	parallely.LaunchThread(0, func_0)
	parallely.LaunchThreadGroup(Q, func_Q, "q")
	
	fmt.Println("Main thread waiting for others to finish");  
	parallely.Wg.Wait()

  elapsed := time.Since(startTime)
  
  fmt.Println("Done!");

  fmt.Println("Elapsed time : ", elapsed.Nanoseconds());

  fmt.Println("Writing to output file");
  f, _ := os.Create("__output__.txt")
  defer f.Close()
  fmt.Fprintln(f, input_array)
}
