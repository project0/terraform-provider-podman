package modifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// RequiresReplaceComputed is a modified version of RequiresPlace, it allows also the replacement when computed values do change!
func RequiresReplaceComputed() tfsdk.AttributePlanModifier {
	return requiresReplaceModifierComputed{}
}

// requiresReplaceModifierComputed is an AttributePlanModifier that sets RequiresReplace on the attribute.
type requiresReplaceModifierComputed struct{}

// Modify fills the AttributePlanModifier interface.
func (r requiresReplaceModifierComputed) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	if req.AttributeConfig == nil || req.AttributePlan == nil || req.AttributeState == nil {
		// shouldn't happen, but let's not panic if it does
		return
	}

	if req.State.Raw.IsNull() {
		// if we're creating the resource, no need to delete and
		// recreate it
		return
	}

	if req.Plan.Raw.IsNull() {
		// if we're deleting the resource, no need to delete and
		// recreate it
		return
	}

	if req.AttributePlan.Equal(req.AttributeState) {
		// if the plan and the state are in agreement, this attribute
		// isn't changing, don't require replace
		return
	}

	resp.RequiresReplace = true
}

// Description returns a human-readable description of the plan modifier.
func (r requiresReplaceModifierComputed) Description(ctx context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (r requiresReplaceModifierComputed) MarkdownDescription(ctx context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}
