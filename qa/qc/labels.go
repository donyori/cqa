package qc

import (
	"errors"
	"strings"
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

	UnknownLabel Label = -1
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

func ParseLabel(labelString string, doesIgnoreCase bool) Label {
	if doesIgnoreCase {
		labelString = strings.ToLower(labelString)
	}
	if labelString == "" || labelString == "invalid" {
		return -2
	}
	if labelString == "unknown" {
		return -1
	}
	for i, s := range labelStrings {
		if labelString == s {
			return Label(i)
		}
	}
	return -1
}

func (l Label) IsKnown() bool {
	return l >= ArrayLabel && l <= OtherLanguageLabel
}

func (l Label) IsValid() bool {
	return l.IsKnown() || l == UnknownLabel
}

func (l Label) String() string {
	if !l.IsValid() {
		return "invalid"
	}
	if !l.IsKnown() {
		return "unknown"
	}
	return labelStrings[l]
}
