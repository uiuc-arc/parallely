// +build !instrument

package dieseldistrel

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

// type ProbInterval struct {
// 	Reliability float32
// }

var noiselevel float32 = 0.9999999

var types = [...]string{"int"}

var ch *amqp.Channel
var conn *amqp.Connection

var pingchannel amqp.Queue

// Regular channels
var preciseChannelMapInt map[int]amqp.Queue
var preciseChannelMapInt32 map[int]amqp.Queue
var preciseChannelMapInt64 map[int]amqp.Queue
var preciseChannelMapFloat64 map[int]amqp.Queue
var preciseChannelMapFloat32 map[int]amqp.Queue

var approxChannelMapInt map[int]amqp.Queue
var approxChannelMapInt32 map[int]amqp.Queue
var approxChannelMapInt64 map[int]amqp.Queue
var approxChannelMapFloat32 map[int]amqp.Queue
var approxChannelMapFloat64 map[int]amqp.Queue

// Array channels
var preciseChannelMapIntArray map[int]amqp.Queue
var preciseChannelMapInt32Array map[int]amqp.Queue
var preciseChannelMapInt64Array map[int]amqp.Queue
var preciseChannelMapFloat64Array map[int]amqp.Queue
var preciseChannelMapFloat32Array map[int]amqp.Queue

var approxChannelMapIntArray map[int]amqp.Queue
var approxChannelMapInt32Array map[int]amqp.Queue
var approxChannelMapInt64Array map[int]amqp.Queue
var approxChannelMapFloat64Array map[int]amqp.Queue
var approxChannelMapFloat32Array map[int]amqp.Queue

var DynamicChannelMap map[int]amqp.Queue
var DynamicChannelMapArray map[int]amqp.Queue

var syncChannelMap map[int]amqp.Queue

var Wg sync.WaitGroup
var Numprocesses int

var debug int = 0

// dyn_rec_str = '''my_chan_index := {} * numprocesses + {};
// __temp_rec_val := <- dynamicChannelMap[my_chan_index];
// DynMap[parallely.DynKey{{Varname: \"{}\", Index: 0}}] = __temp_rec_val;'''

// Dynamic tracking map
type DynKey struct {
	Varname string
	Index   int
}

// var DynMap = map[DynKey] float64{}
var startTime time.Time

func StartTiming() {
	startTime = time.Now()
}

func EndTiming() {
	elapsed := time.Since(startTime)
	fmt.Println("Elapsed time : ", elapsed.Nanoseconds())
}

func ConvBool(x bool) int {
	if x {
		return 1
	} else {
		return 0
	}
}

func Max64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func Min64(a, b float64) float64 {
	if a > b {
		return b
	}
	return a
}

func Abs32(a float32) float32 {
	if a > 0 {
		return a
	}
	return a * -1
}

func AbsInt(a int) int {
	if a > 0 {
		return a
	}
	return a * -1
}

func Max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func Min32(a, b float32) float32 {
	if a > b {
		return b
	}
	return a
}

func intToByte(i int) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(i))
	return bs
}

func intArrayToByte(intarray []int) []byte {
	bs := make([]byte, 8*len(intarray))

	for i, elem := range intarray {
		binary.LittleEndian.PutUint64(bs[8*i:], uint64(elem))
	}
	return bs
}

func intfrombytes(bytes []byte) int {
	return int(int64(binary.LittleEndian.Uint64(bytes)))
}

// func DynCondGeqInt(lvar, rvar int,
// 	DynMap []ProbInterval,
// 	lvarindex, rvarindex int,
// 	op1, op2 int,
// 	op1index, op2index, assignedindex int) int {
// 	condition_rel := DynMap[lvarindex].Reliability + DynMap[rvarindex].Reliability - float32(1.0)
// 	if float64(lvar)-DynMap[lvarindex].Delta >= float64(rvar)+DynMap[lvarindex].Delta {
// 		DynMap[assignedindex] = DynMap[op1index]
// 		DynMap[assignedindex].Reliability *= condition_rel
// 		return op1
// 	} else if float64(lvar)+DynMap[lvarindex].Delta >= float64(rvar)-DynMap[lvarindex].Delta {
// 		DynMap[assignedindex] = DynMap[op2index]
// 		DynMap[assignedindex].Reliability *= condition_rel
// 		return op2
// 	} else {
// 		DynMap[assignedindex].Delta = float64(AbsInt(op1-op2)) * Max64(DynMap[op1index].Delta, DynMap[op2index].Delta)
// 		DynMap[assignedindex].Reliability = Min32(DynMap[op1index].Reliability, DynMap[op2index].Reliability) * condition_rel
// 		if lvar >= rvar {
// 			return op1
// 		} else {
// 			return op2
// 		}
// 	}
// }

// func DynCondFloat32GeqInt(lvar, rvar float32,
// 	DynMap []ProbInterval,
// 	lvarindex, rvarindex int,
// 	op1, op2 int,
// 	op1index, op2index, assignedindex int) int {

// 	f64_1 := float64(lvar)
// 	f64_2 := float64(rvar)
// 	condition_rel := DynMap[lvarindex].Reliability + DynMap[rvarindex].Reliability - float32(1.0)
// 	if f64_1-DynMap[lvarindex].Delta >= f64_2+DynMap[rvarindex].Delta {
// 		DynMap[assignedindex] = DynMap[op1index]
// 		DynMap[assignedindex].Reliability *= condition_rel
// 		return op1
// 	} else if f64_1+DynMap[lvarindex].Delta < f64_2-DynMap[rvarindex].Delta {
// 		DynMap[assignedindex] = DynMap[op2index]
// 		DynMap[assignedindex].Reliability *= condition_rel
// 		return op2
// 	} else {
// 		DynMap[assignedindex].Delta = float64(AbsInt(op1-op2)) * Max64(DynMap[op1index].Delta, DynMap[op2index].Delta)
// 		DynMap[assignedindex].Reliability = Min32(DynMap[op1index].Reliability, DynMap[op2index].Reliability) * condition_rel
// 		if lvar >= rvar {
// 			return op1
// 		} else {
// 			return op2
// 		}
// 	}
// }

