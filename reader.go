package main

import "encoding/binary"

type ClassFileReader struct {
	// the class file binary data
	binData []byte
}

// read a byte
func (self *ClassFileReader) readUint8() uint8 {
	// 读取一个字节
	val := self.binData[0]
	self.binData = self.binData[1:]
	return val
}

// read 2 byte
func (self *ClassFileReader) readUint16() uint16 {
	val := binary.BigEndian.Uint16(self.binData)
	self.binData = self.binData[2:]
	return val
}

// read 4 byte
func (self *ClassFileReader) readUint32() uint32 {
	val := binary.BigEndian.Uint32(self.binData)
	self.binData = self.binData[4:]
	return val
}

// read 8 byte
func (self *ClassFileReader) readUint64() uint64 {
	val := binary.BigEndian.Uint64(self.binData)
	self.binData = self.binData[8:]
	return val
}

// the first 2 bytes represent numbers of bytes to read
func (self *ClassFileReader) readUint16s() []uint16 {
	n := self.readUint16()
	s := make([]uint16, n)
	for i := range s {
		s[i] = self.readUint16()
	}
	return s
}

func (self *ClassFileReader) readBytes(n uint32) []byte {
	bytes := self.binData[:n]
	self.binData = self.binData[n:]
	return bytes
}

func NewReader(bytes []byte) *ClassFileReader {
	return &ClassFileReader{
		binData: bytes,
	}
}
