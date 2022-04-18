package utils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const (
	diagMsgErrorSummary = "Unexpected provider error: %s"
	diagMsgErrorDetail  = "%s This is always a bug in the provider code and should be reported to the provider developers."
)

// AddUnexpectedError adds a diagnostic error with injected bug hint
func AddUnexpectedError(d *diag.Diagnostics, summary, detail string) {
	d.AddError(
		fmt.Sprintf(diagMsgErrorSummary, summary),
		fmt.Sprintf(diagMsgErrorDetail, detail),
	)
}

// AddUnexpectedAttributeError adds a diagnostic error with injected bug hint
func AddUnexpectedAttributeError(path *tftypes.AttributePath, d *diag.Diagnostics, summary, detail string) {
	d.AddAttributeError(
		path,
		fmt.Sprintf(diagMsgErrorSummary, summary),
		fmt.Sprintf(diagMsgErrorDetail, detail),
	)
}
