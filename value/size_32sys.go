//go:build 386 || arm || armbe || mips || mipsle || ppc || s390 || sparc

package value

func (i Size) ToValue() Value {
	return Ref(i)
}
