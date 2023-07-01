package service

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// ExcelManager represents a manager for Excel files, using the excelize package.
type ExcelManager struct {
	File *excelize.File
}

// NewExcel is a function that creates a new ExcelManager instance with a new Excel file
func NewExcel() *ExcelManager {
	return &ExcelManager{File: excelize.NewFile()}
}

// Sheet is a type that represents an Excel sheet
// that corresponds Go flat struct.
type Sheet struct {
	name string
	rows []Row
}

// Row is a type that represents an Excel row.
type Row = []interface{}

// WriteExcel writes the provided collections into Excel sheets and writes the Excel file.
func (em *ExcelManager) WriteExcel(w io.Writer, collections ...interface{}) error {
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
		return fmt.Errorf("writing excel: %w", err)
	}

	if err := em.File.Close(); err != nil {
		return fmt.Errorf("closing excel: %w", err)
	}

	return nil
}

// addSheet adds a new sheet to the Excel file.
func (em *ExcelManager) addSheet(sheet *Sheet) error {
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
	val := reflect.Indirect(reflect.ValueOf(collection))
	if val.Kind() == reflect.Struct {
		val = reflect.Append(reflect.New(reflect.SliceOf(val.Type())).Elem(), val)
	}

	if val.Kind() != reflect.Slice || val.Type().Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("invalid type, expected struct or slice, got: %s", val.Kind())
	}

	if val.Len() == 0 {
		return nil, errors.New("empty collection")
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
