package evident

import(
	"fmt"
)

func GetTestOkResponse(account string, arn string, externalId string,teamId string) string {
	return fmt.Sprintf(`
	{
		"data": {
			"id": %s,
			"type": "external_account_amazon_iam",
			"attributes": {
				"account": "123456789012",
				"arn": "%s",
				"external_id": "%s",
				"created_at": "2017-03-09T20:44:08.569Z",
				"updated_at": "2017-03-09T20:44:08.569Z"
			},
			"relationships": {
				"external_account": {
					"links": {
						"related": "https://api.evident.io/api/v2/external_accounts/1016.json_api"
					},
					"data": {
						"id": %s,
						"type": "external_accounts"
					}
				}
			}
		},
		"included": [
			{
				"data": {
					"id": %s,
					"type": "external_accounts",
					"attributes": {
						"created_at": "2017-03-09T20:44:08.569Z",
						"name": "Account 4",
						"updated_at": "2017-03-09T20:44:08.569Z",
						"provider": "amazon"
					},
					"relationships": {
						"team": {
							"links": {
								"related": "https://api.evident.io/api/v2/teams/4.json_api"
							},
							"data": {
								"id": %s,
								"type": "teams"
							}
						},
						"scan_intervals": {
							"links": {
								"related": "https://api.evident.io/api/v2/external_accounts/1006/scan_intervals.json_api"
							},
							"data": [
								{
									"id": 1,
									"type": "scan_intervals"
								}
							]
						},
						"disabled_signatures": {
							"links": {
								"related": "https://api.evident.io/api/v2/external_accounts/1006/disabled_signatures.json_api"
							},
							"data": [
								{
									"id": 1,
									"type": "signatures"
								}
							]
						},
						"suppressions": {
							"links": {
								"related": "https://api.evident.io/api/v2/suppressions.json?filter"
							},
							"data": [
								{
									"id": 1,
									"type": "suppressions"
								}
							]
						},
						"azure_group": {
							"links": {
								"related": "https://api.evident.io/api/v2/azure_groups/1006.json_api"
							},
							"data": {
								"id": 1,
								"type": "azure_groups"
							}
						},
						"credentials": {
							"links": {
								"related": "https://api.evident.io/api/v2/external_accounts/1006/amazons.json_api"
							},
							"data": {
								"id": 1,
								"type": "external_account_amazon_iam"
							}
						}
					}
				}
			}
		]
	}
	`,account , arn, externalId,externalId,externalId ,teamId)
}


