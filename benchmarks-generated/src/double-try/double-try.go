package main

import (
  "fmt"
  "parallely"
)

/*
Code:

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

Reliability:

Inner block:
p_inner = 1-(1-prob2)^2

Outer block:
p_outer = 1-(1-(prob1*p_inner))^2
*/

func main() {

  //success probs of the two prob choice stmts
  prob1 := float32(0.9)
  prob2 := float32(0.9)

  var flag_outer, flag_inner bool //flags for the tcr blocks

  //garbage initial values
  x := 0
  y := 1

  //start of outer try: reset outer flag
  flag_outer = false
  x = parallely.RandchoiceFlag(prob1, 2, 3, &flag_outer)
  //start of inner try: reset inner flag
  flag_inner = false
  y = parallely.RandchoiceFlag(prob2, 4, 5, &flag_inner)
  //check inner flag - set if inner prob choice failed
  if flag_inner {
    //start of inner redo: reset inner flag
    flag_inner = false
    y = parallely.RandchoiceFlag(prob2, 4, 5, &flag_inner)
  }
  //if inner flag is set, then outer flag must be set too
  //because if inner tcr block fails EVEN AFTER it tries to recover, then outer tcr has failed too
  flag_outer = flag_outer || flag_inner
  if flag_outer {
    //start of outer redo: reset outer flag
    //remaining comments copy of above
    flag_outer = false
    x = parallely.RandchoiceFlag(prob1, 2, 3, &flag_outer)
    flag_inner = false
    y = parallely.RandchoiceFlag(prob2, 4, 5, &flag_inner)
    if flag_inner {
      flag_inner = false
      y = parallely.RandchoiceFlag(prob2, 4, 5, &flag_inner)
    }
    flag_outer = flag_outer || flag_inner
  }

  //print result
  fmt.Println(flag_outer,flag_inner,x,y)

}

