package runtime

func InitGlobalEnvironment() {
	initLexer()
}

func init() {
	InitGlobalEnvironment()
}
