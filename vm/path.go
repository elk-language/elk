package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::FS::Path
func initPath() {
	// Instance methods
	c := &value.PathClass.MethodContainer

	Def(
		c,
		"#init",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Path)(args[0].Pointer())
			str := (string)(args[1].MustReference().(value.String))

			self.Value = str
			return value.Ref(self), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Path)(args[0].Pointer())
			other := (*value.Path)(args[1].Pointer())
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"is_absolute",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Path)(args[0].Pointer())
			return value.BoolVal(self.IsAbsolute()), value.Undefined
		},
	)
	Def(
		c,
		"is_local",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Path)(args[0].Pointer())
			return value.BoolVal(self.IsLocal()), value.Undefined
		},
	)
	Def(
		c,
		"to_string",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Path)(args[0].Pointer())
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)
	Def(
		c,
		"to_slash_string",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Path)(args[0].Pointer())
			return value.Ref(value.String(self.SlashString())), value.Undefined
		},
	)
	Def(
		c,
		"to_backslash_string",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Path)(args[0].Pointer())
			return value.Ref(value.String(self.BackslashString())), value.Undefined
		},
	)
	Def(
		c,
		"volume_name",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Path)(args[0].Pointer())
			return value.Ref(value.String(self.VolumeName())), value.Undefined
		},
	)
	Def(
		c,
		"base",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Path)(args[0].Pointer())
			return value.Ref(value.String(self.Base())), value.Undefined
		},
	)
	Def(
		c,
		"extension",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Path)(args[0].Pointer())
			return value.Ref(value.String(self.Extension())), value.Undefined
		},
	)
	Def(
		c,
		"split",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Path)(args[0].Pointer())
			splitSlice := self.Split()

			result := value.NewArrayListOfValueWithLength(len(splitSlice))
			for i, element := range splitSlice {
				result.SetAt(i, value.Ref(value.String(element)))
			}
			return value.Ref(result), value.Undefined
		},
	)
	Def(
		c,
		"normalize",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Path)(args[0].Pointer())
			return value.Ref(self.Normalize()), value.Undefined
		},
	)
	Def(
		c,
		"dir",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Path)(args[0].Pointer())
			return value.Ref(self.Dir()), value.Undefined
		},
	)
	Def(
		c,
		"to_absolute",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Path)(args[0].Pointer())
			result, err := self.ToAbsolute()
			if err != nil {
				return value.Undefined, value.Ref(value.NewError(value.PathErrorClass, err.Error()))
			}

			return value.Ref(result), value.Undefined
		},
	)
	Def(
		c,
		"to_relative",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Path)(args[0].Pointer())
			target := (*value.Path)(args[1].Pointer())

			result, err := self.ToRelative(target)
			if err != nil {
				return value.Undefined, value.Ref(value.NewError(value.PathErrorClass, err.Error()))
			}

			return value.Ref(result), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"matches_glob",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Path)(args[0].Pointer())
			pattern := (string)(args[1].MustReference().(value.String))

			result, err := self.MatchesGlob(pattern)
			if err != nil {
				return value.Undefined, value.Ref(value.NewError(value.GlobErrorClass, err.Error()))
			}

			return value.BoolVal(result), value.Undefined
		},
	)
	Def(
		c,
		"inspect",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Path)(args[0].Pointer())
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)

	// Singleton methods
	c = &value.PathClass.SingletonClass().MethodContainer

	Def(
		c,
		"from_slash",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			str := (string)(args[1].MustReference().(value.String))

			result := value.NewPathFromSlash(str)
			return value.Ref(result), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"build",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			var elements []string
			for val, err := range Iterate(vm, args[1]) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				elements = append(elements, string(val.MustReference().(value.String)))
			}

			result := value.BuildPath(elements...)
			return value.Ref(result), value.Undefined
		},
		DefWithParameters(1),
	)
}
