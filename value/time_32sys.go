//go:build 386 || arm || armbe || mips || mipsle || ppc || s390 || sparc

package value

func (t Time) ToValue() Value {
	return Ref(t)
}
