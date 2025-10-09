package lavinmq

import (
	"context"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &permissionsDataSource{}
	_ datasource.DataSourceWithConfigure = &permissionsDataSource{}
)

func NewPermissionsDataSource() datasource.DataSource {
	return &permissionsDataSource{}
}

type permissionsDataSource struct {
	services *clientlibrary.Services
}

type permissionsDataSourceModel struct {
	Vhost       types.String                `tfsdk:"vhost"`
	User        types.String                `tfsdk:"user"`
	Permissions []permissionDataSourceModel `tfsdk:"permissions"`
}

type permissionDataSourceModel struct {
	Vhost     types.String `tfsdk:"vhost"`
	User      types.String `tfsdk:"user"`
	Configure types.String `tfsdk:"configure"`
	Read      types.String `tfsdk:"read"`
	Write     types.String `tfsdk:"write"`
}

func (d *permissionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permissions"
}

func (d *permissionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List permissions. Optionally filter by vhost and/or user.",
		Attributes: map[string]schema.Attribute{
			"vhost": schema.StringAttribute{
				Description: "Optional: Filter permissions by vhost.",
				Optional:    true,
			},
			"user": schema.StringAttribute{
				Description: "Optional: Filter permissions by user.",
				Optional:    true,
			},
			"permissions": schema.ListNestedAttribute{
				Description: "List of permissions.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vhost": schema.StringAttribute{
							Description: "Virtual host where the permission is applied.",
							Computed:    true,
						},
						"user": schema.StringAttribute{
							Description: "Name of the user.",
							Computed:    true,
						},
						"configure": schema.StringAttribute{
							Description: "Regular expression pattern for configure permissions.",
							Computed:    true,
						},
						"read": schema.StringAttribute{
							Description: "Regular expression pattern for read permissions.",
							Computed:    true,
						},
						"write": schema.StringAttribute{
							Description: "Regular expression pattern for write permissions.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *permissionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.services = req.ProviderData.(*clientlibrary.Services)
}

func (d *permissionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config permissionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var permissions []clientlibrary.PermissionResponse
	var err error

	// Use appropriate API call based on filters
	if !config.Vhost.IsNull() && !config.User.IsNull() {
		// Both vhost and user specified - get single permission
		permission, err := d.services.Permissions.Get(ctx, config.Vhost.ValueString(), config.User.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Unable to retrieve permission", err.Error())
			return
		}
		if permission != nil {
			permissions = []clientlibrary.PermissionResponse{*permission}
		}
	} else if !config.Vhost.IsNull() {
		permissions, err = d.services.Permissions.ListByVhost(ctx, config.Vhost.ValueString())
	} else if !config.User.IsNull() {
		permissions, err = d.services.Permissions.ListByUser(ctx, config.User.ValueString())
	} else {
		permissions, err = d.services.Permissions.List(ctx)
	}

	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve permissions", err.Error())
		return
	}
	if len(permissions) == 0 {
		tflog.Warn(ctx, "No permissions found")
	}

	config.Permissions = []permissionDataSourceModel{}

	for _, permission := range permissions {
		config.Permissions = append(config.Permissions, permissionDataSourceModel{
			Vhost:     types.StringValue(permission.Vhost),
			User:      types.StringValue(permission.User),
			Configure: types.StringValue(permission.Configure),
			Read:      types.StringValue(permission.Read),
			Write:     types.StringValue(permission.Write),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
