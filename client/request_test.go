package client

import (
	"fmt"
	"net/http"
	"testing"
)

func TestIsUnauthorized(t *testing.T) {
	cases := []struct {
		resp   *serverResponse
		err    error
		unauth bool
	}{
		{&serverResponse{statusCode: http.StatusOK}, nil, false},
		{&serverResponse{statusCode: http.StatusUnauthorized}, nil, true},
		{&serverResponse{statusCode: http.StatusInternalServerError}, nil, false},
		{&serverResponse{statusCode: http.StatusInternalServerError}, fmt.Errorf("unauthorized"), true},
	}

	for _, cs := range cases {
		got := isUnauthorized(cs.resp, cs.err)
		if got != cs.unauth {
			t.Fatalf("expected %v, got %v, with code %v and error %v", cs.unauth, got, cs.resp.statusCode, cs.err)
		}
	}
}
