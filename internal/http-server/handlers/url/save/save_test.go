package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorishak/url-shortener/internal/http-server/handlers/url/save"
	"github.com/gorishak/url-shortener/internal/http-server/handlers/url/save/mocks"
	"github.com/gorishak/url-shortener/internal/lib/logger/handlers/slogdiscard"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://test.com",
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://test.com",
		},
		{
			name:      "Empty url",
			url:       "",
			alias:     "some_alias",
			respError: "field URL is required",
		},
		{
			name:      "Invalid url",
			url:       "invalid url",
			alias:     "some_alias",
			respError: "field URL is not a valid URL",
		},
		{
			name:      "Save url error",
			alias:     "test_alias",
			url:       "https://test.com",
			respError: "failed to save url",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).Return(int64(1), tc.mockError).Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
