package dmlock

import (
	"testing"
)

func TestMarshalCustomerMap(t *testing.T) {

	type args struct {
		m map[string]int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				m: map[string]int64{
					"cus_1": 1710053169605830079,
				},
			},
			want: `{"cus_1":1710053169605830079}`,
		},
		{
			name: "",
			args: args{
				m: map[string]int64{
					"cus_1": 1710053224390141479,
					"cus_2": 1710053224390141754,
				},
			},
			want: `{"cus_1":1710053224390141479,"cus_2":1710053224390141754}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MarshalCustomerMap(tt.args.m); got != tt.want {
				t.Errorf("MarshalCustomerMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
