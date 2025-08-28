package python

const (
	// Flask App Template
	FlaskAppPy = `from flask import Flask, render_template, jsonify
from config.config import Config

def create_app(config_class=Config):
    app = Flask(__name__, template_folder='app/templates', static_folder='app/static')
    app.config.from_object(config_class)

    @app.route('/')
    def index():
        return render_template('index.html', title='<<TILO:.project_name>>')

    @app.route('/api/health')
    def health_check():
        return jsonify({'status': 'healthy', 'service': '<<TILO:.project_name>>'})

    return app

if __name__ == '__main__':
    app = create_app()
    app.run(debug=True, host='0.0.0.0', port=5000)
`

	// Flask Config Template
	FlaskConfigPy = `import os

class Config:
    SECRET_KEY = os.environ.get('SECRET_KEY') or 'dev-secret-key'
    SQLALCHEMY_DATABASE_URI = os.environ.get('DATABASE_URL') or 'sqlite:///<<TILO:.project_name>>.db'
    SQLALCHEMY_TRACK_MODIFICATIONS = False

class DevelopmentConfig(Config):
    DEBUG = True

class ProductionConfig(Config):
    DEBUG = False

config = {
    'development': DevelopmentConfig,
    'production': ProductionConfig,
    'default': DevelopmentConfig
}
`

	// Flask Requirements Template
	FlaskRequirementsTxt = `Flask==<<TILO:.flask_version>>
Flask-SQLAlchemy==3.1.1
gunicorn==21.2.0
`

	// Flask Poetry Template
	FlaskPyprojectToml = `[tool.poetry]
name = "<<TILO:.project_name>>"
version = "0.1.0"
description = "Flask web application"

[tool.poetry.dependencies]
python = "^<<TILO:.python_version>>"
flask = "^<<TILO:.flask_version>>"
flask-sqlalchemy = "^3.1.1"
gunicorn = "^21.2.0"

[tool.poetry.group.dev.dependencies]
pytest = "^7.4.0"
black = "^23.7.0"

[build-system]
requires = ["poetry-core"]
build-backend = "poetry.core.masonry.api"
`

	// Flask Base HTML Template
	FlaskBaseHtml = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title><<TILO:.project_name>></title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
        <div class="container">
            <a class="navbar-brand" href="/"><<TILO:.project_name>></a>
        </div>
    </nav>

    <main class="container mt-4">
        {% block content %}{% endblock %}
    </main>

    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>
`

	// Flask Index HTML Template
	FlaskIndexHtml = `{% extends "base.html" %}

{% block content %}
<div class="row">
    <div class="col-lg-8 mx-auto">
        <div class="jumbotron bg-primary text-white p-5 rounded">
            <h1 class="display-4">Welcome to <<TILO:.project_name>>!</h1>
            <p class="lead">Your Flask application is up and running.</p>
            <a class="btn btn-light btn-lg" href="/api/health" role="button">Check API Health</a>
        </div>
    </div>
</div>
{% endblock %}
`

	// Flask Pipenv Template
	FlaskPipfile = `[[source]]
url = "https://pypi.org/simple"
verify_ssl = true
name = "pypi"

[packages]
flask = "<<TILO:.flask_version>>"
flask-sqlalchemy = "3.1.1"
gunicorn = "21.2.0"

[dev-packages]
pytest = "7.4.0"
black = "23.7.0"

[requires]
python_version = "<<TILO:.python_version>>"
`

	// Flask Conda Environment Template
	FlaskEnvironmentYml = `name: <<TILO:.project_name>>
channels:
  - conda-forge
  - defaults
dependencies:
  - python=<<TILO:.python_version>>
  - flask=<<TILO:.flask_version>>
  - pip
  - pip:
    - flask-sqlalchemy==3.1.1
    - gunicorn==21.2.0
    - pytest==7.4.0
    - black==23.7.0
`

	// Flask README Template
	FlaskReadmeMd = `# <<TILO:.project_name>>

A Flask web application built with TiLoKit.

## Setup

1. Create virtual environment:
` + "```bash" + `
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
` + "```" + `

2. Install dependencies:
` + "```bash" + `
pip install -r requirements.txt
` + "```" + `

3. Run the application:
` + "```bash" + `
python app.py
` + "```" + `

Visit http://localhost:5000 to see your application.

## Features

- Flask web framework
- Bootstrap UI
- SQLAlchemy ORM ready
- Health check endpoint
- Development and production configs
`
)
