package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/project0/terraform-provider-podman/internal/modifier"
	"github.com/project0/terraform-provider-podman/internal/utils"
	"github.com/project0/terraform-provider-podman/internal/validators"
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

// withGenericAttributes returns re-usable standard type definitions
func withGenericAttributes(attributes map[string]schema.Attribute) map[string]schema.Attribute {
	// Name is also used as unique id in podman,
	// IDs itself only exists for docker compatibility and therefore does not make sense to implement
	attributes["name"] = schema.StringAttribute{
		Description: "Name of the resource, also used as ID. If not given a name will be automatically assigned.",
		Required:    false,
		Optional:    true,
		Computed:    true,
		Validators: []validator.String{
			validators.MatchName(),
		},
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
		},
	}

	attributes["labels"] = schema.MapAttribute{
		Description: "Labels is a set of user defined key-value labels of the resource",
		Required:    false,
		Optional:    true,
		Computed:    true,
		ElementType: types.StringType,
		PlanModifiers: []planmodifier.Map{
			modifier.UseDefaultModifier(utils.MapStringEmpty()),
			mapplanmodifier.UseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
		},
	}

	attributes["id"] = schema.StringAttribute{
		Description: "ID of the resource",
		Computed:    true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
			// Podman (go bindingds only?) looks up resource by name or id so it may not trigger replace.
			modifier.RequiresReplaceComputed(),
		},
	}

	return attributes
}
