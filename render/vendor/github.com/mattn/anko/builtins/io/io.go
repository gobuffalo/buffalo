// Package io implements io interface for anko script.
package io

import (
	pkg "io"

	"github.com/mattn/anko/vm"
)

func Import(env *vm.Env) *vm.Env {
	m := env.NewPackage("io")
	m.Define("Copy", pkg.Copy)
	m.Define("CopyN", pkg.CopyN)
	m.Define("EOF", pkg.EOF)
	m.Define("ErrClosedPipe", pkg.ErrClosedPipe)
	m.Define("ErrNoProgress", pkg.ErrNoProgress)
	m.Define("ErrShortBuffer", pkg.ErrShortBuffer)
	m.Define("ErrShortWrite", pkg.ErrShortWrite)
	m.Define("ErrUnexpectedEOF", pkg.ErrUnexpectedEOF)
	m.Define("LimitReader", pkg.LimitReader)
	m.Define("MultiReader", pkg.MultiReader)
	m.Define("MultiWriter", pkg.MultiWriter)
	m.Define("NewSectionReader", pkg.NewSectionReader)
	m.Define("Pipe", pkg.Pipe)
	m.Define("ReadAtLeast", pkg.ReadAtLeast)
	m.Define("ReadFull", pkg.ReadFull)
	m.Define("TeeReader", pkg.TeeReader)
	m.Define("WriteString", pkg.WriteString)
	return m
}
