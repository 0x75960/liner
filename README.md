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

	handler := func(line string) (err error) {
		fmt.Println(line)
		return
	}

	proc := liner.NewLineProcessor(handler, liner.LogWhenError)

	f, err := os.Open("/path/to/files")
	if err != nil {
		log.Fatalln(err)
	}

	defer f.Close()

	proc(f)

}

```
