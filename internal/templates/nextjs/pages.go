package nextjs

// Pages Router Templates (Traditional Next.js with pages directory)
const (
	// TypeScript version for Pages Router
	NextjsPagesPackageJsonTS = `{
  "name": "<<TILO:project_name>>",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint"
  },
  "dependencies": {
    "next": "<<TILO:next_version>>",
    "react": "<<TILO:react_version>>",
    "react-dom": "<<TILO:react_dom_version>>"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "@types/react": "^18.0.0",
    "@types/react-dom": "^18.0.0",
    "autoprefixer": "^10.0.1",
    "eslint": "^8.0.0",
    "eslint-config-next": "<<TILO:next_version>>",
    "postcss": "^8.4.0",
    "tailwindcss": "^3.3.0",
    "typescript": "^5.0.0"
  }
}`

	// JavaScript version removed - Next.js only supports TypeScript
	NextjsPagesPackageJsonJS = ``

	// Main index page for Pages Router
	NextjsPagesIndexPage = `export default function Home() {
  return (
    <main className="container mx-auto px-4 py-8">
      <h1 className="text-4xl font-bold text-center mb-8">
        Welcome to <<TILO:project_name>>
      </h1>
      <p className="text-lg text-center text-gray-600">
        Get started by editing <code className="bg-gray-100 px-2 py-1 rounded">pages/index.tsx</code>
      </p>
      <div className="flex justify-center mt-8">
        <a 
          href="/about" 
          className="bg-blue-500 hover:bg-blue-600 text-white px-6 py-3 rounded-lg transition-colors"
        >
          Learn More
        </a>
      </div>
    </main>
  )
}`

	// App component for Pages Router
	NextjsPagesAppPage = `import type { AppProps } from 'next/app'
import '../styles/globals.css'

export default function App({ Component, pageProps }: AppProps) {
  return <Component {...pageProps} />
}`

	// About page for Pages Router
	NextjsPagesAboutPage = `export default function About() {
  return (
    <main className="container mx-auto px-4 py-8">
      <h1 className="text-4xl font-bold text-center mb-8">
        About <<TILO:project_name>>
      </h1>
      <p className="text-lg text-center text-gray-600 mb-8">
        This is a Next.js application using <strong>Pages Router</strong> with TypeScript and Tailwind CSS.
      </p>
      <div className="flex justify-center">
        <a 
          href="/" 
          className="bg-green-500 hover:bg-green-600 text-white px-6 py-3 rounded-lg transition-colors"
        >
          Go Home
        </a>
      </div>
    </main>
  )
}`

	// Global CSS with Tailwind for Pages Router
	NextjsPagesGlobalCSS = `@tailwind base;
@tailwind components;
@tailwind utilities;

/* Custom styles */
body {
  font-family: 'Inter', sans-serif;
}
`

	// Tailwind config for Pages Router
	NextjsPagesTailwindConfig = `/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './pages/**/*.{js,ts,jsx,tsx}',
    './components/**/*.{js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
`

	// PostCSS config for Pages Router
	NextjsPagesPostCSSConfig = `module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
`
)
