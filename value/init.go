package value

// Initialize the class hierarchy etc
func initBootstrap() {
	ClassClass = &Class{
		ConstructorFunc: ClassConstructor,
		ModulelikeObject: ModulelikeObject{
			Constants: make(SymbolMap),
		},
		Methods: make(MethodMap),
	}
	PrimitiveObjectClass = &Class{
		metaClass:       ClassClass,
		ConstructorFunc: ObjectConstructor,
		ModulelikeObject: ModulelikeObject{
			Constants: make(SymbolMap),
		},
		Methods: make(MethodMap),
	}
	ObjectClass = &Class{
		metaClass:       ClassClass,
		Parent:          PrimitiveObjectClass,
		ConstructorFunc: ObjectConstructor,
		ModulelikeObject: ModulelikeObject{
			Constants: make(SymbolMap),
		},
		Methods: make(MethodMap),
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
	initGlobalObject()
	initUndefined()
	initMixin()
	initNumeric()
	initBool()
	initTrue()
	initFalse()
	initNil()
	initInt()
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
	initList()
	initTuple()
	initSymbolMap()
	initException()
	initRange()
	initArithmeticSequence()
}
