const crypto = require('crypto-js')
const axios = require('axios').default

/*** 
 * This is where the current blog passcode should be provided for testing.
 * The new passcode will be dumped to console, then pasted here for more testing
 ***/
const PASSCODE = 'Tr83OAZ7tfUC'

// Stores url and encryption settings that are mangled away from snoopers
const URL = "http://localhost:8080"
const IV = 'snooping6is9bad0'
const HASH_INDEX = 28

/* Sends an automated new post to backend */
const newPost = async hash => {
    const response = await axios.post(`${URL}/data/post`, {
        title: `test post ${Math.floor(Math.random()*1000000)}`,
        author: 'tester',
        contents: 'an automated test to illustrate and simulate adding a post',
        tag: 0,
        hash
    })
    return response.data
}

/* Unlocks the new raw passcode from server response using hidden security 
settings in frontend */
const unlockRawPasscode = (ciphertext, storedHash) => {
    const cipherHex = crypto.enc.Hex.parse(ciphertext)
    const parsedKey = crypto.enc.Utf8.parse(
        storedHash.substring(HASH_INDEX, HASH_INDEX+32))
    const parsedIv = crypto.enc.Utf8.parse(IV)

    const cipherCp = { ciphertext: cipherHex }
    const decrypted = crypto.AES.decrypt(
        cipherCp, parsedKey, { iv: parsedIv })
    return decrypted.toString(crypto.enc.Utf8)
}

/* Covers the UI workflow */
const client = () => {
    const hash = crypto.SHA256(PASSCODE).toString()
    newPost(hash).then(body => {
        console.log(`new passcode is: ${unlockRawPasscode(body.data, hash)}`)
    }).catch(err => { console.log('invalid request') })
}

client()