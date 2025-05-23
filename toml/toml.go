package toml

import (
	"encoding/json"
	"fmt"
	"strings"

	lib "github.com/pelletier/go-toml"
)

// Toml is a struct
type Toml struct {
	path string
	out  string

	raw []byte

	tree *lib.Tree
}

// NewToml returns the Toml
func NewToml(path string) (Toml, error) {
	toml := Toml{
		path: path,
	}

	err := toml.readFile()
	if err != nil {
		return toml, err
	}

	err = toml.load()
	if err != nil {
		return toml, err
	}

	return toml, nil
}

func (t *Toml) load() error {
	var err error

	t.tree, err = lib.LoadBytes(t.raw)

	if err != nil {
		return err
	}

	return nil
}

// Dest set output given path
func (t *Toml) Out(path string) {
	t.out = path
}

// Get the value at key in the Tree.
// [Wrapped function go-toml.]
func (t *Toml) Get(query string) interface{} {
	v := t.tree.GetPath([]string{query})
	return v
}

// Set the value at key in the Tree.
// [Wrapped function go-toml.]
func (t *Toml) Set(query, attr string, data interface{}) error {
	path := []string{query}
	if attr != "" {
		path = append(path, attr)
	}
	t.tree.SetPath(path, data)
	return nil
}

func (t *Toml) Keys() []string {
	return t.tree.Keys()
}
func (t *Toml) List(query string) []string {
	dst := make([]string, 0)
	for _, k := range t.tree.Keys() {
		var business, comment any = "", ""
		if v := t.tree.GetPath([]string{k, "business"}); v != nil {
			business = v
		}
		if v := t.tree.GetPath([]string{k, "comment"}); v != nil {
			comment = v
		}

		if query == "" {
			dst = append(dst, fmt.Sprintf("%s\t%v\t%v", k, business, comment))
			continue
		}
		if strings.Contains(strings.ToLower(k), strings.ToLower(query)) {
			dst = append(dst, fmt.Sprintf("%s\t%v\t%v", k, business, comment))
		}
	}
	return dst
}
func (t *Toml) Delete(query, attr string) error {
	path := []string{query, attr}
	if !t.tree.HasPath(path) {
		return nil
	}
	return t.tree.DeletePath(path)
}
func (t *Toml) Clear(query string) error {
	path := []string{query}
	if !t.tree.HasPath(path) {
		return nil
	}
	return t.tree.DeletePath(path)
}
func (t *Toml) ToJson() (string, error) {
	res := t.tree.ToMap()
	m, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		return "", err
	}
	return string(m), nil
}

func (t *Toml) ToToml() (string, error) {
	return t.tree.ToTomlString()
}
