class Registration
  include Mongoid::Document
  field :_id, type: Moped::BSON::ObjectId
end

FactoryGirl.define do
  factory :registration, class: Registration do
    message_name "users.new_email"
    health_url  "http://myserver.com/v1/health"
    callback_url "http://myserver.com/v1/newemail"
    creation_date Time.now
  end
end
