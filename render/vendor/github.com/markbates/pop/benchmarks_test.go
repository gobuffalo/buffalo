package pop_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/markbates/pop"
	"github.com/markbates/pop/nulls"
)

func Benchmark_Create_Pop(b *testing.B) {
	transaction(func(tx *pop.Connection) {
		for n := 0; n < b.N; n++ {
			u := &User{
				Name: nulls.NewString("Mark Bates"),
			}
			tx.Create(u)
		}
	})
}

func Benchmark_Create_Raw(b *testing.B) {
	transaction(func(tx *pop.Connection) {
		for n := 0; n < b.N; n++ {
			u := &User{
				Name: nulls.NewString("Mark Bates"),
			}
			q := "INSERT INTO users (alive, bio, birth_date, created_at, name, price, updated_at) VALUES (:alive, :bio, :birth_date, :created_at, :name, :price, :updated_at)"
			tx.Store.NamedExec(q, u)
		}
	})
}

func Benchmark_Update(b *testing.B) {
	transaction(func(tx *pop.Connection) {
		u := &User{
			Name: nulls.NewString("Mark Bates"),
		}
		tx.Create(u)
		for n := 0; n < b.N; n++ {
			tx.Update(u)
		}
	})
}

func Benchmark_Find_Pop(b *testing.B) {
	transaction(func(tx *pop.Connection) {
		u := &User{
			Name: nulls.NewString("Mark Bates"),
		}
		tx.Create(u)
		for n := 0; n < b.N; n++ {
			tx.Find(u, u.ID)
		}
	})
}

func Benchmark_Find_Raw(b *testing.B) {
	transaction(func(tx *pop.Connection) {
		u := &User{
			Name: nulls.NewString("Mark Bates"),
		}
		tx.Create(u)
		for n := 0; n < b.N; n++ {
			tx.Store.Get(u, "select * from users where id = ?", u.ID)
		}
	})
}

func Benchmark_translateOne(b *testing.B) {
	q := "select * from users where id = ? and name = ? and email = ? and a = ? and b = ? and c = ? and d = ? and e = ? and f = ?"
	for n := 0; n < b.N; n++ {
		translateOne(q)
	}
}

func Benchmark_translateTwo(b *testing.B) {
	q := "select * from users where id = ? and name = ? and email = ? and a = ? and b = ? and c = ? and d = ? and e = ? and f = ?"
	for n := 0; n < b.N; n++ {
		translateTwo(q)
	}
}

func translateOne(sql string) string {
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
	return string(out)
}

func translateTwo(sql string) string {
	curr := 1
	csql := ""
	for i := 0; i < len(sql); i++ {
		x := sql[i]
		if x == '?' {
			csql = fmt.Sprintf("%s$%d", csql, curr)
			curr++
		} else {
			csql += string(x)
		}
	}
	return csql
}
