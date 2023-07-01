package service

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name string
	Age  int
}

func TestExcelManager_WriteExcel(t *testing.T) {
	testCases := []struct {
		name        string
		collections []interface{}
		wantErr     error
	}{
		{
			name: "Valid single collection",
			collections: []interface{}{
				[]TestStruct{
					{"Alice", 30},
					{"Alice", 30},
					{"Alice", 30}},
			},
			wantErr: nil,
		},
		{
			name: "Valid multiple collections",
			collections: []interface{}{
				TestStruct{"Alice", 30},
				[]TestStruct{{"Alice", 30}},
				[]TestStruct{{"Bob", 40}},
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			em := NewExcel()

			err := em.WriteExcel(buf, tc.collections...)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			assert.Nil(t, err)
			assert.NotZero(t, buf.Len())
		})
	}
}
