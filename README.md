# river-example

[river](https://github.com/riverqueue/river)

## How to run

```sh
# run compose
docker compose up -d

# migrate db for river 
make db-migrate

# run worker first
make run-worker

# run producer to send job
make run-producer
```
