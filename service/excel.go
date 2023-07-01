package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type ExcelManager struct {
	File *excelize.File
}

func NewExcel() *ExcelManager {
	return &ExcelManager{File: excelize.NewFile()}
}

type Sheet struct {
	name string
	rows []Row
}

type Row = []interface{}

func (em *ExcelManager) WriteExcel(w io.Writer, collections ...interface{}) error {
	for i := range collections {
		sheet, err := toSheet(collections[i])
		if err != nil {
			return err
		}

		log.Printf("+%v", sheet)
		if err := em.addSheet(sheet); err != nil {
			return err
		}
	}

	em.File.SetActiveSheet(1)

	if err := em.File.Write(w); err != nil {
		return fmt.Errorf("writing excel: %w", err)
	}

	if err := em.File.Close(); err != nil {
		return fmt.Errorf("closing excel: %w", err)
	}

	return nil
}

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

	rows := make([]Row, 0, val.Len())
	header := val.Index(0)

	rows = append(rows, makeRow(header, fieldNameFunc(header)))

	for i := 0; i < val.Len(); i++ {
		row := val.Index(i)
		rows = append(rows, makeRow(row, fieldValueFunc(row)))
	}

	name := val.Type().Elem().Name()
	sheet := Sheet{name: name, rows: rows}

	return &sheet, nil
}

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

func makeRow(val reflect.Value, fn fieldFunc) Row {
	row := make(Row, val.Type().NumField())
	for i := 0; i < val.Type().NumField(); i++ {
		row[i] = fn(i)
	}
	return row
}
