// +build !instrument

package diesel

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type ProbInterval struct {
	Reliability float32
	Delta       float64
}

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

// var preciseChannelMapFloat32 map[int]chan float32
// var preciseChannelMapFloat64 map[int]chan float64

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

func Max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func ConvBool(x bool) int {
	if x {
		return 1
	} else {
		return 0
	}
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

// func intArrayFromBytes(bytearray []byte) []int {
// 	outarray := make([]int, 8 * len(intarray))
// 	for i,elem in range(intarray) {
// 		binary.LittleEndian.PutUint64(bs[8*i:], uint64(elem))
// 	}
// 	return bs
// }

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

func float64ArrayToByte(inarray []float64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, inarray)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

// func float64ArrayFromBytes(bytes []byte) []float64 {

// 	bits := binary.LittleEndian.Uint64(bytes)
// 	float := math.Float64frombits(bits)
// 	return float
// }

func probIntervalToBytes(interval ProbInterval) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(interval)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func probIntervalArrayToBytes(intervals []ProbInterval) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(intervals)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func bytesToProbInterval(bytearray []byte) ProbInterval {
	var buf bytes.Buffer
	var temp ProbInterval

	buf.Write(bytearray)
	dec := gob.NewDecoder(&buf)

	err := dec.Decode(&temp)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return temp
}

func bytesToProbIntervalArray(bytearray []byte) []ProbInterval {
	var buf bytes.Buffer
	var temp []ProbInterval

	buf.Write(bytearray)
	dec := gob.NewDecoder(&buf)

	err := dec.Decode(&temp)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return temp
}

func InitDynArray(varname int, size int, DynMap []ProbInterval) {
	// fmt.Println("Initializing dynamic array: ", varname, size)
	for i := 0; i < size; i++ {
		DynMap[varname+i] = ProbInterval{1, 0}
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

func CopyDynArray(array1 int, array2 int, size int, DynMap []ProbInterval) bool {
	for i := 0; i < size; i++ {
		DynMap[array1+size] = DynMap[array2+size]
	}
	return true
}

func CheckArray(start int, limit float32, size int, DynMap []ProbInterval) bool {
	failed := true
	for i := start; i < size; i++ {
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

func PrintWorstElement(DynMap []ProbInterval, start int, end int) {
	var min float32 = DynMap[start].Reliability
	var maxd float64 = DynMap[start].Delta

	for i := start; i < end; i++ {
		if min > DynMap[i].Reliability {
			min = DynMap[i].Reliability
		}
		if maxd < DynMap[i].Delta {
			maxd = DynMap[i].Delta
		}
	}

	fmt.Println("Worst element: epsilon, delta", maxd, min)
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

	preciseChannelMapInt = make(map[int]amqp.Queue)
	for i := 0; i < numprocesses_in*numprocesses_in; i++ {
		channelName := fmt.Sprintf("%d_%s", i, "int")
		q, err := ch.QueueDeclare(
			channelName, // name
			false,       // durable
			true,        // delete when unused
			false,       // exclusive
			false,       // no-wait
			nil,         // arguments
		)
		failOnError(err, "Failed to declare a queue")
		preciseChannelMapInt[i] = q
	}

	preciseChannelMapFloat64 = make(map[int]amqp.Queue)
	for i := 0; i < numprocesses_in*numprocesses_in; i++ {
		channelName := fmt.Sprintf("%d_%s", i, "float64")
		q, err := ch.QueueDeclare(
			channelName, // name
			false,       // durable
			true,        // delete when unused
			false,       // exclusive
			false,       // no-wait
			nil,         // arguments
		)
		failOnError(err, "Failed to declare a queue")
		preciseChannelMapFloat64[i] = q
	}

	preciseChannelMapIntArray = make(map[int]amqp.Queue)
	for i := 0; i < numprocesses_in*numprocesses_in; i++ {
		channelName := fmt.Sprintf("%d_%s", i, "intarray")
		q, err := ch.QueueDeclare(
			channelName, // name
			false,       // durable
			true,        // delete when unused
			false,       // exclusive
			false,       // no-wait
			nil,         // arguments
		)
		failOnError(err, "Failed to declare a queue")
		preciseChannelMapIntArray[i] = q
	}

	approxChannelMapInt = make(map[int]amqp.Queue)
	for i := 0; i < numprocesses_in*numprocesses_in; i++ {
		channelName := fmt.Sprintf("%d_%s", i, "approxint")
		q, err := ch.QueueDeclare(
			channelName, // name
			false,       // durable
			true,        // delete when unused
			false,       // exclusive
			false,       // no-wait
			nil,         // arguments
		)
		failOnError(err, "Failed to declare a queue")
		approxChannelMapInt[i] = q
	}

	approxChannelMapIntArray = make(map[int]amqp.Queue)
	for i := 0; i < numprocesses_in*numprocesses_in; i++ {
		channelName := fmt.Sprintf("%d_%s", i, "approxintarray")
		q, err := ch.QueueDeclare(
			channelName, // name
			false,       // durable
			true,        // delete when unused
			false,       // exclusive
			false,       // no-wait
			nil,         // arguments
		)
		failOnError(err, "Failed to declare a queue")
		approxChannelMapIntArray[i] = q
	}

	preciseChannelMapFloat64Array = make(map[int]amqp.Queue)
	for i := 0; i < numprocesses_in*numprocesses_in; i++ {
		channelName := fmt.Sprintf("%d_%s", i, "float64array")
		q, err := ch.QueueDeclare(
			channelName, // name
			false,       // durable
			true,        // delete when unused
			false,       // exclusive
			false,       // no-wait
			nil,         // arguments
		)
		failOnError(err, "Failed to declare a queue")
		preciseChannelMapInt[i] = q
	}

	approxChannelMapFloat64Array = make(map[int]amqp.Queue)
	for i := 0; i < numprocesses_in*numprocesses_in; i++ {
		channelName := fmt.Sprintf("%d_%s", i, "approxfloat64array")
		q, err := ch.QueueDeclare(
			channelName, // name
			false,       // durable
			true,        // delete when unused
			false,       // exclusive
			false,       // no-wait
			nil,         // arguments
		)
		failOnError(err, "Failed to declare a queue")
		approxChannelMapInt[i] = q
	}

	DynamicChannelMap = make(map[int]amqp.Queue)
	for i := 0; i < numprocesses_in*numprocesses_in; i++ {
		channelName := fmt.Sprintf("%d_%s", i, "dyn")
		q, err := ch.QueueDeclare(
			channelName, // name
			false,       // durable
			true,        // delete when unused
			false,       // exclusive
			false,       // no-wait
			nil,         // arguments
		)
		failOnError(err, "Failed to declare a queue")
		DynamicChannelMap[i] = q
	}

	DynamicChannelMapArray = make(map[int]amqp.Queue)
	for i := 0; i < numprocesses_in*numprocesses_in; i++ {
		channelName := fmt.Sprintf("%d_%s", i, "dynarray")
		q, err := ch.QueueDeclare(
			channelName, // name
			false,       // durable
			true,        // delete when unused
			false,       // exclusive
			false,       // no-wait
			nil,         // arguments
		)
		failOnError(err, "Failed to declare a queue")
		DynamicChannelMapArray[i] = q
	}
}

func CleanupMain() {
	ch.Close()
	conn.Close()

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

// func SendDynVal(value ProbInterval, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	DynamicChannelMap[my_chan_index] <- value
// 	if debug == 1 {
// 		fmt.Printf("%d Sending message in precise int chan : %d (%d * %d + %d)\n",
// 			sender, my_chan_index, sender, Numprocesses, receiver)
// 	}
// }

func PingMain() {
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
}

func WaitForWorkers(numthreads int) {
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

func SendDynVal(value ProbInterval, sender, receiver int) {
	my_chan_index := sender*Numprocesses + receiver
	q2 := DynamicChannelMap[my_chan_index]
	err := ch.Publish(
		"",      // exchange
		q2.Name, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        probIntervalToBytes(value),
		})
	failOnError(err, "Failed to publish a message")
}

func GetDynValue(index int) ProbInterval {
	q2 := DynamicChannelMap[index]
	for {
		msg, ok, err := ch.Get(q2.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			// fmt.Println(len(msg.Body), msg.Body)
			temp := bytesToProbInterval(msg.Body)
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

func SendDynIntArray(value []int, sender, receiver int, DynMap []ProbInterval, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapIntArray[my_chan_index]
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        intArrayToByte(value[:]),
		})
	failOnError(err, "Failed to publish a message")

	q2 := DynamicChannelMapArray[my_chan_index]
	err = ch.Publish(
		"",      // exchange
		q2.Name, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        probIntervalArrayToBytes(DynMap[start : start+len(value)]),
		})
	failOnError(err, "Failed to publish a message")
}

func ReceiveDynIntArray(rec_var []int, receiver, sender int, DynMap []ProbInterval, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapIntArray[my_chan_index]
	temp_array := make([]int, len(rec_var))
	// temp_array2 := make([]int, len(rec_var))

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			// fmt.Println(len(msg.Body), msg.Body)
			for i, _ := range rec_var {
				temp_array[i] = intfrombytes(msg.Body[i*8 : (i+1)*8])
			}
			copy(rec_var, temp_array)
			break
		}
	}

	q2 := DynamicChannelMapArray[my_chan_index]
	for {
		msg, ok, err := ch.Get(q2.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			// fmt.Println(len(msg.Body), msg.Body)
			temp_array2 := bytesToProbIntervalArray(msg.Body)
			for i, _ := range rec_var {
				DynMap[start+i] = temp_array2[i]
			}
			break
		}
	}
}

func SendDynIntArrayO1(value []int, sender, receiver int, DynMap []ProbInterval, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapIntArray[my_chan_index]
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        intArrayToByte(value[:]),
		})
	failOnError(err, "Failed to publish a message")

	var min float32 = DynMap[start].Reliability
	var maxd float64 = DynMap[start].Delta
	for i, _ := range value {
		if min > DynMap[start+i].Reliability {
			min = DynMap[start+i].Reliability
		}
		if maxd < DynMap[start+i].Delta {
			maxd = DynMap[start+i].Delta
		}
	}

	q2 := DynamicChannelMapArray[my_chan_index]
	err = ch.Publish(
		"",      // exchange
		q2.Name, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        probIntervalToBytes(ProbInterval{min, maxd}),
		})
	failOnError(err, "Failed to publish a message")
}

func ReceiveDynIntArrayO1(rec_var []int, receiver, sender int, DynMap []ProbInterval, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapIntArray[my_chan_index]
	temp_array := make([]int, len(rec_var))
	// temp_array2 := make([]int, len(rec_var))

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			// fmt.Println(len(msg.Body), msg.Body)
			for i, _ := range rec_var {
				temp_array[i] = intfrombytes(msg.Body[i*8 : (i+1)*8])
			}
			copy(rec_var, temp_array)
			break
		}
	}

	q2 := DynamicChannelMapArray[my_chan_index]
	for {
		msg, ok, err := ch.Get(q2.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			// fmt.Println(len(msg.Body), msg.Body)
			temp_val := bytesToProbInterval(msg.Body)
			for i, _ := range rec_var {
				DynMap[start+i] = temp_val
			}
			break
		}
	}
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

	q := preciseChannelMapIntArray[my_chan_index]

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
	q := preciseChannelMapIntArray[my_chan_index]

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

