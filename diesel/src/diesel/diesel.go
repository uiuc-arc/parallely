// +build !instrument

package diesel

import "math/rand"
import "fmt"
import "math"
import "sync"
import "os"
import "hash/crc32"
import "encoding/binary"

type ProbInterval struct {
	Reliability float32
  Delta float64
}

// import "encoding/json"

// import "time"

// Regular channels
var preciseChannelMapInt map[int] chan int
var preciseChannelMapInt32 map[int] chan int32
var preciseChannelMapInt64 map[int] chan int64

var preciseChannelMapFloat32 map[int] chan float32
var preciseChannelMapFloat64 map[int] chan float64

var approxChannelMapInt map[int] chan int
var approxChannelMapInt32 map[int] chan int32
var approxChannelMapInt64 map[int] chan int64

var approxChannelMapFloat32 map[int] chan float32
var approxChannelMapFloat64 map[int] chan float64

// Array channels
var preciseChannelMapIntArray map[int] chan []int
var preciseChannelMapInt32Array map[int] chan []int32
var preciseChannelMapInt64Array map[int] chan []int64

var preciseChannelMapFloat64Array map[int] chan []float64
var preciseChannelMapFloat32Array map[int] chan []float32

var approxChannelMapIntArray map[int] chan []int
var approxChannelMapInt32Array map[int] chan []int32
var approxChannelMapInt64Array map[int] chan []int64

var approxChannelMapFloat64Array map[int] chan []float64
var approxChannelMapFloat32Array map[int] chan []float32

var DynamicChannelMap map[int] chan ProbInterval
var ChecksumChannelMap map[int] chan uint32

var Wg sync.WaitGroup
var Numprocesses int

var debug int = 0

// dyn_rec_str = '''my_chan_index := {} * numprocesses + {};
// __temp_rec_val := <- dynamicChannelMap[my_chan_index];
// DynMap[parallely.DynKey{{Varname: \"{}\", Index: 0}}] = __temp_rec_val;'''

// Dynamic tracking map
type DynKey struct {
	Varname string	
  Index int
}

// var DynMap = map[DynKey] float64{}

func Max(a, b float32) float32 {
    if a > b {
        return a
    }
    return b
}

func LaunchThread(threadfunc func()){
	go threadfunc()
}

func LaunchThreadGroup(threadfunc func(tid int), numbers []int){
	for i := range numbers {
		go threadfunc(numbers[i])
	}
}

// func UpdateDynExpression(varname int, index int, expr_list []int, DynMap []ProbInterval){	
// 	// DynMap[DynKey{Varname: varname, Index: index}] = 0;
//   // fmt.Println(varname, index, DynMap, expr_list);
// 	sum_rel := 0.0
// 	for _, vname := range expr_list {
// 		sum_rel = sum_rel + DynMap[vname]
// 	}

// 	DynMap[varname] = Max(0.0, sum_rel - float64(len(expr_list) -1));
// }

func InitDynArray(varname int, size int, DynMap []ProbInterval){
	// fmt.Println("Initializing dynamic array: ", varname, size)
	for i:=0; i<size; i++ {
		DynMap[varname + i] = ProbInterval{1, 0}
	}
}

func SendDynFloat32Array(value []float32, sender, receiver int, DynMap []ProbInterval, start int) {
	my_chan_index := sender * Numprocesses + receiver
	for i:=0; i<len(value); i++ {
		preciseChannelMapFloat32[my_chan_index] <- value[i]
		DynamicChannelMap[my_chan_index] <- DynMap[start + i]
	}
}

func ReceiveDynFloat32Array(rec_var []float32, receiver, sender int, DynMap []ProbInterval, start int) {
	my_chan_index := sender * Numprocesses + receiver

	for i:=0; i<len(rec_var); i++ {
		rec_var[i] = <- preciseChannelMapFloat32[my_chan_index]
		DynMap[start + i] = <- DynamicChannelMap[my_chan_index];
	}

}

func SendDynIntArrayO1(value []int, sender, receiver int, DynMap []ProbInterval, start int) {
	my_chan_index := sender * Numprocesses + receiver
	var min float32 = DynMap[start].Reliability
	var maxd float64 = DynMap[start].Delta
	
	for i, _ := range value {
		preciseChannelMapInt[my_chan_index] <- value[i]

		// This looks wrong. Fix! prob have to get from the same element
		if min > DynMap[start + i].Reliability {
			min = DynMap[start + i].Reliability
		}
		if maxd < DynMap[start + i].Delta {
			maxd = DynMap[start + i].Delta
		}
	}

	DynamicChannelMap[my_chan_index] <- ProbInterval{min, maxd}
}

func ReceiveDynIntArrayO1(rec_var []int, receiver, sender int, DynMap []ProbInterval, start int) {
	my_chan_index := sender * Numprocesses + receiver

	for i:=0; i<len(rec_var); i++ {
		rec_var[i] = <- preciseChannelMapInt[my_chan_index]
	}
	__temp_rec_val := <- DynamicChannelMap[my_chan_index];
	for i:=0; i<len(rec_var); i++ {
		DynMap[start + i] = __temp_rec_val;
	}	
}

