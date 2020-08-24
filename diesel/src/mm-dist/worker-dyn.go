package main

import (
	"dieseldistacc"
	"fmt"
	"math"
	"os"
	"strconv"
)

const ArrayDim = 100
const ArraySize = 10000
const SliceSize = 1000
const Iterations = 10
const RowsPerThread = ArrayDim / 10

var Num_threads int

var A, B []float64

const (
	// single whitespace character
	ws = "[ \n\r\t\v\f]"
	// isolated comment
	cmt = "#[^\n\r]*"
	// comment sub expression
	cmts = "(" + ws + "*" + cmt + "[\n\r])"
	// number with leading comments
	num = "(" + cmts + "+" + ws + "*|" + ws + "+)([0-9]+)"
)

var Q = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func func_0() {
	dieseldistacc.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
	dieseldistacc.WaitForWorkers(Num_threads)
	var DynMap [31002]float64
	var my_chan_index int
	_ = my_chan_index
	_ = DynMap
	var mC [10000]float64
	dieseldistacc.InitDynArray(0, 10000, DynMap[:])
	var mA [10000]float32
	dieseldistacc.InitDynArray(10000, 10000, DynMap[:])
	var mB [10000]float32
	dieseldistacc.InitDynArray(20000, 10000, DynMap[:])
	var slice [1000]float32
	dieseldistacc.InitDynArray(30000, 1000, DynMap[:])
	var idx0 int
	var idx1 int
	var tempF64 float64
	var tempDF64 float64
	DynMap[31000] = 0
	var tempDF32 float32
	DynMap[31001] = 0
	dieseldistacc.StartTiming()
	idx0 = 0
	for __temp_0 := 0; __temp_0 < ArraySize; __temp_0++ {
		_temp_index_1 := idx0
		tempF64 = A[_temp_index_1]
		tempDF64 = tempF64
		DynMap[31000] = 0.0
		tempDF32 = float32(tempDF64)
		DynMap[31001] = dieseldistacc.GetCastingError64to32(tempDF64, tempDF32)
		_temp_index_2 := idx0
		mA[_temp_index_2] = tempDF32
		DynMap[10000+_temp_index_2] = DynMap[31001]
		_temp_index_3 := idx0
		tempF64 = B[_temp_index_3]
		tempDF64 = tempF64
		DynMap[31000] = 0.0
		tempDF32 = float32(tempDF64)
		DynMap[31001] = dieseldistacc.GetCastingError64to32(tempDF64, tempDF32)
		_temp_index_4 := idx0
		mB[_temp_index_4] = tempDF32
		DynMap[20000+_temp_index_4] = DynMap[31001]
		idx0 = idx0 + 1
	}
	for _, q := range Q {
		dieseldistacc.SendDynFloat32ArrayO1(mA[:], 0, q, DynMap[:], 10000)
		dieseldistacc.SendDynFloat32ArrayO1(mB[:], 0, q, DynMap[:], 20000)
	}
	idx0 = 0
	for _, q := range Q {
		dieseldistacc.ReceiveDynFloat32ArrayO1(slice[:], 0, q, DynMap[:], 30000)
		idx1 = 0
		for __temp_1 := 0; __temp_1 < SliceSize; __temp_1++ {
			_temp_index_5 := idx1
			tempDF32 = slice[_temp_index_5]
			DynMap[31001] = DynMap[30000+_temp_index_5]
			tempDF64 = float64(tempDF32)
			DynMap[31000] = DynMap[31001]
			_temp_index_6 := idx0
			mC[_temp_index_6] = tempDF64
			DynMap[0+_temp_index_6] = DynMap[31000]
			idx0 = idx0 + 1
			idx1 = idx1 + 1
		}
	}
	dieseldistacc.EndTiming()

	dieseldistacc.CleanupMain()
	fmt.Println("Ending thread : ", 0)
}

