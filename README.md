[![Build Status](https://travis-ci.org/elliotchance/wikitranslate.svg?branch=master)](https://travis-ci.org/elliotchance/wikitranslate)

`wikitranslate` is a tool specificially designed for aiding the translation of
Wikipedia pages using
[CAT tools](http://www.translationzone.com/solutions/cat-tools/).

Wikipedia pages used a type of markup called
[Wiki markup or wikitext](https://en.wikipedia.org/wiki/Wiki_markup). It is a
pure text based formatting that is similar to
[Markdown](https://en.wikipedia.org/wiki/Markdown) or
[reStructuredText](https://en.wikipedia.org/wiki/ReStructuredText).

This wiki markup (since it's made up of punctuation) makes it very difficult for
CAT tools to understand the difference between the segment text (that needs to
be translated) and formatting. `wikitranslate` converts the wiki markup into
pseudo-HTML that can be ingested and translated. The result HTML can then be
converted back into wiki markup to be uploaded as a new page.

Installation and Updating
=========================

If you have not yet downloaded `wikitransate` you can download the latest binary
from the
[releases page](https://github.com/elliotchance/wikitranslate/releases).

**Note:** You will not be able to open the file downloaded. It will only work
inside your
[Terminal](http://www.howtogeek.com/210147/how-to-open-terminal-in-the-current-os-x-finder-location/).

Once installed you can use the built-in updater:

```bash
$ wiktranslate update
The current version is v0.2.0
Finding the latest version... v0.3.1
Downloading the latest version... Done (5.43 MB)
Installing... Done
```

Usage
=====

To prepare a wiki page for translating, provide the URL:

```bash
wikitranslate https://en.wikipedia.org/wiki/Staffordshire_Bull_Terrier
```

This will generate a `Staffordshire_Bull_Terrier.html` in your Downloads folder.
This is the document to upload or import into you CAT tools.

---

Once the translation is complete you will need to download or export the new
HTML document and use `wikitranslate` to convert it back into the wiki markup by
providing the new HTML file:

```bash
wikitranslate Staffordshire_Bull_Terrier.html
```

**Tip:** You can drag the file into the Terminal to insert the full path to the
HTML document.

This will generate a `Staffordshire_Bull_Terrier.html.txt` in the same folder as
`Staffordshire_Bull_Terrier.html`. You can now open the text file to get the
wiki markup for submission.

Considerations for the Intermediate Markup
==========================================

1. **The HTML file should never be manually edited.** Especially in earlier
versions of `wikitranslate` where the internals expect certain attributes and
may not be able to translate back to wiki makrup accurately, or at all if it has
been edit incorrectly.

2. `wikitranslate` maintains the complete life-cycle of the document and is
intended to only work with wiki markup. It is not made for processing HTML from
other sources.

3. The HTML generated is intended for CAT tools and not viewing directly. While
regular HTML elements are used for most formatting it also uses custom tags and
other markers that make the processing back to wikimarkup possible but also hide
(from view) some of the page elements.

4. The content of references (`<ref>`) and unformatted blocks (`<nowiki>`) are
concealed, they will not appear in your translation but will be returned exactly
as they were and in the same place in the new wikimarkup.

5. The layout of formatting will not be maintained. A good example of this is
[https://en.wikipedia.org/wiki/Help:Table](tables) that use the short-hand `!!`
for adding multiple columns to the same line. This will always be expanded in
the output to use one line per column, however this may change in the future.