func SendDynFloat32ArrayO1(value []float32, sender, receiver int, DynMap []ProbInterval, start int) {
	my_chan_index := sender * Numprocesses + receiver
	var min float32 = DynMap[start].Reliability
	var maxd float64 = DynMap[start].Delta
	
	for i, _ := range value {
		preciseChannelMapFloat32[my_chan_index] <- value[i]

		// This looks wrong. Fix! prob have to get from the same element
		if min > DynMap[start + i].Reliability {
			min = DynMap[start + i].Reliability
		}
		if maxd < DynMap[start + i].Delta {
			maxd = DynMap[start + i].Delta
		}
	}

	DynamicChannelMap[my_chan_index] <- ProbInterval{min, maxd}
}

func ReceiveDynFloat32ArrayO1(rec_var []float32, receiver, sender int, DynMap []ProbInterval, start int) {
	my_chan_index := sender * Numprocesses + receiver

	for i:=0; i<len(rec_var); i++ {
		rec_var[i] = <- preciseChannelMapFloat32[my_chan_index]
	}
	__temp_rec_val := <- DynamicChannelMap[my_chan_index];
	for i:=0; i<len(rec_var); i++ {
		DynMap[start + i] = __temp_rec_val;
	}	
}

func SendDynFloat64Array(value []float64, sender, receiver int, DynMap []ProbInterval, start int) {
	my_chan_index := sender * Numprocesses + receiver

	for i:=0; i<len(value); i++ {
		preciseChannelMapFloat64[my_chan_index] <- value[i]
		DynamicChannelMap[my_chan_index] <- DynMap[start + i]
	}
}

func ReceiveDynFloat64Array(rec_var []float64, receiver, sender int, DynMap []ProbInterval, start int) {
	my_chan_index := sender * Numprocesses + receiver
	// temp_rec_val := <- preciseChannelMapFloat64Array[my_chan_index]
	// fmt.Println("Rec: ", len(rec_var))

	for i:=0; i<len(rec_var); i++ {
		rec_var[i] = <- preciseChannelMapFloat64[my_chan_index]
		DynMap[start + i] = <- DynamicChannelMap[my_chan_index];
		// __temp_rec_val := <- DynamicChannelMap[my_chan_index];
		// DynMap[start + i] = __temp_rec_val;
	}
}

func SendDynFloat64ArrayCustom(value []float64, sender, receiver int, trackedval ProbInterval, start int) {
	my_chan_index := sender * Numprocesses + receiver
	// temp_array := make([]float64, len(value))
	// copy(temp_array, value)

	// preciseChannelMapFloat64Array[my_chan_index] <- temp_array

	// var min float32 = DynMap[start].Reliability
	// var maxd float64 = DynMap[start].Delta
	for i, _ := range value {
		preciseChannelMapFloat64[my_chan_index] <- value[i]
		// This looks wrong. Fix! prob have to get from the same element
		// if min > DynMap[start + i].Reliability {
		// 	min = DynMap[start + i].Reliability
		// }
		// if maxd < DynMap[start + i].Delta {
		// 	maxd = DynMap[start + i].Delta
		// }
	}

	// for i, _ := range value {
	// }
	// for i:=0; i<len(value); i++ {
	// 	DynamicChannelMap[my_chan_index] <- DynMap[start + i]
	// }

	DynamicChannelMap[my_chan_index] <- trackedval
}

func SendDynFloat64ArrayO1(value []float64, sender, receiver int, DynMap []ProbInterval, start int) {
	my_chan_index := sender * Numprocesses + receiver
	// temp_array := make([]float64, len(value))
	// copy(temp_array, value)

	// preciseChannelMapFloat64Array[my_chan_index] <- temp_array

	var min float32 = DynMap[start].Reliability
	var maxd float64 = DynMap[start].Delta
	for i, _ := range value {
		preciseChannelMapFloat64[my_chan_index] <- value[i]
		// This looks wrong. Fix! prob have to get from the same element
		if min > DynMap[start + i].Reliability {
			min = DynMap[start + i].Reliability
		}
		if maxd < DynMap[start + i].Delta {
			maxd = DynMap[start + i].Delta
		}
	}

	// for i, _ := range value {
	// }
	// for i:=0; i<len(value); i++ {
	// 	DynamicChannelMap[my_chan_index] <- DynMap[start + i]
	// }

	DynamicChannelMap[my_chan_index] <- ProbInterval{min, maxd}
}

func ReceiveDynFloat64ArrayO1(rec_var []float64, receiver, sender int, DynMap []ProbInterval, start int) {
	my_chan_index := sender * Numprocesses + receiver
	// temp_rec_val := <- preciseChannelMapFloat64Array[my_chan_index]

	// fmt.Println(len(rec_var))

	for i:=0; i<len(rec_var); i++ {
		rec_var[i] = <- preciseChannelMapFloat64[my_chan_index]
	}
	__temp_rec_val := <- DynamicChannelMap[my_chan_index];
	for i:=0; i<len(rec_var); i++ {
		DynMap[start + i] = __temp_rec_val;
	}	
	// copy(rec_var, temp_rec_val)
}

func Min(array []float64) float64 {
    var min float64 = array[0]
    for _, value := range array {
        if min > value {
            min = value
        }
    }
    return min
}

func SendDynIntArray(value []int, sender, receiver int, DynMap []ProbInterval, start int) {
	my_chan_index := sender * Numprocesses + receiver
	for i := range(value) {
		preciseChannelMapInt[my_chan_index] <- value[i]
		DynamicChannelMap[my_chan_index] <- DynMap[start + i]
	}
}

