@echo off

echo work.bat start.
echo current=%cd%
ping localhost -n 2 > nul
echo work.bat end.

exit /b 3