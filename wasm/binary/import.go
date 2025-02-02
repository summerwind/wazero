package binary

import (
	"bytes"
	"fmt"
	"io"

	"github.com/tetratelabs/wazero/wasm"
	"github.com/tetratelabs/wazero/wasm/internal/leb128"
)

func decodeImport(r *bytes.Reader) (i *wasm.Import, err error) {
	i = &wasm.Import{}
	if i.Module, _, err = decodeUTF8(r, "import module"); err != nil {
		return nil, err
	}

	if i.Name, _, err = decodeUTF8(r, "import name"); err != nil {
		return nil, err
	}

	b := make([]byte, 1)
	if _, err = io.ReadFull(r, b); err != nil {
		return nil, fmt.Errorf("error decoding import kind: %w", err)
	}

	i.Kind = b[0]
	switch i.Kind {
	case wasm.ImportKindFunc:
		if i.DescFunc, _, err = leb128.DecodeUint32(r); err != nil {
			return nil, fmt.Errorf("error decoding import func typeindex: %w", err)
		}
	case wasm.ImportKindTable:
		if i.DescTable, err = decodeTableType(r); err != nil {
			return nil, fmt.Errorf("error decoding import table desc: %w", err)
		}
	case wasm.ImportKindMemory:
		if i.DescMem, err = decodeMemoryType(r); err != nil {
			return nil, fmt.Errorf("error decoding import mem desc: %w", err)
		}
	case wasm.ImportKindGlobal:
		if i.DescGlobal, err = decodeGlobalType(r); err != nil {
			return nil, fmt.Errorf("error decoding import global desc: %w", err)
		}
	default:
		return nil, fmt.Errorf("%w: invalid byte for importdesc: %#x", ErrInvalidByte, b[0])
	}
	return
}

// encodeImport returns the wasm.Import encoded in WebAssembly 1.0 (MVP) Binary Format.
//
// See https://www.w3.org/TR/wasm-core-1/#binary-import
func encodeImport(i *wasm.Import) []byte {
	data := encodeSizePrefixed([]byte(i.Module))
	data = append(data, encodeSizePrefixed([]byte(i.Name))...)
	data = append(data, i.Kind)
	switch i.Kind {
	case wasm.ImportKindFunc:
		data = append(data, leb128.EncodeUint32(i.DescFunc)...)
	case wasm.ImportKindTable:
		panic("TODO: encodeImportKindTable")
	case wasm.ImportKindMemory:
		panic("TODO: encodeImportKindMemory")
	case wasm.ImportKindGlobal:
		panic("TODO: encodeImportKindGlobal")
	default:
		panic(fmt.Errorf("invalid kind: %#x", i.Kind))
	}
	return data
}
