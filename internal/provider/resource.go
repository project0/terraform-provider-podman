package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/project0/terraform-provider-podman/internal/modifier"
	"github.com/project0/terraform-provider-podman/internal/utils"
	"github.com/project0/terraform-provider-podman/internal/validator"
)

type (
	genericResource struct {
		providerData providerData
	}
)

// Configures the podman client
func (g *genericResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	// It looks like this method is called twice, the first time its nil and happen before the provider is initialized.
	if req.ProviderData == nil {
		return
	}

	var ok bool
	g.providerData, ok = req.ProviderData.(providerData)

	if !ok {
		utils.AddUnexpectedError(
			&resp.Diagnostics,
			"Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received.", req.ProviderData),
		)
	}
}

func (g *genericResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// No default labels set, skip merge
	if g.providerData.DefaultLabels.IsNull() {
		return
	}

	var labels types.Map
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("labels"), &labels)...)
	if resp.Diagnostics.HasError() {
		return
	}

	labels.Null = false

	for k, v := range g.providerData.DefaultLabels.Elems {
		if _, exist := labels.Elems[k]; exist {
			continue
		}
		labels.Elems[k] = v
	}

	resp.Plan.SetAttribute(ctx, path.Root("labels"), &labels)
}

func (g genericResource) initClientData(
	ctx context.Context,
	data interface{},
	get func(context.Context, interface{}) diag.Diagnostics,
	diags *diag.Diagnostics,
) context.Context {

	diags.Append(
		get(ctx, data)...,
	)

	if diags.HasError() {
		tflog.Error(ctx, "Failed to retrieve resource data")
		return nil
	}

	return newPodmanClient(ctx, diags, g.providerData)
}

// re-usable type definitions
func withGenericAttributes(attributes map[string]tfsdk.Attribute) map[string]tfsdk.Attribute {
	// Name is also used as unique id in podman,
	// IDs itself only exists for docker compatibility and therefore does not make sense to implement
	attributes["name"] = tfsdk.Attribute{
		Description: "Name of the resource, also used as ID. If not given a name will be automatically assigned.",
		Required:    false,
		Optional:    true,
		Computed:    true,
		Validators:  []tfsdk.AttributeValidator{validator.MatchName()},
		Type:        types.StringType,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			resource.UseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
		},
	}

	attributes["labels"] = tfsdk.Attribute{
		Description: "Labels is a set of user defined key-value labels of the resource",
		Required:    false,
		Optional:    true,
		Computed:    true,
		Type: types.MapType{
			ElemType: types.StringType,
		},
		//		PlanModifiers: tfsdk.AttributePlanModifiers{
		//			modifier.UseDefaultModifier(types.Map{ElemType: types.StringType, Null: false}),
		//			resource.UseStateForUnknown(),
		//			resource.RequiresReplace(),
		//		},
	}

	attributes["id"] = tfsdk.Attribute{
		Description: "ID of the resource",
		Type:        types.StringType,
		Computed:    true,
		PlanModifiers: []tfsdk.AttributePlanModifier{
			resource.UseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
			// note: we may need a custom version of RequiresReplace as it does not support replace on computed attributes.
			// Podman (go bindingds only?) looks up resource by name or id so it may not trigger replace.
			// resource.RequiresReplace(),
		},
	}

	return attributes
}
