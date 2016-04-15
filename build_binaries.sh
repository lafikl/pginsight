env GOOS=darwin GOARCH=386 go build && mv pginsight binaries/pginsight_darwin_x86 &&
env GOOS=linux GOARCH=386 go build && mv pginsight binaries/pginsight_linux_x86 &&
env GOOS=openbsd GOARCH=386 go build && mv pginsight binaries/pginsight_openbsd_x86 &&
env GOOS=freebsd GOARCH=386 go build && mv pginsight binaries/pginsight_freebsd_x86 &&
env GOOS=windows GOARCH=386 go build && mv pginsight.exe binaries/pginsight_window_x86.exe
