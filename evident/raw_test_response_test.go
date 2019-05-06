package evident

import (
	"fmt"
)

func GetTestOkResponse(account string, arn string, externalId string, teamId string, name string) string {
	return fmt.Sprintf(`
	{
		"data": {
		  "id": "12345",
		  "type": "external_accounts",
		  "attributes": {
			"created_at": "2019-03-15T15:35:51.000Z",
			"name": "%s",
			"updated_at": "2019-03-15T15:35:51.000Z",
			"provider": "amazon",
			"arn": "%s",
			"account": "%s",
			"external_id": "%s",
			"cloudtrail_name": null
		  },
		  "relationships": {
			"organization": {
			  "links": {
				"related": "https://esp.evident.io/api/v2/organizations/3072.json"
			  }
			},
			"sub_organization": {
			  "links": {
				"related": "https://esp.evident.io/api/v2/sub_organizations/9638.json"
			  }
			},
			"team": {
			  "links": {
				"related": "https://esp.evident.io/api/v2/teams/14249.json"
			  }
			},
			"scan_intervals": {
			  "links": {
				"related": "https://esp.evident.io/api/v2/external_accounts/27028/scan_intervals.json"
			  }
			},
			"disabled_signatures": {
			  "links": {
				"related": "https://esp.evident.io/api/v2/external_accounts/27028/disabled_signatures.json"
			  }
			},
			"suppressions": {
			  "links": {
				"related": "https://esp.evident.io/api/v2/suppressions.json?filter"
			  }
			},
			"credentials": {
			  "links": {
				"related": "https://esp.evident.io/api/v2/external_accounts/27028/amazon.json"
			  }
			}
		  }
		}
	  }
	`, name, arn, account, externalId)
}
