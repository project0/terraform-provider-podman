package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type testValidatorStringCase struct {
	desc      string
	values    []types.String
	wantFail  bool
	validator validator.String
}

func testValidatorStringExecute(t *testing.T, testCases []testValidatorStringCase) {
	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {

			for _, val := range test.values {
				req := validator.StringRequest{
					ConfigValue: val,
				}
				resp := &validator.StringResponse{}

				test.validator.ValidateString(context.TODO(), req, resp)
				if test.wantFail != resp.Diagnostics.HasError() {
					t.Errorf("%s: value=%[2]q, err: %v", test.desc, val.ValueString(), resp.Diagnostics)
				}
			}

		})
	}
}
