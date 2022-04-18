package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceNetwork_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccResourceNetwork(testName("one")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("podman_network.test", "name", testName("one")),
					resource.TestCheckResourceAttr("podman_network.test", "driver", "bridge"),
					resource.TestCheckResourceAttr("podman_network.test", "internal", "false"),
					resource.TestCheckResourceAttr("podman_network.test", "dns", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "podman_network.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccResourceNetwork(testName("two")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("podman_network.test", "name", testName("two")),
					resource.TestCheckResourceAttr("podman_network.test", "driver", "bridge"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccResourceNetwork_dualStack(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccResourceNetworkDualStack(testName("dual-1")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("podman_network.dualstack", "name", testName("dual-1")),
					resource.TestCheckResourceAttr("podman_network.dualstack", "driver", "bridge"),
					resource.TestCheckResourceAttr("podman_network.dualstack", "internal", "false"),
					resource.TestCheckResourceAttr("podman_network.dualstack", "dns", "true"),
					resource.TestCheckResourceAttr("podman_network.dualstack", "ipv6", "true"),
					resource.TestCheckTypeSetElemNestedAttrs("podman_network.dualstack", "subnets.*",
						map[string]string{
							"subnet":  "2001:db8::/64",
							"gateway": "2001:db8::1",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs("podman_network.dualstack", "subnets.*",
						map[string]string{
							"subnet":  "192.0.2.0/24",
							"gateway": "192.0.2.1",
						},
					),
				),
			},
			// ImportState testing
			{
				ResourceName:      "podman_network.dualstack",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccResourceNetworkDualStack(testName("dual-1")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("podman_network.dualstack", "name", testName("dual-1")),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccResourceNetwork(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "podman_network" "test" {
  name = %[1]q
}
`, configurableAttribute)
}

func testAccResourceNetworkDualStack(configurableAttribute string) string {
	return fmt.Sprintf(`
	resource "podman_network" "dualstack" {
		name        =  %[1]q
		driver      = "bridge"
		ipam_driver = "dhcp"
		options = {
			mtu = 1500
		}
		internal = false
		dns      = true
		# enable dual stack
		ipv6 = true
		subnets = [
			{
				subnet  = "2001:db8::/64"
				gateway = "2001:db8::1"
			},
			{
				subnet  = "192.0.2.0/24"
				gateway = "192.0.2.1"
			}
		]
	}
`, configurableAttribute)
}
