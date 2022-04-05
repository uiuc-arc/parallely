package main

import (
  "math"
  "fmt"
	"time"
  "diesel"
	"math/rand"
)

func Min(x, y int) int {
  if x < y {
    return x
  }
  return y
}

func Max(x, y int) int {
  if x > y {
    return x
  }
  return y
}

func Idx(i, j, width int) int {
  return i*width+j
}

func floorInt(input float64) int {
	return int(math.Floor(input));
}

func ceilInt(input float64) int {
	return int(math.Ceil(input));
}

func convertToFloat(x int) float64 {
	return float64(x)
}

var NumThreads int
var Iterations int
var Q = []int {1,2,3,4,5,6,7,8};
var EdgeDevices = []int {1,2,3,4,5,6,7,8};

func func_0() {
  defer diesel.Wg.Done();
  var DynMap [16432]diesel.ProbInterval;
  _ = DynMap;

	var data [16384]float64;
	diesel.InitDynArray(0, 16384, DynMap[:]);

	var centerIds [8]int;
	var centers [16]float64;
	diesel.InitDynArray(16384, 16, DynMap[:]);
	
	var centerSlice [16]float64;	
	diesel.InitDynArray(16384+16, 16, DynMap[:]);

	var tempcenters [16]float64;
	diesel.InitDynArray(16384+16*2, 16, DynMap[:]);	

	var realCenters  [8]float64
	for i:=0; i<len(realCenters)/2; i++ {
		// diesel.ReceiveFloat64(&temp, 0, q);
		realCenters[2*i] = 30 + rand.Float64() * 5
		// diesel.ReceiveFloat64(&temp, 0, q);
		realCenters[2*i+1] = 40 + rand.Float64() * 10;
	}
	// for i:=len(data)/4; i<len(data)/2; i++ {
	// 	// diesel.ReceiveFloat64(&temp, 0, q);
	// 	data[2*i] = 40 + rand.Float64() * 10
	// 	// diesel.ReceiveFloat64(&temp, 0, q);
	// 	data[2*i+1] = 40 + rand.Float64() * 20;
	// }
	
	for i:=0; i<16384/2; i++ {
		clusterNew := rand.Intn(4)
		data[2*i] = rand.NormFloat64() * 0.5 + realCenters[2*clusterNew]
		DynMap[2*i] = diesel.ProbInterval{1, 1};;
		data[2*i+1] = rand.NormFloat64() * 0.5 + realCenters[2*clusterNew+1]
	  DynMap[2*i+1] = diesel.ProbInterval{1, 2};
	}

  // f, _ := os.Create("output.txt")
  // defer f.Close()
	// f.WriteString(fmt.Sprintln(data))

	for i, _ := range(centerIds) {
		centerIds[i] = rand.Intn(16384/2)
	}
	
	for i, _ := range(centerIds) {
		centers[2*i] = data[2*centerIds[i]];
		DynMap[16384 + 2*i] = DynMap[2*centerIds[i]]
		centers[2*i+1] = data[2*centerIds[i]+1];
		DynMap[16384 + 2*i + 1] = DynMap[2*centerIds[i]+1]
	}
	// fmt.Println("Intial : ", DynMap[16384:16384+1]);

	// fmt.Println("Initial  Centers: ", centers)
	// ficenter, _ := os.Create("init-centers.txt")  
	// ficenter.WriteString(fmt.Sprintln(realCenters))
	// ficenter.Close()

	for _, q := range(EdgeDevices) {
		// diesel.SendFloat64Array(data[:], 0, q);
		diesel.SendDynFloat64ArrayO1(data[:], 0, q, DynMap[:], 0);
	}

	for __temp_0 := 0; __temp_0 < Iterations; __temp_0++ {
		// fmt.Println("Centers Start : ", centers);
		for _, q := range(EdgeDevices) {
			// diesel.SendFloat64Array(centers[:], 0, q);
			diesel.SendDynFloat64ArrayO1(centers[:], 0, q, DynMap[:], 16);
		}

		// Reset the centers
		for i, _ := range(centerIds) {
			tempcenters[2*i] = 0;
			DynMap[2080 + 2*i] = diesel.ProbInterval{1, 0}
			tempcenters[2*i+1] = 0;
			DynMap[2080 + 2*i + 1] = diesel.ProbInterval{1, 0}			
		}

		for _, q := range(EdgeDevices) {
			diesel.NoisyReceiveDynFloat64ArrayO1(centerSlice[:], 0, q, DynMap[:], 16384+16);
			// fmt.Println("(t0) received : ", DynMap[2064:2080]);
			// Reset the centers
			for i, _ := range(centerIds) {
				tempcenters[2*i] = tempcenters[2*i] + centerSlice[2*i];
				DynMap[16384 + 16*2 + 2*i].Reliability = DynMap[16384 + 16*2 + 2*i].Reliability * DynMap[16384 + 16 + 2*i].Reliability
				DynMap[16384 + 16*2 + 2*i].Delta = DynMap[16384 + 16*2 + 2*i].Delta + DynMap[16384 + 16*2 + 2*i].Delta
				
				tempcenters[2*i+1] = tempcenters[2*i+1] + centerSlice[2*i+1];
				DynMap[16384 + 16*2 + 2*i + 1].Reliability = DynMap[16384 + 16*2 + 2*i + 1].Reliability *
					DynMap[16384 + 16 + 2*i + 1].Reliability
				DynMap[16384 + 16*2 + 2*i + 1].Delta = DynMap[16384 + 16*2 + 2*i + 1].Delta + DynMap[16384 + 16 + 2*i + 1].Delta				
			}
		}
		// fmt.Println("(t0) received : ", DynMap[2080:2096]);

		// New centers are the average
		for i, _ := range(centers) {
			tempcenters[i] = tempcenters[i] / float64(len(EdgeDevices));
			// DynMap[2080 + 2*i].Reliability = DynMap[2080 + 2*i].Reliability * DynMap[2064 + 2*i].Reliability
			DynMap[16384 + 16*2 + i].Delta = DynMap[2080 + i].Delta / float64(len(EdgeDevices))
			centers[i] = tempcenters[i]
			DynMap[16384 + 16 + i] = DynMap[16384 + 16*2 + i]
		}
		// fmt.Println("End of Iteration : ", DynMap[2048:2064]);
		// centers = tempcenters
		// fmt.Println("Centers : ", centers);
		diesel.PrintWorstElement(DynMap[:], 16384+16, 16)		
	}
	// fcenter, _ := os.Create("centers.txt")
  // defer fcenter.Close()
	// fcenter.WriteString(fmt.Sprintln(centers))

	// fd, _ := os.Create("dynmap.txt")
  // defer fd.Close()
	// fd.WriteString(fmt.Sprintln(DynMap[16384:16384+16]))
  fmt.Println("Ending thread : ", 0);
}

