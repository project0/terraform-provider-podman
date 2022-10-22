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

func NewVolumeResource() resource.Resource {
	return &volumeResource{}
}

// Metadata returns the data source type name.
func (r volumeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

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
					resource.RequiresReplace(),
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
					resource.RequiresReplace(),
				},
			},
		}),

		Blocks: map[string]tfsdk.Block{},
	}, nil
}

// Configure adds the provider configured client to the data source.
func (r *volumeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.genericResource.Configure(ctx, req, resp)
}

func fromVolumeResponse(v *entities.VolumeConfigResponse, diags *diag.Diagnostics) *volumeResourceData {
	return &volumeResourceData{
		ID:      types.String{Value: v.Name},
		Name:    types.String{Value: v.Name},
		Driver:  types.String{Value: v.Driver},
		Labels:  utils.MapStringToMapType(v.Labels),
		Options: utils.MapStringToMapType(v.Options),
	}
}
