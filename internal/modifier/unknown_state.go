package modifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func AlwaysUseStateForUnknown() tfsdk.AttributePlanModifier {
	return AlwaysUseStateForUnknownModifier{}
}

type AlwaysUseStateForUnknownModifier struct{}

func (m AlwaysUseStateForUnknownModifier) Description(_ context.Context) string {
	return ""
}

func (m AlwaysUseStateForUnknownModifier) MarkdownDescription(_ context.Context) string {
	return ""
}

func (m AlwaysUseStateForUnknownModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	if req.AttributeState == nil || resp.AttributePlan == nil || req.AttributeConfig == nil {
		return
	}

	// if we're creating the resource, no need to modify
	if req.State.Raw.IsNull() {
		return
	}

	// if we're deleting the resource, no need to modify
	if req.Plan.Raw.IsNull() {
		return
	}

	// if it's not planned to be the unknown value, stick with the concrete plan
	if !resp.AttributePlan.IsUnknown() {
		return
	}

	resp.AttributePlan = req.AttributeState
}
