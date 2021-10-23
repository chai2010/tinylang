package main

func main() {
	for _, c := range []byte(s[1:]) {
		println(`write 1024 * 1024 +`, c, ";")
	}
}

const s = `
+---+    +---+
| o |    | o |
|   +----+   |
|            |
|    1024    |
|            |
+------------+
`
