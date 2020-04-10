package main

import (
	"log"
	"testing"
)

func TestCommand(t *testing.T) {
	type args struct {
		channel int64
		who     string
		isAdmin bool
		msg     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{123, "", false, ""}, false},

		{"/hl ok", args{123, "username", false, "/hl XXXX"}, false},
		{"/hl empty", args{123, "username", false, "/hl"}, true},
		{"/hl ok hyperlink", args{123, "username", false, "/hl ok [google](https://www.google.com/)"}, false},
		{"/hl at ok", args{123, "username", false, "/hl@ac_dailybot XXXX"}, false},

		{"/hl_remove", args{123, "username", false, "/hl_remove"}, false},
		{"/hl_remove space", args{123, "username", false, "/hl_remove "}, false},
		{"/hl_remove non admin", args{123, "username", false, "/hl_remove username"}, false},
		{"/hl_remove ok", args{123, "username", true, "/hl_remove username"}, false},

		{"/hl_reset", args{123, "username", false, "/hl_reset"}, true},
		{"/hl_reset emtpy", args{123, "username", false, "/hl_reset "}, true},
		{"/hl_reset wrong name", args{123, "username", false, "/hl_reset not_same"}, true},
		{"/hl_reset non admin", args{123, "username", false, "/hl_reset username"}, true},
		{"/hl_reset ok", args{123, "username", true, "/hl_reset username"}, false},

		{"/hl_list space", args{123, "username", false, "/hl_list "}, false},
		{"/hl_list garbage", args{123, "username", false, "/hl_list abcd"}, false},
		{"/hl_list", args{123, "username", false, "/hl_list"}, false},
	}
	db, err := initDB(false)

	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Command(db, tt.args.channel, tt.args.who, tt.args.isAdmin, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Command() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHighlight(t *testing.T) {
	type args struct {
		who string
		msg string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"add1", args{"useraaa", "add1"}, false},
		{"add2", args{"userbbb", "add2"}, false},
		{"add3", args{"userccc", "add3"}, false},
		{"add4", args{"userddd", "markdown [google](https://www.google.com)"}, false},
	}
	channelID := int64(-436800666)
	db, err := initDB(false)

	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Highlight(db, channelID, tt.args.who, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Highlight() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	type args struct {
		who     string
		user    string
		isAdmin bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"non admin useraaa remove bbb", args{"useraaa", "userbbb", false}, true},
		{"admin useraaa remove bbb", args{"useraaa", "userbbb", true}, false},
		{"userbbb remove bbb", args{"userbbb", "", false}, false},
		{"userccc remove ccc", args{"userccc", "", false}, false},
		{"admin userccc remove ddd", args{"userddd", "", true}, false},
	}
	channelID := int64(-436800666)
	db, err := initDB(false)

	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Remove(db, channelID, tt.args.who, tt.args.user, tt.args.isAdmin); (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHighlightList(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"/hl_list", false},
	}
	channelID := int64(-436800666)
	db, err := initDB(false)

	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := HighlightList(db, channelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("HighlightList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

/*
	{"/hl_reset", args{123, "username", false, "/hl_remove"}, true},
	{"/hl_reset emtpy", args{123, "username", false, "/hl_remove "}, true},
	{"/hl_reset wrong name", args{123, "username", false, "/hl_remove not_same"}, true},
	{"/hl_reset non admin", args{123, "username", false, "/hl_remove username"}, true},
	{"/hl_reset ok", args{123, "username", true, "/hl_remove username"}, false},
*/

func TestResetAll(t *testing.T) {
	type args struct {
		isAdmin bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"reset non admin", args{false}, true},
		{"reset by admin", args{true}, false},
	}
	channelID := int64(-436800666)
	db, err := initDB(false)

	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ResetAll(db, channelID, tt.args.isAdmin); (err != nil) != tt.wantErr {
				t.Errorf("ResetAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
