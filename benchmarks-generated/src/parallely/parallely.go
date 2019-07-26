package parallely

import "math/rand"
import "fmt"
import "sync"
import "time"

var approxChannelMap map[int] chan int
var preciseChannelMap map[int] chan int

var preciseChannelMapFloat map[int] chan int
var preciseChannelMapFloat64Array map[int] chan []float64
var preciseChannelMapFloat32Array map[int] chan []float32

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

func InitChannels(numprocesses_in int){
	rand.Seed(time.Now().UTC().UnixNano())
	// var temp_approxChannelMap map[int] chan int
	numprocesses = numprocesses_in 
	Wg.Add(numprocesses_in)
	if debug==1 {
		fmt.Printf("Starting the wait group for %d threads", numprocesses_in)
	}
	
	approxChannelMap = make(map[int] chan int)
	preciseChannelMap = make(map[int] chan int)
	preciseChannelMapFloat64Array = make(map[int] chan []float64)
	preciseChannelMapFloat32Array = make(map[int] chan []float32)
	
	for i := 0; i < numprocesses_in * numprocesses_in; i++ {
		approxChannelMap[i] = make(chan int, 100)
		preciseChannelMap[i] = make(chan int, 100)
		preciseChannelMapFloat64Array[i] = make(chan []float64)
		preciseChannelMapFloat32Array[i] = make(chan []float32)
	}

	if debug==1 {
		fmt.Println("Initialized Channels : ", len(approxChannelMap))
	}
}

func Send(value, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	preciseChannelMap[my_chan_index] <- value
	if debug==1 {
		fmt.Printf("%d Sending message in precise int chan : %d (%d * %d + %d)\n",
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

func Receive(rec_var *int, receiver, sender int) {
	my_chan_index := sender * numprocesses + receiver
	temp_rec_val := <- preciseChannelMap[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
	}
	*rec_var = temp_rec_val
}

func ReceiveFloat64Array(rec_var []float64, receiver, sender int) {
	my_chan_index := sender * numprocesses + receiver
	temp_rec_val := <- preciseChannelMapFloat64Array[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise float64 array chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
	}
	rec_var = temp_rec_val
}

func ReceiveFloat32Array(rec_var []float32, receiver, sender int) {
	my_chan_index := sender * numprocesses + receiver
	temp_rec_val := <- preciseChannelMapFloat32Array[my_chan_index]
	if debug==1 {
		fmt.Printf("%d Received message in precise float32 array chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
	}
	rec_var = temp_rec_val
}

func Condsend(cond, value, sender, receiver int) {
	my_chan_index := sender * numprocesses + receiver
	if debug==1 {
		fmt.Printf("%d Sending message in approx chan : %d (%d * %d + %d)\n", sender,
			my_chan_index, sender, numprocesses, receiver);
	}
	if cond != 0 {
		approxChannelMap[my_chan_index] <- value
	} else {
		if debug==1 {
			fmt.Printf("[Failure %d] %d Sending message in approx chan : %d (%d * %d + %d)\n", cond, sender,
				my_chan_index, sender, numprocesses, receiver);
		}
		approxChannelMap[my_chan_index] <- -1
	}
}

func Condreceive(rec_cond_var, rec_var *int, receiver, sender int) {
	my_chan_index := sender * numprocesses + receiver

	if debug==1 {
		fmt.Printf("---- %d Waiting to Receive from approx int chan : %d (%d * %d + %d)\n",
			receiver, my_chan_index, sender, numprocesses, receiver);
	}
	
	temp_rec_val := <- approxChannelMap[my_chan_index]

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
	
