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

Install and Usage
=================

Download the latest binary from the
[releases page](https://github.com/elliotchance/wikitranslate/releases).

*Note:* You will not be able to open the file downloaded. You will need to open
the Terminal and use the following commands:

To prepare a wiki page for translating, provide the URL:

```bash
wikitranslate https://en.wikipedia.org/wiki/Staffordshire_Bull_Terrier
```

This will generate a `Staffordshire_Bull_Terrier.html` in your Downloads folder.
This is the documentation you upload or import to you CAT tools.

---

Once the translation is complete you will need to download or export the new
HTML document and use `wikitranslate` to convert it back into the wiki markup:

```bash
wikitranslate Staffordshire_Bull_Terrier.html
```

*Tip:* You can drag the file into the Terminal to insert the full path to the
HTML document.

This will generate a `Staffordshire_Bull_Terrier.html.txt` in the same folder as
`Staffordshire_Bull_Terrier.html`. You can now open the text file to get the
wiki markup for submission.
