// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networksecurity_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/envvar"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func TestAccNetworkSecurityFirewallEndpoints_basic(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		CheckDestroy:             testAccCheckNetworkSecurityFirewallEndpointDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecurityFirewallEndpoints_basic(orgId, randomSuffix),
			},
			{
				ResourceName:            "google_network_security_firewall_endpoint.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkSecurityFirewallEndpoints_update(orgId, randomSuffix),
			},
			{
				ResourceName:            "google_network_security_firewall_endpoint.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkSecurityFirewallEndpoints_basic(orgId string, randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_network_security_firewall_endpoint" "foobar" {
    provider = google-beta
    name     = "tf-test-my-firewall-endpoint%s"
    parent   = "organizations/%s"
    location = "us-central1-a"

    labels = {
        foo = "bar"
    }
}
`, randomSuffix, orgId)
}

func testAccNetworkSecurityFirewallEndpoints_update(orgId string, randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_network_security_firewall_endpoint" "foobar" {
    provider = google-beta
    name     = "tf-test-my-firewall-endpoint%s"
    parent   = "organizations/%s"
    location = "us-central1-a"

    labels = {
        foo = "bar-updated"
    }
}
`, randomSuffix, orgId)
}

func testAccCheckNetworkSecurityFirewallEndpointDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_network_security_firewall_endpoint" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{NetworkSecurityBasePath}}{{parent}}/locations/{{location}}/firewallEndpoints/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("NetworkSecurityFirewallEndpoint still exists at %s", url)
			}
		}

		return nil
	}
}
