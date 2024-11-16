package vm_test

import (
	"testing"

	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/value"
)

func TestVMSource_Exponentiate(t *testing.T) {
	tests := sourceTestTable{
		"Int64 ** Int64": {
			source:       "2i64 ** 10i64",
			wantStackTop: value.Int64(1024),
		},
		"Int64 ** Int32": {
			source: "2i64 ** 10i32",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int64` for parameter `other` in call to `**`, got type `10i32`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_Modulo(t *testing.T) {
	tests := sourceTestTable{
		"Int64 % Int64": {
			source:       "29i64 % 3i64",
			wantStackTop: value.Int64(2),
		},
		"SmallInt % Float": {
			source: "250 % 4.5",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(6, 1, 7), P(8, 1, 9)), "expected type `Std::Int` for parameter `other` in call to `%`, got type `4.5`"),
			},
		},
		"Int64 % Int32": {
			source: "11i64 % 2i32",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(8, 1, 9), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `%`, got type `2i32`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_RightBitshift(t *testing.T) {
	tests := sourceTestTable{
		"Int >> String": {
			source: "3 >> 'foo'",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(5, 1, 6), P(9, 1, 10)), "expected type `Std::AnyInt` for parameter `other` in call to `>>`, got type `\"foo\"`"),
			},
		},
		"UInt16 >> Float": {
			source: "3u16 >> 5.2",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(8, 1, 9), P(10, 1, 11)), "expected type `Std::AnyInt` for parameter `other` in call to `>>`, got type `5.2`"),
			},
		},
		"String >> Int": {
			source: "'36' >> 5",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(8, 1, 9)), "method `>>` is not defined on type `Std::String`"),
			},
		},

		"Int >> Int": {
			source:       "16 >> 2",
			wantStackTop: value.SmallInt(4),
		},
		"-Int >> Int": {
			source:       "-16 >> 2",
			wantStackTop: value.SmallInt(-4),
		},
		"Int >> -Int": {
			source:       "16 >> -2",
			wantStackTop: value.SmallInt(64),
		},
		"Int >> Int32": {
			source:       "39 >> 1i32",
			wantStackTop: value.SmallInt(19),
		},

		"Int64 >> Int64": {
			source:       "16i64 >> 2i64",
			wantStackTop: value.Int64(4),
		},
		"-Int64 >> Int64": {
			source:       "-16i64 >> 2i64",
			wantStackTop: value.Int64(-4),
		},
		"Int64 >> -Int64": {
			source:       "16i64 >> -2i64",
			wantStackTop: value.Int64(64),
		},
		"Int64 >> Int32": {
			source:       "39i64 >> 1i32",
			wantStackTop: value.Int64(19),
		},
		"Int64 >> UInt8": {
			source:       "120i64 >> 5u8",
			wantStackTop: value.Int64(3),
		},
		"Int64 >> Int": {
			source:       "54i64 >> 3",
			wantStackTop: value.Int64(6),
		},

		"Int32 >> Int32": {
			source:       "16i32 >> 2i32",
			wantStackTop: value.Int32(4),
		},
		"-Int32 >> Int32": {
			source:       "-16i32 >> 2i32",
			wantStackTop: value.Int32(-4),
		},
		"Int32 >> -Int32": {
			source:       "16i32 >> -2i32",
			wantStackTop: value.Int32(64),
		},
		"Int32 >> Int16": {
			source:       "39i32 >> 1i16",
			wantStackTop: value.Int32(19),
		},
		"Int32 >> UInt8": {
			source:       "120i32 >> 5u8",
			wantStackTop: value.Int32(3),
		},
		"Int32 >> Int": {
			source:       "54i32 >> 3",
			wantStackTop: value.Int32(6),
		},

		"Int16 >> Int16": {
			source:       "16i16 >> 2i16",
			wantStackTop: value.Int16(4),
		},
		"-Int16 >> Int16": {
			source:       "-16i16 >> 2i16",
			wantStackTop: value.Int16(-4),
		},
		"Int16 >> -Int16": {
			source:       "16i16 >> -2i16",
			wantStackTop: value.Int16(64),
		},
		"Int16 >> Int32": {
			source:       "39i16 >> 1i32",
			wantStackTop: value.Int16(19),
		},
		"Int16 >> UInt8": {
			source:       "120i16 >> 5u8",
			wantStackTop: value.Int16(3),
		},
		"Int16 >> Int": {
			source:       "54i16 >> 3",
			wantStackTop: value.Int16(6),
		},

		"Int8 >> Int8": {
			source:       "16i8 >> 2i8",
			wantStackTop: value.Int8(4),
		},
		"-Int8 >> Int8": {
			source:       "-16i8 >> 2i8",
			wantStackTop: value.Int8(-4),
		},
		"Int8 >> -Int8": {
			source:       "16i8 >> -2i8",
			wantStackTop: value.Int8(64),
		},
		"Int8 >> Int16": {
			source:       "39i8 >> 1i16",
			wantStackTop: value.Int8(19),
		},
		"Int8 >> UInt8": {
			source:       "120i8 >> 5u8",
			wantStackTop: value.Int8(3),
		},
		"Int8 >> Int": {
			source:       "54i8 >> 3",
			wantStackTop: value.Int8(6),
		},

		"UInt64 >> UInt64": {
			source:       "16u64 >> 2u64",
			wantStackTop: value.UInt64(4),
		},
		"UInt64 >> -Int": {
			source:       "16u64 >> -2",
			wantStackTop: value.UInt64(64),
		},
		"UInt64 >> Int32": {
			source:       "39u64 >> 1i32",
			wantStackTop: value.UInt64(19),
		},

		"UInt32 >> UInt32": {
			source:       "16u32 >> 2u32",
			wantStackTop: value.UInt32(4),
		},
		"UInt32 >> -Int": {
			source:       "16u32 >> -2",
			wantStackTop: value.UInt32(64),
		},
		"UInt32 >> Int32": {
			source:       "39u32 >> 1i32",
			wantStackTop: value.UInt32(19),
		},

		"UInt16 >> UInt16": {
			source:       "16u16 >> 2u16",
			wantStackTop: value.UInt16(4),
		},
		"UInt16 >> -Int": {
			source:       "16u16 >> -2",
			wantStackTop: value.UInt16(64),
		},
		"UInt16 >> Int32": {
			source:       "39u16 >> 1i32",
			wantStackTop: value.UInt16(19),
		},

		"UInt8 >> UInt8": {
			source:       "16u8 >> 2u8",
			wantStackTop: value.UInt8(4),
		},
		"UInt8 >> -Int": {
			source:       "16u8 >> -2",
			wantStackTop: value.UInt8(64),
		},
		"UInt8 >> Int32": {
			source:       "39u8 >> 1i32",
			wantStackTop: value.UInt8(19),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_LogicalRightBitshift(t *testing.T) {
	tests := sourceTestTable{
		"Int >>> String": {
			source: "3 >>> 'foo'",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(10, 1, 11)), "method `>>>` is not defined on type `Std::Int`"),
			},
		},
		"Int64 >>> String": {
			source: "3i64 >>> 'foo'",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(9, 1, 10), P(13, 1, 14)), "expected type `Std::AnyInt` for parameter `other` in call to `>>>`, got type `\"foo\"`"),
			},
		},
		"UInt16 >>> Float": {
			source: "3u16 >>> 5.2",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(9, 1, 10), P(11, 1, 12)), "expected type `Std::AnyInt` for parameter `other` in call to `>>>`, got type `5.2`"),
			},
		},
		"String >>> Int": {
			source: "'36' >>> 5",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(9, 1, 10)), "method `>>>` is not defined on type `Std::String`"),
			},
		},

		"Int64 >>> Int64": {
			source:       "16i64 >>> 2i64",
			wantStackTop: value.Int64(4),
		},
		"-Int64 >>> Int64": {
			source:       "-16i64 >>> 2i64",
			wantStackTop: value.Int64(4611686018427387900),
		},
		"Int64 >>> -Int64": {
			source:       "16i64 >>> -2i64",
			wantStackTop: value.Int64(64),
		},
		"Int64 >>> Int32": {
			source:       "39i64 >>> 1i32",
			wantStackTop: value.Int64(19),
		},
		"Int64 >>> UInt8": {
			source:       "120i64 >>> 5u8",
			wantStackTop: value.Int64(3),
		},
		"Int64 >>> Int": {
			source:       "54i64 >>> 3",
			wantStackTop: value.Int64(6),
		},

		"Int32 >>> Int32": {
			source:       "16i32 >>> 2i32",
			wantStackTop: value.Int32(4),
		},
		"-Int32 >>> Int32": {
			source:       "-16i32 >>> 2i32",
			wantStackTop: value.Int32(1073741820),
		},
		"Int32 >>> -Int32": {
			source:       "16i32 >>> -2i32",
			wantStackTop: value.Int32(64),
		},
		"Int32 >>> Int16": {
			source:       "39i32 >>> 1i16",
			wantStackTop: value.Int32(19),
		},
		"Int32 >>> UInt8": {
			source:       "120i32 >>> 5u8",
			wantStackTop: value.Int32(3),
		},
		"Int32 >>> Int": {
			source:       "54i32 >>> 3",
			wantStackTop: value.Int32(6),
		},

		"Int16 >>> Int16": {
			source:       "16i16 >>> 2i16",
			wantStackTop: value.Int16(4),
		},
		"-Int16 >>> Int16": {
			source:       "-16i16 >>> 2i16",
			wantStackTop: value.Int16(16380),
		},
		"Int16 >>> -Int16": {
			source:       "16i16 >>> -2i16",
			wantStackTop: value.Int16(64),
		},
		"Int16 >>> Int32": {
			source:       "39i16 >>> 1i32",
			wantStackTop: value.Int16(19),
		},
		"Int16 >>> UInt8": {
			source:       "120i16 >>> 5u8",
			wantStackTop: value.Int16(3),
		},
		"Int16 >>> Int": {
			source:       "54i16 >>> 3",
			wantStackTop: value.Int16(6),
		},

		"Int8 >>> Int8": {
			source:       "16i8 >>> 2i8",
			wantStackTop: value.Int8(4),
		},
		"-Int8 >>> Int8": {
			source:       "-16i8 >>> 2i8",
			wantStackTop: value.Int8(60),
		},
		"Int8 >>> -Int8": {
			source:       "16i8 >>> -2i8",
			wantStackTop: value.Int8(64),
		},
		"Int8 >>> Int16": {
			source:       "39i8 >>> 1i16",
			wantStackTop: value.Int8(19),
		},
		"Int8 >>> UInt8": {
			source:       "120i8 >>> 5u8",
			wantStackTop: value.Int8(3),
		},
		"Int8 >>> Int": {
			source:       "54i8 >>> 3",
			wantStackTop: value.Int8(6),
		},

		"UInt64 >>> UInt64": {
			source:       "16u64 >>> 2u64",
			wantStackTop: value.UInt64(4),
		},
		"UInt64 >>> -Int": {
			source:       "16u64 >>> -2",
			wantStackTop: value.UInt64(64),
		},
		"UInt64 >>> Int32": {
			source:       "39u64 >>> 1i32",
			wantStackTop: value.UInt64(19),
		},

		"UInt32 >>> UInt32": {
			source:       "16u32 >>> 2u32",
			wantStackTop: value.UInt32(4),
		},
		"UInt32 >>> -Int": {
			source:       "16u32 >>> -2",
			wantStackTop: value.UInt32(64),
		},
		"UInt32 >>> Int32": {
			source:       "39u32 >>> 1i32",
			wantStackTop: value.UInt32(19),
		},

		"UInt16 >>> UInt16": {
			source:       "16u16 >>> 2u16",
			wantStackTop: value.UInt16(4),
		},
		"UInt16 >>> -Int": {
			source:       "16u16 >>> -2",
			wantStackTop: value.UInt16(64),
		},
		"UInt16 >>> Int32": {
			source:       "39u16 >>> 1i32",
			wantStackTop: value.UInt16(19),
		},

		"UInt8 >>> UInt8": {
			source:       "16u8 >>> 2u8",
			wantStackTop: value.UInt8(4),
		},
		"UInt8 >>> -Int": {
			source:       "16u8 >>> -2",
			wantStackTop: value.UInt8(64),
		},
		"UInt8 >>> Int32": {
			source:       "39u8 >>> 1i32",
			wantStackTop: value.UInt8(19),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_LeftBitshift(t *testing.T) {
	tests := sourceTestTable{
		"Int << String": {
			source: "3 << 'foo'",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(5, 1, 6), P(9, 1, 10)), "expected type `Std::AnyInt` for parameter `other` in call to `<<`, got type `\"foo\"`"),
			},
		},
		"UInt16 << Float": {
			source: "3u16 << 5.2",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(8, 1, 9), P(10, 1, 11)), "expected type `Std::AnyInt` for parameter `other` in call to `<<`, got type `5.2`"),
			},
		},
		"String << Int": {
			source: "'36' << 5",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(8, 1, 9)), "method `<<` is not defined on type `Std::String`"),
			},
		},

		"Int << Int": {
			source:       "16 << 2",
			wantStackTop: value.SmallInt(64),
		},
		"-Int << Int": {
			source:       "-16 << 2",
			wantStackTop: value.SmallInt(-64),
		},
		"Int << -Int": {
			source:       "16 << -2",
			wantStackTop: value.SmallInt(4),
		},
		"Int << Int32": {
			source:       "39 << 1i32",
			wantStackTop: value.SmallInt(78),
		},

		"Int64 << Int64": {
			source:       "16i64 << 2i64",
			wantStackTop: value.Int64(64),
		},
		"-Int64 << Int64": {
			source:       "-16i64 << 2i64",
			wantStackTop: value.Int64(-64),
		},
		"Int64 << -Int64": {
			source:       "16i64 << -2i64",
			wantStackTop: value.Int64(4),
		},
		"Int64 << Int32": {
			source:       "39i64 << 1i32",
			wantStackTop: value.Int64(78),
		},
		"Int64 << UInt8": {
			source:       "120i64 << 5u8",
			wantStackTop: value.Int64(3840),
		},
		"Int64 << Int": {
			source:       "54i64 << 3",
			wantStackTop: value.Int64(432),
		},

		"Int32 << Int32": {
			source:       "16i32 << 2i32",
			wantStackTop: value.Int32(64),
		},
		"-Int32 << Int32": {
			source:       "-16i32 << 2i32",
			wantStackTop: value.Int32(-64),
		},
		"Int32 << -Int32": {
			source:       "16i32 << -2i32",
			wantStackTop: value.Int32(4),
		},
		"Int32 << Int16": {
			source:       "39i32 << 1i16",
			wantStackTop: value.Int32(78),
		},
		"Int32 << UInt8": {
			source:       "120i32 << 5u8",
			wantStackTop: value.Int32(3840),
		},
		"Int32 << Int": {
			source:       "54i32 << 3",
			wantStackTop: value.Int32(432),
		},

		"Int16 << Int16": {
			source:       "16i16 << 2i16",
			wantStackTop: value.Int16(64),
		},
		"-Int16 << Int16": {
			source:       "-16i16 << 2i16",
			wantStackTop: value.Int16(-64),
		},
		"Int16 << -Int16": {
			source:       "16i16 << -2i16",
			wantStackTop: value.Int16(4),
		},
		"Int16 << Int32": {
			source:       "39i16 << 1i32",
			wantStackTop: value.Int16(78),
		},
		"Int16 << UInt8": {
			source:       "120i16 << 5u8",
			wantStackTop: value.Int16(3840),
		},
		"Int16 << Int": {
			source:       "54i16 << 3",
			wantStackTop: value.Int16(432),
		},

		"Int8 << Int8": {
			source:       "16i8 << 2i8",
			wantStackTop: value.Int8(64),
		},
		"-Int8 << Int8": {
			source:       "-16i8 << 2i8",
			wantStackTop: value.Int8(-64),
		},
		"Int8 << -Int8": {
			source:       "16i8 << -2i8",
			wantStackTop: value.Int8(4),
		},
		"Int8 << Int16": {
			source:       "39i8 << 1i16",
			wantStackTop: value.Int8(78),
		},
		"Int8 << UInt8": {
			source:       "120i8 << 5u8",
			wantStackTop: value.Int8(0),
		},
		"Int8 << Int": {
			source:       "54i8 << 3",
			wantStackTop: value.Int8(-80),
		},

		"UInt64 << UInt64": {
			source:       "16u64 << 2u64",
			wantStackTop: value.UInt64(64),
		},
		"UInt64 << -Int": {
			source:       "16u64 << -2",
			wantStackTop: value.UInt64(4),
		},
		"UInt64 << Int32": {
			source:       "39u64 << 1i32",
			wantStackTop: value.UInt64(78),
		},

		"UInt32 << UInt32": {
			source:       "16u32 << 2u32",
			wantStackTop: value.UInt32(64),
		},
		"UInt32 << -Int": {
			source:       "16u32 << -2",
			wantStackTop: value.UInt32(4),
		},
		"UInt32 << Int32": {
			source:       "39u32 << 1i32",
			wantStackTop: value.UInt32(78),
		},

		"UInt16 << UInt16": {
			source:       "16u16 << 2u16",
			wantStackTop: value.UInt16(64),
		},
		"UInt16 << -Int": {
			source:       "16u16 << -2",
			wantStackTop: value.UInt16(4),
		},
		"UInt16 << Int32": {
			source:       "39u16 << 1i32",
			wantStackTop: value.UInt16(78),
		},

		"UInt8 << UInt8": {
			source:       "16u8 << 2u8",
			wantStackTop: value.UInt8(64),
		},
		"UInt8 << -Int": {
			source:       "16u8 << -2",
			wantStackTop: value.UInt8(4),
		},
		"UInt8 << Int32": {
			source:       "39u8 << 1i32",
			wantStackTop: value.UInt8(78),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_LogicalLeftBitshift(t *testing.T) {
	tests := sourceTestTable{
		"Int64 <<< String": {
			source: "3i64 <<< 'foo'",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(9, 1, 10), P(13, 1, 14)), "expected type `Std::AnyInt` for parameter `other` in call to `<<<`, got type `\"foo\"`"),
			},
		},
		"UInt16 <<< Float": {
			source: "3u16 <<< 5.2",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(9, 1, 10), P(11, 1, 12)), "expected type `Std::AnyInt` for parameter `other` in call to `<<<`, got type `5.2`"),
			},
		},
		"String <<< Int": {
			source: "'36' <<< 5",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(9, 1, 10)), "method `<<<` is not defined on type `Std::String`"),
			},
		},
		"Int <<< Int": {
			source: "16 <<< 2",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(7, 1, 8)), "method `<<<` is not defined on type `Std::Int`"),
			},
		},

		"Int64 <<< Int64": {
			source:       "16i64 <<< 2i64",
			wantStackTop: value.Int64(64),
		},
		"-Int64 <<< Int64": {
			source:       "-16i64 <<< 2i64",
			wantStackTop: value.Int64(-64),
		},
		"Int64 <<< -Int64": {
			source:       "16i64 <<< -2i64",
			wantStackTop: value.Int64(4),
		},
		"Int64 <<< Int32": {
			source:       "39i64 <<< 1i32",
			wantStackTop: value.Int64(78),
		},
		"Int64 <<< UInt8": {
			source:       "120i64 <<< 5u8",
			wantStackTop: value.Int64(3840),
		},
		"Int64 <<< Int": {
			source:       "54i64 <<< 3",
			wantStackTop: value.Int64(432),
		},

		"Int32 <<< Int32": {
			source:       "16i32 <<< 2i32",
			wantStackTop: value.Int32(64),
		},
		"-Int32 <<< Int32": {
			source:       "-16i32 <<< 2i32",
			wantStackTop: value.Int32(-64),
		},
		"Int32 <<< -Int32": {
			source:       "16i32 <<< -2i32",
			wantStackTop: value.Int32(4),
		},
		"Int32 <<< Int16": {
			source:       "39i32 <<< 1i16",
			wantStackTop: value.Int32(78),
		},
		"Int32 <<< UInt8": {
			source:       "120i32 <<< 5u8",
			wantStackTop: value.Int32(3840),
		},
		"Int32 <<< Int": {
			source:       "54i32 <<< 3",
			wantStackTop: value.Int32(432),
		},

		"Int16 <<< Int16": {
			source:       "16i16 <<< 2i16",
			wantStackTop: value.Int16(64),
		},
		"-Int16 <<< Int16": {
			source:       "-16i16 <<< 2i16",
			wantStackTop: value.Int16(-64),
		},
		"Int16 <<< -Int16": {
			source:       "16i16 <<< -2i16",
			wantStackTop: value.Int16(4),
		},
		"Int16 <<< Int32": {
			source:       "39i16 <<< 1i32",
			wantStackTop: value.Int16(78),
		},
		"Int16 <<< UInt8": {
			source:       "120i16 <<< 5u8",
			wantStackTop: value.Int16(3840),
		},
		"Int16 <<< Int": {
			source:       "54i16 <<< 3",
			wantStackTop: value.Int16(432),
		},

		"Int8 <<< Int8": {
			source:       "16i8 <<< 2i8",
			wantStackTop: value.Int8(64),
		},
		"-Int8 <<< Int8": {
			source:       "-16i8 <<< 2i8",
			wantStackTop: value.Int8(-64),
		},
		"Int8 <<< -Int8": {
			source:       "16i8 <<< -2i8",
			wantStackTop: value.Int8(4),
		},
		"Int8 <<< Int16": {
			source:       "39i8 <<< 1i16",
			wantStackTop: value.Int8(78),
		},
		"Int8 <<< UInt8": {
			source:       "120i8 <<< 5u8",
			wantStackTop: value.Int8(0),
		},
		"Int8 <<< Int": {
			source:       "54i8 <<< 3",
			wantStackTop: value.Int8(-80),
		},

		"UInt64 <<< UInt64": {
			source:       "16u64 <<< 2u64",
			wantStackTop: value.UInt64(64),
		},
		"UInt64 <<< -Int": {
			source:       "16u64 <<< -2",
			wantStackTop: value.UInt64(4),
		},
		"UInt64 <<< Int32": {
			source:       "39u64 <<< 1i32",
			wantStackTop: value.UInt64(78),
		},

		"UInt32 <<< UInt32": {
			source:       "16u32 <<< 2u32",
			wantStackTop: value.UInt32(64),
		},
		"UInt32 <<< -Int": {
			source:       "16u32 <<< -2",
			wantStackTop: value.UInt32(4),
		},
		"UInt32 <<< Int32": {
			source:       "39u32 <<< 1i32",
			wantStackTop: value.UInt32(78),
		},

		"UInt16 <<< UInt16": {
			source:       "16u16 <<< 2u16",
			wantStackTop: value.UInt16(64),
		},
		"UInt16 <<< -Int": {
			source:       "16u16 <<< -2",
			wantStackTop: value.UInt16(4),
		},
		"UInt16 <<< Int32": {
			source:       "39u16 <<< 1i32",
			wantStackTop: value.UInt16(78),
		},

		"UInt8 <<< UInt8": {
			source:       "16u8 <<< 2u8",
			wantStackTop: value.UInt8(64),
		},
		"UInt8 <<< -Int": {
			source:       "16u8 <<< -2",
			wantStackTop: value.UInt8(4),
		},
		"UInt8 <<< Int32": {
			source:       "39u8 <<< 1i32",
			wantStackTop: value.UInt8(78),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_BitwiseAnd(t *testing.T) {
	tests := sourceTestTable{
		"Float & Int": {
			source: "3.6 & 5",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(6, 1, 7)), "method `&` is not defined on type `Std::Float`"),
			},
		},
		"Int64 & String": {
			source: "3i64 & 'foo'",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `&`, got type `\"foo\"`"),
			},
		},
		"Int64 & SmallInt": {
			source: "3i64 & 5",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(7, 1, 8), P(7, 1, 8)), "expected type `Std::Int64` for parameter `other` in call to `&`, got type `5`"),
			},
		},
		"UInt16 & Float": {
			source: "3u16 & 5.2",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(7, 1, 8), P(9, 1, 10)), "expected type `Std::UInt16` for parameter `other` in call to `&`, got type `5.2`"),
			},
		},
		"String & Int": {
			source: "'36' & 5",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(7, 1, 8)), "method `&` is not defined on type `Std::String`"),
			},
		},
		"Int & Int": {
			source:       "25 & 14",
			wantStackTop: value.SmallInt(8),
		},
		"Int & BigInt": {
			source:       "255 & 9223372036857247042",
			wantStackTop: value.SmallInt(66),
		},
		"Int8 & Int8": {
			source:       "59i8 & 122i8",
			wantStackTop: value.Int8(58),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_BitwiseAndNot(t *testing.T) {
	tests := sourceTestTable{
		"Int64 &~ String": {
			source: "3i64 &~ 'foo'",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(8, 1, 9), P(12, 1, 13)), "expected type `Std::Int64` for parameter `other` in call to `&~`, got type `\"foo\"`"),
			},
		},
		"Int64 &~ SmallInt": {
			source: "3i64 &~ 5",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(8, 1, 9), P(8, 1, 9)), "expected type `Std::Int64` for parameter `other` in call to `&~`, got type `5`"),
			},
		},
		"UInt16 &~ Float": {
			source: "3u16 &~ 5.2",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(8, 1, 9), P(10, 1, 11)), "expected type `Std::UInt16` for parameter `other` in call to `&~`, got type `5.2`"),
			},
		},
		"String &~ Int": {
			source: "'36' &~ 5",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(8, 1, 9)), "method `&~` is not defined on type `Std::String`"),
			},
		},
		"Float &~ Int": {
			source: "3.6 &~ 5",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(7, 1, 8)), "method `&~` is not defined on type `Std::Float`"),
			},
		},
		"Int &~ Int": {
			source:       "25 &~ 14",
			wantStackTop: value.SmallInt(17),
		},
		"Int &~ BigInt": {
			source:       "255 &~ 9223372036857247042",
			wantStackTop: value.SmallInt(189),
		},
		"Int8 &~ Int8": {
			source:       "59i8 &~ 122i8",
			wantStackTop: value.Int8(1),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_BitwiseOr(t *testing.T) {
	tests := sourceTestTable{
		"Int64 | String": {
			source: "3i64 | 'foo'",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `|`, got type `\"foo\"`"),
			},
		},
		"Int64 | SmallInt": {
			source: "3i64 | 5",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(7, 1, 8), P(7, 1, 8)), "expected type `Std::Int64` for parameter `other` in call to `|`, got type `5`"),
			},
		},
		"UInt16 | Float": {
			source: "3u16 | 5.2",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(7, 1, 8), P(9, 1, 10)), "expected type `Std::UInt16` for parameter `other` in call to `|`, got type `5.2`"),
			},
		},
		"String | Int": {
			source: "'36' | 5",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(7, 1, 8)), "method `|` is not defined on type `Std::String`"),
			},
		},
		"Float | Int": {
			source: "3.6 | 5",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(6, 1, 7)), "method `|` is not defined on type `Std::Float`"),
			},
		},
		"Int | Int": {
			source:       "25 | 14",
			wantStackTop: value.SmallInt(31),
		},
		"Int | BigInt": {
			source:       "255 | 9223372036857247042",
			wantStackTop: value.ParseBigIntPanic("9223372036857247231", 10),
		},
		"Int8 | Int8": {
			source:       "59i8 | 122i8",
			wantStackTop: value.Int8(123),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_BitwiseXor(t *testing.T) {
	tests := sourceTestTable{
		"Int64 ^ String": {
			source: "3i64 ^ 'foo'",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(7, 1, 8), P(11, 1, 12)), "expected type `Std::Int64` for parameter `other` in call to `^`, got type `\"foo\"`"),
			},
		},
		"Int64 ^ SmallInt": {
			source: "3i64 ^ 5",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(7, 1, 8), P(7, 1, 8)), "expected type `Std::Int64` for parameter `other` in call to `^`, got type `5`"),
			},
		},
		"UInt16 ^ Float": {
			source: "3u16 ^ 5.2",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(7, 1, 8), P(9, 1, 10)), "expected type `Std::UInt16` for parameter `other` in call to `^`, got type `5.2`"),
			},
		},
		"String ^ Int": {
			source: "'36' ^ 5",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(7, 1, 8)), "method `^` is not defined on type `Std::String`"),
			},
		},
		"Float ^ Int": {
			source: "3.6 ^ 5",
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(6, 1, 7)), "method `^` is not defined on type `Std::Float`"),
			},
		},
		"Int ^ Int": {
			source:       "25 ^ 14",
			wantStackTop: value.SmallInt(23),
		},
		"Int ^ BigInt": {
			source:       "255 ^ 9223372036857247042",
			wantStackTop: value.ParseBigIntPanic("9223372036857247165", 10),
		},
		"Int8 ^ Int8": {
			source:       "59i8 ^ 122i8",
			wantStackTop: value.Int8(65),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
