package zos_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

func TestAccCtyunZosBucketObject(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_zos_bucket_object." + rnd
	datasourceName := "data.ctyun_zos_bucket_objects." + dnd
	resourceFile := "resource_ctyun_zos_bucket_object.tf"
	datasourceFile := "datasource_ctyun_zos_bucket_objects.tf"

	key := "example/tf.txt"
	source := "./testdata/tf.txt"

	initAcl := "private"
	initTags := `"a": "b"`

	updatedAcl := "public-read-write"
	updatedTags := `"c": "d"`

	nextAcl := "public-read-write"
	nextTags := ``
	dependenceBucket := dependence.bucket
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
				Config: utils.LoadTestCase(resourceFile, rnd, dependenceBucket, key, source, initAcl, initTags),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bucket", dependenceBucket),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "acl", initAcl),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					func(s *terraform.State) error {
						obj, _ := s.RootModule().Resources[resourceName]
						a := obj.Primary.Attributes["tags.a"]
						if a != "b" {
							return fmt.Errorf("expected tag 'a' to be 'b'")
						}
						return nil
					},
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependenceBucket, key, source, updatedAcl, updatedTags),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bucket", dependenceBucket),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "acl", updatedAcl),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					func(s *terraform.State) error {
						obj, _ := s.RootModule().Resources[resourceName]
						c := obj.Primary.Attributes["tags.c"]
						if c != "d" {
							return fmt.Errorf("expected tag 'c' to be 'd'")
						}
						return nil
					},
				),
			},

			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependenceBucket, key, source, nextAcl, nextTags),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bucket", dependenceBucket),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "acl", nextAcl),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					func(s *terraform.State) error {
						obj, _ := s.RootModule().Resources[resourceName]
						l := obj.Primary.Attributes["tags.%"]
						if l != "0" {
							return fmt.Errorf("expected tag is empty")
						}
						return nil
					},
				),
			},

			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependenceBucket, key, source, nextAcl, nextTags) +
					utils.LoadTestCase(datasourceFile, dnd, dependenceBucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							ds := s.RootModule().Resources[datasourceName].Primary
							count, err := strconv.Atoi(ds.Attributes["objects.#"])
							if err != nil || count == 0 {
								return fmt.Errorf("objects 无效: %v", ds.Attributes)
							}

							for i := 0; i < count; i++ {
								if ds.Attributes[fmt.Sprintf("objects.%d.key", i)] == key {
									return nil
								}
							}
							return fmt.Errorf("未找到目标元素")
						},
					)),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"acl",
					"tags",
					"source",
					"content",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependenceBucket, key, source, nextAcl, nextTags) +
					utils.LoadTestCase(datasourceFile, dnd, dependenceBucket),
				Destroy: true,
			},
		},
	},
	)
}

func TestAccCtyunZosBucketObjectAllField(t *testing.T) {
	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_zos_bucket_object." + rnd
	datasourceName := "data.ctyun_zos_bucket_objects." + dnd
	resourceFile := "resource_ctyun_zos_bucket_object_all_field.tf"
	datasourceFile := "datasource_ctyun_zos_bucket_objects.tf"

	key := "example/tf.txt"
	content := "abc"

	initAcl := "private"
	initTags := `"a": "b"`

	updatedAcl := "public-read-write"
	updatedTags := `"c": "d"`

	nextAcl := "public-read-write"
	nextTags := ``
	dependenceBucket := dependence.bucket
	storageType := "STANDARD_IA"
	cacheControl := "no-cache"
	contentEncoding := "identity"
	contentType := "plain/text"
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
				Config: utils.LoadTestCase(resourceFile, rnd, dependenceBucket, key, content, initAcl, initTags, storageType, cacheControl, contentEncoding, contentType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bucket", dependenceBucket),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "acl", initAcl),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "cache_control", cacheControl),
					resource.TestCheckResourceAttr(resourceName, "content_encoding", contentEncoding),
					resource.TestCheckResourceAttr(resourceName, "content_type", contentType),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					func(s *terraform.State) error {
						obj, _ := s.RootModule().Resources[resourceName]
						a := obj.Primary.Attributes["tags.a"]
						if a != "b" {
							return fmt.Errorf("expected tag 'a' to be 'b'")
						}
						return nil
					},
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependenceBucket, key, content, updatedAcl, updatedTags, storageType, cacheControl, contentEncoding, contentType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bucket", dependenceBucket),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "acl", updatedAcl),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "cache_control", cacheControl),
					resource.TestCheckResourceAttr(resourceName, "content_encoding", contentEncoding),
					resource.TestCheckResourceAttr(resourceName, "content_type", contentType),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					func(s *terraform.State) error {
						obj, _ := s.RootModule().Resources[resourceName]
						c := obj.Primary.Attributes["tags.c"]
						if c != "d" {
							return fmt.Errorf("expected tag 'c' to be 'd'")
						}
						return nil
					},
				),
			},

			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependenceBucket, key, content, nextAcl, nextTags, storageType, cacheControl, contentEncoding, contentType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bucket", dependenceBucket),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "acl", nextAcl),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "cache_control", cacheControl),
					resource.TestCheckResourceAttr(resourceName, "content_encoding", contentEncoding),
					resource.TestCheckResourceAttr(resourceName, "content_type", contentType),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					func(s *terraform.State) error {
						obj, _ := s.RootModule().Resources[resourceName]
						l := obj.Primary.Attributes["tags.%"]
						if l != "0" {
							return fmt.Errorf("expected tag is empty")
						}
						return nil
					},
				),
			},

			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependenceBucket, key, content, nextAcl, nextTags, storageType, cacheControl, contentEncoding, contentType) +
					utils.LoadTestCase(datasourceFile, dnd, dependenceBucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						func(s *terraform.State) error {
							ds := s.RootModule().Resources[datasourceName].Primary
							count, err := strconv.Atoi(ds.Attributes["objects.#"])
							if err != nil || count == 0 {
								return fmt.Errorf("objects 无效: %v", ds.Attributes)
							}

							for i := 0; i < count; i++ {
								if ds.Attributes[fmt.Sprintf("objects.%d.key", i)] == key {
									return nil
								}
							}
							return fmt.Errorf("未找到目标元素")
						},
					)),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"acl",
					"tags",
					"source",
					"content",
				},
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, dependenceBucket, key, content, nextAcl, nextTags, storageType, cacheControl, contentEncoding, contentType) +
					utils.LoadTestCase(datasourceFile, dnd, dependenceBucket),
				Destroy: true,
			},
		},
	},
	)
}
