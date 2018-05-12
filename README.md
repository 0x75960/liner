liner
=====

line processing utility for golang

```sh
go get -u github.com/0x75960
```

usage
------

```go
import (

	// ...

	"github.com/0x75960/liner"
)

func main() {

	f, err := os.Open("/path/to/files")
	if err != nil {
		log.Fatalln(err)
	}

	defer f.Close()

	for line := range liner.LinesIn(f) {
		fmt.Println(line)
	}

	proc(f)

}

```
