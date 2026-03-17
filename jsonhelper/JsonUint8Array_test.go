package jsonhelper

import (
	"encoding/json"
	"testing"
)

func TestJsonUint8Array_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		arr  JsonUint8Array
		want string
	}{
		{"nil", nil, "null"},
		{"empty", JsonUint8Array{}, "[]"},
		{"single", JsonUint8Array{42}, "[42]"},
		{"multiple", JsonUint8Array{1, 2, 255}, "[1,2,255]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.arr)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("Marshal = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestJsonUint8Array_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    JsonUint8Array
		wantErr bool
	}{
		{"null", "null", nil, false},
		{"empty", "[]", JsonUint8Array{}, false},
		{"single", "[42]", JsonUint8Array{42}, false},
		{"multiple", "[1,2,255]", JsonUint8Array{1, 2, 255}, false},
		{"invalid", `"not an array"`, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got JsonUint8Array
			err := json.Unmarshal([]byte(tt.input), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Fatalf("len = %d, want %d", len(got), len(tt.want))
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Errorf("index %d: got %d, want %d", i, got[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestJsonUint8Array_RoundTrip(t *testing.T) {
	original := JsonUint8Array{0, 128, 255}
	b, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	var decoded JsonUint8Array
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(original) != len(decoded) {
		t.Fatalf("round trip: len %d != %d", len(decoded), len(original))
	}
	for i := range original {
		if original[i] != decoded[i] {
			t.Errorf("round trip: index %d: got %d, want %d", i, decoded[i], original[i])
		}
	}
}

func TestJsonUint8Array_InStruct(t *testing.T) {
	type Msg struct {
		Data JsonUint8Array `json:"data"`
	}

	original := Msg{Data: JsonUint8Array{10, 20, 30}}
	b, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}

	var decoded Msg
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatal(err)
	}
	for i := range original.Data {
		if original.Data[i] != decoded.Data[i] {
			t.Errorf("struct round trip: index %d: got %d, want %d", i, decoded.Data[i], original.Data[i])
		}
	}
}
