
## Why?

I needed some sort of exporter that exported some Wordpress post data from the API in a format that I could use in a site that was using Netlify CMS and Gatsby.

I couldn't find exactly what I needed so thought I'd put a very simple / quick project together.

## What does it do?

Export Wordpress data to Netlify CMS markdown files

1. Gets the following from the WP API...
* 1. Posts
* 2. Users
* 3. Categories
* 4. Tags
* 5. Media data

2. Creates Markdown files of...
* 1. Posts - mapping the user, categories and tags
* 2. Authors

3. Downloads the media files...
* 1. Featured Media
* 2. First image in the post
* 3. Author Avatar (96) version

## Why Go?

Why not? It's on my list of languages to get better at so I thought I'd give it a go. Please note the code is probably not as idomatic as it should be and there are a load of things that could make it better, but this was put together quickly and for my purposes it works as required.

## Can I use it?

Of course, please feel free. Please do let me know if you do and if you have any problems. Remember it is tailored for my specific project, so no guarantees that it'll work with yours!

## How do I use it?

The code was put together on a machine running go1.13.4 darwin/amd64 (on a Mac) so I am assuming you have that version of go installed.

1. download the code
2. go build the code
3. execute ./go-wp-ncms (on a mac - windows will be different) with the following flags:
* 1. -siteURL=http://{site URL} - without the trailing /
* 2. -num= int - number of values to get per page from the WP API - see https://developer.wordpress.org/rest-api/reference/posts/#list-posts
* 3. -refresh= bool - true / false - default false - if you want to refresh the data completely and get all from the API again

or just run ./go-wp-ncms and it'll let you know

## Example Netlify CMS config
This is an example of the Netlify CMS config I am using, which maps to the fields output by the tool.
```
media_folder: static/assets
public_folder: /assets

collections:
  - name: blog
    label: "Post"
    folder: "content/posts"
    create: true
    slug: "{{year}}-{{month}}-{{day}}-{{slug}}.md"
    fields:
      - { label: "Title", name: "title", widget: "string" }
      - { label: "Type", name: "type", widget: "hidden", default: "blog" }
      - {
          label: "Author",
          name: "author",
          widget: "relation",
          collection: "authors",
          searchFields: ["name"],
          valueField: "name",
        }
      - { label: "Publish Date", name: "date", widget: "datetime" }
      - { label: "Featured Image", name: "featuredImage", widget: "image" }
      - {
          label: "Featured",
          name: "featured",
          widget: "boolean",
          default: false,
        }
      - { label: "Category", name: "category", widget: "string" }
      - { label: "Tags", name: "tags", widget: "list" }
      - { label: "Body", name: "body", widget: "markdown" }
  - name: authors
    identifier_field: name
    label: "Author"
    folder: content/author
    create: true
    slug: "{{year}}-{{month}}-{{day}}-{{name}}.md"
    fields:
      - { label: "Name", name: "name", widget: "string" }
      - { label: "Type", name: "type", widget: "hidden", default: "author" }
      - { label: "Short Description", name: "short_desc", widget: "string" }
      - { label: "Image", name: "thumbnail", widget: "image" }
      - { label: "Bio", name: "body", widget: "markdown" }

```

## What if I have problems?

Um. This isn't a supported project, but let me know if you have any issues. Remember - It is customised for my specific needs.

## Future?

It would be nice to make this a bit more generic. If I have need to do this in the future then maybe I will, maybe not.
