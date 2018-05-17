package qc

import (
	"errors"
)

type Label int8

const (
	ArrayLabel Label = iota
	PointerLabel
	QtLabel
	LinuxLabel
	TemplateLabel
	StringLabel
	WinLabel
	BoostLabel
	MultithreadingLabel
	OpenCvLabel
	VectorLabel
	StructLabel
	AlgorithmLabel
	OpenGlLabel
	StlLabel
	FunctionLabel
	SocketLabel
	ClassLabel
	MemoryLabel
	IdeLabel
	OtherLanguageLabel
)

var (
	KnownLabels = [...]Label{
		ArrayLabel,
		PointerLabel,
		QtLabel,
		LinuxLabel,
		TemplateLabel,
		StringLabel,
		WinLabel,
		BoostLabel,
		MultithreadingLabel,
		OpenCvLabel,
		VectorLabel,
		StructLabel,
		AlgorithmLabel,
		OpenGlLabel,
		StlLabel,
		FunctionLabel,
		SocketLabel,
		ClassLabel,
		MemoryLabel,
		IdeLabel,
		OtherLanguageLabel,
	}

	ErrUnknownLabel error = errors.New("label is unknown")

	labelStrings = [...]string{
		"array",
		"pointer",
		"qt",
		"linux",
		"template",
		"string",
		"win",
		"boost",
		"multithreading",
		"opencv",
		"vector",
		"struct",
		"algorithm",
		"opengl",
		"stl",
		"function",
		"socket",
		"class",
		"memory",
		"ide",
		"other-language",
	}
)

func (l Label) IsKnown() bool {
	return l >= ArrayLabel && l <= OtherLanguageLabel
}

func (l Label) IsValid() bool {
	return l.IsKnown()
}

func (l Label) String() string {
	if !l.IsKnown() {
		return "unknown"
	}
	if !l.IsValid() {
		return "invalid"
	}
	return labelStrings[l]
}
