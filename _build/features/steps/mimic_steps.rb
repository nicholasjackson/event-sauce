require 'rest_client'

Given(/^Mimic is configured with specification$/) do |specification|
  RestClient.send("post", "#{$MIMIC_SERVER}/api/multi", "#{specification}")
end
