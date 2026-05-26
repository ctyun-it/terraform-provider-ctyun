package faas_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strings"
	"testing"
)

func TestAccCtyunFunctionTrigger_Schedule(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_function_trigger." + rnd
	dataSourceName := "data.ctyun_function_triggers." + rnd

	triggerName := "trigger-schedule-" + rnd
	functionName := dependence.functionName

	resourceFile := "resource_ctyun_function_trigger_schedule.tf"
	dataSourceFile := "datasource_ctyun_function_triggers.tf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					triggerName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "function_name", functionName),
					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName),
					resource.TestCheckResourceAttr(resourceName, "trigger_type", "schedule"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrWith(resourceName, "event_data", func(value string) error {
						if !strings.Contains(value, "cronExpr") {
							return fmt.Errorf("event_data should contain cronExpr")
						}
						return nil
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					name := ds.Attributes["trigger_name"]
					function_name := ds.Attributes["function_name"]
					return fmt.Sprintf("%s,%s", name, function_name), nil
				},
				ImportStateVerifyIgnore: []string{"event_data"},
			},
			{
				Config: utils.LoadTestCase(dataSourceFile,
					rnd,
					functionName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "function_name", functionName),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					triggerName+"-update",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName+"-update"),
					resource.TestCheckResourceAttrWith(resourceName, "event_data", func(value string) error {
						if !strings.Contains(value, "cronExpr") {
							return fmt.Errorf("updated event_data should contain cronExpr")
						}
						return nil
					}),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, functionName, triggerName),
				Destroy: true,
			},
		},
	})
}

// 没有域名
func TestAccCtyunFunctionTrigger_Http(t *testing.T) {
	t.Parallel()
	rnd := utils.GenerateRandomString()
	resourceName := "ctyun_function_trigger." + rnd

	triggerName := "trigger-http-" + rnd
	functionName := dependence.functionName

	resourceFile := "resource_ctyun_function_trigger_http.tf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					triggerName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "function_name", functionName),
					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName),
					resource.TestCheckResourceAttr(resourceName, "trigger_type", "http"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds := s.RootModule().Resources[resourceName].Primary
					name := ds.Attributes["trigger_name"]
					function_name := ds.Attributes["function_name"]
					return fmt.Sprintf("%s,%s", name, function_name), nil
				},
				ImportStateVerifyIgnore: []string{"event_data"},
			},
			{
				Config: utils.LoadTestCase(resourceFile,
					rnd,
					functionName,
					triggerName+"-update",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName+"-update"),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, functionName, triggerName),
				Destroy: true,
			},
		},
	})
}

