package ecs_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunKeyPair(t *testing.T) {
	rnd := utils.GenerateRandomString()

	resourceName := "ctyun_keypair." + rnd

	resourceFile := "resource_ctyun_keypair.tf"

	keyName := "tf-keypair-" + rnd
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjUnAnTid4wmVtajSmElMtH03OvOyY81ybfswbUu9Gt83DVVzDnwb3rcQW1us8SeKm/gRINkgdrRAgfXAmTKR7AorYtWWc/tzb6kcDpL2E8Qk+n6cyFAxXNoX2vXBr4kC9wz1uwjGyxoSlpHLIpscfI0Ef652gMlSyfODehAJHj3JPMr8pvtPIUqsZI3JOGTUzxaA2JVC0LxQegphYYf2TxGd9GLRUv1p/0BUAPCMg1NaITXNVEj3A11hk1nrFoJMmvIwIUkLmRuQcxuNAdxeLB7GXXVjKpnKIJL4L64dyA9GWa3Gb7gCJyRaBc5UhK4hT57wmukCrldHHtdF1IJr"

	resource.Test(t, resource.TestCase{
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 创建
			{
				Config: utils.LoadTestCase(resourceFile, rnd, keyName, publicKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", keyName),
					resource.TestCheckResourceAttr(resourceName, "public_key", publicKey),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},

			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					name := ds.Attributes["name"]
					regionId := ds.Attributes["region_id"]
					return fmt.Sprintf("%s,%s", name, regionId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"project_id",
				},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, keyName, publicKey),
				Destroy: true,
			},
		},
	})
}
