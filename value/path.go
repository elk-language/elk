package value

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

var PathClass *Class // ::Std::FS::Path

func initPath() {
	PathClass = NewClass()
	FSModule.AddConstantString("Path", Ref(PathClass))
}

type Path struct {
	Value string
}

func NewPath(value string) *Path {
	return &Path{
		Value: value,
	}
}

func NewPathFromSlash(value string) *Path {
	return NewPath(filepath.FromSlash(value))
}

func BuildPath(elements ...string) *Path {
	return NewPath(path.Join(elements...))
}

func (*Path) Class() *Class {
	return PathClass
}

func (*Path) DirectClass() *Class {
	return PathClass
}

func (*Path) SingletonClass() *Class {
	return nil
}

func (p *Path) Copy() Reference {
	return p
}

func (p *Path) Inspect() string {
	return fmt.Sprintf("Std::FS::Path(%s)", String(p.Value).Inspect())
}

func (p *Path) Error() string {
	return p.Inspect()
}

func (*Path) InstanceVariables() SymbolMap {
	return nil
}

// Reports whether the path is absolute.
func (p *Path) IsAbsolute() bool {
	return filepath.IsAbs(p.Value)
}

// Reports whether the path is local.
// It is a primitive lexical check it does not take into account
// symbolic links etc
func (p *Path) IsLocal() bool {
	return filepath.IsLocal(p.Value)
}

func (p *Path) String() string {
	return p.Value
}

func (p *Path) Equal(other *Path) bool {
	return p.Value == other.Value
}

func (p *Path) SlashString() string {
	return filepath.ToSlash(p.Value)
}

func (p *Path) BackslashString() string {
	return strings.ReplaceAll(p.Value, "/", "\\")
}

// VolumeName returns leading volume name.
// Given "C:\foo\bar" it returns "C:" on Windows.
// Given "\\host\share\foo" it returns "\\host\share". On other platforms it returns "".
func (p *Path) VolumeName() string {
	return filepath.VolumeName(p.Value)
}

// Split the path into individual elements
// separated by the OS separator (`/` or `\`)
func (p *Path) Split() []string {
	return strings.Split(p.Value, string(filepath.Separator))
}

// Creates a new path based on `self` that is
// the shortest possible version of it based on lexical analysis.
func (p *Path) Normalize() *Path {
	return NewPath(filepath.Clean(p.Value))
}

// Returns a new path based on `self` omitting the last element.
// Typically this would result in the path to the parent directory.
func (p *Path) Dir() *Path {
	return NewPath(filepath.Dir(p.Value))
}

// Returns the absolute version of this path.
// If the path is not absolute it will be joined with the current working directory to turn it into an absolute path.
func (p *Path) ToAbsolute() (*Path, error) {
	if p.IsAbsolute() {
		return p, nil
	}

	result, err := filepath.Abs(p.Value)
	if err != nil {
		return nil, err
	}
	return NewPath(result), nil
}

// Returns the last element of the path.
// Typically this is the name of the file.
func (p *Path) Base() string {
	return filepath.Base(p.Value)
}

// Returns the extension of the file.
// The extension is the suffix beginning at the final dot in the final element of path; it is empty if there is no dot.
//
// "index" => "";
// "index.js" => ".js";
// "index.html.erb" => ".erb";
func (p *Path) Extension() string {
	return filepath.Ext(p.Value)
}

// Creates a new path based on `target` that is relative to `self`.
func (p *Path) ToRelative(target *Path) (*Path, error) {
	result, err := filepath.Rel(p.Value, target.Value)
	if err != nil {
		return nil, err
	}

	return NewPath(result), nil
}

// Checks whether the path matches the given glob pattern.
func (p *Path) MatchesGlob(pattern string) (bool, error) {
	return doublestar.PathMatch(pattern, p.Value)
}
