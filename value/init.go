package value

// Initialize the class hierarchy etc
func initBootstrap() {
	ClassClass = &Class{
		ConstructorFunc: ClassConstructor,
	}
	PrimitiveObjectClass = &Class{
		metaClass:       ClassClass,
		ConstructorFunc: ObjectConstructor,
	}
	ObjectClass = &Class{
		metaClass:       ClassClass,
		Parent:          PrimitiveObjectClass,
		ConstructorFunc: ObjectConstructor,
	}
	ClassClass.Parent = ObjectClass
	ClassClass.metaClass = ClassClass

	ModuleClass = NewClassWithOptions(ClassWithConstructor(ModuleConstructor))
	RootModule = NewModuleWithOptions(ModuleWithName("Root"))
	StdModule = NewModule()

	RootModule.AddConstantString("Root", RootModule)
	RootModule.AddConstantString("Std", StdModule)

	StdModule.AddConstantString("Class", ClassClass)
	StdModule.AddConstantString("Object", ObjectClass)
	StdModule.AddConstantString("PrimitiveObject", PrimitiveObjectClass)
	StdModule.AddConstantString("Module", ModuleClass)
}

// Initialize all built-ins
func init() {
	initBootstrap()
	initNumeric()
	initBool()
	initTrue()
	initFalse()
	initNil()
	initInt()
	initSmallInt()
	initBigInt()
	initInt64()
	initInt32()
	initInt16()
	initInt8()
	initUInt64()
	initUInt32()
	initUInt16()
	initUInt8()
	initFloat()
	initBigFloat()
	initFloat64()
	initFloat32()
	initString()
	initChar()
	initSymbol()
	initSymbolMap()
	initException()
	initRange()
	initArithmeticSequence()
}
