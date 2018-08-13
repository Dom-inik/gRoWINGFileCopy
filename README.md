# gRoWINGFileCopy

Simple tool to copy growing and closed files.

## Build and run from source

If you have already set up your go enviroment run 

```powershell
go get -u github.com/Dom-inik/gRoWINGFileCopy
```

## Use the prebuild executeable

```Powershell
PS D:\> [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
PS D:\> Invoke-WebRequest -Uri https://github.com/Dom-inik/gRoWINGFileCopy/raw/master/growingfilecopy.exe -OutFile growingfilecopy.exe
PS D:\> .\growingfilecopy.exe -h
2018/08/13 14:12:23 Started ...
Usage of D:\growingfilecopy.exe:
  -dst string
        destination file path
  -src string
        source file path
PS D:\>
```