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

	RootModule.AddConstantString("Root", Ref(RootModule))
	RootModule.AddConstantString("Std", Ref(StdModule))

	StdModule.AddConstantString("Class", Ref(ClassClass))
	StdModule.AddConstantString("Object", Ref(ObjectClass))
	StdModule.AddConstantString("Value", Ref(ValueClass))
	StdModule.AddConstantString("Module", Ref(ModuleClass))
}

// Initialize all built-ins
func InitGlobalEnvironment() {
	initBootstrap()
	initGlobalObject()
	initInterface()
	initKernel()
	initDebug()
	initUndefined()
	initFunction()
	initMethod()
	initClosure()
	initGenerator()
	initMixin()
	initComparable()
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
	initIterator()
	initTuple()
	initList()
	initArrayTuple()
	initArrayList()
	initPair()
	initRecord()
	initMap()
	initHashRecord()
	initHashMap()
	initSet()
	initHashSet()
	initSymbolMap()
	initRange()
	initClosedRange()
	initOpenRange()
	initRightOpenRange()
	initLeftOpenRange()
	initBeginlessClosedRange()
	initEndlessClosedRange()
	initBeginlessOpenRange()
	initEndlessOpenRange()
	initDuration()
	initTimezone()
	initTime()
	initPosition()
	initSpan()
	initElk()
	initElkToken()
	initElkAST()
	initThread()
	initChannel()
	initLockable()
	initSync()
	initWaitGroup()
	initMutex()
	initRWMutex()
	initROMutex()
	initOnce()
	initStackTrace()
	initCallFrame()
	initPromise()
	initThreadPool()
	initFS()
	initPath()
	initError()
}

func init() {
	InitGlobalEnvironment()
}
