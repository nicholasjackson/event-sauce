require 'cucumber/rest_api'
require_relative 'cucumber_rest_monkey_patch'
require 'cucumber/pickle_mongodb'
require 'cucumber/mailcatcher'

$SERVER_PATH = ENV['WEB_SERVER_URI']

Mongoid.load!('features/support/localdb.yml', :development)