func ReceiveDynIntArray(rec_var []int, receiver, sender int, DynMap []ProbInterval, start int) {
	my_chan_index := sender * Numprocesses + receiver
	for i:=0; i<len(rec_var); i++ {
		rec_var[i] = <- preciseChannelMapInt[my_chan_index]		
		DynMap[start + i] = <- DynamicChannelMap[my_chan_index];
	}	
	// if debug==1 {
	// 	fmt.Printf("%d Received message in precise float64 array chan : %d (%d * %d + %d)\n",
	// 		receiver, my_chan_index, sender, Numprocesses, receiver);
	// }
}

// func SendDynIntArrayO1(value []int, sender, receiver int, DynMap []float64, start int) {
// 	my_chan_index := sender * Numprocesses + receiver
// 	// temp_array := make([]int, len(value))
// 	// copy(temp_array, value)
// 	var min float64 = DynMap[start]
// 	for i, _ := range value {
// 		preciseChannelMapInt[my_chan_index] <- value[i]
// 		if min > DynMap[start + i] {
// 			min = DynMap[start + i]
// 		}
// 	}
// 	// fmt.Println(min)
// 	DynamicChannelMap[my_chan_index] <- min
// }

// func ReceiveDynIntArrayO1(rec_var []int, receiver, sender int, DynMap []float64, start int) {
// 	my_chan_index := sender * Numprocesses + receiver
// 	for i:= range(rec_var) {
// 		rec_var[i] = <- preciseChannelMapInt[my_chan_index]
// 	}
// 	__temp_dyn := <- DynamicChannelMap[my_chan_index];
// 	// fmt.Println(len(temp_rec_var))
// 	for i:= range(rec_var) {
// 		DynMap[start + i] = __temp_dyn;
// 	}
// }

func CopyDynArray(array1 int, array2 int, size int, DynMap []ProbInterval) bool {
	for i:=0; i<size; i++ {
		DynMap[array1 + size] = DynMap[array2 + size]
	}
	return true
}


func CheckArray(start int, limit float32, size int, DynMap []ProbInterval) bool {
	failed := true
	for i:=start; i<size; i++ {
		if failed && (DynMap[i].Reliability < limit) {
			fmt.Println("Verification failed due to reliability of: ", DynMap[i])
			failed = false
		}
	}
	return failed
}

func CheckFloat64(val float64, PI ProbInterval, epsThresh float32, deltaThresh float64) (result bool) {
	result = (PI.Reliability < epsThresh && PI.Delta < deltaThresh)
	//fmt.Println(result)
	return
}

func DumpDynMap(DynMap []ProbInterval, filename string) {
	f, _ := os.Create(filename)
	defer f.Close()

	// v := make([]float64, 0, len(DynMap))

	// for  _, value := range DynMap {
	// 	v = append(v, value)
	// }

	// jsonString, _ := json.Marshal(DynMap)
	
	f.WriteString(fmt.Sprintln(DynMap))
}

func GetCastingError64to32(original float64, casted float32) float64{
	recasted := float64(casted)
	// fmt.Println(recasted - original)
	return math.Abs(recasted - original)
}

func Cast64to32Array(array32 []float32, array64 []float64){
	for i, f64 := range array64 {
        array32[i] = float32(f64)
	}
}

func Cast32to64Array(array32 []float32, array64 []float64){
	for i, f32 := range array32 {
        array64[i] = float64(f32)
	}
}

func Randchoice(prob float32, option1, option2 int) int {
	failure := rand.Float32()
	if failure < prob {
		// fmt.Println("Pass", failure, option1)
		return option1
	} else {
		// fmt.Println("Fail", failure, option2)
		return option2
	}
}

func RandchoiceFloat64(prob float32, option1, option2 float64) float64 {
	return option1
	// failure := rand.Float32()
	// if failure < prob {
	// 	// fmt.Println("Pass", failure, option1)
	// 	return option1
	// } else {
	// 	// fmt.Println("Fail", failure, option2)
	// 	return option2
	// }
}

func RandchoiceFlag(prob float32, option1, option2 int, flag *bool) int {
	failure := rand.Float32()
	if failure < prob {
		// fmt.Println("Pass", failure, option1)
		return option1
	} else {
		// fmt.Println("Fail", failure, option2)
		*flag = true
		return option2
	}
}

func RandchoiceFlagFloat64(prob float32, option1, option2 float64, flag *bool) float64 {
	failure := rand.Float32()
	if failure < prob {
		// fmt.Println("Pass", failure, option1)
		return option1
	} else {
		// fmt.Println("Fail", failure, option2)
		*flag = true
		return option2
	}
}

