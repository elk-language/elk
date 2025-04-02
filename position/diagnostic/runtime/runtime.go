package runtime

func InitGlobalEnvironment() {
	initDiagnostic()
	initDiagnosticList()
	initDiagnosticListIterator()
	initSyncDiagnosticList()
	initSyncDiagnosticListIterator()
}

func init() {
	InitGlobalEnvironment()
}
