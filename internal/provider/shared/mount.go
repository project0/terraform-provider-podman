package shared

import (
	"context"

	"github.com/containers/podman/v4/libpod/define"
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/project0/terraform-provider-podman/internal/modifier"
	"github.com/project0/terraform-provider-podman/internal/validators"

	"github.com/opencontainers/runtime-spec/specs-go"
)

type (
	Mounts []Mount

	Mount struct {
		Destination types.String `tfsdk:"destination"`

		Volume *MountVolume `tfsdk:"volume"`
		Bind   *MountBind   `tfsdk:"bind"`
		// Tmpfs  *MountTmpfs   `tfsdk:"tmpfs"`
	}

	// MountVolume mounts a named volume
	MountVolume struct {
		Name types.String `tfsdk:"name"`

		ReadOnly types.Bool `tfsdk:"read_only"`
		Chown    types.Bool `tfsdk:"chown"`

		Suid  types.Bool `tfsdk:"suid"`
		Exec  types.Bool `tfsdk:"exec"`
		Dev   types.Bool `tfsdk:"dev"`
		IDmap types.Bool `tfsdk:"idmap"`
	}

	// MountBind mounts host path
	MountBind struct {
		Path types.String `tfsdk:"path"`

		ReadOnly types.Bool `tfsdk:"read_only"`
		Chown    types.Bool `tfsdk:"chown"`

		Suid  types.Bool `tfsdk:"suid"`
		Exec  types.Bool `tfsdk:"exec"`
		Dev   types.Bool `tfsdk:"dev"`
		IDmap types.Bool `tfsdk:"idmap"`

		Propagation types.String `tfsdk:"propagation"`
		Recursive   types.Bool   `tfsdk:"recursive"`
		Relabel     types.Bool   `tfsdk:"relabel"`
	}

	// MountTmpfs mounts a tmpfs
	MountTmpfs struct {
		ReadOnly types.Bool `tfsdk:"read_only"`
		Chown    types.Bool `tfsdk:"chown"`

		Suid types.Bool `tfsdk:"suid"`
		Exec types.Bool `tfsdk:"exec"`
		Dev  types.Bool `tfsdk:"dev"`

		Size      types.String `tfsdk:"size"`
		Mode      types.String `tfsdk:"mode"`
		TmpCopyUp types.Bool   `tfsdk:"tmpcopyup"`
	}
)