func InitChannels(numprocesses_in int){
	// var temp_approxChannelMap map[int] chan int
	Numprocesses = numprocesses_in 
	Wg.Add(numprocesses_in)
	if debug==1 {
		fmt.Printf("Starting the wait group for %d threads", numprocesses_in)
	}
	

	preciseChannelMapInt = make(map[int] chan int)
	preciseChannelMapInt32 = make(map[int] chan int32)
	preciseChannelMapInt64 = make(map[int] chan int64)

	preciseChannelMapFloat32 = make(map[int] chan float32)
	preciseChannelMapFloat64 = make(map[int] chan float64)

	approxChannelMapInt = make(map[int] chan int)
	approxChannelMapInt32 = make(map[int] chan int32)
	approxChannelMapInt64 = make(map[int] chan int64)

	approxChannelMapFloat32 = make(map[int] chan float32)
	approxChannelMapFloat64 = make(map[int] chan float64)

	preciseChannelMapIntArray = make(map[int] chan []int)
	preciseChannelMapInt32Array = make(map[int] chan []int32)
	preciseChannelMapInt64Array = make(map[int] chan []int64)

	preciseChannelMapFloat64Array = make(map[int] chan []float64)
	preciseChannelMapFloat32Array = make(map[int] chan []float32)

	approxChannelMapIntArray = make(map[int] chan []int)
	approxChannelMapInt32Array = make(map[int] chan []int32)
	approxChannelMapInt64Array = make(map[int] chan []int64)

	approxChannelMapFloat64Array = make(map[int] chan []float64)
	approxChannelMapFloat32Array = make(map[int] chan []float32)
	
	DynamicChannelMap = make(map[int] chan ProbInterval)
	ChecksumChannelMap = make(map[int] chan uint32)

	buffer_size := 1000000
	
	for i := 0; i < numprocesses_in * numprocesses_in; i++ {
		preciseChannelMapInt[i] = make(chan int, buffer_size)
		preciseChannelMapInt32[i] = make(chan int32, buffer_size)
		preciseChannelMapInt64[i] = make(chan int64, buffer_size)
		preciseChannelMapIntArray[i] = make(chan []int, buffer_size)
		preciseChannelMapInt32Array[i] = make(chan []int32, buffer_size)
		preciseChannelMapInt64Array[i] = make(chan []int64, buffer_size)

		preciseChannelMapFloat64[i] = make(chan float64, buffer_size)
		preciseChannelMapFloat32[i] = make(chan float32, buffer_size)
		preciseChannelMapFloat64Array[i] = make(chan []float64, buffer_size)
		preciseChannelMapFloat32Array[i] = make(chan []float32, buffer_size)

		approxChannelMapInt[i] = make(chan int)
		approxChannelMapInt32[i] = make(chan int32)
		approxChannelMapInt64[i] = make(chan int64)
		approxChannelMapIntArray[i] = make(chan []int)
		approxChannelMapInt32Array[i] = make(chan []int32)
		approxChannelMapInt64Array[i] = make(chan []int64)

		approxChannelMapFloat64[i] = make(chan float64)
		approxChannelMapFloat32[i] = make(chan float32)
		approxChannelMapFloat64Array[i] = make(chan []float64)
		approxChannelMapFloat32Array[i] = make(chan []float32)

		DynamicChannelMap[i] = make(chan ProbInterval, buffer_size)
		ChecksumChannelMap[i] = make(chan uint32, buffer_size)
	}

	if debug==1 {
		fmt.Println("Initialized Channels : ", len(approxChannelMapInt) * 20)
	}
}

func SendChkFloat64Array(value []float64, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	for i:=0; i<len(value); i++ {
		preciseChannelMapFloat64[my_chan_index] <- value[i]
		bs := make([]byte, 8)
		binary.LittleEndian.PutUint64(bs, math.Float64bits(value[i]))
		ChecksumChannelMap[my_chan_index] <- crc32.ChecksumIEEE(bs)
	}
}

func SendSChkFloat64Array(value []float64, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	bs := make([]byte, 8*len(value))
	for i:=0; i<len(value); i++ {
		preciseChannelMapFloat64[my_chan_index] <- value[i]
		binary.LittleEndian.PutUint64(bs[8*i:], math.Float64bits(value[i]))
	}
	ChecksumChannelMap[my_chan_index] <- crc32.ChecksumIEEE(bs)
}

func ReceiveChkFloat64Array(rec_var []float64, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver
	for i:=0; i<len(rec_var); i++ {
		rec_var[i] = <- preciseChannelMapFloat64[my_chan_index]
		bs := make([]byte, 8)
		binary.LittleEndian.PutUint64(bs, math.Float64bits(rec_var[i]))
		checksumD := crc32.ChecksumIEEE(bs)
		checksumR := <- ChecksumChannelMap[my_chan_index]
		if checksumD != checksumR { fmt.Println("Checksum mismatch!") }
	}
}

func ReceiveSChkFloat64Array(rec_var []float64, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver
	bs := make([]byte, 8*len(rec_var))
	for i:=0; i<len(rec_var); i++ {
		rec_var[i] = <- preciseChannelMapFloat64[my_chan_index]
		binary.LittleEndian.PutUint64(bs[8*i:], math.Float64bits(rec_var[i]))
	}
	checksumD := crc32.ChecksumIEEE(bs)
	checksumR := <- ChecksumChannelMap[my_chan_index]
	if checksumD != checksumR { fmt.Println("Checksum mismatch!") }
}

func SendChkIntArray(value []int, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	for i:=0; i<len(value); i++ {
		preciseChannelMapInt[my_chan_index] <- value[i]
		bs := make([]byte, 8)
		binary.LittleEndian.PutUint64(bs, uint64(value[i]))
		ChecksumChannelMap[my_chan_index] <- crc32.ChecksumIEEE(bs)
	}
}

func SendSChkIntArray(value []int, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	bs := make([]byte, 8*len(value))
	for i:=0; i<len(value); i++ {
		preciseChannelMapInt[my_chan_index] <- value[i]
		binary.LittleEndian.PutUint64(bs[8*i:], uint64(value[i]))
	}
	ChecksumChannelMap[my_chan_index] <- crc32.ChecksumIEEE(bs)
}

