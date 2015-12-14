# event-sauce
NOTES:
Event Sauce is an open source distributed server for Event Sourcing, it performs the role of event store and pub sub queue using REST

## This document is work in progress the current API spec can be found at the below location
()[http://htmlpreview.github.io/?]

The system is made up of a server and either one or more workers.

Server:
The server is responsible for receiving requests and for registering subscribers, exposes a HTTP endpoint

Workers:
The workers publish the events, to the subscribers

API:
HEADERS: X-API-Key: asdfasdasdasd33434

## publish is the public interface which the clients use to distribute a message through the system
```
/publish/{event-name}
{
  "name": {
    "firstname": "Nic",
    "surname": "Jackson"
  },
  "roles": [
    "Admin",
    "Guest"
  ]
}
```
## subscribe is used by a service to register its endpoint in order to receive events of a particular type
```
/subscribe/{event-name}
{
  "server-id": "12345"
  "callback": "http://myserver.com/path" // POST
  "filter": {"data": "name.firstname=jack%"} # see filters
}
```

## unsubscirbe to one or multiple events
```
/unsubscribe/{event-name}/{server-id}
```

```
/unsubscribe/all/{server-id}
```

## events list will list all events sent by the system that match the criteria
```
/events/list/{event-name}/{limit}/{page}/{event_name}/{filter}
```

## replay a single event
```
/events/reply/single/{event_id}
```

## replay multiple events
```
/events/reply/multiple/{event_name}/{filter}/{limit}
```

## filters
```
{
  "date_start": "2015-10-12",
  "date-end": "2015-10-13", //optional
  "data": "name.firstname=nic%&name&roles[any]=Admin"
}
```
### Object Matcher
attempts to match a value on an data object using dot notation e.g. given the below object
```
{
  "name": {
    "firstname": "Nic",
    "surname": "Jackson"
  },
  "roles": [
    "Admin",
    "Guest"
  ]
}
```
name.firstname=nic% would match this message
```
= equality matcher
!= inequality matcher
> greater than (number)
< less than (number)
%wildcard
```
### Array matcher
```
[first] = first index
[any] = match any index
[last] = last index
[2] = 3rd index
```

### Array matchers can be combined with object matchers like...
`data[any]=(name.firstname=nic%)`