func SendDynFloat64Array(value []float64, sender, receiver int, DynMap []ProbInterval, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapIntArray[my_chan_index]
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        float64ArrayToByte(value[:]),
		})
	failOnError(err, "Failed to publish a message")

	q2 := DynamicChannelMapArray[my_chan_index]
	err = ch.Publish(
		"",      // exchange
		q2.Name, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        probIntervalArrayToBytes(DynMap[start : start+len(value)]),
		})
	failOnError(err, "Failed to publish a message")
}

func ReceiveDynFloat64Array(rec_var []float64, receiver, sender int, DynMap []ProbInterval, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapIntArray[my_chan_index]
	temp_array := make([]float64, len(rec_var))
	// temp_array2 := make([]int, len(rec_var))

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			// fmt.Println(len(msg.Body), msg.Body)
			for i, _ := range rec_var {
				temp_array[i] = Float64frombytes(msg.Body[i*8 : (i+1)*8])
			}
			copy(rec_var, temp_array)
			break
		}
	}

	q2 := DynamicChannelMapArray[my_chan_index]
	for {
		msg, ok, err := ch.Get(q2.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			// fmt.Println(len(msg.Body), msg.Body)
			temp_array2 := bytesToProbIntervalArray(msg.Body)
			for i, _ := range rec_var {
				DynMap[start+i] = temp_array2[i]
			}
			break
		}
	}
}

