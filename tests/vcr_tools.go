package tests

import (
	"bytes"
	"io"
	"net/http"
	"reflect"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
)

var HeadersToIgnore = []string{
	"Authorization",
	"X-Amz-Date",
	"X-Amz-Content-Sha256",
	"User-Agent",
	"Amz-Sdk-Request",
	"Amz-Sdk-Invocation-Id",
	"X-Amzn-Requestid",
	"Date",
	"X-Amz-Id-2",
	"X-Amz-Request-Id",
	"Content-Length",
}

func deepEqualContents(x, y any) bool {
	if reflect.ValueOf(x).IsNil() {
		if reflect.ValueOf(y).IsNil() {
			return true
		} else {
			return reflect.ValueOf(y).Len() == 0
		}
	} else {
		if reflect.ValueOf(y).IsNil() {
			return reflect.ValueOf(x).Len() == 0
		} else {
			return reflect.DeepEqual(x, y)
		}
	}
}

// TODO fix rest of matchers

func bodyMatches(r *http.Request, i cassette.Request) bool {
	if r.Body != nil {
		var buffer bytes.Buffer
		if _, err := buffer.ReadFrom(r.Body); err != nil {
			return false
		}

		r.Body = io.NopCloser(bytes.NewBuffer(buffer.Bytes()))
		if buffer.String() != i.Body {
			return false
		}
	} else {
		if len(i.Body) != 0 {
			return false
		}
	}

	return true
}

// modified version from default matcher
var CustomMatcher = func(r *http.Request, i cassette.Request) bool {
	if r.Method != i.Method {
		return false
	}

	if r.URL.String() != i.URL {
		return false
	}

	if r.Proto != i.Proto {
		return false
	}

	if r.ProtoMajor != i.ProtoMajor {
		return false
	}

	if r.ProtoMinor != i.ProtoMinor {
		return false
	}

	requestHeader := r.Header.Clone()
	cassetteRequestHeaders := i.Headers.Clone()

	for _, header := range HeadersToIgnore {
		delete(requestHeader, header)
		delete(cassetteRequestHeaders, header)
	}

	if !deepEqualContents(requestHeader, cassetteRequestHeaders) {
		return false
	}

	if !bodyMatches(r, i) {
		return false
	}

	if !deepEqualContents(r.TransferEncoding, i.TransferEncoding) {
		return false
	}

	if r.Host != i.Host {
		return false
	}

	// Only ParseForm for non-GET requests since that would use query params
	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
		err := r.ParseForm()
		if err != nil {
			return false
		}
	}
	if !deepEqualContents(r.Form, i.Form) {
		return false
	}

	if !deepEqualContents(r.Trailer, i.Trailer) {
		return false
	}

	if r.RemoteAddr != i.RemoteAddr {
		return false
	}

	if r.RequestURI != i.RequestURI {
		return false
	}

	return true
}