func ReceiveChkIntArray(rec_var []int, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver
	for i:=0; i<len(rec_var); i++ {
		rec_var[i] = <- preciseChannelMapInt[my_chan_index]
		bs := make([]byte, 8)
		binary.LittleEndian.PutUint64(bs, uint64(rec_var[i]))
		checksumD := crc32.ChecksumIEEE(bs)
		checksumR := <- ChecksumChannelMap[my_chan_index]
		if checksumD != checksumR { fmt.Println("Checksum mismatch!") }
	}
}

func ReceiveSChkIntArray(rec_var []int, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver
	bs := make([]byte, 8*len(rec_var))
	for i:=0; i<len(rec_var); i++ {
		rec_var[i] = <- preciseChannelMapInt[my_chan_index]
		binary.LittleEndian.PutUint64(bs[8*i:], uint64(rec_var[i]))
	}
	checksumD := crc32.ChecksumIEEE(bs)
	checksumR := <- ChecksumChannelMap[my_chan_index]
	if checksumD != checksumR { fmt.Println("Checksum mismatch!") }
}

func SendDynVal(value ProbInterval, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	DynamicChannelMap[my_chan_index] <- value
	if debug==1 {
		fmt.Printf("%d Sending message in precise int chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, Numprocesses, receiver);
	}
}



func SendInt(value, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	preciseChannelMapInt[my_chan_index] <- value
	if debug==1 {
		fmt.Printf("%d Sending message in precise int chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, Numprocesses, receiver);
	}
}

func SendInt32(value int32, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	preciseChannelMapInt32[my_chan_index] <- value
	if debug==1 {
		fmt.Printf("%d Sending message in precise int32 chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, Numprocesses, receiver);
	}
}

func SendInt64(value int64, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	preciseChannelMapInt64[my_chan_index] <- value
	if debug==1 {
		fmt.Printf("%d Sending message in precise int64 chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, Numprocesses, receiver);
	}
}

func SendFloat32(value float32, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	preciseChannelMapFloat32[my_chan_index] <- value
	if debug==1 {
		fmt.Printf("%d Sending message in precise float32 chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, Numprocesses, receiver);
	}
}

func SendFloat64(value float64, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	preciseChannelMapFloat64[my_chan_index] <- value
	// if debug==1 {
	// 	fmt.Printf("%d Sending message in precise float64 chan : %d (%d * %d + %d)\n",
	// 		sender, my_chan_index, sender, Numprocesses, receiver);
	// }
}

func SendApprox(value, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	approxChannelMapInt[my_chan_index] <- value
	if debug==1 {
		fmt.Printf("%d Sending message in approx int chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, Numprocesses, receiver);
	}
}

func SendInt32Approx(value int32, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	approxChannelMapInt32[my_chan_index] <- value
	if debug==1 {
		fmt.Printf("%d Sending message in approx int32 chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, Numprocesses, receiver);
	}
}

func SendInt64Approx(value int64, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	approxChannelMapInt64[my_chan_index] <- value
	if debug==1 {
		fmt.Printf("%d Sending message in approx int64 chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, Numprocesses, receiver);
	}
}

func SendIntArrayApprox(value []int, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	temp_array := make([]int, len(value))
	copy(temp_array, value)
	approxChannelMapIntArray[my_chan_index] <- temp_array
	// if debug==1 {
	// 	fmt.Printf("%d Sending message in approx int array chan : %d (%d * %d + %d)\n",
	// 		sender, my_chan_index, sender, Numprocesses, receiver);
	// }
}

func SendInt32Array(value []int32, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	temp_array := make([]int32, len(value))
	copy(temp_array, value)
	preciseChannelMapInt32Array[my_chan_index] <- temp_array
	if debug==1 {
		fmt.Printf("%d Sending message in precise float32 array chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, Numprocesses, receiver);
	}
}

func SendInt32ArrayApprox(value []int32, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	temp_array := make([]int32, len(value))
	copy(temp_array, value)
	approxChannelMapInt32Array[my_chan_index] <- temp_array
	if debug==1 {
		fmt.Printf("%d Sending message in precise float32 array chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, Numprocesses, receiver);
	}
}

func SendFloat64Array(value []float64, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	// temp_array := make([]float64, len(value))
	// copy(temp_array, value)
	// preciseChannelMapFloat64Array[my_chan_index] <- temp_array

	for i := range(value){
		preciseChannelMapFloat64[my_chan_index] <- value[i]
	}

}

func ReceiveFloat64Array(rec_var []float64, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver
	// temp_rec_val := <- preciseChannelMapFloat64Array[my_chan_index]

	for i := range(rec_var) {
		rec_var[i] = <- preciseChannelMapFloat64[my_chan_index]
	}
	// if len(rec_var) != len(temp_rec_val) {
	// 	rec_var = make([]float64, len(temp_rec_val))
	// }

	// copy(rec_var, temp_rec_val)
}

// func SendDynFloat64Array(value []float64, sender, receiver int, DynMap []float64, start int) {
// 	my_chan_index := sender * Numprocesses + receiver
// 	temp_array := make([]float64, len(value))
// 	copy(temp_array, value)
// 	preciseChannelMapFloat64Array[my_chan_index] <- temp_array

// 	for i:=0; i<len(value); i++ {
// 		DynamicChannelMap[my_chan_index] <- DynMap[start + i]
// 	}
	
// 	if debug==1 {
// 		fmt.Printf("%d Sending message in precise float64 array chan : %d (%d * %d + %d)\n",
// 			sender, my_chan_index, sender, Numprocesses, receiver);
// 	}
// }

func SendFloat32Array(value []float32, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	temp_array := make([]float32, len(value))
	copy(temp_array, value)
	preciseChannelMapFloat32Array[my_chan_index] <- temp_array
	if debug==1 {
		fmt.Printf("%d Sending message in precise float32 array chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, Numprocesses, receiver);
	}
}


func ReceiveDynVal(rec_var *ProbInterval, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver
	temp_rec_val := <- DynamicChannelMap[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, Numprocesses, receiver);
	}
	*rec_var = temp_rec_val
}



func ReceiveInt(rec_var *int, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver
	temp_rec_val := <- preciseChannelMapInt[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, Numprocesses, receiver);
	}
	*rec_var = temp_rec_val
}

func ReceiveInt32(rec_var *int32, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver
	temp_rec_val := <- preciseChannelMapInt32[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, Numprocesses, receiver);
	}
	*rec_var = temp_rec_val
}

func ReceiveInt64(rec_var *int64, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver
	temp_rec_val := <- preciseChannelMapInt64[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, Numprocesses, receiver);
	}
	*rec_var = temp_rec_val
}

func ReceiveFloat32(rec_var *float32, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver
	temp_rec_val := <- preciseChannelMapFloat32[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise float32 chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, Numprocesses, receiver);
	}
	*rec_var = temp_rec_val
}

func ReceiveFloat64(rec_var *float64, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver
	temp_rec_val := <- preciseChannelMapFloat64[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise float64 chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, Numprocesses, receiver);
	}
	*rec_var = temp_rec_val
}

func ReceiveIntApprox(rec_var *int, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver
	temp_rec_val := <- approxChannelMapInt[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, Numprocesses, receiver);
	}
	*rec_var = temp_rec_val
}

func ReceiveInt32Approx(rec_var *int32, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver
	temp_rec_val := <- approxChannelMapInt32[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, Numprocesses, receiver);
	}
	*rec_var = temp_rec_val
}

func ReceiveInt64Approx(rec_var *int64, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver
	temp_rec_val := <- approxChannelMapInt64[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, Numprocesses, receiver);
	}
	*rec_var = temp_rec_val
}

func SendIntArray(value []int, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	// temp_array := make([]int, len(value))
	// copy(temp_array, value)
	for i := range(value){
		preciseChannelMapInt[my_chan_index] <- value[i]
	}
	// if debug==1 {
	// 	fmt.Printf("%d Sending message in precise float32 array chan : %d (%d * %d + %d)\n",
	// 		sender, my_chan_index, sender, Numprocesses, receiver);
	// }
}

func ReceiveIntArray(rec_var []int, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver
	// temp_rec_var := <- preciseChannelMapIntArray[my_chan_index]
	// copy(rec_var, temp_rec_var)
	for i := range(rec_var) {	
		rec_var[i] = <- preciseChannelMapInt[my_chan_index]
	}
	// if debug==1 {
	// 	fmt.Printf("%d Received message in precise int array chan : %d (%d * %d + %d)\n",
	// 		receiver, my_chan_index, sender, Numprocesses, receiver);
	// }
	// fmt.Println(len(rec_var), len(temp_rec_val))
	// if len(rec_var) != len(temp_rec_val) {
	// 	rec_var = make([]int, len(temp_rec_val))
	// 	fmt.Println("=======>", len(rec_var), len(temp_rec_val))
	// }	
	// copy(rec_var, temp_rec_val)
}

func ReceiveInt32Array(rec_var []int32, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver
	temp_rec_val := <- preciseChannelMapInt32Array[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise int32 array chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, Numprocesses, receiver);
	}
	// if len(rec_var) != len(temp_rec_val) {
	// 	rec_var = make([]int32, len(temp_rec_val))
	// }
	copy(rec_var, temp_rec_val)
}

// func ReceiveFloat64Array(rec_var []float64, receiver, sender int) {
// 	my_chan_index := sender * Numprocesses + receiver
// 	// temp_rec_val := <- preciseChannelMapFloat64Array[my_chan_index]

// 	for i := range(rec_var) {
// 		rec_var[i] = <- preciseChannelMapFloat64[my_chan_index]
// 	}
// 	// if len(rec_var) != len(temp_rec_val) {
// 	// 	rec_var = make([]float64, len(temp_rec_val))
// 	// }
// 	// copy(rec_var, temp_rec_val)
// }

func ReceiveFloat32Array(rec_var []float32, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver
	temp_rec_val := <- preciseChannelMapFloat32Array[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise float32 array chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, Numprocesses, receiver);
	}
	// if len(rec_var) != len(temp_rec_val) {
	// 	rec_var = make([]float32, len(temp_rec_val))
	// }
	copy(rec_var, temp_rec_val)
}

