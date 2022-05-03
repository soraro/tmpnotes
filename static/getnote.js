document.getElementById("dec-button").addEventListener("click", decryptNote)
document.getElementById("copy-note-button").addEventListener("click", copyNote)
document.getElementById("get-note-button").addEventListener("click", getNote)

async function getNote() {
    var opts = {
        method: 'GET',
        headers: {
            'X-Note': 'Destroy'
        },
    }

    var req = new Request(location.href, opts)

    await fetch(req)
        .then(function (response) {
            if (!response.ok) {
                location.reload()
            } else {
                return response.text().then(function (text) {
                    document.getElementById("note-row").style.visibility = "visible"
                    document.getElementById("notification-row").style.visibility = "visible"
                    document.getElementById("note").value = text
                    document.getElementById("get-note-row").remove()
                    checkDecryptionRequired()
                });
            }
        })
        .catch(error => {
            console.log('request failed', error);
            alert("Unable to reach the server. Try again later.")
            throw error;
        });
}

function copyNote() {
    let textarea = document.getElementById("note");
    textarea.select();
    document.execCommand("copy");
}

function checkDecryptionRequired() {
    note = document.getElementById("note");
    if (note.value.startsWith("[ENC]")) {
        e = document.getElementById("dec-button");
        e.style.visibility = "visible";
        e = document.getElementById("dec-input");
        e.style.visibility = "visible";
    } else {
        document.getElementById("decrypt-row").remove()
    }
}

function decryptNote() {
    note = document.getElementById("note");
    notetext = note.value.replace("[ENC]", "");
    key = document.getElementById("dec-key").value;

    try {
        var bytes = CryptoJS.AES.decrypt(notetext, key);
        var res = bytes.toString(CryptoJS.enc.Utf8);

        if (res != "") {
            note.value = res
            document.getElementById("decrypt-row").remove()
        } else {
            alert("Decryption failed")
        }
    } catch (error) {
        console.error(error);
        alert("Decryption failed")
    }
}