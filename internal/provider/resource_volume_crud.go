package provider

import (
	"context"
	"fmt"

	"github.com/project0/terraform-provider-podman/api/client/volumes"
	"github.com/project0/terraform-provider-podman/api/models"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/project0/terraform-provider-podman/internal/utils"
)

func (r volumeResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data volumeResourceData

	client := r.initClientData(ctx, &data, req.Config.Get, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build volume
	var volCreate = &models.VolumeCreateOptions{
		// null values are automatically empty string
		// as we do not use pointer we do not need to distinguish and pass it directly
		Name:   data.Name.Value,
		Driver: data.Driver.Value,
	}

	resp.Diagnostics.Append(data.Labels.ElementsAs(ctx, &volCreate.Labels, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// optional
	if !data.Options.Null {
		resp.Diagnostics.Append(data.Options.ElementsAs(ctx, &volCreate.Options, true)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create
	volResponse, err := client.Volumes.VolumeCreateLibpod(
		volumes.NewVolumeCreateLibpodParamsWithContext(ctx).
			WithCreate(volCreate),
	)
	if err != nil {
		resp.Diagnostics.AddError("Podman client error", fmt.Sprintf("Failed to create volume resource: %s", err.Error()))
		return
	}

	// Set state
	resp.Diagnostics.Append(
		resp.State.Set(ctx, fromVolumeResponse(volResponse.GetPayload(), &resp.Diagnostics))...,
	)
}

func (r volumeResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data volumeResourceData

	client := r.initClientData(ctx, &data, req.State.Get, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	volResponse, err := client.Volumes.VolumeInspectLibpod(
		volumes.NewVolumeInspectLibpodParamsWithContext(ctx).
			WithName(data.Name.Value),
	)
	if err != nil {
		resp.Diagnostics.AddError("Podman client error", fmt.Sprintf("Failed to read volume resource: %s", err.Error()))
		return
	}

	// Set state
	resp.Diagnostics.Append(
		resp.State.Set(ctx, fromVolumeResponse(volResponse.GetPayload(), &resp.Diagnostics))...,
	)
}

// Update is not implemented
func (r volumeResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	utils.AddUnexpectedError(
		&resp.Diagnostics,
		"Update triggered for a volume resource",
		"Volumes are immutable resources and cannot be updated, it always needs to be replaced.",
	)
}

func (r volumeResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data volumeResourceData

	client := r.initClientData(ctx, &data, req.State.Get, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Allow force ?
	err := volumes.Remove(client, data.Name.Value, nil)
	if err != nil {
		resp.Diagnostics.AddError("Podman client error", fmt.Sprintf("Failed to delete volume resource: %s", err.Error()))
	}

	resp.State.RemoveResource(ctx)
}

func (r volumeResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("name"), req, resp)
}
