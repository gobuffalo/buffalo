package render

type Widget struct {
	Name string
}

func (w Widget) ToPath() string {
	return w.Name
}
