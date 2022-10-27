package validator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type (
	genericStringValidator struct {
		tfsdk.AttributeValidator
		description string
		validate    func(context.Context, tfsdk.ValidateAttributeRequest, *tfsdk.ValidateAttributeResponse, string)
	}
)

func (v *genericStringValidator) Description(ctx context.Context) string {
	return v.description
}

func (v *genericStringValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *genericStringValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var str types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &str)
	resp.Diagnostics.Append(diags...)

	if diags.HasError() {
		return
	}

	if str.IsUnknown() || str.IsNull() {
		return
	}

	v.validate(ctx, req, resp, str.ValueString())
}
