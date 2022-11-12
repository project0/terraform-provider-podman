package validator

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func testStringToVals(str ...string) []attr.Value {
	vals := make([]attr.Value, len(str))
	for i := 0; i < len(str); i++ {
		vals[i] = types.StringValue(str[i])
	}
	return vals
}

func TestStringValidator_Octal(t *testing.T) {
	tests := []testValidatorCase{
		{
			desc: "Null and Unknown is valid",
			values: []attr.Value{
				types.StringUnknown(),
				types.StringNull(),
			},
			validator: MatchOctal(),
		},
		{
			desc: "Octal is valid",
			values: testStringToVals(
				"000",
				"644",
				"1755",
				"777",
				"600",
				"0400",
				"0777",
			),
			validator: MatchOctal(),
		},
		{
			desc: "Octal should fail",
			values: testStringToVals(
				"somestring",
				"1 2",
				"0",
				"12345",
				"9999",
				"66",
				"14444",
				"s644",
			),
			wantFail:  true,
			validator: MatchOctal(),
		},
	}
	testValidatorExecute(t, tests)
}

func TestStringValidator_TmpfSize(t *testing.T) {
	tests := []testValidatorCase{
		{
			desc: "Null and Unknown is valid",
			values: []attr.Value{
				types.StringUnknown(),
				types.StringNull(),
			},
			validator: MatchTmpfSize(),
		},
		{
			desc: "TmpfSize is valid",
			values: testStringToVals(
				"4095",
				"666k",
				"200m",
				"100g",
				"1%",
				"20%",
				"90%",
				"100%",
			),
			validator: MatchTmpfSize(),
		},
		{
			desc: "TmpfSize should fail",
			values: testStringToVals(
				"somestring",
				"1 2",
				"100Gi",
				"2Mib",
				"m",
				"k",
			),
			wantFail:  true,
			validator: MatchTmpfSize(),
		},
	}

	testValidatorExecute(t, tests)
}