// func DynCondFloat64GeqInt(lvar, rvar float64,
// 	DynMap []ProbInterval,
// 	lvarindex, rvarindex int,
// 	op1, op2 int,
// 	op1index, op2index, assignedindex int) int {

// 	f64_1 := lvar
// 	f64_2 := rvar
// 	condition_rel := DynMap[lvarindex].Reliability + DynMap[rvarindex].Reliability - float32(1.0)
// 	if f64_1-DynMap[lvarindex].Delta >= f64_2+DynMap[rvarindex].Delta {
// 		DynMap[assignedindex] = DynMap[op1index]
// 		DynMap[assignedindex].Reliability *= condition_rel
// 		return op1
// 	} else if f64_1+DynMap[lvarindex].Delta < f64_2-DynMap[rvarindex].Delta {
// 		DynMap[assignedindex] = DynMap[op2index]
// 		DynMap[assignedindex].Reliability *= condition_rel
// 		return op2
// 	} else {
// 		DynMap[assignedindex].Delta = float64(AbsInt(op1-op2)) * Max64(DynMap[op1index].Delta, DynMap[op2index].Delta)
// 		DynMap[assignedindex].Reliability = Min32(DynMap[op1index].Reliability, DynMap[op2index].Reliability) * condition_rel
// 		if lvar >= rvar {
// 			return op1
// 		} else {
// 			return op2
// 		}
// 	}
// }

// func DynCondFloat32GeqFloat32(lvar, rvar float32,
// 	DynMap []ProbInterval,
// 	lvarindex, rvarindex int,
// 	op1, op2 float32,
// 	op1index, op2index, assignedindex int) float32 {

// 	condition_rel := DynMap[lvarindex].Reliability + DynMap[rvarindex].Reliability - float32(1.0)

// 	f32_1 := float32(DynMap[lvarindex].Delta)
// 	f32_2 := float32(DynMap[rvarindex].Delta)

// 	if lvar-f32_1 >= rvar+f32_2 {
// 		DynMap[assignedindex] = DynMap[op1index]
// 		DynMap[assignedindex].Reliability *= condition_rel
// 		// fmt.Printf("[Debug]: B1 %f+-%f %f+-%f\n", lvar, DynMap[lvarindex].Delta, rvar, DynMap[rvarindex].Delta)
// 		return op1
// 	} else if lvar+f32_1 < rvar-f32_2 {
// 		DynMap[assignedindex] = DynMap[op2index]
// 		DynMap[assignedindex].Reliability *= condition_rel
// 		// fmt.Printf("[Debug]: B2 %f+-%f %f+-%f\n", lvar, DynMap[lvarindex].Delta, rvar, DynMap[rvarindex].Delta)
// 		return op2
// 	} else {
// 		// fmt.Printf("[Debug]: B3 %f+-%f %f+-%f\n", lvar, DynMap[lvarindex].Delta, rvar, DynMap[rvarindex].Delta)
// 		DynMap[assignedindex].Delta = float64(Abs32(op1-op2)) * Max64(DynMap[op1index].Delta, DynMap[op2index].Delta)
// 		DynMap[assignedindex].Reliability = Min32(DynMap[op1index].Reliability, DynMap[op2index].Reliability) * condition_rel
// 		if lvar >= rvar {
// 			return op1
// 		} else {
// 			return op2
// 		}
// 	}
// }

// func DynCondFloat64GeqFloat64(lvar, rvar float64,
// 	DynMap []ProbInterval,
// 	lvarindex, rvarindex int,
// 	op1, op2 float64,
// 	op1index, op2index, assignedindex int) float64 {

// 	condition_rel := DynMap[lvarindex].Reliability + DynMap[rvarindex].Reliability - float32(1.0)

// 	f32_1 := DynMap[lvarindex].Delta
// 	f32_2 := DynMap[rvarindex].Delta

// 	if lvar-f32_1 >= rvar+f32_2 {
// 		DynMap[assignedindex] = DynMap[op1index]
// 		DynMap[assignedindex].Reliability *= condition_rel
// 		// fmt.Printf("[Debug]: B1 %f+-%f %f+-%f\n", lvar, DynMap[lvarindex].Delta, rvar, DynMap[rvarindex].Delta)
// 		return op1
// 	} else if lvar+f32_1 < rvar-f32_2 {
// 		DynMap[assignedindex] = DynMap[op1index]
// 		DynMap[assignedindex].Reliability *= condition_rel
// 		// fmt.Printf("[Debug]: B2 %f+-%f %f+-%f\n", lvar, DynMap[lvarindex].Delta, rvar, DynMap[rvarindex].Delta)
// 		return op2
// 	} else {
// 		// fmt.Printf("[Debug]: B3 %f+-%f %f+-%f\n", lvar, DynMap[lvarindex].Delta, rvar, DynMap[rvarindex].Delta)
// 		DynMap[assignedindex].Delta = float64(math.Abs(op1-op2)) * Max64(DynMap[op1index].Delta, DynMap[op1index].Delta)
// 		DynMap[assignedindex].Reliability = Min32(DynMap[op1index].Reliability, DynMap[op1index].Reliability) * condition_rel
// 		if lvar >= rvar {
// 			return op1
// 		} else {
// 			return op2
// 		}
// 	}
// }

// func intArrayFromBytes(bytearray []byte) []int {
// 	outarray := make([]int, 8 * len(intarray))
// 	for i,elem in range(intarray) {
// 		binary.LittleEndian.PutUint64(bs[8*i:], uint64(elem))
// 	}
// 	return bs
// }

func float32ToByte(f float32) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, f)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func Float32frombytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func float32ArrayToByte(inarray []float32) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, inarray)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func float64ToByte(f float64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, f)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func float64ArrayFromBytes(bytes []byte, len int) []float64 {
	var temp_array []float64
	for i := 0; i < len; i++ {
		temp_array = append(temp_array, Float64frombytes(bytes[i*8:(i+1)*8]))
	}
	return temp_array
}

