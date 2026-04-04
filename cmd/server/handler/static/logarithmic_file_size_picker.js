const fileSizeLabel = document.getElementById("file_size_label");
const fileSizeValue = document.getElementById("file_size_value");
const fileSizePicker = document.getElementById("file_size_picker");

const MIN = 4096; // 4 KiB
const MAX = 128 * 1024 ** 4; // 128 TiB
const steps = Math.log2(MAX) - Math.log2(MIN);
fileSizePicker.setAttribute("max", String(steps));

const snappingTable = [];
for (let i = MIN; i <= MAX; i *= 2) {
    snappingTable.push(i);
}

function bytesToHuman(bytes) {
    const sizes = [
        "bytes",
        "KiB",
        "MiB",
        "GiB",
        "TiB",
        "PiB",
    ];

    let i = 0;
    while (bytes >= 1024 && i < sizes.length - 1) {
        i++;
        bytes /= 1024;
    }

    if (Math.floor(bytes) == bytes) {
        bytes = String(bytes);
    } else {
        bytes = bytes.toFixed(3);
    }
    return `${bytes} ${sizes[i]}`;
}

function update(bytes) {
    fileSizeValue.value = String(bytes);
    fileSizeLabel.textContent = bytesToHuman(bytes);
}

let i = 0;
while (i < snappingTable.length-1) {
    if (snappingTable[i] > Number(fileSizeValue.value)) {
        break
    }
    i++
}
fileSizePicker.value = String(i);

update(fileSizeValue.value);

fileSizePicker.addEventListener("input", (event) => {
    update(snappingTable[Number(event.target.value)]);
});
