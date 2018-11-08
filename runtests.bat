@echo off

go test -tags appveyor -v ./...
cd demos
for /d %%d in (*.*) do call :for_body %%d
echo done
exit /b

:for_body
if %1==demoutils (
  echo Skipping %1
  goto :cont
)
cd %1
echo Verifying %1 ...
go build -tags demo -o %1.exe
cd ..
:cont
exit /b
