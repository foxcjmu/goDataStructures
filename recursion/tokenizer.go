// tokenizer.go: It is convenient to have a uniform mechanism for reading the bytes in a
// string one by one. The strings package provides a Reader for this, but it convenient to
// have an even more abstract view of things. The Tokenizer type packages up a string
// reader and the current byte in the string along with methods to advance or back-up
// one byte.

package recursion

import "strings"

type Tokenizer struct {
	reader *strings.Reader // source for reading chars
	Char   byte            // the current char in string; '$' if no more
}

// Create a new Tokenizer: the char field will contain the first byte in
// string s, or $ if s is empty.
func NewTokenizer(s string) *Tokenizer {
	result := new(Tokenizer)
	result.reader = strings.NewReader(s)
	result.Next()
	return result
}

// Next advances to the next byte in the string and puts it in t.Char.
// If the string is exhausted, then t.Char == '$'
func (t *Tokenizer) Next() {
	if t.reader.Len() == 0 {
		t.Char = '$'
	} else {
		t.Char, _ = t.reader.ReadByte()
	}
}

// Last backs-up to the previous byte in the string and puts it in t.Char.
// Pre: at least two characters have been read
// Pre violation: panic
// Normal return: t.Char is set to the previous character read
func (t *Tokenizer) Last() {
	if err := t.reader.UnreadByte(); err != nil {
		panic(err)
	}
	if err := t.reader.UnreadByte(); err != nil {
		panic(err)
	}
	t.Char, _ = t.reader.ReadByte()
}
