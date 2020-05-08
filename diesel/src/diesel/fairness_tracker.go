// +build !instrument

package diesel

//import "math/rand"
//import "fmt"
import "math"

type BooleanTracker struct {
	successes int 
	totalSamples int 
	mean float64 
	delta float64 
	eps float64 
}


func (b *BooleanTracker) NewBooleanTracker(){
	b.successes = 0
	b.totalSamples  = 0
	b.mean  = 0
	b.delta  = 0.1
	b.eps  = 1
}

func (b *BooleanTracker) addSample(samp int) {
    b.successes = b.successes + samp
    b.totalSamples = b.totalSamples + 1
    b.Hoeffding()
    b.GetMean()
}

func (b *BooleanTracker) Hoeffding() {
	b.eps = math.Sqrt((0.6*math.Log((math.Log(float64(1.1*float64(b.totalSamples+1)))/math.Log(1.10)))+0.555*math.Log(24/b.delta))/float64(b.totalSamples+1))
}

func (b *BooleanTracker) GetMean(){
	b.mean = float64(b.successes)/float64(b.totalSamples)
}

