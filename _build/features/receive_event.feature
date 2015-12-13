@receive_event
Feature: Receive Event
	In order to ensure quality
	As a user
	I want to be able to test functionality of my API

Scenario: Receive a valid event
	Given I send a POST request to "/v1/event"
	Then the response status should be "200"
	And the JSON response should have "$..status_message" with the text "OK"
  And a new message should be registered on the queue

Scenario: Receive a event with no payload
	Given I send a POST request to "/v1/event"
	Then the response status should be "500"
	And the JSON response should have "$..status_message" with the text "No payload"
  And a new message should not be registered on the queue

Scenario: Receive a event with no title
	Given I send a POST request to "/v1/event"
	Then the response status should be "200"
	And the JSON response should have "$..status_message" with the text "No title"
  And a new message should not be registered on the queue
