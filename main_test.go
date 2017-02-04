package main

import (
    "testing"
)

type example struct {
    wiki string
    html string
    newWiki string
}

var examples = []example{
  // <em> and <strong>
  { "foo ''bar'' baz", "foo <em>bar</em> baz", "" },
  { "foo ''bar'' ''baz'' qux", "foo <em>bar</em> <em>baz</em> qux", "" },
  { "foo '''bar''' baz", "foo <strong>bar</strong> baz", "" },
  { "foo '''bar''' '''baz''' qux", "foo <strong>bar</strong> <strong>baz</strong> qux", "" },
  { "foo '''''bar''''' baz", "foo <strong><em>bar</em></strong> baz", "" },
  { "foo ''bar baz", "foo <em>bar baz</em>", "foo ''bar baz''" },
  { "foo '''bar baz", "foo <strong>bar baz</strong>", "foo '''bar baz'''" },

  // Links
  { "foo [[Bar]] baz", `foo <a href="Bar">Bar</a> baz`, "" },
  { "foo [[Bar|some label]] baz", `foo <a href="Bar">some label</a> baz`, "" },
  { "foo [[Bar|some label|foo]] baz", `foo <a href="Bar">some label|foo</a> baz`, "" },

  // Images
  { "foo [[File:filename.extension]] baz", `foo <img src="filename.extension" options="" link=""></img> baz`, "" },
  { "foo [[File:filename.extension|options]] baz", `foo <img src="filename.extension" options="options" link=""></img> baz`, "" },
  { "foo [[File:filename.extension|options|caption words]] baz", `foo <img src="filename.extension" options="options" link="">caption words</img> baz`, "" },
  { "foo [[File:filename.extension|options|link=Internal]] baz", `foo <img src="filename.extension" options="options" link="Internal"></img> baz`, "" },
  { "foo [[File:filename.extension|options|link=http://External]] baz", `foo <img src="filename.extension" options="options" link="http://External"></img> baz`, "" },
  
  // References
  { "foo <ref>[[ABC]]</ref> baz", `foo <ref data="W1tBQkNdXQ=="></ref> baz`, "" },
  { `foo <ref name="qux">[[ABC]]</ref> baz`, `foo <ref data="W1tBQkNdXQ==" name="qux"></ref> baz`, "" },
  
  // <nowiki>
  { "foo <nowiki>''qux''</nowiki> baz", `foo <nowiki data="JydxdXgnJw=="></nowiki> baz`, "" },
  { "foo <nowiki abc>''qux''</nowiki> baz", `foo <nowiki data="JydxdXgnJw==" abc></nowiki> baz`, "" },

  // Templates
  { "foo {{bar}} baz", `foo <template name="bar"></template> baz`, "" },
  { "foo {{bar|qux}} baz", `foo <template name="bar"><arg name="">qux</arg></template> baz`, "" },
  { "foo {{bar|qux|abc}} baz", `foo <template name="bar"><arg name="">qux</arg><arg name="">abc</arg></template> baz`, "" },
  { "foo {{bar|qux=abc}} baz", `foo <template name="bar"><arg name="qux">abc</arg></template> baz`, "" },

  // Headings
  { "====== The Heading ======\nbar", "<h6> The Heading </h6>\nbar", "" },
  { "===== The Heading =====\nbar", "<h5> The Heading </h5>\nbar", "" },
  { "==== The Heading ====\nbar", "<h4> The Heading </h4>\nbar", "" },
  { "=== The Heading ===\nbar", "<h3> The Heading </h3>\nbar", "" },
  { "== The Heading ==\nbar", "<h2> The Heading </h2>\nbar", "" },
  { "= The Heading =\nbar", "<h1> The Heading </h1>\nbar", "" },
  { " ====== The Heading ======\nbar", " <h6> The Heading </h6>\nbar", "" },
  { " ===== The Heading =====\nbar", " <h5> The Heading </h5>\nbar", "" },
  { " ==== The Heading ====\nbar", " <h4> The Heading </h4>\nbar", "" },
  { " === The Heading ===\nbar", " <h3> The Heading </h3>\nbar", "" },
  { " == The Heading ==\nbar", " <h2> The Heading </h2>\nbar", "" },
  { " = The Heading =\nbar", " <h1> The Heading </h1>\nbar", "" },
  { "foo\n====== The Heading ======\nbar", "foo\n<h6> The Heading </h6>\nbar", "" },
  { "foo\n===== The Heading =====\nbar", "foo\n<h5> The Heading </h5>\nbar", "" },
  { "foo\n==== The Heading ====\nbar", "foo\n<h4> The Heading </h4>\nbar", "" },
  { "foo\n=== The Heading ===\nbar", "foo\n<h3> The Heading </h3>\nbar", "" },
  { "foo\n== The Heading ==\nbar", "foo\n<h2> The Heading </h2>\nbar", "" },
  { "foo\n= The Heading =\nbar", "foo\n<h1> The Heading </h1>\nbar", "" },

  // Lists
  { "Foo\n* Bar\n* Baz\nQux", "Foo\n<li> Bar</li>\n<li> Baz</li>\nQux", "" },
  { "Foo\n# Bar\n# Baz\nQux", "Foo\n<oli> Bar</oli>\n<oli> Baz</oli>\nQux", "" },
  { "Foo\n*Bar\n*Baz\nQux", "Foo\n<li>Bar</li>\n<li>Baz</li>\nQux", "" },
  { "Foo\n#Bar\n#Baz\nQux", "Foo\n<oli>Bar</oli>\n<oli>Baz</oli>\nQux", "" },
  
  // Tables
  { "Foo\n{|\n|-\n|Bar\n|}\nQux",
    "Foo\n<table >\n<tr>\n<td >Bar</td>\n</tr>\n</table>\nQux",
    "" },
  { "Foo\n{|\n|-\n|Bar\n|Baz\n|}\nQux",
    "Foo\n<table >\n<tr>\n<td >Bar</td>\n<td >Baz</td>\n</tr>\n</table>\nQux",
    "" },
  { "Foo\n{|\n|-\n|Bar\n|-\n|Baz\n|}\nQux",
    "Foo\n<table >\n<tr>\n<td >Bar</td>\n</tr>\n<tr>\n<td >Baz</td>\n</tr>\n</table>\nQux",
    "" },

  { "Foo\n{|\n|Bar\n|}\nQux",
    "Foo\n<table >\n<tr>\n<td >Bar</td>\n</tr>\n</table>\nQux",
    "Foo\n{|\n|-\n|Bar\n|}\nQux" },
  { "Foo\n{|\n|Bar\n|Baz\n|}\nQux",
    "Foo\n<table >\n<tr>\n<td >Bar</td>\n<td >Baz</td>\n</tr>\n</table>\nQux",
    "Foo\n{|\n|-\n|Bar\n|Baz\n|}\nQux" },
  { "Foo\n{|\n|Bar\n|-\n|Baz\n|}\nQux",
    "Foo\n<table >\n<tr>\n<td >Bar</td>\n</tr>\n<tr>\n<td >Baz</td>\n</tr>\n</table>\nQux",
    "Foo\n{|\n|-\n|Bar\n|-\n|Baz\n|}\nQux" },

  { "Foo\n{|\n!Bar\n|}\nQux",
    "Foo\n<table >\n<tr>\n<th >Bar</th>\n</tr>\n</table>\nQux",
    "Foo\n{|\n|-\n!Bar\n|}\nQux" },
  { "Foo\n{|\n!Bar\n!Baz\n|}\nQux",
    "Foo\n<table >\n<tr>\n<th >Bar</th>\n<th >Baz</th>\n</tr>\n</table>\nQux",
    "Foo\n{|\n|-\n!Bar\n!Baz\n|}\nQux" },
  { "Foo\n{|\n!Bar\n|-\n!Baz\n|}\nQux",
    "Foo\n<table >\n<tr>\n<th >Bar</th>\n</tr>\n<tr>\n<th >Baz</th>\n</tr>\n</table>\nQux",
    "Foo\n{|\n|-\n!Bar\n|-\n!Baz\n|}\nQux" },
}

