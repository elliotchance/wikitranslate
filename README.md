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
