package util

import (
	"encoding/json"
	"sync"
)

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
	avg.total += val
	avg.count++
}

func (avg *Float64_Avg) GetAndReset() float64 {
	res := avg.total / float64(avg.count)
	avg.total = 0
	avg.count = 0
	return res
}

func (avg *Float64_Avg_Sync) Add(val float64) {
	avg.Lock()
	defer avg.Unlock()
	avg.total += val
	avg.count++
}

func (avg *Float64_Avg_Sync) GetAndReset() float64 {
	avg.Lock()
	defer avg.Unlock()
	res := avg.total / float64(avg.count)
	avg.total = 0
	avg.count = 0
	return res
}

// Avg,Total,Count
func (avg *Float64_Avg) Get() (float64, float64, int) {
	return (avg.total / float64(avg.count)), avg.total, avg.count
}

// Avg,Total,Count
func (avg *Float64_Avg_Sync) Get() (float64, float64, int) {
	avg.Lock()
	defer avg.Unlock()
	return (avg.total / float64(avg.count)), avg.total, avg.count
}

func (avg *Float64_Avg) MarshalJSON() ([]byte, error) {
	return json.Marshal(avg.total / float64(avg.count))
}

func (avg *Float64_Avg_Sync) MarshalJSON() ([]byte, error) {
	avg.Lock()
	defer avg.Unlock()
	return json.Marshal(avg.total / float64(avg.count))
}
