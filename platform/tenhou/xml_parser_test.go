package tenhou

import (
	"testing"
	"io/ioutil"
	"encoding/xml"
	"fmt"
)

func Test(t *testing.T) {
	data, err := ioutil.ReadFile("xx.xml")
	if err != nil {
		t.Fatal(err)
	}

	d := Record{}
	if err := xml.Unmarshal(data, &d); err != nil {
		t.Fatal(err)
	}

	for _, action := range d.Actions {
		fmt.Println(*action)
	}
}
