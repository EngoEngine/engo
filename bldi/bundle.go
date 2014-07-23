package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func cp(src, dst string) error {
	cmd := exec.Command("cp", "-R", src, dst)
	return cmd.Run()
}

func buildApp() {
	if err := os.RemoveAll(BUNDLE); err != nil {
		if !os.IsNotExist(err) {
			log.Fatal(err)
		}
	}

	if err := os.Mkdir(BUNDLE, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	if err := os.Mkdir(CONTENTS, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	if err := os.Mkdir(EXE, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	if err := os.Mkdir(RESOURCES, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	err := ioutil.WriteFile(CONTENTS+"/Info.plist", []byte(infoTmpl), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(CONTENTS+"/PkgInfo", []byte(infoTmpl), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	out, err := exec.Command("otool", "-L", "game").Output()
	if err != nil {
		log.Fatal(err)
	}

	err = cp("game", EXE+"/"+NAME)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines[1:] {
		line = strings.Split(strings.TrimSpace(line), " ")[0]
		if len(line) == 0 || strings.HasPrefix(line, "/System/Library") || strings.HasPrefix(line, "/usr/lib") {
			continue
		}
		newpath := EXE + "/" + path.Base(line)
		err = cp(line, newpath)
		if err != nil {
			log.Fatal(err)
		}
		cmd := exec.Command("install_name_tool", "-change", line, newpath, EXE+"/"+NAME)
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	err = cp("data", RESOURCES+"/data")
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("SetFile", "-a", "C", BUNDLE)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

var pkgInfo = "APPL" + NAME

var infoTmpl = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleDevelopmentRegion</key>
    <string>English</string>
    <key>CFBundleExecutable</key>
    <string>Game</string>
    <key>CFBundleGetInfoString</key>
    <string>0.48.2, Copyright 2014 ENGi</string>
    <key>CFBundleIconFile</key>
    <string>Game.icns</string>
    <key>CFBundleIdentifier</key>
    <string>com.ajhager.Game</string>
    <key>CFBundleDocumentTypes</key>
    <array>
    </array>
    <key>CFBundleInfoDictionaryVersion</key>
    <string>6.0</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>0.48.2</string>
    <key>CFBundleSignature</key>
    <string>Game</string>
    <key>CFBundleVersion</key>
    <string>0.48.2</string>
    <key>NSHumanReadableCopyright</key>
    <string>Copyright 2014 ENGi.</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.6</string>
</dict>
</plist>
`
