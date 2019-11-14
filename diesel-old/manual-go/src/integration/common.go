package main

const Max = float32(11)
const Min = float32(1)
const Delta = float32(0.00390625)
const Threads = 10
const IntervalPerThread = (Max-Min)/float32(Threads)
const DivisionsPerThread = int(IntervalPerThread/Delta)
