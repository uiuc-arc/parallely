package main

import (
  "os"
  "strconv"
  "fmt"
  "encoding/gob"
)

func main() {
  height, _ := strconv.Atoi(os.Args[1])
  width, _ := strconv.Atoi(os.Args[2])

  original := []float64{}
  precred := []float64{}

  fo, _ := os.Open("/tmp/original.dat")
  odecoder := gob.NewDecoder(fo)
  odecoder.Decode(&original)
  fp, _ := os.Open("/tmp/precred.dat")
  pdecoder := gob.NewDecoder(fp)
  pdecoder.Decode(&precred)

  acc := 0.0
  for i := 0; i < height*width; i++ {
    diff := original[i]-precred[i]
    acc += diff*diff
  }

  fmt.Println(acc/float64(height*width))
}
