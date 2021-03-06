package main

import (
	"regexp"
	"io"
	"io/ioutil"
	"errors"
	"strconv"
	"os"
	"math"
	"fmt"
)

const (
	// single whitespace character
	ws = "[ \n\r\t\v\f]"
	// isolated comment
	cmt = "#[^\n\r]*"
	// comment sub expression
	cmts = "(" + ws + "*" + cmt + "[\n\r])"
	// number with leading comments
	num = "(" + cmts + "+" + ws + "*|" + ws + "+)([0-9]+)"
)

var rxHeader = regexp.MustCompile("^P6" + num + num + num + "(" + cmts + "*" + ")" + ws)
var rxComment = regexp.MustCompile(cmt)

// ReadFrom constructs a Bitmap object from an io.Reader.
func ReadPpmFrom(r io.Reader) (b []int, x int, y int, err error) {
	var all []byte
	all, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}
	bss := rxHeader.FindSubmatch(all)
	if bss == nil {
		return nil, 0, 0, errors.New("unrecognized ppm header")
	}
	x, _ = strconv.Atoi(string(bss[3]))
	y, _ = strconv.Atoi(string(bss[6]))
	maxval, _ := strconv.Atoi(string(bss[9]))
	if maxval > 255 {
		return nil, 0, 0, errors.New("16 bit ppm not supported")
	}
	b = make([]int, x*y)

	b3 := all[len(bss[0]):]
	var n1 int
	for i := range b {
		b[i] = int(b3[n1]) * 255 / maxval
		n1 += 3
	}
	return
}

// ReadFile writes binary P6 format PPM from the specified filename.
func ReadPpmFile(fn string) (b []int, x int, y int, err error) {
	var f *os.File
	if f, err = os.Open(fn); err != nil {
		return
	}
	if b, x, y, err = ReadPpmFrom(f); err != nil {
		return
	}
	return b, x, y, f.Close()
}

// WriteTo outputs 8-bit P6 PPM format to an io.Writer.
func WritePpmTo(b []int, x int, y int, w io.Writer) (err error) {
	// magic number
	if _, err = fmt.Fprintln(w, "P6"); err != nil {
		return
	}
	// x, y, depth
	_, err = fmt.Fprintf(w, "%d %d\n255\n", x, y)
	if err != nil {
		return
	}
	// raster data in a single write
	b3 := make([]byte, 3*len(b))
	n1 := 0
	for _, px := range b {
		b3[n1] = byte(px)
		b3[n1+1] = byte(px)
		b3[n1+2] = byte(px)
		n1 += 3
	}
	if _, err = w.Write(b3); err != nil {
		return
	}
	return
}

// WriteFile writes to the specified filename.
func WritePpmFile(b []int, x int, y int, fn string) (err error) {
	var f *os.File
	if f, err = os.Create(fn); err != nil {
		return
	}
	if err = WritePpmTo(b, x, y, f); err != nil {
		return
	}
	return f.Close()
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Idx(i, j, width int) int {
	return i*width + j
}

func floorInt(input float64) int {
	return int(math.Floor(input))
}

func ceilInt(input float64) int {
	return int(math.Ceil(input))
}

func convertToFloat(x int) float64 {
	return float64(x)
}
