### smsc

smsc is an HTTP [API](https://smsc.ru/api/) client for
[smsc.ru](https://smsc.ru).

```
go get -u github.com/koorgoo/smsc
```


#### Example

```
package main

import (
	"log"
	"github.com/koorgoo/smsc"
)

func main() {
	c, err := smsc.New(smsc.Config{Login: "me", Password: "secret"})
	if err != nil {
		log.Fatal(err)
	}
	_, err = c.Send("How are you?", []string{"+71234567890"})
	if err != nil {
		log.Fatal(err)
	}
}
```
