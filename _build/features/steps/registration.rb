Then(/^a #{capture_model} should exist with #{capture_fields} "(.*?)"$/) do |model_name, fields|
  u = find_model(model_name, fields)
  puts u
end

Then(/^a registration should not exist$/) do
  pending # express the regexp above with the code you wish you had
end
