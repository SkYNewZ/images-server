package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

func Test_imageIDFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{context.WithValue(context.TODO(), imageIDContextKey, "foo")},
			want: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := imageIDFromContext(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("imageIDFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_respond(t *testing.T) {
	type fields struct {
		router *mux.Router
		Image  ImageService
	}
	type args struct {
		w          http.ResponseWriter
		in1        *http.Request
		data       interface{}
		statusCode int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "No body",
			fields: fields{},
			args: args{
				w:          nil,
				in1:        nil,
				data:       nil,
				statusCode: 300,
			},
			wantErr: false,
		},
		{
			name:   "With body",
			fields: fields{},
			args: args{
				w:          nil,
				in1:        nil,
				data:       map[string]interface{}{"foo": "bar"},
				statusCode: 400,
			},
			wantErr: false,
		},

		{
			name:   "Invalid json",
			fields: fields{},
			args: args{
				w:          nil,
				in1:        nil,
				data:       map[string]interface{}{"foo": make(chan struct{})},
				statusCode: 500,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			s := &server{
				router: tt.fields.router,
				Image:  tt.fields.Image,
			}

			if err := s.respond(rr, tt.args.in1, tt.args.data, tt.args.statusCode); (err != nil) != tt.wantErr {
				t.Errorf("respond() error = %v, wantErr %v", err, tt.wantErr)
			}

			resp := rr.Result()
			defer resp.Body.Close()

			// Test response body
			if tt.args.data != nil {
				got, _ := io.ReadAll(resp.Body)
				got = bytes.TrimSpace(got)
				want, _ := json.Marshal(tt.args.data)
				want = bytes.TrimSpace(want)

				if !bytes.Equal(got, want) {
					t.Errorf("respond() body = %s, want %s", got, want)
				}
			}

			// Test response status code
			if got := resp.StatusCode; got != tt.args.statusCode {
				t.Errorf("respond() statusCode = %d, want %d", got, tt.args.statusCode)
			}
		})
	}
}

func Test_server_handleGetImageName(t *testing.T) {
	wantValue := "bar"
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if value := req.Context().Value(imageIDContextKey); value != wantValue {
			t.Errorf("handleGetImageName() context value = %s, want %s", value, wantValue)
		}
	})

	s := new(server)
	rr := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/foo", nil)
	req = mux.SetURLVars(req, map[string]string{
		"image": wantValue,
	})

	// Run the test
	s.handleGetImageName(handler)(rr, req)

}
