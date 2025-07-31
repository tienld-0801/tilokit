package frameworks

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/ti-lo/tilokit/internal/core/context"
	"github.com/ti-lo/tilokit/internal/utils"
)

// VuePlugin implements Vue framework support
type VuePlugin struct{}

// NewVuePlugin creates a new Vue plugin instance
func NewVuePlugin() *VuePlugin {
	return &VuePlugin{}
}

func (p *VuePlugin) Name() string {
	return "vue-framework"
}

func (p *VuePlugin) Version() string {
	return "1.0.0"
}

func (p *VuePlugin) Description() string {
	return "Vue 3 framework with Composition API and modern setup"
}

func (p *VuePlugin) SupportedFrameworks() []string {
	return []string{"vue"}
}

func (p *VuePlugin) SupportedBuildTools() []string {
	return []string{"vite", "webpack"}
}

func (p *VuePlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	ctx.SetVariable("vue_version", "^3.4.0")
	ctx.SetVariable("typescript_support", true)
	return nil
}

func (p *VuePlugin) Generate(ctx *tilocontext.ExecutionContext) error {
	if err := p.createDirectoryStructure(ctx); err != nil {
		return errors.Wrap(err, "failed to create directory structure")
	}

	if err := p.generatePackageJson(ctx); err != nil {
		return errors.Wrap(err, "failed to generate package.json")
	}

	if err := p.generateSourceFiles(ctx); err != nil {
		return errors.Wrap(err, "failed to generate source files")
	}

	if err := p.generateConfigFiles(ctx); err != nil {
		return errors.Wrap(err, "failed to generate config files")
	}

	return nil
}

func (p *VuePlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	ctx.SetMetadata("framework_generated", true)
	ctx.SetMetadata("install_command", "npm install")
	ctx.SetMetadata("start_command", "npm run dev")
	return nil
}

func (p *VuePlugin) createDirectoryStructure(ctx *tilocontext.ExecutionContext) error {
	dirs := []string{
		"src",
		"src/components",
		"src/composables",
		"src/stores",
		"src/views",
		"src/assets",
		"src/styles",
		"public",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(ctx.ProjectPath, dir)
		if err := utils.EnsureDir(dirPath); err != nil {
			return err
		}
	}

	return nil
}

func (p *VuePlugin) generatePackageJson(ctx *tilocontext.ExecutionContext) error {
	packageJson := `{
  "name": "` + ctx.Config.ProjectName + `",
  "private": true,
  "version": "0.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc --noEmit && vite build",
    "preview": "vite preview",
    "lint": "eslint . --ext .vue,.js,.jsx,.cjs,.mjs,.ts,.tsx,.cts,.mts --fix --ignore-path .gitignore",
    "type-check": "vue-tsc --noEmit"
  },
  "dependencies": {
    "vue": "^3.4.0",
    "vue-router": "^4.2.5",
    "pinia": "^2.1.7"
  },
  "devDependencies": {
    "@rushstack/eslint-patch": "^1.3.3",
    "@tsconfig/node18": "^18.2.2",
    "@types/node": "^18.18.5",
    "@vitejs/plugin-vue": "^4.4.0",
    "@vue/eslint-config-prettier": "^8.0.0",
    "@vue/eslint-config-typescript": "^12.0.0",
    "@vue/tsconfig": "^0.4.0",
    "eslint": "^8.49.0",
    "eslint-plugin-vue": "^9.17.0",
    "npm-run-all2": "^6.1.1",
    "prettier": "^3.0.3",
    "typescript": "~5.2.0",
    "vite": "^5.0.0",
    "vue-tsc": "^1.8.19"
  }
}`

	packageJsonPath := filepath.Join(ctx.ProjectPath, "package.json")
	return utils.WriteFile(packageJsonPath, packageJson)
}

