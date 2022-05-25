package segmenter

import (
	"strings"

	"github.com/huichen/sego"
)

var seg sego.Segmenter

func init() {
	seg.LoadDictionary("dictionary.txt")
}

func Segment(str string) []string {
	var words []string
	segments := seg.InternalSegment([]byte(str), true)
	for _, v := range segments {
		word := v.Token().Text()
		word = strings.Replace(word, " ", "", -1)
		if word != "" {
			words = append(words, word)
		}
	}
	return words
}
