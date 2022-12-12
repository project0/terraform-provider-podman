package validators

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	regexName     = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`)
	regexOctal    = regexp.MustCompile(`^[0-7]{3,4}$`)
	regexTmpfSize = regexp.MustCompile(`^(\d+[kmg]?|\d{1,3}%)$`)
)

// MatchName validates given name to be compatible with podman
func MatchName() validator.String {
	return stringvalidator.RegexMatches(regexName, "")
}

// MatchOctal validates chmod octal representation
func MatchOctal() validator.String {
	return stringvalidator.RegexMatches(regexOctal, "")
}

// MatchTmpfSize validates the tmpfs size option
func MatchTmpfSize() validator.String {
	return stringvalidator.RegexMatches(regexTmpfSize, "")
}
