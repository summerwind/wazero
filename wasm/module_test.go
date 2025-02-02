package wasm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFunctionType_String(t *testing.T) {
	for _, tc := range []struct {
		functype *FunctionType
		exp      string
	}{
		{functype: &FunctionType{}, exp: "null_null"},
		{functype: &FunctionType{Params: []ValueType{ValueTypeI32}}, exp: "i32_null"},
		{functype: &FunctionType{Params: []ValueType{ValueTypeI32, ValueTypeF64}}, exp: "i32f64_null"},
		{functype: &FunctionType{Params: []ValueType{ValueTypeF32, ValueTypeI32, ValueTypeF64}}, exp: "f32i32f64_null"},
		{functype: &FunctionType{Results: []ValueType{ValueTypeI64}}, exp: "null_i64"},
		{functype: &FunctionType{Results: []ValueType{ValueTypeI64, ValueTypeF32}}, exp: "null_i64f32"},
		{functype: &FunctionType{Results: []ValueType{ValueTypeF32, ValueTypeI32, ValueTypeF64}}, exp: "null_f32i32f64"},
		{functype: &FunctionType{Params: []ValueType{ValueTypeI32}, Results: []ValueType{ValueTypeI64}}, exp: "i32_i64"},
		{functype: &FunctionType{Params: []ValueType{ValueTypeI64, ValueTypeF32}, Results: []ValueType{ValueTypeI64, ValueTypeF32}}, exp: "i64f32_i64f32"},
		{functype: &FunctionType{Params: []ValueType{ValueTypeI64, ValueTypeF32, ValueTypeF64}, Results: []ValueType{ValueTypeF32, ValueTypeI32, ValueTypeF64}}, exp: "i64f32f64_f32i32f64"},
	} {
		tc := tc
		t.Run(tc.functype.String(), func(t *testing.T) {
			require.Equal(t, tc.exp, tc.functype.String())
		})
	}
}

func TestSectionIDName(t *testing.T) {
	tests := []struct {
		name     string
		input    SectionID
		expected string
	}{
		{"custom", SectionIDCustom, "custom"},
		{"type", SectionIDType, "type"},
		{"import", SectionIDImport, "import"},
		{"function", SectionIDFunction, "function"},
		{"table", SectionIDTable, "table"},
		{"memory", SectionIDMemory, "memory"},
		{"global", SectionIDGlobal, "global"},
		{"export", SectionIDExport, "export"},
		{"start", SectionIDStart, "start"},
		{"element", SectionIDElement, "element"},
		{"code", SectionIDCode, "code"},
		{"data", SectionIDData, "data"},
		{"unknown", 100, "unknown"},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, SectionIDName(tc.input))
		})
	}
}

func TestExportKindName(t *testing.T) {
	tests := []struct {
		name     string
		input    ExportKind
		expected string
	}{
		{"func", ExportKindFunc, "func"},
		{"table", ExportKindTable, "table"},
		{"mem", ExportKindMemory, "mem"},
		{"global", ExportKindGlobal, "global"},
		{"unknown", 100, "unknown"},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, exportKindName(tc.input))
		})
	}
}

func TestValueTypeName(t *testing.T) {
	tests := []struct {
		name     string
		input    ValueType
		expected string
	}{
		{"i32", ValueTypeI32, "i32"},
		{"i64", ValueTypeI64, "i64"},
		{"f32", ValueTypeF32, "f32"},
		{"f64", ValueTypeF64, "f64"},
		{"unknown", 100, "unknown"},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, ValueTypeName(tc.input))
		})
	}
}

func TestModule_allDeclarations(t *testing.T) {
	for i, tc := range []struct {
		module            *Module
		expectedFunctions []Index
		expectedGlobals   []*GlobalType
		expectedMemories  []*MemoryType
		expectedTables    []*TableType
	}{
		// Functions.
		{
			module: &Module{
				ImportSection:   []*Import{{Kind: ImportKindFunc, DescFunc: 10000}},
				FunctionSection: []Index{10, 20, 30},
			},
			expectedFunctions: []Index{10000, 10, 20, 30},
		},
		{
			module: &Module{
				FunctionSection: []Index{10, 20, 30},
			},
			expectedFunctions: []Index{10, 20, 30},
		},
		{
			module: &Module{
				ImportSection: []*Import{{Kind: ImportKindFunc, DescFunc: 10000}},
			},
			expectedFunctions: []Index{10000},
		},
		// Globals.
		{
			module: &Module{
				ImportSection: []*Import{{Kind: ImportKindGlobal, DescGlobal: &GlobalType{Mutable: false}}},
				GlobalSection: []*Global{{Type: &GlobalType{Mutable: true}}},
			},
			expectedGlobals: []*GlobalType{{Mutable: false}, {Mutable: true}},
		},
		{
			module: &Module{
				GlobalSection: []*Global{{Type: &GlobalType{Mutable: true}}},
			},
			expectedGlobals: []*GlobalType{{Mutable: true}},
		},
		{
			module: &Module{
				ImportSection: []*Import{{Kind: ImportKindGlobal, DescGlobal: &GlobalType{Mutable: false}}},
			},
			expectedGlobals: []*GlobalType{{Mutable: false}},
		},
		// Memories.
		{
			module: &Module{
				ImportSection: []*Import{{Kind: ImportKindMemory, DescMem: &LimitsType{Min: 1}}},
				MemorySection: []*MemoryType{{Min: 100}},
			},
			expectedMemories: []*MemoryType{{Min: 1}, {Min: 100}},
		},
		{
			module: &Module{
				ImportSection: []*Import{{Kind: ImportKindMemory, DescMem: &LimitsType{Min: 1}}},
			},
			expectedMemories: []*MemoryType{{Min: 1}},
		},
		{
			module: &Module{
				MemorySection: []*MemoryType{{Min: 100}},
			},
			expectedMemories: []*MemoryType{{Min: 100}},
		},
		// Tables.
		{
			module: &Module{
				ImportSection: []*Import{{Kind: ImportKindTable, DescTable: &TableType{Limit: &LimitsType{Min: 1}}}},
				TableSection:  []*TableType{{Limit: &LimitsType{Min: 10}}},
			},
			expectedTables: []*TableType{{Limit: &LimitsType{Min: 1}}, {Limit: &LimitsType{Min: 10}}},
		},
		{
			module: &Module{
				ImportSection: []*Import{{Kind: ImportKindTable, DescTable: &TableType{Limit: &LimitsType{Min: 1}}}},
			},
			expectedTables: []*TableType{{Limit: &LimitsType{Min: 1}}},
		},
		{
			module: &Module{
				TableSection: []*TableType{{Limit: &LimitsType{Min: 10}}},
			},
			expectedTables: []*TableType{{Limit: &LimitsType{Min: 10}}},
		},
	} {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			functions, globals, memories, tables := tc.module.allDeclarations()
			require.Equal(t, tc.expectedFunctions, functions)
			require.Equal(t, tc.expectedGlobals, globals)
			require.Equal(t, tc.expectedTables, tables)
			require.Equal(t, tc.expectedMemories, memories)
		})
	}
}
