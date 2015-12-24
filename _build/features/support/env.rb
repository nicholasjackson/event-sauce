require 'cucumber/rest_api'
require 'cucumber/pickle_mongodb'
require 'cucumber/mailcatcher'
require 'cucumber/pickle_mongodb/pickle_steps.rb'

$SERVER_PATH = ENV['WEB_SERVER_URI']

Mongoid.load!('features/support/localdb.yml', :development)
