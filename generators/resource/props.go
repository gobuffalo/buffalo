package resource

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/markbates/going/randx"
	"github.com/markbates/inflect"
)

// Prop of a model. Starts as name:type on the command line.
type Prop struct {
	Name         inflect.Name
	Type         string
	OriginalType string
	TestValue    string
}

// String representation of Prop
func (m Prop) String() string {
	return string(m.Name)
}

func modelPropertiesFromArgs(args []string) []Prop {
	var props []Prop
	if len(args) == 0 {
		return props
	}
	for _, a := range args[1:] {
		ax := strings.Split(a, ":")
		p := Prop{
			Name:         inflect.Name(inflect.ForeignKeyToAttribute(ax[0])),
			Type:         "string",
			OriginalType: "string",
		}
		if len(ax) > 1 {
			p.OriginalType = strings.ToLower(ax[1])
			p.Type = strings.ToLower(strings.TrimPrefix(ax[1], "nulls."))
		}
		p.TestValue = setTypeValue(p.OriginalType)
		props = append(props, p)
	}
	return props
}

func setTypeValue(propType string) string {
	switch propType {
	case "string", "text":
		s := fmt.Sprintf("\"%s\"", randx.String(40))
		return s
	case "float", "float32", "float64":
		r := rand.New(rand.NewSource(0)).Float64()
		s := fmt.Sprintf("%.10f", r)
		return s
	case "int", "int32", "int64":
		r := rand.Int31()
		s := strconv.Itoa(int(r))
		return s
	case "time":
		return "time.Now()"
	case "uuid":
		return `uuid.Must(uuid.NewV4())`
	case "bool":
		return "false"
	case "slices.string":
		rn := rand.Intn(4) + 2
		ss := ""

		for i := 0; i < rn; i++ {
			s := randx.String(rn)
			ss = ss + fmt.Sprintf(`"%s",`, s)
		}
		return fmt.Sprintf("[]string{%s}", ss)
	case "slices.int":
		rn := rand.Intn(4) + 2
		is := ""

		for i := 0; i < rn; i++ {
			n := rand.Int31()
			is = is + fmt.Sprintf(`%v,`, n)
		}
		return fmt.Sprintf("[]int{%s}", is)
	case "slices.float":
		rn := rand.Intn(4) + 2
		fs := ""

		for i := 0; i < rn; i++ {
			f := rand.New(rand.NewSource(0)).Float64()
			fs = fs + fmt.Sprintf(`%v,`, f)
		}
		return fmt.Sprintf("[]float64{%s}", fs)
	case "slices.uuid":
		rn := rand.Intn(1) + 1
		us := ""

		for i := 0; i < rn; i++ {
			us = us + `uuid.Must(uuid.NewV4()),`
		}
		return fmt.Sprintf("[]uuid.UUID{%s}", us)
	}
	if strings.Contains(propType, "null") {
		return "null"
	}
	return "unidentified type"
}
