package nodejs

const (
	FastifyPackageJsonTS = `{
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
    "fastify": "^5.2.0",
    "@fastify/cors": "^10.0.1",
    "@fastify/env": "^5.0.1"
  },
  "devDependencies": {
    "@types/node": "^22.10.5",
    "tsx": "^4.19.2",
    "typescript": "^5.7.2"
  }
}`

	FastifyPackageJsonJS = `{
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
    "fastify": "^5.2.0",
    "@fastify/cors": "^10.0.1",
    "@fastify/env": "^5.0.1"
  },
  "devDependencies": {
    "nodemon": "^3.1.7"
  }
}`

	FastifyIndexTS = `import Fastify from 'fastify';

const fastify = Fastify({
  logger: true
});

// Register CORS plugin
fastify.register(require('@fastify/cors'), {
  origin: true
});

// Register env plugin
fastify.register(require('@fastify/env'), {
  schema: {
    type: 'object',
    required: ['PORT'],
    properties: {
      PORT: {
        type: 'string',
        default: '3000'
      }
    }
  }
});

// Declare a route
fastify.get('/', async (request, reply) => {
  return { message: 'Hello World from <<TILO:.project_name>>!' };
});

// Health check route
fastify.get('/health', async (request, reply) => {
  return { status: 'ok', service: '<<TILO:.project_name>>' };
});

// Run the server!
const start = async () => {
  try {
    const port = Number(process.env.PORT) || 3000;
    await fastify.listen({ port, host: '0.0.0.0' });
    console.log('🚀 <<TILO:.project_name>> is running on http://localhost:' + port);
  } catch (err) {
    fastify.log.error(err);
    process.exit(1);
  }
};

start();`

	FastifyIndexJS = `const fastify = require('fastify')({
  logger: true
});

// Register CORS plugin
fastify.register(require('@fastify/cors'), {
  origin: true
});

// Register env plugin
fastify.register(require('@fastify/env'), {
  schema: {
    type: 'object',
    required: ['PORT'],
    properties: {
      PORT: {
        type: 'string',
        default: '3000'
      }
    }
  }
});

// Declare a route
fastify.get('/', async (request, reply) => {
  return { message: 'Hello World from <<TILO:.project_name>>!' };
});

// Health check route
fastify.get('/health', async (request, reply) => {
  return { status: 'ok', service: '<<TILO:.project_name>>' };
});

// Run the server!
const start = async () => {
  try {
    const port = Number(process.env.PORT) || 3000;
    await fastify.listen({ port, host: '0.0.0.0' });
    console.log('🚀 <<TILO:.project_name>> is running on http://localhost:' + port);
  } catch (err) {
    fastify.log.error(err);
    process.exit(1);
  }
};

start();`

	FastifyTsConfig = `{
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

	FastifyEnv = `PORT=3000
NODE_ENV=development`
)
