const shortFormatter = new Intl.DateTimeFormat("de-DE", {dateStyle: "short", timeStyle: "short"});
const longFormatter = new Intl.DateTimeFormat("en-DE", {dateStyle: "full", timeStyle: "long"});

for (let elem of document.getElementsByClassName("timestamp")) {
    const date = new Date(Number(elem.childNodes[0].data));
    elem.childNodes[0].data = shortFormatter.format(date);
    elem.setAttribute("title", longFormatter.format(date));
}
