package pop

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	_ "github.com/lib/pq"
	. "github.com/markbates/pop/columns"
	"github.com/markbates/pop/fizz"
	"github.com/markbates/pop/fizz/translators"
	"github.com/pkg/errors"
)

var _ dialect = &postgresql{}

type postgresql struct {
	translateCache    map[string]string
	mu                sync.Mutex
	ConnectionDetails *ConnectionDetails
}

func (p *postgresql) Details() *ConnectionDetails {
	return p.ConnectionDetails
}

func (p *postgresql) Create(s store, model *Model, cols Columns) error {
	keyType := model.PrimaryKeyType()
	switch keyType {
	case "int", "int64":
		cols.Remove("id")
		id := struct {
			ID int `db:"id"`
		}{}
		w := cols.Writeable()
		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) returning id", model.TableName(), w.String(), w.SymbolizedString())
		Log(query)
		stmt, err := s.PrepareNamed(query)
		if err != nil {
			return errors.Wrapf(err, "postgres error preparing insert statement %s", query)
		}
		err = stmt.Get(&id, model.Value)
		if err != nil {
			return errors.Wrap(err, "postgres error inserting record")
		}
		model.setID(id.ID)
		return nil
	case "UUID":
		return genericCreate(s, model, cols)
	}
	return errors.Errorf("can not use %s as a primary key type!", keyType)
}

func (p *postgresql) Update(s store, model *Model, cols Columns) error {
	return errors.Wrap(genericUpdate(s, model, cols), "postgres update")
}

func (p *postgresql) Destroy(s store, model *Model) error {
	return errors.Wrap(genericDestroy(s, model), "postgres destroy")
}

func (p *postgresql) SelectOne(s store, model *Model, query Query) error {
	return errors.Wrap(genericSelectOne(s, model, query), "postgres select one")
}

func (p *postgresql) SelectMany(s store, models *Model, query Query) error {
	return errors.Wrap(genericSelectMany(s, models, query), "postgres select many")
}

func (p *postgresql) CreateDB() error {
	// createdb -h db -p 5432 -U postgres enterprise_development
	deets := p.ConnectionDetails
	cmd := exec.Command("createdb", "-e", "-h", deets.Host, "-p", deets.Port, "-U", deets.User, deets.Database)
	Log(strings.Join(cmd.Args, " "))
	comboOut, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%s: %s", err.Error(), string(comboOut))
		return errors.Wrapf(err, "error creating PostgreSQL database %s", deets.Database)
	}

	fmt.Printf("created database %s\n", deets.Database)
	return nil
}

func (p *postgresql) DropDB() error {
	deets := p.ConnectionDetails
	cmd := exec.Command("dropdb", "-e", "-h", deets.Host, "-p", deets.Port, "-U", deets.User, deets.Database)
	Log(strings.Join(cmd.Args, " "))
	comboOut, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%s: %s", err.Error(), string(comboOut))
		return errors.Wrapf(err, "error dropping PostgreSQL database %s", deets.Database)
	}
	fmt.Printf("dropped database %s\n", deets.Database)
	return nil
}

func (m *postgresql) URL() string {
	c := m.ConnectionDetails
	if c.URL != "" {
		return c.URL
	}

	s := "postgres://%s:%s@%s:%s/%s?sslmode=disable"
	return fmt.Sprintf(s, c.User, c.Password, c.Host, c.Port, c.Database)
}

func (m *postgresql) MigrationURL() string {
	return m.URL()
}

func (p *postgresql) TranslateSQL(sql string) string {
	defer p.mu.Unlock()
	p.mu.Lock()

	if csql, ok := p.translateCache[sql]; ok {
		return csql
	}
	curr := 1
	out := make([]byte, 0, len(sql))
	for i := 0; i < len(sql); i++ {
		if sql[i] == '?' {
			str := "$" + strconv.Itoa(curr)
			for _, char := range str {
				out = append(out, byte(char))
			}
			curr += 1
		} else {
			out = append(out, sql[i])
		}
	}
	csql := string(out)
	p.translateCache[sql] = csql
	return csql
}

func (p *postgresql) FizzTranslator() fizz.Translator {
	return translators.NewPostgres()
}

func (p *postgresql) Lock(fn func() error) error {
	return fn()
}

func (p *postgresql) DumpSchema(w io.Writer) error {
	cmd := exec.Command("pg_dump", "-s", fmt.Sprintf("--dbname=%s", p.URL()))
	Log(strings.Join(cmd.Args, " "))
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Printf("dumped schema for %s\n", p.Details().Database)
	return nil
}

func (p *postgresql) LoadSchema(r io.Reader) error {
	cmd := exec.Command("psql", p.URL())
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer in.Close()
		io.Copy(in, r)
	}()
	Log(strings.Join(cmd.Args, " "))
	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}

	fmt.Printf("loaded schema for %s\n", p.Details().Database)
	return nil
}

func (p *postgresql) TruncateAll(tx *Connection) error {
	return tx.RawQuery(pgTruncate).Exec()
}

func newPostgreSQL(deets *ConnectionDetails) dialect {
	cd := &postgresql{
		ConnectionDetails: deets,
		translateCache:    map[string]string{},
		mu:                sync.Mutex{},
	}
	return cd
}

const pgTruncate = `DO
$func$
BEGIN
   EXECUTE
  (SELECT 'TRUNCATE TABLE '
       || string_agg(quote_ident(schemaname) || '.' || quote_ident(tablename), ', ')
       || ' CASCADE'
   FROM   pg_tables
   WHERE  schemaname = 'public'
   );
END
$func$;`
