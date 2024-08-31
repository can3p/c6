package ast

type ComputableValue interface {
	GetValueType() ValueType
}

type ValueType uint16

const (
	NumberValue    ValueType = 0
	HexColorValue  ValueType = 1
	RGBAColorValue ValueType = 2
	RGBColorValue  ValueType = 3
	ListValue      ValueType = 4
	MapValue       ValueType = 5
)
