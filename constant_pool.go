package main

import (
	"fmt"
	"math"
)

// constant pool
// used for next parsing of fields, methods and attributes
type ConstantPool []ReaderParser

type ConstantInfo struct {
	Tag    uint8 // the constant tag
	Reader *ClassFileReader
}

func NewConstantInfo(tag uint8, reader *ClassFileReader) ReaderParser {
	constantInfo := ConstantInfo{tag, reader}
	switch tag {
	case 7:
		// class
		return &ClassInfo{ConstantInfo: constantInfo}
	case 9:
		// field, method, interface method ref
		return &FieldRefInfo{ConstantInfo: constantInfo}
	case 10:
		return &MethodRefInfo{ConstantInfo: constantInfo}
	case 11:
		return &InterfaceMethodRefInfo{ConstantInfo: constantInfo}
	case 8:
		// string
		return &StringInfo{ConstantInfo: constantInfo}
	case 3:
		// integer float
		return &IntegerInfo{ConstantInfo: constantInfo}
	case 4:
		return &FloatInfo{ConstantInfo: constantInfo}
	case 5:
		// long double
		return &LongInfo{ConstantInfo: constantInfo}
	case 6:
		return &DoubleInfo{ConstantInfo: constantInfo}
	case 12:
		// name and type
		return &NameAndTypeInfo{ConstantInfo: constantInfo}
	case 1:
		// utf8
		return &Utf8Info{ConstantInfo: constantInfo}
	case 15:
		// method handle
		return &MethodHandleInfo{ConstantInfo: constantInfo}
	case 16:
		// method type
		return &MethodTypeInfo{ConstantInfo: constantInfo}
	case 18:
		// invoke dynamic
		return &InvokeDynamicInfo{ConstantInfo: constantInfo}
	}
	panic("unsupported tag by jvm8")
}

func (c ConstantPool) GetUTF8(index uint16) string {
	if t := c[index]; t != nil {
		return t.(*Utf8Info).Str
	}

	panic(fmt.Errorf("Invalid constant pool index: %v\n", index))
}

func (c ConstantPool) GetClassName(index uint16) string {
	class := c[index]
	if class == nil {
		panic(fmt.Errorf("Invalid constant pool index: %v\n", index))
	}

	nameIndex := class.(*ClassInfo).NameIndex
	return c.GetUTF8(nameIndex)
}

/**
 * class_info {
 *   u1 tag
 *   u2 name_index
 * }
 */
type ClassInfo struct {
	ConstantInfo
	NameIndex uint16
}

func (c *ClassInfo) Parse() {
	//c.Tag = c.Reader.readUint8()
	c.NameIndex = c.Reader.readUint16()
}

type FieldRefInfo struct {
	ConstantInfo
	NameIndex        uint16
	NameAndTypeIndex uint16
}

func (c *FieldRefInfo) Parse() {
	//c.Tag = c.Reader.readUint8()
	c.NameIndex = c.Reader.readUint16()
	c.NameAndTypeIndex = c.Reader.readUint16()
}

type MethodRefInfo = FieldRefInfo

type InterfaceMethodRefInfo = FieldRefInfo

type StringInfo struct {
	ConstantInfo
	StringIndex uint16
}

func (c *StringInfo) Parse() {
	//c.Tag = c.Reader.readUint8()
	c.StringIndex = c.Reader.readUint16()
}

type IntegerInfo struct {
	ConstantInfo
	Bytes uint32
}

func (c *IntegerInfo) Parse() {
	//c.Tag = c.Reader.readUint8()
	c.Bytes = c.Reader.readUint32()
}

type FloatInfo = IntegerInfo

type LongInfo struct {
	ConstantInfo
	HighBytes uint32
	LowBytes  uint32

	Value int64
}

func (c *LongInfo) Parse() {
	c.HighBytes = c.Reader.readUint32()
	c.LowBytes = c.Reader.readUint32()
	c.Value = int64(c.HighBytes)<<32 | int64(c.LowBytes)
}

type DoubleInfo struct {
	ConstantInfo
	HighBytes uint32
	LowBytes  uint32

	Value float64
}

func (c *DoubleInfo) Parse() {
	c.HighBytes = c.Reader.readUint32()
	c.LowBytes = c.Reader.readUint32()
	c.Value = math.Float64frombits(uint64(c.HighBytes)<<32 | uint64(c.LowBytes))
}

type NameAndTypeInfo struct {
	ConstantInfo
	NameIndex       uint16
	DescriptorIndex uint16
}

func (c *NameAndTypeInfo) Parse() {
	c.NameIndex = c.Reader.readUint16()
	c.DescriptorIndex = c.Reader.readUint16()
}

type Utf8Info struct {
	ConstantInfo
	Bytes []byte
	Str   string
}

func (c *Utf8Info) Parse() {
	length := c.Reader.readUint16()
	c.Bytes = c.Reader.readBytes(uint32(length))
	c.Str = decodeMUTF8(c.Bytes)
}

type MethodHandleInfo struct {
	ConstantInfo
	ReferenceKind  uint8
	ReferenceIndex uint16
}

func (c *MethodHandleInfo) Parse() {
	c.ReferenceKind = c.Reader.readUint8()
	c.ReferenceIndex = c.Reader.readUint16()
}

type MethodTypeInfo struct {
	ConstantInfo
	DescriptorIndex uint16
}

func (c *MethodTypeInfo) Parse() {
	c.DescriptorIndex = c.Reader.readUint16()
}

type InvokeDynamicInfo struct {
	ConstantInfo
	BootstrapMethodAttrIndex uint16
	NameAndTypeIndex         uint16
}

func (c *InvokeDynamicInfo) Parse() {
	c.BootstrapMethodAttrIndex = c.Reader.readUint16()
	c.NameAndTypeIndex = c.Reader.readUint16()
}
