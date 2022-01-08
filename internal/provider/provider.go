package provider

import (
	"context"
	"fmt"
	"github.com/chris922/terraform-provider-openhab/internal/api"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
)

// OpenhabProvider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type OpenhabProvider struct {
	// client can contain the upstream OpenhabProvider SDK or HTTP client used to
	// communicate with the upstream service. Resource and DataSource
	// implementations can then make calls using this client.
	//
	// TODO: If appropriate, implement upstream OpenhabProvider SDK or HTTP client.
	// client vendorsdk.ExampleClient

	// Configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the OpenhabProvider was previously Configured.
	Configured bool

	// Version is set to the OpenhabProvider Version on release, "dev" when the
	// OpenhabProvider is built and ran locally, and "test" when running acceptance
	// testing.
	Version string

	data providerData

	Client *api.Client
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	Endpoint types.String `tfsdk:"endpoint"`
	ApiToken types.String `tfsdk:"api_token"`
}

func (p *OpenhabProvider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Example.Null { /* ... */ }

	// If the upstream OpenhabProvider SDK or HTTP client requires configuration, such
	// as authentication or logging, this is a great opportunity to do so.
	client, err := api.NewClient(data.Endpoint.Value, api.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.SetBasicAuth(data.ApiToken.Value, "")
		return nil
	}))
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create client, got error: %s", err))
		return
	}

	p.Client = client
	p.Configured = true
}

func (p *OpenhabProvider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		//"scaffolding_example": ExampleResourceType{},
		"openhab_item": ItemResourceType{},
	}, nil
}

func (p *OpenhabProvider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		//"scaffolding_example": ExampleDataSourceType{},
	}, nil
}

func (p *OpenhabProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"endpoint": {
				MarkdownDescription: "API endpoint of the target openHAB server, usually the URL with `/rest` suffix, e.g. `https://openhab:8080/rest`",
				Required:            true,
				Type:                types.StringType,
				// TODO: validator
			},
			"api_token": {
				MarkdownDescription: "API token used to authenticate against the openHAB server",
				Required:            true,
				Type:                types.StringType,
				// TODO: validator
			},
		},
	}, nil
}

func NewOpenhabProvider(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &OpenhabProvider{
			Version: version,
		}
	}
}

// ConvertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete OpenhabProvider type. Alternatively,
// this helper can be skipped and the OpenhabProvider type can be directly type
// asserted (e.g. OpenhabProvider: in.(*OpenhabProvider)), however using this can prevent
// potential panics.
func ConvertProviderType(in tfsdk.Provider) (OpenhabProvider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*OpenhabProvider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected OpenhabProvider type (%T) was received. This is always a bug in the OpenhabProvider code and should be reported to the OpenhabProvider developers.", p),
		)
		return OpenhabProvider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty OpenhabProvider instance was received. This is always a bug in the OpenhabProvider code and should be reported to the OpenhabProvider developers.",
		)
		return OpenhabProvider{}, diags
	}

	return *p, diags
}
