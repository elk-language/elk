//go:build 386 || arm || armbe || mips || mipsle || ppc || s390 || sparc

package value

func (d Duration) ToValue() Value {
	return d
}
