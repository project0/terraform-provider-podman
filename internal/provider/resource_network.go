package provider

import (
	"context"
	"fmt"
	"net"

	ntypes "github.com/containers/common/libnetwork/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/project0/terraform-provider-podman/internal/modifier"
	"github.com/project0/terraform-provider-podman/internal/utils"
	"github.com/project0/terraform-provider-podman/internal/validators"
)

type (
	networkResource struct {
		genericResource
	}

	networkResourceData struct {
		ID     types.String `tfsdk:"id"`
		Name   types.String `tfsdk:"name"`
		Labels types.Map    `tfsdk:"labels"`

		DNS      types.Bool `tfsdk:"dns"`
		IPv6     types.Bool `tfsdk:"ipv6"`
		Internal types.Bool `tfsdk:"internal"`

		Driver     types.String `tfsdk:"driver"`
		IPAMDriver types.String `tfsdk:"ipam_driver"`
		Options    types.Map    `tfsdk:"options"`

		Subnets []networkResourceSubnetData `tfsdk:"subnets"`
	}

	networkResourceSubnetData struct {
		Subnet  types.String `tfsdk:"subnet"`
		Gateway types.String `tfsdk:"gateway"`
	}
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &networkResource{}
	_ resource.ResourceWithConfigure   = &networkResource{}
	_ resource.ResourceWithImportState = &networkResource{}
)

// NewNetworkResource creates a new network resource.
func NewNetworkResource() resource.Resource {
	return &networkResource{}
}

// Configure adds the provider configured client to the resource.
func (r *networkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.genericResource.Configure(ctx, req, resp)
}

// Metadata returns the resource type name.
func (r networkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

// Schema returns the resource schema.
func (r networkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage networks for containers and pods",
		Attributes: withGenericAttributes(
			map[string]schema.Attribute{
				"dns": schema.BoolAttribute{
					MarkdownDescription: "Enable the DNS plugin for this network which if enabled, can perform container to container name resolution. Defaults to `false`.",
					Computed:            true,
					Optional:            true,
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.UseStateForUnknown(),
						modifier.RequiresReplaceComputed(),
					},
				},

				"ipv6": schema.BoolAttribute{
					MarkdownDescription: "Enable IPv6 (Dual Stack) networking. If no subnets are given it will allocate a ipv4 and ipv6 subnet. Defaults to `false`.",
					Computed:            true,
					Optional:            true,
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.UseStateForUnknown(),
						modifier.RequiresReplaceComputed(),
					},
				},

				"internal": schema.BoolAttribute{
					MarkdownDescription: "Internal is whether the Network should not have external routes to public or other Networks. Defaults to `false`.",
					Computed:            true,
					Optional:            true,
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.UseStateForUnknown(),
						modifier.RequiresReplaceComputed(),
					},
				},

				"driver": schema.StringAttribute{
					MarkdownDescription: fmt.Sprintf(
						"Driver to manage the network. One of `%s`, `%s`, `%s` are currently supported. By podman defaults to `bridge`.",
						ntypes.BridgeNetworkDriver,
						ntypes.MacVLANNetworkDriver,
						ntypes.IPVLANNetworkDriver,
					),
					Computed: true,
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf(
							ntypes.BridgeNetworkDriver,
							ntypes.MacVLANNetworkDriver,
							ntypes.IPVLANNetworkDriver,
						),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
						modifier.RequiresReplaceComputed(),
					},
				},

				"ipam_driver": schema.StringAttribute{
					Computed: true,
					Optional: true,
					MarkdownDescription: fmt.Sprintf(
						"Set the ipam driver (IP Address Management Driver) for the network. Valid values are `%s`, `%s`, `%s`. When unset podman will choose an ipam driver automatically based on the network driver.",
						ntypes.HostLocalIPAMDriver,
						ntypes.DHCPIPAMDriver,
						"none",
					),
					Validators: []validator.String{
						stringvalidator.OneOf(
							ntypes.HostLocalIPAMDriver,
							ntypes.DHCPIPAMDriver,
							"none",
						),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
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

				"subnets": schema.SetNestedAttribute{
					Description: "Subnets for this network.",
					Required:    false,
					Optional:    true,
					Computed:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"subnet": schema.StringAttribute{
								MarkdownDescription: "The subnet in CIDR notation.",
								Required:            true,
								Optional:            false,
								Validators: []validator.String{
									validators.IsCIDR(),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
							"gateway": schema.StringAttribute{
								MarkdownDescription: "Gateway IP for this Network.",
								Computed:            true,
								Optional:            true,
								Validators: []validator.String{
									validators.IsIpAdress(),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
									modifier.RequiresReplaceComputed(),
								},
							},
						},
					},
				},
			},
		),
	}
}

// toPodmanNetwork converts a resource data to a podman network
func toPodmanNetwork(ctx context.Context, d networkResourceData, diags *diag.Diagnostics) *ntypes.Network {
	var nw = &ntypes.Network{
		Name:        d.Name.ValueString(),
		Driver:      d.Driver.ValueString(),
		IPv6Enabled: d.IPv6.ValueBool(),
		DNSEnabled:  d.DNS.ValueBool(),
		Internal:    d.Internal.ValueBool(),
	}

	// Convert map types
	diags.Append(d.Labels.ElementsAs(ctx, &nw.Labels, true)...)
	diags.Append(d.Options.ElementsAs(ctx, &nw.Options, true)...)

	if !d.IPAMDriver.IsNull() {
		ipam := map[string]string{
			"driver": d.IPAMDriver.ValueString(),
		}
		nw.IPAMOptions = ipam
	}

	// subnet
	for _, s := range d.Subnets {
		_, ipNet, err := net.ParseCIDR(s.Subnet.ValueString())
		if err != nil {
			diags.AddError("Cannot parse subnet CIDR", err.Error())
			continue
		}
		subnet := &ntypes.Subnet{
			Subnet: ntypes.IPNet{
				IPNet: *ipNet,
			},
		}
		if !s.Gateway.IsNull() {
			subnet.Gateway = net.ParseIP(s.Gateway.ValueString())
		}
		nw.Subnets = append(nw.Subnets, *subnet)
	}

	return nw
}

// fromNetwork converts a podman network to a resource data
func fromPodmanNetwork(n ntypes.Network, diags *diag.Diagnostics) *networkResourceData {
	d := &networkResourceData{
		ID:       types.StringValue(n.Name),
		Name:     types.StringValue(n.Name),
		DNS:      types.BoolValue(n.DNSEnabled),
		IPv6:     types.BoolValue(n.IPv6Enabled),
		Internal: types.BoolValue(n.Internal),
		Driver:   types.StringValue(n.Driver),
		Labels:   utils.MapStringToMapType(n.Labels, diags),
		Options:  utils.MapStringToMapType(n.Options, diags),
	}

	d.IPAMDriver = utils.MapStringValueToStringType(n.IPAMOptions, "driver")

	for _, s := range n.Subnets {
		subnet := networkResourceSubnetData{
			Subnet:  types.StringValue(s.Subnet.String()),
			Gateway: types.StringValue(s.Gateway.String()),
		}
		d.Subnets = append(d.Subnets, subnet)
	}
	return d
}
