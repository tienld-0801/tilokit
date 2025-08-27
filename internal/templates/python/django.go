package python

const (
	// Django Manage Template
	DjangoManagePy = `#!/usr/bin/env python
import os
import sys

if __name__ == '__main__':
    os.environ.setdefault('DJANGO_SETTINGS_MODULE', '<<TILO:.project_name>>.settings.development')
    try:
        from django.core.management import execute_from_command_line
    except ImportError as exc:
        raise ImportError("Couldn't import Django.") from exc
    execute_from_command_line(sys.argv)
`

	// Django URLs Template
	DjangoUrlsPy = `from django.contrib import admin
from django.urls import path, include

urlpatterns = [
    path('admin/', admin.site.urls),
    path('', include('<<TILO:.project_name>>.apps.core.urls')),
]
`

	// Django Base Settings Template
	DjangoBasePy = `import os
from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent.parent.parent

SECRET_KEY = os.environ.get('SECRET_KEY', 'django-insecure-change-me')

DJANGO_APPS = [
    'django.contrib.admin',
    'django.contrib.auth',
    'django.contrib.contenttypes',
    'django.contrib.sessions',
    'django.contrib.messages',
    'django.contrib.staticfiles',
]

THIRD_PARTY_APPS = [
    'rest_framework',
    'corsheaders',
]

LOCAL_APPS = [
    '<<TILO:.project_name>>.apps.core',
]

INSTALLED_APPS = DJANGO_APPS + THIRD_PARTY_APPS + LOCAL_APPS

MIDDLEWARE = [
    'corsheaders.middleware.CorsMiddleware',
    'django.middleware.security.SecurityMiddleware',
    'django.contrib.sessions.middleware.SessionMiddleware',
    'django.middleware.common.CommonMiddleware',
    'django.middleware.csrf.CsrfViewMiddleware',
    'django.contrib.auth.middleware.AuthenticationMiddleware',
    'django.contrib.messages.middleware.MessageMiddleware',
    'django.middleware.clickjacking.XFrameOptionsMiddleware',
]

ROOT_URLCONF = '<<TILO:.project_name>>.urls'

TEMPLATES = [
    {
        'BACKEND': 'django.template.backends.django.DjangoTemplates',
        'DIRS': [BASE_DIR / 'templates'],
        'APP_DIRS': True,
        'OPTIONS': {
            'context_processors': [
                'django.template.context_processors.debug',
                'django.template.context_processors.request',
                'django.contrib.auth.context_processors.auth',
                'django.contrib.messages.context_processors.messages',
            ],
        },
    },
]

DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.postgresql',
        'NAME': os.environ.get('DB_NAME', '<<TILO:.project_name>>'),
        'USER': os.environ.get('DB_USER', 'postgres'),
        'PASSWORD': os.environ.get('DB_PASSWORD', 'password'),
        'HOST': os.environ.get('DB_HOST', 'localhost'),
        'PORT': os.environ.get('DB_PORT', '5432'),
    }
}

LANGUAGE_CODE = 'en-us'
TIME_ZONE = 'UTC'
USE_I18N = True
USE_TZ = True

STATIC_URL = '/static/'
STATIC_ROOT = BASE_DIR / 'staticfiles'
STATICFILES_DIRS = [BASE_DIR / 'static']

MEDIA_URL = '/media/'
MEDIA_ROOT = BASE_DIR / 'media'

DEFAULT_AUTO_FIELD = 'django.db.models.BigAutoField'

REST_FRAMEWORK = {
    'DEFAULT_AUTHENTICATION_CLASSES': [
        'rest_framework.authentication.SessionAuthentication',
    ],
    'DEFAULT_PERMISSION_CLASSES': [
        'rest_framework.permissions.IsAuthenticated',
    ],
    'PAGE_SIZE': 20,
}

CORS_ALLOWED_ORIGINS = [
    "http://localhost:3000",
    "http://127.0.0.1:3000",
]
`

	// Django Development Settings Template
	DjangoDevelopmentPy = `from .base import *

DEBUG = True
ALLOWED_HOSTS = ['localhost', '127.0.0.1', '0.0.0.0']

DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.sqlite3',
        'NAME': BASE_DIR / 'db.sqlite3',
    }
}

EMAIL_BACKEND = 'django.core.mail.backends.console.EmailBackend'
`

	// Django Core App Config Template
	DjangoCoreAppsPy = `from django.apps import AppConfig

class CoreConfig(AppConfig):
    default_auto_field = 'django.db.models.BigAutoField'
    name = '<<TILO:.project_name>>.apps.core'
`

	// Django Core Models Template
	DjangoCoreModelsPy = `from django.db import models

class Item(models.Model):
    title = models.CharField(max_length=200)
    description = models.TextField(blank=True)
    created = models.DateTimeField(auto_now_add=True)

    def __str__(self):
        return self.title
`

	// Django Core Views Template
	DjangoCoreViewsPy = `from django.shortcuts import render
from rest_framework import viewsets
from rest_framework.decorators import api_view
from rest_framework.response import Response
from .models import Item

def home(request):
    return render(request, 'home.html', {'title': '<<TILO:.project_name>>'})

@api_view(['GET'])
def health_check(request):
    return Response({'status': 'healthy', 'service': '<<TILO:.project_name>>'})
`

	// Django Core URLs Template
	DjangoCoreUrlsPy = `from django.urls import path
from . import views

app_name = 'core'

urlpatterns = [
    path('', views.home, name='home'),
    path('health/', views.health_check, name='health_check'),
]
`

	// Django Requirements Template
	DjangoRequirementsTxt = `Django==<<TILO:.django_version>>
djangorestframework==3.14.0
django-cors-headers==4.3.1
psycopg2-binary==2.9.9
`

	// Django Poetry Template
	DjangoPyprojectToml = `[tool.poetry]
name = "<<TILO:.project_name>>"
version = "0.1.0"
description = "Django web application"

[tool.poetry.dependencies]
python = "^<<TILO:.python_version>>"
django = "^<<TILO:.django_version>>"
djangorestframework = "^3.14.0"
django-cors-headers = "^4.3.1"
psycopg2-binary = "^2.9.9"

[tool.poetry.group.dev.dependencies]
pytest = "^7.4.0"
pytest-django = "^4.7.0"

[build-system]
requires = ["poetry-core"]
build-backend = "poetry.core.masonry.api"
`

	// Django Home Template
	DjangoHomeHtml = `<!DOCTYPE html>
<html>
<head>
    <title><<TILO:.project_name>></title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <div class="container mt-5">
        <h1>Welcome to <<TILO:.project_name>>!</h1>
        <p>Your Django application is running.</p>
        <a href="/health/" class="btn btn-primary">Check Health</a>
        <a href="/admin/" class="btn btn-secondary">Admin</a>
    </div>
</body>
</html>
`

	// Django Pipenv Template
	DjangoPipfile = `[[source]]
url = "https://pypi.org/simple"
verify_ssl = true
name = "pypi"

[packages]
django = "<<TILO:.django_version>>"
djangorestframework = "3.14.0"
django-cors-headers = "4.3.1"
psycopg2-binary = "2.9.9"

[dev-packages]
pytest = "7.4.0"
pytest-django = "4.7.0"

[requires]
python_version = "<<TILO:.python_version>>"
`

	// Django Conda Environment Template
	DjangoEnvironmentYml = `name: <<TILO:.project_name>>
channels:
  - conda-forge
  - defaults
dependencies:
  - python=<<TILO:.python_version>>
  - django=<<TILO:.django_version>>
  - pip
  - pip:
    - djangorestframework==3.14.0
    - django-cors-headers==4.3.1
    - psycopg2-binary==2.9.9
    - pytest==7.4.0
    - pytest-django==4.7.0
`

	// Django README Template
	DjangoReadmeMd = `# <<TILO:.project_name>>

A Django web application built with TiLoKit.

## Setup

1. Create virtual environment:
` + "```bash" + `
python -m venv venv
source venv/bin/activate
` + "```" + `

2. Install dependencies:
` + "```bash" + `
pip install -r requirements.txt
` + "```" + `

3. Run migrations:
` + "```bash" + `
python manage.py migrate
` + "```" + `

4. Create superuser:
` + "```bash" + `
python manage.py createsuperuser
` + "```" + `

5. Run the server:
` + "```bash" + `
python manage.py runserver
` + "```" + `

Visit http://localhost:8000 to see your application.
`
)
