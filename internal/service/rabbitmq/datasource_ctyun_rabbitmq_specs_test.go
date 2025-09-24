package rabbitmq_test

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunRabbitmqSpecs(t *testing.T) {
	dnd := utils.GenerateRandomString()

	datasourceName := "data.ctyun_rabbitmq_specs." + dnd
	datasourceFile := "datasource_ctyun_rabbitmq_specs.tf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{

				Config: utils.LoadTestCase(datasourceFile, dnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(datasourceName, "specs.#", utils.AtLeastOne),
				),
			},
		},
	})
}
