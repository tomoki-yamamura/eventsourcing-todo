package value

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUserID(t *testing.T) {
	tests := map[string]struct {
		input   string
		want    UserID
		wantErr error
	}{
		"valid user ID": {
			input:   "user123",
			want:    UserID("user123"),
			wantErr: nil,
		},
		"user ID with spaces": {
			input:   "  user123  ",
			want:    UserID("user123"),
			wantErr: nil,
		},
		"empty user ID": {
			input:   "",
			want:    UserID(""),
			wantErr: ErrUserIDEmpty,
		},
		"user ID with only spaces": {
			input:   "   ",
			want:    UserID(""),
			wantErr: ErrUserIDEmpty,
		},
		"too long user ID": {
			input:   strings.Repeat("a", 129),
			want:    UserID(""),
			wantErr: ErrUserIDTooLong,
		},
		"exactly 128 characters": {
			input:   strings.Repeat("a", 128),
			want:    UserID(strings.Repeat("a", 128)),
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := NewUserID(tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantErr, err)
				require.Equal(t, UserID(""), result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}

func TestUserID_String(t *testing.T) {
	userID := UserID("test123")
	require.Equal(t, "test123", userID.String())
}
