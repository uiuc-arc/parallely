package main

func GetIdx(row, col, cols int) int {
  return row*cols + col
}

const Dim = 100
const Bands = 10
const BandW = Dim/Bands
