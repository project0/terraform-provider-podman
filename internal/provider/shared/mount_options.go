package shared

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/project0/terraform-provider-podman/internal/validator"
)

const (
	bindPropagationShared     = "shared"
	bindPropagationSlave      = "slave"
	bindPropagationPrivate    = "private"
	bindPropagationUnbindable = "unbindable"

	bindPropagationSharedRecursive     = "rshared"
	bindPropagationSlaveRecursive      = "rslave"
	bindPropagationPrivateRecursive    = "rprivate"
	bindPropagationUnbindableRecursive = "runbindable"
)

func (m Mounts) attributeSchemaReadOnly() tfsdk.Attribute {
	return tfsdk.Attribute{
		Description: "TODO",
		Computed:    true,
		Optional:    true,
		Type:        types.BoolType,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			tfsdk.UseStateForUnknown(),
			tfsdk.RequiresReplace(),
		},
	}
}

func (m Mounts) attributeSchemaSuid() tfsdk.Attribute {
	return tfsdk.Attribute{
		Description: "TODO",
		Computed:    true,
		Optional:    true,
		Type:        types.BoolType,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			tfsdk.UseStateForUnknown(),
			tfsdk.RequiresReplace(),
		},
	}
}

func (m Mounts) attributeSchemaExec() tfsdk.Attribute {
	return tfsdk.Attribute{
		Description: "Mounting the volume with the exec(true) or noexec(false) option means that no executables on the volume will be able to executed within the pod." +
			"Defaults depends on the mount type or storage driver.",
		Computed: true,
		Optional: true,
		Type:     types.BoolType,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			tfsdk.UseStateForUnknown(),
			tfsdk.RequiresReplace(),
		},
	}
}

func (m Mounts) attributeSchemaDev() tfsdk.Attribute {
	return tfsdk.Attribute{
		Description: "TODO",
		Computed:    true,
		Optional:    true,
		Type:        types.BoolType,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			tfsdk.UseStateForUnknown(),
			tfsdk.RequiresReplace(),
		},
	}
}

func (m Mounts) attributeSchemaChown() tfsdk.Attribute {
	return tfsdk.Attribute{
		Description: "TODO",
		Computed:    true,
		Optional:    true,
		Type:        types.BoolType,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			tfsdk.UseStateForUnknown(),
			tfsdk.RequiresReplace(),
		},
	}

}

func (m Mounts) attributeSchemaIDmap() tfsdk.Attribute {
	return tfsdk.Attribute{
		Description: "TODO",
		Computed:    true,
		Optional:    true,
		Type:        types.BoolType,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			tfsdk.UseStateForUnknown(),
			tfsdk.RequiresReplace(),
		},
	}
}

func (m Mounts) attributeSchemaBindPropagation() tfsdk.Attribute {
	return tfsdk.Attribute{
		Description: "TODO",
		Computed:    true,
		Optional:    true,
		Type:        types.StringType,
		Validators: []tfsdk.AttributeValidator{
			validator.OneOf(
				bindPropagationShared,
				bindPropagationSlave,
				bindPropagationPrivate,
				bindPropagationUnbindable,
				bindPropagationSharedRecursive,
				bindPropagationSlaveRecursive,
				bindPropagationPrivateRecursive,
				bindPropagationUnbindableRecursive,
			),
		},
		PlanModifiers: tfsdk.AttributePlanModifiers{
			tfsdk.UseStateForUnknown(),
			tfsdk.RequiresReplace(),
		},
	}
}

func (m Mounts) attributeSchemaBindRecursive() tfsdk.Attribute {
	return tfsdk.Attribute{
		Description: "TODO",
		Computed:    true,
		Optional:    true,
		Type:        types.BoolType,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			tfsdk.UseStateForUnknown(),
			tfsdk.RequiresReplace(),
		},
	}
}

func (m Mounts) attributeSchemaRelabel() tfsdk.Attribute {
	return tfsdk.Attribute{
		Description: "Labels the volume mounts. Sets the z (true) flag label the content with a shared content label, " +
			"or Z (false) flag to label the content with a private unshared label. " +
			"Default is unset (null).",
		Computed: true,
		Optional: true,
		Type:     types.BoolType,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			tfsdk.UseStateForUnknown(),
			tfsdk.RequiresReplace(),
		},
	}
}

type allMountOptions struct {
	readOnly    types.Bool
	dev         types.Bool
	exec        types.Bool
	suid        types.Bool
	chown       types.Bool
	idmap       types.Bool
	recursive   types.Bool
	relabel     types.Bool
	propagation types.String
}

func parseMountOptions(diags *diag.Diagnostics, options []string) allMountOptions {
	result := allMountOptions{}

	result.readOnly.Null = true
	// result.dev.Null = true
	// result.exec.Null = true
	// result.suid.Null = true
	// result.chown.Null = true
	// result.idmap.Null = true
	// result.recursive.Null = true
	// result.relabel.Null = true
	// result.propagation.Null = true

	for _, o := range options {

		switch o {
		case "ro", "rw":
			result.readOnly.Null = false
			if o == "ro" {
				result.readOnly.Value = true
			}

		case "dev", "nodev":
			result.dev.Null = false
			if o == "dev" {
				result.dev.Value = true
			}

		case "exec", "noexec":
			result.dev.Null = false
			if o == "exec" {
				result.exec.Value = true
			}
		case "suid", "nosuid":
			result.suid.Null = false
			if o == "suid" {
				result.suid.Value = true
			}

		case "bind", "rbind":
			result.recursive.Null = false
			if o == "rbind" {
				result.recursive.Value = true
			}

		// public = z, private = Z
		case "z", "Z":
			result.relabel.Null = false
			if o == "z" {
				result.relabel.Value = true
			}

		case "U":
			result.chown.Null = false
			result.chown.Value = true

		case "idmap":
			result.chown.Null = false
			result.chown.Value = true

		case
			bindPropagationShared,
			bindPropagationSlave,
			bindPropagationPrivate,
			bindPropagationUnbindable,
			bindPropagationSharedRecursive,
			bindPropagationSlaveRecursive,
			bindPropagationPrivateRecursive,
			bindPropagationUnbindableRecursive:
			result.propagation.Null = false
			result.propagation.Value = o

		default:
			diags.AddWarning("Unknown mount option retrieved", o)
		}
	}
	return result
}
