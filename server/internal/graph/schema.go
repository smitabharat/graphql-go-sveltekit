// Package graph defines the GraphQL schema and wires resolvers to the todo store.
package graph

import (
	"github.com/graphql-go/graphql"

	"graphql-poc/internal/todo"
)

var todoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Todo",
	Fields: graphql.Fields{
		"id":    &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
		"title": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"done":  &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
	},
})

var idArg = graphql.FieldConfigArgument{
	"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
}

// NewSchema builds the GraphQL schema backed by the given store.
func NewSchema(store *todo.Store) (graphql.Schema, error) {
	r := &resolver{store: store}

	query := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"todos": &graphql.Field{
				Type:    graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(todoType))),
				Resolve: r.todos,
			},
			"todo": &graphql.Field{
				Type:    todoType,
				Args:    idArg,
				Resolve: r.todo,
			},
		},
	})

	mutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"addTodo": &graphql.Field{
				Type: graphql.NewNonNull(todoType),
				Args: graphql.FieldConfigArgument{
					"title": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: r.addTodo,
			},
			"toggleTodo": &graphql.Field{
				Type:    graphql.NewNonNull(todoType),
				Args:    idArg,
				Resolve: r.toggleTodo,
			},
			"deleteTodo": &graphql.Field{
				Type:    graphql.NewNonNull(graphql.Boolean),
				Args:    idArg,
				Resolve: r.deleteTodo,
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{Query: query, Mutation: mutation})
}
