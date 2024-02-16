package application_test

import (
	"fmt"
	"testing"

	"github.com/cloutierMat/terraform-provider-sendbird/internal/acctest"
	"github.com/cloutierMat/terraform-provider-sendbird/internal/service/names"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccApplicationResource_basic(t *testing.T) {
	randName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	region := names.RegionDefaultKey
	resourceName := "sendbird_application.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationResource_basic(randName, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "api_token"),
					resource.TestCheckResourceAttr(resourceName, "region_key", region),
					resource.TestCheckResourceAttr(resourceName, "region_name", "Central, Canada"),
					resource.TestMatchResourceAttr(resourceName, "created_at", acctest.MatchDateTime),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"created_at"},
			},
		},
	})
}

func TestAccApplicationResource_updateName(t *testing.T) {
	randName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	randNameUpdated := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	region := names.RegionDefaultKey
	resourceName := "sendbird_application.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationResource_basic(randName, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randName),
				),
			},
			{
				Config: testAccApplicationResource_basic(randNameUpdated, region),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randNameUpdated),
				),
			},
		},
	})
}

func TestAccApplicationResource_updateRegion(t *testing.T) {
	randName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	region := names.RegionCanadaCentral1Key
	regionUpdated := names.RegionNorthVirginia2Key
	resourceName := "sendbird_application.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationResource_basic(randName, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "region_key", region),
					resource.TestCheckResourceAttr(resourceName, "region_name", "Central, Canada"),
				),
			},
			{
				Config: testAccApplicationResource_basic(randName, regionUpdated),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "region_key", regionUpdated),
					resource.TestCheckResourceAttr(resourceName, "region_name", "North Virginia 2, USA"),
				),
			},
		},
	})
}

func testAccApplicationResource_basic(name string, region string) string {
	return acctest.ProviderConfig + fmt.Sprintf(`
resource "sendbird_application" "test" {
	name = %[1]q
	region_key = %[2]q
}
`, name, region)
}
