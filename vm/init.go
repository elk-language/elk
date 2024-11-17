package vm

import "github.com/elk-language/elk/value"

func InitGlobalEnvironment() {
	value.InitGlobalEnvironment()

	initArrayList()
	initArrayListIterator()
	initArrayTuple()
	initArrayTupleIterator()
	initBeginlessClosedRange()
	initBeginlessOpenRange()
	initBigFloat()
	initChar()
	initClass()
	initClosedRange()
	initClosedRangeIterator()
	initComparable()
	initEndlessClosedRange()
	initEndlessClosedRangeIterator()
	initEndlessOpenRange()
	initEndlessOpenRangeIterator()
	initFloat()
	initHashMap()
	initHashMapIterator()
	initHashRecord()
	initHashRecordIterator()
	initHashSet()
	initHashSetIterator()
	initInt()
	initKernel()
	initLeftOpenRange()
	initLeftOpenRangeIterator()
	initMethod()
	initMixin()
	initModule()
	initOpenRange()
	initOpenRangeIterator()
	initPair()
	initRegex()
	initRightOpenRange()
	initRightOpenRangeIterator()
	initString()
	initStringCharIterator()
	initStringByteIterator()
	initSymbol()
	initTime()
	initValue()
}

func init() {
	InitGlobalEnvironment()
}