func float64ArrayToByte(inarray []float64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, inarray)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func float32ArrayFromBytes(bytes []byte, len int) []float32 {
	var temp_array []float32
	for i := 0; i < len; i++ {
		temp_array = append(temp_array, Float32frombytes(bytes[i*4:(i+1)*4]))
	}
	return temp_array
}

// func float64ArrayFromBytes(bytes []byte) []float64 {

// 	bits := binary.LittleEndian.Uint64(bytes)
// 	float := math.Float64frombits(bits)
// 	return float
// }

// func probIntervalToBytes(interval ProbInterval) []byte {
// 	// var buf bytes.Buffer
// 	// enc := gob.NewEncoder(&buf)

// 	// err := enc.Encode(interval)
// 	// if err != nil {
// 	// 	fmt.Println("Converting probinterval to bytes failed:", err)
// 	// }

// 	buf := &bytes.Buffer{}
// 	err := binary.Write(buf, binary.BigEndian, interval)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// fmt.Println(buf.Bytes())

// 	return buf.Bytes()
// }

// func probIntervalArrayToBytes(intervals []ProbInterval) []byte {
// 	var buf bytes.Buffer
// 	enc := gob.NewEncoder(&buf)

// 	err := enc.Encode(intervals)
// 	if err != nil {
// 		fmt.Println("binary.Write failed:", err)
// 	}
// 	return buf.Bytes()
// }

// func bytesToProbInterval(bytearray []byte) ProbInterval {
// 	// var buf bytes.Buffer
// 	// var temp ProbInterval

// 	// buf.Write(bytearray)
// 	// dec := gob.NewDecoder(&buf)

// 	// err := dec.Decode(&temp)
// 	// if err != nil {
// 	// 	fmt.Println("Converting bytes to probinterval failed:", err)
// 	// }
// 	// return temp

// 	var temp ProbInterval
// 	err := binary.Read(bytes.NewBuffer(bytearray), binary.BigEndian, &temp)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return temp
// }

// func bytesToProbIntervalArray(bytearray []byte) []ProbInterval {
// 	var buf bytes.Buffer
// 	var temp []ProbInterval

// 	buf.Write(bytearray)
// 	dec := gob.NewDecoder(&buf)

// 	err := dec.Decode(&temp)
// 	if err != nil {
// 		fmt.Println("binary.Write failed:", err)
// 	}

// 	return temp
// }

