package zos_test

import (
	"encoding/json"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

func TestAccCtyunZosBucket(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_zos_bucket." + rnd
	datasourceName := "data.ctyun_zos_buckets." + dnd
	resourceFile := "resource_ctyun_zos_bucket.tf"
	datasourceFile := "datasource_ctyun_zos_buckets.tf"

	bucket := "tf-bucket"
	acl := "public-read"
	azPolicy := "single-az"
	storageType := "STANDARD_IA"

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
				Config: utils.LoadTestCase(resourceFile, rnd, bucket, acl, azPolicy, storageType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bucket", bucket),
					resource.TestCheckResourceAttr(resourceName, "acl", acl),
					resource.TestCheckResourceAttr(resourceName, "az_policy", azPolicy),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "is_encrypted", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, bucket, acl, azPolicy, storageType) +
					utils.LoadTestCase(datasourceFile, dnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							ds := s.RootModule().Resources[datasourceName].Primary

							count, err := strconv.Atoi(ds.Attributes["buckets.#"])
							if err != nil || count == 0 {
								return fmt.Errorf("buckets 无效: %v", ds.Attributes)
							}

							for i := 0; i < count; i++ {
								if ds.Attributes[fmt.Sprintf("buckets.%d.bucket", i)] == bucket {
									if ds.Attributes[fmt.Sprintf("buckets.%d.storage_type", i)] != storageType {
										return fmt.Errorf("storage_type 不符合预期")
									}
									if ds.Attributes[fmt.Sprintf("buckets.%d.az_policy", i)] != azPolicy {
										return fmt.Errorf("az_policy 不符合预期")
									}
									return nil
								}
							}
							return fmt.Errorf("未找到目标元素")
						}),
				),
			},
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
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"project_id",
					"acl",
					"retention_mode",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, bucket, acl, azPolicy, storageType) +
					utils.LoadTestCase(datasourceFile, dnd),
				Destroy: true,
			},
		},
	},
	)
}

func TestAccCtyunZosBucketAllField(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_zos_bucket." + rnd
	datasourceName := "data.ctyun_zos_buckets." + dnd
	resourceFile := "resource_ctyun_zos_bucket_all_field.tf"
	datasourceFile := "datasource_ctyun_zos_buckets_b.tf"

	baseResourceFile := "resource_ctyun_zos_bucket_encrypted.tf"
	nologResourceFile := "resource_ctyun_zos_bucket_no_log.tf"

	bucket := "tf-bucket-all-field"
	azPolicy := "multi-az"
	storageType := "STANDARD"
	initAcl := "public-read"
	initTags := map[string]string{"c": "d"}
	initTagStr, _ := json.Marshal(initTags)
	initLogPrefix := "log/myfile"
	initLogBucket := dependence.bucket
	retention := "retention_year = 2"

	updatedAcl := "private"
	updatedTags := map[string]string{"a": "b"}
	updatedTagStr, _ := json.Marshal(updatedTags)
	updatedLogPrefix := "log/mylog"
	updatedLogBucket := bucket
	updatedRetention := "retention_day = 10"

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
				Config: utils.LoadTestCase(resourceFile, rnd, bucket, initAcl, azPolicy, storageType, initTagStr, initLogBucket, initLogPrefix, retention),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bucket", bucket),
					resource.TestCheckResourceAttr(resourceName, "acl", initAcl),
					resource.TestCheckResourceAttr(resourceName, "az_policy", azPolicy),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "is_encrypted", "true"),
					resource.TestCheckResourceAttr(resourceName, "version_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "log_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "log_prefix", initLogPrefix),
					resource.TestCheckResourceAttr(resourceName, "log_bucket", initLogBucket),
					resource.TestCheckResourceAttr(resourceName, "retention_year", "2"),
					func(s *terraform.State) error {
						obj, _ := s.RootModule().Resources[resourceName]
						c := obj.Primary.Attributes["tags.c"]
						if c != "d" {
							return fmt.Errorf("expected tag 'c' to be 'd'")
						}
						return nil
					},
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, bucket, initAcl, azPolicy, storageType, initTagStr, initLogBucket, initLogPrefix, retention) +
					utils.LoadTestCase(datasourceFile, dnd, bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "buckets.0.bucket", bucket),
					resource.TestCheckResourceAttr(datasourceName, "buckets.0.az_policy", azPolicy),
					resource.TestCheckResourceAttr(datasourceName, "buckets.0.storage_type", storageType),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, bucket, updatedAcl, azPolicy, storageType, updatedTagStr, updatedLogBucket, updatedLogPrefix, updatedRetention),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bucket", bucket),
					resource.TestCheckResourceAttr(resourceName, "acl", updatedAcl),
					resource.TestCheckResourceAttr(resourceName, "az_policy", azPolicy),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "is_encrypted", "true"),
					resource.TestCheckResourceAttr(resourceName, "version_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "log_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "log_prefix", updatedLogPrefix),
					resource.TestCheckResourceAttr(resourceName, "log_bucket", updatedLogBucket),
					resource.TestCheckResourceAttr(resourceName, "retention_day", "10"),
					func(s *terraform.State) error {
						obj, _ := s.RootModule().Resources[resourceName]
						a := obj.Primary.Attributes["tags.a"]
						if a != "b" {
							return fmt.Errorf("expected tag 'a' to be 'b'")
						}
						return nil
					},
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(nologResourceFile, rnd, bucket, updatedAcl, azPolicy, storageType, updatedTagStr, updatedRetention),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bucket", bucket),
					resource.TestCheckResourceAttr(resourceName, "acl", updatedAcl),
					resource.TestCheckResourceAttr(resourceName, "az_policy", azPolicy),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "is_encrypted", "true"),
					resource.TestCheckResourceAttr(resourceName, "version_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "log_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "retention_day", "10"),
					func(s *terraform.State) error {
						obj, _ := s.RootModule().Resources[resourceName]
						a := obj.Primary.Attributes["tags.a"]
						if a != "b" {
							return fmt.Errorf("expected tag 'a' to be 'b'")
						}
						return nil
					},
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(baseResourceFile, rnd, bucket, updatedAcl, azPolicy, storageType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bucket", bucket),
					resource.TestCheckResourceAttr(resourceName, "acl", updatedAcl),
					resource.TestCheckResourceAttr(resourceName, "az_policy", azPolicy),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "is_encrypted", "true"),
					resource.TestCheckResourceAttr(resourceName, "version_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "log_enabled", "false"),
				),
			},
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
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"project_id",
					"acl",
					"retention_mode",
				},
			},
			{
				Config:  utils.LoadTestCase(baseResourceFile, rnd, bucket, updatedAcl, azPolicy, storageType),
				Destroy: true,
			},
		},
	},
	)
}
