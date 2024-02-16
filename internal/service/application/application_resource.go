package application

import (
	"context"
	"fmt"

	"github.com/cloutierMat/terraform-provider-sendbird/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &ApplicationResource{}

func NewApplicationResource() resource.Resource {
	return &ApplicationResource{}
}

type ApplicationResource struct {
	client *client.SendbirdClient
}

type ApplicationResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	ApiToken   types.String `tfsdk:"api_token"`
	CreatedAt  types.String `tfsdk:"created_at"`
	RegionKey  types.String `tfsdk:"region_key"`
	RegionName types.String `tfsdk:"region_name"`
}

func (r *ApplicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *ApplicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		MarkdownDescription: "Application Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "App Id",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "App Name. While this value could change without replacing the application, " +
					"Sendbird isn't exposing an API to modify the name. If a name change is required, it is possible to" +
					"change the value on the console and in the configuration. The next `terraform apply` will appropriately" +
					"detect the change and won't replace the application",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region_name": schema.StringAttribute{
				MarkdownDescription: "Region Name",
				Computed:            true,
			},
		},
	}
}

func (r *ApplicationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.SendbirdClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.SendbirdClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *ApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ApplicationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	application, err := r.client.CreateApplication(plan.Name.ValueString(), plan.RegionKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating application",
			"Could not create application, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(application.Id)
	plan.ApiToken = types.StringValue(application.ApiToken)
	plan.CreatedAt = types.StringValue(application.CreatedAt)
	plan.RegionName = types.StringValue(application.Region.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ApplicationResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading application resource", map[string]interface{}{"app_id": data.Id.ValueString()})
	application, err := r.client.GetApplication(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting application resource", err.Error())
		return
	}

	state := ApplicationResourceModel{
		Id:         types.StringValue(application.Id),
		Name:       types.StringValue(application.Name),
		ApiToken:   types.StringValue(application.ApiToken),
		CreatedAt:  types.StringValue(application.CreatedAt),
		RegionKey:  types.StringValue(application.Region.Key),
		RegionName: types.StringValue(application.Region.Name),
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ApplicationResource) Update(_ context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *ApplicationResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ApplicationResourceModel
	diags := req.State.Get(context.Background(), &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteApplication(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting application", err.Error())
		return
	}
}

func (r *ApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
