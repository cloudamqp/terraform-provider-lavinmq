package lavinmq

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &lavinmqProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(v string, httpClient *http.Client) provider.Provider {
	return &lavinmqProvider{version: v, httpClient: httpClient}
}

// lavinmqProvider is the provider implementation.
type lavinmqProvider struct {
	version    string
	httpClient *http.Client
}

// lavinmqProviderModel maps provider schema data to a Go type.
type lavinmqProviderModel struct {
	BaseURL  types.String `tfsdk:"baseurl"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Metadata returns the provider type name.
func (p *lavinmqProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "lavinmq"
}

// Schema defines the provider-level schema for configuration data.
func (p *lavinmqProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with LavinMQ API.",
		Attributes: map[string]schema.Attribute{
			"baseurl": schema.StringAttribute{
				Description: "BaseURL API.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username to access the API",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password to access the API",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a lavinmq API client for data sources and resources.
func (p *lavinmqProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring lavinmq client")
	var config lavinmqProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.BaseURL.IsNull() {
		baseURL := os.Getenv("LAVINMQ_API_BASEURL")
		if baseURL == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("baseurl"),
				"Missing LavinMQ API BaseURL",
				"The provider cannot configure the lavinmq API client as there is a missing configuration "+
					"value for the LavinMQ BaseURL.",
			)
		} else {
			config.BaseURL = types.StringValue(baseURL)
		}
	}
	if config.BaseURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("baseurl"),
			"Unknown LavinMQ API BaseURL",
			"The provider cannot configure the lavinmq API client as there is an unknown configuration "+
				"value for the lavinmq BaseURL.",
		)
	}

	if config.Username.IsNull() {
		username := os.Getenv("LAVINMQ_API_USERNAME")
		if username == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("username"),
				"Missing LavinMQ API username",
				"The provider cannot configure the lavinmq API client as there is a missing configuration "+
					"value for the LavinMQ username.",
			)
		} else {
			config.Username = types.StringValue(username)
		}
	}
	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown LavinMQ API username",
			"The provider cannot create the lavinmq API client as there is an unknown configuration "+
				"value for the lavinmq username.",
		)
	}

	if config.Password.IsNull() {
		password := os.Getenv("LAVINMQ_API_PASSWORD")
		if password == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("password"),
				"Missing LavinMQ API password",
				"The provider cannot configure the lavinmq API client as there is a missing configuration "+
					"value for the LavinMQ password.",
			)
		} else {
			config.Password = types.StringValue(password)
		}
	}
	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown LavinMQ API password",
			"The provider cannot create the lavinmq API client as there is an unknown configuration "+
				"value for the lavinmq password.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client := clientlibrary.NewClient(
		config.BaseURL.ValueString(),
		fmt.Sprintf("terraform-provider-lavinmq_%s", p.version),
		config.Username.ValueString(),
		config.Password.ValueString(),
		p.httpClient,
	)
	services := clientlibrary.NewServices(client)
	resp.DataSourceData = services
	resp.ResourceData = services
}

// DataSources defines the data sources implemented in the provider.
func (p *lavinmqProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewVhostDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *lavinmqProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserResource,
		NewVhostResource,
		NewQueueResource,
		NewExchangeResource,
	}
}