func InitDynArray(varname int, size int, DynMap []float32) {
	// fmt.Println("Initializing dynamic array: ", varname, size)
	for i := 0; i < size; i++ {
		DynMap[varname+i] = 1.0
	}
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

func CopyDynArray(array1 int, array2 int, size int, DynMap []float32) bool {
	for i := 0; i < size; i++ {
		DynMap[array1+size] = DynMap[array2+size]
	}
	return true
}

func CheckArray(start int, limit float32, size int, DynMap []float32) bool {
	failed := true
	for i := start; i < size; i++ {
		if failed && (DynMap[i] < limit) {
			fmt.Println("Verification failed due to reliability of: ", DynMap[i])
			failed = false
		}
	}
	return failed
}

// func CheckFloat64(val float64, PI ProbInterval, epsThresh float32, deltaThresh float64) (result bool) {
// 	result = (PI.Reliability < epsThresh && PI.Delta < deltaThresh)
// 	//fmt.Println(result)
// 	return
// }

// func PrintWorstElement(DynMap []ProbInterval, start int, end int) {
// 	var min float32 = DynMap[start].Reliability
// 	var maxd float64 = DynMap[start].Delta

// 	for i := start; i < end; i++ {
// 		if min > DynMap[i].Reliability {
// 			min = DynMap[i].Reliability
// 		}
// 		if maxd < DynMap[i].Delta {
// 			maxd = DynMap[i].Delta
// 		}
// 	}

// 	fmt.Println("Worst element: epsilon, delta", maxd, min)
// }

func DumpDynMap(DynMap []float32, filename string) {
	f, _ := os.Create(filename)
	defer f.Close()

	// v := make([]float64, 0, len(DynMap))

	// for  _, value := range DynMap {
	// 	v = append(v, value)
	// }

	// jsonString, _ := json.Marshal(DynMap)

	f.WriteString(fmt.Sprintln(DynMap))
}

func GetCastingError64to32(original float64, casted float32) float64 {
	recasted := float64(casted)
	// fmt.Println(recasted - original)
	return math.Abs(recasted - original)
}

func Cast64to32Array(array32 []float32, array64 []float64) {
	for i, f64 := range array64 {
		array32[i] = float32(f64)
	}
}

func Cast32to64Array(array32 []float32, array64 []float64) {
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

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
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

func createChannels(channelMap map[int]amqp.Queue, numprocesses_in int, name string) {
	for i := 0; i < numprocesses_in; i++ {
		channelName := fmt.Sprintf("%d_%s", i, name)
		// fmt.Println("creating channel: " + channelName)
		q, err := ch.QueueDeclare(
			channelName, // name
			false,       // durable
			true,        // delete when unused
			false,       // exclusive
			false,       // no-wait
			nil,         // arguments
		)
		failOnError(err, "Failed to declare a queue")
		channelMap[i] = q
	}

}

func InitQueues(numprocesses_in int, link string) {
	var err error
	conn, err = amqp.Dial(link)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	pingchannel, err = ch.QueueDeclare(
		"ping", // name
		false,  // durable
		true,   // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	DynamicChannelMap = make(map[int]amqp.Queue)
	createChannels(DynamicChannelMap, numprocesses_in, "dyn")

	DynamicChannelMapArray = make(map[int]amqp.Queue)
	createChannels(DynamicChannelMapArray, numprocesses_in, "dynarray")

	syncChannelMap = make(map[int]amqp.Queue)
	createChannels(syncChannelMap, numprocesses_in, "sync")

	preciseChannelMapInt = make(map[int]amqp.Queue)
	createChannels(preciseChannelMapInt, numprocesses_in, "int")

	preciseChannelMapFloat64 = make(map[int]amqp.Queue)
	createChannels(preciseChannelMapFloat64, numprocesses_in, "float64")

	preciseChannelMapFloat32 = make(map[int]amqp.Queue)
	createChannels(preciseChannelMapFloat32, numprocesses_in, "float32")

	preciseChannelMapIntArray = make(map[int]amqp.Queue)
	createChannels(preciseChannelMapIntArray, numprocesses_in, "intarray")

	approxChannelMapInt = make(map[int]amqp.Queue)
	createChannels(approxChannelMapInt, numprocesses_in, "approxint")

	approxChannelMapIntArray = make(map[int]amqp.Queue)
	createChannels(approxChannelMapIntArray, numprocesses_in, "approxintarray")

	preciseChannelMapFloat64Array = make(map[int]amqp.Queue)
	createChannels(preciseChannelMapFloat64Array, numprocesses_in, "float64array")

	preciseChannelMapFloat32Array = make(map[int]amqp.Queue)
	createChannels(preciseChannelMapFloat32Array, numprocesses_in, "float32array")

	approxChannelMapFloat64Array = make(map[int]amqp.Queue)
	createChannels(approxChannelMapFloat64Array, numprocesses_in, "approxfloat64array")

	approxChannelMapFloat32Array = make(map[int]amqp.Queue)
	createChannels(approxChannelMapFloat32Array, numprocesses_in, "approxfloat32array")
}

func CleanupMain() {
	defer ch.Close()
	defer conn.Close()

	for _, queue := range DynamicChannelMapArray {
		ch.QueueDelete(queue.Name, false, false, false)
	}
	for _, queue := range DynamicChannelMap {
		ch.QueueDelete(queue.Name, false, false, false)
	}

	for _, queue := range approxChannelMapIntArray {
		ch.QueueDelete(queue.Name, false, false, false)
	}
	for _, queue := range approxChannelMapInt {
		ch.QueueDelete(queue.Name, false, false, false)
	}
	for _, queue := range preciseChannelMapIntArray {
		ch.QueueDelete(queue.Name, false, false, false)
	}
	for _, queue := range preciseChannelMapInt {
		ch.QueueDelete(queue.Name, false, false, false)
	}

	for _, queue := range preciseChannelMapFloat64 {
		ch.QueueDelete(queue.Name, false, false, false)
	}
	for _, queue := range approxChannelMapFloat64 {
		ch.QueueDelete(queue.Name, false, false, false)
	}

	for _, queue := range preciseChannelMapFloat64Array {
		ch.QueueDelete(queue.Name, false, false, false)
	}
	for _, queue := range approxChannelMapFloat64Array {
		ch.QueueDelete(queue.Name, false, false, false)
	}

	for _, queue := range preciseChannelMapFloat32Array {
		ch.QueueDelete(queue.Name, false, false, false)
	}
	for _, queue := range approxChannelMapFloat32Array {
		ch.QueueDelete(queue.Name, false, false, false)
	}

	for _, queue := range preciseChannelMapFloat32 {
		ch.QueueDelete(queue.Name, false, false, false)
	}
	for _, queue := range approxChannelMapFloat32 {
		ch.QueueDelete(queue.Name, false, false, false)
	}

	for _, queue := range approxChannelMapInt {
		ch.QueueDelete(queue.Name, false, false, false)
	}
	for _, queue := range preciseChannelMapIntArray {
		ch.QueueDelete(queue.Name, false, false, false)
	}
	for _, queue := range preciseChannelMapInt {
		ch.QueueDelete(queue.Name, false, false, false)
	}

	ch.QueueDelete(pingchannel.Name, false, false, false)
}

func Cleanup() {
	ch.Close()
	conn.Close()
}

func xorSumUpdate(sum uint32, data uint64) uint32 {
	return sum ^ uint32(data>>32) ^ uint32(data&0xFFFFFFFF)
}

func PingMain(tid int) {
	err := ch.Publish(
		"",               // exchange
		pingchannel.Name, // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        intToByte(0),
		})
	failOnError(err, "Failed to publish a message")

	for {
		_, ok, err := ch.Get(syncChannelMap[tid].Name, true)
		failOnError(err, "Failed to register a consumer")
		if ok {
			break
		}
	}
}

func WaitForWorkers(numthreads int) {
	for i := 1; i < numthreads; i++ {
		err := ch.Publish(
			"",                     // exchange
			syncChannelMap[i].Name, // routing key
			false,                  // mandatory
			false,                  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        intToByte(0),
			})
		failOnError(err, "Failed to register a consumer")
	}

	pingcount := 0
	for {
		_, ok, err := ch.Get(pingchannel.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			pingcount += 1
			if pingcount == (numthreads - 1) {
				return
			}
		}
	}
}

func SendDynVal(value float32, sender, receiver int) {
	my_chan_index := sender*Numprocesses + receiver
	q2 := DynamicChannelMap[my_chan_index]
	err := ch.Publish(
		"",      // exchange
		q2.Name, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        float32ToByte(value),
		})
	failOnError(err, "Failed to publish a message")
}

func GetDynValue(index int) float32 {
	q2 := DynamicChannelMap[index]
	for {
		msg, ok, err := ch.Get(q2.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			// fmt.Println(len(msg.Body), msg.Body)
			temp := Float32frombytes(msg.Body)
			return temp
		}
	}
}

func SendInt(value, sender, receiver int) {
	my_chan_index := sender*Numprocesses + receiver
	q := preciseChannelMapInt[my_chan_index]

	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        intToByte(value),
		})
	failOnError(err, "Failed to publish a message")

	if debug == 1 {
		fmt.Printf("%d Sending message in precise int chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, Numprocesses, receiver)
	}
}

