package nodejs

const (
	ExpressPackageJsonTS = `{
  "name": "<<TILO:.project_name>>",
  "version": "1.0.0",
  "description": "",
  "main": "dist/index.js",
  "scripts": {
    "build": "tsc",
    "start": "node dist/index.js",
    "dev": "tsx watch src/index.ts",
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  "dependencies": {
    "express": "^4.21.2",
    "cors": "^2.8.5",
    "dotenv": "^16.4.7"
  },
  "devDependencies": {
    "@types/express": "^5.0.0",
    "@types/cors": "^2.8.17",
    "@types/node": "^22.10.5",
    "tsx": "^4.19.2",
    "typescript": "^5.7.2"
  }
}`

	ExpressPackageJsonJS = `{
  "name": "<<TILO:.project_name>>",
  "version": "1.0.0",
  "description": "",
  "main": "index.js",
  "scripts": {
    "start": "node index.js",
    "dev": "nodemon index.js",
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  "dependencies": {
    "express": "^4.21.2",
    "cors": "^2.8.5",
    "dotenv": "^16.4.7"
  },
  "devDependencies": {
    "nodemon": "^3.1.7"
  }
}`

	ExpressIndexTS = `import express from 'express';
import cors from 'cors';
import dotenv from 'dotenv';

dotenv.config();

const app = express();
const PORT = process.env.PORT || 3000;

// Middleware
app.use(cors());
app.use(express.json());

// Routes
app.get('/', (req, res) => {
  res.json({ message: 'Hello World from <<TILO:>>!' });
});

app.get('/health', (req, res) => {
  res.json({ status: 'ok', service: '<<TILO:.project_name>>' });
});

app.listen(PORT, () => {
  console.log('🚀 <<TILO:.project_name>> is running on http://localhost:' + PORT);
});`

	ExpressIndexJS = `const express = require('express');
const cors = require('cors');
require('dotenv').config();

const app = express();
const PORT = process.env.PORT || 3000;

// Middleware
app.use(cors());
app.use(express.json());

// Routes
app.get('/', (req, res) => {
  res.json({ message: 'Hello World from <<TILO:.project_name>>!' });
});

app.get('/health', (req, res) => {
  res.json({ status: 'ok', service: '<<TILO:.project_name>>' });
});

app.listen(PORT, () => {
  console.log('🚀 <<TILO:.project_name>> is running on http://localhost:' + PORT);
});`

	ExpressTsConfig = `{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "lib": ["ES2020"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist"]
}`

	ExpressEnv = `PORT=3000
NODE_ENV=development`
)
