package types

import "fmt"

type BP struct {
	Systolic  *float64
	Diastolic *float64
}

func (bp BP) isNil() bool {
	return bp.Systolic == nil || bp.Diastolic == nil
}

func (bp BP) Valid() bool {
	return !bp.isNil()
}

func (bp BP) String() string {
	if bp.isNil() {
		return ""
	}
	return fmt.Sprintf("%g/%g", *bp.Systolic, *bp.Diastolic)
}
