// Package strings implements strings interface for anko script.
package strings

import (
	pkg "strings"

	"github.com/mattn/anko/vm"
)

func Import(env *vm.Env) *vm.Env {
	m := env.NewPackage("strings")
	m.Define("Contains", pkg.Contains)
	m.Define("ContainsAny", pkg.ContainsAny)
	m.Define("ContainsRune", pkg.ContainsRune)
	m.Define("Count", pkg.Count)
	m.Define("EqualFold", pkg.EqualFold)
	m.Define("Fields", pkg.Fields)
	m.Define("FieldsFunc", pkg.FieldsFunc)
	m.Define("HasPrefix", pkg.HasPrefix)
	m.Define("HasSuffix", pkg.HasSuffix)
	m.Define("Index", pkg.Index)
	m.Define("IndexAny", pkg.IndexAny)
	m.Define("IndexByte", pkg.IndexByte)
	m.Define("IndexFunc", pkg.IndexFunc)
	m.Define("IndexRune", pkg.IndexRune)
	m.Define("Join", pkg.Join)
	m.Define("LastIndex", pkg.LastIndex)
	m.Define("LastIndexAny", pkg.LastIndexAny)
	m.Define("LastIndexFunc", pkg.LastIndexFunc)
	m.Define("Map", pkg.Map)
	m.Define("NewReader", pkg.NewReader)
	m.Define("NewReplacer", pkg.NewReplacer)
	m.Define("Repeat", pkg.Repeat)
	m.Define("Replace", pkg.Replace)
	m.Define("Split", pkg.Split)
	m.Define("SplitAfter", pkg.SplitAfter)
	m.Define("SplitAfterN", pkg.SplitAfterN)
	m.Define("SplitN", pkg.SplitN)
	m.Define("Title", pkg.Title)
	m.Define("ToLower", pkg.ToLower)
	m.Define("ToLowerSpecial", pkg.ToLowerSpecial)
	m.Define("ToTitle", pkg.ToTitle)
	m.Define("ToTitleSpecial", pkg.ToTitleSpecial)
	m.Define("ToUpper", pkg.ToUpper)
	m.Define("ToUpperSpecial", pkg.ToUpperSpecial)
	m.Define("Trim", pkg.Trim)
	m.Define("TrimFunc", pkg.TrimFunc)
	m.Define("TrimLeft", pkg.TrimLeft)
	m.Define("TrimLeftFunc", pkg.TrimLeftFunc)
	m.Define("TrimPrefix", pkg.TrimPrefix)
	m.Define("TrimRight", pkg.TrimRight)
	m.Define("TrimRightFunc", pkg.TrimRightFunc)
	m.Define("TrimSpace", pkg.TrimSpace)
	m.Define("TrimSuffix", pkg.TrimSuffix)
	return m
}
