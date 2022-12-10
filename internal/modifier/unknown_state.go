package modifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func AlwaysUseStateForUnknown() AlwaysUseStateForUnknownModifier {
	return AlwaysUseStateForUnknownModifier{}
}

type AlwaysUseStateForUnknownModifier struct{}

func (m AlwaysUseStateForUnknownModifier) Description(_ context.Context) string {
	return ""
}

func (m AlwaysUseStateForUnknownModifier) MarkdownDescription(_ context.Context) string {
	return ""
}

func (m AlwaysUseStateForUnknownModifier) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if m.setState(req.State, req.Plan, resp.PlanValue) {
		resp.PlanValue = req.StateValue
	}
}

func (m AlwaysUseStateForUnknownModifier) PlanModifyFloat64(_ context.Context, req planmodifier.Float64Request, resp *planmodifier.Float64Response) {
	if m.setState(req.State, req.Plan, resp.PlanValue) {
		resp.PlanValue = req.StateValue
	}
}

func (m AlwaysUseStateForUnknownModifier) PlanModifyInt64(_ context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if m.setState(req.State, req.Plan, resp.PlanValue) {
		resp.PlanValue = req.StateValue
	}
}

func (m AlwaysUseStateForUnknownModifier) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if m.setState(req.State, req.Plan, resp.PlanValue) {
		resp.PlanValue = req.StateValue
	}
}

func (m AlwaysUseStateForUnknownModifier) PlanModifyMap(_ context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	if m.setState(req.State, req.Plan, resp.PlanValue) {
		resp.PlanValue = req.StateValue
	}
}

func (m AlwaysUseStateForUnknownModifier) PlanModifyNumber(_ context.Context, req planmodifier.NumberRequest, resp *planmodifier.NumberResponse) {
	if m.setState(req.State, req.Plan, resp.PlanValue) {
		resp.PlanValue = req.StateValue
	}
}

func (m AlwaysUseStateForUnknownModifier) PlanModifyObject(_ context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if m.setState(req.State, req.Plan, resp.PlanValue) {
		resp.PlanValue = req.StateValue
	}
}

func (m AlwaysUseStateForUnknownModifier) PlanModifySet(_ context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if m.setState(req.State, req.Plan, resp.PlanValue) {
		resp.PlanValue = req.StateValue
	}
}

func (m AlwaysUseStateForUnknownModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if m.setState(req.State, req.Plan, resp.PlanValue) {
		resp.PlanValue = req.StateValue
	}
}

//func (m AlwaysUseStateForUnknownModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
//	if req.AttributeState == nil || resp.AttributePlan == nil || req.AttributeConfig == nil {
//		return
//	}
//
//	// if we're creating the resource, no need to modify
//	if req.State.Raw.IsNull() {
//		return
//	}
//
//	// if we're deleting the resource, no need to modify
//	if req.Plan.Raw.IsNull() {
//		return
//	}
//
//	// if it's not planned to be the unknown value, stick with the concrete plan
//	if !resp.AttributePlan.IsUnknown() {
//		return
//	}
//
//	resp.AttributePlan = req.AttributeState
//}

func (m AlwaysUseStateForUnknownModifier) setState(reqState tfsdk.State, reqPlan tfsdk.Plan, respAttrPlan attr.Value) bool {
	// if we're creating the resource, no need to modify
	if reqState.Raw.IsNull() {
		return false
	}

	// if we're deleting the resource, no need to modify
	if reqPlan.Raw.IsNull() {
		return false
	}

	// if it's not planned to be the unknown value, stick with the concrete plan
	return respAttrPlan.IsUnknown()
}
