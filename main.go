package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fileName := ""

	flag.StringVar(&fileName, "f", "./Test.class", "-f file")
	flag.Parse()

	bytes, err := readFile(fileName)
	if err != nil {
		fmt.Printf("read file with error: %s", err.Error())
		os.Exit(1)
	}

	reader := NewReader(bytes)
	// magic
	magic := reader.readUint32()
	fmt.Printf("magic: %x\n", magic)

	// minor
	minorVersion := reader.readUint16()
	majorVersion := reader.readUint16()
	fmt.Printf("major: %d, minor: %d\n", majorVersion, minorVersion)

	// constant_pool_count
	// which value is the pool item count + 1
	cpct := reader.readUint16()
	fmt.Printf("constant pool count: %d\n", cpct)

	// constant pool[]
	// store strings, class or interface names
	constantPool := make(ConstantPool, cpct)
	for i := uint16(1); i < cpct; i++ {
		// read tag
		tag := reader.readUint8()
		tagName := GetConstantPoolTagName(tag)
		fmt.Printf("--tag: %s\n", tagName)

		constantInfo := NewConstantInfo(tag, reader)
		constantInfo.Parse()
		switch constantInfo.(type) {
		case *DoubleInfo, *LongInfo:
			i++
		}

		constantPool[i] = constantInfo
		fmt.Printf("--value: %+v\n", constantInfo)
		fmt.Printf("==================================\n")
		//fmt.Printf("--index: %d\n", i+1)
	}

	// access flags
	accessFlags := reader.readUint16()
	fmt.Printf("access flags: %0*b\n", 16, accessFlags)

	// this class
	thisClass := reader.readUint16() // an index that point to constant pool
	fmt.Printf("this class index: %d\n", thisClass)
	fmt.Printf("this class name: %s\n", constantPool.GetClassName(thisClass))

	// super class
	superClass := reader.readUint16() // an index points to constant pool or zero
	fmt.Printf("super class index: %d\n", superClass)
	fmt.Printf("super class name: %s\n", constantPool.GetClassName(superClass))

	// interface count, u2 means a class can implement 2^16-1 interfaces by jvm
	interfaceCount := reader.readUint16()
	fmt.Printf("interface count: %d\n", interfaceCount)

	// interface array
	for i := uint16(0); i < interfaceCount; i++ {
		// each interface is a struct of constant_class_info
		nameIndex := reader.readUint16()
		fmt.Printf("--interface name_index: %d\n", nameIndex)
	}

	// fields count
	fieldsCount := reader.readUint16()
	fmt.Printf("fields count: %d\n", fieldsCount)
	fields := make([]ClassMemberInfo, fieldsCount)
	fmt.Printf("fields:\n")
	for i := range fields {
		field := *NewClassMemberInfo(&constantPool, reader)
		field.Parse()
		fields[i] = field
		fmt.Printf("--field: %+v\n", field)
	}

	methodsCount := reader.readUint16()
	fmt.Printf("methods count: %d\n", methodsCount)
	methods := make([]ClassMemberInfo, methodsCount)
	fmt.Printf("methods:\n")
	for i := range methods {
		method := *NewClassMemberInfo(&constantPool, reader)
		method.Parse()
		methods[i] = method
		fmt.Printf("--method name: %s\n", constantPool.GetUTF8(method.NameIndex))
		fmt.Printf("--method: %+v\n", method)
		fmt.Printf("--method attributes: %v", method.Attributes)
	}

	attributesCount := reader.readUint16()
	fmt.Printf("attributes count: %d\n", attributesCount)
	attrs := make([]AttributeInfo, attributesCount)
	fmt.Printf("attributes:\n")
	for i := range attrs {
		attrInfo := ReadAttribute(reader, &constantPool)
		attrInfo.ParseAttr()
		attrs[i] = attrInfo

		fmt.Printf("--attribute: %+v\n", attrInfo)
	}
}

func readFile(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}
