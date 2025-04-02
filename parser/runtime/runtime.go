package runtime

func InitGlobalEnvironment() {
	initParser()
	initResult()
}

func init() {
	InitGlobalEnvironment()
}
