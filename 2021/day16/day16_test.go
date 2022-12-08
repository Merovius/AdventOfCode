package main

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tcs := []struct {
		input string
		want  Packet
	}{
		{"D2FE28", Literal{6, 2021}},
		{
			"38006F45291200",
			Operator{
				Ver: 1,
				ID:  6,
				Sub: []Packet{
					Literal{6, 10},
					Literal{2, 20},
				},
			},
		},
		{
			"EE00D40C823060",
			Operator{
				Ver: 7,
				ID:  3,
				Sub: []Packet{
					Literal{2, 1},
					Literal{4, 2},
					Literal{1, 3},
				},
			},
		},
		{
			"8A004A801A8002F478",
			Operator{
				Ver: 4,
				ID:  2,
				Sub: []Packet{
					Operator{
						Ver: 1,
						ID:  2,
						Sub: []Packet{
							Operator{
								Ver: 5,
								ID:  2,
								Sub: []Packet{Literal{Ver: 6, Val: 15}},
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range tcs {
		p, err := ParsePacket(tc.input)
		if err != nil {
			t.Fatalf("ParsePacket(%q) = _, %v, want <nil>", tc.input, err)
		}
		if !reflect.DeepEqual(p, tc.want) {
			t.Fatalf("ReadPacket(%q) = %v, want %v", tc.input, p, tc.want)
		}
	}
}

func TestVersionSum(t *testing.T) {
	tcs := []struct {
		input string
		want  int
	}{
		{"8A004A801A8002F478", 16},
		{"620080001611562C8802118E34", 12},
		{"C0015000016115A2E0802F182340", 23},
		{"A0016C880162017C3686B18A3D4780", 31},
	}
	for _, tc := range tcs {
		p, err := ParsePacket(tc.input)
		if err != nil {
			t.Fatalf("ParsePacket(%q) = _, %v, want <nil>", tc.input, err)
		}
		if got := p.VersionSum(); got != tc.want {
			t.Fatalf("VersionSum(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestEval(t *testing.T) {
	tcs := []struct {
		input string
		want  int
	}{
		{"C200B40A82", 3},
		{"04005AC33890", 54},
		{"880086C3E88112", 7},
		{"CE00C43D881120", 9},
		{"D8005AC2A8F0", 1},
		{"F600BC2D8F", 0},
		{"9C005AC2F8F0", 0},
		{"9C0141080250320F1802104A08", 1},
	}
	for _, tc := range tcs {
		p, err := ParsePacket(tc.input)
		if err != nil {
			t.Fatalf("ParsePackage(%q) = _, %v, want <nil>", tc.input, err)
		}
		got := p.Eval()
		if got != tc.want {
			t.Logf("Packet expression: %v", p.Expr())
			t.Fatalf("Eval(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}
