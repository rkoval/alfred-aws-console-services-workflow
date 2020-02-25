# Creating a release after a PR

This is unfortunately a fairly manual process since I'm unaware of any reliable automated way to package/upload an Alfred Workflow to GitHub

1. Merge PR
1. Do modifications from PR if necessary
1. Run `./prepare_package.sh` from the root of this repo.
    - If no change in git, then the PR contributor remembered to do this and everything is good (hooray).
1. Open Alfred and navigate to Workflows
1. Find "AWS Console Services" --> Right Click --> Export...
1. Bump version number up appropriately.
1. Append `v2.6` or whatever version before the file extension to the resultant .alfredworkflow file
1. Run `./commit_version.sh` which will open the GitHub Releases page.
1. Fill out the description and upload the .alfredworkflow file to the prepopulated release