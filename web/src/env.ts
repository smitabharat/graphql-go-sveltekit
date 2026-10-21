import { defineEnvVars } from '@sveltejs/kit/env';

export const variables = defineEnvVars({
	GRAPHQL_URL: {
		description: 'URL of the Go GraphQL API',
		schema: (value) => value || 'http://localhost:8080/graphql'
	}
});
