package opennebula

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/OpenNebula/one/src/oca/go/src/goca"
	ds "github.com/OpenNebula/one/src/oca/go/src/goca/schemas/datastore"
	dskeys "github.com/OpenNebula/one/src/oca/go/src/goca/schemas/datastore/keys"
	"github.com/OpenNebula/one/src/oca/go/src/goca/schemas/shared"
)

func TestAccVirtualMachine(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachineTemplateConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccSetDSdummy(),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "name", "test-virtual_machine"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "permissions", "642"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "memory", "128"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "cpu", "0.1"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "graphics.#", "1"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "graphics.0.keymap", "en-us"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "graphics.0.listen", "0.0.0.0"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "graphics.0.type", "VNC"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "os.#", "1"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "os.0.arch", "x86_64"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "os.0.boot", ""),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "tags.%", "2"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "tags.env", "prod"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "tags.customer", "test"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "timeout", "5"),
					resource.TestCheckResourceAttrSet("opennebula_virtual_machine.test", "uid"),
					resource.TestCheckResourceAttrSet("opennebula_virtual_machine.test", "gid"),
					resource.TestCheckResourceAttrSet("opennebula_virtual_machine.test", "uname"),
					resource.TestCheckResourceAttrSet("opennebula_virtual_machine.test", "gname"),
					testAccCheckVirtualMachinePermissions(&shared.Permissions{
						OwnerU: 1,
						OwnerM: 1,
						GroupU: 1,
						OtherM: 1,
					}),
				),
			},
			{
				Config: testAccVirtualMachineConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccSetDSdummy(),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "name", "test-virtual_machine-renamed"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "permissions", "660"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "group", "oneadmin"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "memory", "196"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "cpu", "0.2"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "graphics.#", "1"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "graphics.0.keymap", "en-us"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "graphics.0.listen", "0.0.0.0"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "graphics.0.type", "VNC"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "os.#", "1"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "os.0.arch", "x86_64"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "os.0.boot", ""),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "tags.%", "3"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "tags.env", "dev"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "tags.customer", "test"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "tags.version", "2"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "timeout", "4"),
					resource.TestCheckResourceAttrSet("opennebula_virtual_machine.test", "uid"),
					resource.TestCheckResourceAttrSet("opennebula_virtual_machine.test", "gid"),
					resource.TestCheckResourceAttrSet("opennebula_virtual_machine.test", "uname"),
					resource.TestCheckResourceAttrSet("opennebula_virtual_machine.test", "gname"),
					testAccCheckVirtualMachinePermissions(&shared.Permissions{
						OwnerU: 1,
						OwnerM: 1,
						OwnerA: 0,
						GroupU: 1,
						GroupM: 1,
					}),
				),
			},
		},
	})
}

func TestAccVirtualMachineDiskUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachineTemplateConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccSetDSdummy(),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "name", "test-virtual_machine"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "disk.#", "0"),
				),
			},
			{
				Config: testAccVirtualMachineTemplateConfigDisk,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "name", "test-virtual_machine"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "disk.#", "1"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "disk.0.target", "vda"),
				),
			},
			{
				Config: testAccVirtualMachineTemplateConfigDiskUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "name", "test-virtual_machine"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "disk.#", "1"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "disk.0.target", "vdb"),
				),
			},
			{
				Config: testAccVirtualMachineTemplateConfigDiskDetached,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "name", "test-virtual_machine"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "disk.#", "0"),
				),
			},
		},
	})
}

