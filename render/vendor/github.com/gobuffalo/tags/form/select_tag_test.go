package form_test

import (
	"testing"

	"github.com/gobuffalo/tags"
	"github.com/gobuffalo/tags/form"
	"github.com/stretchr/testify/require"
)

func Test_SelectTag(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	s := f.SelectTag(tags.Options{})
	r.Equal(`<select></select>`, s.String())
}

func Test_SelectTag_WithSelectOptions(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	st := f.SelectTag(tags.Options{
		"options": form.SelectOptions{
			{Value: 1, Label: "one"},
			{Value: 2, Label: "two"},
		},
	})
	s := st.String()
	r.Contains(s, `<option value="1">one</option>`)
	r.Contains(s, `<option value="2">two</option>`)
}

func Test_SelectTag_WithSelectOptions_Selected(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	st := f.SelectTag(tags.Options{
		"options": form.SelectOptions{
			{Value: 3, Label: "three"},
			{Value: 2, Label: "two"},
		},
		"value": "3",
	})
	s := st.String()
	r.Contains(s, `<option value="3" selected>three</option>`)
	r.Contains(s, `<option value="2">two</option>`)
}

func Test_SelectTag_WithMap(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	st := f.SelectTag(tags.Options{
		"options": map[string]interface{}{
			"one": 1,
			"two": 2,
		},
	})
	s := st.String()
	r.Contains(s, `<option value="1">one</option>`)
	r.Contains(s, `<option value="2">two</option>`)
}

func Test_SelectTag_WithMap_Selected(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	st := f.SelectTag(tags.Options{
		"options": map[string]interface{}{
			"three": 3,
			"two":   2,
		},
		"value": 3,
	})
	s := st.String()
	r.Contains(s, `<option value="3" selected>three</option>`)
	r.Contains(s, `<option value="2">two</option>`)
}

func Test_SelectTag_WithSlice(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	st := f.SelectTag(tags.Options{
		"options": []string{"one", "two"},
	})
	s := st.String()
	r.Contains(s, `<option value="one">one</option>`)
	r.Contains(s, `<option value="two">two</option>`)
}

func Test_SelectTag_WithSlice_Selected(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	st := f.SelectTag(tags.Options{
		"options": []string{"one", "two"},
		"value":   "two",
	})
	s := st.String()
	r.Contains(s, `<option value="one">one</option>`)
	r.Contains(s, `<option value="two" selected>two</option>`)
}

func Test_SelectTag_WithSlice_Selectable(t *testing.T) {
	r := require.New(t)
	f := form.New(tags.Options{})
	st := f.SelectTag(tags.Options{
		"options": []SelectableModel{
			{"John", "1"},
			{"Peter", "2"},
		},
		"value": "1",
	})
	s := st.String()
	r.Contains(s, `<option value="1" selected>John</option>`)
	r.Contains(s, `<option value="2">Peter</option>`)
}

type SelectableModel struct {
	Name string
	ID   string
}

func (sm SelectableModel) SelectLabel() string {
	return sm.Name
}

func (sm SelectableModel) SelectValue() interface{} {
	return sm.ID
}
