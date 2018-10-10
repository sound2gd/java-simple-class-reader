package main

// jvm class stores a lot information in attributes
// such as annotations
//
// @Autowired
// private DemoService demoService;
//
// there are 5 critical attributes to jvm8
// - Code
// - ConstantValue
// - StackMapTable
// - Exceptions
// - BootstrapMethods
//
// and 12 critical attributes to JavaSE library
// and 6 useful attributes to tools
// which are not listed
type AttributeInfo interface {
	ParseAttr()
}

// jvm use attribute name to distinguish all attributes
// attribute name is a Constant_UTF8_info from constant_pool
// each attribute has its own data structure
func ReadAttribute(reader *ClassFileReader, cp *ConstantPool) AttributeInfo {
	attrNameIndex := reader.readUint16()
	attrName := cp.GetUTF8(attrNameIndex)
	attrLen := reader.readUint32()

	general := AttributeInfoGeneral{cp, reader, attrNameIndex, attrLen, attrName}
	switch attrName {
	case "ConstantValue":
		return &ConstantValueAttr{
			AttributeInfoGeneral: general,
		}
	default:
		// unimplemented attributes
		reader.readBytes(attrLen)
		return &general
	}
}

//====================================================================
// attributes implementation
type AttributeInfoGeneral struct {
	CP     *ConstantPool
	Reader *ClassFileReader

	AttributeNameIndex uint16
	AttributeLength    uint32

	AttributeName string // name resolved from constant pool
}

func (c *AttributeInfoGeneral) ParseAttr() {
}

type ConstantValueAttr struct {
	AttributeInfoGeneral
	ConstantValueIndex uint16
}

func (c *ConstantValueAttr) ParseAttr() {
	c.ConstantValueIndex = c.Reader.readUint16()
}
