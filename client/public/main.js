const IV = 'snooping6is9bad0'
const HASH_INDEX = 28

/* Gets latest raw chain data from backend */
const getChainHtml = async () => {
    const posts = await fetch('/html/chain', { method: 'GET' })
    return await posts.text()
}

/* Adds a post to the chain */
const addPostData = async post => {
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
    getChainHtml().then(content => {
        if (content.length > 1) {
            document.getElementById('deck').innerHTML = content
        } else {
            document.getElementById('deck').innerHTML = `<h2 class="text-lg
            text-white">Go make the first post! No passcode is required.</h2>`
        }
    }).catch(err => {
        document.getElementById('deck').innerHTML = `<h2 class="text-lg
        text-red-400">Error fetching chain from server.</h2>`
        console.error(err)
    })
}

/* Processes an attempt to add a new post */
const addPost = () => {
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

    // Communicating with server; attempting *addPost* then showing result
    addPostData(postParams).then(resp => {
        if (resp.marker === 1) {
            const newCode = unlockRawPasscode(resp.data, hash)
            responseBox.textContent = 
                `The new passcode is ${newCode}; ${resp.message}`

            // Refreshing chain html
            imprintChain()
            document.getElementById('add-modal-tr').textContent = 'Show passcode'
        } else if (resp.marker === 2) {
            const newCode = unlockRawPasscode(
                resp.data, CryptoJS.SHA256('gen6si9').toString())
            responseBox.textContent = `${resp.message}. 
                The next passcode is ${newCode}`

            // Refreshing chain html
            imprintChain()
            document.getElementById('add-modal-tr').textContent = 'Show passcode'
        } else {
            responseBox.textContent = `Failure! ${resp.message}`
        }
    }).catch(err => {
        responseBox.textContent = 'Error adding post'
        console.error(err)
    })
}

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

/* Toggles whether the add post modal is shown or not */
const addModalHandler = isOpen => {
    if (isOpen) {
        document.getElementById('add-modal').style.display = 'initial'
    } else {
        document.getElementById('add-modal').style.display = 'none'
    }
}

/* Toggles whether the add reaction modal is shown or not, also gets the modal
contents from server dependent on state */
const reactModalHandler = val => {
    // Handle the case where we want to close the modal
    if (val == 0 || !val) {
        document.getElementById('react-modal').style.display = 'none'
        return
    }

    // getReactionDeckHtml(val).then .....
    document.getElementById('react-modal').style.display = 'initial'
}

window.onload = imprintChain