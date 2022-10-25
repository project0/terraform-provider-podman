package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	podmanDefaultURI = "unix:///run/podman/podman.sock"
)

// provider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type podmanProvider struct {
	// client is the configured connection context for the
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	URI           types.String `tfsdk:"uri"`
	Identity      types.String `tfsdk:"identity"`
	DefaultLabels types.Map    `tfsdk:"default_labels"`
}

func New() provider.Provider {
	return &podmanProvider{}
}

// Metadata returns the provider type name.
func (p *podmanProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "podman"
}

// GetSchema defines the provider-level schema for configuration data.
func (p *podmanProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "The Podman provider provides resource management via the remote API.",
		Attributes: map[string]tfsdk.Attribute{
			"uri": {
				Description: "Connection URI to the podman service.",
				MarkdownDescription: "Connection URI to the podman service. " +
					"A valid URI connection should be of `scheme://`. " +
					"For example `tcp://localhost:<port>`" +
					"or `unix:///run/podman/podman.sock`" +
					"or `ssh://<user>@<host>[:port]/run/podman/podman.sock?secure=True`." +
					"Defaults to `" + podmanDefaultURI + "`.",
				Optional: true,
				Type:     types.StringType,
			},
			"identity": {
				Description: "Local path to the identity file for SSH based connections.",
				Optional:    true,
				Type:        types.StringType,
			},
			"default_labels": {
				Description: "Default labels for all resources. Please note that changing labels will trigger replacement for most of the resources.",
				Optional:    true,
				Type: types.MapType{
					ElemType: types.StringType,
				},
			},
		},
	}, nil
}

// Configure prepares a HashiCups API client for data sources and resources.
func (p *podmanProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Debug(ctx, "Configure Podman client")

	var data providerData
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newPodmanClient(ctx, &resp.Diagnostics, data)
	if resp.Diagnostics.HasError() {
		return
	}

	// make podman clent data available
	resp.DataSourceData = data
	resp.ResourceData = data

	tflog.Info(ctx, "Configured Podman client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *podmanProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *podmanProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewNetworkResource,
		NewPodResource,
		NewVolumeResource,
	}
}

// newPodmanClient initializes a new podman connection for further usage
// The final client is the configured connection context
func newPodmanClient(ctx context.Context, diags *diag.Diagnostics, data providerData) context.Context {
	// set default to local socket
	uri := podmanDefaultURI

	// only used for tests
	if testuri := os.Getenv("TF_ACC_TEST_PROVIDER_PODMAN_URI"); testuri != "" {
		uri = testuri
	}

	if data.URI.Value != "" {
		uri = data.URI.Value
	}

	c, err := bindings.NewConnectionWithIdentity(ctx, uri, data.Identity.Value, false)
	if err != nil {
		diags.AddError("Failed to initialize connection to podman server", fmt.Sprintf("URI: %s, error: %s", uri, err.Error()))
	}

	return c
}
