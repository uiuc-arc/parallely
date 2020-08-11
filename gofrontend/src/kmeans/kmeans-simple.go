package main;

import "fmt";

type process struct {
	Reliability float32;
	Delta       float64;
};

const NUMSENSORS = 64;
const NUMP = 8;
const ITERATIONS = 10;
var Q [64] process;
var R [8] process;

func main() {
	fmt.Println(NUMP);
}
