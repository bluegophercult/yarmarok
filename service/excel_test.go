package service

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

func TestExcelManager_WriteExcel(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	type NotAStruct int

	tests := map[string]struct {
		collections []interface{}
		wantRowsNum int
		wantErr     error
	}{
		"single struct": {
			collections: []interface{}{Person{Name: "Alice", Age: 25}},
			wantRowsNum: 2,
			wantErr:     nil,
		},
		"slice of structs": {
			collections: []interface{}{[]Person{{"Bob", 30}, {"Charlie", 40}}},
			wantRowsNum: 3,
			wantErr:     nil,
		},
		"empty slice": {
			collections: []interface{}{[]Person{}},
			wantErr:     errors.New("empty collection"),
		},
		"non-struct nor slice value": {
			collections: []interface{}{NotAStruct(1)},
			wantErr:     fmt.Errorf("invalid type, expected struct or slice, got: %s", reflect.TypeOf(NotAStruct(1)).Kind()),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			em := NewExcel()

			buf := new(bytes.Buffer)
			err := em.WriteExcel(buf, tc.collections...)
			require.Equal(t, tc.wantErr, err)

			if err != nil {
				return
			}

			f, err := excelize.OpenReader(buf)
			require.NoError(t, err)

			rows, err := f.GetRows(f.GetSheetName(1))
			require.NoError(t, err)

			require.Len(t, rows, tc.wantRowsNum)
		})
	}
}
