package main

const (
	host = "192.168.7.199"
	port = "21"
	user = "admin"
	pwd  = "abc123!"
)

func main() {
	c, err := Link(host, port, user, pwd)
	if err != nil {
		panic(err)
	}
	defer CloseLink(c)

	upload(c, "test.pdf", "test5.pdf")
	download(c, "test5.pdf", "test10.pdf", 4096)
}