func SendDynFloat64ArrayO1(value []float64, sender, receiver int, DynMap []ProbInterval, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapIntArray[my_chan_index]
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        float64ArrayToByte(value[:]),
		})
	failOnError(err, "Failed to publish a message")

	var min float32 = DynMap[start].Reliability
	var maxd float64 = DynMap[start].Delta
	for i, _ := range value {
		if min > DynMap[start+i].Reliability {
			min = DynMap[start+i].Reliability
		}
		if maxd < DynMap[start+i].Delta {
			maxd = DynMap[start+i].Delta
		}
	}

	q2 := DynamicChannelMapArray[my_chan_index]
	err = ch.Publish(
		"",      // exchange
		q2.Name, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        probIntervalToBytes(ProbInterval{min, maxd}),
		})
	failOnError(err, "Failed to publish a message")
}

func ReceiveDynFloat64ArrayO1(rec_var []float64, receiver, sender int, DynMap []ProbInterval, start int) {
	my_chan_index := sender*Numprocesses + receiver

	q := approxChannelMapIntArray[my_chan_index]
	temp_array := make([]float64, len(rec_var))
	// temp_array2 := make([]int, len(rec_var))

	for {
		msg, ok, err := ch.Get(q.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			// fmt.Println(len(msg.Body), msg.Body)
			for i, _ := range rec_var {
				temp_array[i] = Float64frombytes(msg.Body[i*8 : (i+1)*8])
			}
			copy(rec_var, temp_array)
			break
		}
	}

	q2 := DynamicChannelMapArray[my_chan_index]
	for {
		msg, ok, err := ch.Get(q2.Name, true)
		failOnError(err, "Failed to register a consumer")

		if ok {
			// fmt.Println(len(msg.Body), msg.Body)
			temp_val := bytesToProbInterval(msg.Body)
			for i, _ := range rec_var {
				DynMap[start+i] = temp_val
			}
			break
		}
	}
}

