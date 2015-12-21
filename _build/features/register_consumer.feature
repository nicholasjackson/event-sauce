@register_consumer
Feature: Register Consumer
	In order to ensure quality
	As a user
	I want to be able to test functionality of my API

Scenario: Invalid GET request
  Given I send a GET request to "/v1/register"
  Then the response status should be "404"

Scenario: Register endpoint with correct data registers consumer successfully
	Given I send a POST request to "/v1/register" with the following:
	"""
	{
		"message_name": "testmessage.register",
		"health_url": "http://something.something/v1/health",
		"callback_url": "http://something.something/v1/callback"
	}
	"""
	Then the response status should be "200"
	And the JSON response should have "$..status_message" with the text "OK"
  And a registration should exist with message_name: "testmessage.register"

Scenario: Register endpoint with no callback returns error
	Given I send a POST request to "/v1/register" with the following:
	"""
	{
		"message_name": "testmessage.register",
		"health_url": "http://something.something/v1/health"
	}
	"""
	Then the response status should be "500"
	And the JSON response should have "$..status_message" with the text "No callback defined"
  And a registration should not exist

Scenario: Register endpoint with no health check returns error
	Given I send a POST request to "/v1/register" with the following:
	"""
	{
		"message_name": "testmessage.register",
		"callback_url": "http://something.something/v1/callback"
	}
	"""
	Then the response status should be "500"
  And the JSON response should have "$..status_message" with the text "No healthcheck defined"
  And a registration should not exist
