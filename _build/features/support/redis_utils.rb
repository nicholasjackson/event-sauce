def clear_queues
  redis = Redis.new(:host => $REDIS_IP, :port => $REDIS_PORT.to_i, :db => 1)
  redis.ltrim('rmq::queue::[event_queue]::ready', 1, 0)
  redis.ltrim('rmq::queue::[dead_letter_queue]::ready', 1, 0)
end