// func SendInt32(value int32, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	preciseChannelMapInt32[my_chan_index] <- value
// 	if debug == 1 {
// 		fmt.Printf("%d Sending message in precise int32 chan : %d (%d * %d + %d)\n",
// 			sender, my_chan_index, sender, Numprocesses, receiver)
// 	}
// }

// func SendInt64(value int64, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	preciseChannelMapInt64[my_chan_index] <- value
// 	if debug == 1 {
// 		fmt.Printf("%d Sending message in precise int64 chan : %d (%d * %d + %d)\n",
// 			sender, my_chan_index, sender, Numprocesses, receiver)
// 	}
// }

// func SendFloat32(value float32, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	preciseChannelMapFloat32[my_chan_index] <- value
// 	if debug == 1 {
// 		fmt.Printf("%d Sending message in precise float32 chan : %d (%d * %d + %d)\n",
// 			sender, my_chan_index, sender, Numprocesses, receiver)
// 	}
// }

// func SendFloat64(value float64, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	preciseChannelMapFloat64[my_chan_index] <- value
// 	// if debug==1 {
// 	// 	fmt.Printf("%d Sending message in precise float64 chan : %d (%d * %d + %d)\n",
// 	// 		sender, my_chan_index, sender, Numprocesses, receiver);
// 	// }
// }

// func SendApprox(value, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	approxChannelMapInt[my_chan_index] <- value
// 	if debug == 1 {
// 		fmt.Printf("%d Sending message in approx int chan : %d (%d * %d + %d)\n",
// 			sender, my_chan_index, sender, Numprocesses, receiver)
// 	}
// }

// func SendInt32Approx(value int32, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	approxChannelMapInt32[my_chan_index] <- value
// 	if debug == 1 {
// 		fmt.Printf("%d Sending message in approx int32 chan : %d (%d * %d + %d)\n",
// 			sender, my_chan_index, sender, Numprocesses, receiver)
// 	}
// }

// func SendInt64Approx(value int64, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	approxChannelMapInt64[my_chan_index] <- value
// 	if debug == 1 {
// 		fmt.Printf("%d Sending message in approx int64 chan : %d (%d * %d + %d)\n",
// 			sender, my_chan_index, sender, Numprocesses, receiver)
// 	}
// }

// func SendIntArrayApprox(value []int, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	temp_array := make([]int, len(value))
// 	copy(temp_array, value)
// 	approxChannelMapIntArray[my_chan_index] <- temp_array
// 	// if debug==1 {
// 	// 	fmt.Printf("%d Sending message in approx int array chan : %d (%d * %d + %d)\n",
// 	// 		sender, my_chan_index, sender, Numprocesses, receiver);
// 	// }
// }

// func SendInt32Array(value []int32, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	temp_array := make([]int32, len(value))
// 	copy(temp_array, value)
// 	preciseChannelMapInt32Array[my_chan_index] <- temp_array
// 	if debug == 1 {
// 		fmt.Printf("%d Sending message in precise float32 array chan : %d (%d * %d + %d)\n",
// 			sender, my_chan_index, sender, Numprocesses, receiver)
// 	}
// }

// func SendInt32ArrayApprox(value []int32, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	temp_array := make([]int32, len(value))
// 	copy(temp_array, value)
// 	approxChannelMapInt32Array[my_chan_index] <- temp_array
// 	if debug == 1 {
// 		fmt.Printf("%d Sending message in precise float32 array chan : %d (%d * %d + %d)\n",
// 			sender, my_chan_index, sender, Numprocesses, receiver)
// 	}
// }

// func SendFloat64Array(value []float64, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	// temp_array := make([]float64, len(value))
// 	// copy(temp_array, value)
// 	// preciseChannelMapFloat64Array[my_chan_index] <- temp_array

// 	for i := range value {
// 		preciseChannelMapFloat64[my_chan_index] <- value[i]
// 	}

// }

// func ReceiveFloat64Array(rec_var []float64, receiver, sender int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	// temp_rec_val := <- preciseChannelMapFloat64Array[my_chan_index]

// 	for i := range rec_var {
// 		rec_var[i] = <-preciseChannelMapFloat64[my_chan_index]
// 	}
// 	// if len(rec_var) != len(temp_rec_val) {
// 	// 	rec_var = make([]float64, len(temp_rec_val))
// 	// }

// 	// copy(rec_var, temp_rec_val)
// }

// // func SendDynFloat64Array(value []float64, sender, receiver int, DynMap []float64, start int) {
// // 	my_chan_index := sender * Numprocesses + receiver
// // 	temp_array := make([]float64, len(value))
// // 	copy(temp_array, value)
// // 	preciseChannelMapFloat64Array[my_chan_index] <- temp_array

