package provider

import (
	"context"
	"fmt"

	"github.com/containers/podman/v4/pkg/bindings/network"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/project0/terraform-provider-podman/internal/utils"
)

func (r networkResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data networkResourceData

	client := r.initClientData(ctx, &data, req.Config.Get, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	networkCreate := toPodmanNetwork(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	networkResponse, err := network.Create(client, networkCreate)
	if err != nil {
		resp.Diagnostics.AddError("Podman client error", fmt.Sprintf("Failed to create network resource: %s", err.Error()))
		return
	}

	diags := resp.State.Set(ctx, fromPodmanNetwork(networkResponse, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
}

func (r networkResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data networkResourceData

	client := r.initClientData(ctx, &data, req.State.Get, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	networkResponse, err := network.Inspect(client, data.Name.Value, nil)
	if err != nil {
		resp.Diagnostics.AddError("Podman client error", fmt.Sprintf("Failed to read network resource: %s", err.Error()))
		return
	}

	diags := resp.State.Set(ctx, fromPodmanNetwork(networkResponse, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
}

// Update is not implemented
func (r networkResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	utils.AddUnexpectedError(
		&resp.Diagnostics,
		"Update triggered for a network resource",
		"Networks are immutable resources and cannot be updated, it always needs to be replaced.",
	)
}

func (r networkResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data networkResourceData

	client := r.initClientData(ctx, &data, req.State.Get, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Allow force which detaches containers from network?
	rmErrors, err := network.Remove(client, data.Name.Value, nil)
	if err != nil {
		resp.Diagnostics.AddError("Podman client error", fmt.Sprintf("Failed to delete network resource: %s", err.Error()))
	}
	for _, e := range rmErrors {
		if e.Err != nil {
			resp.Diagnostics.AddError("Error report on deletion for "+e.Name, e.Err.Error())
		}
	}

	resp.State.RemoveResource(ctx)
}

func (r networkResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("name"), req, resp)
}
