package validators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type (
	genericStringValidator struct {
		description string
		validate    func(context.Context, validator.StringRequest, *validator.StringResponse)
	}
)

func (v *genericStringValidator) Description(ctx context.Context) string {
	return v.description
}

func (v *genericStringValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *genericStringValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	v.validate(ctx, req, resp)
}
