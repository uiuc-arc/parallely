package main

import "time"

const rows = 100
const cols = 100
const bands = 10
const bandw = rows/bands
const iterations = 1

var seed = 0*time.Now().UnixNano() + 12345

func GetIdx(row, col, cols int) int {
  return row*cols + col
}

func max(array []float64) float64 {
  m := array[0]
  for i := range array {
    if m < array[i] {
      m = array[i]
    }
  }
  return m
}

func fill(array []float64, val float64) {
  for i := range array {
    array[i] = val
  }
}
