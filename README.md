Paul Miller's personal blog, built with Hugo: https://blog.northbriton.net/

Use Hugo Extended 0.139.4 (matching CI) or the included dev container.

Preview locally:

```sh
hugo server --bind 0.0.0.0
```

Open http://localhost:1313. Edit posts in `content/posts/`, templates in
`layouts/`, and the blog stylesheet in `static/css/style.css`.

Build the static site:

```sh
hugo --gc --minify
```

Output goes to `public/`. Pushes to `master` build and deploy through the
Hugo GitHub Pages workflow in `.github/workflows/hugo.yml`.
