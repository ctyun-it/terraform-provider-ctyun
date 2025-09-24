package scaling_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunScalingConfig(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_scaling_config." + rnd
	resourceFile := "resource_ctyun_scaling_config_floating.tf"
	resourceFile1 := "resource_ctyun_scaling_config_disable_floating.tf"

	datasourceName := "data.ctyun_scaling_configs." + dnd
	datasourceFile := "datasource_ctyun_scaling_configs.tf"

	name := "sc-" + utils.GenerateRandomString()
	imageID := dependence.imageID
	flavorName := "s7.large.2"
	useFloating := "auto"
	bandwidth := 1
	loginMode := "password"
	password := "Kyk136@" + utils.GenerateRandomString()
	monitorService := true
	azNames := fmt.Sprintf(`["%s", "%s"]`, "cn-huadong1-jsnj1A-public-ctcloud", "cn-huadong1-jsnj2A-public-ctcloud")
	//azNames := fmt.Sprintf(`["%s"]`, "cn-huadong1-jsnj1A-public-ctcloud")
	tags := fmt.Sprintf(`[{"key":"%s", "value":"%s"}, {"key":"%s", "value":"%s"}]`, "provider", "scaling_conifg", "version", "1.1.1")
	volumes := fmt.Sprintf(`[{"volume_type":"%s", "volume_size":%d, "flag":"%s"}, {"volume_type":"%s", "volume_size":%d, "flag":"%s"}]`,
		"SATA", 40, "OS", "SAS", 100, "DATA")

	updatedName := "scn-" + utils.GenerateRandomString()
	updatedImageID := dependence.imageID1
	updatedFlavorName := "s8e.large.2"
	updateduseFloating := "diable"
	updatedLoginMode := "key_pair"
	keyPairId := dependence.keyPairID
	updatedMonitorService := false
	updatedTags := fmt.Sprintf(`[ {"key":"%s", "value":"%s"}]`, "version", "1.1.1")
	updatedVolumes := fmt.Sprintf(`[{"volume_type":"%s", "volume_size":%d, "flag":"%s"}]`, "SATA", 40, "OS")
	updatedAzName := fmt.Sprintf(`["%s"]`, "cn-huadong1-jsnj2A-public-ctcloud")
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
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, imageID, flavorName, useFloating, bandwidth, loginMode, password, monitorService, azNames, tags, volumes),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "image_id", imageID),
					resource.TestCheckResourceAttr(resourceName, "bandwidth", fmt.Sprintf("%d", bandwidth)),
					resource.TestCheckResourceAttr(resourceName, "password", password),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, imageID, flavorName, useFloating, bandwidth, loginMode, password, monitorService, azNames, tags, volumes) +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf("%s.id", resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "scaling_config_list.0.name", name),
					resource.TestCheckResourceAttr(datasourceName, "scaling_config_list.0.flavor_name", flavorName),
					resource.TestCheckResourceAttr(datasourceName, "scaling_config_list.0.bandwidth", fmt.Sprintf("%d", bandwidth)),
				),
			},
			// 更新
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, updatedName, updatedImageID, updatedFlavorName, updateduseFloating, updatedLoginMode, keyPairId, updatedMonitorService,
					updatedAzName, updatedTags, updatedVolumes),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "image_id", updatedImageID),
					resource.TestCheckResourceAttr(resourceName, "key_pair_id", keyPairId),
					resource.TestCheckResourceAttr(resourceName, "monitor_service", fmt.Sprintf("%t", updatedMonitorService)),
				),
			},
			//  资源导入测试
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					id := ds.ID
					regionId := ds.Attributes["region_id"]
					if id == "" || regionId == "" {
						return "", fmt.Errorf("id or region_id is required")
					}
					return fmt.Sprintf("%s,%s", id, regionId), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_pair_id"},
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, name, imageID, flavorName, useFloating, bandwidth, loginMode, password, monitorService, azNames, tags, volumes),
				Destroy: true,
			},
		},
	})

}
