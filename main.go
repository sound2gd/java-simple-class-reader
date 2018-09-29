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
	for i := uint16(0); i < cpct-1; i++ {
		// read tag
		tag := reader.readUint8()
		tagName := GetConstantPoolTagName(tag)
		fmt.Printf("--index: %d\n", i+1)
		fmt.Printf("--tag: %s\n", tagName)
		dispatchConstantRead(tag, reader)
	}

	// access flags
	accessFlags := reader.readUint16()
	fmt.Printf("access flags: %0*b\n", 16, accessFlags)

	// this class
	thisClass := reader.readUint16() // an index that point to constant pool
	fmt.Printf("this class index: %d\n", thisClass)

	// super class
	superClass := reader.readUint16() // an index points to constant pool or zero
	fmt.Printf("super class index: %d\n", superClass)

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
}

func dispatchConstantRead(tag uint8, reader *ClassFileReader) {
	switch tag {
	case 7:
		// class
		nameIndex := reader.readUint16()
		fmt.Printf("--name_index: %d\n", nameIndex)
	case 9, 10, 11:
		// field, method, interface method ref
		classIndex := reader.readUint16()
		fmt.Printf("--class_index: %d\n", classIndex)
		nameAndTypeIndex := reader.readUint16()
		fmt.Printf("--name_and_type_index: %d\n", nameAndTypeIndex)
	case 8:
		// string
		stringIndex := reader.readUint16()
		fmt.Printf("--string_index: %d\n", stringIndex)
	case 3, 4:
		// integer float
		bytes := reader.readUint32()
		fmt.Printf("--bytes: %d\n", bytes)
	case 5, 6:
		// long double
		highBytes := reader.readUint32()
		lowBytes := reader.readUint32()
		fmt.Printf("--highBytes: %d\n", highBytes)
		fmt.Printf("--lowBytes: %d\n", lowBytes)
		fmt.Printf("--converted: %d\n", int64(highBytes)<<32|int64(lowBytes))
	case 12:
		// name and type
		nameIndex := reader.readUint16()
		fmt.Printf("--name_index: %d\n", nameIndex)
		decriptorIndex := reader.readUint16()
		fmt.Printf("--descriptor_index: %d\n", decriptorIndex)
	case 1:
		// utf8
		length := reader.readUint16()
		fmt.Printf("--lenth: %d\n", length)
		bytes := reader.readBytes(uint32(length))
		fmt.Printf("--bytes: %s\n", decodeMUTF8(bytes))
	case 15:
		// method handle
		referenceKind := reader.readUint8()
		fmt.Printf("--reference_kind: %d\n", referenceKind)
		referenceIndex := reader.readUint16()
		fmt.Printf("--reference_index: %d\n", referenceIndex)
	case 16:
		// method type
		descriptorIndex := reader.readUint8()
		fmt.Printf("--decriptor_index: %d\n", descriptorIndex)
	case 18:
		// invoke dynamic
		bootStrapMethodAttrIndex := reader.readUint16()
		fmt.Printf("--bootstrap_method_attr_index: %d\n", bootStrapMethodAttrIndex)
		nameAndTypeIndex := reader.readUint16()
		fmt.Printf("--name_and_type_index: %d\n", nameAndTypeIndex)
	}
	fmt.Println("========================================")
}

func readFile(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}
