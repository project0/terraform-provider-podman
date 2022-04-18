package validator

import (
	"context"
	"fmt"
	"net"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var (
	regexNetworkInterface = regexp.MustCompile(`^[a-z][a-z0-9]*$`)
)

func MatchNetworkInterfaceName() tfsdk.AttributeValidator {
	return MatchRegex(regexNetworkInterface)
}

func IsCIDR() tfsdk.AttributeValidator {
	return &genericStringValidator{
		description: "",
		validate: func(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse, str string) {
			_, _, err := net.ParseCIDR(str)
			if err != nil {
				resp.Diagnostics.AddAttributeError(
					req.AttributePath,
					"Failed to parse CIDR",
					fmt.Sprintf("invalid value: %s, error: %s", str, err.Error()),
				)
			}
		},
	}
}

func IsIpAdress() tfsdk.AttributeValidator {
	return &genericStringValidator{
		description: "",
		validate: func(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse, str string) {
			if net.ParseIP(str) == nil {
				resp.Diagnostics.AddAttributeError(
					req.AttributePath,
					"Failed to parse IP address",
					fmt.Sprintf("invalid value: %s", str),
				)
			}
		},
	}
}
