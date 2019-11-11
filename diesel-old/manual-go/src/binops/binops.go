package main

import "fmt"
import "math/rand"
import "time"
import ."dynfloats"

func main() {
  randSource := rand.NewSource(time.Now().UnixNano())
  randGen := rand.New(randSource)

  acc  := MakeDynFloat32(0.0)
  //half := MakeDynFloat32(0.5)
  iter := 0

  for ; acc.Delta < 0.001 ; iter++ {
    num := MakeDynFloat32(randGen.Float32()+float32(1.0))
    num.Delta = acc.Delta
    acc = AddDynFloat32(acc, num)
    //acc = MulDynFloat32(acc, half)
    fmt.Println(iter+1, acc.Delta)
  }
}
