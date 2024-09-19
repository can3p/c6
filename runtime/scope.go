package runtime

import (
	"fmt"

	"github.com/c9s/c6/ast"
)

type Scope struct {
	Parent    *Scope
	Variables map[string]ast.Value
	Mixins    map[string]*ast.MixinStmt
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		Parent:    parent,
		Variables: make(map[string]ast.Value, 4),
		Mixins:    make(map[string]*ast.MixinStmt, 4),
	}
}

func (s *Scope) Lookup(name string) (ast.Value, error) {
	if v, ok := s.Variables[name]; ok {
		return v, nil
	} else if s.Parent != nil {
		return s.Parent.Lookup(name)
	}

	return nil, fmt.Errorf("Undefined variable.")
}

func (s *Scope) Insert(name string, obj ast.Value) {
	s.Variables[name] = obj
}

func (s *Scope) LookupMixin(name string) (*ast.MixinStmt, error) {
	if v, ok := s.Mixins[name]; ok {
		return v, nil
	} else if s.Parent != nil {
		return s.Parent.LookupMixin(name)
	}

	return nil, fmt.Errorf("Undefined mixin - [%s]", name)
}

func (s *Scope) InsertMixin(name string, obj *ast.MixinStmt) {
	s.Mixins[name] = obj
}

func (s *Scope) GetGlobal() *Scope {
	scope := s

	for scope.Parent != nil {
		scope = scope.Parent
	}

	return scope
}
