package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func Test_check(t *testing.T) {
	b := int64(10)
	for i := 0; i < 100; i++ {
		a := int64(i)
		if a%10 == 0 {
			fmt.Println()
		}

		f := float64(a) / float64(b)
		r := a * 2 / b

		fmt.Printf("%3d %10.3f [%5s] %3d -> %3d\n", a, f, strconv.FormatInt(r, 2), r*b/2, HalfEvenRounding(a, b))
	}
}

func Test_halfEvenRounding(t *testing.T) {
	type args struct {
		a int64
		b int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "0.4",
			args: args{a: 4, b: 10},
			want: 0,
		},
		{
			name: "0.5",
			args: args{a: 5, b: 10},
			want: 0,
		},
		{
			name: "0.6",
			args: args{a: 6, b: 10},
			want: 1,
		},
		{
			name: "1.4",
			args: args{a: 14, b: 10},
			want: 1,
		},
		{
			name: "1.5",
			args: args{a: 15, b: 10},
			want: 2,
		},
		{
			name: "1.6",
			args: args{a: 15, b: 10},
			want: 2,
		},
		{
			name: "2.4",
			args: args{a: 24, b: 10},
			want: 2,
		}, {
			name: "2.5",
			args: args{a: 25, b: 10},
			want: 2,
		}, {
			name: "2.6",
			args: args{a: 26, b: 10},
			want: 3,
		},
		{
			name: "3.4",
			args: args{a: 34, b: 10},
			want: 3,
		},
		{
			name: "3.5",
			args: args{a: 35, b: 10},
			want: 4,
		},
		{
			name: "3.6",
			args: args{a: 36, b: 10},
			want: 4,
		},
		{
			name: "4.5",
			args: args{a: 45, b: 10},
			want: 4,
		},
		{
			name: "5.5",
			args: args{a: 55, b: 10},
			want: 6,
		},
		{
			name: "6.5",
			args: args{a: 65, b: 10},
			want: 6,
		},
		{
			name: "7.5",
			args: args{a: 75, b: 10},
			want: 8,
		},
		{
			name: "8.5",
			args: args{a: 85, b: 10},
			want: 8,
		},
		{
			name: "9.5",
			args: args{a: 95, b: 10},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, HalfEvenRounding(tt.args.a, tt.args.b), "halfEvenRounding(%v, %v)", tt.args.a, tt.args.b)
		})
	}
}
