### Vaba

Vaba is an app that helps you find and download free ebooks.
It's like a search engine, you enter a search query and you'll
get a list of links to epub or pdf files.

It works this way. Vaba will search duckduckgo for the ebook
and scrape and collect the websites links. Then for each
website link, it'll scrape the associated page and extract
a link to a file from that page. After doing that for each
website, it'll just send a response containing to those links
to the frontend. The frontend will take that and render to the user.

At this point, Vaba is basically "done". Obviously there's more
than can be done, like scraping from Standard Ebooks or something,
but as it is, Vaba works pretty well and I'm moving on.

### Run instructions
```
# First clone the repo
git clone https://github.com/aabiji/vaba.git

# Then build
cd vaba
go run .

# Then use by going to http://localhost:8080 in your browser.
```

Vaba was a fun little project for me and maybe you'll enjoy it too.
