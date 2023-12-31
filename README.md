# Libsql backend for durabletask-go

To learn more about durabletask-go, see [durabletask-go](https://github.com/microsoft/durabletask-go)

**This is a blatant copy of the sqlite backend implementation with some minor changes to make it work with libsql.**

What is durabletask-go?

> The Durable Task Framework is a lightweight, embeddable engine for writing durable, fault-tolerant business logic (orchestrations) as ordinary code. 

What is Libsql?

> [libsql](https://github.com/tursodatabase/libsql) is an open source, open contribution fork of SQLite, created and maintained by Turso.
> Features include:
> * Embedded replicas that allow you to have replicated database inside your app.
> * libSQL server for remote SQLite access, similar to PostgreSQL or MySQL
> * Supports Rust, JavaScript, Python, Go, and more.
> 
> There are also various improvements and extensions to the core SQLite:
> 
> * ALTER TABLE extension for modifying column types and constraints
> * Randomized ROWID
> * WebAssembly User Defined Functions
> * Pass down SQL string to virtual table implementation
> * Virtual write-ahead log interface
>

## Installation

2 ways to start using libsql server backend.

* The easiest to get started is to use turso cloud service.  See [turso](https://turso.io) for more information.

* Spin yp your own libsql server backend.  See [libsql](https://github.com/tursodatabase/libsql)

## Example

The example code can be found in [example](example) folder.

Requires the following environment variables to be set:

```bash
# For local development
export DURABLETASK_DEMO_DB_HOST=localhost:8181
export DURABLETASK_DEMO_DB_TOKEN=abc1234567890
export DURABLETASK_DEMO_DB_HOST_SCHEME=http

go run ./example/
```

If you are using turso, here is the example:

```bash
export DURABLETASK_DEMO_DB_HOST=mydb.turso.io
export DURABLETASK_DEMO_DB_TOKEN=ey......3vV7BW
export DURABLETASK_DEMO_DB_HOST_SCHEME=libsql
```

Sample output:

```bash
{"time":"2023-12-31T09:28:34.95456265+08:00","level":"INFO","msg":"creating the task registry","orchestrator":"SimpleOrchestration"}
{"time":"2023-12-31T09:28:35.361891+08:00","level":"INFO","msg":"worker started with backend host: libsql://maximum-spirit-balchua.turso.io"}
{"time":"2023-12-31T09:28:36.229099278+08:00","level":"INFO","msg":"fa0208b1-a743-400f-a4db-f36c9f64e425: starting new 'SimpleOrchestration' instance with ID = 'fa0208b1-a743-400f-a4db-f36c9f64e425'."}
{"time":"2023-12-31T09:28:40.429663286+08:00","level":"INFO","msg":"fa0208b1-a743-400f-a4db-f36c9f64e425: 'SimpleOrchestration' completed with a COMPLETED status."}
{"time":"2023-12-31T09:28:41.116989021+08:00","level":"INFO","msg":"Orchestration completed: [{1 John Doe 30} {4 Lily Smith 35}]"}
{"time":"2023-12-31T09:28:41.117057759+08:00","level":"INFO","msg":"backend stopping..."}
{"time":"2023-12-31T09:28:41.117087488+08:00","level":"INFO","msg":"workers stopping and draining..."}
{"time":"2023-12-31T09:28:41.117424496+08:00","level":"INFO","msg":"orchestration-processor: received cancellation signal"}
{"time":"2023-12-31T09:28:41.117482421+08:00","level":"INFO","msg":"orchestration-processor: stopped listening for new work items"}
{"time":"2023-12-31T09:28:41.118393535+08:00","level":"INFO","msg":"activity-processor: received cancellation signal"}
```

