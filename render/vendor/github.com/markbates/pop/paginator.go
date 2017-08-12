package pop

import (
	"encoding/json"
	"strconv"

	"github.com/markbates/going/defaults"
)

var PaginatorPerPageDefault = 20
var PaginatorPageKey = "page"
var PaginatorPerPageKey = "per_page"

// Paginator is a type used to represent the pagination of records
// from the database.
type Paginator struct {
	// Current page you're on
	Page int `json:"page"`
	// Number of results you want per page
	PerPage int `json:"per_page"`
	// Page * PerPage (ex: 2 * 20, Offset == 40)
	Offset int `json:"offset"`
	// Total potential records matching the query
	TotalEntriesSize int `json:"total_entries_size"`
	// Total records returns, will be <= PerPage
	CurrentEntriesSize int `json:"current_entries_size"`
	// Total pages
	TotalPages int `json:"total_pages"`
}

func (p Paginator) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

// NewPaginator returns a new `Paginator` value with the appropriate
// defaults set.
func NewPaginator(page int, per_page int) *Paginator {
	if page < 1 {
		page = 1
	}
	if per_page < 1 {
		per_page = 20
	}
	p := &Paginator{Page: page, PerPage: per_page}
	p.Offset = (page - 1) * p.PerPage
	return p
}

type PaginationParams interface {
	Get(key string) string
}

// NewPaginatorFromParams takes an interface of type `PaginationParams`,
// the `url.Values` type works great with this interface, and returns
// a new `Paginator` based on the params or `PaginatorPageKey` and
// `PaginatorPerPageKey`. Defaults are `1` for the page and
// PaginatorPerPageDefault for the per page value.
func NewPaginatorFromParams(params PaginationParams) *Paginator {
	page := defaults.String(params.Get("page"), "1")

	per_page := defaults.String(params.Get("per_page"), strconv.Itoa(PaginatorPerPageDefault))

	p, err := strconv.Atoi(page)
	if err != nil {
		p = 1
	}

	pp, err := strconv.Atoi(per_page)
	if err != nil {
		pp = PaginatorPerPageDefault
	}
	return NewPaginator(p, pp)
}

// Paginate records returned from the database.
//
//	q := c.Paginate(2, 15)
//	q.All(&[]User{})
//	q.Paginator
func (c *Connection) Paginate(page int, per_page int) *Query {
	return Q(c).Paginate(page, per_page)
}

// Paginate records returned from the database.
//
//	q = q.Paginate(2, 15)
//	q.All(&[]User{})
//	q.Paginator
func (q *Query) Paginate(page int, per_page int) *Query {
	q.Paginator = NewPaginator(page, per_page)
	return q
}

// Paginate records returned from the database.
//
//	q := c.PaginateFromParams(req.URL.Query())
//	q.All(&[]User{})
//	q.Paginator
func (c *Connection) PaginateFromParams(params PaginationParams) *Query {
	return Q(c).PaginateFromParams(params)
}

// Paginate records returned from the database.
//
//	q = q.PaginateFromParams(req.URL.Query())
//	q.All(&[]User{})
//	q.Paginator
func (q *Query) PaginateFromParams(params PaginationParams) *Query {
	q.Paginator = NewPaginatorFromParams(params)
	return q
}