func (m Mounts) GetSchema(ctx context.Context) schema.Attribute {
	return schema.SetNestedAttribute{
		Description: "Mounts volume, bind, image, tmpfs, etc..",
		Required:    false,
		Optional:    true,
		Computed:    true,
		PlanModifiers: []planmodifier.Set{
			modifier.RequiresReplaceComputed(),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"destination": schema.StringAttribute{
					Description: "Target path",
					Required:    true,
					Computed:    false,
					Validators:  []validator.String{
						//TODO
					},
				},
				"volume": schema.SingleNestedAttribute{
					Description: "Named Volume",
					Optional:    true,
					Computed:    false,
					Validators: []validator.Object{
						objectvalidator.ConflictsWith(
							path.MatchRelative().AtParent().AtName("bind"),
						),
					},
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.RequiresReplace(),
					},
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the volume",
							Required:    true,
							Computed:    false,
							Validators: []validator.String{
								validators.MatchName(),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"read_only": m.attributeSchemaReadOnly(),
						"dev":       m.attributeSchemaDev(),
						"exec":      m.attributeSchemaExec(),
						"suid":      m.attributeSchemaSuid(),
						"chown":     m.attributeSchemaChown(),
						"idmap":     m.attributeSchemaIDmap(),
					},
				},

				"bind": schema.SingleNestedAttribute{
					Description: "Bind Volume",
					Optional:    true,
					Computed:    false,
					Validators: []validator.Object{
						objectvalidator.ConflictsWith(
							path.MatchRelative().AtParent().AtName("volume"),
						),
					},
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.RequiresReplace(),
					},
					Attributes: map[string]schema.Attribute{
						"path": schema.StringAttribute{
							Description: "Host path",
							Required:    true,
							Computed:    false,
							Validators:  []validator.String{
								// TODO
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"read_only":   m.attributeSchemaReadOnly(),
						"dev":         m.attributeSchemaDev(),
						"exec":        m.attributeSchemaExec(),
						"suid":        m.attributeSchemaSuid(),
						"chown":       m.attributeSchemaChown(),
						"idmap":       m.attributeSchemaIDmap(),
						"propagation": m.attributeSchemaBindPropagation(),
						"recursive":   m.attributeSchemaBindRecursive(),
						"relabel":     m.attributeSchemaBindRelabel(),
					},
				},
			},
		},
	}

	// TODO:
	// While its technically possible, the podman api does not support this case for pods pretty well.
	// The tmpfs mount will be created on the infra container, but it is not exposed on inspect anymore
	// https://github.com/containers/podman/blob/v4.3.1/libpod/container_inspect.go#L276-L280
	//				"tmpfs": {
	//					Description: "Tmpfs Volume",
	//					Optional:    true,
	//					Computed:    false,
	//					Attributes: tfsdk.SingleNestedAttributes(
	//						map[string]tfsdk.Attribute{
	//							"read_only": m.attributeSchemaReadOnly(),
	//							"dev":       m.attributeSchemaDev(),
	//							"exec":      m.attributeSchemaExec(),
	//							"suid":      m.attributeSchemaSuid(),
	//							"chown":     m.attributeSchemaChown(),
	//							"size":      m.attributeSchemaTmpfsSize(),
	//							"mode":      m.attributeSchemaTmpfsMode(),
	//							"tmpcopyup": m.attributeSchemaTmpfsTmpCopyUp(),
	//						},
	//					),
	//				},

}

// ToPodmanSpec creates volume and mounts
func (m Mounts) ToPodmanSpec(diags *diag.Diagnostics) ([]*specgen.NamedVolume, []specs.Mount) {

	specNamedVolumes := make([]*specgen.NamedVolume, 0)
	specMounts := make([]specs.Mount, 0)
	for _, mount := range m {
		if mount.Volume != nil {
			// Named volume mount options
			specVol := specgen.NamedVolume{
				Name: mount.Volume.Name.ValueString(),
				Dest: mount.Destination.ValueString(),
			}

			specVol.Options = appendMountOptBool(specVol.Options, mount.Volume.ReadOnly, "ro", "rw")
			specVol.Options = appendMountOptBool(specVol.Options, mount.Volume.Dev, "dev", "nodev")
			specVol.Options = appendMountOptBool(specVol.Options, mount.Volume.Exec, "exec", "noexec")
			specVol.Options = appendMountOptBool(specVol.Options, mount.Volume.Suid, "sui", "nosuid")

			if !mount.Volume.Chown.IsNull() && mount.Volume.Chown.ValueBool() {
				specVol.Options = append(specVol.Options, "U")
			}

			if !mount.Volume.IDmap.IsNull() && mount.Volume.IDmap.ValueBool() {
				specVol.Options = append(specVol.Options, "idmap")
			}

			specNamedVolumes = append(specNamedVolumes, &specVol)

		} else if mount.Bind != nil {
			// Bind mount options
			specMount := specs.Mount{
				Destination: mount.Destination.ValueString(),
				Type:        "bind",
				Source:      mount.Bind.Path.ValueString(),
			}

			specMount.Options = appendMountOptBool(specMount.Options, mount.Bind.ReadOnly, "ro", "rw")
			specMount.Options = appendMountOptBool(specMount.Options, mount.Bind.Dev, "dev", "nodev")
			specMount.Options = appendMountOptBool(specMount.Options, mount.Bind.Exec, "exec", "noexec")
			specMount.Options = appendMountOptBool(specMount.Options, mount.Bind.Suid, "suid", "nosuid")

			if mount.Bind.Chown.ValueBool() {
				specMount.Options = append(specMount.Options, "U")
			}

			if mount.Bind.IDmap.ValueBool() {
				specMount.Options = append(specMount.Options, "idmap")
			}

			if mount.Bind.Propagation.ValueString() != "" {
				specMount.Options = append(specMount.Options, mount.Bind.Propagation.ValueString())
			}
			specMount.Options = appendMountOptBool(specMount.Options, mount.Bind.Recursive, "rbind", "bind")
			// public = z, private = Z
			specMount.Options = appendMountOptBool(specMount.Options, mount.Bind.Relabel, "z", "Z")

			specMounts = append(specMounts, specMount)
		}
		// else if mount.Tmpfs != nil {
		//
		// 			// Tmpfs mount options
		// 			specMount := specs.Mount{
		// 				Destination: mount.Destination.ValueString(),
		// 				Type:        "tmpfs",
		// 			}
		//
		// 			specMount.Options = appendMountOptBool(specMount.Options, mount.Tmpfs.ReadOnly, "ro", "rw")
		// 			specMount.Options = appendMountOptBool(specMount.Options, mount.Tmpfs.Dev, "dev", "nodev")
		// 			specMount.Options = appendMountOptBool(specMount.Options, mount.Tmpfs.Exec, "exec", "noexec")
		// 			specMount.Options = appendMountOptBool(specMount.Options, mount.Tmpfs.Suid, "suid", "nosuid")
		//
		// 			if mount.Tmpfs.Chown.ValueBool() {
		// 				specMount.Options = append(specMount.Options, "U")
		// 			}
		//
		// 			specMounts = append(specMounts, specMount)
		// 		}
	}
	return specNamedVolumes, specMounts
}

