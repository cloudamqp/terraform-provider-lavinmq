# 0.1.0 (Unreleased)

NOTES:

* Initial commit
* Provider built using Terraform Framework SDK

FEATURES:

**Resources:**

Multiple resources added

* Binding ([#43])
* Exchange ([#12])
* Federation upstream
* Permission ([#33])
* Policy
* Publish message ([#50], [#52])
* Queue ([#9])
* Queue action (pause/resume/purge) ([#32], [#41])
* Shovel ([#44])
* User ([#10])
* VHost ([#2])

**Data Sources:**

Multiple data sources added

* Bindings ([#43])
* Exchanges ([#25])
* Federation upstreams
* Permissions ([#33])
* Policies ([#16], [#30])
* Queues ([#23])
* Shovels ([#44])
* Users ([#24])
* VHosts ([#18])

**Client Library:**

* Initial code to handle communication with HTTP API
* Add net/url package to build URL with path escaping ([#7])
* Multiple endpoints added
  * Bindings ([#43])
  * Exchanges ([#12])
  * Messages service to publish messages ([#50])
  * Parameters ([#44])
  * Permissions ([#33])
  * Policies ([#16])
  * Queue action (purge) ([#32])
  * Queues ([#9])
  * Shovel ([#44])
  * Users ([#10])
  * VHosts ([#2])

**Testing:**

* Added support for VCR testing of the provider ([#3])
* VCR tests run in parallel ([#38])
* VCR test: Use environmental variable sanitizer for configurable URIs ([#44])

**Documentation:**

* Add automatic documentation generation with tfplugindocs
* Add GitHub Actions workflow to validate documentation

IMPROVEMENTS:

* Queue: Make `auto_delete` computed ([#22])
* Policy: Use computed value for `priority` and `apply_to`
* Policy: Add test for invalid `apply_to` ([#20])
* Queue and Exchange: Support `arguments` attribute ([#36])
* Publish message: Add `publish_message_counter` argument ([#52])
* User: Support both password and hash as writeonly attributes
* Data Source Queues: Extend to count consumers and messages
* Data Source Exchanges: Add message_state
* Data Source Policies: Handle vhost filtering ([#30])
* Return nil when vhost doesn't exist ([#27])
* Return empty list instead of nil in data sources
* Handle resource drift for queue and user resources
* Refactor: Replace Client service fields with Services struct pattern ([#11])
* Refactor: Use `any` type alias instead of `interface{}` ([#26])
* Refactor: Use inline config in tests ([#37])
* Refactor: Make Client struct fields private ([#39])
* Fix plural naming for collections to match Go conventions ([#17])
* Update definition handling to use Object type for improved type safety
* Use standard JSON unmarshal instead of GenericUnmarshal ([#48])
* Provider configuration attributes can be set via environment variables
* Sort resources and data sources in ascending order
* Don't hardcode test parallelism, use GOMAXPROCS ([#40])

BUG FIXES:

* Unmarshal failed response body and return error ([#8])
* Fix parallel test race condition in CI ([#49])
* Fix user resource documentation to match actual schema
* Handle non-existing vhost in data sources
* ClientLibrary: Combine List* methods and include vhost filtering

DEPENDENCIES:

* Bump github.com/cloudflare/circl from 1.3.7 to 1.6.1 ([#5])
* Bump golang.org/x/net from 0.34.0 to 0.38.0 ([#1])

[#1]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/1
[#2]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/2
[#3]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/3
[#5]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/5
[#7]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/7
[#8]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/8
[#9]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/9
[#10]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/10
[#11]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/11
[#12]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/12
[#16]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/16
[#17]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/17
[#18]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/18
[#20]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/20
[#22]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/22
[#23]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/23
[#24]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/24
[#25]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/25
[#26]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/26
[#27]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/27
[#30]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/30
[#32]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/32
[#33]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/33
[#36]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/36
[#37]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/37
[#38]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/38
[#39]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/39
[#40]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/40
[#41]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/41
[#43]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/43
[#44]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/44
[#48]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/48
[#49]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/49
[#50]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/50
[#52]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/52
