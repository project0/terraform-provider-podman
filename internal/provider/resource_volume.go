package provider

import (
	"context"

	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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

// Schema returns the resource schema.
func (r volumeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage volumes for containers and pods",
		Attributes: withGenericAttributes(
			map[string]schema.Attribute{
				"driver": schema.StringAttribute{
					MarkdownDescription: "Name of the volume driver. Defaults by podman to `local`.",
					Required:            false,
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						modifier.RequiresReplaceComputed(),
					},
				},
				"options": schema.MapAttribute{
					Description: "Driver specific options.",
					Required:    false,
					Optional:    true,
					Computed:    true,
					ElementType: types.StringType,
					PlanModifiers: []planmodifier.Map{
						modifier.UseDefaultModifier(utils.MapStringEmpty()),
						modifier.RequiresReplaceComputed(),
					},
				},
			},
		),
	}
}

func fromVolumeResponse(v *entities.VolumeConfigResponse, diags *diag.Diagnostics) *volumeResourceData {
	return &volumeResourceData{
		// volumes do not have IDs, it wilbe mapped to the unique name
		ID:      types.StringValue(v.Name),
		Name:    types.StringValue(v.Name),
		Driver:  types.StringValue(v.Driver),
		Labels:  utils.MapStringToMapType(v.Labels, diags),
		Options: utils.MapStringToMapType(v.Options, diags),
	}
}
