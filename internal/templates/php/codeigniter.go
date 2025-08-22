package php

const (
	CodeIgniterComposer = `{
    "name": "<<TILO:.project_name>>",
    "type": "project",
    "description": "CodeIgniter 4 application",
    "license": "MIT",
    "require": {
        "php": "^7.4 || ^8.0",
        "codeigniter4/framework": "^4.4"
    },
    "require-dev": {
        "fakerphp/faker": "^1.9",
        "mikey179/vfsstream": "^1.6",
        "phpunit/phpunit": "^9.1"
    },
    "autoload": {
        "psr-4": {
            "App\\": "app/",
            "Config\\": "app/Config/"
        },
        "exclude-from-classmap": [
            "**/Database/Migrations/**"
        ]
    }
}`

	CodeIgniterController = `<?php

namespace App\Controllers;

class Home extends BaseController
{
    public function index(): string
    {
        return json_encode([
            'message' => 'Hello World from "<<TILO:.project_name>>"!',
            'framework' => 'CodeIgniter 4'
        ]);
    }

    public function health(): string
    {
        return json_encode([
            'status' => 'ok'
        ]);
    }
}`

	CodeIgniterRoutes = `<?php

use CodeIgniter\Router\RouteCollection;

/**
 * @var RouteCollection $routes
 */
$routes->get('/', 'Home::index');
$routes->get('/health', 'Home::health');`

	CodeIgniterEnv = `#--------------------------------------------------------------------
# Example Environment Configuration file
#--------------------------------------------------------------------

CI_ENVIRONMENT = development

#--------------------------------------------------------------------
# APP
#--------------------------------------------------------------------

app.baseURL = 'http://localhost:8080/'
app.indexPage = ''
app.uriProtocol = 'REQUEST_URI'
app.defaultLocale = 'en'
app.negotiateLocale = false
app.supportedLocales = ['en']

#--------------------------------------------------------------------
# DATABASE
#--------------------------------------------------------------------

database.default.hostname = localhost
database.default.database = ci4
database.default.username = root
database.default.password =
database.default.DBDriver = MySQLi`

	CodeIgniterPublicIndex = `<?php

// Path to the front controller (this file)
define('FCPATH', __DIR__ . DIRECTORY_SEPARATOR);

// Ensure the current directory is pointing to the front controller's directory
chdir(FCPATH);

/*
 *---------------------------------------------------------------
 * BOOTSTRAP THE APPLICATION
 *---------------------------------------------------------------
 * This process sets up the path constants, loads and registers
 * our autoloader, along with Composer's, loads our constants
 * and fires up an environment-specific bootstrapping.
 */

// Load our paths config file
// This is the line that might need to be changed, depending on your folder structure.
require FCPATH . '../app/Config/Paths.php';
// ^^^ Change this line if you move your application folder

$paths = new Config\Paths();

// Location of the framework bootstrap file.
require rtrim($paths->systemDirectory, '\\/ ') . DIRECTORY_SEPARATOR . 'bootstrap.php';

// Load environment settings from .env files into $_SERVER and $_ENV
require_once SYSTEMPATH . 'Config/DotEnv.php';
(new CodeIgniter\Config\DotEnv(ROOTPATH))->load();

/*
 * ---------------------------------------------------------------
 * GRAB OUR CODEIGNITER INSTANCE
 * ---------------------------------------------------------------
 *
 * The CodeIgniter class contains the core functionality to make
 * the application run, and does all of the dirty work for us.
 */
$app = Config\Services::codeigniter();
$app->initialize();
$context = is_cli() ? 'php-cli' : 'web';
$app->setContext($context);

/*
 *---------------------------------------------------------------
 * LAUNCH THE APPLICATION
 *---------------------------------------------------------------
 * Now that everything is setup, it's time to actually fire
 * up the engines and make this app do its thang.
 */
$app->run();`
)
