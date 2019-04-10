package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want string
	}{
		{
			name: "default",
			date: time.Date(2019, 4, 8, 10, 0, 0, 0, time.UTC),
			want: "Apr 08 2019, 10:00",
		},
		{
			name: "empty",
			date: time.Time{},
			want: "?",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if fd := formatDate(test.date); fd != test.want {
				t.Errorf("want %q; got %q", test.want, fd)
			}
		})
	}
}

func TestFormatLink(t *testing.T) {
	tests := []struct {
		name    string
		shorty  string
		request *http.Request
		want    string
	}{
		{
			name:    "https",
			shorty:  "AbH",
			request: httptest.NewRequest("GET", "https://test.com", nil),
			want:    "https://test.com/AbH",
		},
		{
			name:    "http",
			shorty:  "zx0",
			request: httptest.NewRequest("GET", "http://test.com", nil),
			want:    "http://test.com/zx0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if link := formatLink(test.shorty, test.request); link != test.want {
				t.Errorf("want %q; got %q", test.want, link)
			}
		})
	}
}
