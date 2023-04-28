# https://github.com/fyne-io/fyne
.PHONY: default build run clean 
default: build run 
build:
	fyne package -os android -appID my.domain.appname -icon icon.png
run:
	fyne install -os android -appID my.domain.appname -icon icon.png
clean:
	rm -f mobile_app.apk 