func (p *VuePlugin) generateSourceFiles(ctx *tilocontext.ExecutionContext) error {
	// Generate main.ts
	mainTs := `import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

import './assets/main.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')`

	// Generate App.vue
	appVue := `<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import HelloWorld from './components/HelloWorld.vue'
</script>

<template>
  <header>
    <div class="wrapper">
      <HelloWorld msg="Welcome to ` + ctx.Config.ProjectName + `" />

      <nav>
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/about">About</RouterLink>
      </nav>
    </div>
  </header>

  <RouterView />
</template>

<style scoped>
header {
  line-height: 1.5;
  max-height: 100vh;
}

nav {
  width: 100%;
  font-size: 12px;
  text-align: center;
  margin-top: 2rem;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: inline-block;
  padding: 0 1rem;
  border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
  border: 0;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }

  nav {
    text-align: left;
    margin-left: -1rem;
    font-size: 1rem;

    padding: 1rem 0;
    margin-top: 1rem;
  }
}
</style>`

	// Generate HelloWorld.vue
	helloWorldVue := `<script setup lang="ts">
defineProps<{
  msg: string
}>()
</script>

<template>
  <div class="greetings">
    <h1 class="green">{{ msg }}</h1>
    <h3>
      You've successfully created a project with
      <a href="https://vitejs.dev/" target="_blank" rel="noopener">Vite</a> +
      <a href="https://vuejs.org/" target="_blank" rel="noopener">Vue 3</a> +
      <a href="https://www.typescriptlang.org/" target="_blank" rel="noopener">TypeScript</a>.
    </h3>
  </div>
</template>

<style scoped>
h1 {
  font-weight: 500;
  font-size: 2.6rem;
  position: relative;
  top: -10px;
}

h3 {
  font-size: 1.2rem;
}

.greetings h1,
.greetings h3 {
  text-align: center;
}

@media (min-width: 1024px) {
  .greetings h1,
  .greetings h3 {
    text-align: left;
  }
}
</style>`

	// Generate router
	routerTs := `import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/AboutView.vue')
    }
  ]
})

export default router`

	// Generate views
	homeViewVue := `<script setup lang="ts">
import TheWelcome from '../components/TheWelcome.vue'
</script>

<template>
  <main>
    <TheWelcome />
  </main>
</template>`

	aboutViewVue := `<template>
  <div class="about">
    <h1>This is an about page</h1>
  </div>
</template>

<style>
@media (min-width: 1024px) {
  .about {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>`

	theWelcomeVue := `<script setup lang="ts">
import WelcomeItem from './WelcomeItem.vue'
import DocumentationIcon from './icons/IconDocumentation.vue'
import ToolingIcon from './icons/IconTooling.vue'
import EcosystemIcon from './icons/IconEcosystem.vue'
import CommunityIcon from './icons/IconCommunity.vue'
import SupportIcon from './icons/IconSupport.vue'
</script>

<template>
  <WelcomeItem>
    <template #icon>
      <DocumentationIcon />
    </template>
    <template #heading>Documentation</template>

    Vue's
    <a href="https://vuejs.org/" target="_blank" rel="noopener">official documentation</a>
    provides you with all information you need to get started.
  </WelcomeItem>

  <WelcomeItem>
    <template #icon>
      <ToolingIcon />
    </template>
    <template #heading>Tooling</template>

    This project is served and bundled with
    <a href="https://vitejs.dev/guide/features.html" target="_blank" rel="noopener">Vite</a>. The
    recommended IDE setup is
    <a href="https://code.visualstudio.com/" target="_blank" rel="noopener">VSCode</a> +
    <a href="https://github.com/johnsoncodehk/volar" target="_blank" rel="noopener">Volar</a>. If
    you need to test your components and web pages, check out
    <a href="https://www.cypress.io/" target="_blank" rel="noopener">Cypress</a> and
    <a href="https://on.cypress.io/component" target="_blank" rel="noopener"
      >Cypress Component Testing</a
    >.
  </WelcomeItem>

  <WelcomeItem>
    <template #icon>
      <EcosystemIcon />
    </template>
    <template #heading>Ecosystem</template>

    Get official tools and libraries for your project:
    <a href="https://pinia.vuejs.org/" target="_blank" rel="noopener">Pinia</a>,
    <a href="https://router.vuejs.org/" target="_blank" rel="noopener">Vue Router</a>,
    <a href="https://test-utils.vuejs.org/" target="_blank" rel="noopener">Vue Test Utils</a>, and
    <a href="https://github.com/vuejs/devtools" target="_blank" rel="noopener">Vue Dev Tools</a>. If
    you need more resources, we suggest paying
    <a href="https://github.com/vuejs/awesome-vue" target="_blank" rel="noopener">Awesome Vue</a>
    a visit.
  </WelcomeItem>

  <WelcomeItem>
    <template #icon>
      <CommunityIcon />
    </template>
    <template #heading>Community</template>

    Got stuck? Ask your question on
    <a href="https://chat.vuejs.org" target="_blank" rel="noopener">Vue Land</a>, our official
    Discord server, or
    <a href="https://stackoverflow.com/questions/tagged/vue.js" target="_blank" rel="noopener"
      >StackOverflow</a
    >. You should also subscribe to
    <a href="https://news.vuejs.org" target="_blank" rel="noopener">our mailing list</a> and follow
    <a href="https://twitter.com/vuejs" target="_blank" rel="noopener">the official Vue Twitter account</a>
    for latest news in the Vue world.
  </WelcomeItem>

  <WelcomeItem>
    <template #icon>
      <SupportIcon />
    </template>
    <template #heading>Support Vue</template>

    As an independent project, Vue relies on community backing for its sustainability. You can help
    us by
    <a href="https://vuejs.org/sponsor/" target="_blank" rel="noopener">becoming a sponsor</a>.
  </WelcomeItem>
</template>`

	// Generate CSS
	mainCss := `@import './base.css';

#app {
  max-width: 1280px;
  margin: 0 auto;
  padding: 2rem;

  font-weight: normal;
}

a,
.green {
  text-decoration: none;
  color: hsla(160, 100%, 37%, 1);
  transition: 0.4s;
}

@media (hover: hover) {
  a:hover {
    background-color: hsla(160, 100%, 37%, 0.2);
  }
}

@media (min-width: 1024px) {
  body {
    display: flex;
    place-items: center;
  }

  #app {
    display: grid;
    grid-template-columns: 1fr 1fr;
    padding: 0 2rem;
  }
}`

	// Write files
	files := map[string]string{
		"src/main.ts":                  mainTs,
		"src/App.vue":                  appVue,
		"src/components/HelloWorld.vue": helloWorldVue,
		"src/router/index.ts":          routerTs,
		"src/views/HomeView.vue":       homeViewVue,
		"src/views/AboutView.vue":      aboutViewVue,
		"src/components/TheWelcome.vue": theWelcomeVue,
		"src/assets/main.css":          mainCss,
	}

	for path, content := range files {
		fullPath := filepath.Join(ctx.ProjectPath, path)
		if err := utils.WriteFile(fullPath, content); err != nil {
			return err
		}
	}

	return nil
}