// // 	for i:=0; i<len(value); i++ {
// // 		DynamicChannelMap[my_chan_index] <- DynMap[start + i]
// // 	}

// // 	if debug==1 {
// // 		fmt.Printf("%d Sending message in precise float64 array chan : %d (%d * %d + %d)\n",
// // 			sender, my_chan_index, sender, Numprocesses, receiver);
// // 	}
// // }

// func SendFloat32Array(value []float32, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	temp_array := make([]float32, len(value))
// 	copy(temp_array, value)
// 	preciseChannelMapFloat32Array[my_chan_index] <- temp_array
// 	if debug == 1 {
// 		fmt.Printf("%d Sending message in precise float32 array chan : %d (%d * %d + %d)\n",
// 			sender, my_chan_index, sender, Numprocesses, receiver)
// 	}
// }

// func ReceiveDynVal(rec_var *ProbInterval, receiver, sender int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	temp_rec_val := <-DynamicChannelMap[my_chan_index]
// 	if debug == 1 {
// 		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
// 			receiver, my_chan_index, sender, Numprocesses, receiver)
// 	}
// 	*rec_var = temp_rec_val
// }

// func ReceiveInt(rec_var *int, receiver, sender int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	temp_rec_val := <-preciseChannelMapInt[my_chan_index]
// 	if debug == 1 {
// 		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
// 			receiver, my_chan_index, sender, Numprocesses, receiver)
// 	}
// 	*rec_var = temp_rec_val
// }

// func ReceiveInt32(rec_var *int32, receiver, sender int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	temp_rec_val := <-preciseChannelMapInt32[my_chan_index]
// 	if debug == 1 {
// 		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
// 			receiver, my_chan_index, sender, Numprocesses, receiver)
// 	}
// 	*rec_var = temp_rec_val
// }

// func ReceiveInt64(rec_var *int64, receiver, sender int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	temp_rec_val := <-preciseChannelMapInt64[my_chan_index]
// 	if debug == 1 {
// 		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
// 			receiver, my_chan_index, sender, Numprocesses, receiver)
// 	}
// 	*rec_var = temp_rec_val
// }

// func ReceiveFloat32(rec_var *float32, receiver, sender int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	temp_rec_val := <-preciseChannelMapFloat32[my_chan_index]
// 	if debug == 1 {
// 		fmt.Printf("%d Received message in precise float32 chan : %d (%d * %d + %d)\n",
// 			receiver, my_chan_index, sender, Numprocesses, receiver)
// 	}
// 	*rec_var = temp_rec_val
// }

// func ReceiveFloat64(rec_var *float64, receiver, sender int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	temp_rec_val := <-preciseChannelMapFloat64[my_chan_index]
// 	if debug == 1 {
// 		fmt.Printf("%d Received message in precise float64 chan : %d (%d * %d + %d)\n",
// 			receiver, my_chan_index, sender, Numprocesses, receiver)
// 	}
// 	*rec_var = temp_rec_val
// }

// func ReceiveIntApprox(rec_var *int, receiver, sender int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	temp_rec_val := <-approxChannelMapInt[my_chan_index]
// 	if debug == 1 {
// 		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
// 			receiver, my_chan_index, sender, Numprocesses, receiver)
// 	}
// 	*rec_var = temp_rec_val
// }

// func ReceiveInt32Approx(rec_var *int32, receiver, sender int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	temp_rec_val := <-approxChannelMapInt32[my_chan_index]
// 	if debug == 1 {
// 		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
// 			receiver, my_chan_index, sender, Numprocesses, receiver)
// 	}
// 	*rec_var = temp_rec_val
// }

// func ReceiveInt64Approx(rec_var *int64, receiver, sender int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	temp_rec_val := <-approxChannelMapInt64[my_chan_index]
// 	if debug == 1 {
// 		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
// 			receiver, my_chan_index, sender, Numprocesses, receiver)
// 	}
// 	*rec_var = temp_rec_val
// }

// func SendIntArray(value []int, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	// temp_array := make([]int, len(value))
// 	// copy(temp_array, value)
// 	for i := range value {
// 		preciseChannelMapInt[my_chan_index] <- value[i]
// 	}
// 	// if debug==1 {
// 	// 	fmt.Printf("%d Sending message in precise float32 array chan : %d (%d * %d + %d)\n",
// 	// 		sender, my_chan_index, sender, Numprocesses, receiver);
// 	// }
// }

