package provider

import (
	"context"
	"fmt"

	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/project0/terraform-provider-podman/internal/modifier"
	"github.com/project0/terraform-provider-podman/internal/provider/shared"
	"github.com/project0/terraform-provider-podman/internal/utils"
)

type (
	podResource struct {
		genericResource
	}
	podResourceData struct {
		ID     types.String `tfsdk:"id"`
		Name   types.String `tfsdk:"name"`
		Labels types.Map    `tfsdk:"labels"`

		CgroupParent types.String `tfsdk:"cgroup_parent"`
		Hostname     types.String `tfsdk:"hostname"`

		Mounts shared.Mounts `tfsdk:"mounts"`
	}
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &podResource{}
	_ resource.ResourceWithSchema      = &podResource{}
	_ resource.ResourceWithConfigure   = &podResource{}
	_ resource.ResourceWithImportState = &podResource{}
)

// NewPodResource creates a new pod resource.
func NewPodResource() resource.Resource {
	return &podResource{}
}

// Configure adds the provider configured client to the resource.
func (r *podResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.genericResource.Configure(ctx, req, resp)
}

// Metadata returns the resource type name.
func (r podResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pod"
}

// Schema returns the resource schema.
func (r podResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	mountsAttr := make(shared.Mounts, 0)
	resp.Schema = schema.Schema{
		Description: "Manage pods for containers",
		Attributes: withGenericAttributes(
			map[string]schema.Attribute{
				"cgroup_parent": schema.StringAttribute{
					MarkdownDescription: "Path to cgroups under which the cgroup for the pod will be created. " +
						"If the path is not absolute, the path is considered to be relative to the cgroups path of the init process. " +
						"Cgroups will be created if they do not already exist.",
					Required: false,
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
						modifier.RequiresReplaceComputed(),
					},
				},
				"hostname": schema.StringAttribute{
					Description: "Hostname is the pod's hostname. " +
						"If not set, the name of the pod will be used (if a name was not provided here, the name auto-generated for the pod will be used). " +
						"This will be used by the infra container and all containers in the pod as long as the UTS namespace is shared.",
					Required: false,
					Optional: true,
					Computed: false,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
						stringplanmodifier.RequiresReplace(),
					},
				},
				"mounts": mountsAttr.GetSchema(ctx),
			},
		),
	}
}

func toPodmanPodSpecGenerator(ctx context.Context, d podResourceData, diags *diag.Diagnostics) *specgen.PodSpecGenerator {
	s := specgen.NewPodSpecGenerator()
	p := &entities.PodCreateOptions{
		Name:         d.Name.ValueString(),
		CgroupParent: d.CgroupParent.ValueString(),
		Hostname:     d.Hostname.ValueString(),
		Infra:        true,
	}

	diags.Append(d.Labels.ElementsAs(ctx, &p.Labels, true)...)
	sp, err := entities.ToPodSpecGen(*s, p)
	if err != nil {
		diags.AddError("Invalid pod configuration", fmt.Sprintf("Cannot build pod configuration: %q", err.Error()))
	}
	// add storage
	sp.Volumes, sp.Mounts = d.Mounts.ToPodmanSpec(diags)
	if err := sp.Validate(); err != nil {
		diags.AddError("Invalid pod configuration", fmt.Sprintf("Cannot build pod configuration: %q", err.Error()))
	}
	return sp
}

func fromPodResponse(p *entities.PodInspectReport, diags *diag.Diagnostics) *podResourceData {
	hostname := types.StringNull()
	if p.Hostname != "" {
		hostname = types.StringValue(p.Hostname)
	}

	d := &podResourceData{
		ID:           types.StringValue(p.ID),
		Name:         types.StringValue(p.Name),
		Labels:       utils.MapStringToMapType(p.Labels, diags),
		Mounts:       shared.FromPodmanToMounts(diags, p.Mounts),
		CgroupParent: types.StringValue(p.CgroupParent),
		Hostname:     hostname,
	}

	return d
}
