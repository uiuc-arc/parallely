package main

import "diesel"

const NUMSENSORS = 64
const NUMP = 8
const NUMCENTERS = 8
const PERTHREAD = 8

const ITERATIONS = 10
var Q [NUMSENSORS] Process
var R [NUMP] Process

func main() {	
	diesel.LaunchProcess(Master)
	diesel.LaunchProcessGroup(Q[:], IoTDevice)
	diesel.LaunchProcessGroup(R[:], Worker)
}

type point struct {
	temperature float64 /*@dynamic*/
	humidity float64 /*@dynamic*/
}

func randomPic() [NUMCENTERS]point {
	var centers [NUMCENTERS] point
	for i, _ := range(centers) {
		centers[i] = point{1, 1}
	}
	return centers
}

func dist() float64 {
	return 1.0
}

func readHumidity() float64 {
	return 1.0
}

func readTemperature() float64 {
	return 1.0
}

func readTempError() (float64, float64) {
	return 0, 1.0
}

func readHumidError() (float64, float64) {
	return 0, 1.0
}

func AverageOverThreads(newcenters [NUMP][PERTHREAD] point) [NUMCENTERS] point {
	var centers [NUMCENTERS] point
	return centers
}

func Master() {
	var data [NUMSENSORS] point
	var centers [NUMCENTERS] point
	var newcenters [NUMP][PERTHREAD] point
	//Setting up data structs
	for i, IoTDevice := range(Q) {
		data[i] = diesel.Receive(IoTDevice)
	}
	
	centers = randomPic()
	
	for i, Worker := range(R) {
		diesel.Send(Worker, data)
	}
	
	for j:=0; j<ITERATIONS; j++ {
		for _, Worker := range(R) {
			diesel.Send(Worker, centers) }
		for i, Worker := range(R) {
			newcenters[i] = diesel.Receive(Worker)
		}
		centers = AverageOverThreads(newcenters)
	}
	checkArray(centers, 1, 0.99, 4, 0.99)
}

func IoTDevice(tid diesel.Process) {
	var /*@dynamic*/ temperature float64
	var humidity float64 /*@dynamic*/
	tempVal := readTemperature()
	tempErr, tempConf := readTempError()
	humidVal := readHumidity()
	humidErr, humidConf := readHumidError()
	temperature = diesel.Track(tempVal, tempErr, tempConf)
	humidity = diesel.Track(humidVal, humidErr, humidConf)
	diesel.Send(Master, point{temperature, humidity})
}

func Worker(tid diesel.Process) {
	var data [NUMSENSORS] point
	var centers [NUMCENTERS] point
	var newcenters [NUMCENTERS] point // asd
	var assign [PERTHREAD] int /*@dynamic*/
	
	data = diesel.Receive(Master)
	for iter:=0; iter<ITERATIONS; iter++ {
		centers = diesel.Receive(Master)
		for i:=0; i<PERTHREAD; i++ {
			mindist := MAX
			for c:=0; c<NUMCENTERS; c++ {
				dist := dist(centers[c], data[i])
				if dist < mindist {
					assign[i] = c
					mindist = dist
				} } }
		newcenters = avgAssigned(data, assign)
		send(Master, newcenters)
	}
}
