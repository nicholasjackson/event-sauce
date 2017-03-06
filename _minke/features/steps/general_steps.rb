Then(/^I wait (\d+) second$/) do |seconds|
  sleep seconds.to_i
end