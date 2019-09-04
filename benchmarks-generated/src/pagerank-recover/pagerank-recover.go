package main

import "fmt"
import "parallely"
import "time"
import "os"
import  "io/ioutil"
import  "strings"
import  "strconv"

func max(x, y int) int {
	if x > y {
		return x;
	} else {
		return y;
	}
}

func min(x, y int) int {
	if x < y {
		return x;
	} else {
		return y;
	}
}

func convertToFloat(x int) float64 {
	return float64(x)
}
	
var Edges []int;
var Inlinks []int;
var Outlinks []int;
var Pageranks []float64;


func func_0() {
  defer parallely.Wg.Done()
  var temppageranks [1000]float64;
var newPagerank float64;
var temp float64;
var outNFloat float64;
var current float64;
var i int;
var j int;
var index int;
var inlink int;
var neighbor int;
var outN int;
for __temp_0 := 0; __temp_0 < 10; __temp_0++ {
 i = 0;
for __temp_1 := 0; __temp_1 < 1000; __temp_1++ {
 newPagerank = 0.15;
inlink=Inlinks[i];
j = 0;
for __temp_2 := 0; __temp_2 < inlink; __temp_2++ {
 index = i*1000+j;
neighbor=Edges[index];
outN=Outlinks[neighbor];
outNFloat=convertToFloat(outN);
current=Pageranks[neighbor];
__flag_1 := false;
 temp = parallely.RandchoiceFlagFloat64(float32(0.99), 0.85*current/outNFloat, 0, &__flag_1);

 if __flag_1 {
 __flag_1 = false;
 temp = parallely.RandchoiceFlagFloat64(float32(0.99), 0.85*current/outNFloat, 0, &__flag_1);

 }
 
newPagerank = newPagerank+temp;
j = j+1;
 }
temppageranks[i]=newPagerank;
i = i+1;
 }
i = 0;
for __temp_3 := 0; __temp_3 < 1000; __temp_3++ {
 temp=temppageranks[i];
Pageranks[i]=temp;
i = i+1;
 }
 }

  fmt.Println("Ending thread : ", 0);
}

func main() {
	fmt.Println("Starting main thread");
	
	parallely.InitChannels(1);
	
  data_bytes, _ := ioutil.ReadFile("../../../benchmarks/inputs/node1000.txt")
  Num_nodes := 1000 // strconv.Atoi(os.Args[2])

  fmt.Println("Starting reading the file")
  data_string := string(data_bytes)
  data_str_array := strings.Split(data_string, "\n")

  fmt.Println("Setting up the data structures")
  Edges = make([]int, Num_nodes * 1000)
  Inlinks = make([]int, Num_nodes)
  Outlinks = make([]int, Num_nodes)
  Pageranks = make([]float64, Num_nodes)

  for i := range Inlinks{
    Inlinks[i] = 0
    Outlinks[i] = 0
    Pageranks[i] = 0.15
  }

  //fmt.Println("Populating the data structures")
  for i := 1; i<len(data_str_array)-1 ; i++ {
    elements := strings.Fields(data_str_array[i])
    index1, _ := strconv.Atoi(elements[0])
    index2, _ := strconv.Atoi(elements[1])

		// fmt.Println(Num_nodes * 1000, index2 * 1000 + Inlinks[index2], index1, index2, Inlinks[index2])
    Edges[(index2 * 1000) + Inlinks[index2]] = index1
    Inlinks[index2]++
    Outlinks[index1]++
		// fmt.Println("---------------")
  }

  fmt.Println("Starting the iterations")
  startTime := time.Now()

	go func_0();


	fmt.Println("Main thread waiting for others to finish");  
	parallely.Wg.Wait()
  elapsed := time.Since(startTime)

	fmt.Println("Done!");
  fmt.Println("Elapsed time : ", elapsed.Nanoseconds());

  f, _ := os.Create("output.txt")
  defer f.Close()

  for i := range Pageranks{
    f.WriteString(fmt.Sprintln(Pageranks[i]))
  }
}
