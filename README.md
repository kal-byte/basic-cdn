Hai so you need a .env file, with "TOKEN" as a var, set it to whatever I guess.

Upload requests go:
http://ip:7070/upload?token=TOKEN_HERE

You're also going to need to `go get` these:

github.com/gin-gonic/gin
github.com/joho/godotenv

I apologize for any newbieness here, I am not quite experience with Go yet, I promise I'll get better!