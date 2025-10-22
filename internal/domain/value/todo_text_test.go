package value

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTodoText(t *testing.T) {
	tests := map[string]struct {
		input   string
		want    TodoText
		wantErr error
	}{
		"valid text": {
			input: "買い物に行く",
			want:  TodoText("買い物に行く"),
		},
		"text with leading and trailing spaces": {
			input: "  勉強する  ",
			want:  TodoText("勉強する"),
		},
		"empty string": {
			input:   "",
			wantErr: ErrTodoTextEmpty,
		},
		"only spaces": {
			input:   "   ",
			wantErr: ErrTodoTextEmpty,
		},
		"exactly 256 characters": {
			input: strings.Repeat("a", 256),
			want:  TodoText(strings.Repeat("a", 256)),
		},
		"over 256 characters": {
			input:   strings.Repeat("a", 257),
			wantErr: ErrTodoTextTooLong,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewTodoText(tt.input)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
