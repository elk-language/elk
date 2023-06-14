package object

// Initialize the class hierarchy etc
func initBootstrap() {
	ClassClass = &Class{}
	PrimitiveObjectClass = &Class{metaClass: ClassClass}
	ObjectClass = &Class{
		metaClass: ClassClass,
		Parent:    PrimitiveObjectClass,
	}
	ClassClass.Parent = ObjectClass
	ClassClass.metaClass = ClassClass

	ModuleClass = NewClass()
	RootModule = NewModule(ModuleWithName("Root"))
	StdModule = NewModule()

	RootModule.AddConstant("Root", RootModule)
	RootModule.AddConstant("Std", StdModule)

	StdModule.AddConstant("Class", ClassClass)
	StdModule.AddConstant("Object", ObjectClass)
	StdModule.AddConstant("PrimitiveObject", PrimitiveObjectClass)
	StdModule.AddConstant("Module", ModuleClass)
}

// Initialize all built-ins
func init() {
	initBootstrap()
	initInt()
}
