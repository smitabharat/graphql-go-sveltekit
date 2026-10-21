# GraphQL (Go) + SvelteKit POC

A tiny Todo app: a Go GraphQL API backed by PostgreSQL, consumed by a SvelteKit UI.

```
server/                     Go GraphQL API (graphql-go), :8080
  main.go                   wiring only: config -> db -> store -> schema -> HTTP server
  internal/db/              pgx connection pool + schema migration (creates `todos` table)
  internal/todo/            domain model + PostgreSQL store (business rules)
  internal/graph/           GraphQL schema (schema.go) + resolvers (resolver.go)
  internal/httpserver/      routes, GraphiQL, CORS middleware
web/                        SvelteKit 3 app – calls the API from server load + form actions, :5173
```

## Database setup (once)

Create an empty database; the server creates the `todos` table itself on startup.

```powershell
& "C:\Program Files\PostgreSQL\17\bin\createdb.exe" -U postgres -p 5432 graphql_poc
```

## Run

```powershell
# terminal 1 – API
cd server
$env:DATABASE_URL = "postgres://postgres:YOUR_PASSWORD@localhost:5432/graphql_poc?sslmode=disable"
go run .            # http://localhost:8080/graphql  (GraphiQL in the browser)

# terminal 2 – UI
cd web
npm install
npm run dev         # http://localhost:5173
```

| Variable       | Used by | Default |
|----------------|---------|---------|
| `DATABASE_URL` | server  | `postgres://postgres:postgres@localhost:5432/graphql_poc?sslmode=disable` |
| `PORT`         | server  | `8080` |
| `GRAPHQL_URL`  | web     | `http://localhost:8080/graphql` |

## Schema

```graphql
type Todo { id: ID!  title: String!  done: Boolean! }

type Query {
  todos: [Todo!]!
  todo(id: ID!): Todo
}

type Mutation {
  addTodo(title: String!): Todo!
  toggleTodo(id: ID!): Todo!
  deleteTodo(id: ID!): Boolean!
}
```

## How the integration works

- `web/src/lib/server/graphql.ts` – minimal `gql()` client (POST `{query, variables}`) plus the query/mutation documents.
- `web/src/routes/+page.server.ts` – `load` runs the `todos` query; the `add` / `toggle` / `delete` form actions run mutations.
- `web/src/routes/+page.svelte` – renders the list; forms use `use:enhance`, so updates happen without a full reload (and still work with JS disabled).

Because the API is called server-side, the browser never talks to Go directly. CORS is still enabled on the Go server in case you want to call it from the client.
