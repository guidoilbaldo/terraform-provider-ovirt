// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ovirtsdk4 "gopkg.in/imjoey/go-ovirt.v4"
)

func TestAccOvirtHost_basic(t *testing.T) {
	var host ovirtsdk4.Host
	clusterID, updateClusterID := "5bc08e5c-00ef-01e3-01dd-0000000001df", "5bc08e5c-00ef-01e3-01dd-0000000001df"
	address, updateAddress := "10.1.110.64", "10.1.110.64"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		Providers:     testAccProviders,
		IDRefreshName: "ovirt_host.host",
		CheckDestroy:  testAccCheckHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHostBasic(address, clusterID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvirtHostExists("ovirt_host.host", &host),
					resource.TestCheckResourceAttr("ovirt_host.host", "name", "host64"),
					resource.TestCheckResourceAttr("ovirt_host.host", "address", address),
				),
			},
			{
				Config: testAccHostBasicUpdate(updateAddress, updateClusterID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvirtHostExists("ovirt_host.host", &host),
					resource.TestCheckResourceAttr("ovirt_host.host", "name", "host64"),
					resource.TestCheckResourceAttr("ovirt_host.host", "address", updateAddress),
				),
			},
		},
	})
}

func testAccCheckHostDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ovirtsdk4.Connection)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovirt_host" {
			continue
		}
		getResp, err := conn.SystemService().HostsService().
			HostService(rs.Primary.ID).
			Get().
			Send()
		if err != nil {
			if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
				continue
			}
			return err
		}
		if _, ok := getResp.Host(); ok {
			return fmt.Errorf("Host %s still exist", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckOvirtHostExists(n string, v *ovirtsdk4.Host) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Host ID is set")
		}
		conn := testAccProvider.Meta().(*ovirtsdk4.Connection)
		getResp, err := conn.SystemService().HostsService().
			HostService(rs.Primary.ID).
			Get().
			Send()
		if err != nil {
			return err
		}
		host, ok := getResp.Host()
		if ok {
			*v = *host
			return nil
		}
		return fmt.Errorf("Host %s not exist", rs.Primary.ID)
	}
}

func testAccHostBasic(address, clusterID string) string {
	return fmt.Sprintf(`
resource "ovirt_host" "host" {
	name        = "host64"
	description = "my new host"
	address		= "%s"
	root_password = "qwer1234"
	cluster_id  = "%s"
}`, address, clusterID)
}

func testAccHostBasicUpdate(address, clusterID string) string {
	return fmt.Sprintf(`
resource "ovirt_host" "host" {
	name        = "host64"
	description = "my updated new host"
	address		= "%s"
	root_password = "qwer1234"
	cluster_id  = "%s"
}`, address, clusterID)
}
