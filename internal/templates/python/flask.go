package python

const (
	FlaskRequirements = `Flask==2.3.0
Flask-CORS==4.0.0
python-decouple==3.8
gunicorn==21.2.0`

	FlaskApp = `from flask import Flask, jsonify
from flask_cors import CORS
from decouple import config

app = Flask(__name__)
CORS(app)

@app.route('/')
def hello_world():
    return jsonify({'message': 'Hello World from {{.ProjectName}}!'})

@app.route('/health')
def health():
    return jsonify({'status': 'ok'})

if __name__ == '__main__':
    debug_mode = config('DEBUG', default=True, cast=bool)
    app.run(debug=debug_mode, host='0.0.0.0', port=5000)`

	FlaskEnv = `DEBUG=True
SECRET_KEY=your-secret-key-here`
)
