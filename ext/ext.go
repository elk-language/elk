// Package ext contains functions that are used to handle Elk native extensions
package ext

import "github.com/elk-language/elk/types"

var Map map[string]*Extension

func init() {
	Map = make(map[string]*Extension)
}

type RuntimeInitialiser func()
type TypecheckerInitialiser func(types.Checker)

type Extension struct {
	Name                string
	GoPackagePath       string
	RuntimeInitFuncName string
	RuntimeInit         RuntimeInitialiser
	TypecheckerInit     TypecheckerInitialiser
	Initialised         bool
}

func (e *Extension) Init(checker types.Checker) {
	if e.Initialised {
		return
	}

	if e.RuntimeInit != nil && !checker.HasNativeGoCompiler() {
		e.RuntimeInit()
	}
	if e.TypecheckerInit != nil {
		e.TypecheckerInit(checker)
	}
	e.Initialised = true
}

func New(name string, goPackagePath string, runtimeInitName string, runtimeInit RuntimeInitialiser, typeInit TypecheckerInitialiser) *Extension {
	return &Extension{
		Name:                name,
		GoPackagePath:       goPackagePath,
		RuntimeInitFuncName: runtimeInitName,
		RuntimeInit:         runtimeInit,
		TypecheckerInit:     typeInit,
	}
}

// Register registers a new native extension
func Register(name string, goPackagePath string, runtimeInitName string, runtimeInit RuntimeInitialiser, typeInit TypecheckerInitialiser) {
	Map[name] = New(
		name,
		goPackagePath,
		runtimeInitName,
		runtimeInit,
		typeInit,
	)
}
