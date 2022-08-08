package rule

import (
	"context"
	"errors"
	"regexp"
)

var (
	ErrUnknown = errors.New("undefined proxy rule")
)

type ConfigRepository interface {
	GetAll(ctx context.Context) ([]Ruleset, error)
}

type Ruleset struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Frontend    Frontend        `yaml:"frontend"`
	Backend     Backend         `yaml:"backend"`
	Middlewares MiddlewareSpecs `yaml:"middlewares"`
	Hooks       HookSpecs       `yaml:"hooks"`
}

type MiddlewareSpec struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"`
}

type MiddlewareSpecs []MiddlewareSpec

func (m MiddlewareSpecs) Get(name string) (MiddlewareSpec, bool) {
	for _, n := range m {
		if n.Name == name {
			return n, true
		}
	}
	return MiddlewareSpec{}, false
}

type HookSpec struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"`
}

type HookSpecs []HookSpec

func (m HookSpecs) Get(name string) (HookSpec, bool) {
	for _, n := range m {
		if n.Name == name {
			return n, true
		}
	}
	return HookSpec{}, false
}

type Frontend struct {
	URL   string         `yaml:"url"`
	URLRx *regexp.Regexp `yaml:"-"`

	Method string `yaml:"method"`
}

type Backend struct {
	URL       string `yaml:"url"`
	Namespace string `yaml:"namespace"`
	Prefix    string `yaml:"prefix"`
}
