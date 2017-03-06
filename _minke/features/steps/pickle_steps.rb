require 'cucumber/pickle_mongodb/pickle_steps.rb'

Then(/^(\d+) #{capture_plural_factory} should exist with event name "(.*?)"$/) do |items, plural_factory, event_name|
  count = 0
  models = find_models(plural_factory.singularize, nil)
  models.each do |model|
    count = count + 1 if model.event.event_name == event_name
  end

  count.should == items.to_i
end
