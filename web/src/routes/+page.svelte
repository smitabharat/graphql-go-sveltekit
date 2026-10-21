<script lang="ts">
	import { enhance } from '$app/forms';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();

	let remaining = $derived(data.todos.filter((t) => !t.done).length);
</script>

<svelte:head>
	<title>GraphQL Todos</title>
</svelte:head>

<main>
	<h1>Todos <small>Go GraphQL + SvelteKit</small></h1>

	<form method="POST" action="?/add" use:enhance class="add">
		<input name="title" placeholder="What needs doing?" autocomplete="off" required />
		<button>Add</button>
	</form>

	{#if form?.message}
		<p class="error">{form.message}</p>
	{/if}

	<ul>
		{#each data.todos as todo (todo.id)}
			<li class:done={todo.done}>
				<form method="POST" action="?/toggle" use:enhance>
					<input type="hidden" name="id" value={todo.id} />
					<button class="check" aria-label="Toggle {todo.title}">{todo.done ? '✔' : ''}</button>
				</form>
				<span>{todo.title}</span>
				<form method="POST" action="?/delete" use:enhance>
					<input type="hidden" name="id" value={todo.id} />
					<button class="delete" aria-label="Delete {todo.title}">✕</button>
				</form>
			</li>
		{:else}
			<li class="empty">Nothing to do 🎉</li>
		{/each}
	</ul>

	<p class="meta">{remaining} of {data.todos.length} remaining</p>
</main>

<style>
	main {
		max-width: 32rem;
		margin: 3rem auto;
		padding: 0 1rem;
		font-family: system-ui, sans-serif;
	}
	h1 small {
		display: block;
		font-size: 0.9rem;
		font-weight: 400;
		color: #888;
	}
	.add {
		display: flex;
		gap: 0.5rem;
		margin-bottom: 1rem;
	}
	.add input {
		flex: 1;
		padding: 0.5rem;
	}
	ul {
		list-style: none;
		padding: 0;
	}
	li {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.5rem 0;
		border-bottom: 1px solid #eee;
	}
	li span {
		flex: 1;
	}
	li.done span {
		text-decoration: line-through;
		color: #999;
	}
	.check {
		width: 1.5rem;
		height: 1.5rem;
		padding: 0;
		cursor: pointer;
	}
	.delete {
		background: none;
		border: none;
		color: #c33;
		cursor: pointer;
	}
	.error {
		color: #c33;
	}
	.meta,
	.empty {
		color: #888;
	}
</style>
