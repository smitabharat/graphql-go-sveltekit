// Package todo contains the Todo domain model and its PostgreSQL-backed storage.
package todo

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a todo with the given ID does not exist.
var ErrNotFound = errors.New("todo not found")

// ErrEmptyTitle is returned when creating a todo without a title.
var ErrEmptyTitle = errors.New("title must not be empty")

// Todo is the domain model exposed over GraphQL.
type Todo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// Store persists todos in the "todos" table.
type Store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) List(ctx context.Context) ([]*Todo, error) {
	rows, err := s.pool.Query(ctx, `SELECT id::text, title, done FROM todos ORDER BY id`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, scanTodo)
}

// Get returns nil (and no error) when no row matches the ID.
func (s *Store) Get(ctx context.Context, id string) (*Todo, error) {
	key, ok := parseID(id)
	if !ok {
		return nil, nil
	}
	rows, err := s.pool.Query(ctx, `SELECT id::text, title, done FROM todos WHERE id = $1`, key)
	if err != nil {
		return nil, err
	}
	t, err := pgx.CollectExactlyOneRow(rows, scanTodo)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return t, err
}

func (s *Store) Add(ctx context.Context, title string) (*Todo, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}
	rows, err := s.pool.Query(ctx,
		`INSERT INTO todos (title) VALUES ($1) RETURNING id::text, title, done`, title)
	if err != nil {
		return nil, err
	}
	return pgx.CollectExactlyOneRow(rows, scanTodo)
}

func (s *Store) Toggle(ctx context.Context, id string) (*Todo, error) {
	key, ok := parseID(id)
	if !ok {
		return nil, ErrNotFound
	}
	rows, err := s.pool.Query(ctx,
		`UPDATE todos SET done = NOT done WHERE id = $1 RETURNING id::text, title, done`, key)
	if err != nil {
		return nil, err
	}
	t, err := pgx.CollectExactlyOneRow(rows, scanTodo)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return t, err
}

func (s *Store) Delete(ctx context.Context, id string) (bool, error) {
	key, ok := parseID(id)
	if !ok {
		return false, nil
	}
	tag, err := s.pool.Exec(ctx, `DELETE FROM todos WHERE id = $1`, key)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func scanTodo(row pgx.CollectableRow) (*Todo, error) {
	var t Todo
	err := row.Scan(&t.ID, &t.Title, &t.Done)
	return &t, err
}

// parseID converts a GraphQL ID to the table's BIGINT key.
func parseID(id string) (int64, bool) {
	n, err := strconv.ParseInt(id, 10, 64)
	return n, err == nil
}
