# Agent Development Guide

* Ensure we handle resource drift, e.g. a resource is saved in our state but removed on the server (404).

## Test

* Run tests with `TF_ACC=1 go test ./lavinmq -v`
* Don't write VCR recordings, use `LAVINMQ_RECORD=1 TF_ACC=1 dotenv -f .env go test ./lavinmq/ -v -run {TestName} -timeout 5s`
* Do not implement or test internal features (internal queue, internal exchange)

## Coding style

* Use `any` instead of `interface{}`
