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

func func_0() {
	defer diesel.Wg.Done()
	var DynMap [5158]diesel.ProbInterval
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
	var tempcentersTemp [1024]float64
	diesel.InitDynArray(2064, 1024, DynMap[:])
	var tempcentersHumid [1024]float64
	diesel.InitDynArray(3088, 1024, DynMap[:])
	var assigned [1024]int
	diesel.InitDynArray(4112, 1024, DynMap[:])
	var countcenters [8]int
	diesel.InitDynArray(5136, 8, DynMap[:])
	var k int
	var i int
	var temp int
	var temp0m float64
	var tempf0 float64
	var tempf float64
	DynMap[5144] = diesel.ProbInterval{1, 0}
	var tempf1 float64
	DynMap[5145] = diesel.ProbInterval{1, 0}
	var temp0 float64
	DynMap[5146] = diesel.ProbInterval{1, 0}
	var mindist float64
	DynMap[5147] = diesel.ProbInterval{1, 0}
	var mincenter int
	DynMap[5148] = diesel.ProbInterval{1, 0}
	var condition int
	DynMap[5149] = diesel.ProbInterval{1, 0}
	var temp4 float64
	DynMap[5150] = diesel.ProbInterval{1, 0}
	var temp1 float64
	DynMap[5151] = diesel.ProbInterval{1, 0}
	var temp2 float64
	DynMap[5152] = diesel.ProbInterval{1, 0}
	var temp23 float64
	DynMap[5153] = diesel.ProbInterval{1, 0}
	var temp24 float64
	DynMap[5154] = diesel.ProbInterval{1, 0}
	var tempi int
	DynMap[5155] = diesel.ProbInterval{1, 0}
	var data1 float64
	DynMap[5156] = diesel.ProbInterval{1, 0}
	var center1 float64
	DynMap[5157] = diesel.ProbInterval{1, 0}
	i = 0
	for __temp_0 := 0; __temp_0 < 1024; __temp_0++ {
		_temp_index_1 := i
		tempf0 = Sensors[_temp_index_1]
		tempf1 = tempf0
		DynMap[5145] = diesel.ProbInterval{1.0, 1.5}
		_temp_index_2 := i
		datatemp[_temp_index_2] = tempf1
		DynMap[0+_temp_index_2] = DynMap[5145]
		_temp_index_3 := i
		tempf0 = Sensorshumid[_temp_index_3]
		tempf1 = tempf0
		DynMap[5145] = diesel.ProbInterval{1.0, 2.0}
		_temp_index_4 := i
		datahumid[_temp_index_4] = tempf1
		DynMap[1024+_temp_index_4] = DynMap[5145]
	}
	i = 0
	for __temp_1 := 0; __temp_1 < 8; __temp_1++ {
		_temp_index_5 := i
		temp = CenterIds[_temp_index_5]
		_temp_index_6 := temp
		tempf = datatemp[_temp_index_6]
		DynMap[5144] = DynMap[0+_temp_index_6]
		_temp_index_7 := i
		centersTemp[_temp_index_7] = tempf
		DynMap[2048+_temp_index_7] = DynMap[5144]
		_temp_index_8 := temp
		tempf = datahumid[_temp_index_8]
		DynMap[5144] = DynMap[1024+_temp_index_8]
		_temp_index_9 := 1
		centersHumid[_temp_index_9] = tempf
		DynMap[2056+_temp_index_9] = DynMap[5144]
		i = i + 1
	}
	for __temp_2 := 0; __temp_2 < Iterations; __temp_2++ {
		temp0m = 0.0
		i = 0
		for __temp_3 := 0; __temp_3 < 8; __temp_3++ {
			_temp_index_10 := i
			tempcentersTemp[_temp_index_10] = temp0m
			DynMap[2064+_temp_index_10] = diesel.ProbInterval{1, 0}
			_temp_index_11 := i
			tempcentersHumid[_temp_index_11] = temp0m
			DynMap[3088+_temp_index_11] = diesel.ProbInterval{1, 0}
			i = i + 1
		}
		i = 0
		for __temp_4 := 0; __temp_4 < 1024; __temp_4++ {
			DynMap[5147] = diesel.ProbInterval{1, 0}
			mindist = 1000000.0
			DynMap[5148] = diesel.ProbInterval{1, 0}
			mincenter = 0
			k = 0
			for __temp_5 := 0; __temp_5 < 8; __temp_5++ {
				_temp_index_12 := i
				data1 = datatemp[_temp_index_12]
				DynMap[5156] = DynMap[0+_temp_index_12]
				_temp_index_13 := k
				center1 = centersTemp[_temp_index_13]
				DynMap[5157] = DynMap[2048+_temp_index_13]
				DynMap[5146].Reliability = DynMap[5156].Reliability + DynMap[5157].Reliability - 1.0
				DynMap[5146].Delta = DynMap[5156].Delta + DynMap[5157].Delta
				temp0 = data1 - center1
				_temp_index_14 := i
				data1 = datahumid[_temp_index_14]
				DynMap[5156] = DynMap[1024+_temp_index_14]
				_temp_index_15 := k
				center1 = centersHumid[_temp_index_15]
				DynMap[5157] = DynMap[2056+_temp_index_15]
				DynMap[5151].Reliability = DynMap[5156].Reliability + DynMap[5157].Reliability - 1.0
				DynMap[5151].Delta = DynMap[5156].Delta + DynMap[5157].Delta
				temp1 = data1 - center1
				DynMap[5153].Reliability = DynMap[5151].Reliability
				DynMap[5153].Delta = math.Abs(float64(temp1))*DynMap[5151].Delta + math.Abs(float64(temp1))*DynMap[5151].Delta + DynMap[5151].Delta*DynMap[5151].Delta
				temp23 = temp1 * temp1
				DynMap[5154].Reliability = DynMap[5152].Reliability
				DynMap[5154].Delta = math.Abs(float64(temp2))*DynMap[5152].Delta + math.Abs(float64(temp2))*DynMap[5152].Delta + DynMap[5152].Delta*DynMap[5152].Delta
				temp24 = temp2 * temp2
				DynMap[5150].Reliability = DynMap[5153].Reliability + DynMap[5154].Reliability - 1.0
				DynMap[5150].Delta = DynMap[5154].Delta + DynMap[5153].Delta
				temp4 = temp24 + temp23
				mindist = diesel.DynCondFloat64GeqFloat64(mindist, temp4, DynMap[:], 5147, 5150, temp4, mindist, 5150, 5147, 5147)
				tempi = k
				DynMap[5155] = diesel.ProbInterval{1.0, 0.0}
				temp_bool_16 := condition
				if temp_bool_16 != 0 {
					mincenter = tempi
				} else {
					mincenter = mincenter
				}
				if temp_bool_16 != 0 {
					DynMap[5148].Reliability = DynMap[5149].Reliability * DynMap[5155].Reliability
				} else {
					DynMap[5148].Reliability = DynMap[5149].Reliability * DynMap[5148].Reliability
				}
				k = k + 1
			}
			_temp_index_17 := i
			assigned[_temp_index_17] = mincenter
			DynMap[4112+_temp_index_17] = DynMap[5148]
			_temp_index_18 := mincenter
			tempi = countcenters[_temp_index_18]
			DynMap[5155] = DynMap[5136+_temp_index_18]
			_temp_index_19 := mincenter
			countcenters[_temp_index_19] = tempi + 1
			DynMap[5136+_temp_index_19] = DynMap[5155]
			i = i + 1
		}
		i = 0
		for __temp_6 := 0; __temp_6 < 1024; __temp_6++ {
			_temp_index_20 := i
			tempi = assigned[_temp_index_20]
			DynMap[5155] = DynMap[4112+_temp_index_20]
			_temp_index_21 := tempi
			temp1 = tempcentersTemp[_temp_index_21]
			DynMap[5151] = DynMap[2064+_temp_index_21]
			_temp_index_22 := i
			temp2 = datatemp[_temp_index_22]
			DynMap[5152] = DynMap[0+_temp_index_22]
			DynMap[5150].Reliability = DynMap[5152].Reliability + DynMap[5151].Reliability - 1.0
			DynMap[5150].Delta = DynMap[5151].Delta + DynMap[5152].Delta
			temp4 = temp1 + temp2
			_temp_index_23 := i
			tempcentersTemp[_temp_index_23] = temp0
			DynMap[2064+_temp_index_23] = DynMap[5146]
			_temp_index_24 := tempi
			temp1 = tempcentersHumid[_temp_index_24]
			DynMap[5151] = DynMap[3088+_temp_index_24]
			_temp_index_25 := i
			temp2 = datahumid[_temp_index_25]
			DynMap[5152] = DynMap[1024+_temp_index_25]
			DynMap[5150].Reliability = DynMap[5152].Reliability + DynMap[5151].Reliability - 1.0
			DynMap[5150].Delta = DynMap[5151].Delta + DynMap[5152].Delta
			temp4 = temp1 + temp2
			_temp_index_26 := i
			tempcentersHumid[_temp_index_26] = temp0
			DynMap[3088+_temp_index_26] = DynMap[5146]
			i = i + 1
		}
		i = 0
		for __temp_7 := 0; __temp_7 < 8; __temp_7++ {
			_temp_index_27 := i
			tempf1 = tempcentersTemp[_temp_index_27]
			DynMap[5145] = DynMap[2064+_temp_index_27]
			DynMap[5144].Reliability = DynMap[5145].Reliability
			DynMap[5144].Delta = DynMap[5145].Delta / math.Abs(float64(8.0))
			tempf = tempf1 / 8.0
			_temp_index_28 := i
			centersTemp[_temp_index_28] = tempf
			DynMap[2048+_temp_index_28] = DynMap[5144]
			_temp_index_29 := i
			tempf1 = tempcentersHumid[_temp_index_29]
			DynMap[5145] = DynMap[3088+_temp_index_29]
			DynMap[5144].Reliability = DynMap[5145].Reliability
			DynMap[5144].Delta = DynMap[5145].Delta / math.Abs(float64(8.0))
			tempf = tempf1 / 8.0
			_temp_index_30 := i
			tempcentersHumid[_temp_index_30] = tempf
			DynMap[3088+_temp_index_30] = DynMap[5144]
			i = i + 1
		}
	}

	fmt.Println("Ending thread : ", 0)
}

func main() {
	Iterations = 20

	fmt.Println("Starting main thread")
	NumThreads = 1

	diesel.InitChannels(1)

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

	fmt.Println("Main thread waiting for others to finish")
	diesel.Wg.Wait()

	end := time.Now()
	elapsed := end.Sub(startTime)
	fmt.Println("Elapsed time :", elapsed.Nanoseconds())
	diesel.PrintMemory()
}
