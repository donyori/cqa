package qc

import (
	"sort"
	"strings"
)

var (
	winStrings []string = []string{
		"win",
		"mfc",
		"directx",
	}
	socketString string   = "sock"
	ideStrings   []string = []string{
		"gcc",
		"visual-c", // Include "visual-c++" and "visual-c#"
		"visual-studio",
		"dll",
		"linker",
		"makefile",
		"assembly",
		"g++",
		"debug",
		"cmake",
		"compile",
		"eclipse",
		"xcode",
		"compilation",
		"gdb",
	}
	otherLanguageStrings []string = []string{
		"c#",
		"python",
		"java",
		"android",
		"php",
		"javascript",
		"cuda",
		".net",
		"sql",
		"ruby",
		"objective-c",
		"matlab",
		"perl",
		"swift",
		"delphi",
		"pascal",
		"directx",
	}
	otherLanguageWholeMatchStrings []string = []string{
		"go",
		"r",
		"scala",
	}
)

func GetLabelsByTags(tags []string) []Label {
	labels := getLabelsByTagsUnsorted(tags)
	if labels == nil {
		return nil
	}
	sort.Slice(labels, func(i, j int) bool {
		return labels[i] < labels[j]
	})
	return labels
}

func GetLabelStringsByTags(tags []string) []string {
	labels := getLabelsByTagsUnsorted(tags)
	if labels == nil {
		return nil
	}
	labelStrings := make([]string, 0, len(labels))
	for _, label := range labels {
		labelStrings = append(labelStrings, strings.ToLower(label.String()))
	}
	sort.Strings(labelStrings)
	return labelStrings
}

func getLabelsByTagsUnsorted(tags []string) []Label {
	if tags == nil {
		return nil
	}
	labelsSet := make(map[Label]bool)
	for _, tag := range tags {
		tag = strings.ToLower(tag)
		for _, label := range KnownLabels {
			containsLabel := false
			labelStringLowerCase := strings.ToLower(label.String())
			switch label {
			case WinLabel:
				for _, winString := range winStrings {
					if strings.Contains(tag, winString) {
						containsLabel = true
						break
					}
				}
			case StlLabel:
				tagWords := strings.Split(tag, "-")
				for _, tagWord := range tagWords {
					if tagWord == labelStringLowerCase {
						containsLabel = true
						break
					}
				}
			case SocketLabel:
				// Sometimes "sock" is used instead of "socket".
				containsLabel = strings.Contains(tag, socketString)
			case IdeLabel:
				for _, ideString := range ideStrings {
					if strings.Contains(tag, ideString) {
						containsLabel = true
						break
					}
				}
			case OtherLanguageLabel:
				for _, ols := range otherLanguageStrings {
					if strings.Contains(tag, ols) {
						containsLabel = true
						break
					}
				}
				if containsLabel {
					break
				}
				for _, olwms := range otherLanguageWholeMatchStrings {
					if tag == olwms {
						containsLabel = true
						break
					}
				}
			default:
				containsLabel = strings.Contains(tag, labelStringLowerCase)
			}
			if containsLabel {
				labelsSet[label] = true
			}
		}
	}
	labels := make([]Label, 0, len(labelsSet))
	for label, ok := range labelsSet {
		if ok {
			labels = append(labels, label)
		}
	}
	return labels
}
