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

func TestAccOvirtVnic_basic(t *testing.T) {
	var nic ovirtsdk4.Nic
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVnicDestroy,
		IDRefreshName: "ovirt_vnic.nic",
		Steps: []resource.TestStep{
			{
				Config: testAccVnicBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvirtVnicExists("ovirt_vnic.nic", &nic),
					resource.TestCheckResourceAttr("ovirt_vnic.nic", "name", "testAccOvirtVnicBasic"),
					resource.TestCheckResourceAttr("ovirt_vnic.nic", "vm_id", "1a4bc4d8-fec7-4fe4-b01a-7d1185854c39"),
					resource.TestCheckResourceAttr("ovirt_vnic.nic", "vnic_profile_id", "0000000a-000a-000a-000a-000000000398"),
				),
			},
			{
				Config: testAccVnicBasicUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvirtVnicExists("ovirt_vnic.nic", &nic),
					resource.TestCheckResourceAttr("ovirt_vnic.nic", "name", "testAccOvirtVnicBasicUpdate"),
					resource.TestCheckResourceAttr("ovirt_vnic.nic", "vm_id", "77f7e0d9-6105-492f-92e8-06b989211e46"),
					resource.TestCheckResourceAttr("ovirt_vnic.nic", "vnic_profile_id", "0000000a-000a-000a-000a-000000000398"),
				),
			},
		},
	})
}

func testAccCheckVnicDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ovirtsdk4.Connection)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovirt_vnic" {
			continue
		}

		vmID, nicID, err := getVMIDAndNicID(rs.Primary.ID)
		if err != nil {
			return err
		}

		getResp, err := conn.SystemService().VmsService().
			VmService(vmID).
			NicsService().
			NicService(nicID).
			Get().
			Send()
		if err != nil {
			if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
				continue
			}
			return err
		}
		if _, ok := getResp.Nic(); ok {
			return fmt.Errorf("Vnic %s still exist", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckOvirtVnicExists(n string, v *ovirtsdk4.Nic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Vnic ID is set")
		}

		vmID, nicID, err := getVMIDAndNicID(rs.Primary.ID)
		if err != nil {
			return err
		}

		conn := testAccProvider.Meta().(*ovirtsdk4.Connection)
		getResp, err := conn.SystemService().VmsService().
			VmService(vmID).
			NicsService().
			NicService(nicID).
			Get().
			Send()
		if err != nil {
			return err
		}
		nic, ok := getResp.Nic()
		if ok {
			*v = *nic
			return nil
		}
		return fmt.Errorf("Vnic %s not exist", rs.Primary.ID)
	}
}

const testAccVnicBasic = `
resource "ovirt_vnic" "nic" {
	name        	= "testAccOvirtVnicBasic"
	vm_id			= "1a4bc4d8-fec7-4fe4-b01a-7d1185854c39"
	vnic_profile_id = "0000000a-000a-000a-000a-000000000398"
}
`

const testAccVnicBasicUpdate = `
resource "ovirt_vnic" "nic" {
	name        	= "testAccOvirtVnicBasicUpdate"
	vm_id			= "77f7e0d9-6105-492f-92e8-06b989211e46"
	vnic_profile_id = "0000000a-000a-000a-000a-000000000398"
}
`