func Condsend(cond, value, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	if debug==1 {
		fmt.Printf("%d Sending message in approx int chan : %d (%d * %d + %d)\n", sender,
			my_chan_index, sender, Numprocesses, receiver);
	}
	if cond != 0 {
		approxChannelMapInt[my_chan_index] <- value
	} else {
		if debug==1 {
			fmt.Printf("[Failure %d] %d Sending message in approx int chan : %d (%d * %d + %d)\n", cond, sender,
				my_chan_index, sender, Numprocesses, receiver);
		}
		approxChannelMapInt[my_chan_index] <- -1
	}
}

func CondsendIntArray(cond int, value []int, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	if debug==1 {
		fmt.Printf("%d Sending message in approx int chan : %d (%d * %d + %d)\n", sender,
			my_chan_index, sender, Numprocesses, receiver);
	}
	if cond != 0 {
		temp_array := make([]int, len(value))
		copy(temp_array, value)
		approxChannelMapIntArray[my_chan_index] <- temp_array
	} else {
		if debug==1 {
			fmt.Printf("[Failure %d] %d Cond Sending message in approx int chan : %d (%d * %d + %d)\n", cond, sender,
				my_chan_index, sender, Numprocesses, receiver);
		}
		approxChannelMapIntArray[my_chan_index] <- []int{}
	}
}

