package main

import (
  "dieseldistacc"
)

const Num_threads = 11

func main() {
	dieseldistacc.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
	dieseldistacc.CleanupMain()
}
