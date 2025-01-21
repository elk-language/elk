package value

func initLockable() {
	StdModule.AddConstantString("Lockable", Ref(NewInterface()))
}