func TestAccVirtualMachinePending(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachinePending,
				Check: resource.ComposeTestCheckFunc(
					testAccSetDSdummy(),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "name", "virtual_machine_pending"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "permissions", "642"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "memory", "128"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "cpu", "0.1"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "pending", "true"),
					resource.TestCheckResourceAttrSet("opennebula_virtual_machine.test", "uid"),
					resource.TestCheckResourceAttrSet("opennebula_virtual_machine.test", "gid"),
					resource.TestCheckResourceAttrSet("opennebula_virtual_machine.test", "uname"),
					resource.TestCheckResourceAttrSet("opennebula_virtual_machine.test", "gname"),
					testAccCheckVirtualMachinePermissions(&shared.Permissions{
						OwnerU: 1,
						OwnerM: 1,
						GroupU: 1,
						OtherM: 1,
					}),
				),
			},
		},
	})
}
func testAccCheckVirtualMachineDestroy(s *terraform.State) error {
	controller := testAccProvider.Meta().(*goca.Controller)

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "opennebula_image":
			imgID, _ := strconv.ParseUint(rs.Primary.ID, 10, 64)
			imgc := controller.Image(int(imgID))
			// Get Virtual Machine Info
			img, _ := imgc.Info(false)
			if img != nil {
				imgState, _ := img.State()
				if imgState != 6 {
					return fmt.Errorf("Expected image %s to have been destroyed. imgState: %v", rs.Primary.ID, imgState)
				}
			}

		case "opennebula_virtual_machine":
			vmID, _ := strconv.ParseUint(rs.Primary.ID, 10, 64)
			vmc := controller.VM(int(vmID))
			// Get Virtual Machine Info
			vm, _ := vmc.Info(false)
			if vm != nil {
				vmState, _, _ := vm.State()
				if vmState != 6 {
					return fmt.Errorf("Expected virtual machine %s to have been destroyed. vmState: %v", rs.Primary.ID, vmState)
				}
			}
		default:
		}

	}

	return nil
}

func testAccSetDSdummy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if v := os.Getenv("TF_ACC_VM"); v == "1" {
			controller := testAccProvider.Meta().(*goca.Controller)

			dstpl := ds.NewTemplate()
			dstpl.Add(dskeys.TMMAD, "dummy")
			dstpl.Add(dskeys.DSMAD, "dummy")
			controller.Datastore(0).Update(dstpl.String(), 1)
			controller.Datastore(1).Update(dstpl.String(), 1)
		}
		return nil
	}
}

func testAccCheckVirtualMachinePermissions(expected *shared.Permissions) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		controller := testAccProvider.Meta().(*goca.Controller)

		for _, rs := range s.RootModule().Resources {
			switch rs.Type {
			case "opennebula_virtual_machine":

				vmID, _ := strconv.ParseUint(rs.Primary.ID, 10, 64)
				vmc := controller.VM(int(vmID))
				// Get Virtual Machine Info
				vm, err := vmc.Info(false)
				if vm == nil {
					return fmt.Errorf("Expected virtual_machine %s to exist when checking permissions: %s", rs.Primary.ID, err)
				}

				if !reflect.DeepEqual(vm.Permissions, expected) {
					return fmt.Errorf(
						"Permissions for virtual_machine %s were expected to be %s. Instead, they were %s",
						rs.Primary.ID,
						permissionsUnixString(*expected),
						permissionsUnixString(*vm.Permissions),
					)
				}
			default:
			}
		}

		return nil
	}
}

var testAccVirtualMachineTemplateConfigBasic = `
resource "opennebula_virtual_machine" "test" {
  name        = "test-virtual_machine"
  group       = "oneadmin"
  permissions = "642"
  memory = 128
  cpu = 0.1

  context = {
    NETWORK  = "YES"
    SET_HOSTNAME = "$NAME"
  }

  graphics {
    type   = "VNC"
    listen = "0.0.0.0"
    keymap = "en-us"
  }

  os {
    arch = "x86_64"
    boot = ""
  }

  tags = {
    env = "prod"
    customer = "test"
  }

  timeout = 5
}
`

var testAccVirtualMachineConfigUpdate = `
resource "opennebula_virtual_machine" "test" {
  name        = "test-virtual_machine-renamed"
  group       = "oneadmin"
  permissions = "660"
  memory = 196
  cpu = 0.2

  context = {
    NETWORK  = "YES"
    SET_HOSTNAME = "$NAME"
  }

  graphics {
    type   = "VNC"
    listen = "0.0.0.0"
    keymap = "en-us"
  }

  os {
    arch = "x86_64"
    boot = ""
  }

  tags = {
    env = "dev"
    customer = "test"
    version = "2"
  }
  timeout = 4
}
`

