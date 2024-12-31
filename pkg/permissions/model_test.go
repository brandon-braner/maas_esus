package permissions

import (
	"testing"
)

func TestValidatePermission(t *testing.T) {
	tests := []struct {
		name       string
		permission string
		want       bool
	}{
		{
			name:       "valid permission",
			permission: "generate_llm_meme",
			want:       true,
		},
		{
			name:       "invalid permission",
			permission: "can_use_ai",
			want:       false,
		},
		{
			name:       "empty permission",
			permission: "",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatePermission(tt.permission); got != tt.want {
				t.Errorf("ValidatePermission() = %v, want %v", got, tt.want)
			}
		})
	}
}
