package service

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// XLSXManager represents a manager for xlsx files, using the excelize package.
type XLSXManager struct {
	File *excelize.File
}

// NewXLSX is a function that creates a new XLSXManager instance with a new xlsx file.
func NewXLSX() *XLSXManager {
	return &XLSXManager{File: excelize.NewFile()}
}

// Sheet is a type that represents an XLSX sheet
// that corresponds Go flat struct.
type Sheet struct {
	name string
	rows []Row
}

// Row is a type that represents an xlsx row.
type Row = []interface{}

// WriteXLSX writes the provided collections into xlsx file.
func (em *XLSXManager) WriteXLSX(w io.Writer, collections ...interface{}) error {
	for i := range collections {
		sheet, err := toSheet(collections[i])
		if err != nil {
			return err
		}

		if err := em.addSheet(sheet); err != nil {
			return err
		}
	}

	if err := em.File.DeleteSheet("Sheet1"); err != nil {
		return err
	}

	if err := em.File.Write(w); err != nil {
		return fmt.Errorf("writing xlsx: %w", err)
	}

	if err := em.File.Close(); err != nil {
		return fmt.Errorf("closing xlsx: %w", err)
	}

	return nil
}

// addSheet adds a new sheet to the xlsx file.
func (em *XLSXManager) addSheet(sheet *Sheet) error {
	if _, err := em.File.NewSheet(sheet.name); err != nil {
		return fmt.Errorf("create sheet %q: %w", sheet.name, err)
	}

	for i, row := range sheet.rows {
		cell := "A" + strconv.Itoa(i+1)
		if err := em.File.SetSheetRow(sheet.name, cell, &row); err != nil {
			return fmt.Errorf("set sheet %q row: %w", sheet.name, err)
		}
	}

	return nil
}

// toSheet converts a collection of structs to a Sheet.
// It verifies that the collection is a slice of structs.
func toSheet(collection interface{}) (*Sheet, error) {
	val, err := extractValue(collection)
	if err != nil {
		return nil, err
	}

	rows := make([]Row, val.Len()+1)
	header := val.Index(0)
	rows[0] = makeRow(header, fieldNameFunc(header))

	for i := 0; i < val.Len(); i++ {
		row := val.Index(i)
		rows[i+1] = makeRow(row, fieldValueFunc(row))
	}

	name := val.Type().Elem().Name()
	sheet := Sheet{name: name, rows: rows}

	return &sheet, nil
}

// extractValue validates a collection if it is a slice of structs
// or struct returning underlying reflected value or error.
func extractValue(collection interface{}) (reflect.Value, error) {
	val := reflect.Indirect(reflect.ValueOf(collection))
	if val.Kind() == reflect.Struct {
		slice := reflect.New(reflect.SliceOf(val.Type())).Elem()
		val = reflect.Append(slice, val)
	}

	if val.Kind() != reflect.Slice || val.Type().Elem().Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("invalid type, expected struct or slice, got: %s", val.Kind())
	}

	if val.Len() == 0 {
		return reflect.Value{}, errors.New("empty collection")
	}

	return val, nil
}

// fieldFunc is used to create a row in a sheet from a struct,
// either by getting the name of each field or by getting the value of each field.
type fieldFunc func(int) interface{}

func fieldNameFunc(val reflect.Value) fieldFunc {
	return func(i int) interface{} {
		return val.Type().Field(i).Name
	}
}

func fieldValueFunc(val reflect.Value) fieldFunc {
	return func(i int) interface{} {
		return val.Field(i).Interface()
	}
}

// makeRow creates a new Row from a struct by applying
// the provided field function to each field in the struct.
func makeRow(val reflect.Value, fn fieldFunc) Row {
	row := make(Row, val.Type().NumField())
	for i := 0; i < val.Type().NumField(); i++ {
		row[i] = fn(i)
	}
	return row
}
