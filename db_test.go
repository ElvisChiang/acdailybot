package main

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func Test_initDB(t *testing.T) {
	type args struct {
		isReset bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"do reset", args{true}, false},
		{"don't reset", args{false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := initDB(tt.args.isReset); (err != nil) != tt.wantErr {
				t.Errorf("initDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