func func_Q(tid int) {
	dieseldistacc.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
	dieseldistacc.PingMain(tid)
	var DynMap [42005]float64
	var my_chan_index int
	_ = my_chan_index
	_ = DynMap
	q := tid
	var mA [10000]float64
	dieseldistacc.InitDynArray(0, 10000, DynMap[:])
	var mB [10000]float64
	dieseldistacc.InitDynArray(10000, 10000, DynMap[:])
	var slice [1000]float64
	dieseldistacc.InitDynArray(20000, 1000, DynMap[:])
	var mA32 [10000]float32
	dieseldistacc.InitDynArray(21000, 10000, DynMap[:])
	var mB32 [10000]float32
	dieseldistacc.InitDynArray(31000, 10000, DynMap[:])
	var slice32 [1000]float32
	dieseldistacc.InitDynArray(41000, 1000, DynMap[:])
	var myStartRow int
	var rowIdx int
	var colIdx int
	var innerIdx int
	var outIdx int
	var sum float64
	DynMap[42000] = 0
	var tempDF0 float64
	DynMap[42001] = 0
	var tempDF1 float64
	DynMap[42002] = 0
	var tempDF2 float64
	DynMap[42003] = 0
	var tempDF32 float32
	DynMap[42004] = 0
	myStartRow = (q - 1) * RowsPerThread
	dieseldistacc.ReceiveDynFloat32ArrayO1(mA32[:], tid, 0, DynMap[:], 21000)
	dieseldistacc.ReceiveDynFloat32ArrayO1(mB32[:], tid, 0, DynMap[:], 31000)
	outIdx = 0
	for __temp_2 := 0; __temp_2 < ArraySize; __temp_2++ {
		_temp_index_1 := outIdx
		tempDF32 = mA32[_temp_index_1]
		DynMap[42004] = DynMap[21000+_temp_index_1]
		tempDF0 = float64(tempDF32)
		DynMap[42001] = DynMap[42004]
		_temp_index_2 := outIdx
		mA[_temp_index_2] = tempDF0
		DynMap[0+_temp_index_2] = DynMap[42001]
		_temp_index_3 := outIdx
		tempDF32 = mB32[_temp_index_3]
		DynMap[42004] = DynMap[31000+_temp_index_3]
		tempDF0 = float64(tempDF32)
		DynMap[42001] = DynMap[42004]
		_temp_index_4 := outIdx
		mB[_temp_index_4] = tempDF0
		DynMap[10000+_temp_index_4] = DynMap[42001]
		outIdx = outIdx + 1
	}
	outIdx = 0
	rowIdx = myStartRow
	for __temp_3 := 0; __temp_3 < RowsPerThread; __temp_3++ {
		colIdx = 0
		for __temp_4 := 0; __temp_4 < ArrayDim; __temp_4++ {
			sum = 0.0
			innerIdx = 0
			for __temp_5 := 0; __temp_5 < ArrayDim; __temp_5++ {
				_temp_index_5 := (rowIdx * ArrayDim) + innerIdx
				tempDF0 = mA[_temp_index_5]
				// DynMap[42001] = DynMap[0+_temp_index_5]

				_temp_index_6 := (innerIdx * ArrayDim) + colIdx
				tempDF1 = mB[_temp_index_6]
				 // DynMap[42002] = DynMap[10000+_temp_index_6]

				 // DynMap[42003] = math.Abs(float64(tempDF0))*DynMap[0+_temp_index_5] +
	// 	math.Abs(float64(tempDF1))*DynMap[10000+_temp_index_6] +
	// 	DynMap[0+_temp_index_5]*DynMap[10000+_temp_index_6]
				tempDF2 = tempDF0 * tempDF1

				// DynMap[42001] = DynMap[42000] + math.Abs(float64(tempDF0))*DynMap[0+_temp_index_5] +
				// 	math.Abs(float64(tempDF1))*DynMap[10000+_temp_index_6] +
				// 	DynMap[0+_temp_index_5]*DynMap[10000+_temp_index_6]
				tempDF0 = sum + tempDF2

				DynMap[42000] = DynMap[42000] + math.Abs(float64(tempDF0))*DynMap[0+_temp_index_5] +
					math.Abs(float64(tempDF1))*DynMap[10000+_temp_index_6] +
					DynMap[0+_temp_index_5]*DynMap[10000+_temp_index_6]
				sum = tempDF0
				innerIdx = innerIdx + 1
			}
			_temp_index_7 := outIdx
			slice[_temp_index_7] = sum
			DynMap[20000+_temp_index_7] = DynMap[42000]
			colIdx = colIdx + 1
			outIdx = outIdx + 1
		}
		rowIdx = rowIdx + 1
	}
	outIdx = 0
	for __temp_6 := 0; __temp_6 < SliceSize; __temp_6++ {
		_temp_index_8 := outIdx
		tempDF0 = slice[_temp_index_8]
		DynMap[42001] = DynMap[20000+_temp_index_8]
		tempDF32 = float32(tempDF0)
		DynMap[42004] = dieseldistacc.GetCastingError64to32(tempDF0, tempDF32)
		_temp_index_9 := outIdx
		slice32[_temp_index_9] = tempDF32
		DynMap[41000+_temp_index_9] = DynMap[42004]
		outIdx = outIdx + 1
	}
	dieseldistacc.SendDynFloat32ArrayO1(slice32[:], tid, 0, DynMap[:], 41000)

	fmt.Println("Ending thread : ", q)
}

func main() {
	// rand.Seed(time.Now().UTC().UnixNano())
	tid, _ := strconv.Atoi(os.Args[1])
	_ = math.Inf

	Num_threads = 11

	func_Q(tid)

}
