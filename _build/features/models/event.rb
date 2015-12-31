class Event
  include Mongoid::Document
  field :event_name, type: String
  field :payload,    type: String

  embedded_in :deadletteritem
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

  embeds_one :event
end

FactoryGirl.define do
  factory :event, class: Event do
    event_name "users.new_email"
  end

  factory :eventstoreitem, class: EventStoreItem do
    event
    creation_date Time.now
  end

  factory :deadletteritem, class: DeadLetterItem do
    association :event, factory: :event, strategy: :build
    failure_count 1
    callback_url "http://myserver.com/v1/newemail"
    first_failure_date Time.now
  end

end
