// Package runtime implements runtime interface for anko script.
package runtime

import (
	"github.com/mattn/anko/vm"
	pkg "runtime"
)

func Import(env *vm.Env) *vm.Env {
	m := env.NewModule("runtime")
	//m.Define("BlockProfile", pkg.BlockProfile)
	//m.Define("Breakpoint", pkg.Breakpoint)
	//m.Define("CPUProfile", pkg.CPUProfile)
	//m.Define("Caller", pkg.Caller)
	//m.Define("Callers", pkg.Callers)
	//m.Define("CallersFrames", pkg.CallersFrames)
	//m.Define("Compiler", pkg.Compiler)
	//m.Define("FuncForPC", pkg.FuncForPC)
	m.Define("GC", pkg.GC)
	m.Define("GOARCH", pkg.GOARCH)
	m.Define("GOMAXPROCS", pkg.GOMAXPROCS)
	m.Define("GOOS", pkg.GOOS)
	m.Define("GOROOT", pkg.GOROOT)
	//m.Define("Goexit", pkg.Goexit)
	//m.Define("GoroutineProfile", pkg.GoroutineProfile)
	//m.Define("Gosched", pkg.Gosched)
	//m.Define("LockOSThread", pkg.LockOSThread)
	//m.Define("MemProfile", pkg.MemProfile)
	//m.Define("MemProfileRate", pkg.MemProfileRate)
	//m.Define("NumCPU", pkg.NumCPU)
	//m.Define("NumCgoCall", pkg.NumCgoCall)
	//m.Define("NumGoroutine", pkg.NumGoroutine)
	//m.Define("ReadMemStats", pkg.ReadMemStats)
	//m.Define("ReadTrace", pkg.ReadTrace)
	//m.Define("SetBlockProfileRate", pkg.SetBlockProfileRate)
	//m.Define("SetCPUProfileRate", pkg.SetCPUProfileRate)
	//m.Define("SetFinalizer", pkg.SetFinalizer)
	//m.Define("Stack", pkg.Stack)
	//m.Define("StartTrace", pkg.StartTrace)
	//m.Define("StopTrace", pkg.StopTrace)
	//m.Define("ThreadCreateProfile", pkg.ThreadCreateProfile)
	//m.Define("UnlockOSThread", pkg.UnlockOSThread)
	//m.Define("Version", pkg.Version)
	return m
}
