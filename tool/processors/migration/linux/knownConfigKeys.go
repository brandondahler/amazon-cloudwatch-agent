// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package linux

import (
	"fmt"
	"strconv"

	"github.com/aws/amazon-cloudwatch-agent/tool/data/config"

	"strings"

	"github.com/bigkevmcd/go-configparser"
)

// Get the keys from https://code.amazon.com/packages/AwsCWLogsPlugin/blobs/9efe1aa104aced70f82bb92ee3657f4598266d77/--/cwlogs/push.py#L382-L415
var knownConfigKeys = []string{
	"file",                     // "file_path"
	"log_group_name",           // "log_group_name"
	"log_stream_name",          // "log_stream_name", currently only single value in the output
	"datetime_format",          // "timestamp_format", refer to https://code.amazon.com/packages/GoAmzn-CWAgentConfigTranslator/blobs/mainline/--/src/translator/translate/logs/logs_collected/files/collect_list/ruleTimestampFormat.go
	"time_zone",                // "timezone", UTC or LOCAL
	"multi_line_start_pattern", // "multi_line_start_pattern"
	"encoding",                 // "encoding"
	"buffer_duration",          // "force_flush_interval", from ms to sec

	"use_gzip_http_content_encoding", // Not used in new agent. Auto choose when the payload is optimized by this
	"queue_size",                     // Not used in new agent
	"initial_position",               // Not really used in new agent. Always set to start from beginning.
	"file_fingerprint_lines",         // Not used in new agent
	"batch_size",                     // Not used in new agent
	"batch_count",                    // Not used in new agent
}

func isUnknownKey(key string) bool {
	for _, knownConfigKey := range knownConfigKeys {
		if key == knownConfigKey {
			return false
		}
	}
	return true
}
func addLogConfig(logsConfig *config.Logs, filePath, section string, p *configparser.ConfigParser) {
	options, err := p.Options(section)
	if err != nil {
		fmt.Printf("Error in fetching options for section %s in file %s:\n%v\n", section, filePath, err)
		return
	}
	for _, k := range options {
		if isUnknownKey(strings.ToLower(strings.TrimSpace(k))) {
			fmt.Printf("Warning: Option key %s for section %s in file %s is unknown.\n", k, section, filePath)
		}
	}
	logFilePath, _ := p.Get(section, "file")
	logGroupName, _ := p.Get(section, "log_group_name")
	logStreamName, _ := p.Get(section, "log_stream_name")
	timestampFormat, _ := p.Get(section, "datetime_format")
	timezone, _ := p.Get(section, "time_zone")
	multiLineStartPattern, _ := p.Get(section, "multi_line_start_pattern")
	if multiLineStartPattern == "{datetime_format}" {
		multiLineStartPattern = "{timestamp_format}"
	}
	encoding, _ := p.Get(section, "encoding")
	if encoding != "" {
		normalized := NormalizeEncoding(encoding)
		if normalized == "" {
			panic(fmt.Sprintf("Encoding %s is not supported.", encoding))
		} else {
			encoding = normalized
		}
	}
	bufferDuration, _ := p.Get(section, "buffer_duration")
	if bufferDuration != "" {
		if forceFlushInterval, err := strconv.Atoi(bufferDuration); err == nil {
			forceFlushInterval /= 1000 // from ms to sec
			if logsConfig.ForceFlushInterval == 0 {
				logsConfig.ForceFlushInterval = forceFlushInterval
			} else if logsConfig.ForceFlushInterval != forceFlushInterval {
				fmt.Printf("Warning: The buffer_duration was set to different values (existing value: %v sec, new value: %v sec) for different files. Use 1st buffer_duration value.",
					logsConfig.ForceFlushInterval, forceFlushInterval)
			}
		}
	}
	logsConfig.AddLogFile(logFilePath, logGroupName, logStreamName, timestampFormat, timezone, multiLineStartPattern, encoding)
}