// func ReceiveIntArray(rec_var []int, receiver, sender int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	// temp_rec_var := <- preciseChannelMapIntArray[my_chan_index]
// 	// copy(rec_var, temp_rec_var)
// 	for i := range rec_var {
// 		rec_var[i] = <-preciseChannelMapInt[my_chan_index]
// 	}
// 	// if debug==1 {
// 	// 	fmt.Printf("%d Received message in precise int array chan : %d (%d * %d + %d)\n",
// 	// 		receiver, my_chan_index, sender, Numprocesses, receiver);
// 	// }
// 	// fmt.Println(len(rec_var), len(temp_rec_val))
// 	// if len(rec_var) != len(temp_rec_val) {
// 	// 	rec_var = make([]int, len(temp_rec_val))
// 	// 	fmt.Println("=======>", len(rec_var), len(temp_rec_val))
// 	// }
// 	// copy(rec_var, temp_rec_val)
// }

// func ReceiveInt32Array(rec_var []int32, receiver, sender int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	temp_rec_val := <-preciseChannelMapInt32Array[my_chan_index]
// 	if debug == 1 {
// 		fmt.Printf("%d Received message in precise int32 array chan : %d (%d * %d + %d)\n",
// 			receiver, my_chan_index, sender, Numprocesses, receiver)
// 	}
// 	// if len(rec_var) != len(temp_rec_val) {
// 	// 	rec_var = make([]int32, len(temp_rec_val))
// 	// }
// 	copy(rec_var, temp_rec_val)
// }

// // func ReceiveFloat64Array(rec_var []float64, receiver, sender int) {
// // 	my_chan_index := sender * Numprocesses + receiver
// // 	// temp_rec_val := <- preciseChannelMapFloat64Array[my_chan_index]

// // 	for i := range(rec_var) {
// // 		rec_var[i] = <- preciseChannelMapFloat64[my_chan_index]
// // 	}
// // 	// if len(rec_var) != len(temp_rec_val) {
// // 	// 	rec_var = make([]float64, len(temp_rec_val))
// // 	// }
// // 	// copy(rec_var, temp_rec_val)
// // }

// func ReceiveFloat32Array(rec_var []float32, receiver, sender int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	temp_rec_val := <-preciseChannelMapFloat32Array[my_chan_index]
// 	if debug == 1 {
// 		fmt.Printf("%d Received message in precise float32 array chan : %d (%d * %d + %d)\n",
// 			receiver, my_chan_index, sender, Numprocesses, receiver)
// 	}
// 	// if len(rec_var) != len(temp_rec_val) {
// 	// 	rec_var = make([]float32, len(temp_rec_val))
// 	// }
// 	copy(rec_var, temp_rec_val)
// }

// func Condsend(cond, value, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	if debug == 1 {
// 		fmt.Printf("%d Sending message in approx int chan : %d (%d * %d + %d)\n", sender,
// 			my_chan_index, sender, Numprocesses, receiver)
// 	}
// 	if cond != 0 {
// 		approxChannelMapInt[my_chan_index] <- value
// 	} else {
// 		if debug == 1 {
// 			fmt.Printf("[Failure %d] %d Sending message in approx int chan : %d (%d * %d + %d)\n", cond, sender,
// 				my_chan_index, sender, Numprocesses, receiver)
// 		}
// 		approxChannelMapInt[my_chan_index] <- -1
// 	}
// }

// func CondsendIntArray(cond int, value []int, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	if debug == 1 {
// 		fmt.Printf("%d Sending message in approx int chan : %d (%d * %d + %d)\n", sender,
// 			my_chan_index, sender, Numprocesses, receiver)
// 	}
// 	if cond != 0 {
// 		temp_array := make([]int, len(value))
// 		copy(temp_array, value)
// 		approxChannelMapIntArray[my_chan_index] <- temp_array
// 	} else {
// 		if debug == 1 {
// 			fmt.Printf("[Failure %d] %d Cond Sending message in approx int chan : %d (%d * %d + %d)\n", cond, sender,
// 				my_chan_index, sender, Numprocesses, receiver)
// 		}
// 		approxChannelMapIntArray[my_chan_index] <- []int{}
// 	}
// }

// func CondsendInt32(cond, value int32, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	if debug == 1 {
// 		fmt.Printf("%d Sending message in approx int32 chan : %d (%d * %d + %d)\n", sender,
// 			my_chan_index, sender, Numprocesses, receiver)
// 	}
// 	if cond != 0 {
// 		approxChannelMapInt32[my_chan_index] <- value
// 	} else {
// 		if debug == 1 {
// 			fmt.Printf("[Failure %d] %d Sending message in approx int32 chan : %d (%d * %d + %d)\n", cond, sender,
// 				my_chan_index, sender, Numprocesses, receiver)
// 		}
// 		approxChannelMapInt32[my_chan_index] <- -1
// 	}
// }

// func CondsendInt64(cond, value int64, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	if debug == 1 {
// 		fmt.Printf("%d Sending message in approx int64 chan : %d (%d * %d + %d)\n", sender,
// 			my_chan_index, sender, Numprocesses, receiver)
// 	}
// 	if cond != 0 {
// 		approxChannelMapInt64[my_chan_index] <- value
// 	} else {
// 		if debug == 1 {
// 			fmt.Printf("[Failure %d] %d Sending message in approx int64 chan : %d (%d * %d + %d)\n", cond, sender,
// 				my_chan_index, sender, Numprocesses, receiver)
// 		}
// 		approxChannelMapInt64[my_chan_index] <- -1
// 	}
// }

