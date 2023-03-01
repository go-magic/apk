package apk

import (
	"fmt"
	"testing"
)

func TestNewXMLFile(t *testing.T) {
	app, err := GetApkInfo("./base1.apk")
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(app.manifest.Package)
	fmt.Println(*app.manifest.Application.Label)
	fmt.Println(app.manifest.Activity)
}
