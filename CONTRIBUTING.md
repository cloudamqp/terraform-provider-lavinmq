# Contributing to terraform-provider-lavinmq

Thank you for your interest in contributing to the LavinMQ Terraform provider.

## Code Contributions

### Getting Started

Before starting work on a new feature or significant change, please discuss your plans on an existing or new GitHub issue. This helps avoid duplicate work and ensures the change aligns with the project's direction.

### Development Workflow

1. Fork the repository and create a feature branch from `main`
2. Make your changes
3. Run tests and linting
4. Submit a pull request

Use GitHub's draft/open PR status:
- **Draft** while implementing changes
- **Open** when ready for review
- Return to **Draft** if revisions are requested

### Prerequisites

- **Go**: Follow [Go's installation guide](https://go.dev/doc/install)
- **Terraform**: Version >= 0.12. Follow [Terraform's installation guide](https://developer.hashicorp.com/terraform/downloads)

Install development tools:

```sh
make tools
```

### Building

```sh
make build
```

To clean and rebuild:

```sh
make clean build
```

### Testing

The provider uses [Go-VCR](https://github.com/dnaeon/go-vcr) for recording and replaying HTTP interactions during tests.

**Run all tests (replay mode):**

```sh
make test
```

**Run a specific test:**

```sh
TF_ACC=1 go test ./lavinmq/ -v -run TestName
```

**Record new VCR cassettes:**

Requires a running LavinMQ instance. Set up a `.env` file:

```dotenv
LAVINMQ_API_BASEURL="http://localhost:15672/"
LAVINMQ_API_USERNAME="guest"
LAVINMQ_API_PASSWORD="guest"
TEST_AMQP_URI="amqp://guest:guest@localhost:5672//"
```

Then record:

```sh
LAVINMQ_RECORD=1 TF_ACC=1 dotenv -f .env go test ./lavinmq/ -v -run TestName -timeout 5s
```

### Linting

```sh
make lint
```

### Formatting

```sh
make fmt
```

Check formatting:

```sh
make fmtcheck
```

### Documentation

Documentation is auto-generated using [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs).

After modifying resources or data sources:

```sh
make docs
```

Commit the generated changes in `docs/`. A CI workflow validates that documentation is up to date on every PR.

## Pull Request Guidelines

- Keep changes focused and atomic
- Include tests for new functionality
- Ensure all tests pass
- Run linting before submitting
- Update documentation
- Handle resource drift (e.g., when a resource is deleted outside Terraform)

## Non-Code Contributions

### Report Issues

Submit bugs or feature requests on [GitHub Issues](https://github.com/cloudamqp/terraform-provider-lavinmq/issues).

### Share Feedback

Discuss your experience, suggest improvements, or ask questions via GitHub Issues or Discussions.
