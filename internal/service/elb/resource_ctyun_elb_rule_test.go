package elb_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccCtyunElbRule(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_elb_rule." + rnd
	resourceFile := "resource_ctyun_elb_rule.tf"

	datasourceName := "data.ctyun_elb_rules." + dnd
	datasourceFile := "datasource_ctyun_elb_rules.tf"

	actionType := "forward"

	listenerId := dependence.listenerID
	//conditions := fmt.Sprintf(`{"type": "%s", "condition_server_name": "%s", "condition_url_paths":"%s","condition_match_type":"%s"}`, "server_name", "terraform-test.com", "/test", "PREFIX")
	conditions := fmt.Sprintf(`{"condition_type": "%s", "condition_server_name": "%s"}`, "server_name", "terraform-test.com")
	//updatedConditions := fmt.Sprintf(`{"type": "%s", "condition_server_name": "%s","condition_url_paths":"%s","condition_match_type":"%s"}`, "server_name", "terraform-test-new.com", "test_new", "PREFIX")
	updatedConditions := fmt.Sprintf(`{"condition_type": "%s", "condition_server_name": "%s"}`, "server_name", "terraform-test-new.com")
	pathConditions := fmt.Sprintf(`{"condition_type": "%s","condition_url_paths":"%s","condition_match_type":"%s"}`, "url_path", "test", "PREFIX")
	//updatedPathConditions := fmt.Sprintf(`{"type": "%s","condition_url_paths":"%s","condition_match_type":"%s"}`, "url_path", "test-new", "PREFIX")
	actionTargetGroups := fmt.Sprintf(`{target_group_id="%s"}`, dependence.targetGroupID4)
	//updatedActionTargetGroups := fmt.Sprintf(`{target_group_id="%s"}`, dependence.targetGroupID)

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
			// 1.1 Create验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, listenerId, conditions, actionType, actionTargetGroups),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "listener_id", listenerId),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.condition_type", "server_name"),
					resource.TestCheckResourceAttr(resourceName, "action_type", actionType),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.condition_server_name", "terraform-test.com"),
					//resource.TestCheckResourceAttr(resourceName, "conditions.0.condition_url_paths", "test"),
					//resource.TestCheckResourceAttr(resourceName, "conditions.0.condition_match_type", "PREFIX"),
					resource.TestCheckResourceAttr(resourceName, "action_target_groups.0.target_group_id", dependence.targetGroupID4),
				),
			},
			// 1.2 update
			{
				Config: utils.LoadTestCase(resourceFile, rnd, listenerId, updatedConditions, actionType, actionTargetGroups),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "listener_id", listenerId),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.condition_type", "server_name"),
					resource.TestCheckResourceAttr(resourceName, "action_type", actionType),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.condition_server_name", "terraform-test-new.com"),
					//resource.TestCheckResourceAttr(resourceName, "conditions.0.condition_url_paths", "test_new"),
					//resource.TestCheckResourceAttr(resourceName, "conditions.0.condition_match_type", "PREFIX"),
					resource.TestCheckResourceAttr(resourceName, "action_target_groups.0.target_group_id", dependence.targetGroupID4),
				),
			},
			// 1.3 datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, listenerId, updatedConditions, actionType, actionTargetGroups) +
					utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "elb_rules.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "elb_rules.0.listener_id", listenerId),
					resource.TestCheckResourceAttr(datasourceName, "elb_rules.0.conditions.0.condition_type", "server_name"),
					resource.TestCheckResourceAttr(datasourceName, "elb_rules.0.conditions.0.server_name", "terraform-test-new.com"),
					resource.TestCheckResourceAttr(datasourceName, "elb_rules.0.action_type", actionType),
					resource.TestCheckResourceAttr(datasourceName, "elb_rules.0.action_target_groups.0.target_group_id", dependence.targetGroupID4),
				),
			},
			//1.4 destroy
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, listenerId, updatedConditions, actionType, actionTargetGroups) + utils.LoadTestCase(datasourceFile, dnd, resourceName+".id"),
				Destroy: true,
			},

			// 2 type=url_path 验证
			// 2.1 Create
			{
				Config: utils.LoadTestCase(resourceFile, rnd, listenerId, conditions, actionType, actionTargetGroups),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "listener_id", listenerId),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.condition_type", "server_name"),
					resource.TestCheckResourceAttr(resourceName, "action_type", actionType),
					//resource.TestCheckResourceAttr(resourceName, "conditions.0.condition_url_paths", "test"),
					//resource.TestCheckResourceAttr(resourceName, "conditions.0.condition_match_type", "PREFIX"),
					resource.TestCheckResourceAttr(resourceName, "action_target_groups.0.target_group_id", dependence.targetGroupID4),
				),
			},

			{
				Config: utils.LoadTestCase(resourceFile, rnd, listenerId, pathConditions, actionType, actionTargetGroups),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "listener_id", listenerId),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.condition_type", "url_path"),
					resource.TestCheckResourceAttr(resourceName, "action_type", actionType),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.condition_url_paths", "test"),
					resource.TestCheckResourceAttr(resourceName, "conditions.0.condition_match_type", "PREFIX"),
					resource.TestCheckResourceAttr(resourceName, "action_target_groups.0.target_group_id", dependence.targetGroupID4),
				),
			},
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, listenerId, pathConditions, actionType, actionTargetGroups),
				Destroy: true,
			},
		},
	})
}
