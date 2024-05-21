
function renderLinks(links) {
    const div = document.getElementById("links");
    for (let link of links) {
        const a = document.createElement("a");
        a.href = link;
        a.innerHTML = "link"; // TODO: change this
        div.appendChild(a);
    }
}

// Call the /search endpoint on button click
const button = document.getElementById("search");
const query = document.getElementById("query");
button.onclick = () => {
    // TODO: change this back to localhost
    const url = "https://cuddly-barnacle-45xjjqj5x763jr5r-8080.app.github.dev/search";
    const data = {
        method: "POST",
        body: JSON.stringify({query: query.value})
    };

    fetch(url, data)
        .then((response) => response.json())
        .then((json) => {
            const errorMsg = document.getElementById("error");
            if (json.error == undefined) {
                errorMsg.innerHTML = "";
                renderLinks(json.links);
                return;
            }
            errorMsg.innerHTML = json.error;
        });
}