package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type ApplicationModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	ApiToken   types.String `tfsdk:"api_token"`
	CreatedAt  types.String `tfsdk:"created_at"`
	RegionKey  types.String `tfsdk:"region_key"`
	RegionName types.String `tfsdk:"region_name"`
}
