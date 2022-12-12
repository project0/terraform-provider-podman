package modifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UseDefaultModifier returns a new DefaultModifier
func UseDefaultModifier(value attr.Value) DefaultModifier {
	return DefaultModifier{
		value: value,
	}
}

// DefaultModifier ensures a new value with null value will always be set to default
type DefaultModifier struct {
	value attr.Value
}

func (m DefaultModifier) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.PlanValue.IsNull() {
		resp.PlanValue = m.value.(types.Bool)
	}
}

func (m DefaultModifier) PlanModifyFloat64(_ context.Context, req planmodifier.Float64Request, resp *planmodifier.Float64Response) {
	if req.PlanValue.IsNull() {
		resp.PlanValue = m.value.(types.Float64)
	}
}

func (m DefaultModifier) PlanModifyInt64(_ context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if req.PlanValue.IsNull() {
		resp.PlanValue = m.value.(types.Int64)
	}
}

func (m DefaultModifier) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.PlanValue.IsNull() {
		resp.PlanValue = m.value.(types.List)
	}
}
func (m DefaultModifier) PlanModifyMap(_ context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	if req.PlanValue.IsNull() {
		resp.PlanValue = m.value.(types.Map)
	}
}

func (m DefaultModifier) PlanModifyNumber(_ context.Context, req planmodifier.NumberRequest, resp *planmodifier.NumberResponse) {
	if req.PlanValue.IsNull() {
		resp.PlanValue = m.value.(types.Number)
	}
}

func (m DefaultModifier) PlanModifyObject(_ context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if req.PlanValue.IsNull() {
		resp.PlanValue = m.value.(types.Object)
	}
}

func (m DefaultModifier) PlanModifySet(_ context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if req.PlanValue.IsNull() {
		resp.PlanValue = m.value.(types.Set)
	}
}

func (m DefaultModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsNull() {
		resp.PlanValue = m.value.(types.String)
	}
}

// Modify sets the default value if not known or empty
//func (m DefaultModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
//	val, err := req.AttributePlan.ToTerraformValue(ctx)
//	if err != nil {
//		utils.AddUnexpectedAttributeError(meq.AttributePath, &resp.Diagnostics, "Failed to retrieve value", err.Error())
//		return
//	}
//
//	if val.IsNull() {
//		resp.AttributePlan = r.value
//	}
//}

// Description returns a human-readable description of the plan modifier.
func (m DefaultModifier) Description(_ context.Context) string {
	return "Ensure null values are replaced by given default attribute."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m DefaultModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}
