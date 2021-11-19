// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file defines the IDs for PCDATA and FUNCDATA instructions
// in Go binaries. It is included by assembly sources, so it must
// be written using #defines.
//
// These must agree with symtab.go and ../cmd/internal/objabi/funcdata.go.

#define PCDATA_UnsafePoint 0
#define PCDATA_StackMapIndex 1
#define PCDATA_InlTreeIndex 2

#define FUNCDATA_ArgsPointerMaps 0 /* garbage collector blocks */
#define FUNCDATA_LocalsPointerMaps 1
#define FUNCDATA_StackObjects 2
#define FUNCDATA_InlTree 3
#define FUNCDATA_OpenCodedDeferInfo 4 /* info for func with open-coded defers */
#define FUNCDATA_ArgInfo 5

// Pseudo-assembly statements.

// GO_ARGS, GO_RESULTS_INITIALIZED, and NO_LOCAL_POINTERS are macros
// that communicate to the runtime information about the location and liveness
// of pointers in an assembly function's arguments, results, and stack frame.
// This communication is only required in assembly functions that make calls
// to other functions that might be preempted or grow the stack.
// NOSPLIT functions that make no calls do not need to use these macros.

// GO_ARGS indicates that the Go prototype for this assembly function
// defines the pointer map for the function's arguments.
// GO_ARGS should be the first instruction in a function that uses it.
// It can be omitted if there are no arguments at all.
// GO_ARGS is inserted implicitly by the linker for any function whose
// name starts with a middle-dot and that also has a Go prototype; it
// is therefore usually not necessary to write explicitly.
#define GO_ARGS	FUNCDATA $FUNCDATA_ArgsPointerMaps, go_args_stackmap(SB)

// GO_RESULTS_INITIALIZED indicates that the assembly function
// has initialized the stack space for its results and that those results
// should be considered live for the remainder of the function.
#define GO_RESULTS_INITIALIZED	PCDATA $PCDATA_StackMapIndex, $1

// NO_LOCAL_POINTERS indicates that the assembly function stores
// no pointers to heap objects in its local stack variables.
#define NO_LOCAL_POINTERS	FUNCDATA $FUNCDATA_LocalsPointerMaps, no_pointers_stackmap(SB)

// ArgsSizeUnknown is set in Func.argsize to mark all functions
// whose argument size is unknown (C vararg functions, and
// assembly code without an explicit specification).
// This value is generated by the compiler, assembler, or linker.
#define ArgsSizeUnknown 0x80000000
