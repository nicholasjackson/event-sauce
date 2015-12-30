require 'cucumber/rest_api'
require_relative 'cucumber_rest_monkey_patch'
require 'cucumber/pickle_mongodb'
require 'cucumber/mailcatcher'

$SERVER_PATH = ENV['WEB_SERVER_URI']
$REDIS_IP = ENV['REDIS_IP']
$REDIS_PORT = ENV['REDIS_PORT']
$MIMIC_SERVER = ENV['MIMIC_SERVER']

Mongoid.load!('features/support/localdb.yml', :development)
