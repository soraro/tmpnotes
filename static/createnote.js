document.addEventListener('keydown', countInput);
document.addEventListener('keyup', countInput);
document.getElementById("send-button").addEventListener("click", postInfo)
document.getElementById("copy-button").addEventListener("click", copyLink)
const maxChars = Number(document.getElementById("ui-expire").textContent)

function countInput(e) {
    const input = document.getElementById('notes-input');
    input.labels[0].innerText = input.textLength + "/" + maxChars
}

function postInfo() {
    var e = document.getElementById("hours-input");
    var hours = e.options[e.selectedIndex].text;
    var e = document.getElementById("notes-input");
    var key = document.getElementById("enc-key");
    var note = encryptNote(e.value, key.value);
    e.readOnly = true;
    document.getElementById("send-button").style.visibility = "hidden";
    var data = { "message": note, "expire": parseInt(hours) }
    var opts = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    }

    var req = new Request(location.origin + "/new", opts)

    fetch(req).then(function (response) {
        return response.text().then(function (text) {
            e = document.getElementById("note-id");
            e.value = location.origin + "/id/" + text;
        });
    })
        .catch(error => {
            console.log('request failed', error);
            alert("Unable to reach the server. Try again later.")
            throw error;
        });

    e = document.getElementById("note-id");
    e.style.visibility = "visible";
    document.getElementById("copy-button").style.visibility = "visible";
    document.getElementById("encrypt-row").remove();
}

function copyLink() {
    let noteid = document.getElementById("note-id");
    noteid.select();
    document.execCommand('copy');
}

// encrypt the note if an encryption key is provided
function encryptNote(note, key) {
    if (key.length > 0) {
        var ciphertext = CryptoJS.AES.encrypt(note, key);
        return "[ENC]" + ciphertext.toString()
    } else {
        return note;
    }
}