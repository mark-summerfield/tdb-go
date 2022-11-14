// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

type fieldFlag uint8

const (
	notNullFlag  fieldFlag = 0b00
	nullableFlag fieldFlag = 0b01
	uniqueFlag   fieldFlag = 0b10
)

func (me fieldFlag) with(flag fieldFlag) fieldFlag {
	return me | flag
}

func (me fieldFlag) isNullable() bool {
	return me&nullableFlag != 0
}

func (me fieldFlag) isUnique() bool {
	return me&uniqueFlag != 0
}
