package template

import _ "embed"

//go:embed otp-email-template.html
var otpEmailTemplate string

//go:embed financial-report-template.html
var financialReportTemplate string

var Template = map[string]string{
	"otp-email-template.html":        otpEmailTemplate,
	"financial-report-template.html": financialReportTemplate,
} 