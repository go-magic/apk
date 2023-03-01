package apk

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/avast/apkparser"
)

func GetXmlByApk(apk string) (string, error) {
	buf := bytes.Buffer{}
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "\t")
	zipErr, resErr, manErr := apkparser.ParseApk(apk, enc)
	if zipErr != nil {
		return "", zipErr
	}
	if resErr != nil {
		return "", resErr
	}
	if manErr != nil {
		return "", manErr
	}
	return buf.String(), nil
}

func GetApkInfo(apk string) (*APK, error) {
	appXml, err := GetXmlByApk(apk)
	if err != nil {
		return nil, err
	}
	app := &APK{}
	if err = xml.Unmarshal([]byte(appXml), &app.manifest); err != nil {
		if app.manifest.Package == "" || app.manifest.Application.Label == nil || *app.manifest.Application.Label == "" {
			return nil, err
		}
	}
	activity, err := MainActivity(app.manifest)
	if err != nil {
		return nil, err
	}
	app.manifest.Activity = activity
	return app, nil
}

func MainActivity(manifest *Manifest) (activity string, err error) {
	if manifest.Application.Activities == nil {
		return "", fmt.Errorf("activities invalid")
	}
	for _, act := range *manifest.Application.Activities {
		if act.IntentFilters == nil {
			continue
		}
		for _, intent := range *act.IntentFilters {
			if isMainIntentFilter(intent) {
				return act.Name, nil
			}
		}
	}
	if manifest.Application.ActivityAliases == nil {
		return "", fmt.Errorf("activityAliases invalid")
	}
	for _, act := range *manifest.Application.ActivityAliases {
		if act.IntentFilters == nil {
			continue
		}
		for _, intent := range *act.IntentFilters {
			if isMainIntentFilter(intent) {
				return act.TargetActivity, nil
			}
		}
	}

	return "", fmt.Errorf("no main activity found")
}

func isMainIntentFilter(intent IntentFilter) bool {
	ok := false
	for _, action := range intent.Actions {
		if action.Name == "android.intent.action.MAIN" {
			ok = true
			break
		}
	}
	if !ok {
		return false
	}
	if intent.Categories == nil {
		return false
	}
	ok = false
	for _, category := range *intent.Categories {
		if category.Name == "android.intent.category.LAUNCHER" {
			ok = true
			break
		}
	}
	return ok
}
