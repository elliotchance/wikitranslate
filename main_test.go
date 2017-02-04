package main

import (
	"testing"
)

type example struct {
	name    string
	wiki    string
	html    string
	newWiki string
}

var examples = []example{
	// <em> and <strong>
	{"f101", "foo ''bar'' baz", "foo <em>bar</em> baz", ""},
	{"f102", "foo ''bar'' ''baz'' qux", "foo <em>bar</em> <em>baz</em> qux", ""},
	{"f103", "foo '''bar''' baz", "foo <strong>bar</strong> baz", ""},
	{"f104", "foo '''bar''' '''baz''' qux", "foo <strong>bar</strong> <strong>baz</strong> qux", ""},
	{"f105", "foo '''''bar''''' baz", "foo <strong><em>bar</em></strong> baz", ""},
	{"f106", "foo ''bar baz", "foo <em>bar baz</em>", "foo ''bar baz''"},
	{"f107", "foo '''bar baz", "foo <strong>bar baz</strong>", "foo '''bar baz'''"},

	// Links
	{"l101", "foo [[Bar]] baz", `foo <a href="Bar">Bar</a> baz`, ""},
	{"l102", "foo [[Bar|some label]] baz", `foo <a href="Bar">some label</a> baz`, ""},
	{"l103", "foo [[Bar|some label|foo]] baz", `foo <a href="Bar">some label|foo</a> baz`, ""},

	// Images
	{"i101", "foo [[File:filename.extension]] baz", `foo <img src="filename.extension" options="" link=""></img> baz`, ""},
	{"i102", "foo [[File:filename.extension|options]] baz", `foo <img src="filename.extension" options="options" link=""></img> baz`, ""},
	{"i103", "foo [[File:filename.extension|options|caption words]] baz", `foo <img src="filename.extension" options="options" link="">caption words</img> baz`, ""},
	{"i104", "foo [[File:filename.extension|options|link=Internal]] baz", `foo <img src="filename.extension" options="options" link="Internal"></img> baz`, ""},
	{"i105", "foo [[File:filename.extension|options|link=http://External]] baz", `foo <img src="filename.extension" options="options" link="http://External"></img> baz`, ""},

	// References
	{"r101", "foo <ref>[[ABC]]</ref> baz", `foo <ref data="W1tBQkNdXQ=="></ref> baz`, ""},
	{"r102", `foo <ref name="qux">[[ABC]]</ref> baz`, `foo <ref data="W1tBQkNdXQ==" name="qux"></ref> baz`, ""},

	// <nowiki>
	{"w101", "foo <nowiki>''qux''</nowiki> baz", `foo <nowiki data="JydxdXgnJw=="></nowiki> baz`, ""},
	{"w102", "foo <nowiki abc>''qux''</nowiki> baz", `foo <nowiki data="JydxdXgnJw==" abc></nowiki> baz`, ""},

	// Templates
	{"t101", "foo {{bar}} baz", `foo <template name="bar"></template> baz`, ""},
	{"t102",
		"foo {{bar|qux}} baz",
		`foo <template name="bar"><arg name="">qux</arg></template> baz`,
		""},
	{"t103",
		"foo {{bar|qux|abc}} baz",
		`foo <template name="bar"><arg name="">qux</arg><arg name="">abc</arg></template> baz`,
		""},
	{"t104",
		"foo {{bar|qux=abc}} baz",
		`foo <template name="bar"><arg name="qux">abc</arg></template> baz`,
		""},
	{"t105",
		"foo {{bar|\nqux=abc}} baz",
		`foo <template name="bar"><arg name="qux">abc</arg></template> baz`,
		"foo {{bar|qux=abc}} baz"},
	{"t106",
		"foo {{bar| qux =abc}} baz",
		`foo <template name="bar"><arg name="qux">abc</arg></template> baz`,
		"foo {{bar|qux=abc}} baz"},
	{"t107",
		"foo {{bar\n|qux=abc}} baz",
		`foo <template name="bar"><arg name="qux">abc</arg></template> baz`,
		"foo {{bar|qux=abc}} baz"},
	{"t108",
		"foo {{bar\n|qux=[[abc|foo]]}} baz",
		`foo <template name="bar"><arg name="qux"><a href="abc">foo</a></arg></template> baz`,
		"foo {{bar|qux=[[abc|foo]]}} baz"},

	// Nested templates
	{"t201",
		"foo {{bar|{{qux|xyz}}|a=c}} baz",
		`foo <template name="bar"><arg name=""><template name="qux"><arg name="">xyz</arg></template></arg><arg name="a">c</arg></template> baz`,
		""},

	// Headings
	{"h101", "====== The Heading ======\nbar", "<h6> The Heading </h6>\nbar", ""},
	{"h102", "===== The Heading =====\nbar", "<h5> The Heading </h5>\nbar", ""},
	{"h103", "==== The Heading ====\nbar", "<h4> The Heading </h4>\nbar", ""},
	{"h104", "=== The Heading ===\nbar", "<h3> The Heading </h3>\nbar", ""},
	{"h105", "== The Heading ==\nbar", "<h2> The Heading </h2>\nbar", ""},
	{"h106", "= The Heading =\nbar", "<h1> The Heading </h1>\nbar", ""},

	{"h201", " ====== The Heading ======\nbar", " <h6> The Heading </h6>\nbar", ""},
	{"h202", " ===== The Heading =====\nbar", " <h5> The Heading </h5>\nbar", ""},
	{"h203", " ==== The Heading ====\nbar", " <h4> The Heading </h4>\nbar", ""},
	{"h204", " === The Heading ===\nbar", " <h3> The Heading </h3>\nbar", ""},
	{"h205", " == The Heading ==\nbar", " <h2> The Heading </h2>\nbar", ""},
	{"h206", " = The Heading =\nbar", " <h1> The Heading </h1>\nbar", ""},

	{"h301", "foo\n====== The Heading ======\nbar", "foo\n<h6> The Heading </h6>\nbar", ""},
	{"h302", "foo\n===== The Heading =====\nbar", "foo\n<h5> The Heading </h5>\nbar", ""},
	{"h303", "foo\n==== The Heading ====\nbar", "foo\n<h4> The Heading </h4>\nbar", ""},
	{"h304", "foo\n=== The Heading ===\nbar", "foo\n<h3> The Heading </h3>\nbar", ""},
	{"h305", "foo\n== The Heading ==\nbar", "foo\n<h2> The Heading </h2>\nbar", ""},
	{"h306", "foo\n= The Heading =\nbar", "foo\n<h1> The Heading </h1>\nbar", ""},

	// Lists
	{"o101", "Foo\n* Bar\n* Baz\nQux", "Foo\n<li> Bar</li>\n<li> Baz</li>\nQux", ""},
	{"o102", "Foo\n# Bar\n# Baz\nQux", "Foo\n<oli> Bar</oli>\n<oli> Baz</oli>\nQux", ""},
	{"o103", "Foo\n*Bar\n*Baz\nQux", "Foo\n<li>Bar</li>\n<li>Baz</li>\nQux", ""},
	{"o104", "Foo\n#Bar\n#Baz\nQux", "Foo\n<oli>Bar</oli>\n<oli>Baz</oli>\nQux", ""},

	// Tables
	{"g101", "Foo\n{|\n|-\n|Bar\n|}\nQux",
		"Foo\n<table >\n<tr >\n<td >Bar</td>\n</tr>\n</table>\nQux",
		""},
	{"g102", "Foo\n{|\n|-\n|Bar\n|Baz\n|}\nQux",
		"Foo\n<table >\n<tr >\n<td >Bar</td>\n<td >Baz</td>\n</tr>\n</table>\nQux",
		""},
	{"g103", "Foo\n{|\n|-\n|Bar\n|-\n|Baz\n|}\nQux",
		"Foo\n<table >\n<tr >\n<td >Bar</td>\n</tr>\n<tr >\n<td >Baz</td>\n</tr>\n</table>\nQux",
		""},

	{"g201", "Foo\n{|\n|Bar\n|}\nQux",
		"Foo\n<table >\n<tr>\n<td >Bar</td>\n</tr>\n</table>\nQux",
		"Foo\n{|\n|-\n|Bar\n|}\nQux"},
	{"g202", "Foo\n{|\n|Bar\n|Baz\n|}\nQux",
		"Foo\n<table >\n<tr>\n<td >Bar</td>\n<td >Baz</td>\n</tr>\n</table>\nQux",
		"Foo\n{|\n|-\n|Bar\n|Baz\n|}\nQux"},
	{"g203", "Foo\n{|\n|Bar\n|-\n|Baz\n|}\nQux",
		"Foo\n<table >\n<tr>\n<td >Bar</td>\n</tr>\n<tr >\n<td >Baz</td>\n</tr>\n</table>\nQux",
		"Foo\n{|\n|-\n|Bar\n|-\n|Baz\n|}\nQux"},

	{"g301", "Foo\n{|\n!Bar\n|}\nQux",
		"Foo\n<table >\n<tr>\n<th >Bar</th>\n</tr>\n</table>\nQux",
		"Foo\n{|\n|-\n!Bar\n|}\nQux"},
	{"g302", "Foo\n{|\n!Bar\n!Baz\n|}\nQux",
		"Foo\n<table >\n<tr>\n<th >Bar</th>\n<th >Baz</th>\n</tr>\n</table>\nQux",
		"Foo\n{|\n|-\n!Bar\n!Baz\n|}\nQux"},
	{"g303", "Foo\n{|\n!Bar\n|-\n!Baz\n|}\nQux",
		"Foo\n<table >\n<tr>\n<th >Bar</th>\n</tr>\n<tr >\n<th >Baz</th>\n</tr>\n</table>\nQux",
		"Foo\n{|\n|-\n!Bar\n|-\n!Baz\n|}\nQux"},
}

