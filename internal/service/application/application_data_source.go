package application

import (
	"context"
	"fmt"

	"github.com/cloutierMat/terraform-provider-sendbird/internal/client"
	"github.com/cloutierMat/terraform-provider-sendbird/internal/service/models"
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
			"region_key": schema.StringAttribute{
				MarkdownDescription: "Region Key",
				Computed:            true,
			},
			"region_name": schema.StringAttribute{
				MarkdownDescription: "Region Name",
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
	var data models.ApplicationModel

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

	state := models.ApplicationModel{
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
