//go:build 386 || arm || armbe || mips || mipsle || ppc || s390 || sparc

package value

import "math"

type UInt uint32

const MaxInt64ForUInt = math.MaxUint32