func (p *VuePlugin) generateConfigFiles(ctx *tilocontext.ExecutionContext) error {
	// Generate tsconfig.json
	tsConfig := `{
  "files": [],
  "references": [
    {
      "path": "./tsconfig.node.json"
    },
    {
      "path": "./tsconfig.app.json"
    }
  ]
}`

	// Generate tsconfig.app.json
	tsConfigApp := `{
  "extends": "@vue/tsconfig/tsconfig.dom.json",
  "include": ["env.d.ts", "src/**/*", "src/**/*.vue"],
  "exclude": ["src/**/__tests__/*"],
  "compilerOptions": {
    "composite": true,
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    },
    "noEmit": true
  }
}`

	// Generate tsconfig.node.json
	tsConfigNode := `{
  "extends": "@tsconfig/node18/tsconfig.json",
  "include": [
    "vite.config.*",
    "vitest.config.*",
    "cypress.config.*",
    "nightwatch.conf.*",
    "playwright.config.*"
  ],
  "compilerOptions": {
    "composite": true,
    "module": "ESNext",
    "types": ["node"]
  }
}`

	// Generate env.d.ts
	envDts := `/// <reference types="vite/client" />`

	// Generate index.html
	indexHtml := `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <link rel="icon" href="/favicon.ico">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + ctx.Config.ProjectName + `</title>
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/src/main.ts"></script>
  </body>
</html>`

	// Write config files
	configs := map[string]string{
		"tsconfig.json":      tsConfig,
		"tsconfig.app.json":  tsConfigApp,
		"tsconfig.node.json": tsConfigNode,
		"src/env.d.ts":       envDts,
		"index.html":         indexHtml,
	}

	for path, content := range configs {
		fullPath := filepath.Join(ctx.ProjectPath, path)
		if err := utils.WriteFile(fullPath, content); err != nil {
			return err
		}
	}

	return nil
}
