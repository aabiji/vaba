
const button  = document.getElementById("search");
const query   = document.getElementById("query");
const message = document.getElementById("error");
const results = document.getElementById("results");
const loader  = document.getElementById("loader");

const show = (element) => element.style.display = "block";
const hide = (element) => element.style.display = "none";

const renderLinks = (links) => {
    for (let link of links) {
        const a = document.createElement("a");
        a.href = link.Href;
        a.innerHTML = link.Name;
        results.appendChild(a);
    }
}

const search = () => {
    if (query.value.length == 0) return;

    hide(message);
    hide(results);
    show(loader);

    const url = "http://localhost:8080/search";
    const data = {
        method: "POST",
        body: JSON.stringify({query: query.value})
    };

    fetch(url, data)
        .then((response) => response.json())
        .then((json) => {
            hide(loader);

            if (json.error != undefined) {
                show(message);
                message.innerHTML = json.error;
                return;
            }

            show(results);
            renderLinks(json.links);
        });
}

// Call the /search endpoint on button click
// or on Enter
button.onclick = () => search();
query.onkeydown = (event) => {
    if (event.key == "Enter") {
        search();
    }
}