package main

import "fmt"
import "parallely"
import "math/rand"
import "time"
import "os"
import "strconv"

var input_array [1000000]float64;

var Q = []int {1,2,3,4,5,6,7,8,9,10};


func func_0() {
  // defer parallely.Wg.Done()
  var slice [100000]float32
	var i int
	var idx int
	var elem float32
	var tempelem float64
	var lowarray [1000000]float32
	
	lowarray = ParallelyCastF32(input_array)
	
	for _, q := range(Q) {
		send(q, lowarray);
	}
	for _, q := range(Q) {
		slice = receive(q)
		for i := 0; i < 100000; i++ /*maxiterations=100000*/ {
			idx = (q-1)*100000+i;
			elem=slice[i];
			tempelem = float64(elem);
			input_array[idx]=tempelem;
			i = i+1;
		}
	}
}

func func_Q(q int) {
  // defer parallely.Wg.Done()
	var array [1000000]float32;
	var slice [100000]float32;
	var i int;
	var iter int;
	var j int;
	var k int;
	var conditional int;
	var point float32;
	var tempi int;	

	array = receive(0)
	i = (q-1)*100
	k = 0
	for iter = 0; iter < 10; iter++ /*maxiterations=10*/ {
		conditional = ((i<999)&&(i>0))
		if conditional != 0 {
			for j := 0; j < 998; j++ /*maxiterations=998*/ {
				tempi = (i-1)*1000+j-1
				lm = array[tempi]
				tempi = (i-1)*1000+j
				m = array[tempi]
				tempi = (i-1)*1000+j+1
				rm = array[tempi]
				tempi = (i+1)*1000+j-1
				rb = array[tempi]
				tempi = (i+1)*1000+j
				mb = array[tempi]
				tempi = (i+1)*1000+j
				lb = array[tempi]
				
				point = lm + m + rm + rb + mb + lb
				slice[k]=point
				k = k+1
			}
		} else {
			for j := 0; j < 998; j++ /*maxiterations=998*/ {
				tempi = i*1000+j
				point =	array[tempi]
				slice[k] = point
				j = j+1
				k = k+1
			}
		}
		i = i+1
	}
	send(0, slice)
}


func initArray() {
  for i := 0; i < 1000000; i++ {
     input_array[i] = rand.Float64()
  }
}

func main() {
	fmt.Println("Starting main thread");

  inputseed, _ := strconv.Atoi(os.Args[1])
	rand.Seed(int64(inputseed))
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
