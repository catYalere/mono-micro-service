package utils

import (
	"testing"
)

type Test struct {
	Field1 string `bson:"fieldbson1,omitempty" json:"fieldjson1,omitempty"`
	Field2 string `bson:"fieldbson2,omitempty" json:"fieldjson2,omitempty"`
}

func TestGetReflection(t *testing.T) {
	type args struct {
		entity *Test
		tag    string
		value  string
		set    string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Get error no reflection found",
			args: args{
				entity: &Test{
					Field1: "value1",
					Field2: "value2",
				},
				tag:   "something",
				value: "something",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "GetValue - Check if reflection is filled correctly when tag is bson",
			args: args{
				entity: &Test{
					Field1: "value1",
					Field2: "value2",
				},
				tag:   "bson",
				value: "fieldbson1",
			},
			want:    "value1",
			wantErr: false,
		},
		{
			name: "GetValue - Check if reflection is filled correctly when tag is json",
			args: args{
				entity: &Test{
					Field1: "value1",
					Field2: "value2",
				},
				tag:   "json",
				value: "fieldjson2",
			},
			want:    "value2",
			wantErr: false,
		},
		{
			name: "SetValue - Check if reflection is filled correctly when tag is json",
			args: args{
				entity: &Test{
					Field1: "value1",
					Field2: "value2",
				},
				tag:   "json",
				value: "fieldjson2",
				set:   "value3",
			},
			want:    "value3",
			wantErr: false,
		},
		{
			name: "SetValue - Check if reflection is filled correctly when tag is bson",
			args: args{
				entity: &Test{
					Field1: "value1",
					Field2: "value2",
				},
				tag:   "bson",
				value: "fieldbson1",
				set:   "value3",
			},
			want:    "value3",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetReflection[Test](tt.args.entity, tt.args.tag, tt.args.value)
			if err != nil {
				if !tt.wantErr {
					t.Errorf(`expected "%s", got error "%v"`, tt.want, err)
				}
				return
			}
			if tt.args.set != "" {
				got.SetValue(tt.args.set)
			}
			if got.GetStringValue() != tt.want {
				t.Errorf(`expected "%s", got "%v"`, tt.want, got.GetStringValue())
			}
		})
	}
}
