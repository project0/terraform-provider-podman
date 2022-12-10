package modifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// RequiresReplaceComputed is a modified version of RequiresPlace, it allows also the replacement when computed values do change!
func RequiresReplaceComputed() requiresReplaceModifierComputed {
	return requiresReplaceModifierComputed{}
}

// requiresReplaceModifierComputed is an AttributePlanModifier that sets RequiresReplace on the attribute.
type requiresReplaceModifierComputed struct{}

func (r requiresReplaceModifierComputed) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	resp.RequiresReplace = r.replace(req.State, req.Plan, req.StateValue, req.PlanValue)
}

func (r requiresReplaceModifierComputed) PlanModifyFloat64(_ context.Context, req planmodifier.Float64Request, resp *planmodifier.Float64Response) {
	resp.RequiresReplace = r.replace(req.State, req.Plan, req.StateValue, req.PlanValue)
}

func (r requiresReplaceModifierComputed) PlanModifyInt64(_ context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	resp.RequiresReplace = r.replace(req.State, req.Plan, req.StateValue, req.PlanValue)
}

func (r requiresReplaceModifierComputed) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	resp.RequiresReplace = r.replace(req.State, req.Plan, req.StateValue, req.PlanValue)
}

func (r requiresReplaceModifierComputed) PlanModifyMap(_ context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	resp.RequiresReplace = r.replace(req.State, req.Plan, req.StateValue, req.PlanValue)
}

func (r requiresReplaceModifierComputed) PlanModifyNumber(_ context.Context, req planmodifier.NumberRequest, resp *planmodifier.NumberResponse) {
	resp.RequiresReplace = r.replace(req.State, req.Plan, req.StateValue, req.PlanValue)
}

func (r requiresReplaceModifierComputed) PlanModifyObject(_ context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	resp.RequiresReplace = r.replace(req.State, req.Plan, req.StateValue, req.PlanValue)
}

func (r requiresReplaceModifierComputed) PlanModifySet(_ context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	resp.RequiresReplace = r.replace(req.State, req.Plan, req.StateValue, req.PlanValue)
}

func (r requiresReplaceModifierComputed) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	resp.RequiresReplace = r.replace(req.State, req.Plan, req.StateValue, req.PlanValue)
}

// Modify fills the AttributePlanModifier interface.
// func (r requiresReplaceModifierComputed) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
// 	if req.AttributeConfig == nil || req.AttributePlan == nil || req.AttributeState == nil {
// 		// shouldn't happen, but let's not panic if it does
// 		return
// 	}
//
// 	if req.State.Raw.IsNull() {
// 		// if we're creating the resource, no need to delete and
// 		// recreate it
// 		return
// 	}
//
// 	if req.Plan.Raw.IsNull() {
// 		// if we're deleting the resource, no need to delete and
// 		// recreate it
// 		return
// 	}
//
// 	if req.AttributePlan.Equal(req.AttributeState) {
// 		// if the plan and the state are in agreement, this attribute
// 		// isn't changing, don't require replace
// 		return
// 	}
//
// 	resp.RequiresReplace = true
// }

func (r requiresReplaceModifierComputed) replace(reqState tfsdk.State, reqPlan tfsdk.Plan, reqAttrState attr.Value, reqAttrPlan attr.Value) bool {
	if reqState.Raw.IsNull() {
		// if we're creating the resource, no need to delete and
		// recreate it
		return false
	}

	if reqPlan.Raw.IsNull() {
		// if we're deleting the resource, no need to delete and
		// recreate it
		return false
	}

	// if the plan and the state are in agreement, this attribute
	// isn't changing, don't require replace
	return !reqAttrPlan.Equal(reqAttrState)
}

// Description returns a human-readable description of the plan modifier.
func (r requiresReplaceModifierComputed) Description(ctx context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (r requiresReplaceModifierComputed) MarkdownDescription(ctx context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}
