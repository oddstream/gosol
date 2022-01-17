TARGET=gosol
NDK=~/Android/Sdk/ndk-bundle/
ANDROID_HOME=~/Android/Sdk/
GOMOBILE_PATH=~/go/bin/gomobile

#ANDROID_HOME="$HOME/android-studio"
#ANDROID_NDK_HOME="$HOME/android-ndk-r23b"

android: Makefile
	ANDROID_HOME=$(ANDROID_HOME) ebitenmobile bind -target android -javapkg games.oddstream.$(TARGET) -o ~/gomps/5/android/$(TARGET).aar .
#	ANDROID_HOME=$(ANDROID_HOME) gomobile build -v -target android

server:
	python3 -m http.server
