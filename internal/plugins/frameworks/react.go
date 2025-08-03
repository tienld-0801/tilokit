package frameworks

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/ti-lo/tilokit/internal/core/context"
	"github.com/ti-lo/tilokit/internal/utils"
)

// ReactPlugin implements React framework support
type ReactPlugin struct{}

// NewReactPlugin creates a new React plugin instance
func NewReactPlugin() *ReactPlugin {
	return &ReactPlugin{}
}

func (p *ReactPlugin) Name() string {
	return "react-framework"
}

func (p *ReactPlugin) Version() string {
	return "1.0.0"
}

func (p *ReactPlugin) Description() string {
	return "React framework with modern setup and best practices"
}

func (p *ReactPlugin) SupportedFrameworks() []string {
	return []string{"react"}
}

func (p *ReactPlugin) SupportedBuildTools() []string {
	return []string{"vite", "webpack", "rollup"}
}

func (p *ReactPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set React-specific variables
	ctx.SetVariable("react_version", "^18.2.0")
	ctx.SetVariable("react_dom_version", "^18.2.0")
	ctx.SetVariable("typescript_support", true)
	
	return nil
}

func (p *ReactPlugin) Generate(ctx *tilocontext.ExecutionContext) error {
	// Create directory structure
	if err := p.createDirectoryStructure(ctx); err != nil {
		return errors.Wrap(err, "failed to create directory structure")
	}

	// Generate package.json
	if err := p.generatePackageJson(ctx); err != nil {
		return errors.Wrap(err, "failed to generate package.json")
	}

	// Generate source files
	if err := p.generateSourceFiles(ctx); err != nil {
		return errors.Wrap(err, "failed to generate source files")
	}

	// Generate configuration files
	if err := p.generateConfigFiles(ctx); err != nil {
		return errors.Wrap(err, "failed to generate config files")
	}

	return nil
}

func (p *ReactPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set post-generation metadata
	ctx.SetMetadata("framework_generated", true)
	ctx.SetMetadata("install_command", "npm install")
	ctx.SetMetadata("start_command", "npm run dev")
	
	return nil
}

