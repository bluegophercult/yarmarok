package service_test

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/kaznasho/yarmarok/service"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

func TestExcelManagerWriteXLSX(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	type notAStruct int

	type customString string

	var tests = map[string]struct {
		collections []interface{}
		sheetIdx    int
		wantRowsNum int
		wantErr     error
	}{
		"single struct": {
			collections: []interface{}{person{Name: "Alice", Age: 25}},
			sheetIdx:    0,
			wantRowsNum: 2,
			wantErr:     nil,
		},
		"slice of structs": {
			collections: []interface{}{
				&service.Raffle{"raffle_id", "organizer_id", "Raffle", "Wow wow wow", time.Now()},
				service.Prize{"prize_id", "Super prize", 42, "cat in the bag", time.Now()},
				[]service.Participant{
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
			collections: []interface{}{[]person{}},
			wantErr:     errors.New("empty collection"),
		},
		"invalid collections": {
			collections: []interface{}{[]interface{}{nil, nil, []struct{}{}, &[]***struct{}{}, &[]struct{}{}}},
			wantErr:     errors.New("invalid type, expected ..."),
		},
		"invalid ": {
			collections: []interface{}{1234},
			wantErr:     errors.New("invalid type, expected ..."),
		},
		"non-struct nor slice value": {
			collections: []interface{}{notAStruct(1)},
			wantErr:     fmt.Errorf("invalid type, expected struct or slice, got: %s", reflect.TypeOf(notAStruct(1)).Kind()),
		},
		"non-exported field struct": {
			collections: []interface{}{&struct{ name string }{"bob"}},
			wantErr:     errors.New("invalid collection: ..."),
		},
		"custom type field struct": {
			collections: []interface{}{&struct{ Age notAStruct }{Age: notAStruct(25)}},
			wantErr:     errors.New("invalid collection: ..."),
		},
		"pointer type field struct": {
			collections: []interface{}{&struct{ Age *int }{Age: new(int)}},
			wantErr:     errors.New("invalid collection: ..."),
		},
		"interface type field struct": {
			collections: []interface{}{&struct{ Age interface{} }{Age: 25}},
			wantErr:     errors.New("invalid collection: ..."),
		},
		"slice type field struct": {
			collections: []interface{}{&struct{ Ages []int }{Ages: []int{25, 26, 27}}},
			wantErr:     errors.New("invalid collection: ..."),
		},
		"map type field struct": {
			collections: []interface{}{&struct{ Ages map[string]int }{Ages: map[string]int{"Alice": 25}}},
			wantErr:     errors.New("invalid collection: ..."),
		},
		"function type field struct": {
			collections: []interface{}{&struct{ Age func() int }{Age: func() int { return 25 }}},
			wantErr:     errors.New("invalid collection: ..."),
		},
		"channel type field struct": {
			collections: []interface{}{&struct{ Age chan int }{Age: make(chan int)}},
			wantErr:     errors.New("invalid collection: ..."),
		},
		"array type field struct": {
			collections: []interface{}{&struct{ Ages [3]int }{Ages: [3]int{25, 26, 27}}},
			wantErr:     errors.New("invalid collection: ..."),
		},
		"struct type field struct": {
			collections: []interface{}{&struct{ Age struct{ Value int } }{Age: struct{ Value int }{Value: 25}}},
			wantErr:     errors.New("invalid collection: ..."),
		},
		"empty stringer interface field struct": {
			collections: []interface{}{&struct{ Stringer fmt.Stringer }{*new(interface{ String() string })}},
			wantErr:     errors.New("invalid collection: ..."),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			em := service.NewXLSX()

			buf := new(bytes.Buffer)
			err := em.WriteXLSX(buf, tc.collections...)

			if tc.wantErr != nil {
				require.Error(t, err)
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
