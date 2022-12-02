package validator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type (
	conflictsWithValidator struct {
		path func(currPath path.Path) path.Path
		attr func() attr.Value
	}
)

// ConflictsWith validates if a given path is also set
func ConflictsWith(p func(currPath path.Path) path.Path, a func() attr.Value) tfsdk.AttributeValidator {
	return &conflictsWithValidator{
		path: p,
		attr: a,
	}
}

func (v *conflictsWithValidator) Description(ctx context.Context) string {
	return "Conflicts with another attribute"
}

func (v *conflictsWithValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *conflictsWithValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	a := v.attr()
	p := v.path(req.AttributePath)

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, p, a)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if req.AttributeConfig.IsNull() || req.AttributeConfig.IsUnknown() {
		return
	}

	if a.IsNull() || a.IsUnknown() {
		return
	}

	resp.Diagnostics.AddAttributeError(req.AttributePath, "attribute conflict", "conflcits with: "+p.String())
}
