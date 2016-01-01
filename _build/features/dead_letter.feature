@dead_letters
Feature: Dead Letter Queue
	In order to ensure quality
	As a user
	I want to be able to test functionality of my API

Background:
	Given Mimic is configured with specification
		"""
			{"stubs":[
				{
					"method": "POST",
					"path": "/v1/helloworld"
				},
				{
					"method": "POST",
					"path": "/v1/othercallback"
				}
			]}
		"""

Scenario: When an event exists on the dead letter queue and the endpoint is healthy a callback should occur
	Given 1 registrations exist with callback_url: "http://callbackserver:11988/v1/helloworld"
  And 1 deadletteritems exist with callback_url: "http://callbackserver:11988/v1/helloworld"
	And I wait 2 second
	Then I expect 1 callbacks to have been received with the correct payload
