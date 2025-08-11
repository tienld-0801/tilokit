package nextjs

const (
	NextjsPackageJson = `{
  "name": "{{.ProjectName}}",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint"
  },
  "dependencies": {
    "next": "14.0.0",
    "react": "^18.0.0",
    "react-dom": "^18.0.0"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "@types/react": "^18.0.0",
    "@types/react-dom": "^18.0.0",
    "eslint": "^8.0.0",
    "eslint-config-next": "14.0.0",
    "typescript": "^5.0.0"
  }
}`

	NextjsIndexPage = `export default function Home() {
  return (
    <main>
      <h1>Welcome to {{.ProjectName}}</h1>
      <p>Get started by editing pages/index.tsx</p>
    </main>
  )
}`

	NextjsAppPage = `import type { AppProps } from 'next/app'

export default function App({ Component, pageProps }: AppProps) {
  return <Component {...pageProps} />
}`
)
