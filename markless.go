package main

import (
	"regexp"
	"bytes"
	"bufio"
	"strings"
)

// headings: 1-6 # to equivalent h1-h6 tags
var h6Re = regexp.MustCompile(`^###### ([A-z0-9].*)$`)
var h5Re = regexp.MustCompile(`^##### (.*)$`)
var h4Re = regexp.MustCompile(`^#### (.*)$`)
var h3Re = regexp.MustCompile(`^### (.*)$`)
var h2Re = regexp.MustCompile(`^## (.*)$`)
var h1Re = regexp.MustCompile(`^# (.*)$`)
var mailRe = regexp.MustCompile(`<(.*?)@(.*?)>`)
var urlRe = regexp.MustCompile(`<(.*?)>`)
var mdashRe = regexp.MustCompile(`---`)
var ndashRe = regexp.MustCompile(`--`)
var threeAsteriskRe = regexp.MustCompile(`\*\*\*(.*?)\*\*\*`)
var twoAsteriskRe = regexp.MustCompile(`\*\*(.*?)\*\*`)
var oneAsteriskRe = regexp.MustCompile(`\*(.*?)\*`)
var brRe = regexp.MustCompile(`  $`)
var citeRe = regexp.MustCompile(`^>\[(.*?)\] (.*)$`)
var quoteRe = regexp.MustCompile(`^> (.*)$`)
var codeRe = regexp.MustCompile(`^(	|    )(.*)$`)
var footnoteRe = regexp.MustCompile(`\[(.*?)\]\[(.*?)\]`)
var ulRe = regexp.MustCompile(`^(-|\*|.) (.*)$`)
var olRe = regexp.MustCompile(`^1. (.*)$`)
var hrRe = regexp.MustCompile(`^___$`)
var subRe = regexp.MustCompile(`\^\^(.*?)\^\^`)
var supRe = regexp.MustCompile(`\^(.*?)\^`)
var imgRe = regexp.MustCompile(`!\[(.*?)\]\((.*?) "(.*?)"\)`)
var img2Re = regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
var linkRe = regexp.MustCompile(`\[(.*?)\]\((.*?) "(.*?)"\)`)
var link2Re = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
var wikiRe = regexp.MustCompile(`\[(.*?)\]`)

func markupHeadings(line *[]byte) {
	*line = h6Re.ReplaceAll(*line, []byte("<h6>$1</h6>"))
	*line = h5Re.ReplaceAll(*line, []byte("<h5>$1</h5>"))
	*line = h4Re.ReplaceAll(*line, []byte("<h4>$1</h4>"))
	*line = h3Re.ReplaceAll(*line, []byte("<h3>$1</h3>"))
	*line = h2Re.ReplaceAll(*line, []byte("<h2>$1</h2>"))
	*line = h1Re.ReplaceAll(*line, []byte("<h1>$1</h1>"))
}

func markupMail(line *[]byte) {
	regex := "[$1@$2](mailto:$1@$2 \"Email $1 at $2\")"
	*line = mailRe.ReplaceAll(*line, []byte(regex))
}

func markupUrl(line *[]byte) {
	regex := "<a href=\"$1\">$1</a>"
	*line = urlRe.ReplaceAll(*line, []byte(regex))
}

func markupDashes(line *[]byte) {
	// TeX-style dashes
	*line = mdashRe.ReplaceAll(*line, []byte("&mdash;"))
	*line = ndashRe.ReplaceAll(*line, []byte("&ndash;"))
}

func markupAsterisks(line *[]byte) {
	// triple asterisk: bold italic
	// double asterisk: bold
	// single asterisk: italic
	r_emstrong := "<em><strong>$1</strong></em>"
	r_strong := "<strong>$1</strong>"
	r_em := "<em>$1</em>"
	*line = threeAsteriskRe.ReplaceAll(*line, []byte(r_emstrong))
	*line = twoAsteriskRe.ReplaceAll(*line, []byte(r_strong))
	*line = oneAsteriskRe.ReplaceAll(*line, []byte(r_em))
}

func markupCite(line *[]byte) {
	citeText := `<figure>
  <blockquote>
    <p>"$2"</p>
  </blockquote>
  <figcaption>
    &mdash; $1
  </figcaption>
</figure>`
	*line = citeRe.ReplaceAll(*line, []byte(citeText))
}

func markupQuote(line *[]byte) {
	quoteText := `<figure>
  <blockquote>
    <p>"$1"</p>
  </blockquote>
</figure>`
	*line = quoteRe.ReplaceAll(*line, []byte(quoteText))
}

