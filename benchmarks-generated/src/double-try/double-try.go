package main

import (
  "fmt"
  "parallely"
)

/*
try {
  x = 2 [prob1] 3
  try {
    y = 4 [prob2] 5
  } recover {
    redo
  }
} recover {
  redo
}
*/

func main() {

  prob1 := float32(0.9)
  prob2 := float32(0.9)

  flag1 := false
  flag2 := false

  x := 0
  y := 1

  flag1 = false
  x = parallely.RandchoiceFlag(prob1, 2, 3, &flag1)
  flag2 = false
  y = parallely.RandchoiceFlag(prob2, 4, 5, &flag2)
  if flag2 {
    flag2 = false
    y = parallely.RandchoiceFlag(prob2, 4, 5, &flag2)
  }
  if flag1 {
    flag1 = false
    x = parallely.RandchoiceFlag(prob1, 2, 3, &flag1)
    flag2 = false
    y = parallely.RandchoiceFlag(prob2, 4, 5, &flag2)
    if flag2 {
      flag2 = false
      y = parallely.RandchoiceFlag(prob2, 4, 5, &flag2)
    }
  }

  fmt.Println(flag1,flag2,x,y)

}

