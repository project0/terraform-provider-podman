package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceVolume_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccResourceVolumeConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("podman_volume.test", "name", "one"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "podman_volume.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"configurable_attribute"},
			},
			// Update and Read testing
			{
				Config: testAccResourceVolumeConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("podman_volume.test", "name", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccResourceVolume_local(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccResourceVolumeConfigFull("local", "type", "tmpfs"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("podman_volume.test", "driver", "local"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "podman_volume.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccResourceVolumeConfigFull("local", "o", "noexec"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("podman_volume.test", "driver", "local"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccResourceVolumeConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "podman_volume" "test" {
  name = %[1]q
}
`, configurableAttribute)
}

func testAccResourceVolumeConfigFull(driver, optkey, optvalue string) string {
	return fmt.Sprintf(`
resource "podman_volume" "test" {
	driver = %[1]q
	options = {
     %[2]q = %[3]q
	}
}
`, driver, optkey, optvalue)
}
