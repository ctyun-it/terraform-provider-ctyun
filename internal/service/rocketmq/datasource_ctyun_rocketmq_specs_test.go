package rocketmq_test

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunRocketmqSpecs(t *testing.T) {
	t.Parallel()
	dnd := utils.GenerateRandomString()

	//datasourceName := "data.ctyun_rocketmq_specs." + dnd
	datasourceFile := "datasource_ctyun_rocketmq_specs.tf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: utils.LoadTestCase(datasourceFile, dnd),
			},
		},
	})
}
