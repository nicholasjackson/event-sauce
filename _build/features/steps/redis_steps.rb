require 'redis'

Then(/^(\d+) (?:message|messages) should be registered on the queue$/) do |count|
  redis = Redis.new(:host => $REDIS_IP, :port => $REDIS_PORT.to_i, :db => 1)
  elements = redis.lrange('rmq::queue::[message_queue]::ready', 0, -1)
  elements.length.should == count.to_i
end
