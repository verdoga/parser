package diagnostics

import "testing"

// TestCodeValues проверяет полный закрытый набор стабильных кодов parser.
func TestCodeValues(t *testing.T) {
	tests := []struct {
		name string
		code Code
		want string
	}{
		{name: "P001", code: P001, want: "P001"},
		{name: "P002", code: P002, want: "P002"},
		{name: "P003", code: P003, want: "P003"},
		{name: "P004", code: P004, want: "P004"},
		{name: "P005", code: P005, want: "P005"},
		{name: "P006", code: P006, want: "P006"},
		{name: "P007", code: P007, want: "P007"},
		{name: "P008", code: P008, want: "P008"},
		{name: "P009", code: P009, want: "P009"},
		{name: "P010", code: P010, want: "P010"},
		{name: "P011", code: P011, want: "P011"},
		{name: "P012", code: P012, want: "P012"},
		{name: "P013", code: P013, want: "P013"},
		{name: "P014", code: P014, want: "P014"},
		{name: "IO001", code: IO001, want: "IO001"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if string(test.code) != test.want {
				t.Fatalf("code = %q, want %q", test.code, test.want)
			}
		})
	}
}
