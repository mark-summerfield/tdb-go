// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: Apache-2.0

package tdb

type fieldFlag uint8

const (
	notNullFlag  fieldFlag = 0b000
	nullableFlag fieldFlag = 0b001
	uniqueFlag   fieldFlag = 0b010
	autoFlag     fieldFlag = 0b110
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

func (me fieldFlag) isAuto() bool {
	return me&autoFlag == autoFlag
}
