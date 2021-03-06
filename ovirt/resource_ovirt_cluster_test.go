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

func TestAccOvirtCluster_basic(t *testing.T) {
	datacenterID := "5bc08e5b-03ab-0194-03cb-000000000289"
	networkID := "00000000-0000-0000-0000-000000000009"
	var cluster ovirtsdk4.Cluster
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		Providers:     testAccProviders,
		IDRefreshName: "ovirt_cluster.cluster",
		CheckDestroy:  testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterBasic(datacenterID, networkID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvirtClusterExists("ovirt_cluster.cluster", &cluster),
					resource.TestCheckResourceAttr("ovirt_cluster.cluster", "name", "testAccOvirtClusterBasic"),
					resource.TestCheckResourceAttr("ovirt_cluster.cluster", "datacenter_id", datacenterID),
					resource.TestCheckResourceAttr("ovirt_cluster.cluster", "management_network_id", networkID),
				),
			},
			{
				Config: testAccClusterBasicUpdate(datacenterID, networkID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvirtClusterExists("ovirt_cluster.cluster", &cluster),
					resource.TestCheckResourceAttr("ovirt_cluster.cluster", "name", "testAccOvirtClusterBasicUpdate"),
					resource.TestCheckResourceAttr("ovirt_cluster.cluster", "datacenter_id", datacenterID),
					resource.TestCheckResourceAttr("ovirt_cluster.cluster", "management_network_id", networkID),
					// resource.TestCheckNoResourceAttr("ovirt_cluster.cluster", "description"),
					resource.TestCheckResourceAttr("ovirt_cluster.cluster", "description", ""),
				),
			},
		},
	})
}

func testAccCheckClusterDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ovirtsdk4.Connection)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovirt_cluster" {
			continue
		}
		getResp, err := conn.SystemService().ClustersService().
			ClusterService(rs.Primary.ID).
			Get().
			Send()
		if err != nil {
			if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
				continue
			}
			return err
		}
		if _, ok := getResp.Cluster(); ok {
			return fmt.Errorf("Cluster %s still exist", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckOvirtClusterExists(n string, v *ovirtsdk4.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cluster ID is set")
		}
		conn := testAccProvider.Meta().(*ovirtsdk4.Connection)
		getResp, err := conn.SystemService().ClustersService().
			ClusterService(rs.Primary.ID).
			Get().
			Send()
		if err != nil {
			return err
		}
		cluster, ok := getResp.Cluster()
		if ok {
			*v = *cluster
			return nil
		}
		return fmt.Errorf("Cluster %s not exist", rs.Primary.ID)
	}
}

func testAccClusterBasic(datacenterID, networkID string) string {
	return fmt.Sprintf(`
resource "ovirt_cluster" "cluster" {
	name	    						= "testAccOvirtClusterBasic"
	description 						= "Desc of cluster"
	datacenter_id						= "%s"
	management_network_id				= "%s"
	memory_policy_over_commit_percent   = 100
	ballooning							= true
	gluster								= true
	threads_as_cores					= true
	cpu_arch							= "x86_64"
	cpu_type							= "Intel SandyBridge Family"
	compatibility_version				= "4.1"
}`, datacenterID, networkID)
}

func testAccClusterBasicUpdate(datacenterID, networkID string) string {
	return fmt.Sprintf(`
resource "ovirt_cluster" "cluster" {
	name	    						= "testAccOvirtClusterBasicUpdate"
	datacenter_id						= "%s"
	management_network_id				= "%s"
	memory_policy_over_commit_percent   = 100
	ballooning							= true
	gluster								= true
	threads_as_cores					= true
	cpu_arch							= "x86_64"
	cpu_type							= "Intel SandyBridge Family"
	compatibility_version				= "4.1"
}`, datacenterID, networkID)
}
