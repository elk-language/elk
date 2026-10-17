package value

// Initialize the class hierarchy etc
func initBootstrap() {
	NamespaceClass = &Class{
		ConstructorFunc: UndefinedConstructor,
		ConstantContainer: ConstantContainer{
			constants: make(SymbolMap),
		},
		MethodContainer: MethodContainer{
			methods: make(MethodMap),
		},
	}
	ClassClass = &Class{
		ConstructorFunc: ClassConstructor,
		ConstantContainer: ConstantContainer{
			constants: make(SymbolMap),
		},
		MethodContainer: MethodContainer{
			methods: make(MethodMap),
			parent:  NamespaceClass,
		},
	}
	ValueClass = &Class{
		metaClass:       ClassClass,
		ConstructorFunc: ObjectConstructor,
		ConstantContainer: ConstantContainer{
			constants: make(SymbolMap),
		},
		MethodContainer: MethodContainer{
			methods: make(MethodMap),
		},
	}
	ObjectClass = &Class{
		metaClass: ClassClass,
		MethodContainer: MethodContainer{
			methods: make(MethodMap),
			parent:  ValueClass,
		},
		ConstructorFunc: ObjectConstructor,
		ConstantContainer: ConstantContainer{
			constants: make(SymbolMap),
		},
	}
	NamespaceClass.parent = ObjectClass
	NamespaceClass.metaClass = ClassClass
	ClassClass.metaClass = ClassClass

	ModuleClass = NewClassWithOptions(ClassWithConstructor(ModuleConstructor))
	RootModule = NewModuleWithOptions(ModuleWithName("Root"))
	StdModule = NewModule()

	RootModule.AddConstantString("Root", Ref(RootModule))
	RegisterNativeModule("Root", "value.RootModule")

	RootModule.AddConstantString("Std", Ref(StdModule))
	RegisterNativeModule("Std", "value.StdModule")

	StdModule.AddConstantString("Namespace", Ref(NamespaceClass))
	RegisterNativeClass("Std::Namespace", "value.NamespaceClass")

	StdModule.AddConstantString("Class", Ref(ClassClass))
	RegisterNativeClass("Std::Class", "value.ClassClass")

	StdModule.AddConstantString("Object", Ref(ObjectClass))
	RegisterNativeClass("Std::Object", "value.ObjectClass")

	StdModule.AddConstantString("Value", Ref(ValueClass))
	RegisterNativeClass("Std::Value", "value.ValueClass")

	StdModule.AddConstantString("Module", Ref(ModuleClass))
	RegisterNativeClass("Std::Module", "value.ModuleClass")
}

// Initialize all built-ins
func InitGlobalEnvironment() {
	initBootstrap()
	initGlobalObject()
	initFS()
	initPath()
	initError()
	initInterface()
	initAborter()
	initKernel()
	initRuntime()
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
	initUInt()
	initFloat()
	initBigFloat()
	initFloat64()
	initFloat32()
	initString()
	initChar()
	initSymbol()
	initRegex()
	initIterator()
	initIterable()
	initImmutableCollection()
	initCollection()
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
	initPosition()
	initSpan()
	initMacro()
	initElk()
	initElkToken()
	initElkAST()
	initThread()
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
	initChannel()
	initReadChannel()
	initWriteChannel()
	initImmutableBox()
	initWeak()
	initBox()
	initResult()
	initColorizer()
	initDiagnostic()
	initDiagnosticList()
	initSyncDiagnosticList()
	initElkLexer()
	initElkParser()
	initElkType()
	initLocation()
	initDuration()
	initDate()
	initDateSpan()
	initTime()
	initTimeSpan()
	initDateTime()
	initDateTimeSpan()
	initTimezone()
}
