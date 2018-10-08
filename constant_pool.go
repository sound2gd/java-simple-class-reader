package main

// constant pool
// used for next parsing of fields, methods and attributes
type ConstantPool []ConstantInfo

type ConstantInfo struct {
	Tag uint8 // the constant tag
	Reader ClassFileReader
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
	c.Tag = c.Reader.readUint8()
	c.NameIndex = c.Reader.readUint16()
}




