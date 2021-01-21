package util

import (
	"encoding/json"
	"fmt"
	"strings"

	jsondiff "github.com/yazgazan/jaydiff/diff"

	"github.com/vetyy/kubetools/internal/functions"
)

type outputOptions struct {
	Indent     string `long:"indent" description:"indent string" default:"\t"`
	ShowTypes  bool   `long:"show-types" short:"t" description:"show types"`
	Colorized  bool
	JSON       bool `long:"json" description:"json-style output"`
	JSONValues bool
}

func Diff(actual []byte, expected []byte, encryptedVariables map[string]interface{}) ([]string, error) {
	var actualJson interface{}
	var expectedJson interface{}
	err := json.Unmarshal(actual, &actualJson)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(expected, &expectedJson)
	if err != nil {
		return nil, err
	}

	diff, err := jsondiff.Diff(actualJson, expectedJson)
	if err != nil {
		return nil, err
	}

	diffReport, err := jsondiff.Report(diff, jsondiff.Output(outputOptions{Colorized: true}))
	if err != nil {
		return nil, err
	}

	// TODO: very inefficient but works, should be reworked when possible
	redactIndexes := map[int]string{}
	redactB64Indexes := map[int]string{}
	for _, v := range encryptedVariables {
		vs := fmt.Sprintf("%v", v)
		b64vs := functions.Base64encode(vs)
		for i, d := range diffReport {
			if strings.Contains(d, vs) {
				redactIndexes[i] = vs
			}
			if strings.Contains(d, b64vs) {
				redactB64Indexes[i] = b64vs
			}
		}
	}

	var diffReportRedacted []string
	for i, d := range diffReport {
		if redactedVal, ok := redactIndexes[i]; ok {
			d = strings.Replace(d, redactedVal, "<REDACTED>", -1)
		}
		if redactedVal, ok := redactB64Indexes[i]; ok {
			d = strings.Replace(d, redactedVal, "<REDACTED>", -1)
		}
		diffReportRedacted = append(diffReportRedacted, d)
	}

	return diffReportRedacted, nil
}
