// Package sort implements sort interface for anko script.
package sort

import (
	s "sort"

	"github.com/mattn/anko/vm"
)

type is []interface{}

func (p is) Len() int           { return len(p) }
func (p is) Less(i, j int) bool { return p[i].(int64) < p[j].(int64) }
func (p is) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type fs []interface{}

func (p fs) Len() int           { return len(p) }
func (p fs) Less(i, j int) bool { return p[i].(float64) < p[j].(float64) }
func (p fs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type ss []interface{}

func (p ss) Len() int           { return len(p) }
func (p ss) Less(i, j int) bool { return p[i].(string) < p[j].(string) }
func (p ss) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func Import(env *vm.Env) *vm.Env {
	m := env.NewPackage("sort")
	m.Define("Ints", func(arr interface{}) interface{} {
		if iarr, ok := arr.([]int); ok {
			s.Ints(iarr)
		} else {
			s.Sort(is(arr.([]interface{})))
		}
		return arr
	})
	m.Define("Float64s", func(arr interface{}) interface{} {
		if farr, ok := arr.([]float64); ok {
			s.Float64s(farr)
		} else {
			s.Sort(fs(arr.([]interface{})))
		}
		return arr
	})
	m.Define("Strings", func(arr interface{}) interface{} {
		if sarr, ok := arr.([]string); ok {
			s.Strings(sarr)
		} else {
			s.Sort(ss(arr.([]interface{})))
		}
		return arr
	})
	handleGo18(m)
	return m
}
