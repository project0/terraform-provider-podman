package provider

import (
	"context"

	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/project0/terraform-provider-podman/internal/modifier"
	"github.com/project0/terraform-provider-podman/internal/utils"
)

type (
	volumeResource struct {
		genericResource
	}
	volumeResourceData struct {
		ID     types.String `tfsdk:"id"`
		Name   types.String `tfsdk:"name"`
		Labels types.Map    `tfsdk:"labels"`

		Driver  types.String `tfsdk:"driver"`
		Options types.Map    `tfsdk:"options"`
	}
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &volumeResource{}
	_ resource.ResourceWithConfigure   = &volumeResource{}
	_ resource.ResourceWithImportState = &volumeResource{}
)

// NewVolumeResource creates a new volume resource.
func NewVolumeResource() resource.Resource {
	return &volumeResource{}
}

// Configure adds the provider configured client to the resource.
func (r *volumeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.genericResource.Configure(ctx, req, resp)
}

// Metadata returns the resource type name.
func (r volumeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

// GetSchema returns the resource schema.
func (t volumeResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Manage volumes for containers and pods",
		Attributes: withGenericAttributes(map[string]tfsdk.Attribute{
			"driver": {
				MarkdownDescription: "Name of the volume driver. Defaults by podman to `local`.",
				Required:            false,
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifier.RequiresReplaceComputed(),
				},
			},
			"options": {
				Description: "Driver specific options.",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Type: types.MapType{
					ElemType: types.StringType,
				},
				PlanModifiers: tfsdk.AttributePlanModifiers{
					modifier.UseDefaultModifier(types.Map{ElemType: types.StringType, Null: false}),
					modifier.RequiresReplaceComputed(),
				},
			},
		}),

		Blocks: map[string]tfsdk.Block{},
	}, nil
}

func fromVolumeResponse(v *entities.VolumeConfigResponse, diags *diag.Diagnostics) *volumeResourceData {
	return &volumeResourceData{
		// volumes do not have IDs, it wilbe mapped to the unique name
		ID:      types.String{Value: v.Name},
		Name:    types.String{Value: v.Name},
		Driver:  types.String{Value: v.Driver},
		Labels:  utils.MapStringToMapType(v.Labels),
		Options: utils.MapStringToMapType(v.Options),
	}
}
