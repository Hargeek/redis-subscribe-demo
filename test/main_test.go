package test

import (
	"fmt"
	"github.com/imroc/req/v3"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	fmt.Println("TestMain exit with code:", code)
	os.Exit(code)
}

func TestCreateSubscription(t *testing.T) {
	t.Log("TestCreateSubscription start")
	type args struct {
		UserID  string `json:"user_id"`
		Service string `json:"service"`
		Branch  string `json:"branch"`
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestCreateSubscription",
			args: args{
				UserID:  "123",
				Service: "test",
				Branch:  "dev",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("TestCreateSubscription")
			client := req.C().DevMode()
			resp, err := client.R().
				SetBody(&tt.args).
				Post("http://localhost:8080/api/v1/subscription")
			if err != nil {
				t.Errorf("create subscription failed: %v", err)
			}
			t.Log("create subscription response:", resp)
		})
	}
}

func TestCreateNotification(t *testing.T) {
	t.Log("TestCreateNotification start")
	type args struct {
		Service     string `json:"service"`
		Branch      string `json:"branch"`
		TriggerUser string `json:"trigger_user"`
		Status      string `json:"status"`
		Result      string `json:"result"`
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestCreateNotification",
			args: args{
				Service:     "test",
				Branch:      "dev",
				TriggerUser: "123",
				Status:      "success",
				Result:      "ok",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("TestCreateNotification")
			client := req.C().DevMode()
			resp, err := client.R().
				SetBody(&tt.args).
				Post("http://localhost:8080/api/v1/notification")
			if err != nil {
				t.Errorf("create notification failed: %v", err)
			}
			t.Log("create notification response:", resp)
		})
	}
}
