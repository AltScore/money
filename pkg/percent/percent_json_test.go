package percent

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPercent_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Percent
		wantErr bool
	}{
		{"zero", "0", MustParse("0.0"), false},
		{"int", "42", MustParse("42.0"), false},
		{"minimum", "0.001", MustParse("0.001"), false},
		{"large", "123456.789", MustParse("123456.789"), false},
		{"large negative", "-123456.789", MustParse("-123456.789"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p Percent
			if err := p.UnmarshalJSON([]byte(tt.args)); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			} else if !p.Equal(tt.want) {
				t.Errorf("UnmarshalJSON() actual = %v, want %v", p, tt.want)
			}
		})
	}
}

func TestPercent_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		p       Percent
		want    string
		wantErr bool
	}{
		{"zero", MustParse("0.0"), `"0"`, false},
		{"positive", MustParse("122.35"), `"122.35"`, false},
		{"negative", MustParse("-9854.658"), `"-9854.658"`, false},
		{"0 decimals", MustParse("6587.200"), `"6587.2"`, false},
		{"1 decimals", MustParse("6587.200"), `"6587.2"`, false},
		{"2 decimals", MustParse("6587.270"), `"6587.27"`, false},
		{"3 decimals", MustParse("6587.271"), `"6587.271"`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("MarshalJSON() got = %s, want %v", got, tt.want)
			}
		})
	}
}

func Test_shows_that_map_of_percents_keys_serialize_as_strings(t *testing.T) {
	m := map[Percent]string{
		MustParse("10.0"):   "10",
		MustParse("42.032"): "42.032",
		MustParse("67.007"): "67.007",
		MustParse("43.700"): "43.7",
	}

	bytes, _ := json.Marshal(m)

	fmt.Printf("%s\n", bytes)

	m2 := map[string]string{}

	_ = json.Unmarshal(bytes, &m2)

	for k, v := range m2 {
		assert.Equal(t, v, k)
	}

}
