package webgen

// NewHTML creates a new blank HTML element with the specified language.
func NewHTML(lang string) *Node {
	if lang == "" {
		return HTML()
	}
	return HTML(Attrs{"lang": lang})
}

// NewHTMLEnglish creates a new English HTML element.
func NewHTMLEnglish() *Node { return NewHTML("en") }

// NewHTMLChinese creates a new Chinese HTML element.
func NewHTMLChinese() *Node { return NewHTML("zh") }
