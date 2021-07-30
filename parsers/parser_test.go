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
		rawQuery: "ec2 ",
	},
	{
		rawQuery: " ec2",
	},
	{
		rawQuery: " ec2 ",
	},
	{
		rawQuery: "      ec2      ",
	},
	{
		rawQuery: "ec2 instances ",
	},
	{
		rawQuery: "ec2 instances i-0000",
	},
	{
		rawQuery: "ec2 OPEN_ALL ",
	},
	{
		rawQuery: "ec2 OPEN_ALL instances",
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
		rawQuery: "ec2 $us-west-2",
	},
	{
		rawQuery: "ec2 $us-west-2 ",
	},
	{
		rawQuery: "ec2 $us-whoops-2 ",
	},
	{
		rawQuery: "ec2 ,search $us-west-2",
	},
	{
		rawQuery: "ec2 ,search $us-west-2 ",
	},
	{
		rawQuery: "OPEN_ALL ec2 instances",
	},
	{
		rawQuery: "asdf asdf",
	},
	{
		rawQuery: "asdf asdf asdf ",
	},
	{
		rawQuery: "ec2 ,search",
	},
	{
		rawQuery: "ec2 ,search term more hello",
	},
	{
		rawQuery: "ec2 ,search term more hello ",
	},
	{
		rawQuery: " ec2 ,search term more hello ",
	},
	{
		rawQuery: "ec2 subservice search term more hello",
	},
	{
		rawQuery: " ec2 subservice search term more hello",
	},
	{
		rawQuery: "ec2 subservice search term more hello ",
	},
	{
		rawQuery: " ec2 subservice search term more hello ",
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
