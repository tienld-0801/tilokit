package vue

const (
	ViteTsPackageJson = `{
  "name": "{{.ProjectName}}",
  "private": true,
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "vue": "^3.3.0"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^4.0.0",
    "typescript": "^5.0.0",
    "vue-tsc": "^1.4.0",
    "vite": "^4.4.0"
  }
}`

	ViteTsMainFile = `import { createApp } from 'vue'
import App from './App.vue'

createApp(App).mount('#app')`

	ViteTsAppFile = `<template>
  <div id="app">
    <h1>Hello World</h1>
  </div>
</template>

<script setup lang="ts">
// TypeScript setup syntax
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>`

	ViteTsIndexHtml = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="/vite.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{.ProjectName}}</title>
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/src/main.ts"></script>
  </body>
</html>`
)
