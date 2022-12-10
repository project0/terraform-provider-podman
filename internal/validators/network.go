package validators

import (
	"context"
	"fmt"
	"net"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	regexNetworkInterface = regexp.MustCompile(`^[a-z][a-z0-9]*$`)
)

func MatchNetworkInterfaceName() validator.String {
	return stringvalidator.RegexMatches(regexNetworkInterface, "")
}

func IsCIDR() validator.String {
	return &genericStringValidator{
		description: "",
		validate: func(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
			_, _, err := net.ParseCIDR(req.ConfigValue.ValueString())
			if err != nil {
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Failed to parse CIDR",
					fmt.Sprintf("invalid value: %s, error: %s", req.ConfigValue.String(), err.Error()),
				)
			}
		},
	}
}

func IsIpAdress() validator.String {
	return &genericStringValidator{
		description: "",
		validate: func(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
			if net.ParseIP(req.ConfigValue.ValueString()) == nil {
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Failed to parse IP address",
					fmt.Sprintf("invalid value: %s", req.ConfigValue.String()),
				)
			}
		},
	}
}
