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
To use, activate Alfred and type in `aws`. From there, type to query any of the services offered on the AWS homepage dashboard. You can hit `Tab` to populate sub-services, if they exist (for example, navigate to "Security Groups" within the "EC2" service). If the service is configured At any time, hit `Enter` to navigate to your result or `Cmd+Enter` to copy the URL to clipboard.

*Note that you must be logged in for the page to open directly to your service*. See [this config file](console-services.yml) for the full list of supported services and their sub-services and [this file](https://github.com/rkoval/alfred-aws-console-services-workflow/blob/master/workflow/searchers_by_service_id.go) for the list of supported searchers.

## Contributing

See [this README](CONTRIBUTING.md)

## Packaging for Release

See [this README](release_tools/README.md)
