package wasm

import (
	"github.com/tetratelabs/wazero/wasm/leb128"
)

// nameSectionPrefix a length-prefixed 'name'. subsection data should follow.
var nameSectionPrefix = encodeName("name")

// encodeCustomNameSection encodes a possibly empty buffer representing the "name" wasm.Module CustomSection.
// See https://www.w3.org/TR/wasm-core-1/#name-section%E2%91%A0
func encodeCustomNameSection(m *Module) []byte {
	funcCount, funcNameCount := uint32(0), uint32(0)
	var funcNameEntries []byte
	for idx, i := range m.ImportSection {
		if i.Desc.Kind != ImportKindFunction {
			continue
		}
		funcCount++
		if i.Desc.FuncName != "" {
			funcNameCount++
			funcNameEntries = append(funcNameEntries, encodeNameMapEntry(uint32(idx), i.Desc.FuncName)...)
		}
	}
	var data = nameSectionPrefix
	if m.Name != "" {
		// See https://www.w3.org/TR/wasm-core-1/#binary-modulenamesec
		data = append(data, encodeNameSubsection(uint8(0), encodeName(m.Name))...)
	}
	if funcNameCount > 0 {
		// See https://www.w3.org/TR/wasm-core-1/#binary-funcnamesec
		content := leb128.EncodeUint32(funcNameCount)
		content = append(content, funcNameEntries...)
		data = append(data, encodeNameSubsection(uint8(1), content)...)
	}
	if funcCount > 0 {
		// See https://www.w3.org/TR/wasm-core-1/#binary-localnamesec
		content := leb128.EncodeUint32(funcCount)
		for i := uint32(0); i < funcCount; i++ {
			content = append(content, leb128.EncodeUint32(i)...)
			content = append(content, 0) // TODO: actually append locals!
		}
		data = append(data, encodeNameSubsection(uint8(2), content)...)
	}

	// Finally, make the header
	dataSize := leb128.EncodeUint32(uint32(len(data)))
	header := append([]byte{SectionIDCustom}, dataSize...)
	return append(header, data...)
}

// This returns a buffer encoding the given subsection
// See https://www.w3.org/TR/wasm-core-1/#subsections%E2%91%A0
func encodeNameSubsection(subsectionID uint8, content []byte) []byte {
	contentSizeInBytes := leb128.EncodeUint32(uint32(len(content)))
	result := []byte{subsectionID}
	result = append(result, contentSizeInBytes...)
	result = append(result, content...)
	return result
}

// encodeNameMapEntry encodes the index and name prefixed by their size.
// See https://www.w3.org/TR/wasm-core-1/#binary-namemap
func encodeNameMapEntry(i uint32, name string) []byte {
	return append(leb128.EncodeUint32(i), encodeName(name)...)
}

// encodeName encodes the name prefixed by its size.
func encodeName(name string) []byte {
	nameBytes := []byte(name)
	nameSize := leb128.EncodeUint32(uint32(len(nameBytes)))
	content := append(nameSize, nameBytes...)
	return content
}