func TestExamples(t *testing.T) {
    for _, test := range examples {
        if test.newWiki == "" {
            test.newWiki = test.wiki
        }

        html := WikiToHtml(test.wiki)
        if html != test.html {
            t.Errorf("Expected HTML '%v' from wiki '%v', got '%v'", test.html, test.wiki, html)
        }

        wiki := HtmlToWiki(test.html)
        if wiki != test.newWiki {
            t.Errorf("Expected Wiki '%v' from HTML '%v', got '%v'", test.newWiki, test.html, wiki)
        }
    }
}

type balanceHtmlTagsExample struct {
    before string
    expected string
}

var balanceHtmlTagsExamples = []balanceHtmlTagsExample{
  // Already balanced
  { "foo bar", "foo bar" },
  { "foo <bar>bar</bar> baz", "foo <bar>bar</bar> baz" },
  { "foo <bar a=1>bar</bar> baz", "foo <bar a=1>bar</bar> baz" },
  { "foo <qux><bar>bar</bar><bar>quxx</bar></qux> baz", "foo <qux><bar>bar</bar><bar>quxx</bar></qux> baz" },

  // Too many opens
  { "foo <bar>bar baz", "foo <bar>bar baz</bar>" },
  { "foo <bar a=1>bar baz", "foo <bar a=1>bar baz</bar>" },
  { "foo <bar>bar <qux>baz", "foo <bar>bar <qux>baz</bar></qux>" },
  { "foo <bar><qux>bar baz", "foo <bar><qux>bar baz</bar></qux>" },
  { "foo <bar><qux>bar</bar> baz", "foo <bar><qux>bar</qux></bar> baz" },
  { "foo <bar><qux>bar<abc></bar> baz", "foo <bar><qux>bar<abc></abc></qux></bar> baz" },

  // To many closes
  { "foo bar</bar> baz", "foo bar baz" },
  { "foo bar</bar> baz</foo>", "foo bar baz" },
  { "foo bar</bar> <abc>baz</foo>", "foo bar <abc>baz</abc>" },

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
