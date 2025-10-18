package lavinmq

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &federationUpstreamSetResource{}
	_ resource.ResourceWithConfigure   = &federationUpstreamSetResource{}
	_ resource.ResourceWithImportState = &federationUpstreamSetResource{}
)

func NewFederationUpstreamSetResource() resource.Resource {
	return &federationUpstreamSetResource{}
}

type federationUpstreamSetResource struct {
	services *clientlibrary.Services
}

type federationUpstreamSetResourceModel struct {
	Name      types.String `tfsdk:"name"`
	Vhost     types.String `tfsdk:"vhost"`
	Upstreams types.List   `tfsdk:"upstreams"`
}

func (r *federationUpstreamSetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_federation_upstream_set"
}

func (r *federationUpstreamSetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a federation upstream set for high availability federation configurations.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the federation upstream set.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vhost": schema.StringAttribute{
				Description: "Virtual host where the federation upstream set is defined.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"upstreams": schema.ListAttribute{
				Description: "List of upstream names that belong to this set.",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *federationUpstreamSetResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.services = req.ProviderData.(*clientlibrary.Services)
}

func (r *federationUpstreamSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan federationUpstreamSetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var upstreams []string
	resp.Diagnostics.Append(plan.Upstreams.ElementsAs(ctx, &upstreams, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upstreamItems := make([]clientlibrary.FederationUpstreamSetItem, len(upstreams))
	for i, upstream := range upstreams {
		upstreamItems[i] = clientlibrary.FederationUpstreamSetItem{Upstream: upstream}
	}

	createReq := clientlibrary.ParameterRequest{
		Value: upstreamItems,
	}

	err := r.services.Parameters.CreateOrUpdate(ctx, "federation-upstream-set", plan.Vhost.ValueString(), plan.Name.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating federation upstream set", err.Error())
		return
	}

	parameter, err := r.services.Parameters.Get(ctx, "federation-upstream-set", plan.Vhost.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read federation upstream set data", err.Error())
		return
	}

	if err := updateFederationUpstreamSetStateFromParameter(ctx, &plan, parameter); err != nil {
		resp.Diagnostics.AddError("Failed to update state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "create diag failed")
	}
}

func (r *federationUpstreamSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state federationUpstreamSetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parameter, err := r.services.Parameters.Get(ctx, "federation-upstream-set", state.Vhost.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read federation upstream set data", err.Error())
		return
	}
	if parameter == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if err := updateFederationUpstreamSetStateFromParameter(ctx, &state, parameter); err != nil {
		resp.Diagnostics.AddError("Failed to update state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *federationUpstreamSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan federationUpstreamSetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var upstreams []string
	resp.Diagnostics.Append(plan.Upstreams.ElementsAs(ctx, &upstreams, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upstreamItems := make([]clientlibrary.FederationUpstreamSetItem, len(upstreams))
	for i, upstream := range upstreams {
		upstreamItems[i] = clientlibrary.FederationUpstreamSetItem{Upstream: upstream}
	}

	updateReq := clientlibrary.ParameterRequest{
		Value: upstreamItems,
	}

	err := r.services.Parameters.CreateOrUpdate(ctx, "federation-upstream-set", plan.Vhost.ValueString(), plan.Name.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating federation upstream set", err.Error())
		return
	}

	parameter, err := r.services.Parameters.Get(ctx, "federation-upstream-set", plan.Vhost.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read federation upstream set data", err.Error())
		return
	}

	if err := updateFederationUpstreamSetStateFromParameter(ctx, &plan, parameter); err != nil {
		resp.Diagnostics.AddError("Failed to update state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *federationUpstreamSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan federationUpstreamSetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.services.Parameters.Delete(ctx, "federation-upstream-set", plan.Vhost.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting federation upstream set", err.Error())
		return
	}
}

func (r *federationUpstreamSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "@")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected format: vhost@federation_upstream_set_name",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vhost"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), parts[1])...)
}

func updateFederationUpstreamSetStateFromParameter(ctx context.Context, state *federationUpstreamSetResourceModel, parameter *clientlibrary.ParameterResponse) error {
	state.Name = types.StringValue(parameter.Name)
	state.Vhost = types.StringValue(parameter.Vhost)

	valueBytes, err := json.Marshal(parameter.Value)
	if err != nil {
		return err
	}

	var upstreamItems []clientlibrary.FederationUpstreamSetItem
	if err := json.Unmarshal(valueBytes, &upstreamItems); err != nil {
		return err
	}

	upstreams := make([]string, len(upstreamItems))
	for i, item := range upstreamItems {
		upstreams[i] = item.Upstream
	}

	upstreamsList, diags := types.ListValueFrom(ctx, types.StringType, upstreams)
	if diags.HasError() {
		return &ValidationError{"Failed to convert upstreams to list"}
	}

	state.Upstreams = upstreamsList

	return nil
}
