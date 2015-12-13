@register_consuimer
Feature: Register Consumer
	In order to ensure quality
	As a user
	I want to be able to test functionality of my API

Scenario: Register endpoint with correct data registers consumer successfully
	Given I send a POST request to "/v1/register"
	Then the response status should be "200"
	And the JSON response should have "$..status_message" with the text "OK"
  And a new consumer should be registered in the database

Scenario: Register endpoint with no callback returns error
	Given I send a POST request to "/v1/register"
	Then the response status should be "500"
	And the JSON response should have "$..status_message" with the text "No callback defined"
  And a new consumer should not be registered in the database

Scenario: Register endpoint with no health check returns error
	Given I send a POST request to "/v1/register"
	Then the response status should be "500"
	And the JSON response should have "$..status_message" with the text "No healthcheck defined"
  And a new consumer should not be registered in the database
