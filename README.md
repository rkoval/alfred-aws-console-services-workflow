# <img src="icon.png" width="26"> AWS Console Services - Alfred Workflow

Very simple workflow to quickly open up AWS Console Services in your browser.

![AWS Console Services - Alfred Workflow Demo](demo.gif)

### Installation
- [Download the latest release](https://github.com/rkoval/alfred-aws-console-services-workflow/releases)
- Open the downloaded file in Finder

### Usage
To use, activate Alfred and type in `aws`. From there, query any of the services offered on the AWS homepage dashboard. *Note that you must be logged in in order for the page to open directly to your service*. Your region will automatically be populated to the default tied to your account.

Open up the Alfred Workflow configuration to see the full list of supported services and their identifiers.

### Contributing

#### Requirements
- go 1.14.0 (or later)

#### Installation
From the root of this repo, run:

```sh
go get github.com/rkoval/alfred-aws-console-services-workflow
```

#### Adding Entries

To add entries to the Alfred Workflow, modify [the .yml config file](console-services.yml). This file is used by the executable to populate entries in Alfred.

#### Packaging for Release

See [this README](release_tools/README.md)