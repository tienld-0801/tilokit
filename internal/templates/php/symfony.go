package php

const (
	SymfonyComposer = `{
    "name": "<<TILO:.project_name>>",
    "type": "project",
    "description": "Symfony application",
    "license": "MIT",
    "minimum-stability": "stable",
    "prefer-stable": true,
    "require": {
        "php": ">=8.1",
        "ext-ctype": "*",
        "ext-iconv": "*",
        "symfony/console": "6.4.*",
        "symfony/dotenv": "6.4.*",
        "symfony/flex": "^2",
        "symfony/framework-bundle": "6.4.*",
        "symfony/runtime": "6.4.*",
        "symfony/yaml": "6.4.*"
    },
    "require-dev": {
        "symfony/maker-bundle": "^1.0"
    },
    "config": {
        "allow-plugins": {
            "php-http/discovery": true,
            "symfony/flex": true,
            "symfony/runtime": true
        },
        "sort-packages": true
    },
    "autoload": {
        "psr-4": {
            "App\\": "src/"
        }
    },
    "autoload-dev": {
        "psr-4": {
            "App\\Tests\\": "tests/"
        }
    }
}`

	SymfonyController = `<?php

namespace App\Controller;

use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\Routing\Annotation\Route;

class HomeController extends AbstractController
{
    #[Route('/', name: 'app_home')]
    public function index(): JsonResponse
    {
        return $this->json([
            'message' => 'Hello World from "<<TILO:.project_name>>"!',
            'framework' => 'Symfony'
        ]);
    }

    #[Route('/health', name: 'app_health')]
    public function health(): JsonResponse
    {
        return $this->json([
            'status' => 'ok'
        ]);
    }
}`

	SymfonyKernel = `<?php

namespace App;

use Symfony\Bundle\FrameworkBundle\Kernel\MicroKernelTrait;
use Symfony\Component\HttpKernel\Kernel as BaseKernel;

class Kernel extends BaseKernel
{
    use MicroKernelTrait;
}`

	SymfonyEnv = `APP_ENV=dev
APP_SECRET=<<TILO:.app_secret>>
DATABASE_URL="sqlite:///%kernel.project_dir%/var/data.db"`

	SymfonyPublicIndex = `<?php

use App\Kernel;

require_once dirname(__DIR__).'/vendor/autoload_runtime.php';

return function (array $context) {
    return new Kernel($context['APP_ENV'], (bool) $context['APP_DEBUG']);
};`
)
