# tea
An experimental web-framework in Go.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/marconi/tea"
)

func main() {
	t := tea.New()
	t.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World!")
	})

	t.Logger.Fatalln(t.Start("localhost:8080"))
}
```

Output:

```
INFO[2019-03-05T09:25:17+08:00] 🍵  is served on localhost:8080  
```
