class Event
  include Mongoid::Document
  field :event_name, type: String
  field :payload,    type: String
end

class EventStoreItem
  include Mongoid::Document
  store_in collection: "events"

  field :_id,           type: Moped::BSON::ObjectId
  field :event,         type: Event
  field :creation_date, type: DateTime
end

class DeadLetterItem
  include Mongoid::Document
  store_in collection: "dead_letters"

  field :_id,                type: Moped::BSON::ObjectId
  field :event,              type: Event
  field :creation_date,      type: DateTime
  field :first_failure_date, type: DateTime
  field :next_retry_date,    type: DateTime
  field :failure_count,      type: Integer
  field :callback_url,       type: String
end

FactoryGirl.define do
  factory :event, class: Event do
    event_name "users.new_email"
    callback_url "http://myserver.com/v1/newemail"
  end

  factory :eventstoreitem, class: EventStoreItem do
    event
    creation_date Time.now
  end

  factory :deadletteritem, class: DeadLetterItem do
    event
  end
end
