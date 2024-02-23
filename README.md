# Terraform Sendbird Provider

Initial Project to develop a Terraform provider for Senbird.
Contributions are welcome and if some resources would be of a higher priority for you, please create an issue
and I will add it to my priority list.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/install)
- [Go](https://go.dev/doc/install)
- [GNU Make](https://www.gnu.org/software/make/)

## Development

### Building

1. `git clone` this repository and `cd` into its directory
2. `make` will trigger the Golang build
The provided `makefile` defines additional commands generally useful during development, like for running tests, generating documentation, code formatting and linting. Taking a look at it's content is recommended.

## Testing

Testing of Terraform Resources and Data Sources are implemented as acceptance testing.

### Running a Acceptance Tests

Acceptance testing will create and destroy real resources against your senbird account.
At the moment there is no known costs associated with any of the resources created, but it is recommended
to run in a development environment to avoid potential side effects.

It is required to set your api key in you environment variables.

```sh
SENDBIRD_API_KEY=XXX
```

Once done, execute one of the following commands to run the test

```sh
# To run all tests
make testacc

# To run only a specific package
make testacc PKG=application

# To run only tests matching a Regex
make testacc TESTS=TestAccApplicationDataSource_

# or combined
make testacc TESTS=_create PKG=application
```

## Generating the documentation

You can genarate the documentation for the Sendbird provider by using the following make command.

```sh
make generate
```

## Using a development build

If [running tests and acceptance tests](#testing) isn't enough, it's possible to set up a local terraform configuration
to use a development builds of the provider. This can be achieved by leveraging the Terraform CLI
[configuration file development overrides](https://www.terraform.io/cli/config/config-file#development-overrides-for-provider-developers).

First, use `make install` to place a fresh development build of the provider in your
[`${GOBIN}`](https://pkg.go.dev/cmd/go#hdr-Compile_and_install_packages_and_dependencies)
(defaults to `${GOPATH}/bin` or `${HOME}/go/bin` if `${GOPATH}` is not set). Repeat
this every time you make changes to the provider locally.

Then, setup your environment following [these instructions](https://www.terraform.io/plugin/debugging#terraform-cli-development-overrides)
to make your local terraform use your local build.
