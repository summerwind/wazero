package wasm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeNameSection(t *testing.T) {
	m := &Module{
		Name:        "simple",
		TypeSection: []*FunctionType{{}},
		ImportSection: []*ImportSegment{
			{
				Module: "",
				Name:   "Hello",
				Desc: &ImportDesc{
					Kind:          ImportKindFunction,
					FuncTypeIndex: 0,
					FuncName:      "hello",
				},
			},
		},
	}

	// TIP: the below is the binary suffix of `wat2wasm --debug-names --debug-parser -v simple.wat` where simple.wat
	// contains the same text as the above inlined text format. Ex.
	//	(module $simple
	//		(import "" "Hello" (func $hello))
	//		(start $hello)
	//	)
	require.Equal(t, []byte{
		SectionIDCustom,
		0x1d, // 24 bytes after this point
		0x04, // the custom second name "name" is 4 bytes long
		'n', 'a', 'm', 'e',
		0x00, // module subsection ID zero
		0x07, // 7 bytes to follow
		0x06, // the module name simple is 6 bytes long
		's', 'i', 'm', 'p', 'l', 'e',
		0x01, // function subsection ID one
		0x08, // 8 bytes to follow
		0x01, // one function name
		0x00, // the function index is zero
		0x05, // the function name hello is 5 bytes long
		'h', 'e', 'l', 'l', 'o',
		0x02, // local subsection ID two
		0x03, // 3 bytes to follow
		0x01, // one function
		0x00, // the function index is zero
		0x00, // no locals
	}, encodeCustomNameSection(m))
}

// TestEncodeNameSection_OnlyFuncName shows that we don't rely on the module name being present. For example, this isn't
// encoded in TinyGo.
func TestEncodeNameSection_OnlyFuncName(t *testing.T) {
	func0, func1 := "runtime.args_sizes_get", "runtime.fd_write"
	i32 := ValueTypeI32
	type1 := &FunctionType{Params: []ValueType{i32, i32}, Results: []ValueType{i32}}
	type2 := &FunctionType{Params: []ValueType{i32, i32, i32, i32}, Results: []ValueType{i32}}

	m := &Module{
		TypeSection: []*FunctionType{type1, type2},
		ImportSection: []*ImportSegment{
			{
				Module: "wasi_snapshot_preview1",
				Name:   "args_sizes_get",
				Desc: &ImportDesc{
					Kind:          ImportKindFunction,
					FuncTypeIndex: 0,
					FuncName:      func0,
				},
			},
			{
				Module: "wasi_snapshot_preview1",
				Name:   "fd_write",
				Desc: &ImportDesc{
					Kind:          ImportKindFunction,
					FuncTypeIndex: 1,
					FuncName:      func1,
				},
			},
		},
	}

	expected := append(append(append([]byte{
		SectionIDCustom,
		0x39, // 57 bytes after this point
		0x04, // the custom second name "name" is 4 bytes long
		'n', 'a', 'm', 'e',
		0x01, // function subsection ID one
		// length includes overhead for size in bytes of the function name count, plus index + length prefix per name
		byte(1 + 2 + 2 + len(func0) + len(func1)),
		0x02, // two function names
	},
		append([]byte{0x00 /* funcIndex */, byte(len(func0))}, func0...)...),
		append([]byte{0x01 /* funcIndex */, byte(len(func1))}, func1...)...),
		0x02, // local subsection ID two
		0x05, // 5 bytes to follow
		0x02, // two functions
		0x00, // function index is zero
		0x00, // no locals
		0x01, // function index is one
		0x00, // no locals
	)

	require.Equal(t, expected, encodeCustomNameSection(m))
}

func TestEncodeNameSubsection(t *testing.T) {
	subsectionID := uint8(1)
	name := "simple"
	require.Equal(t, []byte{
		subsectionID,
		byte(1 + 6), // 1 is the size of 6 in LEB128 encoding
		6, 's', 'i', 'm', 'p', 'l', 'e'}, encodeNameSubsection(subsectionID, encodeName(name)))
}

func TestEncodeNameMapEntry(t *testing.T) {
	index := uint32(1)
	name := "hello"
	require.Equal(t, []byte{byte(index), 5, 'h', 'e', 'l', 'l', 'o'}, encodeNameMapEntry(index, name))
}

func TestEncodeName(t *testing.T) {
	// We expect a length (in LEB128) prefixed string encoding
	require.Equal(t, []byte{5, 'h', 'e', 'l', 'l', 'o'}, encodeName("hello"))
}
