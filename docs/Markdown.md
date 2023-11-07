---
title: Markdown Tutorial
layout: default
parent: Howto document
nav_order: 3
---

# Comprehensive Markdown Tutorial

## Headers
Create headers by prefacing the text with one or more `#` symbols. The number of `#` you use determines the size of the header.

```markdown
# Header 1
## Header 2
### Header 3
#### Header 4
##### Header 5
###### Header 6
```

## Emphasis
Emphasize text with bold or italics.

```markdown
*italicized text*
**bold text**
***bold and italicized text***
```

## Lists
Create ordered and unordered lists.

```markdown
- Unordered list item 1
- Unordered list item 2
  - Nested unordered list item

1. Ordered list item 1
2. Ordered list item 2
   1. Nested ordered list item
```

- Unordered list item 1
- Unordered list item 2
  - Nested unordered list item

1. Ordered list item 1
2. Ordered list item 2
   1. Nested ordered list item

## Links
Include hyperlinks with text.

```markdown
[GitHub](http://github.com)
```

[GitHub](http://github.com)

## Images
Embed images using the following syntax.

```markdown
![alt text for the image](image-url.jpg)
```

![alt text for the image](https://de.wikipedia.org/wiki/Siemens#/media/Datei:Siemens-logo.svg)

## Code
Add inline code with single backticks, and code blocks with triple backticks.

```markdown
`inline code`

```

`for i < 5`

# code block
print('Hello, world!')
```
```

## Tables
Organize data into tables.

```markdown
| Header 1    | Header 2    |
| ----------- | ----------- |
| Row 1 Col 1 | Row 1 Col 2 |
| Row 2 Col 1 | Row 2 Col 2 |
```

| Header 1    | Header 2    |
| ----------- | ----------- |
| Row 1 Col 1 | Row 1 Col 2 |
| Row 2 Col 1 | Row 2 Col 2 |

## Blockquotes
Use blockquotes to quote text.

```markdown
> This is a blockquote.
```

> This is a blockquote.

## Horizontal Rules
Create a horizontal line or page break.

```markdown
---
```

---

## Extended Syntax

### Strikethrough
```markdown
~~Strikethrough text~~
```

~~Strikethrough text~~

### Fenced Code Blocks
```markdown
```
{
  "firstName": "John",
  "lastName": "Smith",
  "age": 25
}
```
```

```
{
  "firstName": "John",
  "lastName": "Smith",
  "age": 25
}
```

### Footnotes
Create a footnote like this.[^1]

[^1]: This is the footnote.

### Heading IDs
```markdown
### My Great Heading {#custom-id}
```

### Definition Lists
```markdown
term
: definition

term2
: definition2
```

### Task Lists
```markdown
- [x] Write the press release
- [ ] Update the website
- [ ] Contact the media
```

- [x] Write the press release
- [ ] Update the website
- [ ] Contact the media

## More

For more advanced Markdown syntax, check the [official documentation](https://daringfireball.net/projects/markdown/).

## Conclusion

Now you can use Markdown to format your documentation content. Remember to preview your changes locally before pushing them to your GitHub repository.
