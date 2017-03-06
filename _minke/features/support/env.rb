require 'cucumber/rest_api'
require_relative 'cucumber_rest_monkey_patch'
require 'cucumber/pickle_mongodb'
require 'cucumber/mailcatcher'
require 'minke'

discovery = Minke::Docker::ServiceDiscovery.new(
  ENV['DOCKER_PROJECT'],
  Minke::Docker::DockerRunner.new(Minke::Logging.create_logger(STDOUT,true)),
  ENV['DOCKER_NETWORK']
)
$SERVER_PATH = "http://#{discovery.bridge_address_for 'sorcery', '8001'}"

$REDIS_SERVER = "#{discovery.bridge_address_for 'redis', '6379'}"
puts $REDIS_SERVER
$REDIS_IP = $REDIS_SERVER.split(":")[0]
$REDIS_PORT = $REDIS_SERVER.split(":")[1]

$MIMIC_SERVER = "http://#{discovery.bridge_address_for 'callbackserver', '11988'}"

$MONGO_SERVER = "#{discovery.bridge_address_for 'mongo', '27017'}"

Mongoid.load!('features/support/localdb.yml', :development)
