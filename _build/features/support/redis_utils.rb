def clear_queue
  redis = Redis.new(:host => $REDIS_IP, :port => $REDIS_PORT.to_i, :db => 1)
  redis.ltrim('rmq::queue::[message_queue]::ready', 1, 0)
end
