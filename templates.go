package main

import (
	"net/http"
	"text/template"
	"regexp"
	"bytes"
	"fmt"
	"bufio"
	"strings"
)

// https://go.dev/doc/articles/wiki/

var templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html"))
var statics = template.Must(template.ParseFiles("templates/header.html", "templates/footer.html"))
var staticHead bytes.Buffer
var staticFoot bytes.Buffer

func initStaticTemplates() error {
	//fmt.Println("Loading static templates.")
	err := statics.ExecuteTemplate(&staticHead, "header.html", nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	//fmt.Println("HEAD: ", staticHead.String())
	err = statics.ExecuteTemplate(&staticFoot, "footer.html", nil)
	return err
}

// TODO: break up this function.
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

		// headings: 1-6 # to equivalent h1-h6 tags
		var h6Re = regexp.MustCompile(`^###### ([A-z0-9].*)$`)
		var h5Re = regexp.MustCompile(`^##### (.*)$`)
		var h4Re = regexp.MustCompile(`^#### (.*)$`)
		var h3Re = regexp.MustCompile(`^### (.*)$`)
		var h2Re = regexp.MustCompile(`^## (.*)$`)
		var h1Re = regexp.MustCompile(`^# (.*)$`)

		var mailRe = regexp.MustCompile(`<(.*)@(.*)>`)
		line = mailRe.ReplaceAll(line, []byte("[$1@$2](mailto:$1@$2 \"Email $1 at $2\")"))

		var urlRe = regexp.MustCompile(`<(.*)>`)
		line = urlRe.ReplaceAll(line, []byte("<a href=\"$1\">$1</a>"))

		line = h6Re.ReplaceAll(line, []byte("<h6>$1</h6>"))
		line = h5Re.ReplaceAll(line, []byte("<h5>$1</h5>"))
		line = h4Re.ReplaceAll(line, []byte("<h4>$1</h4>"))
		line = h3Re.ReplaceAll(line, []byte("<h3>$1</h3>"))
		line = h2Re.ReplaceAll(line, []byte("<h2>$1</h2>"))
		line = h1Re.ReplaceAll(line, []byte("<h1>$1</h1>"))

		// TeX-style dashes
		var mdashRe = regexp.MustCompile(`---`)
		line = mdashRe.ReplaceAll(line, []byte("&mdash;"))
		var ndashRe = regexp.MustCompile(`--`)
		line = ndashRe.ReplaceAll(line, []byte("&ndash;"))

		// triple asterisk: bold italic
		// double asterisk: bold
		// single asterisk: italic
		var threeAsteriskRe = regexp.MustCompile(`\*\*\*(.*)\*\*\*`)
		line = threeAsteriskRe.ReplaceAll(line, []byte("<em><strong>$1</strong></em>"))
		var twoAsteriskRe = regexp.MustCompile(`\*\*(.*)\*\*`)
		line = twoAsteriskRe.ReplaceAll(line, []byte("<strong>$1</strong>"))
		var oneAsteriskRe = regexp.MustCompile(`\*(.*)\*`)
		line = oneAsteriskRe.ReplaceAll(line, []byte("<em>$1</em>"))

		var brRe = regexp.MustCompile(`  $`)
		line = brRe.ReplaceAll(line, []byte("<br/>"))

		var citeRe = regexp.MustCompile(`^>\[(.*)\] (.*)$`)
		citeText := `<figure>
  <blockquote>
    <p>"$2"</p>
  </blockquote>
  <figcaption>
    &mdash; $1
  </figcaption>
</figure>`
		line = citeRe.ReplaceAll(line, []byte(citeText))

		var quoteRe = regexp.MustCompile(`^> (.*)$`)
		quoteText := `<figure>
  <blockquote>
    <p>"$1"</p>
  </blockquote>
</figure>`
		line = quoteRe.ReplaceAll(line, []byte(quoteText))

		var codeRe = regexp.MustCompile(`^(	|    )(.*)$`)
		line = codeRe.ReplaceAll(line, []byte("<tt>$2</tt>"))

		var footnoteRe = regexp.MustCompile(`\[(.*?)\]\[(.*?)\]`)
		for footnoteRe.Match(line) {
			sline := string(line)
			m := footnoteRe.FindStringSubmatch(sline)
			name := string(m[1])
			s := footnoteRe.Split(sline, 2)
			rest := string(s[1])
			addFootnotes.Write([]byte("<ol>\n<li><a name=\"" + m[1] + "\">" + m[2] + "</a></li>\n</ol>\n"))
			sline = s[0] + name + "^<a href=\"#" + name + "\">&#91;v&#93;</a>^" + rest
			line = []byte(sline)
		}

		var ulRe = regexp.MustCompile(`^(-|\*|.) (.*)$`)
		line = ulRe.ReplaceAll(line, []byte("<ul>\n<li>$2</li>\n</ul>"))

		var olRe = regexp.MustCompile(`^1. (.*)$`)
		line = olRe.ReplaceAll(line, []byte("<ol>\n<li>$1</li>\n</ol>"))

		var hrRe = regexp.MustCompile(`^___$`)
		line = hrRe.ReplaceAll(line, []byte("<hr/>"))

		var subRe = regexp.MustCompile(`\^\^(.*?)\^\^`)
		line = subRe.ReplaceAll(line, []byte("<sub>$1</sub>"))

		var supRe = regexp.MustCompile(`\^(.*?)\^`)
		line = supRe.ReplaceAll(line, []byte("<sup>$1</sup>"))

		var imgRe = regexp.MustCompile(`!\[(.*)\]\((.*) "(.*)"\)`)
		imgText := `<figure>
  <img src="$2" alt="&#91;IMAGE: $3&#93;" />
  <figcaption>$1</figcaption>
</figure>`
		line = imgRe.ReplaceAll(line, []byte(imgText))

		var img2Re = regexp.MustCompile(`!\[(.*)\]\((.*)\)`)
		img2Text := `<figure>
  <img src="$2" alt="&#91;IMAGE: No description provided&#93;" />
  <figcaption>$1</figcaption>
</figure>`
		line = img2Re.ReplaceAll(line, []byte(img2Text))

		var linkRe = regexp.MustCompile(`\[(.*)\]\((.*) "(.*)"\)`)
		line = linkRe.ReplaceAll(line, []byte("<a href=\"$2\" title=\"$3\">$1</a>"))

		var link2Re = regexp.MustCompile(`\[(.*)\]\((.*)\)`)
		line = link2Re.ReplaceAll(line, []byte("<a href=\"$2\">$1</a>"))

		var wikiRe = regexp.MustCompile(`\[(.*)\]`)
		line = wikiRe.ReplaceAll(line, []byte("<a href=\"/view/$1\">$1</a>"))

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

	// add a header and footer which are not marked up
	// so we can have actual HTML tags in there
	var capped bytes.Buffer
	capped.Write(staticHead.Bytes())
	capped.Write(smoothOs)
	capped.Write(staticFoot.Bytes())
	return capped.Bytes(), err
}

func captureTemplate(tmpl string, p *Page) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := templates.ExecuteTemplate(buf, tmpl + ".html", p)
	return buf.Bytes(), err
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	output, err := captureTemplate(tmpl, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if tmpl == "view" {
		output, err = markupOutput(output)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(output)
}
