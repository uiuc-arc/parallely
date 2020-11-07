package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var Num_threads int

func max(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func convertToFloat(x int) float64 {
	return float64(x)
}

// __GLOBAL_DECS__

// __FUNC_DECS__

func parseOutput(outstr string) (result int, confidence float64) {
	r := regexp.MustCompile(`Prediction: .*`)
	matches := r.FindAllString(outstr, 1)
	if len(matches) == 0 {
		fmt.Println("could not read the output")
		return -1, -1.0
	}
	outparts := strings.Fields(matches[0])
	cat, err1 := strconv.Atoi(outparts[1])
	conf, err2 := strconv.ParseFloat(outparts[2], 64)
	if err1 != nil || err2 != nil {
		fmt.Println("could not read the output")
		return -1, -1.0
	}
	return cat, conf
}

func readCamera() (result int, confidence float64) {
	cmd := exec.Command("python3", "mio_inference_single.py")
	cmd.Dir = "./CNN_MIO_KERAS/"

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running the python script")
	}

	return parseOutput(string(out))
}

func main() {
	fmt.Println("Starting main thread")

	// Num_threads = __NUM_THREADS__

	// diesel.InitChannels(__NUM_THREADS__)

	// fmt.Println("Starting the iterations")
	// startTime := time.Now()

	// __START__THREADS__

	// for i := 0; i < 100; i++ {
	cmd := exec.Command("python3", "mio_inference_single.py")
	cmd.Dir = "./CNN_MIO_KERAS/"

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running the python script")
	}

	// matched, err := regexp.MatchString("Prediction: .*", out)
	// if err != nil {
	// 	fmt.Println("Error matching regex")
	// }

	fmt.Println(parseOutput(string(out)))
	// }

	// fmt.Println("Main thread waiting for others to finish")
	// diesel.Wg.Wait()
	// elapsed := time.Since(startTime)

	// fmt.Println("Done!")
	// fmt.Println("Elapsed time : ", elapsed.Nanoseconds())

	// f, _ := os.Create("output.txt")
	// defer f.Close()

	// for i := range DistGlobal {
	// 	f.WriteString(fmt.Sprintln(DistGlobal[i]))
	// }
}
