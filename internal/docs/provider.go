package docs

const (
	ProviderApiKey                = "API Key for the Sendbird Organisation. This value can also be passed through SENDBIRD_API_KEY environment variable."
	ProviderHost                  = "Host url for the Sendbird Organisation. This value can also be passed through SENDBIRD_HOST environment variable."
	ProviderApplicationRateLimit  = "Rate limit per seconds for the Sendbird Application's queries. Defaults to 5/seconds. This value can also be passed through SENDBIRD_APPLICATION_RATE_LIMIT."
	ProviderOrganizationRateLimit = "Rate limit per minutes for the Sendbird Organization's queries. Defaults to 10/minutes. This value can also be passed through SENDBIRD_ORGANIZATION_RATE_LIMIT. Sendbird considers creating and deleting applications to be at the organization level."
)
