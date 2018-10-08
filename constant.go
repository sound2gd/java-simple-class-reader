package main

import (
	"fmt"
	"unicode/utf16"
)

func GetConstantPoolTagName(value uint8) string {
	switch value {
	case 7:
		return "CONSTANT_Class"
	case 9:
		return "CONSTANT_FieldRef"
	case 10:
		return "CONSTANT_MethodRef"
	case 11:
		return "CONSTANT_InterfaceMethodRef"
	case 8:
		return "CONSTANT_String"
	case 3:
		return "CONSTANT_Integer"
	case 4:
		return "CONSTANT_Float"
	case 5:
		return "CONSTANT_Long"
	case 6:
		return "CONSTANT_Double"
	case 12:
		return "CONSTANT_NameAndType"
	case 1:
		return "CONSTANT_Utf8"
	case 15:
		return "CONSTANT_MethodHandle"
	case 16:
		return "CONSTANT_MethodType"
	case 18:
		return "CONSTANT_InvokeDynamic"
	}
	panic("unsupported tag by openjdk8")
}

func decodeMUTF8(bytes []byte) string {
	utflen := len(bytes)
	chars := make([]uint16, utflen)

	var c, char2, char3 uint16
	count := 0
	charsCount := 0

	for count < utflen {
		c = uint16(bytes[count])
		if c > 127 {
			break
		}
		count++
		chars[charsCount] = c
		charsCount++
	}

	for count < utflen {
		c = uint16(bytes[count])
		switch c >> 4 {
		case 0, 1, 2, 3, 4, 5, 6, 7:
			count++
			chars[charsCount] = c
			charsCount++
		case 12, 13:
			count += 2
			if count > utflen {
				panic("malformed input: partial character at end!")
			}
			char2 = uint16(bytes[count-1])
			if char2&0xC0 != 0x80 {
				panic(fmt.Errorf("malformed input around byte %v", count))
			}
			chars[charsCount] = c&0x1F<<6 | char2&0x3F
			charsCount++
		case 14:
			count += 3
			if count > utflen {
				panic("malformed input: partial character at end")
			}
			char2 = uint16(bytes[count-2])
			char3 = uint16(bytes[count-1])
			if char2&0xC0 != 0x80 || char3&0xC0 != 0x80 {
				panic(fmt.Errorf("malformed input around byte %v", (count - 1)))
			}
			chars[charsCount] = c&0x0F<<12 | char2&0x3F<<6 | char3&0x3F<<0
			charsCount++
		default:
			panic(fmt.Errorf("malformed input around byte %v", count))
		}
	}

	chars = chars[0:charsCount]
	runes := utf16.Decode(chars)
	return string(runes)
}
