package main

import (
	"diesel"
	"fmt"
	"math"
	"math/rand"
	"time"
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
	return i*width + j
}

func floorInt(input float64) int {
	return int(math.Floor(input))
}

func ceilInt(input float64) int {
	return int(math.Ceil(input))
}

func convertToFloat(x int) float64 {
	return float64(x)
}

var NumThreads int
var Iterations int
var Sensors [1024]float64
var Sensorshumid [1024]float64
var CenterIds [8]int

var Q = []int{1, 2, 3, 4, 5, 6, 7, 8}

func func_0() {
	defer diesel.Wg.Done()
	var DynMap [4124]diesel.ProbInterval
	var my_chan_index int
	_ = my_chan_index
	_ = DynMap
	var datatemp [1024]float64
	diesel.InitDynArray(0, 1024, DynMap[:])
	var datahumid [1024]float64
	diesel.InitDynArray(1024, 1024, DynMap[:])
	var centersTemp [8]float64
	diesel.InitDynArray(2048, 8, DynMap[:])
	var centersHumid [8]float64
	diesel.InitDynArray(2056, 8, DynMap[:])
	var centerSlice [8]float64
	diesel.InitDynArray(2064, 8, DynMap[:])
	var tempcentersTemp [1024]float64
	diesel.InitDynArray(2072, 1024, DynMap[:])
	var tempcentersHumid [1024]float64
	diesel.InitDynArray(3096, 1024, DynMap[:])
	var i int
	var temp int
	var temp0 float64
	var tempf0 float64
	var tempf float64
	DynMap[4120] = diesel.ProbInterval{1, 0}
	var tempf1 float64
	DynMap[4121] = diesel.ProbInterval{1, 0}
	var tempf2 float64
	DynMap[4122] = diesel.ProbInterval{1, 0}
	var temp1 float64
	DynMap[4123] = diesel.ProbInterval{1, 0}
	i = 0
	for __temp_0 := 0; __temp_0 < 1024; __temp_0++ {
		_temp_index_1 := i
		tempf0 = Sensors[_temp_index_1]
		tempf1 = tempf0
		DynMap[4121] = diesel.ProbInterval{1.0, 1.5}
		_temp_index_2 := i
		datatemp[_temp_index_2] = tempf1
		DynMap[0+_temp_index_2] = DynMap[4121]
		_temp_index_3 := i
		tempf0 = Sensorshumid[_temp_index_3]
		tempf1 = tempf0
		DynMap[4121] = diesel.ProbInterval{1.0, 2.0}
		_temp_index_4 := i
		datahumid[_temp_index_4] = tempf1
		DynMap[1024+_temp_index_4] = DynMap[4121]
	}
	i = 0
	for __temp_1 := 0; __temp_1 < 8; __temp_1++ {
		_temp_index_5 := i
		temp = CenterIds[_temp_index_5]
		_temp_index_6 := temp
		tempf = datatemp[_temp_index_6]
		DynMap[4120] = DynMap[0+_temp_index_6]
		_temp_index_7 := i
		centersTemp[_temp_index_7] = tempf
		DynMap[2048+_temp_index_7] = DynMap[4120]
		_temp_index_8 := temp
		tempf = datahumid[_temp_index_8]
		DynMap[4120] = DynMap[1024+_temp_index_8]
		_temp_index_9 := 1
		centersHumid[_temp_index_9] = tempf
		DynMap[2056+_temp_index_9] = DynMap[4120]
		i = i + 1
	}
	for _, q := range Q {
		diesel.SendDynFloat64ArrayO1(datatemp[:], 0, q, DynMap[:], 0)
		diesel.SendDynFloat64ArrayO1(datahumid[:], 0, q, DynMap[:], 1024)
	}
	for __temp_2 := 0; __temp_2 < Iterations; __temp_2++ {
		for _, q := range Q {
			diesel.SendDynFloat64ArrayO1(centersTemp[:], 0, q, DynMap[:], 2048)
			diesel.SendDynFloat64ArrayO1(centersHumid[:], 0, q, DynMap[:], 2056)
		}
		temp0 = 0.0
		i = 0
		for __temp_3 := 0; __temp_3 < 8; __temp_3++ {
			_temp_index_10 := i
			tempcentersTemp[_temp_index_10] = temp0
			DynMap[2072+_temp_index_10] = diesel.ProbInterval{1, 0}
			_temp_index_11 := i
			tempcentersHumid[_temp_index_11] = temp0
			DynMap[3096+_temp_index_11] = diesel.ProbInterval{1, 0}
			i = i + 1
		}
		for _, q := range Q {
			diesel.ReceiveDynFloat64ArrayO1(centerSlice[:], 0, q, DynMap[:], 2064)
			i = 0
			for __temp_4 := 0; __temp_4 < 8; __temp_4++ {
				_temp_index_12 := i
				tempf = tempcentersTemp[_temp_index_12]
				DynMap[4120] = DynMap[2072+_temp_index_12]
				_temp_index_13 := i
				tempf1 = centerSlice[_temp_index_13]
				DynMap[4121] = DynMap[2064+_temp_index_13]
				DynMap[4122].Reliability = DynMap[4123].Reliability + DynMap[4120].Reliability - 1.0
				DynMap[4122].Delta = DynMap[4120].Delta + DynMap[4123].Delta
				tempf2 = tempf + temp1
				_temp_index_14 := i
				tempcentersTemp[_temp_index_14] = tempf2
				DynMap[2072+_temp_index_14] = DynMap[4122]
				i = i + 1
			}
			diesel.ReceiveDynFloat64ArrayO1(centerSlice[:], 0, q, DynMap[:], 2064)
			i = 0
			for __temp_5 := 0; __temp_5 < 8; __temp_5++ {
				_temp_index_15 := i
				tempf = tempcentersHumid[_temp_index_15]
				DynMap[4120] = DynMap[3096+_temp_index_15]
				_temp_index_16 := i
				tempf1 = centerSlice[_temp_index_16]
				DynMap[4121] = DynMap[2064+_temp_index_16]
				DynMap[4122].Reliability = DynMap[4123].Reliability + DynMap[4120].Reliability - 1.0
				DynMap[4122].Delta = DynMap[4120].Delta + DynMap[4123].Delta
				tempf2 = tempf + temp1
				_temp_index_17 := i
				tempcentersHumid[_temp_index_17] = tempf2
				DynMap[3096+_temp_index_17] = DynMap[4122]
				i = i + 1
			}
		}
		i = 0
		for __temp_6 := 0; __temp_6 < 8; __temp_6++ {
			_temp_index_18 := i
			tempf1 = tempcentersTemp[_temp_index_18]
			DynMap[4121] = DynMap[2072+_temp_index_18]
			DynMap[4120].Reliability = DynMap[4121].Reliability
			DynMap[4120].Delta = DynMap[4121].Delta / math.Abs(float64(8.0))
			tempf = tempf1 / 8.0
			_temp_index_19 := i
			centersTemp[_temp_index_19] = tempf
			DynMap[2048+_temp_index_19] = DynMap[4120]
			_temp_index_20 := i
			tempf1 = tempcentersHumid[_temp_index_20]
			DynMap[4121] = DynMap[3096+_temp_index_20]
			DynMap[4120].Reliability = DynMap[4121].Reliability
			DynMap[4120].Delta = DynMap[4121].Delta / math.Abs(float64(8.0))
			tempf = tempf1 / 8.0
			_temp_index_21 := i
			tempcentersHumid[_temp_index_21] = tempf
			DynMap[3096+_temp_index_21] = DynMap[4120]
			i = i + 1
		}
		fmt.Println(DynMap[3096:4004])

	}

	fmt.Println("Ending thread : ", 0)
}
func func_Q(tid int) {
	defer diesel.Wg.Done()
	var DynMap [0]diesel.ProbInterval
	var my_chan_index int
	_ = my_chan_index
	_ = DynMap
	q := tid
	var datatemp [1024]float64
	var datahumid [1024]float64
	var centersTemp [8]float64
	var centersHumid [8]float64
	var tempcentersTemp [1024]float64
	var tempcentersHumid [1024]float64
	var countcenters [8]int
	var assigned [1024]int
	var mystart int
	var myend int
	var perthread int
	var mypoints int
	var i int
	var k int
	var temp0 float64
	var mindist float64
	var mincenter int
	var condition int
	var temp4 float64
	var temp1 float64
	var temp2 float64
	var temp23 float64
	var temp24 float64
	var tempi int
	var data1 float64
	var center1 float64
	perthread = 1024 / 8
	mystart = (q - 1) * perthread
	myend = mystart + perthread
	mypoints = myend - mystart
	diesel.ReceiveFloat64Array(datatemp[:], tid, 0)
	diesel.ReceiveFloat64Array(datahumid[:], tid, 0)
	for __temp_7 := 0; __temp_7 < Iterations; __temp_7++ {
		diesel.ReceiveFloat64Array(centersTemp[:], tid, 0)
		diesel.ReceiveFloat64Array(centersHumid[:], tid, 0)
		temp0 = 0.0
		i = 0
		for __temp_8 := 0; __temp_8 < 8; __temp_8++ {
			_temp_index_1 := i
			tempcentersTemp[_temp_index_1] = temp0
			_temp_index_2 := i
			tempcentersHumid[_temp_index_2] = temp0
			i = i + 1
		}
		i = mystart
		for __temp_9 := 0; __temp_9 < mypoints; __temp_9++ {
			mindist = 1000000.0
			mincenter = 0
			k = 0
			for __temp_10 := 0; __temp_10 < 8; __temp_10++ {
				_temp_index_3 := i
				data1 = datatemp[_temp_index_3]
				_temp_index_4 := k
				center1 = centersTemp[_temp_index_4]
				temp0 = data1 - center1
				_temp_index_5 := i
				data1 = datahumid[_temp_index_5]
				_temp_index_6 := k
				center1 = centersHumid[_temp_index_6]
				temp1 = data1 - center1
				temp23 = temp1 * temp1
				temp24 = temp2 * temp2
				temp4 = temp24 + temp23
				if mindist >= temp4 {
					mindist = temp4
				} else {
					mindist = mindist
				}
				tempi = k
				_ = 1.0
				_ = 0.0
				temp_bool_7 := condition
				if temp_bool_7 != 0 {
					mincenter = tempi
				} else {
					mincenter = mincenter
				}
				k = k + 1
			}
			_temp_index_8 := i
			assigned[_temp_index_8] = mincenter
			_temp_index_9 := mincenter
			tempi = countcenters[_temp_index_9]
			_temp_index_10 := mincenter
			countcenters[_temp_index_10] = tempi + 1
			i = i + 1
		}
		i = mystart
		for __temp_11 := 0; __temp_11 < mypoints; __temp_11++ {
			_temp_index_11 := i
			tempi = assigned[_temp_index_11]
			_temp_index_12 := tempi
			temp1 = tempcentersTemp[_temp_index_12]
			_temp_index_13 := i
			temp2 = datatemp[_temp_index_13]
			temp4 = temp1 + temp2
			_temp_index_14 := i
			tempcentersTemp[_temp_index_14] = temp0
			_temp_index_15 := tempi
			temp1 = tempcentersHumid[_temp_index_15]
			_temp_index_16 := i
			temp2 = datahumid[_temp_index_16]
			temp4 = temp1 + temp2
			_temp_index_17 := i
			tempcentersHumid[_temp_index_17] = temp0
			i = i + 1
		}
		diesel.SendFloat64Array(tempcentersTemp[:], tid, 0)
		diesel.SendFloat64Array(tempcentersHumid[:], tid, 0)
	}

	fmt.Println("Ending thread : ", q)
}

func main() {
	Iterations = 20

	fmt.Println("Starting main thread")
	NumThreads = 9

	diesel.InitChannels(9)

	var realCenters [8]float64

	for i := 0; i < len(realCenters)/2; i++ {
		realCenters[2*i] = 30 + rand.Float64()*5
		realCenters[2*i+1] = 40 + rand.Float64()*10
	}

	for i := 0; i < 1024; i++ {
		clusterNew := rand.Intn(4)
		Sensors[i] = rand.NormFloat64()*0.5 + realCenters[2*clusterNew]
		Sensorshumid[i] = rand.NormFloat64()*0.5 + realCenters[2*clusterNew+1]
	}

	for i, _ := range CenterIds {
		CenterIds[i] = rand.Intn(1024)
	}

	startTime := time.Now()
	go func_0()
	for _, index := range Q {
		go func_Q(index)
	}

	fmt.Println("Main thread waiting for others to finish")
	diesel.Wg.Wait()

	end := time.Now()
	elapsed := end.Sub(startTime)
	fmt.Println("Elapsed time :", elapsed.Nanoseconds())
	diesel.PrintMemory()
}
