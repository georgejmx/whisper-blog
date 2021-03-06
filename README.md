# WhisperBlog

### https://whisperblog.xyz

---

Interactive blog, with a unique chain-like structure

## Usage

*WhisperBlog* is a basic social media site that works entirely sequentially,
where making a new post can only happen with the passcode. The initial post,
from **AddPost** is made with creation of the chain. 
From this point, making a post randomly generates the new passcode. This
can then be given to someone new, who then creates the next post.

This was designed to be a digital chinese whispers,
posts then circulate in turn around a small
social group. The chain could also move around society if the passcode
is given to a more broad range of people. Uniquely, *WhisperBlog*
does not require user accounts like mainstream social media, but also
facilitates a chain of trust between posters unlike anonymous social
media platforms.

Clicking **Vote** on each post allows
reacting to this post, where possible reactions are random English
adjectives. There is room for 6 anonymous reactions, where
further reactions can be made by providing a passcode that was valid for
the 3 previous posts. This means that attributed reactions, that carry
more weight, can outvote any attempt to spam reactions.

After 5 days of inactivity, the previous passcode can also make a new post. This
ensures that the chain does not get stuck with one person. Then every 2 days,
the next previous person can also make a post using their previous passcode.

## How to build and run

- Will need sqlite3 and go1.18 installed
- Clone the repo into your GOROOT, then modify your config package to look like;

```
package config

import "os"

var (
	DB_FILEPATH      string
	AES_IV           string // must be of length 32
	AES_SPLICE_INDEX string // must be a string parsable to >=0 and <= 31
)

/* Sets environment variables ued by program. Will be different for integration
tests than in production */
func SetupEnv(isProduction bool) {
	if isProduction {
		DB_FILEPATH = "./data/blog.db"
		AES_IV = "[YOUR IV]"
		AES_SPLICE_INDEX = "[YOUR SPLICE INDEX]"
	} else {
		DB_FILEPATH = "./data/blog_test.db"
		AES_IV = "snooping6is9bad0"
		AES_SPLICE_INDEX = "28"
	}
	os.Setenv("DB_FILEPATH", DB_FILEPATH)
	os.Setenv("AES_IV", AES_IV)
	os.Setenv("AES_SPLICE_INDEX", AES_SPLICE_INDEX)
}
```

Also adjust the constants at the top of *client/src/main.js* to match the above
values for `[YOUR IV]` and `[YOUR SPLICE INDEX]`. This means rebuilding your
obfuscated *client/public/main.js*

- Run `go test ./...` to ensure your build is stable
- Execute `go build -o wb` to generate a binary
- Put this binary wherever you like with a valid HTTPS key pair, next to a
blank *data/* directory in which the database will be generated.
**./wb** will spin the whole thing up with only sqlite needed
