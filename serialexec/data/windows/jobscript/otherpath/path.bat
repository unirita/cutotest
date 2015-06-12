@echo off

echo path.bat start.
echo filepath=%~dpf0
ping localhost -n 2 > nul
echo path.bat end.

exit /b 4
