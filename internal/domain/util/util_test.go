package util

import (
	"testing"
)

func TestInt64SlicesAreSame(t *testing.T) {
	type args struct {
		a []int64
		b []int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				a: []int64{1, 2, 3},
				b: []int64{3, 2, 1},
			},
			want: true,
		},
		{
			args: args{
				a: []int64{1, 2, 3, 3, 1},
				b: []int64{3, 2, 1},
			},
			want: true,
		},
		{
			args: args{
				a: []int64{1, 2},
				b: []int64{1, 2},
			},
			want: true,
		},
		{
			args: args{
				a: []int64{1},
				b: []int64{1, 1, 1, 1},
			},
			want: true,
		},
		{
			args: args{
				a: []int64{},
				b: []int64{},
			},
			want: true,
		},
		{
			args: args{
				a: []int64{1},
				b: []int64{},
			},
			want: false,
		},
		{
			args: args{
				a: []int64{1, 2, 3},
				b: []int64{4, 5, 6},
			},
			want: false,
		},
		{
			args: args{
				a: []int64{1, 2, 3},
				b: []int64{4, 1, 2},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64SlicesAreSame(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Int64SlicesAreSame() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64SlicesIntersection(t *testing.T) {
	type args struct {
		sl1 []int64
		sl2 []int64
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		{
			args: args{
				sl1: []int64{1, 2, 3},
				sl2: []int64{3, 2, 1},
			},
			want: []int64{2, 1, 3},
		},
		{
			args: args{
				sl1: []int64{1, 2, 3},
				sl2: []int64{4, 2, 5},
			},
			want: []int64{2},
		},
		{
			args: args{
				sl1: []int64{1, 2, 3},
				sl2: []int64{3, 2, 7},
			},
			want: []int64{3, 2},
		},
		{
			args: args{
				sl1: []int64{1, 2, 3},
				sl2: []int64{4, 5, 6},
			},
			want: []int64{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64SlicesIntersection(tt.args.sl1, tt.args.sl2); !Int64SlicesAreSame(tt.want, got) {
				t.Errorf("Int64SlicesIntersection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64SliceExcludeValues(t *testing.T) {
	type args struct {
		sl []int64
		vs []int64
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		{
			args: args{
				sl: []int64{1, 2, 3},
				vs: []int64{3, 2},
			},
			want: []int64{1},
		},
		{
			args: args{
				sl: []int64{1, 2, 3, 4},
				vs: []int64{5, 6, 7},
			},
			want: []int64{4, 2, 1, 3},
		},
		{
			args: args{
				sl: []int64{1, 2},
				vs: []int64{},
			},
			want: []int64{1, 2},
		},
		{
			args: args{
				sl: []int64{},
				vs: []int64{1, 2},
			},
			want: []int64{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64SliceExcludeValues(tt.args.sl, tt.args.vs); !Int64SlicesAreSame(tt.want, got) {
				t.Errorf("Int64SliceExcludeValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
