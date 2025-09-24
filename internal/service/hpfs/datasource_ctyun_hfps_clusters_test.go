package hpfs_test

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccCtyunHpfsClusters(t *testing.T) {

	dnd := utils.GenerateRandomString()
	datasourceName := "data.ctyun_hpfs_clusters." + dnd
	datasourceFile := "datasource_ctyun_hpfs_clusters.tf"

	sfsType := "hpfs_perf"
	azName := "cn-huadong1-jsnj1A-public-ctcloud"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			//查询datasource
			{
				Config: utils.LoadTestCase(datasourceFile, dnd, sfsType, azName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "hpfs_clusters.0.sfs_type", sfsType),
					resource.TestCheckResourceAttr(datasourceName, "hpfs_clusters.0.az_name", azName),
				),
			},
		},
	})
}
