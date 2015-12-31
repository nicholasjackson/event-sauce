@register_consumer
Feature: Register Consumer
	In order to ensure quality
	As a user
	I want to be able to test functionality of my API

Scenario: Invalid GET request
  Given I send a GET request to "/v1/register"
  Then the response status should be "404"

Scenario: Invalid POST request
	Given I send a POST request to "/v1/register"
	Then the response status should be "400"

Scenario: Register endpoint with correct data registers consumer successfully
	Given I send a POST request to "/v1/register" with the following:
		| event_name | testevent.register                   |
		| callback_url | http://something.something/v1/callback |
	Then the response status should be "200"
	And the JSON response should have "$..status_event" with the text "OK"
  And registration should exist with event_name: "testevent.register"

Scenario: Register endpoint with no callback_url returns error
	Given I send a POST request to "/v1/register" with the following:
		| event_name | testevent.register                   |
	Then the response status should be "400"
  And 0 registrations should exist

Scenario: Register endpoint with no event returns error
	Given I send a POST request to "/v1/register" with the following:
		| callback_url | http://something.something/v1/callback |
	Then the response status should be "400"
  And 0 registrations should exist

Scenario: Register endpoint with correct data when registration exists returns 304
	Given the following registrations exist
		| event_name         | callback_url                           |
		| testevent.register | http://something.something/v1/callback |
	When I send a POST request to "/v1/register" with the following:
		| event_name | testevent.register                   |
		| callback_url | http://something.something/v1/callback |
	Then the response status should be "304"
	And 1 registrations should exist with event_name: "testevent.register"

	Scenario: Invalid DELETE request
		Given I send a DELETE request to "/v1/register"
		Then the response status should be "400"

	Scenario: Delete registration with no callback_url returns error
		Given the following registrations exist
			| event_name         | callback_url                           |
			| testevent.register | http://something.something/v1/callback |
		When I send a DELETE request to "/v1/register" with the following:
			| event_name | testevent.register                   |
		Then the response status should be "400"
	  And 1 registrations should exist

	Scenario: Delete registration with no event returns error
		Given the following registrations exist
			| event_name         | callback_url                           |
			| testevent.register | http://something.something/v1/callback |
		When I send a DELETE request to "/v1/register" with the following:
			| callback_url | http://something.something/v1/callback |
		Then the response status should be "400"
	  And 1 registrations should exist

Scenario: Delete registration with correct data when registration exists returns 200
	Given the following registrations exist
		| event_name         | callback_url                           |
		| testevent.register | http://something.something/v1/callback |
	When I send a DELETE request to "/v1/register" with the following:
		| event_name | testevent.register                   |
		| callback_url | http://something.something/v1/callback |
		Then the response status should be "200"
		And the JSON response should have "$..status_event" with the text "OK"
	  And 0 registrations should exist

Scenario: Delete registration with correct data when not registration exists returns 304
	Given I send a DELETE request to "/v1/register" with the following:
		| event_name | testevent.register                   |
		| callback_url | http://something.something/v1/callback |
	Then the response status should be "404"
	And 0 registrations should exist
