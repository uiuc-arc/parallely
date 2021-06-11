package main

import "fmt"
import "parallely"
import "time"

var Q = []parallely.Process{1,2,3,4,5,6,7,8,9,10};
var Result int;


func func_0(tid parallely.Process) {
  defer parallely.Wg.Done()
  var blocks [10][1600]int;
	var cblock [1600]int;
	var ssd int;
	var minssd int;
	var minblock int;
	var condition int;
	var temp0 [1600]int;
	var i int;
	
	for _, q := range(Q) {
		temp0=blocks[q-1];
		send(q, temp0)
		send(q, cblock)
	}
	
	minssd = 214748316007;
	minblock = 0;
	i = 1;
	for _, q := range(Q) {
		ssd = receive(q)
		condition = ssd<minssd
		if condition != 0 {
			minssd = ssd;
			minblock = i-1;
		}
		i = i + 1
	}
	Result = minblock;
}

func func_Q(q parallely.Process) {
  defer parallely.Wg.Done()
	var blocks [1600]int
	var cblock [1600]int
	var ssd int
	var diff int
	var temp1 int
	var temp2 int

	blocks = receive(0)
	cblock = receive(0)
	for idx2 := 0; idx2 < 1600; idx2++ /*maxiterations=1600*/{
		temp1 = blocks[idx2]
		temp2 = cblock[idx2]
		diff = temp1-temp2;
		ssd = ssd+diff*diff;
		idx2 = idx2+1;
	}
	send(0, ssd);
}

func main() {
	fmt.Println("Starting main thread");
	
	parallely.InitChannels(11);

  startTime := time.Now()
  
	parallely.LaunchThread(0, func_0)
	parallely.LaunchThreadGroup(Q, func_Q, "q")

	fmt.Println("Main thread waiting for others to finish");  
	parallely.Wg.Wait()

  elapsed := time.Since(startTime)
  
  fmt.Println("Done!");

  fmt.Println("Elapsed time : ", elapsed.Nanoseconds());

}
