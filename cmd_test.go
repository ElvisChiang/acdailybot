package main

import "testing"

func TestCommand(t *testing.T) {
	type args struct {
		who string
		msg string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{"", ""}, true},
		{"/hl ok", args{"elvis", "/hl XXXX"}, false},
		{"/hl empty", args{"elvis", "/hl"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Command(tt.args.who, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Command() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