//func TestAccCtyunFunctionTrigger_Kafka(t *testing.T) {
//	t.Parallel()
//	rnd := utils.GenerateRandomString()
//	resourceName := "ctyun_function_trigger." + rnd
//
//	triggerName := "trigger-kafka-" + rnd
//	functionName :=dependence.functionName
//
//	regionId := os.Getenv("CTYUN_REGION_ID")
//	kafkaInstanceId := "7ddfb1f6c8b849d680141ec97d7cb31f"
//
//	resourceFile := "resource_ctyun_function_trigger_kafka.tf"
//
//	resource.Test(t, resource.TestCase{
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		CheckDestroy: func(s *terraform.State) error {
//			_, exists := s.RootModule().Resources[resourceName]
//			if exists {
//				return fmt.Errorf("resource destroy failed")
//			}
//			return nil
//		},
//		Steps: []resource.TestStep{
//			{
//				Config: utils.LoadTestCase(resourceFile,
//					rnd,
//					functionName,
//					triggerName,
//					regionId,
//					kafkaInstanceId,
//				),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "function_name", functionName),
//					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName),
//					resource.TestCheckResourceAttr(resourceName, "trigger_type", "kafka"),
//					resource.TestCheckResourceAttr(resourceName, "region_id", regionId),
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttrSet(resourceName, "status"),
//					resource.TestCheckResourceAttrWith(resourceName, "event_data", func(value string) error {
//						if !strings.Contains(value, kafkaInstanceId) {
//							return fmt.Errorf("event_data should contain instanceId")
//						}
//						return nil
//					}),
//				),
//			},
//			{
//				ResourceName:      resourceName,
//				ImportState:       true,
//				ImportStateVerify: true,
//				ImportStateIdFunc: func(s *terraform.State) (string, error) {
//					ds := s.RootModule().Resources[resourceName].Primary
//					region_id := ds.Attributes["region_id"]
//					name := ds.Attributes["trigger_name"]
//					function_name := ds.Attributes["function_name"]
//					return fmt.Sprintf("%s,%s,%s", name, function_name, region_id), nil
//				},
//				ImportStateVerifyIgnore: []string{"event_data"},
//			},
//			{
//				Config: utils.LoadTestCase(resourceFile,
//					rnd,
//					functionName,
//					triggerName+"-update",
//					regionId,
//					kafkaInstanceId,
//				),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName+"-update"),
//					resource.TestCheckResourceAttrWith(resourceName, "event_data", func(value string) error {
//						if !strings.Contains(value, kafkaInstanceId) {
//							return fmt.Errorf("updated event_data should contain instanceId")
//						}
//						return nil
//					}),
//				),
//			},
//			{
//				Config:  utils.LoadTestCase(resourceFile, rnd, functionName, triggerName, regionId, kafkaInstanceId),
//				Destroy: true,
//			},
//		},
//	})
//}
//
//func TestAccCtyunFunctionTrigger_Rocketmq(t *testing.T) {
//	t.Parallel()
//	rnd := utils.GenerateRandomString()
//	resourceName := "ctyun_function_trigger." + rnd
//
//	triggerName := "trigger-rocketmq-" + rnd
//	functionName :=dependence.functionName
//
//	regionId := os.Getenv("CTYUN_REGION_ID")
//	rocketmqInstanceId := os.Getenv("CTYUN_ROCKETMQ_INSTANCE_ID")
//
//	resourceFile := "resource_ctyun_function_trigger_rocketmq.tf"
//
//	resource.Test(t, resource.TestCase{
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		CheckDestroy: func(s *terraform.State) error {
//			_, exists := s.RootModule().Resources[resourceName]
//			if exists {
//				return fmt.Errorf("resource destroy failed")
//			}
//			return nil
//		},
//		Steps: []resource.TestStep{
//			{
//				Config: utils.LoadTestCase(resourceFile,
//					rnd,
//					functionName,
//					triggerName,
//					regionId,
//					rocketmqInstanceId,
//				),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "function_name", functionName),
//					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName),
//					resource.TestCheckResourceAttr(resourceName, "trigger_type", "rocketmq"),
//					resource.TestCheckResourceAttr(resourceName, "region_id", regionId),
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttrSet(resourceName, "status"),
//					resource.TestCheckResourceAttrWith(resourceName, "event_data", func(value string) error {
//						if !strings.Contains(value, rocketmqInstanceId) {
//							return fmt.Errorf("event_data should contain instanceId")
//						}
//						return nil
//					}),
//				),
//			},
//			{
//				ResourceName:      resourceName,
//				ImportState:       true,
//				ImportStateVerify: true,
//				ImportStateIdFunc: func(s *terraform.State) (string, error) {
//					ds := s.RootModule().Resources[resourceName].Primary
//					region_id := ds.Attributes["region_id"]
//					name := ds.Attributes["trigger_name"]
//					function_name := ds.Attributes["function_name"]
//					return fmt.Sprintf("%s,%s,%s", name, function_name, region_id), nil
//				},
//				ImportStateVerifyIgnore: []string{"event_data"},
//			},
//			{
//				Config: utils.LoadTestCase(resourceFile,
//					rnd,
//					functionName,
//					triggerName+"-update",
//					regionId,
//					rocketmqInstanceId,
//				),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName+"-update"),
//					resource.TestCheckResourceAttrWith(resourceName, "event_data", func(value string) error {
//						if !strings.Contains(value, rocketmqInstanceId) {
//							return fmt.Errorf("updated event_data should contain instanceId")
//						}
//						return nil
//					}),
//				),
//			},
//			{
//				Config:  utils.LoadTestCase(resourceFile, rnd, functionName, triggerName, regionId, rocketmqInstanceId),
//				Destroy: true,
//			},
//		},
//	})
//}
//
//func TestAccCtyunFunctionTrigger_Rabbitmq(t *testing.T) {
//	t.Parallel()
//	rnd := utils.GenerateRandomString()
//	resourceName := "ctyun_function_trigger." + rnd
//
//	triggerName := "trigger-rabbitmq-" + rnd
//	functionName :=dependence.functionName
//
//	regionId := os.Getenv("CTYUN_REGION_ID")
//	rabbitmqInstanceId := os.Getenv("CTYUN_RABBITMQ_INSTANCE_ID")
//
//	resourceFile := "resource_ctyun_function_trigger_rabbitmq.tf"
//
//	resource.Test(t, resource.TestCase{
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		CheckDestroy: func(s *terraform.State) error {
//			_, exists := s.RootModule().Resources[resourceName]
//			if exists {
//				return fmt.Errorf("resource destroy failed")
//			}
//			return nil
//		},
//		Steps: []resource.TestStep{
//			{
//				Config: utils.LoadTestCase(resourceFile,
//					rnd,
//					functionName,
//					triggerName,
//					regionId,
//					rabbitmqInstanceId,
//				),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "function_name", functionName),
//					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName),
//					resource.TestCheckResourceAttr(resourceName, "trigger_type", "rabbitmq"),
//					resource.TestCheckResourceAttr(resourceName, "region_id", regionId),
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttrSet(resourceName, "status"),
//					resource.TestCheckResourceAttrWith(resourceName, "event_data", func(value string) error {
//						if !strings.Contains(value, rabbitmqInstanceId) {
//							return fmt.Errorf("event_data should contain instanceId")
//						}
//						return nil
//					}),
//				),
//			},
//			{
//				ResourceName:      resourceName,
//				ImportState:       true,
//				ImportStateVerify: true,
//				ImportStateIdFunc: func(s *terraform.State) (string, error) {
//					ds := s.RootModule().Resources[resourceName].Primary
//					region_id := ds.Attributes["region_id"]
//					name := ds.Attributes["trigger_name"]
//					function_name := ds.Attributes["function_name"]
//					return fmt.Sprintf("%s,%s,%s", name, function_name, region_id), nil
//				},
//				ImportStateVerifyIgnore: []string{"event_data"},
//			},
//			{
//				Config: utils.LoadTestCase(resourceFile,
//					rnd,
//					functionName,
//					triggerName+"-update",
//					regionId,
//					rabbitmqInstanceId,
//				),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName+"-update"),
//					resource.TestCheckResourceAttrWith(resourceName, "event_data", func(value string) error {
//						if !strings.Contains(value, rabbitmqInstanceId) {
//							return fmt.Errorf("updated event_data should contain instanceId")
//						}
//						return nil
//					}),
//				),
//			},
//			{
//				Config:  utils.LoadTestCase(resourceFile, rnd, functionName, triggerName, regionId, rabbitmqInstanceId),
//				Destroy: true,
//			},
//		},
//	})
//}
//
//func TestAccCtyunFunctionTrigger_Mqtt(t *testing.T) {
//	t.Parallel()
//	rnd := utils.GenerateRandomString()
//	resourceName := "ctyun_function_trigger." + rnd
//
//	triggerName := "trigger-mqtt-" + rnd
//	functionName :=dependence.functionName
//
//	regionId := os.Getenv("CTYUN_REGION_ID")
//	mqttInstanceId := os.Getenv("CTYUN_MQTT_INSTANCE_ID")
//
//	resourceFile := "resource_ctyun_function_trigger_mqtt.tf"
//
//	resource.Test(t, resource.TestCase{
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		CheckDestroy: func(s *terraform.State) error {
//			_, exists := s.RootModule().Resources[resourceName]
//			if exists {
//				return fmt.Errorf("resource destroy failed")
//			}
//			return nil
//		},
//		Steps: []resource.TestStep{
//			{
//				Config: utils.LoadTestCase(resourceFile,
//					rnd,
//					functionName,
//					triggerName,
//					regionId,
//					mqttInstanceId,
//				),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "function_name", functionName),
//					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName),
//					resource.TestCheckResourceAttr(resourceName, "trigger_type", "mqtt"),
//					resource.TestCheckResourceAttr(resourceName, "region_id", regionId),
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttrSet(resourceName, "status"),
//					resource.TestCheckResourceAttrWith(resourceName, "event_data", func(value string) error {
//						if !strings.Contains(value, mqttInstanceId) {
//							return fmt.Errorf("event_data should contain instanceId")
//						}
//						return nil
//					}),
//				),
//			},
//			{
//				ResourceName:      resourceName,
//				ImportState:       true,
//				ImportStateVerify: true,
//				ImportStateIdFunc: func(s *terraform.State) (string, error) {
//					ds := s.RootModule().Resources[resourceName].Primary
//					region_id := ds.Attributes["region_id"]
//					name := ds.Attributes["trigger_name"]
//					function_name := ds.Attributes["function_name"]
//					return fmt.Sprintf("%s,%s,%s", name, function_name, region_id), nil
//				},
//				ImportStateVerifyIgnore: []string{"event_data"},
//			},
//			{
//				Config: utils.LoadTestCase(resourceFile,
//					rnd,
//					functionName,
//					triggerName+"-update",
//					regionId,
//					mqttInstanceId,
//				),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName+"-update"),
//					resource.TestCheckResourceAttrWith(resourceName, "event_data", func(value string) error {
//						if !strings.Contains(value, mqttInstanceId) {
//							return fmt.Errorf("updated event_data should contain instanceId")
//						}
//						return nil
//					}),
//				),
//			},
//			{
//				Config:  utils.LoadTestCase(resourceFile, rnd, functionName, triggerName, regionId, mqttInstanceId),
//				Destroy: true,
//			},
//		},
//	})
//}
//
//func TestAccCtyunFunctionTrigger_Als(t *testing.T) {
//	t.Parallel()
//	rnd := utils.GenerateRandomString()
//	resourceName := "ctyun_function_trigger." + rnd
//
//	triggerName := "trigger-als-" + rnd
//	functionName :=dependence.functionName
//
//	regionId := os.Getenv("CTYUN_REGION_ID")
//	logProjectCode := os.Getenv("CTYUN_LOG_PROJECT_CODE")
//	logProjectName := os.Getenv("CTYUN_LOG_PROJECT_NAME")
//	logUnitCode := os.Getenv("CTYUN_LOG_UNIT_CODE")
//	logUnitName := os.Getenv("CTYUN_LOG_UNIT_NAME")
//
//	resourceFile := "resource_ctyun_function_trigger_als.tf"
//
//	resource.Test(t, resource.TestCase{
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		CheckDestroy: func(s *terraform.State) error {
//			_, exists := s.RootModule().Resources[resourceName]
//			if exists {
//				return fmt.Errorf("resource destroy failed")
//			}
//			return nil
//		},
//		Steps: []resource.TestStep{
//			{
//				Config: utils.LoadTestCase(resourceFile,
//					rnd,
//					functionName,
//					triggerName,
//					regionId,
//					logProjectCode,
//					logProjectName,
//					logUnitCode,
//					logUnitName,
//				),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "function_name", functionName),
//					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName),
//					resource.TestCheckResourceAttr(resourceName, "trigger_type", "als"),
//					resource.TestCheckResourceAttr(resourceName, "region_id", regionId),
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttrSet(resourceName, "status"),
//					resource.TestCheckResourceAttrWith(resourceName, "event_data", func(value string) error {
//						if !strings.Contains(value, logProjectCode) {
//							return fmt.Errorf("event_data should contain logProjectCode")
//						}
//						return nil
//					}),
//				),
//			},
//			{
//				ResourceName:      resourceName,
//				ImportState:       true,
//				ImportStateVerify: true,
//				ImportStateIdFunc: func(s *terraform.State) (string, error) {
//					ds := s.RootModule().Resources[resourceName].Primary
//					region_id := ds.Attributes["region_id"]
//					name := ds.Attributes["trigger_name"]
//					function_name := ds.Attributes["function_name"]
//					return fmt.Sprintf("%s,%s,%s", name, function_name, region_id), nil
//				},
//				ImportStateVerifyIgnore: []string{"event_data"},
//			},
//			{
//				Config: utils.LoadTestCase(resourceFile,
//					rnd,
//					functionName,
//					triggerName+"-update",
//					regionId,
//					logProjectCode,
//					logProjectName,
//					logUnitCode,
//					logUnitName,
//				),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName+"-update"),
//					resource.TestCheckResourceAttrWith(resourceName, "event_data", func(value string) error {
//						if !strings.Contains(value, logProjectCode) {
//							return fmt.Errorf("updated event_data should contain logProjectCode")
//						}
//						return nil
//					}),
//				),
//			},
//			{
//				Config:  utils.LoadTestCase(resourceFile, rnd, functionName, triggerName, regionId, logProjectCode, logProjectName, logUnitCode, logUnitName),
//				Destroy: true,
//			},
//		},
//	})
//}
//
//func TestAccCtyunFunctionTrigger_Zos(t *testing.T) {
//	t.Parallel()
//	rnd := utils.GenerateRandomString()
//	resourceName := "ctyun_function_trigger." + rnd
//
//	triggerName := "trigger-zos-" + rnd
//	functionName :=dependence.functionName
//
//	regionId := os.Getenv("CTYUN_REGION_ID")
//	zosBucket := os.Getenv("CTYUN_ZOS_BUCKET")
//
//	resourceFile := "resource_ctyun_function_trigger_zos.tf"
//
//	resource.Test(t, resource.TestCase{
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		CheckDestroy: func(s *terraform.State) error {
//			_, exists := s.RootModule().Resources[resourceName]
//			if exists {
//				return fmt.Errorf("resource destroy failed")
//			}
//			return nil
//		},
//		Steps: []resource.TestStep{
//			{
//				Config: utils.LoadTestCase(resourceFile,
//					rnd,
//					functionName,
//					triggerName,
//					regionId,
//					zosBucket,
//				),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "function_name", functionName),
//					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName),
//					resource.TestCheckResourceAttr(resourceName, "trigger_type", "zos"),
//					resource.TestCheckResourceAttr(resourceName, "region_id", regionId),
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttrSet(resourceName, "status"),
//					resource.TestCheckResourceAttrWith(resourceName, "event_data", func(value string) error {
//						if !strings.Contains(value, zosBucket) {
//							return fmt.Errorf("event_data should contain bucket")
//						}
//						return nil
//					}),
//				),
//			},
//			{
//				ResourceName:      resourceName,
//				ImportState:       true,
//				ImportStateVerify: true,
//				ImportStateIdFunc: func(s *terraform.State) (string, error) {
//					ds := s.RootModule().Resources[resourceName].Primary
//					region_id := ds.Attributes["region_id"]
//					name := ds.Attributes["trigger_name"]
//					function_name := ds.Attributes["function_name"]
//					return fmt.Sprintf("%s,%s,%s", name, function_name, region_id), nil
//				},
//				ImportStateVerifyIgnore: []string{"event_data"},
//			},
//			{
//				Config: utils.LoadTestCase(resourceFile,
//					rnd,
//					functionName,
//					triggerName+"-update",
//					regionId,
//					zosBucket,
//				),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName+"-update"),
//					resource.TestCheckResourceAttrWith(resourceName, "event_data", func(value string) error {
//						if !strings.Contains(value, zosBucket) {
//							return fmt.Errorf("updated event_data should contain bucket")
//						}
//						return nil
//					}),
//				),
//			},
//			{
//				Config:  utils.LoadTestCase(resourceFile, rnd, functionName, triggerName, regionId, zosBucket),
//				Destroy: true,
//			},
//		},
//	})
//}
//
//func TestAccCtyunFunctionTrigger_Apigateway(t *testing.T) {
//	t.Parallel()
//	rnd := utils.GenerateRandomString()
//	resourceName := "ctyun_function_trigger." + rnd
//
//	triggerName := "trigger-apigateway-" + rnd
//	functionName :=dependence.functionName
//
//	regionId := os.Getenv("CTYUN_REGION_ID")
//	gatewayInstanceId := os.Getenv("CTYUN_API_GATEWAY_INSTANCE_ID")
//	vpceId := os.Getenv("CTYUN_VPCE_ID")
//
//	resourceFile := "resource_ctyun_function_trigger_apigateway.tf"
//
//	resource.Test(t, resource.TestCase{
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		CheckDestroy: func(s *terraform.State) error {
//			_, exists := s.RootModule().Resources[resourceName]
//			if exists {
//				return fmt.Errorf("resource destroy failed")
//			}
//			return nil
//		},
//		Steps: []resource.TestStep{
//			{
//				Config: utils.LoadTestCase(resourceFile,
//					rnd,
//					functionName,
//					triggerName,
//					regionId,
//					gatewayInstanceId,
//					vpceId,
//				),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "function_name", functionName),
//					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName),
//					resource.TestCheckResourceAttr(resourceName, "trigger_type", "apigateway"),
//					resource.TestCheckResourceAttr(resourceName, "region_id", regionId),
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttrSet(resourceName, "status"),
//					resource.TestCheckResourceAttrWith(resourceName, "event_data", func(value string) error {
//						if !strings.Contains(value, gatewayInstanceId) {
//							return fmt.Errorf("event_data should contain gatewayInstanceId")
//						}
//						return nil
//					}),
//				),
//			},
//			{
//				ResourceName:      resourceName,
//				ImportState:       true,
//				ImportStateVerify: true,
//				ImportStateIdFunc: func(s *terraform.State) (string, error) {
//					ds := s.RootModule().Resources[resourceName].Primary
//					region_id := ds.Attributes["region_id"]
//					name := ds.Attributes["trigger_name"]
//					function_name := ds.Attributes["function_name"]
//					return fmt.Sprintf("%s,%s,%s", name, function_name, region_id), nil
//				},
//				ImportStateVerifyIgnore: []string{"event_data"},
//			},
//			{
//				Config: utils.LoadTestCase(resourceFile,
//					rnd,
//					functionName,
//					triggerName+"-update",
//					regionId,
//					gatewayInstanceId,
//					vpceId,
//				),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "trigger_name", triggerName+"-update"),
//					resource.TestCheckResourceAttrWith(resourceName, "event_data", func(value string) error {
//						if !strings.Contains(value, gatewayInstanceId) {
//							return fmt.Errorf("updated event_data should contain gatewayInstanceId")
//						}
//						return nil
//					}),
//				),
//			},
//			{
//				Config:  utils.LoadTestCase(resourceFile, rnd, functionName, triggerName, regionId, gatewayInstanceId, vpceId),
//				Destroy: true,
//			},
//		},
//	})
//}
