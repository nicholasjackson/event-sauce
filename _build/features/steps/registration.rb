Then(/^a #{capture_model} should exist with #{capture_fields} "(.*?)"$/) do |model_name, fields|
  u = find_model(model_name, fields)
  puts u
end

Then(/^a #{capture_model} should not exist$/) do |model_name|
  u = find_model(model_name)
  puts u
end
