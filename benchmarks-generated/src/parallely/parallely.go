package parallely

import "math/rand"
import "fmt"
import "sync"
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


var Wg sync.WaitGroup
var numprocesses int

var debug int = 0

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
	numprocesses = numprocesses_in 
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
	

	
	for i := 0; i < numprocesses_in * numprocesses_in; i++ {
		preciseChannelMapInt[i] = make(chan int)
		preciseChannelMapInt32[i] = make(chan int32)
		preciseChannelMapInt64[i] = make(chan int64)
		preciseChannelMapIntArray[i] = make(chan []int)
		preciseChannelMapInt32Array[i] = make(chan []int32)
		preciseChannelMapInt64Array[i] = make(chan []int64)

		preciseChannelMapFloat64[i] = make(chan float64)
		preciseChannelMapFloat32[i] = make(chan float32)
		preciseChannelMapFloat64Array[i] = make(chan []float64)
		preciseChannelMapFloat32Array[i] = make(chan []float32)

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
	}

	if debug==1 {
		fmt.Println("Initialized Channels : ", len(approxChannelMapInt) * 20)
	}
}

func SendInt(value, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	preciseChannelMapInt[my_chan_index] <- value
	if debug==1 {
		fmt.Printf("%d Sending message in precise int chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, numprocesses, receiver);
	}
}

func SendInt32(value int32, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	preciseChannelMapInt32[my_chan_index] <- value
	if debug==1 {
		fmt.Printf("%d Sending message in precise int32 chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, numprocesses, receiver);
	}
}

func SendInt64(value int64, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	preciseChannelMapInt64[my_chan_index] <- value
	if debug==1 {
		fmt.Printf("%d Sending message in precise int64 chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, numprocesses, receiver);
	}
}

func SendFloat32(value float32, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	preciseChannelMapFloat32[my_chan_index] <- value
	if debug==1 {
		fmt.Printf("%d Sending message in precise float32 chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, numprocesses, receiver);
	}
}

func SendFloat64(value float64, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	preciseChannelMapFloat64[my_chan_index] <- value
	if debug==1 {
		fmt.Printf("%d Sending message in precise float64 chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, numprocesses, receiver);
	}
}

func SendApprox(value, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	approxChannelMapInt[my_chan_index] <- value
	if debug==1 {
		fmt.Printf("%d Sending message in approx int chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, numprocesses, receiver);
	}
}

func SendInt32Approx(value int32, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	approxChannelMapInt32[my_chan_index] <- value
	if debug==1 {
		fmt.Printf("%d Sending message in approx int32 chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, numprocesses, receiver);
	}
}

func SendInt64Approx(value int64, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	approxChannelMapInt64[my_chan_index] <- value
	if debug==1 {
		fmt.Printf("%d Sending message in approx int64 chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, numprocesses, receiver);
	}
}

func SendIntArray(value []int, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	temp_array := make([]int, len(value))
	copy(temp_array, value)
	preciseChannelMapIntArray[my_chan_index] <- temp_array
	if debug==1 {
		fmt.Printf("%d Sending message in precise float32 array chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, numprocesses, receiver);
	}
}

func SendIntArrayApprox(value []int, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	temp_array := make([]int, len(value))
	copy(temp_array, value)
	approxChannelMapIntArray[my_chan_index] <- temp_array
	if debug==1 {
		fmt.Printf("%d Sending message in approx int array chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, numprocesses, receiver);
	}
}

func SendInt32Array(value []int32, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	temp_array := make([]int32, len(value))
	copy(temp_array, value)
	preciseChannelMapInt32Array[my_chan_index] <- temp_array
	if debug==1 {
		fmt.Printf("%d Sending message in precise float32 array chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, numprocesses, receiver);
	}
}

func SendInt32ArrayApprox(value []int32, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	temp_array := make([]int32, len(value))
	copy(temp_array, value)
	approxChannelMapInt32Array[my_chan_index] <- temp_array
	if debug==1 {
		fmt.Printf("%d Sending message in precise float32 array chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, numprocesses, receiver);
	}
}

func SendFloat64Array(value []float64, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	temp_array := make([]float64, len(value))
	copy(temp_array, value)
	preciseChannelMapFloat64Array[my_chan_index] <- temp_array
	if debug==1 {
		fmt.Printf("%d Sending message in precise float64 array chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, numprocesses, receiver);
	}
}

func SendFloat32Array(value []float32, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	temp_array := make([]float32, len(value))
	copy(temp_array, value)
	preciseChannelMapFloat32Array[my_chan_index] <- temp_array
	if debug==1 {
		fmt.Printf("%d Sending message in precise float32 array chan : %d (%d * %d + %d)\n",
			sender, my_chan_index, sender, numprocesses, receiver);
	}
}

func ReceiveInt(rec_var *int, receiver, sender int) {
	my_chan_index := sender * numprocesses + receiver
	temp_rec_val := <- preciseChannelMapInt[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
	}
	*rec_var = temp_rec_val
}

func ReceiveInt32(rec_var *int32, receiver, sender int) {
	my_chan_index := sender * numprocesses + receiver
	temp_rec_val := <- preciseChannelMapInt32[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
	}
	*rec_var = temp_rec_val
}

func ReceiveInt64(rec_var *int64, receiver, sender int) {
	my_chan_index := sender * numprocesses + receiver
	temp_rec_val := <- preciseChannelMapInt64[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
	}
	*rec_var = temp_rec_val
}

func ReceiveFloat32(rec_var *float32, receiver, sender int) {
	my_chan_index := sender * numprocesses + receiver
	temp_rec_val := <- preciseChannelMapFloat32[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise float32 chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
	}
	*rec_var = temp_rec_val
}

func ReceiveFloat64(rec_var *float64, receiver, sender int) {
	my_chan_index := sender * numprocesses + receiver
	temp_rec_val := <- preciseChannelMapFloat64[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise float64 chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
	}
	*rec_var = temp_rec_val
}

func ReceiveIntApprox(rec_var *int, receiver, sender int) {
	my_chan_index := sender * numprocesses + receiver
	temp_rec_val := <- approxChannelMapInt[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
	}
	*rec_var = temp_rec_val
}

func ReceiveInt32Approx(rec_var *int32, receiver, sender int) {
	my_chan_index := sender * numprocesses + receiver
	temp_rec_val := <- approxChannelMapInt32[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
	}
	*rec_var = temp_rec_val
}

func ReceiveInt64Approx(rec_var *int64, receiver, sender int) {
	my_chan_index := sender * numprocesses + receiver
	temp_rec_val := <- approxChannelMapInt64[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
	}
	*rec_var = temp_rec_val
}

func ReceiveIntArray(rec_var []int, receiver, sender int) {
	my_chan_index := sender * numprocesses + receiver
	temp_rec_val := <- preciseChannelMapIntArray[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise int array chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
	}
	copy(rec_var, temp_rec_val)
}

func ReceiveInt32Array(rec_var []int32, receiver, sender int) {
	my_chan_index := sender * numprocesses + receiver
	temp_rec_val := <- preciseChannelMapInt32Array[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise int32 array chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
	}
	copy(rec_var, temp_rec_val)
}

func ReceiveFloat64Array(rec_var []float64, receiver, sender int) {
	my_chan_index := sender * numprocesses + receiver
	temp_rec_val := <- preciseChannelMapFloat64Array[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise float64 array chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
	}
	copy(rec_var, temp_rec_val)
}

func ReceiveFloat32Array(rec_var []float32, receiver, sender int) {
	my_chan_index := sender * numprocesses + receiver
	temp_rec_val := <- preciseChannelMapFloat32Array[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise float32 array chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
	}
	copy(rec_var, temp_rec_val)
}

func Condsend(cond, value, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	if debug==1 {
		fmt.Printf("%d Sending message in approx int chan : %d (%d * %d + %d)\n", sender,
			my_chan_index, sender, numprocesses, receiver);
	}
	if cond != 0 {
		approxChannelMapInt[my_chan_index] <- value
	} else {
		if debug==1 {
			fmt.Printf("[Failure %d] %d Sending message in approx int chan : %d (%d * %d + %d)\n", cond, sender,
				my_chan_index, sender, numprocesses, receiver);
		}
		approxChannelMapInt[my_chan_index] <- -1
	}
}

func CondsendIntArray(cond int, value []int, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	if debug==1 {
		fmt.Printf("%d Sending message in approx int chan : %d (%d * %d + %d)\n", sender,
			my_chan_index, sender, numprocesses, receiver);
	}
	if cond != 0 {
		temp_array := make([]int, len(value))
		copy(temp_array, value)
		approxChannelMapIntArray[my_chan_index] <- temp_array
	} else {
		if debug==1 {
			fmt.Printf("[Failure %d] %d Cond Sending message in approx int chan : %d (%d * %d + %d)\n", cond, sender,
				my_chan_index, sender, numprocesses, receiver);
		}
		approxChannelMapIntArray[my_chan_index] <- []int{}
	}
}

func CondsendInt32(cond, value int32, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	if debug==1 {
		fmt.Printf("%d Sending message in approx int32 chan : %d (%d * %d + %d)\n", sender,
			my_chan_index, sender, numprocesses, receiver);
	}
	if cond != 0 {
		approxChannelMapInt32[my_chan_index] <- value
	} else {
		if debug==1 {
			fmt.Printf("[Failure %d] %d Sending message in approx int32 chan : %d (%d * %d + %d)\n", cond, sender,
				my_chan_index, sender, numprocesses, receiver);
		}
		approxChannelMapInt32[my_chan_index] <- -1
	}
}

func CondsendInt64(cond, value int64, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	if debug==1 {
		fmt.Printf("%d Sending message in approx int64 chan : %d (%d * %d + %d)\n", sender,
			my_chan_index, sender, numprocesses, receiver);
	}
	if cond != 0 {
		approxChannelMapInt64[my_chan_index] <- value
	} else {
		if debug==1 {
			fmt.Printf("[Failure %d] %d Sending message in approx int64 chan : %d (%d * %d + %d)\n", cond, sender,
				my_chan_index, sender, numprocesses, receiver);
		}
		approxChannelMapInt64[my_chan_index] <- -1
	}
}

func CondsendFloat32(cond, value float32, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	if debug==1 {
		fmt.Printf("%d Sending message in approx chan : %d (%d * %d + %d)\n", sender,
			my_chan_index, sender, numprocesses, receiver);
	}
	if cond != 0 {
		approxChannelMapFloat32[my_chan_index] <- value
	} else {
		if debug==1 {
			fmt.Printf("[Failure %d] %d Sending message in approx chan : %d (%d * %d + %d)\n", cond, sender,
				my_chan_index, sender, numprocesses, receiver);
		}
		approxChannelMapFloat32[my_chan_index] <- -1
	}
}

func CondsendFloat64(cond, value float64, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	if debug==1 {
		fmt.Printf("%d Sending message in approx chan : %d (%d * %d + %d)\n", sender,
			my_chan_index, sender, numprocesses, receiver);
	}
	if cond != 0 {
		approxChannelMapFloat64[my_chan_index] <- value
	} else {
		if debug==1 {
			fmt.Printf("[Failure %d] %d Sending message in approx chan : %d (%d * %d + %d)\n", cond, sender,
				my_chan_index, sender, numprocesses, receiver);
		}
		approxChannelMapFloat64[my_chan_index] <- -1
	}
}

func Condreceive(rec_cond_var, rec_var *int, receiver, sender int) {
	my_chan_index := sender * numprocesses + receiver

	if debug==1 {
		fmt.Printf("---- %d Waiting to Receive from approx int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
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
	my_chan_index := sender * numprocesses + receiver

	if debug==1 {
		fmt.Printf("---- %d Waiting to Receive from approx int array chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
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
	my_chan_index := sender * numprocesses + receiver

	if debug==1 {
		fmt.Printf("---- %d Waiting to Receive from approx int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
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


func ConvBool(x bool) int {
	if x {
		return 1;
	} else {
		return 0;
	}
}
	