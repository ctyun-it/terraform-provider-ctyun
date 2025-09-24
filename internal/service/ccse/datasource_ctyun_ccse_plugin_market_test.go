package ccse_test

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCtyunCcsePluginMarket(t *testing.T) {
	datasourceFile := "datasource_ctyun_ccse_plugin_market.tf"
	dnd1 := utils.GenerateRandomString()
	dnd1Name := "data.ctyun_ccse_plugin_market." + dnd1
	dnd2 := utils.GenerateRandomString()
	dnd2Name := "data.ctyun_ccse_plugin_market." + dnd2
	dnd3 := utils.GenerateRandomString()
	dnd3Name := "data.ctyun_ccse_plugin_market." + dnd3
	valuesYaml := "YAML"
	valuesJson := "JSON"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: utils.LoadTestCase(datasourceFile, dnd1, dnd2, dnd3, valuesYaml),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dnd3Name, "values_type", valuesYaml),
					resource.TestCheckResourceAttrWith(dnd1Name, "records.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrWith(dnd2Name, "records.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrWith(dnd2Name, "versions.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrWith(dnd3Name, "records.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrWith(dnd3Name, "versions.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrSet(dnd3Name, "values"),
				),
			},
			{
				Config: utils.LoadTestCase(datasourceFile, dnd1, dnd2, dnd3, valuesJson),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dnd3Name, "values_type", valuesJson),
					resource.TestCheckResourceAttrWith(dnd1Name, "records.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrWith(dnd2Name, "records.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrWith(dnd2Name, "versions.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrWith(dnd3Name, "records.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrWith(dnd3Name, "versions.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrSet(dnd3Name, "values"),
				),
			},
		},
	})
}

func TestAccCtyunCcseTemplateMarket(t *testing.T) {
	datasourceFile := "datasource_ctyun_ccse_template_market.tf"
	dnd1 := utils.GenerateRandomString()
	dnd1Name := "data.ctyun_ccse_template_market." + dnd1
	dnd2 := utils.GenerateRandomString()
	dnd2Name := "data.ctyun_ccse_template_market." + dnd2
	dnd3 := utils.GenerateRandomString()
	dnd3Name := "data.ctyun_ccse_template_market." + dnd3
	valuesYaml := "YAML"
	valuesJson := "JSON"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: utils.LoadTestCase(datasourceFile, dnd1, dnd2, dnd3, valuesYaml),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dnd3Name, "values_type", valuesYaml),
					resource.TestCheckResourceAttrWith(dnd1Name, "records.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrWith(dnd2Name, "records.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrWith(dnd2Name, "versions.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrWith(dnd3Name, "records.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrWith(dnd3Name, "versions.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrSet(dnd3Name, "values"),
				),
			},
			{
				Config: utils.LoadTestCase(datasourceFile, dnd1, dnd2, dnd3, valuesJson),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dnd3Name, "values_type", valuesJson),
					resource.TestCheckResourceAttrWith(dnd1Name, "records.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrWith(dnd2Name, "records.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrWith(dnd2Name, "versions.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrWith(dnd3Name, "records.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrWith(dnd3Name, "versions.#", utils.AtLeastOne),
					resource.TestCheckResourceAttrSet(dnd3Name, "values"),
				),
			},
		},
	})
}
