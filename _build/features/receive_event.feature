@receive_event
Feature: Receive Event
	In order to ensure quality
	As a user
	I want to be able to test functionality of my API

Scenario: Receive a valid event
	Given I send a POST request to "/v1/event" with the following:
		"""
			{
				"event_name": "mytest.event",
				"payload": {
					"something": "something"
				}
			}
		"""
	Then the response status should be "200"
	And the JSON response should have "$..status_event" with the text "OK"
	And I wait just a second
  And 1 events should exist with event_name: "mytest.event"

Scenario: Receive a event with no payload
	Given I send a POST request to "/v1/event" with the following:
		| event_name | mytest.event |
	Then the response status should be "400"
  And 0 events should be registered on the queue

Scenario: Receive a event with no event_name
	Given I send a POST request to "/v1/event" with the following:
	"""
		{
			"payload": {
				"something": "something"
			}
		}
	"""
	Then the response status should be "400"
  And 0 events should be registered on the queue
