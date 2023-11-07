---
title: How to document
layout: default
has_children: true
nav_order: 2
---

# Adding New Pages and Documentation

To add new pages and create documentation with the Just the Docs theme on GitHub Pages, follow these steps:

## Create Markdown Files

For each new page or section of documentation, create a Markdown (.md) file in your repository. The file name will typically correspond to the URL path (e.g., `installation.md` for an Installation page).

## Front Matter

At the top of each Markdown file, include YAML front matter to specify layout and title. Optionally, you can add other front matter as needed.

```yaml
---
layout: default
title: Page Title
nav_order: 1
---
```

The `nav_order` determines the order of the page in the navigation.

## Nested Navigation
To define a nested navigation structure using the front matter in your Markdown files for Just the Docs, you'll need to specify hierarchy keywords like parent, grand_parent, etc., to establish relationships between pages. Here's how you can do it:

In your main page (e.g., docs.md), you might have:

```yaml
---
layout: default
title: Docs
nav_order: 1
---
```
For a child page (e.g., installation.md), you'll reference its parent:

```yaml
---
layout: default
title: Installation
parent: Docs
nav_order: 1
---
```
If you have a sub-page under Installation (e.g., windows.md), you would reference both the parent and grandparent:

```yaml
---
layout: default
title: Windows Installation
parent: Installation
grand_parent: Docs
nav_order: 1
---
```
Adjust nav_order to set the order of pages within the same level of hierarchy. This front matter will tell Just the Docs how to construct your sidebar navigation by creating nested lists that reflect your documentation's structure.

## Write Content

Below the front matter, write your documentation content using Markdown syntax. You can include headings, lists, links, images, code blocks, and other standard Markdown elements.

## Navigation Structure

To define the structure of your documentation, use the front matter in each Markdown file. You can set `parent`, `grand_parent`, etc., to create a nested navigation structure.

## Push Changes

After creating or updating your documentation, push the changes to GitHub. Your GitHub Actions workflow should automatically rebuild the site and publish it to GitHub Pages.

## Local Testing (Optional)

Before pushing to GitHub, you can test your site locally by running `bundle exec jekyll serve` and navigating to `localhost:4000` in your web browser.

As you add more content, the Just the Docs theme will automatically generate a navigation structure based on your front matter and the `nav_order`. Your documentation will be styled according to the theme's design, and it will be responsive and searchable.
