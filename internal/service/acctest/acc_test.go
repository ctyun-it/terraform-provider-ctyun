package acctest

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("Starting Terraform AccTest")
	exitCode := m.Run()
	os.Exit(exitCode)
}

func runTest(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
			},
		},
	})
}

func TestAccUnifiedGateway(t *testing.T) {
	runTest(t)
}

func TestAccVpcInterconnection(t *testing.T) {
	runTest(t)
}

func TestAccHA(t *testing.T) {
	runTest(t)
}

func TestAccRouteRules(t *testing.T) {
	runTest(t)
}

func TestAccBatchEcs(t *testing.T) {
	runTest(t)
}

func TestAccRegion(t *testing.T) {
	runTest(t)
}
