package provider

import (
	"context"
	"fmt"

	"github.com/containers/podman/v4/pkg/bindings/volumes"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/project0/terraform-provider-podman/internal/utils"
)

func (r volumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data volumeResourceData

	client := r.initClientData(ctx, &data, req.Config.Get, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build volume
	var volCreate = &entities.VolumeCreateOptions{
		// null values are automatically empty string
		// as we do not use pointer we do not need to distinguish and pass it directly
		Name:   data.Name.ValueString(),
		Driver: data.Driver.ValueString(),
	}

	resp.Diagnostics.Append(data.Labels.ElementsAs(ctx, &volCreate.Labels, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// optional
	if !data.Options.IsNull() {
		resp.Diagnostics.Append(data.Options.ElementsAs(ctx, &volCreate.Options, true)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create
	volResponse, err := volumes.Create(client, *volCreate, nil)
	if err != nil {
		resp.Diagnostics.AddError("Podman client error", fmt.Sprintf("Failed to create volume resource: %s", err.Error()))
		return
	}

	// Set state
	resp.Diagnostics.Append(
		resp.State.Set(ctx, fromVolumeResponse(volResponse, &resp.Diagnostics))...,
	)
}

func (r volumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data volumeResourceData

	client := r.initClientData(ctx, &data, req.State.Get, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	volResponse, err := volumes.Inspect(client, data.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Podman client error", fmt.Sprintf("Failed to read volume resource: %s", err.Error()))
		return
	}

	// Set state
	resp.Diagnostics.Append(
		resp.State.Set(ctx, fromVolumeResponse(volResponse, &resp.Diagnostics))...,
	)
}

// Update is not implemented
func (r volumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	utils.AddUnexpectedError(
		&resp.Diagnostics,
		"Update triggered for a volume resource",
		"Volumes are immutable resources and cannot be updated, it always needs to be replaced.",
	)
}

func (r volumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data volumeResourceData

	client := r.initClientData(ctx, &data, req.State.Get, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Allow force ?
	err := volumes.Remove(client, data.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Podman client error", fmt.Sprintf("Failed to delete volume resource: %s", err.Error()))
	}

	resp.State.RemoveResource(ctx)
}

func (r volumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
