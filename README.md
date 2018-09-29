# Simple Java Class Reader

> learning by inventing 

this project aimed at learning the `class file format` in [jvm8 spec](https://docs.oracle.com/javase/specs/jvms/se8/jvms8.pdf).
main code written in `golang`

## Class File Format

```
ClassFile {
    u4             magic;
    u2             minor_version;
    u2             major_version;
    u2             constant_pool_count;
    cp_info        constant_pool[constant_pool_count-1];
    u2             access_flags;
    u2             this_class;
    u2             super_class;
    u2             interfaces_count;
    u2             interfaces[interfaces_count];
    u2             fields_count;
    field_info     fields[fields_count];
    u2             methods_count;
    method_info    methods[methods_count];
    u2             attributes_count;
    attribute_info attributes[attributes_count];
}
```

## project progress

- [x] magic
- [x] major and minor version
- [x] constant pool
- [x] access flags
- [x] this class
- [x] super class
- [x] interfaces
- [ ] fields
- [ ] methods
- [ ] attributes


## how to run

```bash
$ go build .
$ ./java-simple-class-reader -f /path/to/class/file
```

Of course, there is a example class file named `Test.class`.

