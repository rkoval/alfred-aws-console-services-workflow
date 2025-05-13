# Contributing

## Requirements

- go 1.24.3 (or later)

## Installation

Clone the repository into your `$GOPATH` and add that directory as a workflow within Alfred (make sure to disable any other versions of this workflow in Alfred so that the `aws` keyword doesn't conflict). Then, from the root of this repo, run:

```sh
# you must run this whenever you make changes to the go files
./build.sh
```

You should now be able to invoke the Alfred workflow from your cloned repo.

## Adding services or sub-services

To add/update/remove services or sub-services, simply update [the .yml config file](console-services.yml) with your changes. This file is used by the executable to populate entries in Alfred (for a list of valid properties, see [the models file](awsworkflow/aws_service.go)). Once you've done that, run `./test.sh` to update any snapshots with the data changes you've made.

### Bulk Verifying URLs

If you're wanting to verify all sub-service URLs for a particular service, you can use `OPEN_ALL` as a sub-service term within the Alfred workflow. For example, to open all EC2 sub-service URLs, use `aws ec2 OPEN_ALL`. This will allow you to just go through the tabs opened in your browser to verify that the links are not broken.

## Adding a searcher

Generally, if you're just wanting to add a searcher, you can follow the patterns/examples from already implemented searchers (like [the EC2 Instance searcher](https://github.com/rkoval/alfred-aws-console-services-workflow/blob/master/searchers/ec2_instances.go)). Some things to keep in mind:

- **Each new searcher must have tests**. The tests in this repo use the following libraries/patterns:
  - [go-vcr](https://github.com/dnaeon/go-vcr) for recording AWS requests to be re-used deterministically in CI. **Ensure that any new AWS fixtures that you add are purged of any sensitive information related to your account.** There are [sanitizers](https://github.com/rkoval/alfred-aws-console-services-workflow/blob/1178d7c9ff81e763e4898dd1450f642974e3b5c7/tests/test_tools.go#L52-L112) implemented to scrub some of this data (please feel free to add to them), though this is hard to exhaustively do automatically.
  - [cupaloy](https://github.com/bradleyjkemp/cupaloy) for snapshot testing to assert results in the Alfred Workflow. The `./test.sh` script should automatically update snapshots on every run, so make sure that the snapshots look the way that they should before committing them.
- Each new searcher should have a test case in [this file](https://github.com/rkoval/alfred-aws-console-services-workflow/blob/master/workflow/workflow_test.go). This file is an integration-style test that will massage logic across the entire workflow. The new test case should have a query that matches the searcher that you added.

## Advanced

For more rapid development, you can also run:

```sh
./run.sh -query='ec2 instances' # example query; don't prefix string with `aws` here!
```

This will build and run a single execution of the workflow and log any relevant output to the console. Keep in that the `./run.sh` environment is not 100% the same as running within Alfred, so you should always verify your changes via `./build.sh` and invoking this workflow manually in Alfred.
