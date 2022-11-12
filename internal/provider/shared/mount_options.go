package shared

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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

var (
	bindPropagations = []string{
		bindPropagationShared,
		bindPropagationSlave,
		bindPropagationPrivate,
		bindPropagationUnbindable,
		bindPropagationSharedRecursive,
		bindPropagationSlaveRecursive,
		bindPropagationPrivateRecursive,
		bindPropagationUnbindableRecursive,
	}
)

func (m Mounts) attributeSchemaReadOnly() tfsdk.Attribute {
	return tfsdk.Attribute{
		Description: "Mount as read only. Default depends on the mount type.",
		Computed:    true,
		Optional:    true,
		Type:        types.BoolType,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			modifier.AlwaysUseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
		},
	}
}

func (m Mounts) attributeSchemaSuid() tfsdk.Attribute {
	return tfsdk.Attribute{
		Description: "Mounting the volume with the nosuid(false) options means that SUID applications on the volume will not be able to change their privilege." +
			"By default volumes are mounted with nosuid.",
		Computed: true,
		Optional: true,
		Type:     types.BoolType,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			modifier.AlwaysUseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
		},
	}
}

func (m Mounts) attributeSchemaExec() tfsdk.Attribute {
	return tfsdk.Attribute{
		Description: "Mounting the volume with the noexec(false) option means that no executables on the volume will be able to executed within the pod." +
			"Defaults depends on the mount type or storage driver.",
		Computed: true,
		Optional: true,
		Type:     types.BoolType,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			modifier.AlwaysUseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
		},
	}
}

func (m Mounts) attributeSchemaDev() tfsdk.Attribute {
	return tfsdk.Attribute{
		Description: "Mounting the volume with the nodev(false) option means that no devices on the volume will be able to be used by processes within the container." +
			"By default volumes are mounted with nodev.",
		Computed: true,
		Optional: true,
		Type:     types.BoolType,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			modifier.AlwaysUseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
		},
	}
}

func (m Mounts) attributeSchemaChown() tfsdk.Attribute {
	return tfsdk.Attribute{
		Description: "Change recursively the owner and group of the source volume based on the UID and GID of the container.",
		Computed:    true,
		Optional:    true,
		Type:        types.BoolType,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			modifier.AlwaysUseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
		},
	}
}

func (m Mounts) attributeSchemaIDmap() tfsdk.Attribute {
	return tfsdk.Attribute{
		Description: "If specified, create an idmapped mount to the target user namespace in the container.",
		Computed:    true,
		Optional:    true,
		Type:        types.BoolType,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			modifier.AlwaysUseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
		},
	}
}

func (m Mounts) attributeSchemaBindPropagation() tfsdk.Attribute {
	return tfsdk.Attribute{
		Description: fmt.Sprintf("One of %s.", strings.Join(bindPropagations, ",")),
		Computed:    true,
		Optional:    true,
		Type:        types.StringType,
		Validators: []tfsdk.AttributeValidator{
			validator.OneOf(bindPropagations...),
		},
		PlanModifiers: tfsdk.AttributePlanModifiers{
			modifier.AlwaysUseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
		},
	}
}

func (m Mounts) attributeSchemaBindRecursive() tfsdk.Attribute {
	return tfsdk.Attribute{
		Description: "Set up a recursive bind mount. By default it is recursive.",
		Computed:    true,
		Optional:    true,
		Type:        types.BoolType,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			modifier.AlwaysUseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
		},
	}
}

func (m Mounts) attributeSchemaBindRelabel() tfsdk.Attribute {
	return tfsdk.Attribute{
		Description: "Labels the volume mounts. Sets the z (true) flag label the content with a shared content label, " +
			"or Z (false) flag to label the content with a private unshared label. " +
			"Default is unset (null).",
		Computed: true,
		Optional: true,
		Type:     types.BoolType,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			modifier.AlwaysUseStateForUnknown(),
			modifier.RequiresReplaceComputed(),
		},
	}
}

// func (m Mounts) attributeSchemaTmpfsSize() tfsdk.Attribute {
// 	return tfsdk.Attribute{
// 		Description: "Size of the tmpfs mount in bytes or units. Unlimited by default in Linux.",
// 		Computed:    true,
// 		Optional:    true,
// 		Type:        types.StringType,
// 		Validators: []tfsdk.AttributeValidator{
// 			validator.MatchTmpfSize(),
// 		},
// 		PlanModifiers: tfsdk.AttributePlanModifiers{
// 			modifier.AlwaysUseStateForUnknown(),
// 			modifier.RequiresReplaceComputed(),
// 		},
// 	}
// }
//
// func (m Mounts) attributeSchemaTmpfsMode() tfsdk.Attribute {
// 	return tfsdk.Attribute{
// 		Description: "File mode of the tmpfs in octal (e.g. 700 or 0700). Defaults to 1777 in Linux.",
// 		Computed:    true,
// 		Optional:    true,
// 		Type:        types.StringType,
// 		Validators: []tfsdk.AttributeValidator{
// 			validator.MatchOctal(),
// 		},
// 		PlanModifiers: tfsdk.AttributePlanModifiers{
// 			modifier.AlwaysUseStateForUnknown(),
// 			modifier.RequiresReplaceComputed(),
// 		},
// 	}
// }
//
// func (m Mounts) attributeSchemaTmpfsTmpCopyUp() tfsdk.Attribute {
// 	return tfsdk.Attribute{
// 		Description: "Enable copyup from the image directory at the same location to the tmpfs. Used by default.",
// 		Computed:    true,
// 		Optional:    true,
// 		Type:        types.BoolType,
// 		PlanModifiers: tfsdk.AttributePlanModifiers{
// 			modifier.AlwaysUseStateForUnknown(),
// 			modifier.RequiresReplaceComputed(),
// 		},
// 	}
// }

type allMountOptions struct {
	readOnly types.Bool
	dev      types.Bool
	exec     types.Bool
	suid     types.Bool
	chown    types.Bool
	idmap    types.Bool
	// bind
	recursive   types.Bool
	relabel     types.Bool
	propagation types.String
	// tmpfs
	size      types.String
	mode      types.String
	tmpcopyup types.Bool
}

func parseMountOptions(diags *diag.Diagnostics, options []string) allMountOptions {
	result := allMountOptions{
		readOnly: types.BoolNull(),
		dev:      types.BoolNull(),
		exec:     types.BoolNull(),
		suid:     types.BoolNull(),

		// chown and idmap is only present when flag is set,
		// consider it false when not present (default)
		chown: types.BoolValue(false),
		idmap: types.BoolValue(false),

		// bind
		recursive:   types.BoolNull(),
		relabel:     types.BoolNull(),
		propagation: types.StringNull(),

		// tmpfs
		size:      types.StringNull(),
		mode:      types.StringNull(),
		tmpcopyup: types.BoolNull(),
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

		case "tmpcopyup", "notmpcopyup":
			result.recursive = types.BoolValue((o == "tmpcopyup"))

		// public = z (relabel), private = Z (no relabel)
		case "z", "Z":
			result.relabel = types.BoolValue((o == "z"))

		case "U":
			result.chown = types.BoolValue(true)

		case "idmap":
			result.idmap = types.BoolValue(true)

		case "size":
			result.size = types.StringValue(o)

		case "mode":
			result.mode = types.StringValue(o)

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
