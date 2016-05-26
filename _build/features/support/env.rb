require 'cucumber/rest_api'
require_relative 'cucumber_rest_monkey_patch'
require 'cucumber/pickle_mongodb'
require 'cucumber/mailcatcher'
require 'minke'

discovery = Minke::Docker::ServiceDiscovery.new 'config.yml'
$SERVER_PATH = "http://#{discovery.public_address_for 'sorcery', '8001', :cucumber}"

$REDIS_IP = ENV['DOCKER_IP']
$REDIS_PORT = 16379
$MIMIC_SERVER = "http://#{ENV['DOCKER_IP']}:11988"

Mongoid.load!('features/support/localdb.yml', :development)
