package parsers

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

type testCase struct {
	rawQuery string
}

var tcs []testCase = []testCase{
	{
		rawQuery: "",
	},
	{
		rawQuery: " ",
	},
	{
		rawQuery: "   ",
	},
	{
		rawQuery: "e",
	},
	{
		rawQuery: "elasticbeanstalk ",
	},
	{
		rawQuery: " elasticbeanstalk",
	},
	{
		rawQuery: " elasticbeanstalk ",
	},
	{
		rawQuery: "      elasticbeanstalk      ",
	},
	{
		rawQuery: "elasticbeanstalk environments ",
	},
	{
		rawQuery: "elasticbeanstalk environments e-0000",
	},
	{
		rawQuery: "elasticbeanstalk OPEN_ALL ",
	},
	{
		rawQuery: "elasticbeanstalk OPEN_ALL environments",
	},
	{
		rawQuery: "$",
	},
	{
		rawQuery: "$us",
	},
	{
		rawQuery: "$us-west-2",
	},
	{
		rawQuery: "elasticbeanstalk $us-west-2",
	},
	{
		rawQuery: "elasticbeanstalk $us-west-2 ",
	},
	{
		rawQuery: "elasticbeanstalk $us-whoops-2 ",
	},
	{
		rawQuery: "elasticbeanstalk ,search $us-west-2",
	},
	{
		rawQuery: "elasticbeanstalk ,search $us-west-2 ",
	},
	{
		rawQuery: "OPEN_ALL elasticbeanstalk environments",
	},
	{
		rawQuery: "asdf asdf",
	},
	{
		rawQuery: "asdf asdf asdf ",
	},
	{
		rawQuery: "elasticbeanstalk ,search",
	},
	{
		rawQuery: "elasticbeanstalk ,search term more hello",
	},
	{
		rawQuery: "elasticbeanstalk ,search term more hello ",
	},
	{
		rawQuery: " elasticbeanstalk ,search term more hello ",
	},
	{
		rawQuery: "elasticbeanstalk subservice search term more hello",
	},
	{
		rawQuery: " elasticbeanstalk subservice search term more hello",
	},
	{
		rawQuery: "elasticbeanstalk subservice search term more hello ",
	},
	{
		rawQuery: " elasticbeanstalk subservice search term more hello ",
	},
}

func TestParser(t *testing.T) {
	for _, tc := range tcs {
		t.Run(tc.rawQuery, func(t *testing.T) {
			parser := NewParser(tc.rawQuery)
			query, _ := parser.Parse("../console-services.yml")
			cupaloy.SnapshotT(t, query)
		})
	}
}
