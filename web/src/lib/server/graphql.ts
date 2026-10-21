import { GRAPHQL_URL } from '$app/env/private';

export type Todo = {
	id: string;
	title: string;
	done: boolean;
};

type GraphQLResponse<T> = {
	data?: T;
	errors?: { message: string }[];
};

/** Minimal GraphQL client: POSTs a query + variables and returns `data` or throws. */
export async function gql<T>(
	fetchFn: typeof fetch,
	query: string,
	variables: Record<string, unknown> = {}
): Promise<T> {
	const res = await fetchFn(GRAPHQL_URL, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ query, variables })
	});

	if (!res.ok) throw new Error(`GraphQL HTTP ${res.status}`);

	const body = (await res.json()) as GraphQLResponse<T>;
	if (body.errors?.length) throw new Error(body.errors.map((e) => e.message).join('; '));
	return body.data as T;
}

export const TODOS_QUERY = /* GraphQL */ `
	query Todos {
		todos { id title done }
	}
`;

export const ADD_TODO = /* GraphQL */ `
	mutation AddTodo($title: String!) {
		addTodo(title: $title) { id title done }
	}
`;

export const TOGGLE_TODO = /* GraphQL */ `
	mutation ToggleTodo($id: ID!) {
		toggleTodo(id: $id) { id done }
	}
`;

export const DELETE_TODO = /* GraphQL */ `
	mutation DeleteTodo($id: ID!) {
		deleteTodo(id: $id)
	}
`;
