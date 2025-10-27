package lavinmq

import (
	"context"
	"math/big"
	"strings"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &bindingResource{}
	_ resource.ResourceWithConfigure   = &bindingResource{}
	_ resource.ResourceWithImportState = &bindingResource{}
)

func NewBindingResource() resource.Resource {
	return &bindingResource{}
}

type bindingResource struct {
	services *clientlibrary.Services
}

type bindingResourceModel struct {
	Vhost           types.String  `tfsdk:"vhost"`
	Source          types.String  `tfsdk:"source"`
	Destination     types.String  `tfsdk:"destination"`
	DestinationType types.String  `tfsdk:"destination_type"`
	RoutingKey      types.String  `tfsdk:"routing_key"`
	Arguments       types.Dynamic `tfsdk:"arguments"`
	PropertiesKey   types.String  `tfsdk:"properties_key"`
}

func (r *bindingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_binding"
}

func (r *bindingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a binding between an exchange and a queue or exchange.",
		Attributes: map[string]schema.Attribute{
			"vhost": schema.StringAttribute{
				Description: "The vhost the binding is located in.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source": schema.StringAttribute{
				Description: "The source exchange name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"destination": schema.StringAttribute{
				Description: "The destination queue or exchange name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"destination_type": schema.StringAttribute{
				Description: "The destination type: 'queue' or 'exchange'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("queue"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"routing_key": schema.StringAttribute{
				Description: "The routing key for the binding.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"arguments": schema.DynamicAttribute{
				Description: "Optional binding arguments.",
				Optional:    true,
				PlanModifiers: []planmodifier.Dynamic{
					dynamicRequiresReplace(),
				},
			},
			"properties_key": schema.StringAttribute{
				Description: "Unique properties key for this binding.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *bindingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.services = req.ProviderData.(*clientlibrary.Services)
}

func (r *bindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importIDParts := strings.Split(req.ID, "@")

	if len(importIDParts) != 5 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected format: vhost@source@destination@destination_type@properties_key",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vhost"), importIDParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("source"), importIDParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("destination"), importIDParts[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("destination_type"), importIDParts[3])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("properties_key"), importIDParts[4])...)
}

func (r *bindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan bindingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var request clientlibrary.BindingRequest
	request.RoutingKey = plan.RoutingKey.ValueString()

	argumentsMap := make(map[string]any)
	if !plan.Arguments.IsNull() && !plan.Arguments.IsUnknown() {
		switch v := plan.Arguments.UnderlyingValue().(type) {
		case types.Object:
			for key, value := range v.Attributes() {
				switch val := value.(type) {
				case types.String:
					argumentsMap[key] = val.ValueString()
				case types.Bool:
					argumentsMap[key] = val.ValueBool()
				case types.Number:
					if bigFloat := val.ValueBigFloat(); bigFloat != nil {
						if intVal, accuracy := bigFloat.Int64(); accuracy == big.Exact {
							argumentsMap[key] = intVal
						} else if floatVal, accuracy := bigFloat.Float64(); accuracy == big.Exact {
							argumentsMap[key] = floatVal
						}
					}
				}
			}
		}
	}
	if len(argumentsMap) > 0 {
		request.Arguments = argumentsMap
	}

	err := r.services.Bindings.Create(
		ctx,
		plan.Vhost.ValueString(),
		plan.Source.ValueString(),
		plan.Destination.ValueString(),
		plan.DestinationType.ValueString(),
		request,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error creating binding", err.Error())
		return
	}

	bindings, err := r.services.Bindings.List(ctx, plan.Vhost.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading binding", err.Error())
		return
	}

	var found *clientlibrary.BindingResponse
	for _, binding := range bindings {
		if binding.Source == plan.Source.ValueString() &&
			binding.Destination == plan.Destination.ValueString() &&
			binding.DestinationType == plan.DestinationType.ValueString() &&
			binding.RoutingKey == plan.RoutingKey.ValueString() {
			found = &binding
			break
		}
	}

	if found == nil {
		resp.Diagnostics.AddError("Error reading binding", "Binding not found after creation")
		return
	}

	plan.PropertiesKey = types.StringValue(found.PropertiesKey)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *bindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state bindingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	binding, err := r.services.Bindings.Get(
		ctx,
		state.Vhost.ValueString(),
		state.Source.ValueString(),
		state.Destination.ValueString(),
		state.DestinationType.ValueString(),
		state.PropertiesKey.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error reading binding", err.Error())
		return
	}
	if binding == nil {
		tflog.Info(ctx, "Binding not found on server, removing from state", map[string]any{
			"vhost":       state.Vhost.ValueString(),
			"source":      state.Source.ValueString(),
			"destination": state.Destination.ValueString(),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	state.RoutingKey = types.StringValue(binding.RoutingKey)

	if len(binding.Arguments) > 0 {
		attributes := make(map[string]attr.Value)
		for key, value := range binding.Arguments {
			switch v := value.(type) {
			case int64:
				attributes[key] = types.NumberValue(new(big.Float).SetInt64(v))
			case float64:
				attributes[key] = types.NumberValue(new(big.Float).SetFloat64(v))
			case bool:
				attributes[key] = types.BoolValue(v)
			case string:
				attributes[key] = types.StringValue(v)
			}
		}

		attributeTypes := make(map[string]attr.Type)
		for key := range attributes {
			attributeTypes[key] = attributes[key].Type(ctx)
		}

		argumentsObject, diags := types.ObjectValue(attributeTypes, attributes)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		state.Arguments = types.DynamicValue(argumentsObject)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *bindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"Bindings cannot be updated. They must be replaced.",
	)
}

func (r *bindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state bindingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.services.Bindings.Delete(
		ctx,
		state.Vhost.ValueString(),
		state.Source.ValueString(),
		state.Destination.ValueString(),
		state.DestinationType.ValueString(),
		state.PropertiesKey.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting binding", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func dynamicRequiresReplace() planmodifier.Dynamic {
	return &dynamicRequiresReplaceModifier{}
}

type dynamicRequiresReplaceModifier struct{}

func (m *dynamicRequiresReplaceModifier) Description(ctx context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

func (m *dynamicRequiresReplaceModifier) MarkdownDescription(ctx context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

func (m *dynamicRequiresReplaceModifier) PlanModifyDynamic(ctx context.Context, req planmodifier.DynamicRequest, resp *planmodifier.DynamicResponse) {
	if req.State.Raw.IsNull() {
		return
	}

	if req.Plan.Raw.IsNull() {
		resp.RequiresReplace = true
		return
	}

	if !req.ConfigValue.Equal(req.StateValue) {
		resp.RequiresReplace = true
	}
}
