# Agent Development Guide

* Ensure we handle resource drift, e.g. a resource is saved in our state but removed on the server (404).

## Test

* Run tests with `TF_ACC=1 go test ./lavinmq -v`
* Don't write VCR recordings, use `LAVINMQ_RECORD=1 TF_ACC=1 dotenv -f .env go test ./lavinmq/ -v -run {TestName} -timeout 5s`
* If you can't record recordings due to LavinMQ not running, stop and ask operator to start the broker.
* Do not implement or test internal features (internal queue, internal exchange)

## Examples

* Always recompile with `make clean build` before trying an example, to ensure we're running the latest code

## Coding style

* Use `any` instead of `interface{}`

## Changelog

* Section Headers (in order):
  * NOTES
  * FEATURES
  * IMPROVEMENTS
  * BUG FIXES
  * DEPENDENCIES
  * DEPRECATED
  * BREAKING CHANGES

* Sub headers to above headers except notes (in order):
  * **Provider:**
  * **Resources:**
  * **Data Sources:**
  * **Client Library:**
  * **Testing:**
  * **Documentation:**

* Formatting Rules:
  * Past tense verbs: "Added", "Fixed", "Updated", "Removed", "Deprecated", "Bumped"
  * PR references: Always include ([#X]) at the end of each entry
  * PR links: Add reference links at the bottom of each release section
  * Consistent structure: Bullet points for all entries
  * Clear descriptions: Brief but informative
