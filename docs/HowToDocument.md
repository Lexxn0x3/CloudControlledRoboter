---
title: How to document
layout: default
nav_order: 1
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

## Write Content

Below the front matter, write your documentation content using Markdown syntax. You can include headings, lists, links, images, code blocks, and other standard Markdown elements.

## Navigation Structure

To define the structure of your documentation, use the front matter in each Markdown file. You can set `parent`, `grand_parent`, etc., to create a nested navigation structure.

## Push Changes

After creating or updating your documentation, push the changes to GitHub. Your GitHub Actions workflow should automatically rebuild the site and publish it to GitHub Pages.

## Local Testing (Optional)

Before pushing to GitHub, you can test your site locally by running `bundle exec jekyll serve` and navigating to `localhost:4000` in your web browser.

As you add more content, the Just the Docs theme will automatically generate a navigation structure based on your front matter and the `nav_order`. Your documentation will be styled according to the theme's design, and it will be responsive and searchable.
