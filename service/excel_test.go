package service

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

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
		sheetIdx    int
		wantRowsNum int
		wantErr     error
	}{
		"single struct": {
			collections: []interface{}{Person{Name: "Alice", Age: 25}},
			sheetIdx:    0,
			wantRowsNum: 2,
			wantErr:     nil,
		},
		"slice of structs": {
			collections: []interface{}{
				&Raffle{"raffle_id", "organizer_id", "Raffle", "Wow wow wow", time.Now()},
				Prize{"prize_id", "Super prize", 42, "cat in the bag", time.Now()},
				[]Participant{
					{"participant_id_1", "Bob George", "323421341", "nope", time.Now()},
					{"participant_id_2", "Mr Kitty", "123455", "mew mew", time.Now()},
					{"participant_id_3", "Mr Cat", "123456", "mew mew", time.Now()},
				},
			},
			sheetIdx:    2,
			wantRowsNum: 4,
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

			if tc.wantErr != nil {
				return
			}

			f, err := excelize.OpenReader(buf)
			require.NoError(t, err)

			rows, err := f.GetRows(f.GetSheetName(tc.sheetIdx))
			require.NoError(t, err)

			require.Len(t, rows, tc.wantRowsNum)

			err = f.Close()
			require.NoError(t, err)
		})
	}
}
