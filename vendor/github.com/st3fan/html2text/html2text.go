package html2text

import (
	"bytes"
	"regexp"
	"strings"
)

var badTagnamesRE = regexp.MustCompile(`^(head|script|style|a)($|\s*)`)
var linkTagRE = regexp.MustCompile(`a.*href=('([^']*?)'|"([^"]*?)")`)
var badLinkHrefRE = regexp.MustCompile(`#|javascript:`)

func parseHTMLEntity(entName string) (string, bool) {
	entName = strings.ToLower(entName)

	// possible entities
	switch entName {
	case "nbsp":
		return " ", true
	case "gt":
		return ">", true
	case "lt":
		return "<", true
	case "amp":
		return "&", true
	case "quot":
		return "\"", true
	case "apos":
		return "'", true
	case "cent":
		return "¢", true
	case "pound":
		return "£", true
	case "yen":
		return "¥", true
	case "euro":
		return "€", true
	case "copy":
		return "©", true
	case "reg":
		return "®", true
	default:
		return "", false
	}

}

// HTMLEntitiesToText decodes HTML entities inside a provided
// string and returns decoded text
func HTMLEntitiesToText(htmlEntsText string) string {
	outBuf := bytes.NewBufferString("")
	inEnt := false

	for i, r := range htmlEntsText {
		switch {
		case r == ';' && inEnt:
			inEnt = false
			continue

		case r == '&': //possible html entity
			entName := ""
			isEnt := false

			// parse the entity name - max 10 chars
			chars := 0
			for _, er := range htmlEntsText[i+1:] {
				if er == ';' {
					isEnt = true
					break
				} else {
					entName += string(er)
				}

				chars++
				if chars == 10 {
					break
				}
			}

			if isEnt {
				if ent, isEnt := parseHTMLEntity(entName); isEnt {
					outBuf.WriteString(ent)
					inEnt = true
					continue
				}
			}
		}

		if !inEnt {
			outBuf.WriteRune(r)
		}
	}

	return outBuf.String()
}

// HTML2Text converts html into a text form
func HTML2Text(html string) string {
	inLen := len(html)
	tagStart := 0
	inEnt := false
	badTagStackDepth := 0 // if == 1 it means we are inside <head>...</head>
	shouldOutput := true
	// new line cannot be printed at the beginning or
	// for <p> after a new line created by previous <p></p>
	canPrintNewline := false

	outBuf := bytes.NewBufferString("")

	for i, r := range html {
		if inLen > 0 && i == inLen-1 {
			// prevent new line at the end of the document
			canPrintNewline = false
		}

		switch {
		case r <= 0xD, r == 0x85, r == 0x2028, r == 0x2029: // skip new lines
			continue

		case r == ';' && inEnt: // end of html entity
			inEnt = false
			shouldOutput = true
			continue

		case r == '&' && shouldOutput: // possible html entity
			entName := ""
			isEnt := false

			// parse the entity name - max 10 chars
			chars := 0
			for _, er := range html[i+1:] {
				if er == ';' {
					isEnt = true
					break
				} else {
					entName += string(er)
				}

				chars++
				if chars == 10 {
					break
				}
			}

			if isEnt {
				if ent, isEnt := parseHTMLEntity(entName); isEnt {
					outBuf.WriteString(ent)
					inEnt = true
					shouldOutput = false
					continue
				}
			}

		case r == '<': // start of a tag
			tagStart = i + 1
			shouldOutput = false
			continue

		case r == '>': // end of a tag
			shouldOutput = true
			tagName := strings.ToLower(html[tagStart:i])

			if tagName == "br" || tagName == "br/" {
				// new line
				outBuf.WriteString("\r\n")
			} else if tagName == "p" || tagName == "/p" {
				if canPrintNewline {
					outBuf.WriteString("\r\n")
				}
				canPrintNewline = false
			} else if badTagnamesRE.MatchString(tagName) {
				// unwanted block
				badTagStackDepth++

				// parse link href
				m := linkTagRE.FindStringSubmatch(tagName)
				if len(m) == 4 {
					link := m[2]
					if len(link) == 0 {
						link = m[3]
					}

					if !badLinkHrefRE.MatchString(link) {
						outBuf.WriteString(HTMLEntitiesToText(link))
					}
				}
			} else if len(tagName) > 0 && tagName[0] == '/' &&
				badTagnamesRE.MatchString(tagName[1:]) {
				// end of unwanted block
				badTagStackDepth--
			}
			continue

		} // switch end

		if shouldOutput && badTagStackDepth == 0 && !inEnt {
			canPrintNewline = true
			outBuf.WriteRune(r)
		}
	}

	return outBuf.String()
}
