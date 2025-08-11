package svelte

const (
	SvelteJsPackageJson = `{
  "name": "{{.ProjectName}}",
  "private": true,
  "scripts": {
    "build": "vite build",
    "dev": "vite dev",
    "preview": "vite preview"
  },
  "devDependencies": {
    "@sveltejs/adapter-auto": "^2.0.0",
    "@sveltejs/kit": "^1.20.4",
    "svelte": "^4.0.5",
    "vite": "^4.4.2"
  },
  "type": "module"
}`

	SvelteJsApp = `<script>
	let name = '{{.ProjectName}}';
</script>

<main>
	<h1>Hello {name}!</h1>
	<p>Welcome to your Svelte app</p>
</main>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

	h1 {
		color: #ff3e00;
		text-transform: uppercase;
		font-size: 4em;
		font-weight: 100;
	}
</style>`

	SvelteJsIndexHtml = `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8" />
		<link rel="icon" type="image/svg+xml" href="/vite.svg" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>{{.ProjectName}}</title>
	</head>
	<body>
		<div id="app"></div>
		<script type="module" src="/src/main.js"></script>
	</body>
</html>`
)
