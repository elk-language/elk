package value

var SyncModule *Module // ::Std::Sync

func initSync() {
	SyncModule = NewModule()
	StdModule.AddConstantString("Sync", Ref(SyncModule))
}
