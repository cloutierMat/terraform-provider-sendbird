# Terraform Sendbird Provider

Initial Project ot develop a Terraform provider for Senbird.

## Contributing

### Testing

Testing of Terraform Resources and Data Sources are implemented as acceptance testing.

#### Running a Acceptance Tests

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
