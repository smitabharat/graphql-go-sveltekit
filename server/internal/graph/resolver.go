package graph

import (
	"github.com/graphql-go/graphql"

	"graphql-poc/internal/todo"
)

// resolver holds the dependencies GraphQL field resolvers need.
type resolver struct {
	store *todo.Store
}

func (r *resolver) todos(p graphql.ResolveParams) (interface{}, error) {
	return r.store.List(p.Context)
}

func (r *resolver) todo(p graphql.ResolveParams) (interface{}, error) {
	t, err := r.store.Get(p.Context, p.Args["id"].(string))
	if t == nil {
		return nil, err // return an untyped nil so GraphQL renders null
	}
	return t, err
}

func (r *resolver) addTodo(p graphql.ResolveParams) (interface{}, error) {
	return r.store.Add(p.Context, p.Args["title"].(string))
}

func (r *resolver) toggleTodo(p graphql.ResolveParams) (interface{}, error) {
	return r.store.Toggle(p.Context, p.Args["id"].(string))
}

func (r *resolver) deleteTodo(p graphql.ResolveParams) (interface{}, error) {
	return r.store.Delete(p.Context, p.Args["id"].(string))
}
