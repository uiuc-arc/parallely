package dynfloats

import "fmt"

func SendF32Arr(arr []float32, num int, chout chan float32) {
  for i:=0; i<num; i++ {
    chout <- arr[i]
  }
}

func RecvF32Arr(arr []float32, num int, chin chan float32) {
  for i:=0; i<num; i++ {
    arr[i] = <- chin
  }
}

func RecvF32ArrI(arr []float32, num int, chin chan float32) {
  receivedData := 0
  for i:=0; i<num; i++ {
    arr[i] = <- chin
    receivedData += 4
  }
  fmt.Println(receivedData)
}

func SendF32ArrAcc(arr []float32, darr []float64, num, doff int, chout chan float32, dchout chan float64, opt bool) {
  dmax := 0.0
  for i:=0; i<num; i++ {
    chout <- arr[i]
    if !opt {
      dchout <- darr[doff + i]
    } else if dmax < darr[doff + i] {
      dmax = darr[doff + i]
    }
  }
  if opt {
    dchout <- dmax
  }
}

func RecvF32ArrAcc(arr []float32, darr []float64, num, doff int, chin chan float32, dchin chan float64, opt bool) {
  for i:=0; i<num; i++ {
    arr[i] = <- chin
    if !opt {
      darr[doff + i] = <- dchin
    }
  }
  if opt {
    dmax := <- dchin
    for i:=0; i<num; i++ {
      darr[doff + i] = dmax
    }
  }
}

func RecvF32ArrAccI(arr []float32, darr []float64, num, doff int, chin chan float32, dchin chan float64, opt bool) {
  receivedData := 0
  for i:=0; i<num; i++ {
    arr[i] = <- chin
    receivedData += 4
    if !opt {
      darr[doff + i] = <- dchin
      receivedData += 8
    }
  }
  if opt {
    dmax := <- dchin
    receivedData += 8
    for i:=0; i<num; i++ {
      darr[doff + i] = dmax
    }
  }
  fmt.Println(receivedData)
}
