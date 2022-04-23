package validator

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var (
	regexName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`)
)

func MatchName() tfsdk.AttributeValidator {
	return MatchRegex(regexName)
}

func MatchRegex(regex *regexp.Regexp) tfsdk.AttributeValidator {
	return &genericStringValidator{
		description: "string must match pattern " + regex.String(),
		validate: func(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse, str string) {
			if !regex.MatchString(str) {
				resp.Diagnostics.AddAttributeError(req.AttributePath, "String did not match pattern "+regex.String(), "")
			}
		},
	}
}

func OneOf(values ...string) tfsdk.AttributeValidator {
	return &genericStringValidator{
		description: "string must be one of " + strings.Join(values, ","),
		validate: func(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse, str string) {
			for _, v := range values {
				if v == str {
					return
				}
			}
			resp.Diagnostics.AddAttributeError(req.AttributePath, "Incompatible value", fmt.Sprintf("%s must be one of %s", str, strings.Join(values, ",")))
		},
	}
}
