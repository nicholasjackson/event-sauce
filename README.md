# event-sauce
Event Sauce is an open source distributed server for Event Sourcing, it performs the role of event store and pub sub queue using REST

The system is made up of a server and either one or more workers.

Server:
The server is responsible for receiving requests and for registering subscribers

Workers:
The workers publish the events, to the subscribers

API:


HEADERS: X-API-Key: asdfasdasdasd33434
/publish/{event-name}
{
  "data": "value"}
}

HEADERS: X-API-Key: asdfasdasdasd33434
/subscribe/{event-name}
{
  "callback": "http://myserver.com/path" // POST
}

/events/list/{event-name}/{page}/{limit}
