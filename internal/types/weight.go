package types

import "fmt"

type Weight struct {
	Value *float64
}

func (w Weight) isNil() bool {
	return w.Value == nil
}

func (w Weight) Valid() bool {
	return !w.isNil()
}

func (w Weight) Float64() float64 {
	if w.isNil() {
		return 0
	}
	return *w.Value
}

func (w Weight) String() string {
	if w.isNil() {
		return ""
	}
	return fmt.Sprintf("%.1f", w.Float64())
}

func NewWeight(val float64) Weight {
	return Weight{Value: &val}
}
