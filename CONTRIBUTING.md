# Contributing

## Adding services or sub-services

If you're just wanting to add/modify/remove services or sub-services, you shouldn't need to make changes to any of the go files. You can simply update [the .yml config file](console-services.yml) with your modifications. This file is used by the executable to populate entries in Alfred (for a list of valid properties, see [the models file](awsworkflow/aws_service.go))

You can then simply submit a pull request with your changes to the .yml file for it to be reviewed and accepted.

## Adding a searcher

If you're wanting to add a searcher that will query specific AWS entities for navigating directly to them, you'll need to write some go for this.

### Requirements
- go 1.14.0 (or later)

### Installation
Clone the repository into your `$GOPATH` and add that directory as a workflow within Alfred (make sure to disable any other versions of this workflow in Alfred so that the `aws` keyword doesn't conflict). Then, from the root of this repo, run:

```sh
# you must run this whenever you make changes to the go files
./build.sh
```

You should now be able to invoke the Alfred workflow from your cloned repo.

For more rapid development, you can also run:

```sh
./run.sh -query='ec2 instances' # example query; don't prefix string with `aws` here!
```

This will build and run a single execution of the workflow and log any relevant output to the console. Keep in that the `./run.sh` environment is not 100% the same as running within Alfred, so you should always verify your changes via `./build.sh` and invoking this workflow manually in Alfred.

### Guidelines

Generally, if you're just wanting to add a searcher, you can follow the patterns/examples from already implemented searchers (like [the EC2 Instance searcher](https://github.com/rkoval/alfred-aws-console-services-workflow/blob/master/searchers/ec2_instances.go)). Some notes to keep in mind:

- **Each new searcher must have tests**. The tests in this repo use the following libraries/patterns:
  - [go-vcr](https://github.com/dnaeon/go-vcr) for recording AWS requests to be re-used deterministically in CI. **Ensure that any new AWS fixtures that you add are purged of any sensitive information related to your account.** There are [sanitizers](https://github.com/rkoval/alfred-aws-console-services-workflow/blob/1178d7c9ff81e763e4898dd1450f642974e3b5c7/tests/test_tools.go#L52-L112) implemented to scrub some of this data (please feel free to add to them), though this is hard to exhaustively do automatically.
  - [cupaloy](https://github.com/bradleyjkemp/cupaloy) for snapshot testing to assert results in the Alfred Workflow. The `./test.sh` script should automatically update snapshots on every run, so make sure that the snapshots look the way that they should before committing them.
- Each new searcher should have a test case in [this file](https://github.com/rkoval/alfred-aws-console-services-workflow/blob/1178d7c9ff81e763e4898dd1450f642974e3b5c7/workflow/workflow_test.go). This file is an integration-style test that will massage logic across the entire workflow. The new test case should have a query that matches the searcher that you added.
- Please verify that your searcher/querying works in both Alfred 3 and Alfred 4
