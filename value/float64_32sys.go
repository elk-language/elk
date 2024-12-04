//go:build 386 || arm || armbe || mips || mipsle || ppc || s390 || sparc

package value

func (f Float64) ToValue() Value {
	return Ref(f)
}
