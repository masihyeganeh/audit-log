# Audit-Log

A simple audit-log system to log events of other services and provide a way to query them.

---

## How to run

Simply clone the project and use docker-compose to run it:
```bash
git clone https://github.com/masihyeganeh/audit-log.git
cd audit-log
docker-compose run --build
```

## Authentication

There are 3 predefined users in system:

| Username | Password | Has Read Permission | Has Write Permission |
|----------|----------|---------------------|----------------------|
| admin    | admin    | yes                 | yes                  |
| reader   | reader   | yes                 | no                   |
| writer   | writer   | no                  | yes                  |

Before being able to use the system, you need to log in with one of the users first:
```bash
curl http://localhost:8088/login \
  --header 'Content-Type: application/json' \
  --data '{
	"username": "admin",
	"password": "admin"
}'
```

It gives a response like:
```json
{
	"status": "ok",
	"response": {
		"jwt_token": "JWT_TOKEN"
	},
	"error": ""
}
```

You need to provide this token in header of every log or query request to system.

## How to ingest logs

Each event should have some common fields named `common_field_1` and `common_field_2` and it can have any number of additional variable fields.

_For now, every key and value should be string, it is explained why in [design decisions](#data-store)._
```bash
curl http://localhost:8088/log \
  --header 'Authorization: Bearer JWT_TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{
	"event_type": "customer_created",
	"common_field_1": "something",
	"common_field_2": "something else",
	"fields": {
		"identity": "1"
	}
}'
```
or
```bash
curl http://localhost:8088/log \
  --header 'Authorization: Bearer JWT_TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{
	"event_type": "customer_did_action",
	"common_field_1": "something",
	"common_field_2": "something else",
	"fields": {
		"action_name": "something",
		"resource": "the resource",
	}
}'
```
or ...

## How to query

To query events, you need to provide type of the event and optionally any other fields of it as filters.

_For now, every key and value should be string, it is explained why in [design decisions](#data-store)._

Query on common fields:
```bash
curl http://localhost:8088/query \
  --header 'Authorization: Bearer JWT_TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{
	"event_type": "customer_created",
	"filters": {
		"common_field_1": "something"
	}
}'
```

Query on variable fields:
```bash
curl http://localhost:8088/query \
  --header 'Authorization: Bearer JWT_TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{
	"event_type": "customer_created",
	"filters": {
		"identity": "1"
	}
}'
```

Query on both common and variable fields:
```bash
curl http://localhost:8088/query \
  --header 'Authorization: Bearer JWT_TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{
	"event_type": "customer_created",
	"filters": {
	    "common_field_1": "something",
	    "common_field_2": "something else",
		"identity": "1"
	}
}'
```

---

## Design decisions

### Code architecture
Because of nature of this project, you may need to support many ingestion protocols, and you may need to test many data stores, so I chose Hexagonal (Ports & Adapters) architecture.
By taking advantage of layering of this architecture, the main logic of the program can be used with any type of inputs and storages.

Because it is stated that the system is write-intensive, I spin up many workers in goroutines to accept events and return immediately, then batch those events in background and write them in data store when batch size is large enough or a deadline is met.   

And as a side-note, there are all kinds of shortcuts, and dirty implementations in the code because it is just a simple demo project.

### Data store
Because it is write-intensive, I selected Clickhouse as data store, because it is an immutable columnar insert-only database and has no lock when inserting data.
But on the other hand, this is not the best database to handle variable data, so I made two implementation of it using two separate methods for handling variable data, `Nested data` and `Map`.

A trade-off of using this database for variable fields is that all keys and values should be string. There is a simple way of implementing many Maps or Nested data, one for each type of data (one for strings, one for integers, one for floats, ...). It can be implemented in a future version.

Another trade-off of using Clickhouse is that it works better with batch inserts instead of single inserts. But as explained before, the system is already batching events, so it is not a problem here.

### Ingest protocol
For this first iteration, events are ingested to the system using simple json requests, but there are other ways for that.
For example, it can accept events from other microservices using sockets (either tcp or udp).
Another way to handle events with specific fields to them is using a data serialization system (like Apache Avro) to be able to define schema for each event type in the 3rd-party system.
It can even parse structured log files and extract events from them.

For now, it is using json, and it would be better to have a Swagger documentation for it, but it was overkill for this simple API.

### Authentication
Because it is a simple project, authentication is implemented in a simple and dirty way.
There are predefined users that can have read, write or read&write permission.
There could be an access level for each event type and user groups and all sorts of fancy solutions, but not for a simple demo project.
Because no framework is used to implement it, and it is using basic web-server functionality of Golang, authentication is handled in a silly way.
It could be a full-fledged solution, and it could be handled in middleware that handles everything better.
