@dispatch_event
Feature: Dispatch Event
	In order to ensure quality
	As a user
	I want to be able to test functionality of my API

Scenario: Dispatch a message to a single healthy consumers, should result in one message being received and the event stored in the eventstore
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
	And I send a POST request to "/v1/register" with the following:
		| message_name | mytest.event                              |
		| callback_url | http://callbackserver:11988/v1/helloworld |
	And the response status should be "200"
	When I send a POST request to "/v1/event" with the following:
		"""
			{
				"message_name": "mytest.event",
				"payload": {
					"something": "something"
				}
			}
		"""
	Then I expect 1 callbacks to have been received with the correct payload
	And 1 events should exist with message_name: "mytest.event"

Scenario: Dispatch a message to a multiple healthy consumers, should result in two messages being received.
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
	And I send a POST request to "/v1/register" with the following:
		| message_name | mytest.event                              |
		| callback_url | http://callbackserver:11988/v1/helloworld |
	And I send a POST request to "/v1/register" with the following:
		| message_name | mytest.event                                 |
		| callback_url | http://callbackserver:11988/v1/othercallback |
	And the response status should be "200"
	When I send a POST request to "/v1/event" with the following:
		"""
			{
				"message_name": "mytest.event",
				"payload": {
					"something": "something"
				}
			}
		"""
	Then I expect 2 callbacks to have been received with the correct payload

Scenario: Dispatch a message to a one healthy one nonexistent consumers, should result in a message received by the healthy endpoint and the unhealthy endpoint being unregistered
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
	And I send a POST request to "/v1/register" with the following:
		| message_name | mytest.event                              |
		| callback_url | http://callbackserver:11988/v1/helloworld |
	And I send a POST request to "/v1/register" with the following:
		| message_name | mytest.event                                |
		| callback_url | http://callbackserver:11988/v1/doesnotexist |
	And the response status should be "200"
	When I send a POST request to "/v1/event" with the following:
		"""
			{
				"message_name": "mytest.event",
				"payload": {
					"something": "something"
				}
			}
		"""
	Then I expect 1 callbacks to have been received with the correct payload
	And 1 registrations should exist with message_name: "mytest.event", callback_url: "http://callbackserver:11988/v1/helloworld"
	And 0 registrations should exist with message_name: "mytest.event", callback_url: "http://callbackserver:11988/v1/doesnotexist"

Scenario: Dispatch an event to a one unhealthy consumers, should result in an event added to the dead letter queue
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
	And I send a POST request to "/v1/register" with the following:
		| message_name | mytest.event                        |
		| callback_url | http://callbackservers/v1/unhealthy |
	And the response status should be "200"
	When I send a POST request to "/v1/event" with the following:
		"""
			{
				"message_name": "mytest.event",
				"payload": {
					"something": "something"
				}
			}
		"""
	Then I expect 1 event on the dead letter queue
	And 1 registrations should exist with message_name: "mytest.event", callback_url: "http://callbackservers/v1/unhealthy"
