const formatter = new Intl.DateTimeFormat("en-DE", {dateStyle: "full", timeStyle: "long"});

for (let elem of document.getElementsByClassName("timestamp")) {
    elem.childNodes[0].data = formatter.format(new Date(Number(elem.childNodes[0].data)));
}
