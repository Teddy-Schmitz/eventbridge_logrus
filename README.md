# Amazon EventBridge Hook for [Logrus](https://github.com/sirupsen/logrus)

### Install
> $ go get github.com/teddy-schmitz/eventbridge_logrus

### Usage
```go
package main

import (
	"github.com/sirupsen/logrus"
	"github.com/teddy-schmitz/eventbridge_logrus"
)

func main() {
    hook, err := logrus_eventbridge.NewEventbridgeHook("region", "source", "eventbus")
    if err != nil {
        panic(err)
    }
    logrus.AddHook(hook)
    logrus.WithField("hello", "test").Errorln("an error")
}
```

### EventBridge Patterns
#### Filter by log level
```
{
  "detail": {
    "level": [
      "info"
    ]
  }
}
```


Pull requests welcome!
