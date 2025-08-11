package python

const (
	DjangoRequirements = `Django==4.2.0
djangorestframework==3.14.0
python-decouple==3.8
gunicorn==21.2.0`

	DjangoSettings = `"""
Django settings for {{.ProjectName}} project.
"""

from pathlib import Path
from decouple import config

BASE_DIR = Path(__file__).resolve().parent.parent

SECRET_KEY = config('SECRET_KEY', default='your-secret-key-here')

DEBUG = config('DEBUG', default=True, cast=bool)

ALLOWED_HOSTS = []

INSTALLED_APPS = [
    'django.contrib.admin',
    'django.contrib.auth',
    'django.contrib.contenttypes',
    'django.contrib.sessions',
    'django.contrib.messages',
    'django.contrib.staticfiles',
    'rest_framework',
]

MIDDLEWARE = [
    'django.middleware.security.SecurityMiddleware',
    'django.contrib.sessions.middleware.SessionMiddleware',
    'django.middleware.common.CommonMiddleware',
    'django.middleware.csrf.CsrfViewMiddleware',
    'django.contrib.auth.middleware.AuthenticationMiddleware',
    'django.contrib.messages.middleware.MessageMiddleware',
    'django.middleware.clickjacking.XFrameOptionsMiddleware',
]

ROOT_URLCONF = '{{.ProjectName}}.urls'

DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.sqlite3',
        'NAME': BASE_DIR / 'db.sqlite3',
    }
}

STATIC_URL = '/static/'
`

	DjangoUrls = `from django.contrib import admin
from django.urls import path
from django.http import JsonResponse

def hello_world(request):
    return JsonResponse({'message': 'Hello World from {{.ProjectName}}!'})

urlpatterns = [
    path('admin/', admin.site.urls),
    path('', hello_world, name='hello'),
]`
)
