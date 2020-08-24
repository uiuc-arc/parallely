package main

import (
	"dieselacc"
	"fmt"
	"math"
	"math/rand"
)

const ArrayDim = 100
const ArraySize = 10000
const SliceSize = 1000
const Iterations = 10
const RowsPerThread = ArrayDim / 10

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
	defer dieselacc.Wg.Done()
	var DynMap [30004]float64
	var my_chan_index int
	_ = my_chan_index
	_ = DynMap
	var mC [10000]float64
	dieselacc.InitDynArray(0, 10000, DynMap[:])
	var mA [10000]float64
	dieselacc.InitDynArray(10000, 10000, DynMap[:])
	var mB [10000]float64
	dieselacc.InitDynArray(20000, 10000, DynMap[:])
	var tempDF0 float64
	DynMap[30000] = 0.0
	var tempDF1 float64
	DynMap[30001] = 0.0
	var tempDF2 float64
	DynMap[30002] = 0.0
	var rowIdx int
	var colIdx int
	var innerIdx int
	var outIdx int
	var sum float64
	DynMap[30003] = 0.0
	outIdx = 0
	rowIdx = 0

	dieselacc.StartTiming()

	for __temp_0 := 0; __temp_0 < 100; __temp_0++ {
		colIdx = 0
		for __temp_1 := 0; __temp_1 < ArrayDim; __temp_1++ {
			sum = 0.0
			// DynMap[30003] = 0.0
			innerIdx = 0
			for __temp_2 := 0; __temp_2 < ArrayDim; __temp_2++ {
				_temp_index_1 := (rowIdx * ArrayDim) + innerIdx
				tempDF0 = mA[_temp_index_1]
				// DynMap[30000] = DynMap[10000+_temp_index_1]

				_temp_index_2 := (innerIdx * ArrayDim) + colIdx
				tempDF1 = mB[_temp_index_2]
				// DynMap[30001] = DynMap[20000+_temp_index_2]

				// DynMap[30002] = math.Abs(float64(tempDF0))*DynMap[30000] + math.Abs(float64(tempDF1))*DynMap[30001] + DynMap[30000]*DynMap[30001]
				tempDF2 = tempDF0 * tempDF1
				// DynMap[30000] = DynMap[30003] + DynMap[30002]
				tempDF0 = sum + tempDF2
				// DynMap[30003] = DynMap[30000]
				DynMap[30003] = DynMap[30003] + tempDF0*DynMap[10000+_temp_index_1] + tempDF1*DynMap[20000+_temp_index_2] // + DynMap[10000+_temp_index_1]*DynMap[20000+_temp_index_2]
				sum = tempDF0
				innerIdx = innerIdx + 1
			}
			_temp_index_3 := outIdx
			mC[_temp_index_3] = sum
			DynMap[0+_temp_index_3] = DynMap[30003]
			colIdx = colIdx + 1
			outIdx = outIdx + 1
		}
		rowIdx = rowIdx + 1
	}

	dieselacc.EndTiming()

	fmt.Println("Ending thread : ", 0)
}

func main() {
	// rand.Seed(time.Now().UTC().UnixNano())
	seed := int64(12345)
	rand.Seed(seed) // deterministic seed for reproducibility

	A = make([]float64, ArraySize)
	B = make([]float64, ArraySize)

	fmt.Println("Generating matrices of size", ArraySize, "using random seed", seed)

	for i := 0; i < ArraySize; i++ {
		A[i] = rand.NormFloat64() * math.Abs(1.0)
		B[i] = rand.NormFloat64() * math.Abs(1.0)
	}

	fmt.Println("Starting program")

	dieselacc.InitChannels(1)

	go func_0()

	fmt.Println("Main thread waiting for others to finish")
	dieselacc.Wg.Wait()
}
