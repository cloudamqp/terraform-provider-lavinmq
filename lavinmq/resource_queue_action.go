package lavinmq

import (
	"context"
	"strings"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &queueActionResource{}
	_ resource.ResourceWithConfigure = &queueActionResource{}
)

// NewQueueActionResource is a helper function to simplify the provider implementation.
func NewQueueActionResource() resource.Resource {
	return &queueActionResource{}
}

// queueActionResource is the resource implementation.
type queueActionResource struct {
	services *clientlibrary.Services
}

// queueActionResourceModel is the
type queueActionResourceModel struct {
	Name   types.String `tfsdk:"name"`
	Vhost  types.String `tfsdk:"vhost"`
	Action types.String `tfsdk:"action"`
}

// Metadata returns the data source type name.
func (r *queueActionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_queue_action"
}

// Schema defines the schema for the resource.
func (r *queueActionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"action": schema.StringAttribute{
				Description: "Action to perform on the queue. Valid values are `purge`.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("purge"),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *queueActionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	services, ok := req.ProviderData.(*clientlibrary.Services)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *clientlibrary.Services type for provider data but got a different type.",
		)
		return
	}

	r.services = services
}

func (r *queueActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan queueActionResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queue, err := r.services.Queues.Get(ctx, plan.Vhost.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Queue",
			"Could not read queue with name "+plan.Name.ValueString()+": "+err.Error(),
		)
		return
	}
	if queue == nil {
		tflog.Warn(ctx, "Queue not found", map[string]any{
			"vhost": plan.Vhost.ValueString(),
			"name":  plan.Name.ValueString(),
		})
	} else {
		switch strings.ToLower(plan.Action.ValueString()) {
		case "purge":
			err = r.services.Queues.Purge(ctx, plan.Vhost.ValueString(), plan.Name.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Purging Queue",
					"Could not purge queue with name "+plan.Name.ValueString()+": "+err.Error(),
				)
				return
			}
		default:
			resp.Diagnostics.AddError(
				"Invalid Action",
				"Action must be one of: purge.",
			)
			return
		}
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *queueActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data queueActionResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *queueActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// This resource does not implement the Update function
}

func (r *queueActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource does not implement the Delete function
}
