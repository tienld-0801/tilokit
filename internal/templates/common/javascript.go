package common

const (
	Gitignore = `
# dependencies
/node_modules
/.pnp
.pnp.js

# build output
/dist

# testing
/coverage

# environment variables
.env
.env.*.local

# logs
npm-debug.log*
yarn-debug.log*
yarn-error.log*
pnpm-debug.log*

# misc
.DS_Store
Thumbs.db
.idea/
.vscode/
*.local`
)