package client

import (
	"net/http/cookiejar"
	"net/url"
	"testing"

	"github.com/docker/engine-api/client/authn"
	"github.com/docker/engine-api/client/transport"
)

func ExampleClient_AuthenticateWith() {
	client, _ := NewEnvClient()
	basicAuthCallback := func(realm string) (user, password string, err error) {
		return "admin", "password", nil
	}
	basicAuth := authn.NewBasicAuth(basicAuthCallback)

	bearerAuthCallback := func(challenge string) (token string, err error) {
		return "token", nil
	}
	bearerAuth := authn.NewBearerAuth(bearerAuthCallback)

	client.AuthenticateWith(basicAuth, bearerAuth)
}

func ExampleClient_AddMiddlewares() {
	client, _ := NewEnvClient()
	cookieJar, _ := cookiejar.New(nil)
	cookies := transport.NewCookieJarMiddleware(cookieJar)

	client.AddMiddlewares(cookies)
}

func TestGetAPIPath(t *testing.T) {
	cases := []struct {
		v string
		p string
		q url.Values
		e string
	}{
		{"", "/containers/json", nil, "/containers/json"},
		{"", "/containers/json", url.Values{}, "/containers/json"},
		{"", "/containers/json", url.Values{"s": []string{"c"}}, "/containers/json?s=c"},
		{"1.22", "/containers/json", nil, "/v1.22/containers/json"},
		{"1.22", "/containers/json", url.Values{}, "/v1.22/containers/json"},
		{"1.22", "/containers/json", url.Values{"s": []string{"c"}}, "/v1.22/containers/json?s=c"},
		{"v1.22", "/containers/json", nil, "/v1.22/containers/json"},
		{"v1.22", "/containers/json", url.Values{}, "/v1.22/containers/json"},
		{"v1.22", "/containers/json", url.Values{"s": []string{"c"}}, "/v1.22/containers/json?s=c"},
	}

	for _, cs := range cases {
		c, err := NewClient("unix:///var/run/docker.sock", cs.v, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		g := c.getAPIPath(cs.p, cs.q)
		if g != cs.e {
			t.Fatalf("Expected %s, got %s", cs.e, g)
		}
	}
}

func TestParseHost(t *testing.T) {
	cases := []struct {
		host  string
		proto string
		addr  string
		base  string
		err   bool
	}{
		{"", "", "", "", true},
		{"foobar", "", "", "", true},
		{"foo://bar", "foo", "bar", "", false},
		{"tcp://localhost:2476", "tcp", "localhost:2476", "", false},
		{"tcp://localhost:2476/path", "tcp", "localhost:2476", "/path", false},
	}

	for _, cs := range cases {
		p, a, b, e := parseHost(cs.host)
		if cs.err && e == nil {
			t.Fatalf("expected error, got nil")
		}
		if !cs.err && e != nil {
			t.Fatal(e)
		}
		if cs.proto != p {
			t.Fatalf("expected proto %s, got %s", cs.proto, p)
		}
		if cs.addr != a {
			t.Fatalf("expected addr %s, got %s", cs.addr, a)
		}
		if cs.base != b {
			t.Fatalf("expected base %s, got %s", cs.base, b)
		}
	}
}

func TestAddMiddlewares(t *testing.T) {
	m1 := func(n transport.Sender) transport.Sender {
		return n
	}
	m2 := m1
	client := &Client{}
	client.AddMiddlewares(m1, m2)
	if len(client.middlewares) != 2 {
		t.Fatalf("expected 2 middlewares, got %v", len(client.middlewares))
	}
}

func TestAuthenticateWith(t *testing.T) {
	basicAuthCallback := func(realm string) (user, password string, err error) {
		return "admin", "password", nil
	}
	basicAuth := authn.NewBasicAuth(basicAuthCallback)

	client := &Client{}
	client.AuthenticateWith(basicAuth)
	if len(client.middlewares) != 1 {
		t.Fatalf("expected 1 middlewares, got %v", len(client.middlewares))
	}
}
