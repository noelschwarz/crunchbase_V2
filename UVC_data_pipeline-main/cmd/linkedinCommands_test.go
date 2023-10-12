package main

import (
	"reflect"
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	type args struct {
		date string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "Correct date format",
			args: args{
				date: "2020-Jan-08",
			},
			want:    time.Date(2020, time.January, 8, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name: "Invalid month in date input",
			args: args{
				date: "2029-Flop-22",
			},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name: "Invalid year in date input",
			args: args{
				date: "20a-Dec-22",
			},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name: "Invalid year and month in date input",
			args: args{
				date: "eduardo-5g-82",
			},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name: "Invalid day in date input",
			args: args{
				date: "1998-Jun-82",
			},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name: "Correct date format with day with two decimals",
			args: args{
				date: "2001-Mar-12",
			},
			want:    time.Date(2001, time.March, 12, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name: "Invalid month in date input (month in decimals)",
			args: args{
				date: "2001-09-11",
			},
			want:    time.Time{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDate(tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