func (p *ReactPlugin) createDirectoryStructure(ctx *tilocontext.ExecutionContext) error {
	dirs := []string{
		"src",
		"src/components",
		"src/hooks",
		"src/utils",
		"src/styles",
		"src/assets",
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

func (p *ReactPlugin) generatePackageJson(ctx *tilocontext.ExecutionContext) error {
	packageJson := `{
  "name": "` + ctx.Config.ProjectName + `",
  "private": true,
  "version": "0.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "lint": "eslint . --ext ts,tsx --report-unused-disable-directives --max-warnings 0",
    "preview": "vite preview",
    "test": "vitest",
    "test:ui": "vitest --ui"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.20.1"
  },
  "devDependencies": {
    "@types/react": "^18.2.37",
    "@types/react-dom": "^18.2.15",
    "@typescript-eslint/eslint-plugin": "^6.10.0",
    "@typescript-eslint/parser": "^6.10.0",
    "@vitejs/plugin-react": "^4.1.1",
    "eslint": "^8.53.0",
    "eslint-plugin-react-hooks": "^4.6.0",
    "eslint-plugin-react-refresh": "^0.4.4",
    "typescript": "^5.2.2",
    "vite": "^5.0.0",
    "vitest": "^1.0.0",
    "@vitest/ui": "^1.0.0"
  }
}`

	packageJsonPath := filepath.Join(ctx.ProjectPath, "package.json")
	return utils.WriteFile(packageJsonPath, packageJson)
}

func (p *ReactPlugin) generateSourceFiles(ctx *tilocontext.ExecutionContext) error {
	// Generate main.tsx
	mainTsx := `import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './styles/index.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)`

	// Generate App.tsx
	appTsx := `import { useState } from 'react'
import './styles/App.css'

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="app">
      <header className="app-header">
        <h1>Welcome to ` + ctx.Config.ProjectName + `</h1>
        <p>Built with React + TypeScript + Vite</p>
        <div className="card">
          <button onClick={() => setCount((count) => count + 1)}>
            count is {count}
          </button>
          <p>
            Edit <code>src/App.tsx</code> and save to test HMR
          </p>
        </div>
      </header>
    </div>
  )
}

export default App`

	// Generate CSS files
	indexCss := `:root {
  font-family: Inter, system-ui, Avenir, Helvetica, Arial, sans-serif;
  line-height: 1.5;
  font-weight: 400;

  color-scheme: light dark;
  color: rgba(255, 255, 255, 0.87);
  background-color: #242424;

  font-synthesis: none;
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  -webkit-text-size-adjust: 100%;
}

a {
  font-weight: 500;
  color: #646cff;
  text-decoration: inherit;
}
a:hover {
  color: #535bf2;
}

body {
  margin: 0;
  display: flex;
  place-items: center;
  min-width: 320px;
  min-height: 100vh;
}

h1 {
  font-size: 3.2em;
  line-height: 1.1;
}

button {
  border-radius: 8px;
  border: 1px solid transparent;
  padding: 0.6em 1.2em;
  font-size: 1em;
  font-weight: 500;
  font-family: inherit;
  background-color: #1a1a1a;
  color: inherit;
  cursor: pointer;
  transition: border-color 0.25s;
}
button:hover {
  border-color: #646cff;
}
button:focus,
button:focus-visible {
  outline: 4px auto -webkit-focus-ring-color;
}

@media (prefers-color-scheme: light) {
  :root {
    color: #213547;
    background-color: #ffffff;
  }
  a:hover {
    color: #747bff;
  }
  button {
    background-color: #f9f9f9;
  }
}`

	appCss := `.app {
  max-width: 1280px;
  margin: 0 auto;
  padding: 2rem;
  text-align: center;
}

.app-header {
  padding: 20px;
}

.card {
  padding: 2em;
}

.read-the-docs {
  color: #888;
}`

	// Write files
	files := map[string]string{
		"src/main.tsx":        mainTsx,
		"src/App.tsx":         appTsx,
		"src/styles/index.css": indexCss,
		"src/styles/App.css":   appCss,
	}

	for path, content := range files {
		fullPath := filepath.Join(ctx.ProjectPath, path)
		if err := utils.WriteFile(fullPath, content); err != nil {
			return err
		}
	}

	return nil
}

func (p *ReactPlugin) generateConfigFiles(ctx *tilocontext.ExecutionContext) error {
	// Generate tsconfig.json
	tsConfig := `{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,

    /* Bundler mode */
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",

    /* Linting */
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,

    /* Path mapping */
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["src"],
  "references": [{ "path": "./tsconfig.node.json" }]
}`

	// Generate tsconfig.node.json
	tsConfigNode := `{
  "compilerOptions": {
    "composite": true,
    "skipLibCheck": true,
    "module": "ESNext",
    "moduleResolution": "bundler",
    "allowSyntheticDefaultImports": true
  },
  "include": ["vite.config.ts"]
}`

	// Generate index.html
	indexHtml := `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="/vite.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>` + ctx.Config.ProjectName + `</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>`

	// Generate .eslintrc.cjs
	eslintConfig := `module.exports = {
  root: true,
  env: { browser: true, es2020: true },
  extends: [
    'eslint:recommended',
    '@typescript-eslint/recommended',
    'plugin:react-hooks/recommended',
  ],
  ignorePatterns: ['dist', '.eslintrc.cjs'],
  parser: '@typescript-eslint/parser',
  plugins: ['react-refresh'],
  rules: {
    'react-refresh/only-export-components': [
      'warn',
      { allowConstantExport: true },
    ],
  },
}`

	// Write config files
	configs := map[string]string{
		"tsconfig.json":      tsConfig,
		"tsconfig.node.json": tsConfigNode,
		"index.html":         indexHtml,
		".eslintrc.cjs":      eslintConfig,
	}

	for path, content := range configs {
		fullPath := filepath.Join(ctx.ProjectPath, path)
		if err := utils.WriteFile(fullPath, content); err != nil {
			return err
		}
	}

	return nil
}
