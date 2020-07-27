# <img src="icon.png" width="26"> AWS Console Services – Alfred Workflow

![build](https://github.com/rkoval/alfred-aws-console-services-workflow/workflows/build/badge.svg)

A powerful workflow for quickly opening up AWS Console Services in your browser or searching for entities within them.

Supports Alfred 3 and 4

![AWS Console Services - Alfred Workflow Demo](demo.gif)

## Installation
- [Download the latest release](https://github.com/rkoval/alfred-aws-console-services-workflow/releases)
- Open the downloaded file in Finder
- Make sure your AWS Credentials and Region are set in your `~/.aws/credentials` and `~/.aws/config` files, respectively. This workflow will use your `default` profile by default within these files. See [the official AWS docs](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-the-region) for more info on how to configure these
- If running on macOS Catalina or later, you _**MUST**_ add Alfred to the list of security exceptions for running unsigned software. See [this guide](https://github.com/deanishe/awgo/wiki/Catalina) for instructions on how to do this.
  - <sub>Yes, this sucks and is annoying, but there is unfortunately is no easy way around this. macOS requires a paying Developer account for proper app notarization. I'm afraid I'm not willing to pay a yearly subscription fee to Apple just so that this (free and open source) project doesn't pester macOS Gatekeeper.</sub>

## Usage
To use, activate Alfred and type `aws` to trigger this workflow. From there:

- type any search term to search for services
- if the current service result has a 🗂 in the subtitle, press <kbd>Tab</kbd> to autocomplete into sub-services (for example, navigate to "Security Groups" within the "EC2" service)
- keep typing after autocompleting to filter sub-services
- if the current sub-service result has a 🔎 in the subtitle, press <kbd>Tab</kbd> again to start searching for its entities (for example, you can search for EC2 Instances when tabbed to `aws ec2 instances `)

At any time:
- press <kbd>Enter</kbd> to open the current result in your default browser
- press <kbd>⌘</kbd>+<kbd>Enter</kbd> to copy the result's URL to clipboard.

*Note that you must be logged in for the page to open directly to your service*. See [this config file](console-services.yml) for the full list of supported services and their sub-services and [this file](https://github.com/rkoval/alfred-aws-console-services-workflow/blob/master/searchtypes/search_types.go) for the list of supported searchers.

## Advanced Features

- [Fuzzy filtering](https://godoc.org/github.com/deanishe/awgo/fuzzy) a la Sublime Text is supported
- Configurable [workflow environment variables](https://www.alfredapp.com/help/workflows/advanced/variables/#environment)
  - Search alias – If a sub-service has a ⭐ in the subtitle, you can use `,` as an alias for it to more quickly search for that entity. For example, in this workflow, the EC2 service's default entity is an EC2 instance, so `aws ec2 ,searchterm` is a shorter alias for `aws ec2 instances searchterm`. You can customize this alias by setting the `ALFRED_AWS_CONSOLE_SERVICES_WORKFLOW_SEARCH_ALIAS` environment variable to any other string.
  - Cache expiration age – Sub-service entity searching makes heavy use of caching to make filtering performant and to prevent handling big requests/responses to/from AWS on every execution. The cache expiration age for each entity is set to 3 minutes by default. If you find that this is too short/long for your usage, you can set the `ALFRED_AWS_CONSOLE_SERVICES_WORKFLOW_MAX_CACHE_AGE_SECONDS` environment variable to the number of seconds that better suits your need.
  - AWS settings – You can override any/all AWS configuration values which the underlying AWS library should respect.

## Contributing

See [this README](CONTRIBUTING.md)

## Packaging for Release

See [this README](release_tools/README.md)

## Troubleshooting

- "I'm seeing the following dialog when running the workflow"

  ![image](https://user-images.githubusercontent.com/1282943/88503823-6eda4b80-cf98-11ea-9a4b-f2a5bdb8a1cc.png)

  Per [the installation steps](https://github.com/rkoval/alfred-aws-console-services-workflow#installation), you **_MUST_** add Alfred to the list of Developer Tool exceptions for Alfred to run any workflow that contains an executable (like this one)


