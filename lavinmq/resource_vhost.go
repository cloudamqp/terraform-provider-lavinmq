package lavinmq

import (
	"context"
	"fmt"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vhostResource{}
	_ resource.ResourceWithConfigure   = &vhostResource{}
	_ resource.ResourceWithImportState = &vhostResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewVhostResource() resource.Resource {
	return &vhostResource{}
}

// userResource is the resource implementation.
type vhostResource struct {
	client *clientlibrary.Client
}

type vhostResourceModel struct {
	Name                   types.String          `tfsdk:"name"`
	Dir                    types.String          `tfsdk:"dir"`
	Tracing                types.Bool            `tfsdk:"tracing"`
	Messages               types.Int64           `tfsdk:"messages"`
	MessagesUnacknowledged types.Int64           `tfsdk:"messages_unacknowledged"`
	MessagesReady          types.Int64           `tfsdk:"messages_ready"`
	MessagesStats          basetypes.ObjectValue `tfsdk:"message_stats"`
}

// type vhostMessageStatsResourceModel struct {
// 	Ack              types.Int64 `tfsdk:"ack"`
// 	Confirm          types.Int64 `tfsdk:"confirm"`
// 	Deliver          types.Int64 `tfsdk:"deliver"`
// 	Get              types.Int64 `tfsdk:"get"`
// 	GetNoAck         types.Int64 `tfsdk:"get_no_ack"`
// 	Publish          types.Int64 `tfsdk:"publish"`
// 	Redeliver        types.Int64 `tfsdk:"redeliver"`
// 	ReturnUnroutable types.Int64 `tfsdk:"return_unroutable"`
// }

// Metadata returns the data source type name.
func (r *vhostResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vhost"
}

// Schema defines the schema for the resource.
func (r *vhostResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a vhost.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the managed vhost.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dir": schema.StringAttribute{
				Description: "Directory of the vhost.",
				Computed:    true,
			},
			"tracing": schema.BoolAttribute{
				Description: "Enable or disable tracing for the vhost.",
				Computed:    true,
			},
			"messages": schema.Int64Attribute{
				Description: "Number of messages in the vhost.",
				Computed:    true,
			},
			"messages_unacknowledged": schema.Int64Attribute{
				Description: "Number of unacknowledged messages in the vhost.",
				Computed:    true,
			},
			"messages_ready": schema.Int64Attribute{
				Description: "Number of ready messages in the vhost.",
				Computed:    true,
			},
			"message_stats": schema.SingleNestedAttribute{
				Description: "Statistics about messages in the vhost.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"ack": schema.Int64Attribute{
						Description: "Number of acknowledged messages.",
						Computed:    true,
					},
					"confirm": schema.Int64Attribute{
						Description: "Number of confirmed messages.",
						Computed:    true,
					},
					"deliver": schema.Int64Attribute{
						Description: "Number of delivered messages.",
						Computed:    true,
					},
					"get": schema.Int64Attribute{
						Description: "Number of messages retrieved with 'get'.",
						Computed:    true,
					},
					"get_no_ack": schema.Int64Attribute{
						Description: "Number of messages retrieved with 'get' without acknowledgment.",
						Computed:    true,
					},
					"publish": schema.Int64Attribute{
						Description: "Number of published messages.",
						Computed:    true,
					},
					"redeliver": schema.Int64Attribute{
						Description: "Number of redelivered messages.",
						Computed:    true,
					},
					"return_unroutable": schema.Int64Attribute{
						Description: "Number of messages returned as unroutable.",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vhostResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*clientlibrary.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *vhostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vhostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Vhosts.CreateOrUpdate(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", err.Error())
		return
	}

	// Read out computed values
	vhost, err := r.client.Vhosts.Get(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read user state", err.Error())
		return
	}

	plan.Dir = types.StringValue(vhost.Dir)
	plan.Messages = types.Int64Value(vhost.Messages)
	plan.MessagesUnacknowledged = types.Int64Value(vhost.MessagesUnacknowledged)
	plan.MessagesReady = types.Int64Value(vhost.MessagesReady)
	plan.Tracing = types.BoolValue(vhost.Tracing)

	elementTypes := map[string]attr.Type{
		"ack":               types.Int64Type,
		"confirm":           types.Int64Type,
		"deliver":           types.Int64Type,
		"get":               types.Int64Type,
		"get_no_ack":        types.Int64Type,
		"publish":           types.Int64Type,
		"redeliver":         types.Int64Type,
		"return_unroutable": types.Int64Type,
	}
	elements := map[string]attr.Value{
		"ack":               types.Int64Value(vhost.MessagesStats.Ack),
		"confirm":           types.Int64Value(vhost.MessagesStats.Confirm),
		"deliver":           types.Int64Value(vhost.MessagesStats.Deliver),
		"get":               types.Int64Value(vhost.MessagesStats.Get),
		"get_no_ack":        types.Int64Value(vhost.MessagesStats.GetNoAck),
		"publish":           types.Int64Value(vhost.MessagesStats.Publish),
		"redeliver":         types.Int64Value(vhost.MessagesStats.Redeliver),
		"return_unroutable": types.Int64Value(vhost.MessagesStats.ReturnUnroutable),
	}
	plan.MessagesStats, _ = basetypes.NewObjectValue(elementTypes, elements)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "created diag failed")
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *vhostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vhostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Name.IsUnknown() {
		tflog.Info(ctx, fmt.Sprintf("import resource with name identifier %s", state.Name))
	}

	vhost, err := r.client.Vhosts.Get(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read user data", err.Error())
		return
	}

	tflog.Warn(ctx, fmt.Sprintf("Vhost response: %v", vhost))

	state.Name = types.StringValue(vhost.Name)
	state.Dir = types.StringValue(vhost.Dir)
	state.Messages = types.Int64Value(vhost.Messages)
	state.MessagesUnacknowledged = types.Int64Value(vhost.MessagesUnacknowledged)
	state.MessagesReady = types.Int64Value(vhost.MessagesReady)
	state.Tracing = types.BoolValue(vhost.Tracing)

	elementTypes := map[string]attr.Type{
		"ack":               types.Int64Type,
		"confirm":           types.Int64Type,
		"deliver":           types.Int64Type,
		"get":               types.Int64Type,
		"get_no_ack":        types.Int64Type,
		"publish":           types.Int64Type,
		"redeliver":         types.Int64Type,
		"return_unroutable": types.Int64Type,
	}
	elements := map[string]attr.Value{
		"ack":               types.Int64Value(vhost.MessagesStats.Ack),
		"confirm":           types.Int64Value(vhost.MessagesStats.Confirm),
		"deliver":           types.Int64Value(vhost.MessagesStats.Deliver),
		"get":               types.Int64Value(vhost.MessagesStats.Get),
		"get_no_ack":        types.Int64Value(vhost.MessagesStats.GetNoAck),
		"publish":           types.Int64Value(vhost.MessagesStats.Publish),
		"redeliver":         types.Int64Value(vhost.MessagesStats.Redeliver),
		"return_unroutable": types.Int64Value(vhost.MessagesStats.ReturnUnroutable),
	}
	state.MessagesStats, _ = basetypes.NewObjectValue(elementTypes, elements)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vhostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Note: Even tho the API state update, nothing to update since the name is
	// part of the endpoint. One idea is to include vhost-limits.
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vhostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan vhostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Vhosts.Delete(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting user", err.Error())
		return
	}
}

func (r *vhostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import resource by name argument
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
