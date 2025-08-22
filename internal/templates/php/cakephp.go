package php

const (
	CakePHPComposer = `{
    "name": "<<TILO:.project_name>>",
    "type": "project",
    "description": "CakePHP application",
    "license": "MIT",
    "require": {
        "php": ">=8.1",
        "cakephp/cakephp": "^5.0",
        "cakephp/migrations": "^4.0",
        "cakephp/plugin-installer": "^2.0",
        "mobiledetect/mobiledetectlib": "^4.8"
    },
    "require-dev": {
        "cakephp/bake": "^3.0",
        "cakephp/cakephp-codesniffer": "^5.0",
        "cakephp/debug_kit": "^5.0",
        "phpunit/phpunit": "^10.1.0"
    },
    "autoload": {
        "psr-4": {
            "App\\": "src/"
        }
    },
    "autoload-dev": {
        "psr-4": {
            "App\\Test\\": "tests/",
            "Cake\\Test\\": "vendor/cakephp/cakephp/tests/"
        }
    }
}`

	CakePHPController = `<?php
declare(strict_types=1);

namespace App\Controller;

class HomeController extends AppController
{
    public function index()
    {
        $this->autoRender = false;
        $this->response = $this->response->withType('application/json');

        echo json_encode([
            'message' => 'Hello World from "<<TILO:.project_name>>"!',
            'framework' => 'CakePHP'
        ]);
    }

    public function health()
    {
        $this->autoRender = false;
        $this->response = $this->response->withType('application/json');

        echo json_encode([
            'status' => 'ok'
        ]);
    }
}`

	CakePHPRoutes = `<?php
use Cake\Routing\RouteBuilder;

return static function (RouteBuilder $routes) {
    $routes->setRouteClass('DashedRoute');

    $routes->scope('/', function (RouteBuilder $builder) {
        $builder->connect('/', ['controller' => 'Home', 'action' => 'index']);
        $builder->connect('/health', ['controller' => 'Home', 'action' => 'health']);

        $builder->fallbacks('DashedRoute');
    });
};`

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
