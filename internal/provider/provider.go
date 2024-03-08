package provider

import (
	"context"
	"os"
	"strconv"

	"github.com/cloutierMat/terraform-provider-sendbird/internal/client"
	"github.com/cloutierMat/terraform-provider-sendbird/internal/docs"
	"github.com/cloutierMat/terraform-provider-sendbird/internal/service/application"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ provider.Provider = &SendbirdProvider{}

type SendbirdProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type SendbirdProviderModel struct {
	ApiKey                types.String `tfsdk:"api_key"`
	Host                  types.String `tfsdk:"host"`
	ApplicationRateLimit  types.Int64  `tfsdk:"application_rate_limit"`
	OrganizationRateLimit types.Int64  `tfsdk:"organization_rate_limit"`
}

func (p *SendbirdProvider) Metadata(_ context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sendbird"
	resp.Version = p.version
}

func (p *SendbirdProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: docs.ProviderApiKey,
				Optional:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: docs.ProviderHost,
				Optional:            true,
			},
			"application_rate_limit": schema.Int64Attribute{
				MarkdownDescription: docs.ProviderApplicationRateLimit,
				Optional:            true,
			},
			"organization_rate_limit": schema.Int64Attribute{
				MarkdownDescription: docs.ProviderOrganizationRateLimit,
				Optional:            true,
			},
		},
	}
}

func (p *SendbirdProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config SendbirdProviderModel
	tflog.Debug(ctx, "Configuring Terraform Sendbird Provider")
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown API Key",
			"The provider cannot create the Sendbird API Client as there is an unknown configuration value for the api key"+
				"Either apply the source of the value first, set the value statically in the configuration or use SENDBIRD_API_KEY environment variable.",
		)
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Host",
			"The provider cannot create the Sendbird API Client as there is an unknown configuration value for the host"+
				"Either apply the source of the value first, set the value statically in the configuration or use SENDBIRD_HOST environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := os.Getenv("SENDBIRD_API_KEY")
	host := os.Getenv("SENDBIRD_HOST")

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	applicationRateLimit, err := strconv.Atoi(os.Getenv("SENDBIRD_APPLICATION_RATE_LIMIT"))
	if err != nil {
		applicationRateLimit = client.ApplicationDefaultRateLimit
	}

	organizationRateLimit, err := strconv.Atoi(os.Getenv("SENDBIRD_ORGANIZATION_RATE_LIMIT"))
	if err != nil {
		organizationRateLimit = client.OrganizationDefaultRateLimit
	}

	if !config.ApplicationRateLimit.IsNull() {
		applicationRateLimit = int(config.ApplicationRateLimit.ValueInt64())
	}

	if !config.OrganizationRateLimit.IsNull() {
		organizationRateLimit = int(config.OrganizationRateLimit.ValueInt64())
	}

	rateLimiter := client.GetRateLimiter(applicationRateLimit, organizationRateLimit)

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing API Key",
			"The Sendbird API Client cannot be created as there is no API Key configured.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client := client.New(host, apiKey, rateLimiter)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *SendbirdProvider) Resources(ctx context.Context) []func() resource.Resource {
	tflog.Debug(ctx, "Registering Terraform Sendbird Provider resources")
	return []func() resource.Resource{
		application.NewApplicationResource,
	}
}

func (p *SendbirdProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	tflog.Debug(ctx, "Registering Terraform Sendbird Provider data sources")
	return []func() datasource.DataSource{
		application.NewApplicationDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SendbirdProvider{
			version: version,
		}
	}
}
