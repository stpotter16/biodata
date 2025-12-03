package types

import "fmt"

type Waist struct {
	Value *float64
}

func (w Waist) isNil() bool {
	return w.Value == nil
}

func (w Waist) Valid() bool {
	return !w.isNil()
}

func (w Waist) Float64() float64 {
	if w.isNil() {
		return 0
	}
	return *w.Value
}

func (w Waist) String() string {
	if w.isNil() {
		return ""
	}
	return fmt.Sprintf("%.1f", w.Float64())
}

func NewWaist(val float64) Waist {
	return Waist{Value: &val}
}