func CondsendInt32(cond, value int32, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	if debug==1 {
		fmt.Printf("%d Sending message in approx int32 chan : %d (%d * %d + %d)\n", sender,
			my_chan_index, sender, Numprocesses, receiver);
	}
	if cond != 0 {
		approxChannelMapInt32[my_chan_index] <- value
	} else {
		if debug==1 {
			fmt.Printf("[Failure %d] %d Sending message in approx int32 chan : %d (%d * %d + %d)\n", cond, sender,
				my_chan_index, sender, Numprocesses, receiver);
		}
		approxChannelMapInt32[my_chan_index] <- -1
	}
}

func CondsendInt64(cond, value int64, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	if debug==1 {
		fmt.Printf("%d Sending message in approx int64 chan : %d (%d * %d + %d)\n", sender,
			my_chan_index, sender, Numprocesses, receiver);
	}
	if cond != 0 {
		approxChannelMapInt64[my_chan_index] <- value
	} else {
		if debug==1 {
			fmt.Printf("[Failure %d] %d Sending message in approx int64 chan : %d (%d * %d + %d)\n", cond, sender,
				my_chan_index, sender, Numprocesses, receiver);
		}
		approxChannelMapInt64[my_chan_index] <- -1
	}
}

func CondsendFloat32(cond, value float32, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	if debug==1 {
		fmt.Printf("%d Sending message in approx chan : %d (%d * %d + %d)\n", sender,
			my_chan_index, sender, Numprocesses, receiver);
	}
	if cond != 0 {
		approxChannelMapFloat32[my_chan_index] <- value
	} else {
		if debug==1 {
			fmt.Printf("[Failure %d] %d Sending message in approx chan : %d (%d * %d + %d)\n", cond, sender,
				my_chan_index, sender, Numprocesses, receiver);
		}
		approxChannelMapFloat32[my_chan_index] <- -1
	}
}

func CondsendFloat64(cond, value float64, sender, receiver int) {
	my_chan_index := sender * Numprocesses + receiver
	if debug==1 {
		fmt.Printf("%d Sending message in approx chan : %d (%d * %d + %d)\n", sender,
			my_chan_index, sender, Numprocesses, receiver);
	}
	if cond != 0 {
		approxChannelMapFloat64[my_chan_index] <- value
	} else {
		if debug==1 {
			fmt.Printf("[Failure %d] %d Sending message in approx chan : %d (%d * %d + %d)\n", cond, sender,
				my_chan_index, sender, Numprocesses, receiver);
		}
		approxChannelMapFloat64[my_chan_index] <- -1
	}
}

