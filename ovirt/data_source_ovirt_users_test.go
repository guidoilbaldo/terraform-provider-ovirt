// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOvirtUsersDataSource_nameRegexFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOvirtUsersDataSourceNameRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvirtDataSourceID("data.ovirt_users.name_regex_filtered_user"),
					resource.TestCheckResourceAttr("data.ovirt_users.name_regex_filtered_user", "users.#", "1"),
					resource.TestMatchResourceAttr("data.ovirt_users.name_regex_filtered_user", "users.0.name", regexp.MustCompile("^admin*")),
				),
			},
		},
	})
}

func TestAccOvirtUsersDataSource_searchFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOvirtUsersDataSourceSearchConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvirtDataSourceID("data.ovirt_users.search_filtered_user"),
					resource.TestCheckResourceAttr("data.ovirt_users.search_filtered_user", "users.#", "1"),
					resource.TestCheckResourceAttr("data.ovirt_users.search_filtered_user", "users.0.name", "admin"),
				),
			},
		},
	})
}

var testAccCheckOvirtUsersDataSourceNameRegexConfig = `
data "ovirt_users" "name_regex_filtered_user" {
	name_regex = "^admin*"
}
`

var testAccCheckOvirtUsersDataSourceSearchConfig = `
data "ovirt_users" "search_filtered_user" {
	search {
		max      = 1
		criteria = "name = admin"
	}
}
`
