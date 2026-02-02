# 0.1.0 (2025-11-04)

NOTES:

* Initial commit
* Provider built using Terraform Framework SDK

FEATURES:

**Resources:**

Multiple resources added

* Binding ([#43])
* Exchange ([#12])
* Federation upstream ([#45])
* Permission ([#33])
* Policy ([#14], [#21])
* Publish message ([#50], [#52])
* Queue ([#9])
* Queue action (pause/resume/purge) ([#32], [#41])
* Shovel ([#44])
* User ([2a272b4])
* VHost ([#2])

**Data Sources:**

Multiple data sources added

* Bindings ([#43])
* Exchanges ([#25])
* Federation upstreams ([#45])
* Permissions ([#33])
* Policies ([#16], [#30])
* Queues ([#23])
* Shovels ([#44])
* Users ([#24])
* VHosts ([#2], [#18])

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

[#2]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/2
[#3]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/3
[#7]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/7
[#9]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/9
[#10]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/10
[#12]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/12
[#14]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/14
[#16]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/16
[#18]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/18
[#21]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/21
[#23]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/23
[#24]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/24
[#25]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/25
[#30]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/30
[#32]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/32
[#33]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/33
[#38]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/38
[#41]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/41
[#43]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/43
[#44]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/44
[#45]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/45
[#50]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/50
[#52]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/52
[2a272b4]: https://github.com/cloudamqp/terraform-provider-lavinmq/commit/2a272b4
