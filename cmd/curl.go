package cmd

import (
	"bytes"
	"encoding/json"
	"strings"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/uaa"
	"github.com/spf13/cobra"
)

func GetCurlValidations(cfg uaa.Config, args []string) error {
	if err := EnsureTargetInConfig(cfg); err != nil {
		return err
	}
	if len(args) < 1 {
		return MissingArgumentError("path")
	}
	return nil
}

func CurlCmd(cm uaa.CurlManager, logger cli.Logger, path, method, data string, headers []string) error {
	resHeaders, resBody, err := cm.Curl(path, method, data, headers)
	if err != nil {
		return err
	}
	bodyPrinter := logger.Info
	if strings.Contains(resHeaders, "application/json") {
		buffer := bytes.Buffer{}
		if err := json.Indent(&buffer, []byte(resBody), "", "  "); err == nil {
			resBody = buffer.String()
			bodyPrinter = logger.Robots
		}
	}
	bodyPrinter(resBody)
	return nil
}

var curlCmd = &cobra.Command{
	Use:   "curl",
	Short: "CURL to a UAA endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyValidationErrors(GetCurlValidations(cfg, args), cmd, log)
		cm := uaa.CurlManager{GetHttpClient(), cfg}
		err := CurlCmd(cm, log, args[0], method, data, headers)
		NotifyErrorsWithRetry(err, cfg, log)
	},
}

func init() {
	RootCmd.AddCommand(curlCmd)
	curlCmd.Annotations = make(map[string]string)
	curlCmd.Annotations[MISC_CATEGORY] = "true"

	curlCmd.Flags().StringVarP(&method, "method", "X", "GET", "HTTP method (GET, POST, PUT, DELETE, etc)")
	curlCmd.Flags().StringVarP(&data, "data", "d", "", "HTTP data to include in the request body")
	curlCmd.Flags().StringSliceVarP(&headers, "header", "H", []string{}, "Custom headers to include in the request, flag can be specified multiple times")
}
