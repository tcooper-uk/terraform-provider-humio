// Copyright © 2020 Humio Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package humio

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	humio "github.com/humio/cli/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccViewRequiredFields(t *testing.T) {
	config := viewEmpty
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found.`)},
		{Config: config, ExpectError: regexp.MustCompile(`Insufficient repository blocks`)},
	}, nil)
}

func TestAccViewInvalidInputs(t *testing.T) {
	config := viewInvalidInputs
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`At least 1 "repository" blocks are required.`)},
	}, nil)
}

func TestAccViewBasic(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: viewBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_view.test", "repository.0.name", "first"),
				resource.TestCheckResourceAttr("humio_view.test", "repository.0.filter", "*"),
				resource.TestCheckResourceAttr("humio_view.test", "name", "simple-view"),
				resource.TestCheckResourceAttr("humio_view.test", "description", "a description"),
			),
		},
	}, testAccCheckAlertDestroy)
}

func TestAccViewBasicToFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: viewBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_view.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_view.test", "name", "alert-test"),
				resource.TestCheckResourceAttr("humio_view.test", "throttle_time_millis", "3600000"),
				resource.TestCheckResourceAttr("humio_view.test", "start", "24h"),
				resource.TestCheckResourceAttr("humio_view.test", "query", "loglevel=ERROR"),
				resource.TestCheckResourceAttr("humio_view.test", "description", "some text"),
				resource.TestCheckResourceAttr("humio_view.test", "silenced", "true"),
				resource.TestCheckResourceAttr("humio_view.test", "labels.#", "2"),
				resource.TestCheckResourceAttr("humio_view.test", "labels.0", "errors"),
				resource.TestCheckResourceAttr("humio_view.test", "labels.1", "important"),
				resource.TestCheckResourceAttr("humio_view.test", "actions.#", "1"),
				resource.TestCheckResourceAttrSet("humio_view.test", "actions.0"),
			),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
		{
			Config: viewBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_view.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_view.test", "name", "alert-test"),
				resource.TestCheckResourceAttr("humio_view.test", "throttle_time_millis", "3600000"),
				resource.TestCheckResourceAttr("humio_view.test", "start", "24h"),
				resource.TestCheckResourceAttr("humio_view.test", "query", "loglevel=ERROR"),
				resource.TestCheckResourceAttr("humio_view.test", "description", "some text"),
				resource.TestCheckResourceAttr("humio_view.test", "silenced", "true"),
				resource.TestCheckResourceAttr("humio_view.test", "labels.#", "2"),
				resource.TestCheckResourceAttr("humio_view.test", "labels.0", "errors"),
				resource.TestCheckResourceAttr("humio_view.test", "labels.1", "important"),
				resource.TestCheckResourceAttr("humio_view.test", "actions.#", "1"),
				resource.TestCheckResourceAttrSet("humio_view.test", "actions.0"),
			),
		},
	}, testAccCheckViewDestroy)
}

func TestAccViewFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: viewBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_view.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_view.test", "name", "alert-test"),
				resource.TestCheckResourceAttr("humio_view.test", "throttle_time_millis", "3600000"),
				resource.TestCheckResourceAttr("humio_view.test", "start", "24h"),
				resource.TestCheckResourceAttr("humio_view.test", "query", "loglevel=ERROR"),
				resource.TestCheckResourceAttr("humio_view.test", "description", "some text"),
				resource.TestCheckResourceAttr("humio_view.test", "silenced", "true"),
				resource.TestCheckResourceAttr("humio_view.test", "labels.#", "2"),
				resource.TestCheckResourceAttr("humio_view.test", "labels.0", "errors"),
				resource.TestCheckResourceAttr("humio_view.test", "labels.1", "important"),
				resource.TestCheckResourceAttr("humio_view.test", "actions.#", "1"),
				resource.TestCheckResourceAttrSet("humio_view.test", "actions.0"),
			),
		},
	}, testAccCheckViewDestroy)
}

func testAccCheckViewDestroy(s *terraform.State) error {
	conn := testAccProviders["humio"].Meta().(*humio.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "humio_view" {
			continue
		}
		// TODO: Use rs.Primary.ID to figure out if alert exists, and not just list all alerts.
		resp, err := conn.Views().List()
		if err == nil {
			if len(resp) > 0 {
				return fmt.Errorf("view still exist: %#+v", resp)
			}
		}
	}
	return nil
}

const viewEmpty = `
resource "humio_view" "test" {}
`

const viewInvalidInputs = `
resource "humio_view" "test" {
	name            = "invalid"
	repository		= "invlid"
	description     = "invalid"
}
`

const viewBasic = `
resource "humio_view" "test" {
	name            = "simple-view"
	description     = "a description"

	repository {
		name	= "sandbox"
		filter	= "*"
	}

	repository {
		name	= "sandbox-2"
		filter	= "test=test"
	}
}
`

var wantView = humio.View{
	Name:        "simple-view",
	Description: "a description",
	Connections: []humio.ViewConnection{{
		RepoName: "sandbox",
		Filter:   "*",
	}},
}

func TestEncodeDecodeViewResource(t *testing.T) {
	res := resourceView()
	data := res.TestResourceData()
	resourceDataFromView(&wantView, data)
	got, err := viewFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantView, got) {
		t.Error(cmp.Diff(wantView, got))
	}
}
