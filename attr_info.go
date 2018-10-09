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
