package lavinmq

import (
	"context"
	"strings"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &queueResource{}
	_ resource.ResourceWithConfigure   = &queueResource{}
	_ resource.ResourceWithImportState = &queueResource{}
)

// NewQueueResource is a helper function to simplify the provider implementation.
func NewQueueResource() resource.Resource {
	return &queueResource{}
}

// queueResource is the resource implementation.
type queueResource struct {
	services *clientlibrary.Services
}

// queueResourceModel is the
type queueResourceModel struct {
	Name       types.String `tfsdk:"name"`
	Vhost      types.String `tfsdk:"vhost"`
	AutoDelete types.Bool   `tfsdk:"auto_delete"`
	Durable    types.Bool   `tfsdk:"durable"`
	Pause      types.Bool   `tfsdk:"pause"`
	State      types.String `tfsdk:"state"`
}

// Metadata returns the data source type name.
func (r *queueResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_queue"
}

// Schema defines the schema for the resource.
func (r *queueResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a queue.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the managed queue.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vhost": schema.StringAttribute{
				Description: "The vhost the queue is located in.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"auto_delete": schema.BoolAttribute{
				Description: "Whether the queue is automatically deleted when no longer used.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"durable": schema.BoolAttribute{
				Description: "Whether the queue should survive a broker restart.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"pause": schema.BoolAttribute{
				Description: "Queue action, when true, the queue will be paused.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Description: "State of the queue: 'running', 'paused', 'flow', 'closed', or 'deleted'.",
				Computed:    true,
			},
			// "arguments": schema.MapAttribute{
			// 	Description: "Arguments for the queue.",
			// 	Optional:    true,
			// 	PlanModifiers: []planmodifier.Map{
			// 		mapplanmodifier.RequiresReplace(),
			// 	},
			// 	ElementType: types.StringType,
			// },
		},
	}
}

// Configure adds the provider configured services to the resource.
func (r *queueResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.services = req.ProviderData.(*clientlibrary.Services)
}

func (r *queueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importIDParts := strings.Split(req.ID, "@")

	if len(importIDParts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected format: vhost@queue_name",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vhost"), importIDParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), importIDParts[1])...)
}

// Create creates the resource and sets the initial Terraform state.
func (r *queueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan queueResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var request clientlibrary.QueueRequest
	if !plan.AutoDelete.IsUnknown() {
		request.AutoDelete = plan.AutoDelete.ValueBoolPointer()
	}
	if !plan.Durable.IsUnknown() {
		request.Durable = plan.Durable.ValueBoolPointer()
	}

	err := r.services.Queues.CreateOrUpdate(ctx, plan.Vhost.ValueString(), plan.Name.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError("Error creating queue", err.Error())
		return
	}

	queue, err := r.services.Queues.Get(ctx, plan.Vhost.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading queue", err.Error())
		return
	}

	plan.AutoDelete = types.BoolValue(queue.AutoDelete)
	plan.Durable = types.BoolValue(queue.Durable)
	plan.State = types.StringValue(queue.State)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *queueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state queueResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	queue, err := r.services.Queues.Get(ctx, state.Vhost.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading queue", err.Error())
		return
	}
	if queue == nil {
		tflog.Info(ctx, "Queue not found on server, removing from state", map[string]any{
			"vhost": state.Vhost.ValueString(),
			"name":  state.Name.ValueString(),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	state.AutoDelete = types.BoolValue(queue.AutoDelete)
	state.Durable = types.BoolValue(queue.Durable)
	state.State = types.StringValue(queue.State)
	state.Pause = types.BoolValue(queue.State == "paused")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *queueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan queueResourceModel
	var state queueResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Pause.ValueBool() != state.Pause.ValueBool() {
		err := r.services.Queues.Pause(ctx, state.Vhost.ValueString(), state.Name.ValueString(), plan.Pause.ValueBool())
		if err != nil {
			resp.Diagnostics.AddError("Error updating queue pause state", err.Error())
			return
		}

		queue, err := r.services.Queues.Get(ctx, plan.Vhost.ValueString(), plan.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error reading queue", err.Error())
			return
		}

		plan.State = types.StringValue(queue.State)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *queueResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state queueResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.services.Queues.Delete(ctx, state.Vhost.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting queue", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
