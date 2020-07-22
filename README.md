# <img src="icon.png" width="26"> AWS Console Services - Alfred Workflow

Powerful workflow for quickly opening up AWS Console Services in your browser or searching within them.

![AWS Console Services - Alfred Workflow Demo](demo.gif)

## Installation
- [Download the latest release](https://github.com/rkoval/alfred-aws-console-services-workflow/releases)
- Open the downloaded file in Finder

## Usage
To use, activate Alfred and type in `aws`. From there, type to query any of the services offered on the AWS homepage dashboard. You can hit `Tab` to populate service sections, if they exist (for example, navigate to the "Security Groups" section within the "EC2" service). At any time, hit `Enter` to navigate to your result or `Cmd+Enter` to copy the URL to clipboard.

*Note that you must be logged in for the page to open directly to your service*. Your region will automatically be populated to the default tied to your account. See [this config file](console-services.yml) for the full list of supported services and their sections.

## Contributing

![tests](https://github.com/rkoval/alfred-aws-console-services-workflow/workflows/test/badge.svg)

### Adding services or service sections

If you're just wanting to add/modify/remove services or service sections, you shouldn't need to make changes to any of the go files. You can simply update [the .yml config file](console-services.yml) with your modifications. This file is used by the executable to populate entries in Alfred (for a list of valid properties, see [the models file](core/aws_service.go))

You can then simply submit a pull request with your changes to the .yml file for it to be reviewed and accepted.

### Adding a searcher

If you're wanting to add a searcher that will query specific AWS entities for navigating directly to them, you'll need to write some go for this.

#### Requirements
- go 1.14.0 (or later)

#### Installation
Clone the repository into your `$GOPATH` and add that directory as a workflow within Alfred (make sure to disable any other versions of this workflow in Alfred so that the `aws` keyword doesn't conflict). Then, from the root of this repo, run:

```sh
# you must run this whenever you make changes to the go files
go build
```

You should now be able to invoke the Alfred workflow from your cloned repo.

For more rapid development, you can also run:

```sh
./test.sh 'ec2 instances' # example query; don't prefix string with `aws` here!
```

This will build and run a single execution of the workflow and log any relevant output to the console. Keep in that the `./test.sh` environment is not 100% the same as running within Alfred, so you should always verify your changes via `go build` and invoking this workflow manually in Alfred.

### Packaging for Release

See [this README](release_tools/README.md)