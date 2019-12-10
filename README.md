
## Why?

I needed some sort of exported that exported some Wordpress post data from the API in a format that I could use in a site that was using Netlify CMS and Gatsby.

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

Why not? It's on my list of languages to get better at so I thought I'd give it a go.

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

## What if I have problems?

Um. This isn't a supported project, but let me know if you have any issues. Remember - It is customised for my specific needs.

## Future?

It would be nice to make this a bit more generic. If I have need to do this in the future then maybe I will, maybe not.
