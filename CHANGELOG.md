## 1.0.0 (2026-02-06)

NOTES:

The LavinMQ Terraform Provider has matured from its initial release (0.1.0) through iterative improvements (0.2.0) and
is now production-ready with version 1.0.0.

**v0.2.0 (Stability & Quality):**

* Enhanced resource drift handling and 404 error recovery
* Improved debugging capabilities with client library logging
* Significantly increased test coverage with automated reporting
* Refined import operations (fixed exchange import separator)
* Updated dependencies for better compatibility and security
* Comprehensive project documentation (contribution guidelines, security policy, issue/PR templates)

**v0.1.0 (Initial Release):**

* Comprehensive resource and data source support for LavinMQ management
* Full CRUD operations for bindings, exchanges, federation upstreams, permissions, policies, queues, shovels, users,
and vhosts
* Built on Terraform Framework Plugin with VCR testing infrastructure

## 0.2.0 (2026-02-02)

NOTES:

* Run tests against latest lavinmq in CI on release ([#92])

IMPROVEMENTS:

**Provider:**

* Use default values for provider arguments ([#75])

**Resources:**

* Fixed resource drift handling when VHost is deleted ([#73])
* Fixed handling of 404 errors for missing parent resources to prevent nil pointer crashes ([#91])

**Client Library:**

* Refactored client library code for better maintainability ([#66])
* Added debug logging for path and data ([#67])

**Testing:**

* Increased test coverage ([#81])
* Added test coverage report ([#80])
* Enabled support to run VCR in pass through mode ([#86])

**Documentation:**

* Added contribution guidelines ([#76])
* Added pull request template ([#77])
* Added security policy ([#78])
* Added issue templates ([#79])

BUG FIXES:

**Resources:**

* Fixed Exchange import separator - changed from ',' to '@' ([#90])

DEPENDENCIES:

* Bumped github.com/hashicorp/terraform-plugin-log from 0.9.0 to 0.10.0 ([#68])
* Bumped golang.org/x/crypto from 0.42.0 to 0.45.0 ([#69])
* Bumped actions/checkout from 5 to 6 ([#70])
* Bumped github.com/hashicorp/terraform-plugin-testing from 1.13.3 to 1.14.0 ([#71])
* Bumped github.com/hashicorp/terraform-plugin-framework ([#72])

[#66]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/66
[#67]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/67
[#68]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/68
[#69]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/69
[#70]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/70
[#71]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/71
[#72]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/72
[#73]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/73
[#75]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/75
[#76]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/76
[#77]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/77
[#78]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/78
[#79]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/79
[#80]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/80
[#81]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/81
[#86]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/86
[#90]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/90
[#91]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/91
[#92]: https://github.com/cloudamqp/terraform-provider-lavinmq/pull/92

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
