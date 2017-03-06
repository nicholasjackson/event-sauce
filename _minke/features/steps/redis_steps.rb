require 'redis'

# refactor this and mimic steps
Then(/^(\d+) (?:event|events) should be registered on the queue$/) do |count|
  checkQueueLength 'rmq::queue::[event_queue]::ready', count
end

Then(/^I expect (\d+) event on the dead letter queue$/) do |count|
  checkQueueLength 'rmq::queue::[dead_letter_queue]::ready', count
end

def checkQueueLength queue, count
  timer = 0

  while true
    redis = Redis.new(:host => $REDIS_IP, :port => $REDIS_PORT.to_i, :db => 1)
    elements = redis.lrange(queue, 0, -1)

    qlength = elements.length

    raise "expected #{count} events, received #{qlength}" unless timer < 5

    if qlength != count.to_i
      timer = timer + 1
      sleep 1
    else
      break
    end
  end
end
