const fileInput = document.getElementById("file");

function checkFileSize() {
    const maxSize = Number(fileInput.getAttribute("data-maxsize"));

    for (const file of fileInput.files) {
        if (file.size > maxSize) {
            fileInput.setCustomValidity("The file is too big");
            fileInput.reportValidity();
            return;
        }
    }

    fileInput.setCustomValidity("");
}

fileInput.addEventListener("change", checkFileSize);