func ReceiveInt(rec_var *int, receiver, sender int) {
	my_chan_index := sender*Numprocesses + receiver
	q := preciseChannelMapInt[my_chan_index]

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			*rec_var = intfrombytes(msg.Body)
			return
		}
	}
}

func SendIntArray(value []int, sender, receiver int) {
	my_chan_index := sender*Numprocesses + receiver

	q := preciseChannelMapIntArray[my_chan_index]

	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        intArrayToByte(value[:]),
		})
	failOnError(err, "Failed to send int array")
}

func ReceiveIntArray(rec_var []int, receiver, sender int) {
	my_chan_index := sender*Numprocesses + receiver
	q := preciseChannelMapIntArray[my_chan_index]

	temp_array := make([]int, len(rec_var))

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			// fmt.Println(len(msg.Body), msg.Body)
			for i, _ := range rec_var {
				temp_array[i] = intfrombytes(msg.Body[i*8 : (i+1)*8])
			}
			copy(rec_var, temp_array)
			return
		}
	}
}

func SendDynIntArray(value []int, sender, receiver int, DynMap []float32, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapIntArray[my_chan_index]
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        append(intArrayToByte(value[:]), float32ArrayToByte(DynMap[start:start+len(value)])...),
		})
	failOnError(err, "Failed to publish a message")

	// q2 := DynamicChannelMapArray[my_chan_index]
	// err = ch.Publish(
	// 	"",      // exchange
	// 	q2.Name, // routing key
	// 	false,   // mandatory
	// 	false,   // immediate
	// 	amqp.Publishing{
	// 		ContentType: "text/plain",
	// 		Body:        probIntervalArrayToBytes(DynMap[start : start+len(value)]),
	// 	})
	// failOnError(err, "Failed to publish a message")
}

func ReceiveDynIntArray(rec_var []int, receiver, sender int, DynMap []float32, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapIntArray[my_chan_index]
	temp_array := make([]int, len(rec_var))
	// temp_array2 := make([]int, len(rec_var))

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			temp_array2 := float32ArrayFromBytes(msg.Body[8*len(rec_var):], len(rec_var))
			// fmt.Println(len(msg.Body), msg.Body)
			for i, _ := range rec_var {
				temp_array[i] = intfrombytes(msg.Body[i*8 : (i+1)*8])
				DynMap[start+i] = temp_array2[i]
			}
			copy(rec_var, temp_array)
			break
		}
	}

	// q2 := DynamicChannelMapArray[my_chan_index]
	// for {
	// 	msg, ok, err := ch.Get(q2.Name, true)
	// 	failOnError(err, "Failed to register a consumer")
	// 	if ok {
	// 		// fmt.Println(len(msg.Body), msg.Body)
	// 		temp_array2 := bytesToProbIntervalArray(msg.Body)
	// 		for i, _ := range rec_var {
	// 			DynMap[start+i] = temp_array2[i]
	// 		}
	// 		break
	// 	}
	// }
}

func NoisyReceiveDynIntArray(rec_var []int, receiver, sender int, DynMap []float32, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapIntArray[my_chan_index]
	temp_array := make([]int, len(rec_var))
	// temp_array2 := make([]int, len(rec_var))

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			temp_array2 := float32ArrayFromBytes(msg.Body[8*len(rec_var):], len(rec_var))
			// fmt.Println(len(msg.Body), msg.Body)
			for i, _ := range rec_var {
				temp_array[i] = intfrombytes(msg.Body[i*8 : (i+1)*8])
				DynMap[start+i] = temp_array2[i] * noiselevel
			}
			copy(rec_var, temp_array)
			break
		}
	}

	// q2 := DynamicChannelMapArray[my_chan_index]
	// for {
	// 	msg, ok, err := ch.Get(q2.Name, true)
	// 	failOnError(err, "Failed to register a consumer")
	// 	if ok {
	// 		// fmt.Println(len(msg.Body), msg.Body)
	// 		temp_array2 := bytesToProbIntervalArray(msg.Body)
	// 		for i, _ := range rec_var {
	// 			DynMap[start+i] = temp_array2[i]
	// 			DynMap[start+i].Reliability *= noiselevel
	// 		}
	// 		break
	// 	}
	// }
}

func SendDynIntArrayO1(value []int, sender, receiver int, DynMap []float32, start int) {
	my_chan_index := sender*Numprocesses + receiver

	var min float32 = DynMap[start]
	// var maxd float64 = DynMap[start].Delta
	for i, _ := range value {
		if min > DynMap[start+i] {
			min = DynMap[start+i]
		}
	}

	q := approxChannelMapIntArray[my_chan_index]
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        append(intArrayToByte(value[:]), float32ToByte(min)...),
		})
	failOnError(err, "Failed to publish a message")

	// q2 := DynamicChannelMapArray[my_chan_index]
	// err = ch.Publish(
	// 	"",      // exchange
	// 	q2.Name, // routing key
	// 	false,   // mandatory
	// 	false,   // immediate
	// 	amqp.Publishing{
	// 		ContentType: "text/plain",
	// 		Body:        probIntervalToBytes(ProbInterval{min, maxd}),
	// 	})
	// failOnError(err, "Failed to publish a message")
}

func ReceiveDynIntArrayO1(rec_var []int, receiver, sender int, DynMap []float32, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapIntArray[my_chan_index]
	temp_array := make([]int, len(rec_var))
	// temp_array2 := make([]int, len(rec_var))

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			temp_val := Float32frombytes(msg.Body[8*len(rec_var):])
			// fmt.Println(len(msg.Body), msg.Body)
			for i, _ := range rec_var {
				temp_array[i] = intfrombytes(msg.Body[i*8 : (i+1)*8])
				DynMap[start+i] = temp_val
			}
			copy(rec_var, temp_array)
			break
		}
	}

	// q2 := DynamicChannelMapArray[my_chan_index]
	// for {
	// 	msg, ok, err := ch.Get(q2.Name, true)
	// 	failOnError(err, "Failed to register a consumer")

	// 	if ok {
	// 		// fmt.Println(len(msg.Body), msg.Body)
	// 		temp_val := bytesToProbInterval(msg.Body)
	// 		for i, _ := range rec_var {
	// 			DynMap[start+i] = temp_val
	// 		}
	// 		break
	// 	}
	// }
}

