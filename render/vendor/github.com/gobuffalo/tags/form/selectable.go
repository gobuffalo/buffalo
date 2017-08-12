package form

//Selectable allows any struct to become an option in the select tag.
type Selectable interface {
	SelectValue() interface{}
	SelectLabel() string
}

//Selectables is the plural for selectable
type Selectables []Selectable
