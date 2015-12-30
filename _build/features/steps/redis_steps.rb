require 'redis'

# refactor this and mimic steps

Then(/^(\d+) (?:message|messages) should be registered on the queue$/) do |count|
  timer = 0
  while true
    qlength = checkQueueLength 'rmq::queue::[message_queue]::ready'

    raise "expected #{count} messages, received #{requests.length}" unless timer < 5

    if qlength != count.to_i
      timer = timer + 1
      sleep 1
    else
      break
    end
  end
end

Then(/^I expect (\d+) event on the dead letter queue$/) do |count|
  timer = 0

  while true
    qlength = checkQueueLength 'rmq::queue::[dead_letter_queue]::ready'

    raise "expected #{count} messages, received #{requests.length}" unless timer < 5

    if qlength != count.to_i
      timer = timer + 1
      sleep 1
    else
      break
    end
  end
end

def checkQueueLength queue
  redis = Redis.new(:host => $REDIS_IP, :port => $REDIS_PORT.to_i, :db => 1)
  elements = redis.lrange(queue, 0, -1)
  elements.length
end
