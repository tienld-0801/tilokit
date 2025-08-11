package php

const (
	LaravelComposer = `{
    "name": "{{.ProjectName}}",
    "type": "project",
    "description": "Laravel application",
    "keywords": ["framework", "laravel"],
    "license": "MIT",
    "require": {
        "php": "^8.1",
        "laravel/framework": "^10.0",
        "laravel/sanctum": "^3.2",
        "laravel/tinker": "^2.8"
    },
    "require-dev": {
        "fakerphp/faker": "^1.9.1",
        "laravel/pint": "^1.0",
        "laravel/sail": "^1.18",
        "mockery/mockery": "^1.4.4",
        "nunomaduro/collision": "^7.0",
        "phpunit/phpunit": "^10.1",
        "spatie/laravel-ignition": "^2.0"
    },
    "autoload": {
        "psr-4": {
            "App\\": "app/",
            "Database\\Factories\\": "database/factories/",
            "Database\\Seeders\\": "database/seeders/"
        }
    }
}`

	LaravelController = `<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;

class HomeController extends Controller
{
    public function index()
    {
        return response()->json([
            'message' => 'Hello World from {{.ProjectName}}!'
        ]);
    }

    public function health()
    {
        return response()->json([
            'status' => 'ok'
        ]);
    }
}`

	LaravelRoutes = `<?php

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;
use App\Http\Controllers\HomeController;

Route::get('/', [HomeController::class, 'index']);
Route::get('/health', [HomeController::class, 'health']);`
)
