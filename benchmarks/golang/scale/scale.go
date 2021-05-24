package main

import (
	"diesel"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"
)

var Src [262144]float64
var Dest [4194304]float64
var ImgSize int
var NumThreads int
var SWidth int
var SHeight int
var DestSize int

var Q = []int{1, 2, 3, 4, 5, 6, 7, 8}

func func_0() {
	// defer diesel.Wg.Done()

	var s_height int
	var s_width int
	var d_height int
	var d_width int
	var t_height int
	var ts_height int
	var i int
	var j int
	var k int
	var myrows int
	var lastthread int
	var te_height int
	var dest_slice [524288]float64
	var outImage [4194304]float64
	var temp float64
	s_width = SWidth
	s_height = SHeight
	d_width = 4 * s_width
	d_height = 4 * s_height
	t_height = d_height / (NumThreads - 1)
	i = 0
	for _, q := range Q {
		send(q, s_height)
		send(q, s_width)
		send(q, d_width)
		ts_height = i * t_height
		lastthread = (i == (NumThreads - 1))
		if lastthread != 0 {
			te_height = d_height
		} else {
			te_height = (i + 1) * t_height
		}
		send(q, ts_height)
		send(q, te_height)
		i = i + 1
	}
	i = 0
	for _, q := range Q {
		_, dest_slice = cond-receive(0)
		ts_height = i * t_height
		lastthread = (i == (NumThreads - 1))
		if lastthread != 0 {
			te_height = d_height
		} else {
			te_height = (i + 1) * t_height
		}
		myrows = te_height - ts_height
		for j := 0; j < myrows; j++ /*maxiterations=10*/ {
			k = 0
			for k := 0; k < d_width; k++ /*maxiterations=10*/ {
				tempi = j*d_width + k
				temp = dest_slice[tempi]
				tempi = (ts_height+j)*d_width + k
				outImage[tempi] = temp
			}
		}
		i = i + 1
	}
	Dest = outImage
}

func func_Q(q int) {
	// defer diesel.Wg.Done()
	
	var image [262144]float64
	var dest [524288]float64
	var ts_height int
	var i int
	var j int
	var myrows int
	var si float64
	var sj float64
	var delta float64
	var s_height int
	var s_width int
	var te_height int
	var d_width int
	var cond int
	var previ int
	var prevj int
	var nexti int
	var nextj int
	var ul float64
	var ll float64
	var ur float64
	var lr float64
	var u_w float64
	var l_w float64
	var ul_w float64
	var ll_w float64
	var ur_w float64
	var lr_w float64
	var tempf float64
	var tempf1 float64
	var tempi int
	image = Src

	s_height = receive(0)
	s_width = receive(0)
	d_width = receive(0)
	ts_height = receive(0)
	te_height = receive(0)
	
	myrows = te_height - ts_height

	delta = 1 / 4.0
	tempf = convertToFloat(ts_height)
	si = tempf * delta
	for i := 0; i < myrows; i++ /*maxiterations=10*/ {
		sj = 0.0
		for j := 0; j < d_width; j++ /*maxiterations=10*/ {
			previ = floorInt(si)
			nexti = ceilInt(si)
			prevj = floorInt(sj)
			nextj = ceilInt(sj)
			cond = (s_height <= nexti)
			if cond != 0 {
				previ = s_height - 2
				nexti = s_height - 1
			}
			cond = (s_width <= nextj)
			if cond != 0 {
				prevj = s_width - 2
				nextj = s_width - 1
			}
			tempi = previ*s_width + prevj
			ul = image[tempi]
			tempi = nexti*s_width + prevj
			ll = image[tempi]
			tempi = previ*s_width + nextj
			ur = image[tempi]
			tempi = nexti*s_width + nextj
			lr = image[tempi]
			tempf = convertToFloat(nexti)
			u_w = tempf - si
			tempf = convertToFloat(nextj)
			l_w = tempf - sj
			ul_w = u_w*l_w
			ll_w = (1.0-u_w)*l_w
			ur_w = u_w*(1.0-l_w)
			lr_w = (1.0-u_w)*(1.0-l_w)
			tempf1 = ul*ul_w
			tempf1 = tempf1+ur*ur_w
			tempf1 = tempf1+ll*ll_w
			tempf1 = tempf1+lr*lr_w
			
			tempi = i*d_width + j
			dest[tempi] = tempf1
			
			sj = sj + delta
		}
		si = si + delta
	}
	tempi = pchoice(1, 0, 0.99)
	cond-send(tempi, 0, dest)
}

func main() {
	// rand.Seed(time.Now().UTC().UnixNano())

	fmt.Println("Starting main thread")
	NumThreads = 9

	diesel.InitChannels(9)

	iFile := os.Args[1]
	oFile := os.Args[2]

	src_tmp, s_width, s_height, _ := ReadPpmFile(iFile)
	SHeight = s_height
	SWidth = s_width
	DestSize = len(src_tmp) * 4 * 4

	for i, _ := range src_tmp {
		Src[i] = float64(src_tmp[i])
	}

	ImgSize = len(src_tmp)

	startTime := time.Now()
	parallely.LaunchThread(0, func_0)
	parallely.LaunchThreadGroup(Q, func_Q, "q")

	fmt.Println("Main thread waiting for others to finish")
	diesel.Wg.Wait()

	end := time.Now()
	elapsed := end.Sub(startTime)
	fmt.Println("Elapsed time :", elapsed.Nanoseconds())

	tmp_dest := make([]int, len(Dest))
	for i, _ := range Dest {
		tmp_dest[i] = int(Dest[i])
	}

	WritePpmFile(tmp_dest, s_width*4, s_height*4, oFile)
}