func NoisyReceiveDynIntArrayO1(rec_var []int, receiver, sender int, DynMap []float32, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapIntArray[my_chan_index]
	temp_array := make([]int, len(rec_var))
	// temp_array2 := make([]int, len(rec_var))

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			// fmt.Println(len(msg.Body), msg.Body)
			temp_val := Float32frombytes(msg.Body[len(rec_var)*8:])
			for i, _ := range rec_var {
				temp_array[i] = intfrombytes(msg.Body[i*8 : (i+1)*8])
				DynMap[start+i] = temp_val * noiselevel
			}
			copy(rec_var, temp_array)
			break
		}
	}

	// q2 := DynamicChannelMapArray[my_chan_index]
	// for {
	// 	msg, ok, err := ch.Get(q2.Name, true)
	// 	failOnError(err, "Failed to register a consumer")

	// 	if ok {
	// 		// fmt.Println(len(msg.Body), msg.Body)
	// 		temp_val := bytesToProbInterval(msg.Body)
	// 		for i, _ := range rec_var {
	// 			DynMap[start+i] = temp_val
	// 			DynMap[start+i].Reliability *= noiselevel
	// 		}
	// 		break
	// 	}
	// }
}

////////////////////////////////////////
////////////////////////////////////////
/////////////  FLOAT64 /////////////////
////////////////////////////////////////
////////////////////////////////////////
func SendFloat64(value float64, sender, receiver int) {
	my_chan_index := sender*Numprocesses + receiver
	q := preciseChannelMapFloat64[my_chan_index]

	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        float64ToByte(value),
		})
	failOnError(err, "Failed to publish a message")

	if debug == 1 {
		fmt.Printf("%d Sending message in precise int chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, Numprocesses, receiver)
	}
}

func ReceiveFloat64(rec_var *float64, receiver, sender int) {
	my_chan_index := sender*Numprocesses + receiver
	q := preciseChannelMapFloat64[my_chan_index]

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			*rec_var = Float64frombytes(msg.Body)
			return
		}
	}
}

func SendFloat64Array(value []float64, sender, receiver int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapFloat64Array[my_chan_index]

	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        float64ArrayToByte(value[:]),
		})
	failOnError(err, "Failed to send int array")
}

func ReceiveFloat64Array(rec_var []float64, receiver, sender int) {
	my_chan_index := sender*Numprocesses + receiver
	q := approxChannelMapFloat64Array[my_chan_index]

	temp_array := make([]float64, len(rec_var))

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			for i, _ := range rec_var {
				temp_array[i] = Float64frombytes(msg.Body[i*8 : (i+1)*8])
			}
			copy(rec_var, temp_array)
			return
		}
	}
}

func SendDynFloat64Array(value []float64, sender, receiver int, DynMap []float32, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapFloat64Array[my_chan_index]
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body: append(float64ArrayToByte(value[:]),
				float32ArrayToByte(DynMap[start:start+len(value)])...),
		})
	failOnError(err, "Failed to publish a message")

	// q2 := DynamicChannelMapArray[my_chan_index]
	// err = ch.Publish(
	// 	"",      // exchange
	// 	q2.Name, // routing key
	// 	false,   // mandatory
	// 	false,   // immediate
	// 	amqp.Publishing{
	// 		ContentType: "text/plain",
	// 		Body:        probIntervalArrayToBytes(DynMap[start : start+len(value)]),
	// 	})
	// failOnError(err, "Failed to publish a message")
}

func ReceiveDynFloat64Array(rec_var []float64, receiver, sender int, DynMap []float32, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapFloat64Array[my_chan_index]
	temp_array := make([]float64, len(rec_var))
	// temp_array2 := make([]int, len(rec_var))

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")
		if ok {
			temp_array2 := float32ArrayFromBytes(msg.Body[8*len(rec_var):], len(rec_var))
			for i, _ := range rec_var {
				temp_array[i] = Float64frombytes(msg.Body[i*8 : (i+1)*8])
				DynMap[start+i] = temp_array2[i]
			}
			copy(rec_var, temp_array)
			break
		}
	}

	// q2 := DynamicChannelMapArray[my_chan_index]
	// for {
	// 	msg, ok, err := ch.Get(q2.Name, true)
	// 	failOnError(err, "Failed to register a consumer")

	// 	if ok {
	// 		// fmt.Println(len(msg.Body), msg.Body)
	// 		temp_array2 := bytesToProbIntervalArray(msg.Body)
	// 		for i, _ := range rec_var {
	// 			DynMap[start+i] = temp_array2[i]
	// 		}
	// 		break
	// 	}
	// }
}

func NoisyReceiveDynFloat64Array(rec_var []float64, receiver, sender int, DynMap []float32, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapFloat64Array[my_chan_index]
	temp_array := make([]float64, len(rec_var))
	// temp_array2 := make([]int, len(rec_var))

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			temp_array2 := float32ArrayFromBytes(msg.Body[8*len(rec_var):], len(rec_var))
			for i, _ := range rec_var {
				temp_array[i] = Float64frombytes(msg.Body[i*8 : (i+1)*8])
				DynMap[start+i] = temp_array2[i] * noiselevel
			}
			copy(rec_var, temp_array)
			break
		}
	}

	// q2 := DynamicChannelMapArray[my_chan_index]
	// for {
	// 	msg, ok, err := ch.Get(q2.Name, true)
	// 	failOnError(err, "Failed to register a consumer")

	// 	if ok {
	// 		// fmt.Println(len(msg.Body), msg.Body)
	// 		temp_array2 := bytesToProbIntervalArray(msg.Body)
	// 		for i, _ := range rec_var {
	// 			DynMap[start+i] = temp_array2[i]
	// 			DynMap[start+i].Reliability *= noiselevel
	// 		}
	// 		break
	// 	}
	// }
}

