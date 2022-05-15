const IV = 'snooping6is9bad0'
const HASH_INDEX = 28
const postBtn = document.getElementById('post-btn')

/* Gets latest chain data from backend */
const getChain = async () => {
    const posts = await fetch('/data', { method: 'GET' })
    const responseBody = await posts.json()

    // Returning the chain from request body, or empty if failure response
    if (responseBody.marker === 1) {
      return responseBody.chain
    }
    return []
}

/* Adds a post to the chain */
const addPost = async post => {
    const response = await fetch('/data/post', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(post)
    })
    return await response.json()
}

/* Adds current chain to frontend */
const imprintChain = () => {
    getChain().then(chain => {
        if (chain.length === 0) {
            document.getElementById('chainground').innerHTML = `
                there are no posts to display.
            `
        } else {
            document.getElementById('chainground').innerHTML = `
                we now have ${chain.length} posts :) 
            `
            console.log(chain)
        }
    })
}

/* Processes an attempt to add a new post */
postBtn.addEventListener('click', async event => {
    event.preventDefault()
    const responseBox = document.getElementById('post-response')
    let tag

    // Parsing tag and hash
    const options = ['pr1', 'pr2', 'pr3', 'pr4', 'pr5', 'pr6', 'pr7']
    for (let option of options) {
        if (document.getElementById(option).checked)
            tag = parseInt(option[2])
    }
    const hash = CryptoJS.SHA256(
        document.getElementById('post-passcode').value).toString()

    // Formatting request body
    const postParams = {
        title: document.getElementById('post-title').value,
        contents: document.getElementById('post-contents').value,
        author: document.getElementById('post-author').value,
        tag,
        hash
    }
    addPost(postParams).then(resp => {
        if (resp.marker === 1) {
            const newCode = unlockRawPasscode(resp.data, hash)
            responseBox.textContent = 
                `The new passcode is ${newCode}; ${resp.message}`
        } else if (resp.marker === 2) {
            const newCode = unlockRawPasscode(
                resp.data, CryptoJS.SHA256('gen6si9').toString())
            responseBox.textContent = `${resp.message}. 
                The next passcode is ${newCode}`
        } else {
            responseBox.textContent = `Failure! ${resp.message}`
        }
    }).catch(err => {
        responseBox.textContent = "Software error occured"
        console.log(err)
    })
})

window.onload = imprintChain

/* Unlocks the new raw passcode from server response using hidden security 
settings in frontend */
const unlockRawPasscode = (ciphertext, storedHash) => {
    const cipherHex = CryptoJS.enc.Hex.parse(ciphertext)
    const parsedKey = CryptoJS.enc.Utf8.parse(
        storedHash.substring(HASH_INDEX, HASH_INDEX+32))
    const parsedIv = CryptoJS.enc.Utf8.parse(IV)

    const cipherCp = { ciphertext: cipherHex }
    const decrypted = CryptoJS.AES.decrypt(
        cipherCp, parsedKey, { iv: parsedIv })
    return decrypted.toString(CryptoJS.enc.Utf8)
}

let modal = document.getElementById("modal");
function modalHandler(val) {
    if (val) {
        fadeIn(modal);
    } else {
        fadeOut(modal);
    }
}
function fadeOut(el) {
    el.style.opacity = 1;
    (function fade() {
        if ((el.style.opacity -= 0.1) < 0) {
            el.style.display = "none";
        } else {
            requestAnimationFrame(fade);
        }
    })();
    el.style.display = 'none'
}
function fadeIn(el, display) {
    el.style.opacity = 0;
    el.style.display = display || "flex";
    (function fade() {
        let val = parseFloat(el.style.opacity);
        if (!((val += 0.2) > 1)) {
            el.style.opacity = val;
            requestAnimationFrame(fade);
        }
    })();
    el.style.display = 'initial'
}