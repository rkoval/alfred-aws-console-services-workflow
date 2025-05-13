read @CONTRIBUTING.md to understand how to add functionality to this project. look around directories and add a searcher for `ecs clusters`. please add the relevant searchers and tests to ensure these work properly

you should not have to make any changes to core files. look to @ec2_instances.go and @ec2_instances_test.go for examples for how to implement this for ECS clusters.

1. @ec2_instances.go contains the bulk of the AWS logic for interfacing with AWS SDK
2. @ec2_instances_test.go is a very simple reference to set up a test. ECS clusters test will be just as simple
3. you will need to add this to @searchers_by_service_id.go
4. you will also need to add test cases to @workflow_test.go. almost no logic to write though; just simple queries that you will append as a `testCase` to `var tcs`
5. you will need to add the right type to @caching.go for `Entity`

once you've done all of this, run `./build.sh` and make sure there are no compile issues. if there are, fix them and re-run `./build.sh` until fixed.

then, run `./test.sh` from the root of the repo. this should run the new tests you added and generate a fixture from interacting with AWS SDK. make sure there was a successful response in the `*aws_fixture.yaml` file that's generated. obfuscate all other info that is contained within so that we can commit this file into source control without private info leaking out