func SendDynFloat64ArrayO1(value []float64, sender, receiver int, DynMap []float32, start int) {
	my_chan_index := sender*Numprocesses + receiver

	var min float32 = DynMap[start]
	for i, _ := range value {
		if min > DynMap[start+i] {
			min = DynMap[start+i]
		}
	}

	q := approxChannelMapFloat64Array[my_chan_index]
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        append(float64ArrayToByte(value[:]), float32ToByte(min)...),
		})
	failOnError(err, "Failed to publish a message")

	// q2 := DynamicChannelMapArray[my_chan_index]
	// err = ch.Publish(
	// 	"",      // exchange
	// 	q2.Name, // routing key
	// 	false,   // mandatory
	// 	false,   // immediate
	// 	amqp.Publishing{
	// 		ContentType: "text/plain",
	// 		Body:        probIntervalToBytes(ProbInterval{min, maxd}),
	// 	})
	// failOnError(err, "Failed to publish a message")
}

func ReceiveDynFloat64ArrayO1(rec_var []float64, receiver, sender int, DynMap []float32, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapFloat64Array[my_chan_index]
	temp_array := make([]float64, len(rec_var))
	// temp_array2 := make([]int, len(rec_var))

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			// fmt.Println(len(msg.Body), msg.Body)
			temp_val := Float32frombytes(msg.Body[len(rec_var)*8:])
			for i, _ := range rec_var {
				temp_array[i] = Float64frombytes(msg.Body[i*8 : (i+1)*8])
				DynMap[start+i] = temp_val
			}
			copy(rec_var, temp_array)
			break
		}
	}

	// q2 := DynamicChannelMapArray[my_chan_index]
	// for {
	// 	msg, ok, err := ch.Get(q2.Name, true)
	// 	failOnError(err, "Failed to register a consumer")

	// 	if ok {
	// 		// fmt.Println(len(msg.Body), msg.Body)
	// 		temp_val := bytesToProbInterval(msg.Body)
	// 		for i, _ := range rec_var {
	// 			DynMap[start+i] = temp_val
	// 		}
	// 		break
	// 	}
	// }
}

func NoisyReceiveDynFloat64ArrayO1(rec_var []float64, receiver, sender int, DynMap []float32, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapFloat64Array[my_chan_index]
	temp_array := make([]float64, len(rec_var))
	// temp_array2 := make([]int, len(rec_var))

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			// fmt.Println(len(msg.Body), msg.Body)
			temp_val := Float32frombytes(msg.Body[len(rec_var)*8:])
			for i, _ := range rec_var {
				temp_array[i] = Float64frombytes(msg.Body[i*8 : (i+1)*8])
				DynMap[start+i] = temp_val * noiselevel
			}
			copy(rec_var, temp_array)
			break
		}
	}

	// q2 := DynamicChannelMapArray[my_chan_index]
	// for {
	// 	msg, ok, err := ch.Get(q2.Name, true)
	// 	failOnError(err, "Failed to register a consumer")

	// 	if ok {
	// 		// fmt.Println(len(msg.Body), msg.Body)
	// 		temp_val := bytesToProbInterval(msg.Body)
	// 		for i, _ := range rec_var {
	// 			DynMap[start+i] = temp_val
	// 			DynMap[start+i].Reliability *= noiselevel
	// 		}
	// 		break
	// 	}
	// }
}

// ////////////////////////////////////////
// ////////////////////////////////////////
// /////////////  FLOAT32 /////////////////
// ////////////////////////////////////////
// ////////////////////////////////////////
func SendFloat32(value float32, sender, receiver int) {
	my_chan_index := sender*Numprocesses + receiver
	q := preciseChannelMapFloat32[my_chan_index]

	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        float32ToByte(value),
		})
	failOnError(err, "Failed to publish a message")

	if debug == 1 {
		fmt.Printf("%d Sending message in precise int chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, Numprocesses, receiver)
	}
}

func ReceiveFloat32(rec_var *float32, receiver, sender int) {
	my_chan_index := sender*Numprocesses + receiver
	q := preciseChannelMapFloat32[my_chan_index]

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			*rec_var = Float32frombytes(msg.Body)
			return
		}
	}
}

func SendFloat32Array(value []float32, sender, receiver int) {
	my_chan_index := sender*Numprocesses + receiver

	q := preciseChannelMapFloat32Array[my_chan_index]

	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        float32ArrayToByte(value[:]),
		})
	failOnError(err, "Failed to send int array")
}

func ReceiveFloat32Array(rec_var []float32, receiver, sender int) {
	my_chan_index := sender*Numprocesses + receiver
	q := preciseChannelMapFloat32Array[my_chan_index]

	temp_array := make([]float32, len(rec_var))

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			for i, _ := range rec_var {
				temp_array[i] = Float32frombytes(msg.Body[i*4 : (i+1)*4])
			}
			copy(rec_var, temp_array)
			return
		}
	}
}

func SendDynFloat32Array(value []float32, sender, receiver int, DynMap []float32, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapFloat32Array[my_chan_index]
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        append(float32ArrayToByte(value[:]), float32ArrayToByte(DynMap[start:start+len(value)])...),
		})
	failOnError(err, "Failed to publish a message")

	// q2 := DynamicChannelMapArray[my_chan_index]
	// err = ch.Publish(
	// 	"",      // exchange
	// 	q2.Name, // routing key
	// 	false,   // mandatory
	// 	false,   // immediate
	// 	amqp.Publishing{
	// 		ContentType: "text/plain",
	// 		Body:        probIntervalArrayToBytes(DynMap[start : start+len(value)]),
	// 	})
	// failOnError(err, "Failed to publish a message")
}