// func func_Q(tid int) {
//   defer diesel.Wg.Done();
//   var DynMap [0]diesel.ProbInterval;
//   _ = DynMap;
//   q := tid;
	
// 	temperature := rand.Float64();
// 	humidity := rand.Float64();

// 	diesel.SendInt(temperature, q, 0);
// 	diesel.SendInt(humidity, q, 0);
// }

func func_Q(tid int) {
  defer diesel.Wg.Done();
  var DynMap [16418]diesel.ProbInterval;
  _ = DynMap;
	q := tid;

	var data [16384]float64;
	diesel.InitDynArray(0, 16384, DynMap[:]);
	
	var centers [16]float64;
	diesel.InitDynArray(16384, 16, DynMap[:]);

	var mindist float64;
	var dist float64;
	DynMap[16384 + 16 + 1] = diesel.ProbInterval{1, 0};	
	
	var avgcenters [16]float64;
	diesel.InitDynArray(16384 + 16 + 1 + 1, 16, DynMap[:]);
	
	var countcenters [8]float64;

	var assignment [1024]int;

	mystart := (tid - 1) * 1024;
	myend := mystart + 1024;
	// fmt.Println(tid, mystart, myend);

	diesel.NoisyReceiveDynFloat64ArrayO1(data[:], tid, 0, DynMap[:], 0);
	
	for __temp_i := 0; __temp_i < Iterations; __temp_i++ {
		diesel.NoisyReceiveDynFloat64ArrayO1(centers[:], q, 0, DynMap[:], 16384);

		// if q == 1 {
		// 	fmt.Println("(t0) : ", DynMap[2048:2064]);
		// }

		for __temp_1 := 0; __temp_1 < 8; __temp_1++ {
			avgcenters[__temp_1*2] = 0
			DynMap[16384 + 16 + 2 +__temp_1*2] = diesel.ProbInterval{0.999999, 0}
			avgcenters[__temp_1*2 + 1] = 0
			DynMap[16384 + 16 + 2 +__temp_1*2 + 1] = diesel.ProbInterval{0.999999, 0}			
			countcenters[__temp_1] = 0
		}

		for __temp_0 := mystart; __temp_0 < myend; __temp_0++ {			
			mindist = 1000000;
			mincenter := -1;
			for __temp_1 := 0; __temp_1 < 8; __temp_1++ {
				temp1 := (data[2*__temp_0]-centers[2*__temp_1])
				d_temp1 := (DynMap[0 + 2*__temp_0].Delta + DynMap[0 + 2*__temp_1].Delta)
				temp2 := (data[2*__temp_0+1]-centers[2*__temp_1+1])
				d_temp2 := (DynMap[0 + 2*__temp_0].Delta + DynMap[0 +
					2*__temp_1].Delta)

				temp3 := temp1 * temp1
				d_temp3 := 2*d_temp1*temp3 + d_temp1*d_temp1
				temp4 := temp2 * temp2
				d_temp4 := 2*d_temp2*temp4 + d_temp3*d_temp3
				
				dist = temp3 + temp4
				DynMap[16384 + 16 + 1].Reliability = DynMap[0 + 2*__temp_0].Reliability * DynMap[16384 + 2*__temp_1].Reliability *DynMap[0 + 2*__temp_0+1].Reliability * DynMap[16384 + 2*__temp_1+1].Reliability
				DynMap[16384 + 16 + 1].Delta =  d_temp3 + d_temp4
				// if q == 1 {
				// 	fmt.Println("(t0) : ", DynMap[2065].Delta);
				// }
				
				if dist<mindist {
					mindist = dist
					mincenter = __temp_1
				}
			}
			
			assignment[__temp_0-mystart] = mincenter
			// fmt.Println(mincenter, mindist)
			avgcenters[mincenter*2] += data[2*__temp_0]
			DynMap[16384 + 16 + 1 + 1 + mincenter*2].Delta += DynMap[2*__temp_0].Delta
			DynMap[16384 + 16 + 1 + 1 + mincenter*2].Reliability = DynMap[2*__temp_0].Reliability * DynMap[2066+mincenter*2].Reliability
			
			avgcenters[mincenter*2 + 1] += data[2*__temp_0+1]
			DynMap[16384 + 16 + 1 + 1 + mincenter*2+1].Delta += DynMap[2*__temp_0+1].Delta
			DynMap[16384 + 16 + 1 + 1 + mincenter*2+1].Reliability =
			DynMap[2*__temp_0+1].Reliability *
			DynMap[16384 + 16 + 1 + 1 +mincenter*2+1].Reliability
			
			countcenters[mincenter] += 1
		}

		// if q == 1 {
		// 	fmt.Println("(t0) avgcenters: ", DynMap[2066:2082]);
		// }		
		

		for __temp_1 := 0; __temp_1 < 8; __temp_1++ {
			if countcenters[__temp_1] > 0 {
				avgcenters[2*__temp_1] = avgcenters[2*__temp_1] / countcenters[__temp_1]
				DynMap[16384 + 16 + 1 + 1 +2*__temp_1].Delta = DynMap[16384 + 16 + 1 + 1+2*__temp_1].Delta / countcenters[__temp_1]
				// DynMap[2066+mincenter*2+1].Reliability += DynMap[2*__temp_0+1].Reliability * DynMap[2066+mincenter*2+1].Reliability
				
				avgcenters[2*__temp_1 + 1] = avgcenters[2*__temp_1+1] / countcenters[__temp_1]
				DynMap[16384 + 16 + 1 + 1+2*__temp_1+1].Delta = DynMap[16384 + 16 + 1 + 1+2*__temp_1+1].Delta / countcenters[__temp_1]				
			} else {
				// fmt.Println(tid, __temp_i, countcenters[__temp_1])
				// os.Exit(-1)
				avgcenters[2*__temp_1] = centers[2*__temp_1]
				DynMap[16384 + 16 + 1 + 1+2*__temp_1] = DynMap[16384+2*__temp_1]
				avgcenters[2*__temp_1 + 1] = centers[2*__temp_1+1]
				DynMap[16384 + 16 + 1 + 1+2*__temp_1+1] = DynMap[16384+2*__temp_1+1]				
			}
		}

		// if q == 1 {
		// 	fmt.Println("(t0) avgcenters - 2: ", DynMap[2066:2082]);
		// }
		// fmt.Println(tid, countcenters)
		// fmt.Println(tid, avgcenters)
		diesel.SendDynFloat64ArrayO1(avgcenters[:], q, 0, DynMap[:], 16384 + 16 + 1 + 1);
	}
	fmt.Println("Ending thread : ", q);
}

func main() {
	Iterations = 20
	// rand.Seed(time.Now().UTC().UnixNano())

  fmt.Println("Starting main thread");
  NumThreads = 9;
	
	diesel.InitChannels(9);

  // Sensors, _ = strconv.Atoi(os.Args[1])
  // oFile := os.Args[2]
  
  // src_tmp, s_width, s_height, _ := ReadPpmFile(iFile)
  // SHeight = s_height
  // SWidth = s_width
  // DestSize = len(src_tmp)*4*4

  // for i, _ := range src_tmp {
	// 	Src[i] = float64(src_tmp[i])
  // }

  // ImgSize = len(src_tmp)

	startTime := time.Now()
	go func_0();
	for _, index := range EdgeDevices {
		go func_Q(index);
	}


	fmt.Println("Main thread waiting for others to finish");  
	diesel.Wg.Wait()

	end := time.Now()
	elapsed := end.Sub(startTime)
	fmt.Println("Elapsed time :", elapsed.Nanoseconds())
  diesel.PrintMemory() 

  // tmp_dest := make([]int, len(Dest))
  // for i, _ := range Dest {
	// 	tmp_dest[i] = int(Dest[i])
  // }

  // WritePpmFile(tmp_dest, s_width*4, s_height*4, oFile)
}