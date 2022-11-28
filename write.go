// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

import (
	"encoding/hex"
	"fmt"
	"github.com/mark-summerfield/gong"
	"io"
	"strconv"
	"time"
)

func (me *Tdb) Write(out io.Writer) error {
	return me.WriteDecimals(out, -1)
}

func (me *Tdb) WriteDecimals(out io.Writer, decimals int) error {
	if !(0 < decimals && decimals < 20) {
		decimals = -1
	}
	var err error
	nl := []byte{'\n'}
	for _, tableName := range me.TableNames {
		table := me.Tables[tableName]
		if err = writeTableMetaData(out, table); err != nil {
			return err
		}
		for _, record := range table.Records {
			sep := ""
			for column, value := range record {
				_, err = out.Write([]byte(sep))
				if err != nil {
					return err
				}
				sep = " "
				kind := table.Fields[column].Kind
				switch kind {
				case BoolField:
					err = writeBool(out, value, kind)
				case BytesField:
					err = writeBytes(out, value, kind)
				case DateField:
					err = writeDateTime(out, value, kind, DateFormat,
						DateStrSentinal)
				case DateTimeField:
					err = writeDateTime(out, value, kind, DateTimeFormat,
						DateTimeStrSentinal)
				case IntField:
					err = writeInt(out, value, kind)
				case RealField:
					err = writeReal(out, value, kind, decimals)
				case StrField:
					err = writeStr(out, value, kind)
				default: // should never happen
					return fmt.Errorf("e%d:invalid kind %q", e142, kind)
				}
				if err != nil {
					return err
				}
			}
			_, err = out.Write(nl)
			if err != nil {
				return err
			}
		}
		_, err = out.Write(nl)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeTableMetaData(out io.Writer, table *Table) error {
	_, err := out.Write([]byte{'['})
	if err != nil {
		return err
	}
	_, err = out.Write([]byte(table.Name))
	if err != nil {
		return err
	}
	for _, field := range table.Fields {
		s := fmt.Sprintf(" %s %s", field.Name, field.Kind)
		_, err = out.Write([]byte(s))
		if err != nil {
			return err
		}
	}
	_, err = out.Write([]byte("\n%\n"))
	return err
}

func writeBool(out io.Writer, value any, kind FieldKind) error {
	v, ok := value.(bool)
	if !ok {
		return fmt.Errorf("e%d:invalid value %v for %q", e143, value, kind)
	}
	t := 'F'
	if v {
		t = 'T'
	}
	_, err := out.Write([]byte{byte(t)})
	return err
}

func writeBytes(out io.Writer, value any, kind FieldKind) error {
	v, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("e%d:invalid value %v for %q", e144, value, kind)
	}
	_, err := out.Write([]byte{'('})
	if err != nil {
		return err
	}
	_, err = out.Write([]byte(hex.EncodeToString(v)))
	if err != nil {
		return err
	}
	_, err = out.Write([]byte{')'})
	return err
}

func writeDateTime(out io.Writer, value any, kind FieldKind, format,
	sentinal string) error {
	v, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("e%d:invalid value %v for %q", e144, value, kind)
	}
	s := v.Format(format)
	if s == sentinal {
		s = "!"
	}
	_, err := out.Write([]byte(s))
	return err
}

func writeInt(out io.Writer, value any, kind FieldKind) error {
	v, ok := value.(int)
	if !ok {
		return fmt.Errorf("e%d:invalid value %v for %q", e145, value, kind)
	}
	var err error
	if v == IntSentinal {
		_, err = out.Write([]byte{'!'})
	} else {
		_, err = out.Write([]byte(strconv.Itoa(v)))
	}
	return err
}

func writeReal(out io.Writer, value any, kind FieldKind,
	decimals int) error {
	v, ok := value.(float64)
	if !ok {
		return fmt.Errorf("e%d:invalid value %v for %q", e141, value, kind)
	}
	var err error
	if gong.IsRealClose(v, RealSentinal) {
		_, err = out.Write([]byte{'!'})
	} else {
		_, err = out.Write([]byte(strconv.FormatFloat(v, 'f', decimals,
			64)))
	}
	return err
}

func writeStr(out io.Writer, value any, kind FieldKind) error {
	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("e%d:invalid value %v for %q", e140, value, kind)
	}
	_, err := out.Write([]byte{'<'})
	if err != nil {
		return err
	}
	_, err = out.Write([]byte(Escape(v)))
	if err != nil {
		return err
	}
	_, err = out.Write([]byte{'>'})
	return err
}
