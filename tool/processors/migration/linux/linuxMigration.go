// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package linux

import (
	"fmt"

	"github.com/aws/amazon-cloudwatch-agent/tool/data"
	"github.com/aws/amazon-cloudwatch-agent/tool/data/config"
	"github.com/aws/amazon-cloudwatch-agent/tool/processors"
	"github.com/aws/amazon-cloudwatch-agent/tool/processors/question/logs"
	"github.com/aws/amazon-cloudwatch-agent/tool/runtime"
	"github.com/aws/amazon-cloudwatch-agent/tool/util"

	"github.com/bigkevmcd/go-configparser"
)

const (
	genericSectionName                = "general"
	anyExistingLinuxConfigQuestion    = "Do you have any existing CloudWatch Log Agent (http://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/AgentReference.html) configuration file to import for migration?"
	filePathLinuxConfigQuestion       = "What is the file path for the existing cloudwatch log agent configuration file?"
	DefaultFilePathLinuxConfiguration = "/var/awslogs/etc/awslogs.conf"
)

var Processor processors.Processor = &processor{}

type processor struct{}

func (p *processor) Process(ctx *runtime.Context, config *data.Config) {
	if ctx.HasExistingLinuxConfig || util.No(anyExistingLinuxConfigQuestion) {
		filePath := ctx.ConfigFilePath
		if filePath == "" {
			filePath = util.AskWithDefault(filePathLinuxConfigQuestion, DefaultFilePathLinuxConfiguration)
		}
		processConfigFromPythonConfigParserFile(filePath, config.LogsConf())
	}
}

func (p *processor) NextProcessor(ctx *runtime.Context, config *data.Config) interface{} {
	return logs.Processor
}

func processConfigFromPythonConfigParserFile(filePath string, logsConfig *config.Logs) {
	p, err := configparser.NewConfigParserFromFile(filePath)
	if err != nil {
		fmt.Printf("Error in reading old python config from file %s: %v\n", filePath, err)
		panic(err)
	}
	if p.HasSection(genericSectionName) {
		err := p.RemoveSection(genericSectionName)
		if err != nil {
			fmt.Printf("Error in removing generic section from the config file %s:\n %v\n", filePath, err)
			panic(err)
		}
	}
	for _, section := range p.Sections() {
		addLogConfig(logsConfig, filePath, section, p)
	}
}
