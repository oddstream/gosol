TARGET=gosol
NDK=~/Android/Sdk/ndk-bundle/
ANDROID_HOME=~/Android/Sdk/
GOMOBILE_PATH=~/go/bin/gomobile

#ANDROID_HOME="$HOME/android-studio"
#ANDROID_NDK_HOME="$HOME/android-ndk-r23b"

wasm: Makefile
	GOOS=js GOARCH=wasm go build -v -o $(TARGET).wasm -ldflags="-s -w"

linux: Makefile
	go build -v -o $(TARGET) -ldflags="-s -w"

windows: Makefile
	GOOS=windows GOARCH=amd64 go build -v -o $(TARGET).exe -ldflags="-s -w"

#android: Makefile
#	ANDROID_HOME=$(ANDROID_HOME) ebitenmobile bind -target android -javapkg games.oddstream.$(TARGET) -o ~/gomps/5/android/$(TARGET).aar .

server:
	python3 -m http.server
