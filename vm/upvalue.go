package vm

import "github.com/elk-language/elk/bitfield"

const (
	UpvalueLongIndexFlag bitfield.BitFlag8 = 1 << iota
	UpvalueLocalFlag
)