var testAccVirtualMachinePending = `
resource "opennebula_virtual_machine" "test" {
  name        = "virtual_machine_pending"
  group       = "oneadmin"
  permissions = "642"
  memory = 128
  cpu = 0.1
  pending = true

  context = {
    NETWORK  = "YES"
    SET_HOSTNAME = "$NAME"
  }

  graphics {
    type   = "VNC"
    listen = "0.0.0.0"
    keymap = "en-us"
  }

  os {
    arch = "x86_64"
    boot = ""
  }
}
`

var testAccVirtualMachineTemplateConfigDisk = `

resource "opennebula_image" "img1" {
  name             = "image1"
  type             = "DATABLOCK"
  size             = "16"
  datastore_id     = 1
  persistent       = false
  permissions      = "660"
}

resource "opennebula_image" "img2" {
  name             = "image2"
  type             = "DATABLOCK"
  size             = "8"
  datastore_id     = 1
  persistent       = false
  permissions      = "660"
}

resource "opennebula_virtual_machine" "test" {
	name        = "test-virtual_machine"
	group       = "oneadmin"
	permissions = "642"
	memory = 128
	cpu = 0.1
  
	context = {
	  NETWORK  = "YES"
	  SET_HOSTNAME = "$NAME"
	}
  
	graphics {
	  type   = "VNC"
	  listen = "0.0.0.0"
	  keymap = "en-us"
	}
  
	os {
	  arch = "x86_64"
	  boot = ""
	}
  
	tags = {
	  env = "prod"
	  customer = "test"
	}

	disk {
		image_id = opennebula_image.img2.id
		target = "vda"
	}
  
	timeout = 5
}
`

var testAccVirtualMachineTemplateConfigDiskUpdate = `

resource "opennebula_image" "img1" {
	name             = "image1"
	type             = "DATABLOCK"
	size             = "16"
	datastore_id     = 1
	persistent       = false
	permissions      = "660"
  }
  
  resource "opennebula_image" "img2" {
	name             = "image2"
	type             = "DATABLOCK"
	size             = "8"
	datastore_id     = 1
	persistent       = false
	permissions      = "660"
  }
  
  resource "opennebula_virtual_machine" "test" {
	  name        = "test-virtual_machine"
	  group       = "oneadmin"
	  permissions = "642"
	  memory = 128
	  cpu = 0.1
	
	  context = {
		NETWORK  = "YES"
		SET_HOSTNAME = "$NAME"
	  }
	
	  graphics {
		type   = "VNC"
		listen = "0.0.0.0"
		keymap = "en-us"
	  }
	
	  os {
		arch = "x86_64"
		boot = ""
	  }
	
	  tags = {
		env = "prod"
		customer = "test"
	  }
  
	  disk {
		  image_id = opennebula_image.img1.id
		  target = "vdb"
	  }
	
	  timeout = 5
}
`

var testAccVirtualMachineTemplateConfigDiskDetached = `

resource "opennebula_image" "img1" {
  name             = "image1"
  type             = "DATABLOCK"
  size             = "16"
  datastore_id     = 1
  persistent       = false
  permissions      = "660"
}

resource "opennebula_image" "img2" {
  name             = "image2"
  type             = "DATABLOCK"
  size             = "8"
  datastore_id     = 1
  persistent       = false
  permissions      = "660"
}

resource "opennebula_virtual_machine" "test" {
	name        = "test-virtual_machine"
	group       = "oneadmin"
	permissions = "642"
	memory = 128
	cpu = 0.1
  
	context = {
	  NETWORK  = "YES"
	  SET_HOSTNAME = "$NAME"
	}
  
	graphics {
	  type   = "VNC"
	  listen = "0.0.0.0"
	  keymap = "en-us"
	}
  
	os {
	  arch = "x86_64"
	  boot = ""
	}
  
	tags = {
	  env = "prod"
	  customer = "test"
	}
  
	timeout = 5
}
`
