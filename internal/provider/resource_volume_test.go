package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceVolume_basic(t *testing.T) {
	name1 := generateResourceName()
	name2 := generateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccResourceVolumeConfig(name1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("podman_volume.test", "name", name1),
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
				Config: testAccResourceVolumeConfig(name2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("podman_volume.test", "name", name2),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccResourceVolume_local(t *testing.T) {
	name1 := generateResourceName()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccResourceVolumeConfigFull(name1, "local", "type", "tmpfs"),
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
				Config: testAccResourceVolumeConfigFull(name1, "local", "o", "noexec"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("podman_volume.test", "driver", "local"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccResourceVolumeConfig(name string) string {
	return fmt.Sprintf(`
resource "podman_volume" "test" {
  name = %[1]q
}
`, name)
}

func testAccResourceVolumeConfigFull(name, driver, optkey, optvalue string) string {
	return fmt.Sprintf(`
resource "podman_volume" "test" {
	name = %[1]q

	driver = %[1]q
	options = {
     %[2]q = %[3]q
	}
}
`, driver, optkey, optvalue)
}
