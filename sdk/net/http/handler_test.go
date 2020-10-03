package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traceableai/goagent/sdk/internal/mock"
)

var _ http.Handler = &mockHandler{}

type mockHandler struct {
	baseHandler http.Handler
	spans       []*mock.Span
}

func (h *mockHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	span := &mock.Span{}
	ctx := mock.ContextWithSpan(context.Background(), span)
	h.spans = append(h.spans, span)
	h.baseHandler.ServeHTTP(rw, r.WithContext(ctx))
}

func TestServerRequestIsSuccessfullyTraced(t *testing.T) {
	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("request_id", "abc123xyz")
		rw.WriteHeader(202)
		rw.Write([]byte("test_res"))
		rw.Write([]byte("ponse_body"))
	})

	ih := &mockHandler{baseHandler: EnrichHandler(h, mock.SpanFromContext)}

	r, _ := http.NewRequest("GET", "http://traceable.ai/foo?user_id=1", strings.NewReader("test_request_body"))
	r.Header.Add("api_key", "xyz123abc")
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)
	assert.Equal(t, "test_response_body", w.Body.String())

	spans := ih.spans
	assert.Equal(t, 1, len(spans))

	assert.Equal(t, "http://traceable.ai/foo?user_id=1", spans[0].Attributes["http.url"].(string))
	assert.Equal(t, "xyz123abc", spans[0].Attributes["http.request.header.Api_key"].(string))
	assert.Equal(t, "abc123xyz", spans[0].Attributes["http.response.header.Request_id"].(string))
}

func TestServerRecordsRequestAndResponseBodyAccordingly(t *testing.T) {
	tCases := map[string]struct {
		requestBody                    string
		requestContentType             string
		shouldHaveRecordedRequestBody  bool
		responseBody                   string
		responseContentType            string
		shouldHaveRecordedResponseBody bool
	}{
		"no content type headers and empty body": {
			shouldHaveRecordedRequestBody:  false,
			shouldHaveRecordedResponseBody: false,
		},
		"no content type headers and non empty body": {
			requestBody:                    "{}",
			responseBody:                   "{}",
			shouldHaveRecordedRequestBody:  false,
			shouldHaveRecordedResponseBody: false,
		},
		"content type headers but empty body": {
			requestContentType:             "application/json",
			responseContentType:            "application/x-www-form-urlencoded",
			shouldHaveRecordedRequestBody:  false,
			shouldHaveRecordedResponseBody: false,
		},
		"content type and body": {
			requestBody:                    "test_request_body",
			responseBody:                   "test_response_body",
			requestContentType:             "application/x-www-form-urlencoded",
			responseContentType:            "Application/JSON",
			shouldHaveRecordedRequestBody:  true,
			shouldHaveRecordedResponseBody: true,
		},
	}

	for name, tCase := range tCases {
		t.Run(name, func(t *testing.T) {
			h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.Header().Add("Content-Type", tCase.responseContentType)
				rw.Header().Add("Content-Type", "charset=UTF-8")
				rw.Write([]byte(tCase.responseBody))
			})

			ih := &mockHandler{baseHandler: EnrichHandler(h, mock.SpanFromContext)}

			r, _ := http.NewRequest("GET", "http://traceable.ai/foo", strings.NewReader(tCase.requestBody))
			r.Header.Add("content-type", tCase.requestContentType)

			w := httptest.NewRecorder()

			ih.ServeHTTP(w, r)

			span := ih.spans[0]

			if tCase.shouldHaveRecordedRequestBody {
				assert.Equal(t, tCase.requestBody, span.Attributes["http.request.body"].(string))
			}

			if tCase.shouldHaveRecordedResponseBody {
				assert.Equal(t, tCase.responseBody, span.Attributes["http.response.body"].(string))
			}
		})
	}
}