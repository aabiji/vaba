Vaba

**Next steps**
- Speed up searching
  - Spawn a new goroutine for each match we find
  - Use sync.WaitGroup
- Is there a way to specify what duckduckgo result page?
- Make searching DuckDuckGo deterministic
- Make the UI better
  - Make sure styling doesn't break on a mobile view
  - Maybe center the search initially then move it to the top when
    the results is shown
- Scrape from Standard Ebooks
  - Abstract away the GetDownloadLinks() function -- use an interface

Vaba is the app that helps you find and download free ebooks.
It's like a search engine, you enter a search query and you'll
get a list of links to epub or pdf files.

Vaba means "free" in Estonian. I don't know, I think it sounds cool.

How it's going to work:
- The frontend will just act as a wrapper around the api
- We'll download the search results page for the specific query
    - We'll use DuckDuckGo as our search engine, obviously because I
      don't want to get banned by Google.
    - We might use the user's query as a base for our queries
      So for example, if the user inputted "hello" and we want to
      see if someone posted that book on vk, we'd end up searching
      "hello epub vk" or something like that
- Then we'll try to find links within the downloaded html to sites
  we want to check
- Then we'll scrape each site and extract links to epub or pdf files
- We'll repeat the entire process for all the sites we want to search
  (VK, Standard Ebooks, Project Gutenberg, Internet Archive, Github, etc, etc)
- We'll then pool our extracted links and send that over in our response.
- The frontend will take that response and render a nice ui.

Components:
- Frontend -- very simple, just html, css and js
  I don't think we'll need tailwind or react...
  maybe instead of separating our frontend and backend,
  our backend can serve the frontend files and the frontend
  will use a form.
- Backend -- written in Go.
  probably 2 endpoints, one for serving the frontend,
  the other for handling the search endpoint
- Web scraper -- preferably only using Go's standard library -- "from scratch"
  things like find by class or element name...