import { error, fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import {
	gql,
	TODOS_QUERY,
	ADD_TODO,
	TOGGLE_TODO,
	DELETE_TODO,
	type Todo
} from '#lib/server/graphql.ts';

export const load: PageServerLoad = async ({ fetch }) => {
	try {
		const data = await gql<{ todos: Todo[] }>(fetch, TODOS_QUERY);
		return { todos: data.todos };
	} catch (e) {
		error(502, `Could not reach GraphQL API: ${(e as Error).message}`);
	}
};

export const actions: Actions = {
	add: async ({ request, fetch }) => {
		const title = String((await request.formData()).get('title') ?? '').trim();
		if (!title) return fail(400, { message: 'Title is required' });

		try {
			await gql(fetch, ADD_TODO, { title });
		} catch (e) {
			return fail(500, { message: (e as Error).message });
		}
	},

	toggle: async ({ request, fetch }) => {
		const id = String((await request.formData()).get('id'));
		try {
			await gql(fetch, TOGGLE_TODO, { id });
		} catch (e) {
			return fail(500, { message: (e as Error).message });
		}
	},

	delete: async ({ request, fetch }) => {
		const id = String((await request.formData()).get('id'));
		try {
			await gql(fetch, DELETE_TODO, { id });
		} catch (e) {
			return fail(500, { message: (e as Error).message });
		}
	}
};
