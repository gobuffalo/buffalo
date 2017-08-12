package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gobuffalo/makr"
	"github.com/pkg/errors"

	"github.com/markbates/going/defaults"
	"github.com/markbates/inflect"
	"github.com/markbates/pop"
	"github.com/spf13/cobra"
)

var skipMigration bool

func init() {
	ModelCmd.Flags().BoolVarP(&skipMigration, "skip-migration", "s", false, "Skip creating a new fizz migration for this model.")
}

var nrx = regexp.MustCompile(`^nulls\.(.+)`)

type names struct {
	Original string
	Table    string
	Proper   string
	File     string
	Plural   string
	Char     string
}

func newName(name string) names {
	return names{
		Original: name,
		File:     inflect.Underscore(inflect.Singularize(name)),
		Table:    inflect.Tableize(inflect.Pluralize(name)),
		Proper:   inflect.ForeignKeyToAttribute(name),
		Plural:   inflect.Pluralize(inflect.Camelize(name)),
		Char:     strings.ToLower(string([]byte(name)[0])),
	}
}

type attribute struct {
	Names        names
	OriginalType string
	GoType       string
	Nullable     bool
}

func (a attribute) String() string {
	return fmt.Sprintf("\t%s %s `json:\"%s\" db:\"%s\"`", a.Names.Proper, a.GoType, a.Names.Original, a.Names.Original)
}

type model struct {
	Package               string
	Imports               []string
	Names                 names
	Attributes            []attribute
	ValidatableAttributes []attribute
}

func (m model) Generate() error {
	g := makr.New()
	ctx := makr.Data{}
	ctx["model"] = m
	ctx["plural_model_name"] = m.Names.Plural
	ctx["model_name"] = m.Names.Proper
	ctx["package_name"] = m.Package
	ctx["char"] = m.Names.Char

	fname := filepath.Join(m.Package, m.Names.File+".go")
	g.Add(makr.NewFile(fname, modelTemplate))
	tfname := filepath.Join(m.Package, m.Names.File+"_test.go")
	g.Add(makr.NewFile(tfname, modelTestTemplate))
	return g.Run(".", ctx)
}

func (m model) Fizz() string {
	s := []string{fmt.Sprintf("create_table(\"%s\", func(t) {", m.Names.Table)}
	for _, a := range m.Attributes {
		switch a.Names.Original {
		case "created_at", "updated_at":
		case "id":
			s = append(s, fmt.Sprintf("\tt.Column(\"id\", \"%s\", {\"primary\": true})", fizzColType(a.OriginalType)))
		default:
			x := fmt.Sprintf("\tt.Column(\"%s\", \"%s\", {})", a.Names.Original, fizzColType(a.OriginalType))
			if a.Nullable {
				x = strings.Replace(x, "{}", `{"null": true}`, -1)
			}
			s = append(s, x)
		}
	}
	s = append(s, "})")
	return strings.Join(s, "\n")
}

func newModel(name string) model {
	id := newName("id")
	id.Proper = "ID"
	m := model{
		Package: "models",
		Imports: []string{"time", "encoding/json", "github.com/markbates/pop", "github.com/markbates/validate", "github.com/satori/go.uuid"},
		Names:   newName(name),
		Attributes: []attribute{
			{Names: id, OriginalType: "uuid.UUID", GoType: "uuid.UUID"},
			{Names: newName("created_at"), OriginalType: "time.Time", GoType: "time.Time"},
			{Names: newName("updated_at"), OriginalType: "time.Time", GoType: "time.Time"},
		},
		ValidatableAttributes: []attribute{},
	}
	m.Names.Proper = inflect.Singularize(m.Names.Proper)
	return m
}

var ModelCmd = &cobra.Command{
	Use:     "model [name]",
	Aliases: []string{"m"},
	Short:   "Generates a model for your database",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("You must supply a name for your model!")
		}

		model := newModel(args[0])

		hasNulls := false
		for _, def := range args[1:] {
			col := strings.Split(def, ":")
			if len(col) == 1 {
				col = append(col, "string")
			}
			nullable := nrx.MatchString(col[1])
			if !hasNulls && nullable {
				hasNulls = true
				model.Imports = append(model.Imports, "github.com/markbates/pop/nulls")
			}
			if strings.HasPrefix(col[1], "slices.") {
				model.Imports = append(model.Imports, "github.com/markbates/pop/slices")
			}

			a := attribute{
				Names:        newName(col[0]),
				OriginalType: col[1],
				GoType:       colType(col[1]),
				Nullable:     nullable,
			}
			model.Attributes = append(model.Attributes, a)
			if !a.Nullable {
				if a.GoType == "string" || a.GoType == "time.Time" || a.GoType == "int" {
					if a.GoType == "time.Time" {
						a.GoType = "Time"
					}
					model.ValidatableAttributes = append(model.ValidatableAttributes, a)
				}
			}
		}

		err := os.MkdirAll(model.Package, 0766)
		if err != nil {
			return errors.Wrapf(err, "couldn't create folder %s", model.Package)
		}

		err = model.Generate()
		if err != nil {
			return err
		}

		if !skipMigration {
			cflag := cmd.Flag("path")
			migrationPath := defaults.String(cflag.Value.String(), "./migrations")
			err = pop.MigrationCreate(migrationPath, fmt.Sprintf("create_%s", model.Names.Table), "fizz", []byte(model.Fizz()), []byte(fmt.Sprintf("drop_table(\"%s\")", model.Names.Table)))
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func colType(s string) string {
	switch strings.ToLower(s) {
	case "text":
		return "string"
	case "time", "timestamp", "datetime":
		return "time.Time"
	case "nulls.text":
		return "nulls.String"
	case "uuid":
		return "uuid.UUID"
	case "json", "jsonb":
		return "slices.Map"
	case "[]string":
		return "slices.String"
	case "[]int":
		return "slices.Int"
	case "slices.float", "[]float", "[]float32", "[]float64":
		return "slices.Float"
	default:
		return s
	}
}

func fizzColType(s string) string {
	switch strings.ToLower(s) {
	case "int":
		return "integer"
	case "time", "datetime":
		return "timestamp"
	case "uuid.uuid", "uuid":
		return "uuid"
	case "nulls.float32", "nulls.float64":
		return "float"
	case "slices.string", "[]string":
		return "varchar[]"
	case "slices.float", "[]float", "[]float32", "[]float64":
		return "numeric[]"
	case "slices.int":
		return "int[]"
	case "slices.map":
		return "jsonb"
	default:
		if nrx.MatchString(s) {
			return fizzColType(strings.Replace(s, "nulls.", "", -1))
		}
		return strings.ToLower(s)
	}
}
