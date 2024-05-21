
// Call the /search endpoint on button click
const button = document.getElementById("search");
const query = document.getElementById("query");
button.onclick = () => {
    const url = "http://localhost:8080/search";
    const data = {
        method: "POST",
        body: JSON.stringify({query: query.value})
    };

    fetch(url, data)
        .then((response) => response.json())
        .then((json) => console.log(json));
}