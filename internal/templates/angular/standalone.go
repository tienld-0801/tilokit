package angular

// Standalone Components Templates (Angular 17+)
const (
	// CSR Standalone Templates
	StandaloneAppComponent = `import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = '<<TILO:.project_name>>';
}`

	StandaloneMainTs = `import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';

bootstrapApplication(AppComponent)
  .catch(err => console.error(err));`

	// SSR Standalone Templates
	StandaloneSSRAppComponent = `import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule],
  template: ` + "`" + `
    <div class="container">
      <h1>Welcome to <<TILO:.project_name>>!</h1>
      <p>Hello World from Angular SSR with Standalone Components!</p>
    </div>
  ` + "`" + `,
  styles: [` + "`" + `
    .container {
      text-align: center;
      margin-top: 100px;
    }
    h1 {
      color: #dd0031;
    }
  ` + "`" + `]
})
export class AppComponent {
  title = '<<TILO:.project_name>>';
}`

	StandaloneSSRMainTs = `import { bootstrapApplication } from '@angular/platform-browser';
import { provideRouter } from '@angular/router';
import { AppComponent } from './app/app.component';

bootstrapApplication(AppComponent, {
  providers: [
    provideRouter([])
  ]
}).catch(err => console.error(err));`

	StandaloneSSRMainServerTs = `import { platformServer } from '@angular/platform-server';

import { AppServerModule } from './app/app.server.module';

const bootstrap = () => platformServer().bootstrapModule(AppServerModule);

export default bootstrap;`

	StandaloneSSRServerTs = `import express from 'express';
import 'zone.js/node';
import { ngExpressEngine } from '@nguniversal/express-engine';
import { APP_BASE_HREF } from '@angular/common';
import { existsSync, readFileSync } from 'fs';
import { join } from 'path';

import bootstrap from './src/main.server';

const app = express();
const PORT = process.env['PORT'] || 4000;
const DIST_FOLDER = join(process.cwd(), 'dist', '<<TILO:.project_name>>', 'browser');
const indexHtml = existsSync(join(DIST_FOLDER, 'index.original.html'))
  ? 'index.original.html'
  : 'index';

app.engine('html', ngExpressEngine({
  bootstrap,
}));

app.set('view engine', 'html');
app.set('views', DIST_FOLDER);

// Serve static files
app.get('*.*', express.static(DIST_FOLDER, {
  maxAge: '1y'
}));

// All regular routes use the Universal engine
app.get('*', (req, res) => {
  res.render(indexHtml, {
    req,
    providers: [{ provide: APP_BASE_HREF, useValue: req.baseUrl }]
  });
});

app.listen(PORT, () => {
  console.log('🚀 Server ready at http://localhost:' + PORT);
});`

	// Angular.json configuration for Standalone SSR
	StandaloneSSRAngularJson = `{
  "$schema": "./node_modules/@angular/cli/lib/config/schema.json",
  "version": 1,
  "newProjectRoot": "projects",
  "projects": {
    "<<TILO:.project_name>>": {
      "projectType": "application",
      "schematics": {},
      "root": "",
      "sourceRoot": "src",
      "prefix": "app",
      "architect": {
        "build": {
          "builder": "@angular-devkit/build-angular:browser",
          "options": {
            "outputPath": "dist/<<TILO:.project_name>>",
            "index": "src/index.html",
            "main": "src/main.ts",
            "polyfills": [
              "zone.js"
            ],
            "tsConfig": "tsconfig.app.json",
            "assets": [
              "src/favicon.ico",
              "src/assets"
            ],
            "styles": [
              "src/styles.css"
            ],
            "scripts": []
          },
          "configurations": {
            "production": {
              "budgets": [
                {
                  "type": "initial",
                  "maximumWarning": "500kb",
                  "maximumError": "1mb"
                },
                {
                  "type": "anyComponentStyle",
                  "maximumWarning": "2kb",
                  "maximumError": "4kb"
                }
              ],
              "outputHashing": "all"
            },
            "development": {
              "buildOptimizer": false,
              "optimization": false,
              "vendorChunk": true,
              "extractLicenses": false,
              "sourceMap": true,
              "namedChunks": true
            }
          },
          "defaultConfiguration": "production"
        },
        "serve": {
          "builder": "@angular-devkit/build-angular:dev-server",
          "configurations": {
            "production": {
              "browserTarget": "<<TILO:.project_name>>:build:production"
            },
            "development": {
              "browserTarget": "<<TILO:.project_name>>:build:development"
            }
          },
          "defaultConfiguration": "development"
        },
        "extract-i18n": {
          "builder": "@angular-devkit/build-angular:extract-i18n",
          "options": {
            "browserTarget": "<<TILO:.project_name>>:build"
          }
        },
        "test": {
          "builder": "@angular-devkit/build-angular:karma",
          "options": {
            "polyfills": [
              "zone.js",
              "zone.js/testing"
            ],
            "tsConfig": "tsconfig.spec.json",
            "assets": [
              "src/favicon.ico",
              "src/assets"
            ],
            "styles": [
              "src/styles.css"
            ],
            "scripts": []
          }
        },
        "server": {
          "builder": "@angular-devkit/build-angular:server",
          "options": {
            "outputPath": "dist/<<TILO:.project_name>>/server",
            "main": "src/main.server.ts",
            "tsConfig": "tsconfig.server.json"
          },
          "configurations": {
            "production": {
              "outputHashing": "media"
            },
            "development": {
              "optimization": false,
              "sourceMap": true,
              "extractLicenses": false
            }
          },
          "defaultConfiguration": "production"
        },
        "serve-ssr": {
          "builder": "@nguniversal/builders:ssr-dev-server",
          "configurations": {
            "development": {
              "browserTarget": "<<TILO:.project_name>>:build:development",
              "serverTarget": "<<TILO:.project_name>>:server:development"
            },
            "production": {
              "browserTarget": "<<TILO:.project_name>>:build:production",
              "serverTarget": "<<TILO:.project_name>>:server:production"
            }
          },
          "defaultConfiguration": "development"
        },
        "prerender": {
          "builder": "@nguniversal/builders:prerender",
          "options": {
            "routes": [
              "/"
            ]
          },
          "configurations": {
            "production": {
              "browserTarget": "<<TILO:.project_name>>:build:production",
              "serverTarget": "<<TILO:.project_name>>:server:production"
            },
            "development": {
              "browserTarget": "<<TILO:.project_name>>:build:development",
              "serverTarget": "<<TILO:.project_name>>:server:development"
            }
          },
          "defaultConfiguration": "production"
        }
      }
    }
  }
}`
)
