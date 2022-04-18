package provider

import (
	"context"

	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/project0/terraform-provider-podman/internal/modifier"
	"github.com/project0/terraform-provider-podman/internal/utils"
)

type (
	volumeResource struct {
		genericResource
	}
	volumeResourceType struct{}
	volumeResourceData struct {
		Name   types.String `tfsdk:"name"`
		Labels types.Map    `tfsdk:"labels"`

		Driver  types.String `tfsdk:"driver"`
		Options types.Map    `tfsdk:"options"`
	}
)

func (t volumeResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Volume",
		Attributes: withGenericAttributes(map[string]tfsdk.Attribute{
			"driver": {
				MarkdownDescription: "Name of the volume driver. Defaults by podman to `local`.",
				Required:            false,
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
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
					tfsdk.RequiresReplace(),
				},
			},
		}),

		Blocks: map[string]tfsdk.Block{},
	}, nil
}

func (t volumeResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return volumeResource{
		genericResource: genericResource{
			provider: provider,
		},
	}, diags
}

func fromVolumeResponse(v *entities.VolumeConfigResponse, diags *diag.Diagnostics) *volumeResourceData {
	return &volumeResourceData{
		Name:    types.String{Value: v.Name},
		Driver:  types.String{Value: v.Driver},
		Labels:  utils.MapStringToMapType(v.Labels),
		Options: utils.MapStringToMapType(v.Options),
	}
}
