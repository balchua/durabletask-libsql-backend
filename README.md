# Libsql backend for durabletask-go

To learn more about durabletask-go, see [durabletask-go](https://github.com/microsoft/durabletask-go)

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

