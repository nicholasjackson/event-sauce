class Event
  include Mongoid::Document
  field :_id,           type: Moped::BSON::ObjectId
  field :event_name,  type: String
  field :payload,       type: String
  field :creation_date, type: DateTime
end

FactoryGirl.define do
  factory :event, class: Event do
    event_name "users.new_email"
    callback_url "http://myserver.com/v1/newemail"
    creation_date Time.now
  end
end
