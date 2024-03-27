package value

// Initialize the class hierarchy etc
func initBootstrap() {
	ClassClass = &Class{
		ConstructorFunc: ClassConstructor,
		ConstantContainer: ConstantContainer{
			Constants: make(SymbolMap),
		},
		MethodContainer: MethodContainer{
			Methods: make(MethodMap),
		},
	}
	ValueClass = &Class{
		metaClass:       ClassClass,
		ConstructorFunc: ObjectConstructor,
		ConstantContainer: ConstantContainer{
			Constants: make(SymbolMap),
		},
		MethodContainer: MethodContainer{
			Methods: make(MethodMap),
		},
	}
	ValueClass.SetAbstract()
	ObjectClass = &Class{
		metaClass: ClassClass,
		MethodContainer: MethodContainer{
			Methods: make(MethodMap),
			Parent:  ValueClass,
		},
		ConstructorFunc: ObjectConstructor,
		ConstantContainer: ConstantContainer{
			Constants: make(SymbolMap),
		},
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
	StdModule.AddConstantString("Value", ValueClass)
	StdModule.AddConstantString("Module", ModuleClass)
}

// Initialize all built-ins
func init() {
	initBootstrap()
	initGlobalObject()
	initUndefined()
	initMethod()
	initMixin()
	initComparable()
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
	initRegex()
	initTuple()
	initList()
	initArrayTuple()
	initArrayList()
	initPair()
	initRecord()
	initMap()
	initHashRecord()
	initHashMap()
	initSymbolMap()
	initException()
	initRange()
	initClosedRange()
	initOpenRange()
	initRightOpenRange()
	initLeftOpenRange()
	initBeginlessClosedRange()
	initEndlessClosedRange()
	initBeginlessOpenRange()
	initEndlessOpenRange()
	initTimezone()
	initTime()
}