// func CondsendFloat32(cond, value float32, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	if debug == 1 {
// 		fmt.Printf("%d Sending message in approx chan : %d (%d * %d + %d)\n", sender,
// 			my_chan_index, sender, Numprocesses, receiver)
// 	}
// 	if cond != 0 {
// 		approxChannelMapFloat32[my_chan_index] <- value
// 	} else {
// 		if debug == 1 {
// 			fmt.Printf("[Failure %d] %d Sending message in approx chan : %d (%d * %d + %d)\n", cond, sender,
// 				my_chan_index, sender, Numprocesses, receiver)
// 		}
// 		approxChannelMapFloat32[my_chan_index] <- -1
// 	}
// }

// func CondsendFloat64(cond, value float64, sender, receiver int) {
// 	my_chan_index := sender*Numprocesses + receiver
// 	if debug == 1 {
// 		fmt.Printf("%d Sending message in approx chan : %d (%d * %d + %d)\n", sender,
// 			my_chan_index, sender, Numprocesses, receiver)
// 	}
// 	if cond != 0 {
// 		approxChannelMapFloat64[my_chan_index] <- value
// 	} else {
// 		if debug == 1 {
// 			fmt.Printf("[Failure %d] %d Sending message in approx chan : %d (%d * %d + %d)\n", cond, sender,
// 				my_chan_index, sender, Numprocesses, receiver)
// 		}
// 		approxChannelMapFloat64[my_chan_index] <- -1
// 	}
// }

// func Condreceive(rec_cond_var, rec_var *int, receiver, sender int) {
// 	my_chan_index := sender*Numprocesses + receiver

// 	if debug == 1 {
// 		fmt.Printf("---- %d Waiting to Receive from approx int chan : %d (%d * %d + %d)\n",
// 			receiver, my_chan_index, sender, Numprocesses, receiver)
// 	}

// 	temp_rec_val := <-approxChannelMapInt[my_chan_index]

// 	if debug == 1 {
// 		fmt.Printf("%d Recieved message in approx chan : %d (%d)\n", receiver, my_chan_index,
// 			temp_rec_val)
// 	}

// 	if temp_rec_val != -1 {
// 		*rec_var = temp_rec_val
// 		*rec_cond_var = 1
// 	} else {
// 		*rec_cond_var = 0
// 	}
// }

// func CondreceiveIntArray(rec_cond_var *int, rec_var []int, receiver, sender int) {
// 	my_chan_index := sender*Numprocesses + receiver

// 	if debug == 1 {
// 		fmt.Printf("---- %d Waiting to Receive from approx int array chan : %d (%d * %d + %d)\n",
// 			receiver, my_chan_index, sender, Numprocesses, receiver)
// 	}

// 	temp_rec_val := <-approxChannelMapIntArray[my_chan_index]

// 	if len(temp_rec_val) != 0 {
// 		rec_var = temp_rec_val
// 		*rec_cond_var = 1
// 	} else {
// 		// if debug==0 {
// 		// 	fmt.Printf("%d Recieved failed message in approx chan : %d (%d, %d)\n", receiver, my_chan_index,
// 		// 		temp_rec_val, len(temp_rec_val));
// 		// }
// 		*rec_cond_var = 0
// 	}
// }

// func CondreceiveInt32(rec_cond_var, rec_var *int32, receiver, sender int) {
// 	my_chan_index := sender*Numprocesses + receiver

// 	if debug == 1 {
// 		fmt.Printf("---- %d Waiting to Receive from approx int chan : %d (%d * %d + %d)\n",
// 			receiver, my_chan_index, sender, Numprocesses, receiver)
// 	}

// 	temp_rec_val := <-approxChannelMapInt32[my_chan_index]

// 	if debug == 1 {
// 		fmt.Printf("%d Recieved message in approx chan : %d (%d)\n", receiver, my_chan_index,
// 			temp_rec_val)
// 	}

// 	if temp_rec_val != -1 {
// 		*rec_var = temp_rec_val
// 		*rec_cond_var = 1
// 	} else {
// 		*rec_cond_var = 0
// 	}
// }

// func Hoeffding(n int, delta float64) (eps float64) {
// 	eps = math.Sqrt((0.6*math.Log((math.Log(float64(1.1*float64(n+1)))/math.Log(1.10))) + 0.555*math.Log(24/delta)) / float64(n+1))
// 	return
// }

// func FuseFloat64(arr []float64, dynMap []ProbInterval) (mean float64, newInterval ProbInterval) {
// 	var ns []int
// 	var totalN int = 0
// 	var sum float64 = 0

