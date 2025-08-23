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

use CodeIgniter\HTTP\ResponseInterface;

class Home extends BaseController
{
    public function index(): ResponseInterface
    {
        return $this->response->setJSON([
            'message' => 'Hello World from "<<TILO:.project_name>>"!',
            'framework' => 'CodeIgniter 4'
        ]);
    }

    public function health(): ResponseInterface
    {
        return $this->response->setJSON([
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
require rtrim($paths->systemDirectory, '\\/ ') . DIRECTORY_SEPARATOR . 'bootstrap.php';`

	CodeIgniterPaths = `<?php

namespace Config;

class Paths
{
    public $systemDirectory = __DIR__ . '/../../vendor/codeigniter4/framework/system';
    public $appDirectory = __DIR__ . '/..';
    public $writableDirectory = __DIR__ . '/../../writable';
    public $testsDirectory = __DIR__ . '/../../tests';
    public $viewDirectory = __DIR__ . '/../Views';
}`

	CodeIgniterApp = `<?php

namespace Config;

use CodeIgniter\Config\BaseConfig;

class App extends BaseConfig
{
    public $baseURL = 'http://localhost:8080/';
    public $indexPage = '';
    public $uriProtocol = 'REQUEST_URI';
    public $defaultLocale = 'en';
    public $negotiateLocale = false;
    public $supportedLocales = ['en'];
    public $appTimezone = 'UTC';
    public $charset = 'UTF-8';
    public $forceGlobalSecureRequests = false;
    public $sessionDriver = 'CodeIgniter\Session\Handlers\FileHandler';
    public $sessionCookieName = 'ci_session';
    public $sessionExpiration = 7200;
    public $sessionSavePath = WRITEPATH . 'session';
    public $sessionMatchIP = false;
    public $sessionTimeToUpdate = 300;
    public $sessionRegenerateDestroy = false;
    public $cookiePrefix = '';
    public $cookieDomain = '';
    public $cookiePath = '/';
    public $cookieSecure = false;
    public $cookieHTTPOnly = false;
    public $cookieSameSite = 'Lax';
    public $proxyIPs = '';
    public $CSRFTokenName = 'csrf_test_name';
    public $CSRFHeaderName = 'X-CSRF-TOKEN';
    public $CSRFCookieName = 'csrf_cookie_name';
    public $CSRFExpiration = 7200;
    public $CSRFRegenerate = true;
    public $CSRFExcludeURIs = [];
    public $CSPEnabled = false;
}`

	CodeIgniterPublicHtaccess = `RewriteEngine On
RewriteCond %{REQUEST_FILENAME} !-f
RewriteCond %{REQUEST_FILENAME} !-d
RewriteRule ^(.*)$ index.php/$1 [L]`
)