func TestExamples(t *testing.T) {
	for _, test := range examples {
		if test.newWiki == "" {
			test.newWiki = test.wiki
		}

		html := WikiToHtml(test.wiki)
		if html != test.html {
			t.Errorf("%v:\n  expected HTML: '%v'\n      from wiki: '%v'\n            got: '%v'\n\n",
				test.name, test.html, test.wiki, html)
		}

		wiki := HtmlToWiki(test.html)
		if wiki != test.newWiki {
			t.Errorf("%v:\n  expected wiki: '%v'\n      from HTML: '%v'\n            got: '%v'\n\n",
				test.name, test.newWiki, test.html, wiki)
		}
	}
}

type balanceHtmlTagsExample struct {
	before   string
	expected string
}

var balanceHtmlTagsExamples = []balanceHtmlTagsExample{
	// Already balanced
	{"foo bar", "foo bar"},
	{"foo <bar>bar</bar> baz", "foo <bar>bar</bar> baz"},
	{"foo <bar a=1>bar</bar> baz", "foo <bar a=1>bar</bar> baz"},
	{"foo <qux><bar>bar</bar><bar>quxx</bar></qux> baz", "foo <qux><bar>bar</bar><bar>quxx</bar></qux> baz"},

	// Too many opens
	{"foo <bar>bar baz", "foo <bar>bar baz</bar>"},
	{"foo <bar a=1>bar baz", "foo <bar a=1>bar baz</bar>"},
	{"foo <bar>bar <qux>baz", "foo <bar>bar <qux>baz</bar></qux>"},
	{"foo <bar><qux>bar baz", "foo <bar><qux>bar baz</bar></qux>"},
	{"foo <bar><qux>bar</bar> baz", "foo <bar><qux>bar</qux></bar> baz"},
	{"foo <bar><qux>bar<abc></bar> baz", "foo <bar><qux>bar<abc></abc></qux></bar> baz"},

	// To many closes
	{"foo bar</bar> baz", "foo bar baz"},
	{"foo bar</bar> baz</foo>", "foo bar baz"},
	{"foo bar</bar> <abc>baz</foo>", "foo bar <abc>baz</abc>"},

	// "<>"
}

func TestBalanceHtmlTags(t *testing.T) {
	for _, test := range balanceHtmlTagsExamples {
		result := BalanceHtmlTags(test.before)
		if test.expected != result {
			t.Errorf("Expected '%v', got '%v'", test.expected, result)
		}
	}
}
