package ruby

const (
	RailsGemfile = `source 'https://rubygems.org'
git_source(:github) { |repo| "https://github.com/#{repo}.git" }

ruby '3.2.0'

gem 'rails', '~> 7.0.0'
gem 'sqlite3', '~> 1.4'
gem 'puma', '~> 5.0'
gem 'bootsnap', '>= 1.4.4', require: false
gem 'rack-cors'

group :development, :test do
  gem 'byebug', platforms: [:mri, :mingw, :x64_mingw]
end

group :development do
  gem 'listen', '~> 3.3'
  gem 'spring'
end`

	RailsController = `class ApplicationController < ActionController::API
  def index
    render json: { message: 'Hello World from {{.ProjectName}}!' }
  end

  def health
    render json: { status: 'ok' }
  end
end`

	RailsRoutes = `Rails.application.routes.draw do
  root 'application#index'
  get '/health', to: 'application#health'
end`
)
