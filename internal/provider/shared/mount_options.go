package shared

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/project0/terraform-provider-podman/internal/modifier"
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
			resource.UseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
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
			resource.UseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
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
			resource.UseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
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
			resource.UseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
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
			resource.UseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
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
			resource.UseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
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
			resource.UseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
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
			resource.UseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
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
			resource.UseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
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
	result := allMountOptions{
		readOnly:    types.BoolNull(),
		exec:        types.BoolNull(),
		suid:        types.BoolNull(),
		chown:       types.BoolNull(),
		idmap:       types.BoolNull(),
		recursive:   types.BoolNull(),
		relabel:     types.BoolNull(),
		propagation: types.StringNull(),
	}

	for _, o := range options {

		switch o {
		case "ro", "rw":
			result.readOnly = types.BoolValue((o == "ro"))

		case "dev", "nodev":
			result.dev = types.BoolValue((o == "dev"))

		case "exec", "noexec":
			result.exec = types.BoolValue((o == "exec"))

		case "suid", "nosuid":
			result.suid = types.BoolValue((o == "suid"))

		case "bind", "rbind":
			result.recursive = types.BoolValue((o == "rbind"))

		// public = z (relabel), private = Z (no relabel)
		case "z", "Z":
			result.relabel = types.BoolValue((o == "z"))

		case "U":
			result.chown = types.BoolValue(true)

		case "idmap":
			result.idmap = types.BoolValue(true)

		case
			bindPropagationShared,
			bindPropagationSlave,
			bindPropagationPrivate,
			bindPropagationUnbindable,
			bindPropagationSharedRecursive,
			bindPropagationSlaveRecursive,
			bindPropagationPrivateRecursive,
			bindPropagationUnbindableRecursive:
			result.propagation = types.StringValue(o)

		default:
			diags.AddWarning("Unknown mount option retrieved", o)
		}
	}
	return result
}
