package meta

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Name_Title(t *testing.T) {
	r := require.New(t)
	table := []struct {
		V string
		E string
	}{
		{V: "foo_bar", E: "Foo Bar"},
		{V: "admin/widget", E: "Admin Widget"},
		{V: "widget", E: "Widget"},
	}
	for _, tt := range table {
		r.Equal(tt.E, Name(tt.V).Title())
	}
}

func Test_Name_Model(t *testing.T) {
	r := require.New(t)
	table := []struct {
		V string
		E string
	}{
		{V: "foo_bar", E: "FooBar"},
		{V: "admin/widget", E: "AdminWidget"},
		{V: "widget", E: "Widget"},
		{V: "widgets", E: "Widget"},
	}
	for _, tt := range table {
		r.Equal(tt.E, Name(tt.V).Model())
	}
}

func Test_Name_Resource(t *testing.T) {
	r := require.New(t)
	table := []struct {
		V string
		E string
	}{
		{V: "foo_bar", E: "FooBars"},
		{V: "admin/widget", E: "AdminWidgets"},
		{V: "widget", E: "Widgets"},
		{V: "widgets", E: "Widgets"},
	}
	for _, tt := range table {
		r.Equal(tt.E, Name(tt.V).Resource())
	}
}

func Test_Name_ModelPlural(t *testing.T) {
	r := require.New(t)
	table := []struct {
		V string
		E string
	}{
		{V: "foo_bar", E: "FooBars"},
		{V: "admin/widget", E: "AdminWidgets"},
		{V: "widget", E: "Widgets"},
		{V: "widgets", E: "Widgets"},
	}
	for _, tt := range table {
		r.Equal(tt.E, Name(tt.V).ModelPlural())
	}
}

func Test_Name_File(t *testing.T) {
	r := require.New(t)
	table := []struct {
		V string
		E string
	}{
		{V: "foo_bar", E: "foo_bar"},
		{V: "admin/widget", E: "admin/widget"},
		{V: "widget", E: "widget"},
		{V: "widgets", E: "widgets"},
		{V: "User", E: "user"},
	}
	for _, tt := range table {
		r.Equal(tt.E, Name(tt.V).File())
	}
}

func Test_Name_VarCaseSingular(t *testing.T) {
	r := require.New(t)
	table := []struct {
		V string
		E string
	}{
		{V: "foo_bar", E: "fooBar"},
		{V: "admin/widget", E: "adminWidget"},
		{V: "widget", E: "widget"},
		{V: "widgets", E: "widget"},
		{V: "User", E: "user"},
		{V: "FooBar", E: "fooBar"},
	}
	for _, tt := range table {
		r.Equal(tt.E, Name(tt.V).VarCaseSingular())
	}
}

func Test_Name_VarCasePlural(t *testing.T) {
	r := require.New(t)
	table := []struct {
		V string
		E string
	}{
		{V: "foo_bar", E: "fooBars"},
		{V: "admin/widget", E: "adminWidgets"},
		{V: "widget", E: "widgets"},
		{V: "widgets", E: "widgets"},
		{V: "User", E: "users"},
		{V: "FooBar", E: "fooBars"},
	}
	for _, tt := range table {
		r.Equal(tt.E, Name(tt.V).VarCasePlural())
	}
}