func FromPodmanToMounts(diags *diag.Diagnostics, specMounts []define.InspectMount) Mounts {
	mounts := make(Mounts, 0)

	for _, specMount := range specMounts {
		opts := parseMountOptions(diags, specMount.Options)
		if opts.readOnly.IsNull() {
			opts.readOnly = types.BoolValue(!specMount.RW)
		}

		switch specMount.Type {
		case "volume":
			mounts = append(mounts, Mount{
				Destination: types.StringValue(specMount.Destination),
				Volume: &MountVolume{
					Name:     types.StringValue(specMount.Name),
					ReadOnly: opts.readOnly,
					Dev:      opts.dev,
					Exec:     opts.exec,
					Suid:     opts.suid,
					Chown:    opts.chown,
					IDmap:    opts.idmap,
				},
			})

		case "bind":
			mounts = append(mounts, Mount{
				Destination: types.StringValue(specMount.Destination),
				Bind: &MountBind{
					Path:        types.StringValue(specMount.Source),
					ReadOnly:    opts.readOnly,
					Dev:         opts.dev,
					Exec:        opts.exec,
					Suid:        opts.suid,
					Chown:       opts.chown,
					IDmap:       opts.idmap,
					Propagation: opts.propagation,
					Recursive:   opts.recursive,
					Relabel:     opts.relabel,
				},
			})

			//		case "tmpfs":
			//			mounts = append(mounts, Mount{
			//				Destination: types.String{Value: specMount.Destination},
			//				Tmpfs: &MountTmpfs{
			//					ReadOnly: opts.readOnly,
			//					Dev:      opts.dev,
			//					Exec:     opts.exec,
			//					Suid:     opts.suid,
			//					Chown:    opts.chown,
			//					Size:     opts.size,
			//				},
			//			})

		default:
			diags.AddError("Unknown mount type retrieved", specMount.Type)
		}
	}
	if len(mounts) == 0 {
		return nil
	}
	return mounts
}

// appendMountOptBool appends a mapped boolen value
func appendMountOptBool(opts []string, v types.Bool, trueVal string, falseVal string) []string {
	if !v.IsNull() {
		if v.ValueBool() {
			opts = append(opts, trueVal)
		} else {
			opts = append(opts, falseVal)
		}
	}
	return opts
}
