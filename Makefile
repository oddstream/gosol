NDK=~/Android/Sdk/ndk-bundle/
ANDROID_HOME=~/Android/Sdk/
GOMOBILE_PATH=~/go/bin/gomobile

#ANDROID_HOME="$HOME/android-studio"
#ANDROID_NDK_HOME="$HOME/android-ndk-r23b"

wasm:
	GOOS=js GOARCH=wasm go build -v -o solitaire.wasm -ldflags="-s -w"

linux:
	go build -v -o solitaire -ldflags="-s -w"

windows:
	GOOS=windows GOARCH=amd64 go build -v -o solitaire.exe -ldflags="-s -w"

android:
	ANDROID_HOME=$(ANDROID_HOME) ebitenmobile bind -target android -javapkg games.oddstream.gomps5 -o ~/gomps/5/android/gomps5.aar .

server:
	python3 -m http.server
