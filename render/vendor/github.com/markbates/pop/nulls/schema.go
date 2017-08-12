package nulls

import "reflect"

type register func(interface{}, func(string) reflect.Value)

// RegisterWithSchema allows for the nulls package to be used with http://www.gorillatoolkit.org/pkg/schema#Converter
func RegisterWithSchema(reg register) {
	reg(String{}, func(s string) reflect.Value {
		ns := String{}
		ns.Scan(s)
		return reflect.ValueOf(ns)
	})
	reg(Bool{}, func(s string) reflect.Value {
		ns := Bool{}
		ns.Scan(s)
		return reflect.ValueOf(ns)
	})
	reg(ByteSlice{}, func(s string) reflect.Value {
		ns := ByteSlice{}
		ns.Scan(s)
		return reflect.ValueOf(ns)
	})
	reg(Float32{}, func(s string) reflect.Value {
		ns := Float32{}
		ns.Scan(s)
		return reflect.ValueOf(ns)
	})
	reg(Float64{}, func(s string) reflect.Value {
		ns := Float64{}
		ns.Scan(s)
		return reflect.ValueOf(ns)
	})
	reg(Int{}, func(s string) reflect.Value {
		ns := Int{}
		ns.Scan(s)
		return reflect.ValueOf(ns)
	})
	reg(Int32{}, func(s string) reflect.Value {
		ns := Int32{}
		ns.Scan(s)
		return reflect.ValueOf(ns)
	})
	reg(Int64{}, func(s string) reflect.Value {
		ns := Int64{}
		ns.Scan(s)
		return reflect.ValueOf(ns)
	})
	reg(Time{}, func(s string) reflect.Value {
		ns := Time{}
		ns.Scan(s)
		return reflect.ValueOf(ns)
	})
	reg(UInt32{}, func(s string) reflect.Value {
		ns := UInt32{}
		ns.Scan(s)
		return reflect.ValueOf(ns)
	})
}
