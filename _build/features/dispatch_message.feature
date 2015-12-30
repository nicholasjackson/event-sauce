@dispatch_message
Feature: Dispatch Message
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

Scenario: Dispatch a message to a single healthy consumers, should result in one message being revceived.
	Given I send a POST request to "/v1/register" with the following:
		| message_name | mytest.event                              |
		| callback_url | http://192.168.99.100:11988/v1/helloworld |
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

Scenario: Dispatch a message to a multiple healthy consumers, should result in two messages being received.
	Given I send a POST request to "/v1/register" with the following:
		| message_name | mytest.event                              |
		| callback_url | http://192.168.99.100:11988/v1/helloworld |
	And I send a POST request to "/v1/register" with the following:
		| message_name | mytest.event                                 |
		| callback_url | http://192.168.99.100:11988/v1/othercallback |
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

Scenario: Dispatch a message to a one healthy one unhealthy consumers, should result in a message received by the healthy endpoint and the unhealthy endpoint being unregistered
	Given I send a POST request to "/v1/register" with the following:
		| message_name | mytest.event                              |
		| callback_url | http://192.168.99.100:11988/v1/helloworld |
	And I send a POST request to "/v1/register" with the following:
		| message_name | mytest.event                                  |
		| callback_url | http://badserver1231sdsd.com/v1/doesnotexist  |
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
	And 1 registrations should exist with message_name: "testmessage.register", callback_url: "http://192.168.99.100:11988/v1/helloworld"
	And 0 registrations should exist with message_name: "testmessage.register", callback_url: "http://badserver1231sdsd.com/v1/doesnotexist"
