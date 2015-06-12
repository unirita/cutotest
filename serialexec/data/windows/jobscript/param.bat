@echo off

echo param.bat start.
echo param1=%1
echo param2=%2
ping localhost -n 2 > nul
echo param.bat end.

exit /b 1