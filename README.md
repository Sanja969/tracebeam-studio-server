# Tracebeam Studio Server

Local ingestion server for Tracebeam events.

## What it does

- Receives events from the Tracebeam SDK
- Stores events in SQLite
- Broadcasts realtime updates over WebSocket
- Supports filtering by type, traceId, sessionId and limit

---

## Run locally

```bash
go run ./cmd/api
```

Server runs on:

http://localhost:8080

---

## API

### POST /events

Receives a Tracebeam event.

### GET /events

Query params:


| Param | Example |
|---|---|
| limit | /events?limit=100 |
| type | /events?type=error |
| traceId | /events?traceId=checkout-flow |
| sessionId | /events?sessionId=session_123 |


### DELETE /events

Clears all events.

### WS /ws

Realtime event stream.

## Status

Tracebeam Studio Server is under active development.

## License

MIT
