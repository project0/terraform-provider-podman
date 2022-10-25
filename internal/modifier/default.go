package modifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/project0/terraform-provider-podman/internal/utils"
)

// UseDefaultModifier returns a new DefaultModifier
func UseDefaultModifier(value attr.Value) tfsdk.AttributePlanModifier {
	return &DefaultModifier{
		value: value,
	}
}

// DefaultModifier ensures a new value with null value will always be set to default
type DefaultModifier struct {
	tfsdk.AttributePlanModifier
	value attr.Value
}

// Modify sets the default value if not known or empty
func (r DefaultModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	val, err := req.AttributePlan.ToTerraformValue(ctx)
	if err != nil {
		utils.AddUnexpectedAttributeError(req.AttributePath, &resp.Diagnostics, "Failed to retrieve value", err.Error())
		return
	}
	if val.IsNull() {
		resp.AttributePlan = r.value
	}
}

// Description returns a human-readable description of the plan modifier.
func (r DefaultModifier) Description(ctx context.Context) string {
	return "Ensure null values are replaced by given default attribute."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (r DefaultModifier) MarkdownDescription(ctx context.Context) string {
	return r.Description(ctx)
}
