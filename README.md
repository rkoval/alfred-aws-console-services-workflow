# <img src="icon.png" width="26"> AWS Console Services - Alfred Workflow

![build](https://github.com/rkoval/alfred-aws-console-services-workflow/workflows/build/badge.svg)

Powerful workflow for quickly opening up AWS Console Services in your browser or searching for entities within them.

![AWS Console Services - Alfred Workflow Demo](demo.gif)

## Installation
- [Download the latest release](https://github.com/rkoval/alfred-aws-console-services-workflow/releases)
- Open the downloaded file in Finder
- Make sure your AWS Credentials and Region are set in your `~/.aws/credentials` and `~/.aws/config` files, respectively. This workflow will use your `default` profile by default within these files.
  - You can override any/all configuration values in [the workflow environment variables](https://www.alfredapp.com/help/workflows/advanced/variables/#environment). See [the official AWS docs](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-the-region) for more info on how to configure these

## Usage
To use, activate Alfred and type `aws` to trigger this workflow. From there:

- type any search term to search for services
- press <kbd>Tab</kbd> to autocomplete into sub-services, if they exist (for example, navigate to "Security Groups" within the "EC2" service)
- keep typing after autocompleting to filter sub-services
- press <kbd>Tab</kbd> again within sub-services to autocomplete the sub-service and start searching for its entities (for example, you can search for EC2 Instances when tabbed to `aws ec2 instances `)

At any time:
- press <kbd>Enter</kbd> to open the result in your default browser
- press <kbd>⌘</kbd>+<kbd>Enter</kbd> to copy the result's URL to clipboard.

*Note that you must be logged in for the page to open directly to your service*. See [this config file](console-services.yml) for the full list of supported services and their sub-services and [this file](https://github.com/rkoval/alfred-aws-console-services-workflow/blob/master/searchtypes/search_types.go) for the list of supported searchers.

## Advanced Usage

- [Fuzzy filtering](https://godoc.org/github.com/deanishe/awgo/fuzzy) a la Sublime Text is supported
- `,` is a sub-section alias to start searching for the default entity type associated with a service. For example, in this workflow, the EC2 service's default entity is an EC2 instance, so `aws ec2 ,searchterm` is a shorter alias for `aws ec2 instances searchterm`. You can customize this alias by setting the `ALFRED_AWS_CONSOLE_SERVICES_WORKFLOW_SEARCH_ALIAS` environment variable to any other string.

## Contributing

See [this README](CONTRIBUTING.md)

## Packaging for Release

See [this README](release_tools/README.md)
