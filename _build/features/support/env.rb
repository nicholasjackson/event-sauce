require 'cucumber/rest_api'
require_relative 'cucumber_rest_monkey_patch'
require 'cucumber/pickle_mongodb'
require 'cucumber/mailcatcher'

$SERVER_PATH = "http://#{ENV['DOCKER_IP']}:8001"
$REDIS_IP = ENV['DOCKER_IP']
$REDIS_PORT = 16379
$MIMIC_SERVER = "http://#{ENV['DOCKER_IP']}:11988"

Mongoid.load!('features/support/localdb.yml', :development)