func ReceiveDynFloat32Array(rec_var []float32, receiver, sender int, DynMap []float32, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapFloat32Array[my_chan_index]
	temp_array := make([]float32, len(rec_var))
	// temp_array2 := make([]int, len(rec_var))

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			// fmt.Println(len(msg.Body), msg.Body)
			temp_array2 := float32ArrayFromBytes(msg.Body[len(rec_var)*4:], len(rec_var))
			for i, _ := range rec_var {
				temp_array[i] = Float32frombytes(msg.Body[i*4 : (i+1)*4])
				DynMap[start+i] = temp_array2[i]
			}
			copy(rec_var, temp_array)
			break
		}
	}

	// q2 := DynamicChannelMapArray[my_chan_index]
	// for {
	// 	msg, ok, err := ch.Get(q2.Name, true)
	// 	failOnError(err, "Failed to register a consumer")

	// 	if ok {
	// 		// fmt.Println(len(msg.Body), msg.Body)
	// 		temp_array2 := bytesToProbIntervalArray(msg.Body)
	// 		for i, _ := range rec_var {
	// 			DynMap[start+i] = temp_array2[i]
	// 		}
	// 		break
	// 	}
	// }
}

func NoisyReceiveDynFloat32Array(rec_var []float32, receiver, sender int, DynMap []float32, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapFloat32Array[my_chan_index]
	temp_array := make([]float32, len(rec_var))
	// temp_array2 := make([]int, len(rec_var))

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			// fmt.Println(len(msg.Body), msg.Body)
			temp_array2 := float32ArrayFromBytes(msg.Body[len(rec_var)*4:], len(rec_var))
			for i, _ := range rec_var {
				temp_array[i] = Float32frombytes(msg.Body[i*4 : (i+1)*4])
				DynMap[start+i] = temp_array2[i] * noiselevel
			}
			copy(rec_var, temp_array)
			break
		}
	}

	// q2 := DynamicChannelMapArray[my_chan_index]
	// for {
	// 	msg, ok, err := ch.Get(q2.Name, true)
	// 	failOnError(err, "Failed to register a consumer")

	// 	if ok {
	// 		// fmt.Println(len(msg.Body), msg.Body)
	// 		temp_array2 := bytesToProbIntervalArray(msg.Body)
	// 		for i, _ := range rec_var {
	// 			DynMap[start+i] = temp_array2[i]
	// 			DynMap[start+i].Reliability *= noiselevel
	// 		}
	// 		break
	// 	}
	// }
}

func SendDynFloat32ArrayO1(value []float32, sender, receiver int, DynMap []float32, start int) {
	my_chan_index := sender*Numprocesses + receiver

	var min float32 = DynMap[start]
	for i, _ := range value {
		if min > DynMap[start+i] {
			min = DynMap[start+i]
		}
	}

	q := approxChannelMapFloat32Array[my_chan_index]
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        append(float32ArrayToByte(value[:]), float32ToByte(min)...),
		})
	failOnError(err, "Failed to publish a message")

	// q2 := DynamicChannelMapArray[my_chan_index]
	// err = ch.Publish(
	// 	"",      // exchange
	// 	q2.Name, // routing key
	// 	false,   // mandatory
	// 	false,   // immediate
	// 	amqp.Publishing{
	// 		ContentType: "text/plain",
	// 		Body:        probIntervalToBytes(ProbInterval{min, maxd}),
	// 	})
	// failOnError(err, "Failed to publish a message")
}

func ReceiveDynFloat32ArrayO1(rec_var []float32, receiver, sender int, DynMap []float32, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapFloat32Array[my_chan_index]
	temp_array := make([]float32, len(rec_var))
	// temp_array2 := make([]int, len(rec_var))

	msg, ok, err := ch.Get(q.Name, true)
	failOnError(err, "Failed to register a consumer")

	for {
		if ok {
			temp_val := Float32frombytes(msg.Body[len(rec_var)*4:])
			for i, _ := range rec_var {
				temp_array[i] = Float32frombytes(msg.Body[i*4 : (i+1)*4])
				DynMap[start+i] = temp_val
			}
			copy(rec_var, temp_array)
			break
		}
		msg, ok, err = ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")
	}

	// q2 := DynamicChannelMapArray[my_chan_index]
	// msg2, ok2, err := ch.Get(q2.Name, true)
	// failOnError(err, "Failed to register a consumer")
	// for {
	// 	if ok2 {
	// 		// fmt.Println(len(msg.Body), msg.Body)
	// 		temp_val := bytesToProbInterval(msg2.Body)
	// 		for i, _ := range rec_var {
	// 			temp_array[i] = Float32frombytes(msg.Body[i*4 : (i+1)*4])
	// 			DynMap[start+i] = temp_val
	// 		}
	// 		copy(rec_var, temp_array)
	// 		break
	// 	}
	// 	msg2, ok2, err = ch.Get(q2.Name, true)
	// 	failOnError(err, "Failed to register a consumer")
	// }
}

func NoisyReceiveDynFloat32ArrayO1(rec_var []float32, receiver, sender int, DynMap []float32, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapFloat32Array[my_chan_index]
	temp_array := make([]float32, len(rec_var))
	// temp_array2 := make([]int, len(rec_var))

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			temp_val := Float32frombytes(msg.Body[len(rec_var)*4:])
			// fmt.Println(len(msg.Body), msg.Body)
			for i, _ := range rec_var {
				temp_array[i] = Float32frombytes(msg.Body[i*4 : (i+1)*4])
				DynMap[start+i] = temp_val * noiselevel
			}
			copy(rec_var, temp_array)
			break
		}
	}

	// q2 := DynamicChannelMapArray[my_chan_index]
	// for {
	// 	msg, ok, err := ch.Get(q2.Name, true)
	// 	failOnError(err, "Failed to register a consumer")

	// 	if ok {
	// 		// fmt.Println(len(msg.Body), msg.Body)
	// 		temp_val := bytesToProbInterval(msg.Body)
	// 		for i, _ := range rec_var {
	// 			DynMap[start+i] = temp_val
	// 			DynMap[start+i].Reliability *= noiselevel
	// 		}
	// 		break
	// 	}
	// }
}