func markupBr(line *[]byte) {
	*line = brRe.ReplaceAll(*line, []byte("<br/>"))
}

func markupCode(line *[]byte) {
	*line = codeRe.ReplaceAll(*line, []byte("<tt>$2</tt>"))
}

func markupFootnote(line *[]byte, addFootnotes *bytes.Buffer) {
	for footnoteRe.Match(*line) {
		sline := string(*line)
		m := footnoteRe.FindStringSubmatch(sline)
		name := string(m[1])
		s := footnoteRe.Split(sline, 2)
		rest := string(s[1])
		(*addFootnotes).Write([]byte("<ol>\n<li><a name=\"" + m[1] + "\">" + m[2] + "</a></li>\n</ol>\n"))
		sline = s[0] + name + "^<a href=\"#" + name + "\">&#91;v&#93;</a>^" + rest
		*line = []byte(sline)
	}
}

func markupLists(line *[]byte) {
	*line = ulRe.ReplaceAll(*line, []byte("<ul>\n<li>$2</li>\n</ul>"))
	*line = olRe.ReplaceAll(*line, []byte("<ol>\n<li>$1</li>\n</ol>"))
}

func markupRule(line *[]byte) {
	*line = hrRe.ReplaceAll(*line, []byte("<hr/>"))
}

func markupSub(line *[]byte) {
	*line = subRe.ReplaceAll(*line, []byte("<sub>$1</sub>"))
}

func markupSup(line *[]byte) {
	*line = supRe.ReplaceAll(*line, []byte("<sup>$1</sup>"))
}

func markupImages(line *[]byte) {
	imgText := `<figure>
  <img src="$2" alt="&#91;IMAGE: $3&#93;" />
  <figcaption>$1</figcaption>
</figure>`
	*line = imgRe.ReplaceAll(*line, []byte(imgText))
	img2Text := `<figure>
  <img src="$2" alt="&#91;IMAGE: No description provided&#93;" />
  <figcaption>$1</figcaption>
</figure>`
	*line = img2Re.ReplaceAll(*line, []byte(img2Text))
}

func markupLinks(line *[]byte) {
	*line = linkRe.ReplaceAll(*line, []byte("<a href=\"$2\" title=\"$3\">$1</a>"))
	*line = link2Re.ReplaceAll(*line, []byte("<a href=\"$2\">$1</a>"))
}

func markupWiki(line *[]byte) {
	*line = wikiRe.ReplaceAll(*line, []byte("<a href=\"/view/$1\">$1</a>"))
}

func markup(line *[]byte, addFootnotes *bytes.Buffer) {
	markupMail(line)
	markupUrl(line)
	markupHeadings(line)
	markupDashes(line)
	markupAsterisks(line)
	markupBr(line)
	markupCite(line)
	markupQuote(line)
	markupCode(line)
	markupFootnote(line, addFootnotes)
	markupLists(line)
	markupRule(line)
	markupSub(line)
	markupSup(line)
	markupImages(line)
	markupLinks(line)
	markupWiki(line)
}

// I'm not doing nesting of anything. If you try nesting and it works,
// great. If not, you get to keep all the pieces.
func markupOutput(o []byte) ([]byte, error) {
	var err error
	// Break into lines and iterate
	var os bytes.Buffer
	var addFootnotes bytes.Buffer
	// This looks complicated but handles newlines correctly

	scanner := bufio.NewScanner(strings.NewReader(string(o)))
	for scanner.Scan() {
		err = scanner.Err()
		if err != nil {
			return []byte(""), err
		}
		line := []byte(scanner.Text())
		markup(&line, &addFootnotes)
		os.Write(line)
		os.Write([]byte("\n"))
	}
	if len(addFootnotes.Bytes()) > 0 {
		os.Write([]byte("Footnotes:<br/>\n"))
		os.Write(addFootnotes.Bytes())
		os.Write([]byte("<hr/>\n"))
	}
	// smooth over lists
	smoothOs := os.Bytes()
	var olSmoothRe = regexp.MustCompile(`</ol>\n<ol>`)
	smoothOs = olSmoothRe.ReplaceAll(smoothOs, []byte(""))

	var ulSmoothRe = regexp.MustCompile(`</ul>\n<ul>`)
	smoothOs = ulSmoothRe.ReplaceAll(smoothOs, []byte(""))
	return smoothOs, err
}
