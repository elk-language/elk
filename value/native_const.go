package value

type NativeConstant struct {
	ElkName string
	GoExpr  string
	GoType  *GoType
}

var NativeConstantMap = map[string]*NativeConstant{}

// Register a native Elk constant that lives in a Go global variable
func RegisterNativeConstant(elkName string, goExpr string, goType *GoType) {
	NativeConstantMap[elkName] = &NativeConstant{
		ElkName: elkName,
		GoExpr:  goExpr,
		GoType:  goType,
	}
}

func RegisterNativeClass(elkName, goName string) {
	RegisterNativeConstant(elkName, goName, NewGoType("*value.Class"))
}

func RegisterNativeMixin(elkName, goName string) {
	RegisterNativeConstant(elkName, goName, NewGoType("*value.Mixin"))
}

func RegisterNativeModule(elkName, goName string) {
	RegisterNativeConstant(elkName, goName, NewGoType("*value.Module"))
}

func RegisterNativeInterface(elkName, goName string) {
	RegisterNativeConstant(elkName, goName, NewGoType("*value.Interface"))
}
