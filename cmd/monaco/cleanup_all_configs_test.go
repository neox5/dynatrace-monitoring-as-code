// +build cleanup

/**
 * @license
 * Copyright 2020 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"regexp"
	"strings"
	"testing"

	"github.com/dynatrace-oss/dynatrace-monitoring-as-code/pkg/api"
	"github.com/dynatrace-oss/dynatrace-monitoring-as-code/pkg/environment"
	"github.com/dynatrace-oss/dynatrace-monitoring-as-code/pkg/rest"
	"github.com/dynatrace-oss/dynatrace-monitoring-as-code/pkg/util"
	"gotest.tools/assert"
)

func TestDoCleanup(t *testing.T) {

	environments, errs := environment.LoadEnvironmentList("", "test-resources/integration-multi-environment/environments.yaml", util.NewFileReader())
	for _, err := range errs {
		assert.NilError(t, err)
	}

	apis := api.NewApis()

	r, _ := regexp.Compile(`^.+_.*\d+.*$`)

	for _, environment := range environments {
		token, err := environment.GetToken()
		assert.NilError(t, err)
		for _, api := range apis {

			_, values, err := rest.GetExistingValuesFromEndpoint(api.GetId(), api.GetUrl(environment), token)
			assert.NilError(t, err)

			for _, value := range values {
				if r.MatchString(value.Name) || r.MatchString(value.Id) || strings.HasSuffix(value.Name, "_") {
					util.Log.Info("Deleting %s (%s)\n", value.Name, api.GetId())
					rest.Delete(api.GetUrl(environment), token, value.Id)
				}
			}
		}
	}
}
