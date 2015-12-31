@dispatch_event
Feature: Dispatch Event
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

Scenario: Dispatch an event to a single healthy consumers, should result in one event being received and the event stored in the eventstore
	Given I send a POST request to "/v1/register" with the following:
		| event_name | mytest.event                              |
		| callback_url | http://callbackserver:11988/v1/helloworld |
	And the response status should be "200"
	When I send a POST request to "/v1/event" with the following:
		"""
			{
				"event_name": "mytest.event",
				"payload": {
					"something": "something"
				}
			}
		"""
	Then I expect 1 callbacks to have been received with the correct payload
	And 1 eventstoreitems should exist with event name "mytest.event"

Scenario: Dispatch an event to a multiple healthy consumers, should result in two events being received.
	Given I send a POST request to "/v1/register" with the following:
		| event_name | mytest.event                              |
		| callback_url | http://callbackserver:11988/v1/helloworld |
	And I send a POST request to "/v1/register" with the following:
		| event_name | mytest.event                                 |
		| callback_url | http://callbackserver:11988/v1/othercallback |
	And the response status should be "200"
	When I send a POST request to "/v1/event" with the following:
		"""
			{
				"event_name": "mytest.event",
				"payload": {
					"something": "something"
				}
			}
		"""
	Then I expect 2 callbacks to have been received with the correct payload

Scenario: Dispatch an event to a one healthy one nonexistent consumers, should result in a event received by the healthy endpoint and the unhealthy endpoint being unregistered
	Given I send a POST request to "/v1/register" with the following:
		| event_name | mytest.event                              |
		| callback_url | http://callbackserver:11988/v1/helloworld |
	And I send a POST request to "/v1/register" with the following:
		| event_name | mytest.event                                |
		| callback_url | http://callbackserver:11988/v1/doesnotexist |
	And the response status should be "200"
	When I send a POST request to "/v1/event" with the following:
		"""
			{
				"event_name": "mytest.event",
				"payload": {
					"something": "something"
				}
			}
		"""
	Then I expect 1 callbacks to have been received with the correct payload
	And 1 registrations should exist with event_name: "mytest.event", callback_url: "http://callbackserver:11988/v1/helloworld"
	And 0 registrations should exist with event_name: "mytest.event", callback_url: "http://callbackserver:11988/v1/doesnotexist"

Scenario: Dispatch an event to a one unhealthy consumers, should result in an event added to the dead letter queue
	Given I send a POST request to "/v1/register" with the following:
		| event_name | mytest.event                        |
		| callback_url | http://callbackservers/v1/unhealthy |
	And the response status should be "200"
	When I send a POST request to "/v1/event" with the following:
		"""
			{
				"event_name": "mytest.event",
				"payload": {
					"something": "something"
				}
			}
		"""
	Then I expect 1 event on the dead letter queue
	And 1 registrations should exist with event_name: "mytest.event", callback_url: "http://callbackservers/v1/unhealthy"