func Condreceive(rec_cond_var, rec_var *int, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver

	if debug==1 {
		fmt.Printf("---- %d Waiting to Receive from approx int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, Numprocesses, receiver);
	}
	
	temp_rec_val := <- approxChannelMapInt[my_chan_index]

	if debug==1 {
		fmt.Printf("%d Recieved message in approx chan : %d (%d)\n", receiver, my_chan_index,
			temp_rec_val);
	}
	
	if temp_rec_val != -1 {
		*rec_var = temp_rec_val
		*rec_cond_var = 1
	} else {
		*rec_cond_var = 0
	}
}

func CondreceiveIntArray(rec_cond_var *int, rec_var []int, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver

	if debug==1 {
		fmt.Printf("---- %d Waiting to Receive from approx int array chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, Numprocesses, receiver);
	}
	
	temp_rec_val := <- approxChannelMapIntArray[my_chan_index]
	
	if len(temp_rec_val) != 0 {
		rec_var = temp_rec_val
		*rec_cond_var = 1
	} else {
		// if debug==0 {
		// 	fmt.Printf("%d Recieved failed message in approx chan : %d (%d, %d)\n", receiver, my_chan_index,
		// 		temp_rec_val, len(temp_rec_val));
		// }
		*rec_cond_var = 0
	}
}

func CondreceiveInt32(rec_cond_var, rec_var *int32, receiver, sender int) {
	my_chan_index := sender * Numprocesses + receiver

	if debug==1 {
		fmt.Printf("---- %d Waiting to Receive from approx int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, Numprocesses, receiver);
	}
	
	temp_rec_val := <- approxChannelMapInt32[my_chan_index]

	if debug==1 {
		fmt.Printf("%d Recieved message in approx chan : %d (%d)\n", receiver, my_chan_index,
			temp_rec_val);
	}
	
	if temp_rec_val != -1 {
		*rec_var = temp_rec_val
		*rec_cond_var = 1
	} else {
		*rec_cond_var = 0
	}
}

func Hoeffding(n int, delta float64) (eps float64) {
	eps = math.Sqrt((0.6*math.Log((math.Log(float64(1.1*float64(n+1)))/math.Log(1.10)))+0.555*math.Log(24/delta))/float64(n+1))
	return
}


func FuseFloat64(arr []float64, dynMap []ProbInterval)(mean float64 , newInterval ProbInterval){
	var ns [] int
	var totalN int = 0
	var sum float64 = 0	

	//var mean float64
	for i:=0;i<len(dynMap);i++{
		ns = append(ns,ComputeN(dynMap[i]))
		totalN = totalN + ns[i]
		sum = sum + (arr[i]*float64(ns[i]))
	}
		
	mean = sum / float64(totalN)
	var eps float64 = Hoeffding(totalN,dynMap[0].Delta)
	newInterval.Reliability = float32(eps)
	newInterval.Delta = dynMap[0].Delta
	return
}


func ComputeN(ui ProbInterval)(n int){
	var eps float64 = float64(ui.Reliability)
	var delta float64 = ui.Delta
	n = int(0.5*(1/(eps*eps))*math.Log((2/(1-delta))))
	return n

}

func AddProbInterval(val1, val2 float64, fst, snd ProbInterval)(retval float64, out ProbInterval){
	out.Reliability = fst.Reliability + snd.Reliability
	out.Delta = fst.Delta + snd.Delta
	retval = val1+val2
	return
} 

func MulProbInterval(val1, val2 float64, fst, snd ProbInterval)(retval float64, out ProbInterval){
	retval = val1 * val2
	out.Reliability = (float32(math.Abs(val1)) * snd.Reliability) + (float32(math.Abs(val2)) * fst.Reliability) + (fst.Reliability * snd.Reliability)
	out.Delta = fst.Delta + snd.Delta
	return
} 

func DivProbInterval(val1, val2 float64, fst, snd ProbInterval)(retval float64, out ProbInterval){
	retval = val1 / val2
	out.Reliability = (float32(math.Abs(val1)) * snd.Reliability) + (float32(math.Abs(val2)) * fst.Reliability) / (float32(math.Abs(val2)) * (float32(math.Abs(val2)) - snd.Reliability))
	out.Delta = fst.Delta + snd.Delta
	return
}


func ConvBool(x bool) int {
	if x {
		return 1;
	} else {
		return 0;
	}
}

func PrintMemory() {
	fmt.Println("Memory not Instrumented")
}


//Sasa's proposed addition to the runtime: add a custome class for tracking means of Boolean Indicator Random Vars
type BooleanTracker struct {
	successes int 
	totalSamples int 
	mean float64 
	//delta float64 
	//eps float64 
	meanProbInt ProbInterval
}


func NewBooleanTracker() (b BooleanTracker) {
	b.successes = 0
	b.totalSamples  = 0
	b.mean  = 0
	b.meanProbInt.Delta = 1
	b.meanProbInt.Reliability = 1
	return
}

func (b *BooleanTracker) SetDelta(d float64)  {
	b.meanProbInt.Delta = d
}

func (b *BooleanTracker) GetMean() float64 {
	return b.mean
}

func (b *BooleanTracker) GetInterval() (res ProbInterval) {
	return b.meanProbInt
}



func (b *BooleanTracker) AddSample(samp int) {
    b.successes = b.successes + samp
    b.totalSamples = b.totalSamples + 1
    b.Hoeffding()
    b.ComputeMean()
}

func (b *BooleanTracker) Hoeffding() {
	b.meanProbInt.Reliability = float32(math.Sqrt((0.6*math.Log((math.Log(float64(1.1*float64(b.totalSamples+1)))/math.Log(1.10)))+0.555*math.Log(24/b.meanProbInt.Delta))/float64(b.totalSamples+1)))
}

func (b *BooleanTracker) ComputeMean(){
	b.mean = float64(b.successes)/float64(b.totalSamples)
}


//func (b *BooleanTracker) Check(c float32) bool{
//	CheckFloat64(val float64, PI ProbInterval, epsThresh float32, deltaThresh float64)
//}

func FuseBooleanTrackers(arr [] BooleanTracker) (res BooleanTracker){
	res = NewBooleanTracker()
	for i:=0; i < len(arr); i++ {
		res.totalSamples = res.totalSamples + arr[i].totalSamples
		res.successes =  res.successes + arr[i].successes
	}

	res.Hoeffding()
	res.GetMean()
	return

}


func FuseFloat64IntoBooleanTracker(arr []float64, dynMap []ProbInterval)(res BooleanTracker){

	
	var ns [] int
	var totalN int = 0
	var sum float64 = 0	

	//var mean float64
	for i:=0;i<len(dynMap);i++{
		ns = append(ns,ComputeN(dynMap[i]))
		totalN = totalN + ns[i]
		sum = sum + (arr[i]*float64(ns[i]))
	}
		

	res = NewBooleanTracker()
	res.successes = int(sum)
	res.totalSamples = totalN
	res.Hoeffding()
	res.ComputeMean()
	return

}





