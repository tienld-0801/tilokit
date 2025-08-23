package php

const (
	CakePHPComposer = `{
    "name": "<<TILO:.project_name>>",
    "description": "CakePHP skeleton app",
    "homepage": "https://cakephp.org",
    "type": "project",
    "license": "MIT",
    "require": {
        "php": ">=8.1",
        "cakephp/cakephp": "^5.0"
    },
    "require-dev": {
        "cakephp/bake": "^3.4.0",
        "cakephp/cakephp-codesniffer": "^5.0",
        "cakephp/debug_kit": "^5.1.3",
        "phpunit/phpunit": "^10.1.0"
    },
    "suggest": {
        "markstory/asset_compress": "An asset compression plugin which provides file concatenation and a flexible filter system for preprocessing and minification.",
        "dereuromark/cakephp-ide-helper": "After baking your code, this keeps your annotations in sync with the code evolving from there on for maximum IDE and PHPStan/Psalm compatibility."
    },
    "autoload": {
        "psr-4": {
            "App\\": "src/"
        }
    },
    "autoload-dev": {
        "psr-4": {
            "App\\Test\\": "tests/"
        }
    }
}`

	CakePHPController = `<?php
declare(strict_types=1);

namespace App\Controller;

use Cake\Controller\Controller;

class HomeController extends AppController
{
    public function index()
    {
        $data = [
            'message' => 'Hello World from "<<TILO:.project_name>>"!',
            'framework' => 'CakePHP'
        ];
        return $this->response
            ->withType('application/json')
            ->withStringBody(json_encode($data, JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE));
    }

    public function health()
    {
        $data = ['status' => 'ok'];
        return $this->response
            ->withType('application/json')
            ->withStringBody(json_encode($data, JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE));
    }
}`

	CakePHPRoutes = `<?php
use Cake\Routing\RouteBuilder;
use Cake\Routing\Route\DashedRoute;

return static function (RouteBuilder $routes) {
    $routes->setRouteClass(DashedRoute::class);

    $routes->scope('/', function (RouteBuilder $builder) {
        $builder->connect('/', ['controller' => 'Home', 'action' => 'index']);
        $builder->connect('/health', ['controller' => 'Home', 'action' => 'health']);

        $builder->fallbacks(DashedRoute::class);
    });
};`

	// #nosec G101 - This is a template string, not hardcoded credentials
	CakePHPWebroot = `<?php

use App\Application;
use Cake\Http\Server;

require dirname(__DIR__) . '/vendor/autoload.php';

$server = new Server(new Application(dirname(__DIR__) . '/config'));
$server->emit($server->run());`

	CakePHPBootstrap = `<?php
declare(strict_types=1);

use Cake\Cache\Cache;
use Cake\Core\Configure;
use Cake\Core\Configure\Engine\PhpConfig;
use Cake\Datasource\ConnectionManager;
use Cake\Error\ErrorTrap;
use Cake\Error\ExceptionTrap;
use Cake\Http\Exception\NotFoundException;
use Cake\Log\Log;
use Cake\Utility\Security;

if (!defined('DS')) {
    define('DS', DIRECTORY_SEPARATOR);
}

define('ROOT', dirname(__DIR__));
define('APP_DIR', 'src');
define('APP', ROOT . DS . APP_DIR . DS);
define('CONFIG', ROOT . DS . 'config' . DS);
define('WWW_ROOT', ROOT . DS . 'webroot' . DS);
define('TESTS', ROOT . DS . 'tests' . DS);
define('TMP', ROOT . DS . 'tmp' . DS);
define('LOGS', ROOT . DS . 'logs' . DS);
define('CACHE', TMP . 'cache' . DS);
define('RESOURCES', ROOT . DS . 'resources' . DS);
define('CAKE_CORE_INCLUDE_PATH', ROOT . DS . 'vendor' . DS . 'cakephp' . DS . 'cakephp');
define('CORE_PATH', CAKE_CORE_INCLUDE_PATH . DS);
define('CAKE', CORE_PATH . 'src' . DS);

Configure::config('default', new PhpConfig());
Configure::load('app', 'default', false);`

	CakePHPBincake = `#!/usr/bin/env php
<?php
declare(strict_types=1);

use App\Application;
use Cake\Console\CommandRunner;

require dirname(__DIR__) . '/vendor/autoload.php';

exit((new CommandRunner(new Application(dirname(__DIR__) . '/config'), 'cake'))->run($argv));`

	CakePHPAppConfig = `<?php
return [
    'debug' => filter_var(env('DEBUG', true), FILTER_VALIDATE_BOOLEAN),
    'App' => [
        'namespace' => 'App',
        'encoding' => env('APP_ENCODING', 'UTF-8'),
        'defaultLocale' => env('APP_DEFAULT_LOCALE', 'en_US'),
        'defaultTimezone' => env('APP_DEFAULT_TIMEZONE', 'UTC'),
        'base' => false,
        'dir' => 'src',
        'webroot' => 'webroot',
        'wwwRoot' => WWW_ROOT,
        'fullBaseUrl' => false,
        'imageBaseUrl' => 'img/',
        'cssBaseUrl' => 'css/',
        'jsBaseUrl' => 'js/',
        'paths' => [
            'plugins' => [ROOT . DS . 'plugins' . DS],
            'templates' => [ROOT . DS . 'templates' . DS],
            'locales' => [RESOURCES . 'locales' . DS],
        ],
    ],
];`

	CakePHPEnvExample = `#!/usr/bin/env php
# When you add to this file, also add to config/app_local.example.php

export APP_NAME="<<TILO:.project_name>>"
export DEBUG="true"
export APP_ENCODING="UTF-8"
export APP_DEFAULT_LOCALE="en_US"
export APP_DEFAULT_TIMEZONE="UTC"
export SECURITY_SALT=""

export DATABASE_URL="sqlite:///tmp/data.db"

export EMAIL_TRANSPORT_DEFAULT_URL=""

export CACHE_CACHECORE_DEFAULT_URL="File://tmp/cache/persistent/"
export CACHE_CAKEMODEL_DEFAULT_URL="File://tmp/cache/models/"
export CACHE_CAKECORE_DEFAULT_URL="File://tmp/cache/persistent/"

export LOG_DEBUG_URL="file://logs/"
export LOG_ERROR_URL="file://logs/"`

	CakePHPAppController = `<?php
declare(strict_types=1);

namespace App\Controller;

use Cake\Controller\Controller;

class AppController extends Controller
{
    public function initialize(): void
    {
        parent::initialize();

        $this->loadComponent('RequestHandler');
        $this->loadComponent('Flash');
    }
}`

	CakePHPEnv = `# Database Configuration
export DATABASE_URL="mysql://root:password@localhost/<<TILO:.project_name>>?encoding=utf8&timezone=UTC&cacheMetadata=true&quoteIdentifiers=false&persistent=false"

# Security salt
export SECURITY_SALT="<<TILO:.project_name>>_security_salt_change_this"

# Debug mode
export DEBUG="true"

# App configuration
export APP_NAME="<<TILO:.project_name>>"
export APP_ENCODING="UTF-8"
export APP_DEFAULT_LOCALE="en_US"
export APP_DEFAULT_TIMEZONE="UTC"`
)
