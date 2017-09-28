package middleware

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/markbates/willie"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

type widget struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func tx(fn func(tx *pop.Connection)) error {
	pop.Debug = true
	defer func() { pop.Debug = false }()
	d, err := ioutil.TempDir("", "")
	if err != nil {
		return errors.WithStack(err)
	}
	path := filepath.Join(d, "pt_test.sqlite")
	defer os.RemoveAll(path)

	db, err := pop.NewConnection(&pop.ConnectionDetails{
		Dialect: "sqlite",
		URL:     path,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	if err := db.Dialect.CreateDB(); err != nil {
		return errors.WithStack(err)
	}
	if err := db.Open(); err != nil {
		return err
	}
	if err := db.RawQuery(mig).Exec(); err != nil {
		return err
	}
	fn(db)
	return nil
}

func app(db *pop.Connection) *buffalo.App {
	app := buffalo.New(buffalo.Options{})
	app.Use(PopTransaction(db))
	app.GET("/success", func(c buffalo.Context) error {
		w := &widget{}
		tx := c.Value("tx").(*pop.Connection)
		if err := tx.Create(w); err != nil {
			return err
		}
		return c.Render(201, nil)
	})
	app.GET("/non-success", func(c buffalo.Context) error {
		w := &widget{}
		tx := c.Value("tx").(*pop.Connection)
		if err := tx.Create(w); err != nil {
			return err
		}
		return c.Render(301, nil)
	})
	app.GET("/error", func(c buffalo.Context) error {
		w := &widget{}
		tx := c.Value("tx").(*pop.Connection)
		if err := tx.Create(w); err != nil {
			return err
		}
		return errors.New("boom")
	})
	return app
}

func Test_PopTransaction(t *testing.T) {
	r := require.New(t)
	err := tx(func(db *pop.Connection) {
		w := willie.New(app(db))
		res := w.HTML("/success").Get()
		r.Equal(201, res.Code)
		count, err := db.Count("widgets")
		r.NoError(err)
		r.Equal(1, count)
	})
	r.NoError(err)
}

func Test_PopTransaction_Error(t *testing.T) {
	r := require.New(t)
	err := tx(func(db *pop.Connection) {
		w := willie.New(app(db))
		res := w.HTML("/error").Get()
		r.Equal(500, res.Code)
		count, err := db.Count("widgets")
		r.NoError(err)
		r.Equal(0, count)
	})
	r.NoError(err)
}

func Test_PopTransaction_NonSuccess(t *testing.T) {
	r := require.New(t)
	err := tx(func(db *pop.Connection) {
		w := willie.New(app(db))
		res := w.HTML("/non-success").Get()
		r.Equal(301, res.Code)
		count, err := db.Count("widgets")
		r.NoError(err)
		r.Equal(1, count)
	})
	r.NoError(err)
}

const mig = `CREATE TABLE "widgets" (
  "created_at" DATETIME NOT NULL,
  "updated_at" DATETIME NOT NULL,
  "id" TEXT PRIMARY KEY
);`
