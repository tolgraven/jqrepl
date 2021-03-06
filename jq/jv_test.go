package jq_test

import (
	"testing"

	"github.com/ashb/jqrepl/jq"
	"github.com/cheekybits/is"
)

func TestJvKind(t *testing.T) {
	is := is.New(t)

	cases := []struct {
		*jq.Jv
		jq.JvKind
		string
	}{
		{jq.JvNull(), jq.JV_KIND_NULL, "null"},
		{jq.JvFromString("a"), jq.JV_KIND_STRING, "string"},
	}

	for _, c := range cases {
		defer c.Free()
		is.Equal(c.Kind(), c.JvKind)
		is.Equal(c.Kind().String(), c.string)
	}
}

func TestJvString(t *testing.T) {
	is := is.New(t)

	jv := jq.JvFromString("test")
	defer jv.Free()

	str, err := jv.String()

	is.Equal(str, "test")
	is.NoErr(err)

	i := jv.ToGoVal()

	is.Equal(i, "test")
}

func TestJvStringOnNonStringType(t *testing.T) {
	is := is.New(t)

	// Test that on a non-string value we get a go error, not a C assert
	jv := jq.JvNull()
	defer jv.Free()

	_, err := jv.String()
	is.Err(err)
}

func TestJvFromJSONString(t *testing.T) {
	is := is.New(t)

	jv, err := jq.JvFromJSONString("[]")
	is.NoErr(err)
	is.OK(jv)
	is.Equal(jv.Kind(), jq.JV_KIND_ARRAY)

	jv, err = jq.JvFromJSONString("not valid")
	is.Err(err)
	is.Nil(jv)
}

func TestJvDump(t *testing.T) {
	is := is.New(t)

	jv := jq.JvFromString("test")
	defer jv.Free()

	dump := jv.Copy().Dump(jq.JvPrintNone)

	is.Equal(`"test"`, dump)
	dump = jv.Copy().Dump(jq.JvPrintColour)

	is.Equal([]byte("\x1b[0;32m"+`"test"`+"\x1b[0m"), []byte(dump))
}

func TestJvInvalid(t *testing.T) {
	is := is.New(t)

	jv := jq.JvInvalid()

	is.False(jv.IsValid())

	_, ok := jv.Copy().GetInvalidMessageAsString()
	is.False(ok) // "Expected no Invalid message"

	jv = jv.GetInvalidMessage()
	is.Equal(jv.Kind(), jq.JV_KIND_NULL)
}

func TestJvInvalidWithMessage_string(t *testing.T) {
	is := is.New(t)

	jv := jq.JvInvalidWithMessage(jq.JvFromString("Error message 1"))

	is.False(jv.IsValid())

	msg := jv.Copy().GetInvalidMessage()
	is.Equal(msg.Kind(), jq.JV_KIND_STRING)
	msg.Free()

	str, ok := jv.GetInvalidMessageAsString()
	is.True(ok)
	is.Equal("Error message 1", str)
}

func TestJvInvalidWithMessage_object(t *testing.T) {
	is := is.New(t)

	jv := jq.JvInvalidWithMessage(jq.JvObject())

	is.False(jv.IsValid())

	msg := jv.Copy().GetInvalidMessage()
	is.Equal(msg.Kind(), jq.JV_KIND_OBJECT)
	msg.Free()

	str, ok := jv.GetInvalidMessageAsString()
	is.True(ok)
	is.Equal("{}", str)

}
