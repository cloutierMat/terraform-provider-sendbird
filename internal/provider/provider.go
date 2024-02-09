package provider

import (
	"context"
	"os"

	"github.com/cloutierMat/terraform-provider-sendbird/internal/client"
	"github.com/cloutierMat/terraform-provider-sendbird/internal/service/organization"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &SendbirdProvider{}

type SendbirdProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type SendbirdProviderModel struct {
	ApiKey types.String `tfsdk:"api_key"`
	Host   types.String `tfsdk:"host"`
}

func (p *SendbirdProvider) Metadata(_ context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sendbird"
	resp.Version = p.version
}

func (p *SendbirdProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API Key for the Sendbird Organisation",
				Optional:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Sendbird Organisation host api",
				Optional:            true,
			},
		},
	}
}

func (p *SendbirdProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config SendbirdProviderModel

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

	client := client.New(host, apiKey)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *SendbirdProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// NewSendbirdUserResource,
	}
}

func (p *SendbirdProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		organization.NewApplicationDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SendbirdProvider{
			version: version,
		}
	}
}
