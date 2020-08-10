package main

import "diesel"

var NUMSENSORS = 64
var NUMP = 8
var ITERATIONS = 10
var Q [NUMSENSORS] process
var R [NUMP] process

func main() {
	diesel.LaunchProcess(Master)
	diesel.LaunchProcessGroup(Q, IoTDevice)
	diesel.LaunchProcessGroup(R, Worker)
}

type point struct {
	/*@dynamic*/ temperature float64
	/*@dynamic*/ humidity float64
}

func randomPic() {
	var centers [NUMCENTERS] point
	for i, _ := range(point) {
		centers[i] = point{1, 1}
	}
	return 
}

func dist() float64 {
	return 1.0
}

func Master() {
	var data [NUMSENSORS] point
	var centers [NUMCENTERS] point
	var newcenters [NUMP][PERTHREAD] point
	//Setting up data structs
	for i, IoTDevice := range(Q) {
		data[i] = receive(IoTDevice)
	}
	
	centers = randomPic()
	
	for i, Worker := range(R) {
		send(Worker, data)
	}
	
	for j:=0; j<ITERATIONS; j++ {
		for _, Worker := range(R) {
			send(Worker, centers) }
		for i, Worker := range(R) {
			newcenters[i] = receive(Worker)
		}
		centers = AverageOverThreads(newcenters)
	}
	checkArray(centers, 1, 0.99, 4, 0.99)
}

func IoTDevice() {
	/*@dynamic*/ var temperature float64
	/*@dynamic*/ var humidity float64
	tempVal := readTemperature()
	tempErr, tempConf := readTempError()
	humidVal := readHumidity()
	humidErr, humidConf := readHumidError()
	temperature = track(tempVal, tempErr, tempConf)
	humidity = track(humidVal, humidErr, humidConf)
	send(Master, point{temperature, humidity})
}

func Worker() {
	var data [NUMSENSORS] point
	var centers [NUMCENTERS] point
	var newcenters [NUMCENTERS] point
	/*@dynamic*/ var assign [PERTHREAD] int
	data = receive(Master)
	for iter:=0; iter<ITERATIONS; iter++ {
		centers = receive(Master)
		for i:=0; i<PERTHREAD; i++ {
			mindist := MAX
			for c:=0; c<NUMCENTERS; c++ {
				dist := dist()
				if dist < mindist {
					assign[i] = c
					mindist = dist
				} } }
		newcenters = avgAssigned(data, assign)
		send(Master, newcenters)
	}
}
