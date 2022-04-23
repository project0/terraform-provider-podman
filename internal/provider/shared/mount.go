package shared

import (
	"github.com/containers/podman/v4/libpod/define"
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/project0/terraform-provider-podman/internal/validator"
)

type (
	Mounts []Mount
	Mount  struct {
		Destination types.String `tfsdk:"destination"`

		Volume *MountVolume `tfsdk:"volume"`
		Bind   *MountBind   `tfsdk:"bind"`
		//Tmpfs  *MountTmpfs
	}

	//MountTmpfs struct{}
	MountVolume struct {
		Name types.String `tfsdk:"name"`

		ReadOnly types.Bool `tfsdk:"read_only"`
		Suid     types.Bool `tfsdk:"suid"`
		Exec     types.Bool `tfsdk:"exec"`
		Dev      types.Bool `tfsdk:"dev"`
		Chown    types.Bool `tfsdk:"chown"`
		IDmap    types.Bool `tfsdk:"idmap"`
	}

	// aka host path
	MountBind struct {
		Path types.String `tfsdk:"path"`

		ReadOnly types.Bool `tfsdk:"read_only"`
		Suid     types.Bool `tfsdk:"suid"`
		Exec     types.Bool `tfsdk:"exec"`
		Dev      types.Bool `tfsdk:"dev"`
		Chown    types.Bool `tfsdk:"chown"`
		IDmap    types.Bool `tfsdk:"idmap"`

		Propagation types.String `tfsdk:"propagation"`
		Recursive   types.Bool   `tfsdk:"recursive"`
		Relabel     types.Bool   `tfsdk:"relabel"`
	}
)

