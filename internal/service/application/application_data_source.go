package application

import (
	"context"
	"fmt"

	"github.com/cloutierMat/terraform-provider-sendbird/internal/client"
	"github.com/cloutierMat/terraform-provider-sendbird/internal/docs"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &ApplicationDataSource{}

const ServiceName = "application"

type ApplicationDataSourceModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	ApiToken   types.String `tfsdk:"api_token"`
	CreatedAt  types.String `tfsdk:"created_at"`
	RegionKey  types.String `tfsdk:"region_key"`
	RegionName types.String `tfsdk:"region_name"`
}

func NewApplicationDataSource() datasource.DataSource {
	return &ApplicationDataSource{}
}

type ApplicationDataSource struct {
	client *client.SendbirdClient
}

func (d *ApplicationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ServiceName
}

func (d *ApplicationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Application data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: docs.ApplicationDataSourceId,
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: docs.ApplicationDataSourceName,
				Computed:            true,
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: docs.ApplicationDataSourceApiToken,
				Computed:            true,
				Sensitive:           true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: docs.ApplicationDataSourceCreatedAt,
				Computed:            true,
			},
			"region_key": schema.StringAttribute{
				MarkdownDescription: docs.ApplicationDataSourceRegionKey,
				Computed:            true,
			},
			"region_name": schema.StringAttribute{
				MarkdownDescription: docs.ApplicationDataSourceRegionName,
				Computed:            true,
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
			fmt.Sprintf("Expected *client.SendbirdClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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
		Id:         types.StringValue(application.Id),
		Name:       types.StringValue(application.Name),
		ApiToken:   types.StringValue(application.ApiToken),
		CreatedAt:  types.StringValue(application.CreatedAt),
		RegionKey:  types.StringValue(application.Region.Key),
		RegionName: types.StringValue(application.Region.Name),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