// 	//var mean float64
// 	for i := 0; i < len(dynMap); i++ {
// 		ns = append(ns, ComputeN(dynMap[i]))
// 		totalN = totalN + ns[i]
// 		sum = sum + (arr[i] * float64(ns[i]))
// 	}

// 	mean = sum / float64(totalN)
// 	var eps float64 = Hoeffding(totalN, dynMap[0].Delta)
// 	newInterval.Reliability = float32(eps)
// 	newInterval.Delta = dynMap[0].Delta
// 	return
// }

// func ComputeN(ui ProbInterval) (n int) {
// 	var eps float64 = float64(ui.Reliability)
// 	var delta float64 = ui.Delta
// 	n = int(0.5 * (1 / (eps * eps)) * math.Log((2 / (1 - delta))))
// 	return n

// }

// func AddProbInterval(val1, val2 float64, fst, snd ProbInterval) (retval float64, out ProbInterval) {
// 	out.Reliability = fst.Reliability + snd.Reliability
// 	out.Delta = fst.Delta + snd.Delta
// 	retval = val1 + val2
// 	return
// }

// func MulProbInterval(val1, val2 float64, fst, snd ProbInterval) (retval float64, out ProbInterval) {
// 	retval = val1 * val2
// 	out.Reliability = (float32(math.Abs(val1)) * snd.Reliability) + (float32(math.Abs(val2)) * fst.Reliability) + (fst.Reliability * snd.Reliability)
// 	out.Delta = fst.Delta + snd.Delta
// 	return
// }

// func DivProbInterval(val1, val2 float64, fst, snd ProbInterval) (retval float64, out ProbInterval) {
// 	retval = val1 / val2
// 	out.Reliability = (float32(math.Abs(val1)) * snd.Reliability) + (float32(math.Abs(val2))*fst.Reliability)/(float32(math.Abs(val2))*(float32(math.Abs(val2))-snd.Reliability))
// 	out.Delta = fst.Delta + snd.Delta
// 	return
// }

// func PrintMemory() {
// 	fmt.Println("Memory not Instrumented")
// }

// //Sasa's proposed addition to the runtime: add a custome class for tracking means of Boolean Indicator Random Vars
// type BooleanTracker struct {
// 	successes    int
// 	totalSamples int
// 	mean         float64
// 	//delta float64
// 	//eps float64
// 	meanProbInt ProbInterval
// }

// func NewBooleanTracker() (b BooleanTracker) {
// 	b.successes = 0
// 	b.totalSamples = 0
// 	b.mean = 0
// 	b.meanProbInt.Delta = 1
// 	b.meanProbInt.Reliability = 1
// 	return
// }

// func (b *BooleanTracker) SetDelta(d float64) {
// 	b.meanProbInt.Delta = d
// }

// func (b *BooleanTracker) GetMean() float64 {
// 	return b.mean
// }

// func (b *BooleanTracker) GetInterval() (res ProbInterval) {
// 	return b.meanProbInt
// }

// func (b *BooleanTracker) AddSample(samp int) {
// 	b.successes = b.successes + samp
// 	b.totalSamples = b.totalSamples + 1
// 	b.Hoeffding()
// 	b.ComputeMean()
// }

// func (b *BooleanTracker) Hoeffding() {
// 	b.meanProbInt.Reliability = float32(math.Sqrt((0.6*math.Log((math.Log(float64(1.1*float64(b.totalSamples+1)))/math.Log(1.10))) + 0.555*math.Log(24/b.meanProbInt.Delta)) / float64(b.totalSamples+1)))
// }

// func (b *BooleanTracker) ComputeMean() {
// 	b.mean = float64(b.successes) / float64(b.totalSamples)
// }

// //func (b *BooleanTracker) Check(c float32) bool{
// //	CheckFloat64(val float64, PI ProbInterval, epsThresh float32, deltaThresh float64)
// //}

// func FuseBooleanTrackers(arr []BooleanTracker) (res BooleanTracker) {
// 	res = NewBooleanTracker()
// 	for i := 0; i < len(arr); i++ {
// 		res.totalSamples = res.totalSamples + arr[i].totalSamples
// 		res.successes = res.successes + arr[i].successes
// 	}

// 	res.Hoeffding()
// 	res.GetMean()
// 	return

// }

// func FuseFloat64IntoBooleanTracker(arr []float64, dynMap []ProbInterval) (res BooleanTracker) {

// 	var ns []int
// 	var totalN int = 0
// 	var sum float64 = 0

// 	//var mean float64
// 	for i := 0; i < len(dynMap); i++ {
// 		ns = append(ns, ComputeN(dynMap[i]))
// 		totalN = totalN + ns[i]
// 		sum = sum + (arr[i] * float64(ns[i]))
// 	}

// 	res = NewBooleanTracker()
// 	res.successes = int(sum)
// 	res.totalSamples = totalN
// 	res.Hoeffding()
// 	res.ComputeMean()
// 	return

// }