func (m Mounts) AttributeSchema() tfsdk.Attribute {

	return tfsdk.Attribute{
		Description: "Mounts volume, bind, image, tmpfs, etc..",
		Required:    false,
		Optional:    true,
		Computed:    false,
		// TODO/IDEA: Change SetNested to MapNested and use key as destination?
		Attributes: tfsdk.SetNestedAttributes(
			map[string]tfsdk.Attribute{
				"destination": {
					Description: "Target path",
					Required:    true,
					Computed:    false,
					Type:        types.StringType,
					Validators:  []tfsdk.AttributeValidator{
						//TODO
					},
					PlanModifiers: tfsdk.AttributePlanModifiers{
						tfsdk.RequiresReplace(),
					},
				},
				// TODO validate conflicts with
				"volume": {
					Description: "Named Volume",
					Optional:    true,
					Computed:    true,
					Attributes: tfsdk.SingleNestedAttributes(
						map[string]tfsdk.Attribute{
							"name": {
								Description: "Name of the volume",
								Required:    true,
								Computed:    false,
								Type:        types.StringType,
								Validators: []tfsdk.AttributeValidator{
									validator.MatchName(),
								},
							},
							"read_only": m.attributeSchemaReadOnly(),
							"dev":       m.attributeSchemaDev(),
							"exec":      m.attributeSchemaExec(),
							"suid":      m.attributeSchemaSuid(),
							"chown":     m.attributeSchemaChown(),
							"idmap":     m.attributeSchemaIDmap(),
						},
					),
					PlanModifiers: tfsdk.AttributePlanModifiers{
						tfsdk.RequiresReplace(),
					},
				},
				"bind": {
					Description: "Named Volume",
					Optional:    true,
					Computed:    true,
					Attributes: tfsdk.SingleNestedAttributes(
						map[string]tfsdk.Attribute{
							"path": {
								Description: "Host path",
								Required:    true,
								Computed:    false,
								Type:        types.StringType,
								Validators:  []tfsdk.AttributeValidator{
									//TODO
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
							"relabel":     m.attributeSchemaRelabel(),
						},
					),
					PlanModifiers: tfsdk.AttributePlanModifiers{
						tfsdk.UseStateForUnknown(),
						tfsdk.RequiresReplace(),
					},
				},
			},
			tfsdk.SetNestedAttributesOptions{},
		),
		PlanModifiers: tfsdk.AttributePlanModifiers{
			tfsdk.RequiresReplace(),
		},
	}

}

// ToPodmanSpec creates volume and mounts
func (m Mounts) ToPodmanSpec(diags *diag.Diagnostics) ([]*specgen.NamedVolume, []specs.Mount) {

	specNamedVolumes := make([]*specgen.NamedVolume, 0)
	specMounts := make([]specs.Mount, 0)

	for _, mount := range m {
		if mount.Volume != nil {
			// Named volume mount options
			specVol := specgen.NamedVolume{
				Name: mount.Volume.Name.Value,
				Dest: mount.Destination.Value,
			}

			specVol.Options = appendMountOptBool(specVol.Options, mount.Volume.ReadOnly, "ro", "rw")
			specVol.Options = appendMountOptBool(specVol.Options, mount.Volume.Dev, "dev", "nodev")
			specVol.Options = appendMountOptBool(specVol.Options, mount.Volume.Exec, "exec", "noexec")
			specVol.Options = appendMountOptBool(specVol.Options, mount.Volume.Suid, "sui", "nosuid")

			if !mount.Volume.Chown.Null && mount.Volume.Chown.Value {
				specVol.Options = append(specVol.Options, "U")
			}

			if !mount.Volume.IDmap.Null && mount.Volume.IDmap.Value {
				specVol.Options = append(specVol.Options, "idmap")
			}

			specNamedVolumes = append(specNamedVolumes, &specVol)

		} else if mount.Bind != nil {
			// Bind mount options
			specMount := specs.Mount{
				Destination: mount.Destination.Value,
				Type:        "bind",
				Source:      mount.Bind.Path.Value,
			}

			specMount.Options = appendMountOptBool(specMount.Options, mount.Bind.ReadOnly, "ro", "rw")
			specMount.Options = appendMountOptBool(specMount.Options, mount.Bind.Dev, "dev", "nodev")
			specMount.Options = appendMountOptBool(specMount.Options, mount.Bind.Exec, "exec", "noexec")
			specMount.Options = appendMountOptBool(specMount.Options, mount.Bind.Suid, "suid", "nosuid")

			if !mount.Bind.Chown.Null && mount.Bind.Chown.Value {
				specMount.Options = append(specMount.Options, "U")
			}

			if !mount.Bind.IDmap.Null && mount.Bind.IDmap.Value {
				specMount.Options = append(specMount.Options, "idmap")
			}

			if !mount.Bind.Propagation.Null && mount.Bind.Propagation.Value != "" {
				specMount.Options = append(specMount.Options, mount.Bind.Propagation.Value)
			}
			specMount.Options = appendMountOptBool(specMount.Options, mount.Bind.Recursive, "rbind", "bind")
			// public = z, private = Z
			specMount.Options = appendMountOptBool(specMount.Options, mount.Bind.Relabel, "z", "Z")

			specMounts = append(specMounts, specMount)
		}
	}
	return specNamedVolumes, specMounts
}

func FromPodmanToMounts(diags *diag.Diagnostics, specMounts []define.InspectMount) Mounts {
	mounts := make(Mounts, 0)

	for _, specMount := range specMounts {

		opts := parseMountOptions(diags, specMount.Options)
		readOnly := types.Bool{Value: !specMount.RW}

		switch specMount.Type {
		case "volume":
			mounts = append(mounts, Mount{
				Destination: types.String{Value: specMount.Destination},
				Volume: &MountVolume{
					Name:     types.String{Value: specMount.Name},
					ReadOnly: readOnly,
					Dev:      opts.dev,
					Exec:     opts.exec,
					Suid:     opts.suid,
					Chown:    opts.chown,
					IDmap:    opts.idmap,
				},
			})
		case "bind":
			mounts = append(mounts, Mount{
				Destination: types.String{Value: specMount.Destination},
				Bind: &MountBind{
					Path:        types.String{Value: specMount.Source},
					ReadOnly:    readOnly,
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
	if !v.Null {
		if v.Value {
			opts = append(opts, trueVal)
		} else {
			opts = append(opts, falseVal)
		}
	}
	return opts
}
