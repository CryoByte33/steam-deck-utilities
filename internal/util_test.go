package internal

import "testing"

func TestGetHumanVRAMSize(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test 1",
			args: args{
				size: 256,
			},
			want: "256MB",
		},
		{
			name: "Test 2",
			args: args{
				size: 1024,
			},
			want: "1GB",
		},
		{
			name: "Test 3",
			args: args{
				size: 2048,
			},
			want: "2GB",
		},
		{
			name: "Test 4",
			args: args{
				size: 4096,
			},
			want: "4GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHumanVRAMSize(tt.args.size); got != tt.want {
				t.Errorf("getHumanVRAMSize() = %v, want %v", got, tt.want)
			}
		})
	}
}