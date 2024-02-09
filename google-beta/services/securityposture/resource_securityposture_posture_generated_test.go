// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package securityposture_test

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

func TestAccSecurityposturePosture_securityposturePostureBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecurityposturePostureDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityposturePosture_securityposturePostureBasicExample(context),
			},
			{
				ResourceName:            "google_securityposture_posture.posture1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "location", "posture_id"},
			},
		},
	})
}

func testAccSecurityposturePosture_securityposturePostureBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_securityposture_posture" "posture1"{
  posture_id  = "posture_1"
  parent      = "organizations/%{org_id}"
  location    = "global"
  state       = "ACTIVE"
  description = "a new posture"
  policy_sets {
    policy_set_id = "org_policy_set"
    description   = "set of org policies"
    policies {
      policy_id = "canned_org_policy"
      constraint {
        org_policy_constraint {
          canned_constraint_id = "storage.uniformBucketLevelAccess"
          policy_rules {
            enforce = true
          }
        }
      }
    }
    policies {
      policy_id = "custom_org_policy"
      constraint {
        org_policy_constraint_custom {
          custom_constraint {
            name         = "organizations/%{org_id}/customConstraints/custom.disableGkeAutoUpgrade"
            display_name = "Disable GKE auto upgrade"
            description  = "Only allow GKE NodePool resource to be created or updated if AutoUpgrade is not enabled where this custom constraint is enforced."
            action_type    = "ALLOW"
            condition      = "resource.management.autoUpgrade == false"
            method_types   = ["CREATE", "UPDATE"]
            resource_types = ["container.googleapis.com/NodePool"]
          }
          policy_rules {
            enforce = true
          }
        }
      }
    }
  }
  policy_sets {
    policy_set_id = "sha_policy_set"
    description   = "set of sha policies"
    policies {
      policy_id = "sha_builtin_module"
      constraint {
        security_health_analytics_module {
          module_name             = "BIGQUERY_TABLE_CMEK_DISABLED"
          module_enablement_state = "ENABLED"
        }
      }
      description = "enable BIGQUERY_TABLE_CMEK_DISABLED"
    }
    policies {
      policy_id = "sha_custom_module"
      constraint {
        security_health_analytics_custom_module {
          display_name = "custom_SHA_policy"
          config {
            predicate {
              expression = "resource.rotationPeriod > duration('2592000s')"
            }
            custom_output {
              properties {
                name = "duration"
                value_expression {
                  expression = "resource.rotationPeriod"
                }
              }
            }
            resource_selector {
              resource_types = ["cloudkms.googleapis.com/CryptoKey"]
            }
            severity       = "LOW"
            description    = "Custom Module"
            recommendation = "Testing custom modules"
          }
          module_enablement_state = "ENABLED"
        }
      }
    }
  }
}
`, context)
}

func testAccCheckSecurityposturePostureDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_securityposture_posture" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{SecuritypostureBasePath}}{{parent}}/locations/{{location}}/postures/{{posture_id}}")
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
				return fmt.Errorf("SecurityposturePosture still exists at %s", url)
			}
		}

		return nil
	}
}
