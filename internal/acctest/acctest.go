package acctest

import (
	"os"
	"regexp"
	"sync"
	"testing"

	"github.com/cloutierMat/terraform-provider-sendbird/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	ProviderConfig = `
	provider "sendbird" {}
	`
	ResourcePrefix = "tf-acc-test"
)

var (
	MatchDateTime            = regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]+)?([Zz]|([\+-])([01]\d|2[0-3]):?([0-5]\d)?)?`)
	ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"sendbird": providerserver.NewProtocol6WithError(&provider.SendbirdProvider{}),
	}
)

// testAccProviderConfigure ensures Provider is only configured once
//
// The PreCheck(t) function is invoked for every test and this prevents
// extraneous reconfiguration to the same values each time.
var testAccProviderConfigure sync.Once

// PreCheck verifies and sets required provider testing configuration
//
// This PreCheck function should be present in every acceptance test. It allows
// test configurations to omit a provider configuration with region and ensures
// testing functions that attempt to call Sendbird API are previously configured.
//
// These verifications and configuration are preferred at this level to prevent
// provider developer from experiencing less clear errors for every test.
func PreCheck(t *testing.T) {
	testAccProviderConfigure.Do(func() {
		if os.Getenv("SENDBIRD_API_KEY") == "" {
			t.Fatal("SENDBIRD_API_KEY must be set for acceptance tests")
		}
	})
}
