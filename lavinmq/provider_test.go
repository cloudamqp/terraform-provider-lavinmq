package lavinmq

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

var (
	testAccProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)
	mode                            = recorder.ModeReplayOnly
)

func testAccPreCheck(t *testing.T) {
	if os.Getenv("LAVINMQ_RECORD") != "" {
		if v := os.Getenv("LAVINMQ_API_BASEURL"); v == "" {
			t.Fatal("baseurl must be set for acceptence test.")
		}
		if v := os.Getenv("LAVINMQ_API_USERNAME"); v == "" {
			t.Fatal("username must be set for acceptence test.")
		}
		if v := os.Getenv("LAVINMQ_API_PASSWORD"); v == "" {
			t.Fatal("password must be set for acceptence test.")
		}
	} else {
		os.Setenv("LAVINMQ_API_BASEURL", "http://localhost:15672/")
		os.Setenv("LAVINMQ_API_USERNAME", "not-used")
		os.Setenv("LAVINMQ_API_PASSWORD", "not-used")
	}
}

func TestMain(m *testing.M) {
	if os.Getenv("LAVINMQ_RECORD") != "" {
		mode = recorder.ModeRecordOnly
	}
	resource.TestMain(m)
}

func lavinMQResourceTest(t *testing.T, c resource.TestCase) {
	rec, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName:       fmt.Sprintf("../test/fixtures/vcr/%s", t.Name()),
		Mode:               mode,
		SkipRequestLatency: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Stop()

	sanitizeHook := func(i *cassette.Interaction) error {
		i.Request.Headers["Authorization"] = []string{"REDACTED"}
		// Filter sensitive data API keys, secrects and tokens from request and response bodies
		i.Request.Body = sanitizeSensistiveData(i.Request.Body)
		i.Response.Body = sanitizeSensistiveData(i.Response.Body)
		return nil
	}
	rec.SetMatcher(requestURIMatcher)
	rec.AddHook(sanitizeHook, recorder.AfterCaptureHook)

	shouldSaveHook := func(i *cassette.Interaction) error {
		if t.Failed() {
			i.DiscardOnSave = true
			return nil
		}
		return nil
	}

	rec.AddHook(shouldSaveHook, recorder.BeforeSaveHook)

	rec.AddPassthrough(func(req *http.Request) bool {
		return req.URL.Path == "/login"
	})

	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"lavinmq": providerserver.NewProtocol6WithError(New("vcr-test", rec.GetDefaultClient())),
	}
	c.ProtoV6ProviderFactories = testAccProtoV6ProviderFactories

	resource.Test(t, c)
}

func requestURIMatcher(request *http.Request, interaction cassette.Request) bool {
	interactionURI, err := url.Parse(interaction.URL)
	if err != nil {
		panic(err)
	}

	// https://pkg.go.dev/net/url#URL.RequestURI
	// only match on path?query URL parts
	return request.Method == interaction.Method && request.URL.RequestURI() == interactionURI.RequestURI()
}

func sanitizeSensistiveData(body string) string {
	return body
}
