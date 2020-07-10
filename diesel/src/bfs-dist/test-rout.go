package main

import (  
    "fmt"
    "time"
)

func hello() {  
    fmt.Println("Hello world goroutine")
}

func main() {  
    go hello()
        time.Sleep(3000 * time.Millisecond)
    fmt.Println("main function")
}