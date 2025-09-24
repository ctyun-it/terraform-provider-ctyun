package kafka_test

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunKafkaSpecs(t *testing.T) {
	dnd := utils.GenerateRandomString()

	datasourceName := "data.ctyun_kafka_specs." + dnd
	datasourceFile := "datasource_ctyun_kafka_specs.tf"

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
