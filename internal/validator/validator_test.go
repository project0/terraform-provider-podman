package validator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type testValidatorCase struct {
	desc      string
	values    []attr.Value
	wantFail  bool
	validator tfsdk.AttributeValidator
}

func testValidatorExecute(t *testing.T, testCases []testValidatorCase) {
	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {

			for _, val := range test.values {
				req := tfsdk.ValidateAttributeRequest{
					AttributeConfig: val,
				}
				resp := &tfsdk.ValidateAttributeResponse{}

				test.validator.Validate(context.TODO(), req, resp)
				if test.wantFail != resp.Diagnostics.HasError() {
					t.Errorf("%s: value=%[2]q", test.desc, val)
				}
			}

		})
	}
}
