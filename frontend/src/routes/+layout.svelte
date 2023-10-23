<script lang="ts">
	import type { LayoutData } from './$types';
	import '../app.css';
	import { url } from '$lib/url';
	import { onMount } from 'svelte';
	export let data: LayoutData;

	onMount(() => {
		console.log(url);
		fetch(url.produtos, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		})
			.then((res) => res.json())
			.then((data) => {
				console.log(data);
			})
			.catch((err) => {
				console.log(err);
			});
	});
</script>

<slot />
