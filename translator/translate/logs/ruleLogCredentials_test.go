// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package logs

import (
	"encoding/json"
	"testing"

	"github.com/aws/amazon-cloudwatch-agent/translator/context"
	"github.com/aws/amazon-cloudwatch-agent/translator/translate/agent"
	"github.com/stretchr/testify/assert"
)

func TestWithAgentConfig(t *testing.T) {
	agent.Global_Config.Credentials = map[string]interface{}{}
	ctx := context.CurrentContext()
	ctx.SetCredentials(map[string]string{})
	c := new(LogCreds)
	var input interface{}
	e := json.Unmarshal([]byte(`{ "credentials" : {"access_key":"metric_ak", "secret_key":"metric_sk", "token": "dummy_token", "profile": "dummy_profile", "role_arn": "role_value"}}`), &input)
	if e == nil {
		_, returnVal := c.ApplyRule(input)
		assert.Equal(t, "role_value", returnVal.(map[string]interface{})["role_arn"], "Expected to be equal")
	} else {
		panic(e)
	}

	agent.Global_Config.Role_arn = "global_role_arn_test"
	e = json.Unmarshal([]byte(`{ "credentials" : {"access_key":"metric_ak", "secret_key":"metric_sk", "token": "dummy_token", "profile": "dummy_profile", "role_arn": "role_value"}}`), &input)
	if e == nil {
		_, returnVal := c.ApplyRule(input)
		assert.Equal(t, "role_value", returnVal.(map[string]interface{})["role_arn"], "Expected to be equal")
	} else {
		panic(e)
	}

	agent.Global_Config.Role_arn = "global_role_arn_test"
	e = json.Unmarshal([]byte(`{ "credentials" : {"access_key":"metric_ak", "secret_key":"metric_sk", "token": "dummy_token", "profile": "dummy_profile"}}`), &input)
	if e == nil {
		_, returnVal := c.ApplyRule(input)
		assert.Equal(t, "global_role_arn_test", returnVal.(map[string]interface{})["role_arn"], "Expected to be equal")
	} else {
		panic(e)
	}

	agent.Global_Config.Role_arn = ""
}
