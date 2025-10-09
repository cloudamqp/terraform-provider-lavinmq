package lavinmq

import (
	"context"
	"strings"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"
	"github.com/cloudamqp/terraform-provider-lavinmq/lavinmq/converters"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &usersDataSource{}
	_ datasource.DataSourceWithConfigure = &usersDataSource{}
)

func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

type usersDataSource struct {
	services *clientlibrary.Services
}

type usersDataSourceModel struct {
	Users []userDataSourceModel `tfsdk:"users"`
}

type userDataSourceModel struct {
	Name types.String `tfsdk:"name"`
	Tags types.List   `tfsdk:"tags"`
}

func (d *usersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *usersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"users": schema.ListNestedAttribute{
				Description: "List of users.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the user.",
							Computed:    true,
						},
						"tags": schema.ListAttribute{
							Description: "List of tags associated with the user.",
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.services = req.ProviderData.(*clientlibrary.Services)
}

func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state usersDataSourceModel

	users, err := d.services.Users.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve users", err.Error())
		return
	}

	if users == nil {
		tflog.Warn(ctx, "no users found")
		return
	}

	for _, user := range users {
		var tags types.List
		if user.Tags != "" {
			tagList := strings.Split(user.Tags, ",")
			tags, _ = types.ListValue(types.StringType, converters.StringsToAttrValues(tagList))
		} else {
			tags, _ = types.ListValue(types.StringType, []attr.Value{})
		}

		state.Users = append(state.Users, userDataSourceModel{
			Name: types.StringValue(user.Name),
			Tags: tags,
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
