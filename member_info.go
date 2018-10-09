package main

// class member
// a java class may contains fields and methods
// which by jvm8 stored in two data structures
//
//  field_info {
//    u2             access_flags;
//    u2             name_index;
//    u2             descriptor_index;
//    u2             attributes_count;
//    attribute_info attributes[attributes_count];
//  }
//
//  method_info {
//    u2             access_flags;
//    u2             name_index;
//    u2             descriptor_index;
//    u2             attributes_count;
//    attribute_info attributes[attributes_count];
//  }
type ClassMemberInfo struct {
	CP     *ConstantPool
	Reader *ClassFileReader

	AccessFlags     uint16
	NameIndex       uint16
	DescriptorIndex uint16
	AttributesCount uint16
	Attributes      []AttributeInfo
}

func (c *ClassMemberInfo) Parse() {
	c.AccessFlags = c.Reader.readUint16()
	c.NameIndex = c.Reader.readUint16()
	c.DescriptorIndex = c.Reader.readUint16()
	c.AttributesCount = c.Reader.readUint16()

	// TODO attributes read
}

func NewClassMemberInfo(cp *ConstantPool, reader *ClassFileReader) *ClassMemberInfo {
	return &ClassMemberInfo{
		CP:     cp,
		Reader: reader,
	}
}
