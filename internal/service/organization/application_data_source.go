package organization

import (
	"context"
	"fmt"

	"github.com/cloutierMat/terraform-provider-sendbird/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &ApplicationDataSource{}

func NewApplicationDataSource() datasource.DataSource {
	return &ApplicationDataSource{}
}

type ApplicationDataSource struct {
	client *client.SendbirdClient
}

type RegionModel struct {
	Name types.String `tfsdk:"name"`
	Key  types.String `tfsdk:"key"`
}

type ApplicationDataSourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	ApiToken  types.String `tfsdk:"api_token"`
	CreatedAt types.String `tfsdk:"created_at"`
	Region    *RegionModel `tfsdk:"region"`
}

func (d *ApplicationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (d *ApplicationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Application data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "App Id",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "App Name",
				Computed:            true,
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "Api Token",
				Computed:            true,
				Sensitive:           true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Created At",
				Computed:            true,
			},
			"region": schema.SingleNestedAttribute{
				MarkdownDescription: "Region",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"key": schema.StringAttribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Computed: true,
					},
				},
			},
		},
	}
}

func (d *ApplicationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.SendbirdClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *ApplicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ApplicationDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading application", map[string]interface{}{"app_id": data.Id.ValueString()})
	application, err := d.client.GetApplication(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting application", err.Error())
		return
	}

	state := ApplicationDataSourceModel{
		Id:        types.StringValue(application.Id),
		Name:      types.StringValue(application.Name),
		ApiToken:  types.StringValue(application.ApiToken),
		CreatedAt: types.StringValue(application.CreatedAt),
		Region: &RegionModel{
			Name: types.StringValue(application.Region.Name),
			Key:  types.StringValue(application.Region.Key),
		},
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
