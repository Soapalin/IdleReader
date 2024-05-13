# Compilation of go executables for different platforms

By setting the `GOOS` and the `GOARCH` environment variables, target the appropriate platform.


linux and mac
```bash
env GOOS=linux GOARCH=arm go build -v github.com/path/to/your/app
``` 

Note: env lets you set the environment variables for just thsi command

windows: 
```bash
set GOOS=windows& set GOARCH=amd64& go build -v .
set GOOS=linux& set GOARCH=arm& go build -v .
```

Note: double-clicking the executable on windows 10 will run it on cmd.exe. To have the best experience (emojis and such), run it with Windows Terminal, which should be the default command line experience for windows 11.


# Check architecture on Windows platform

```bash
echo %PROCESSOR_ARCHITECTURE%
```


