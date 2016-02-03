package transport

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
)

func cookieMock(req *http.Request) (*http.Response, error) {
	c, err := req.Cookie("engine-api-test-cookie")
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, fmt.Errorf("unable to find engine-api-test-cookie in the request")
	}
	resp := &http.Response{
		Header:     http.Header{},
		StatusCode: http.StatusNoContent,
	}
	resp.Header.Add("Set-Cookie", `engine-api-test-cookie-response=true; max-age=30`)

	return resp, nil
}

func TestCookieMiddleware(t *testing.T) {
	cookie := &http.Cookie{
		Name:  "engine-api-test-cookie",
		Value: "true",
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	u, err := url.Parse("http://localhost/info")
	if err != nil {
		t.Fatal(err)
	}
	jar.SetCookies(u, []*http.Cookie{cookie})

	m := NewCookieJarMiddleware(jar)
	c := NewMockClient(nil, cookieMock)

	chain := m(c)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = chain.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	cookies := jar.Cookies(req.URL)
	if len(cookies) != 2 {
		t.Fatalf("expected 2 cookies, got %v\n", len(cookies))
	}
	if cookies[1].Name != "engine-api-test-cookie-response" {
		t.Fatalf("response cookie not found, got %q\n", cookies)
	}
}
