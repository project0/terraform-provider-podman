package provider

import (
	"context"
	"fmt"
	"os"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	apiclient "github.com/project0/terraform-provider-podman/api/client"
	"github.com/project0/terraform-provider-podman/api/client/system_compat"
	"github.com/project0/terraform-provider-podman/internal/utils"
)

const (
	podmanDefaultURI = "unix:///run/podman/podman.sock"
)

// provider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type provider struct {
	data providerData

	// client is the configured connection context for the
	client *apiclient.Podman

	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	URI      types.String `tfsdk:"uri"`
	Identity types.String `tfsdk:"identity"`
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
		},
	}, nil
}

// Client initializes a new podman connection for further usage
func (p *provider) Client(ctx context.Context, diags *diag.Diagnostics) *apiclient.Podman {
	// set default to local socket
	uri := podmanDefaultURI

	// only used for tests
	if testuri := os.Getenv("TF_ACC_TEST_PROVIDER_PODMAN_URI"); testuri != "" {
		uri = testuri
	}

	if p.data.URI.Value != "" {
		uri = p.data.URI.Value
	}

	transport := httptransport.New(uri, "", nil)
	cl := apiclient.New(transport, strfmt.Default)

	// Probe client
	resp, err := cl.SystemCompat.SystemPing(system_compat.NewSystemPingParams(), nil)
	if err != nil {
		diags.AddError("Failed to initialize connection to podman server", fmt.Sprintf("URI: %s, error: %s", uri, err.Error()))
		return nil
	}
	tflog.Info(ctx, "Podman service successfully pinged",
		map[string]interface{}{
			"api_version":         resp.APIVersion,
			"lib_pod_api_version": resp.LibpodAPIVersion,
		},
	)

	return cl

}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	resp.Diagnostics.Append(req.Config.Get(ctx, &p.data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	p.client = p.Client(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	p.configured = true
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"podman_network": networkResourceType{},
		"podman_pod":     podResourceType{},
		"podman_volume":  volumeResourceType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		utils.AddUnexpectedError(
			&diags,
			"Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received.", p),
		)

		return provider{}, diags
	}

	if p == nil {
		utils.AddUnexpectedError(
			&diags,
			"Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received.",
		)
		return provider{}, diags
	}

	return *p, diags
}
