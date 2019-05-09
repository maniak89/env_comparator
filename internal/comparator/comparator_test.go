package comparator

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseYaml(t *testing.T) {
	Convey("parseYaml", t, func() {
		ctx := context.Background()
		data, err := parseYaml(ctx, []byte(`---
spec:
  template:
    spec:
      containers:
      - name: service1
        env:
        - name: key1
          value: "val"
        - name: key2
          value: "val"
---
spec:
  template:
    spec:
      containers:
      - name: service2
        env:
        - name: key1
          value: "val"
        - name: key2
          value: "val"`))
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
	})
}
