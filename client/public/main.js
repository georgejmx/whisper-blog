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
    const responseBody = await response.json()
    if (responseBody.marker == 1) {
        return responseBody.data
    }
    return null
}

/* Adds current chain to frontend */
const imprintChain = () => {
    getChain().then(chain => {
        document.getElementById('chainground').innerHTML = `
            there are no posts to display!
        `
    })
}

/* Processes adding a post */
postBtn.addEventListener('click', async event => {
    event.preventDefault()
    const postParams = {
        title: 'hello' // title box value
    }
    addPost(postParams).then(data => {
        if (data === null) {
            console.log('nully')
        }
    })
})

window.onload = imprintChain


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