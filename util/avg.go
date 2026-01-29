package util

import "sync"

type Float64_Avg struct {
	total float64
count int
}
type Float64_Avg_Sync struct {
	sync.Mutex
	total float64
count int
}

func (avg *Float64_Avg) Add(val float64) {
	avg.total+=val
	avg.count++
}

func (avg *Float64_Avg) GetAndReset() float64 {
	res := avg.total / float64(avg.count)
	avg.total=0
	avg.count=0
	return res
}

func (avg *Float64_Avg_Sync) Add(val float64) {
	avg.Lock()
	defer avg.Unlock()
	avg.total+=val
	avg.count++
}

func (avg *Float64_Avg_Sync) GetAndReset() float64 {
	avg.Lock()
	defer avg.Unlock()
	res := avg.total / float64(avg.count)
	avg.total=0
	avg.count=0
	return res
}