require 'rest_client'
require 'json'

Given(/^Mimic is configured with specification$/) do |specification|
  clearMimicStubs
  RestClient.send("post", "#{$MIMIC_SERVER}/api/multi", "#{specification}")
end

Then(/^I expect (\d+) callbacks to have been received with the correct payload$/) do |count|
  timer = 0
  while true
    requests = checkResponse

    raise "expected #{count} callbacks received #{requests.length}" unless timer < 5

    if requests.length != count.to_i
      timer = timer + 1
      sleep 1
    else
      break
    end
  end
end

def checkResponse
  response = RestClient.send("get", "#{$MIMIC_SERVER}/api/requests")
  body = JSON.parse(response.body)
  return body["requests"]
end

def clearMimicStubs
  response = RestClient.send("post", "#{$MIMIC_SERVER}/api/clear", "")
end
