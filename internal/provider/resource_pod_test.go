package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePod_basic(t *testing.T) {
	name1 := generateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccResourcePod(name1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("podman_pod.test", "name", name1),
				),
			},
			// ImportState testing
			{
				ResourceName:      "podman_pod.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccResourcePod(name1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("podman_pod.test", "name", name1),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccResourcePod_volume(t *testing.T) {

	name1 := generateResourceName()
	name2 := generateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccResourcePodVolume(name1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("podman_pod.test", "name", name1),
				),
			},
			// ImportState testing
			{
				ResourceName:      "podman_pod.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccResourcePodVolume(name1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("podman_pod.test", "name", name1),
				),
			},

			// Test replace
			{
				Config: testAccResourcePodVolume(name2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("podman_pod.test", "name", name2),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccResourcePod(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "podman_pod" "test" {
  name = %[1]q
}
`, configurableAttribute)
}

func testAccResourcePodVolume(name string) string {
	return fmt.Sprintf(`
resource "podman_volume" "test" {
	name = %[1]q
}

resource "podman_pod" "test" {
	name = %[1]q
  mounts = [
		{
    	destination = "/data/one"
    	volume = {
      	name = podman_volume.test.name
    	}
  	},
		{
    	destination = "/data/two"
    	volume = {
      	name = podman_volume.test.name
				read_only = true
				exec = true
				suid = false
				chown = false
				idmap = false
				dev = false
    	}
  	}
	]
}
`, name)
}
