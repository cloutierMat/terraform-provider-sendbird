package application_test

import (
	"fmt"
	"testing"

	"github.com/cloutierMat/terraform-provider-sendbird/internal/acctest"
	"github.com/cloutierMat/terraform-provider-sendbird/internal/service/names"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplicationDataSource(t *testing.T) {
	randName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "sendbird_application.test"
	dataSourceName := "data.sendbird_application.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationDataSourceConfig(randName, names.RegionDefaultKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "api_token", dataSourceName, "api_token"),
					resource.TestCheckResourceAttrPair(resourceName, "region_key", dataSourceName, "region_key"),
					resource.TestCheckResourceAttrPair(resourceName, "region_name", dataSourceName, "region_name"),
					resource.TestMatchResourceAttr(dataSourceName, "created_at", acctest.MatchDateTime),
				),
			},
		},
	})
}

func testAccApplicationDataSourceConfig(name string, region string) string {
	return fmt.Sprintf(`
resource "sendbird_application" "test" {
	name = %[1]q
	region_key = %[2]q
}

data "sendbird_application" "test" {id = resource.sendbird_application.test.id}`, name, region)
}
