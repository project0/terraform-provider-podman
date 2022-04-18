package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/project0/terraform-provider-podman/internal/modifier"
	"github.com/project0/terraform-provider-podman/internal/validator"
)

type (
	genericResource struct {
		provider provider
	}
)

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

	return g.provider.Client(ctx, diags)
}

// re-usable type definitions
func withGenericAttributes(attributes map[string]tfsdk.Attribute) map[string]tfsdk.Attribute {
	// Name is also used as unique id in podman,
	// IDs itself only exists for docker compatibility and therefore does not make sense to implement
	attributes["name"] = tfsdk.Attribute{
		MarkdownDescription: "Name of the resource, also used as ID. If not given a name will be automatically assigned.",
		Required:            false,
		Optional:            true,
		Computed:            true,
		Validators:          []tfsdk.AttributeValidator{validator.MatchName()},
		Type:                types.StringType,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			tfsdk.RequiresReplace(),
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
		PlanModifiers: tfsdk.AttributePlanModifiers{
			modifier.UseDefaultModifier(types.Map{ElemType: types.StringType, Null: false}),
			tfsdk.RequiresReplace(),
		},
	}

	return attributes
}
