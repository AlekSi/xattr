package xattr_test

import (
        "io/ioutil"
        "os"
        "testing"

        . "."
        . "launchpad.net/gocheck"
)

func TestXattr(t *testing.T) { TestingT(t) }

type F struct {
        f    string
        attr string
}

var _ = Suite(&F{})

func (f *F) SetUpTest(c *C) {
        file, err := ioutil.TempFile("", "test_xattr_")
        c.Assert(err, IsNil)
        err = file.Close()
        c.Assert(err, IsNil)
        f.f = file.Name()
        f.attr = "test xattr"
}

func (f *F) TearDownTest(c *C) {
        if !c.Failed() {
                err := os.Remove(f.f)
                c.Assert(err, IsNil)
        }
}

func (f *F) TestFlow(c *C) {
        data := []byte("test xattr data")

        attrs, err := Listxattr(f.f)
        c.Check(err, IsNil)
        c.Check(attrs, DeepEquals, []string{})

        err = Setxattr(f.f, f.attr, data)
        c.Check(err, IsNil)

        attrs, err = Listxattr(f.f)
        c.Check(err, IsNil)
        c.Check(attrs, DeepEquals, []string{f.attr})

        data1, err := Getxattr(f.f, f.attr)
        c.Check(err, IsNil)
        c.Check(data1, DeepEquals, data)

        data1, err = Getxattr(f.f, "test other xattr")
        c.Check(err, ErrorMatches, "*. attribute not found")
        c.Check(err, FitsTypeOf, &XAttrError{})
        c.Check(data1, IsNil)

        err = Removexattr(f.f, f.attr)
        c.Check(err, IsNil)

        attrs, err = Listxattr(f.f)
        c.Check(err, IsNil)
        c.Check(attrs, DeepEquals, []string{})
}

func (f *F) TestEmptyAttr(c *C) {
        data := []byte{}

        err := Setxattr(f.f, f.attr, data)
        c.Check(err, IsNil)

        attrs, err := Listxattr(f.f)
        c.Check(err, IsNil)
        c.Check(attrs, DeepEquals, []string{f.attr})

        data1, err := Getxattr(f.f, f.attr)
        c.Check(err, IsNil)
        c.Check(data1, DeepEquals, data)

        err = Removexattr(f.f, f.attr)
        c.Check(err, IsNil)

        attrs, err = Listxattr(f.f)
        c.Check(err, IsNil)
        c.Check(attrs, DeepEquals, []string{})
}

func (f *F) TestNoFile(c *C) {
        fn := "no-such-file"
        data := []byte("test_xattr data")

        attrs, err := Listxattr(fn)
        c.Check(err, ErrorMatches, "*. no such file or directory")
        c.Check(err, FitsTypeOf, &XAttrError{})
        c.Check(attrs, IsNil)

        err = Setxattr(fn, f.attr, data)
        c.Check(err, ErrorMatches, "*. no such file or directory")
        c.Check(err, FitsTypeOf, &XAttrError{})

        attrs, err = Listxattr(fn)
        c.Check(err, ErrorMatches, "*. no such file or directory")
        c.Check(err, FitsTypeOf, &XAttrError{})
        c.Check(attrs, IsNil)

        data1, err := Getxattr(fn, f.attr)
        c.Check(err, ErrorMatches, "*. no such file or directory")
        c.Check(err, FitsTypeOf, &XAttrError{})
        c.Check(data1, IsNil)

        err = Removexattr(fn, f.attr)
        c.Check(err, ErrorMatches, "*. no such file or directory")
        c.Check(err, FitsTypeOf, &XAttrError{})

        attrs, err = Listxattr(fn)
        c.Check(err, ErrorMatches, "*. no such file or directory")
        c.Check(err, FitsTypeOf, &XAttrError{})
        c.Check(attrs, IsNil)
}
