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
	"context"
	"fmt"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	humio "github.com/humio/cli/api"
)

func resourceView() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceViewCreate,
		ReadContext:   resourceViewRead,
		UpdateContext: resourceViewUpdate,
		DeleteContext: resourceViewDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"filter": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceViewCreate(ctx context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	view, err := viewFromResourceData(d)

	if err != nil {
		return diag.Errorf("count not obtain view from resource data: %s", err)
	}

	client.(*humio.Client).Views().Create(view.Name, view.Description, []humio.ViewConnectionInput{{
		RepositoryName: graphql.String(view.Connections[0].RepoName),
		Filter:         graphql.String(view.Connections[0].Filter),
	}})

	return resourceViewRead(ctx, d, client)
}

func resourceViewRead(_ context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {

	// we can't find a view without a name
	if _, ok := d.GetOk("name"); !ok {
		return diag.Errorf("must specify name for view")
	}

	name := d.Get("name").(string)

	view, err := client.(*humio.Client).Views().Get(name)
	if err != nil {
		return diag.Errorf("Unable to find view %s: %s", name, err)
	}

	return resourceDataFromView(view, d)
}

func resourceDataFromView(a *humio.View, d *schema.ResourceData) diag.Diagnostics {
	err := d.Set("name", a.Name)
	if err != nil {
		return diag.Errorf("error setting view name for resource %s: %s", d.Id(), err)
	}
	err = d.Set("description", a.Description)
	if err != nil {
		return diag.Errorf("error setting view description for resource %s: %s", d.Id(), err)
	}

	if len(a.Connections) == 0 {
		return diag.Errorf("must specify at least one connection for view %s", a.Name)
	}

	var repositories []interface{}
	for _, connection := range a.Connections {
		repo := make(map[string]interface{})
		repo["name"] = connection.RepoName
		repo["filter"] = connection.Filter
		repositories = append(repositories, repo)
	}

	err = d.Set("repository", repositories)
	if err != nil {
		return diag.Errorf("error setting view repository for resource %s: %s", d.Id(), err)
	}

	return nil
}

func resourceViewUpdate(ctx context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	view, err := viewFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain view from resource data: %s", err)
	}

	err = client.(*humio.Client).Views().UpdateDescription(view.Name, view.Description)
	if err != nil {
		return diag.Errorf("error updating view description %s: %s", d.Id(), err)
	}

	err = client.(*humio.Client).Views().UpdateConnections(view.Name, []humio.ViewConnectionInput{{
		RepositoryName: graphql.String(view.Connections[0].RepoName),
		Filter:         graphql.String(view.Connections[0].Filter),
	}})
	if err != nil {
		return diag.Errorf("error updating view connections: %s", err)
	}

	return resourceViewRead(ctx, d, client)
}

func resourceViewDelete(_ context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	view, err := viewFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain view from resource data: %s", err)
	}

	err = client.(*humio.Client).Views().Delete(view.Name, "Resource destruction from Terraform provider.")
	if err != nil {
		return diag.Errorf("error deleting view %s: %s", d.Id(), err)
	}

	return nil
}

func viewFromResourceData(d *schema.ResourceData) (humio.View, error) {

	_, ok := d.GetOk("repository")
	if !ok {
		return humio.View{}, fmt.Errorf("must specify repository")
	}

	connections, err := viewConnectionsFromResourceData(d)
	if err != nil {
		return humio.View{}, err
	}

	return humio.View{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Connections: connections,
	}, nil
}

func viewConnectionsFromResourceData(d *schema.ResourceData) ([]humio.ViewConnection, error) {
	repos := d.Get("repository").([]interface{})

	if len(repos) == 0 {
		return nil, fmt.Errorf("must specify at least one repository")
	}

	var connections []humio.ViewConnection

	for _, repo := range repos {
		repo := repo.(map[string]interface{})

		connections = append(connections, humio.ViewConnection{
			RepoName: repo["name"].(string),
			Filter:   repo["filter"].(string),
		})
	}

	return connections, nil
}
