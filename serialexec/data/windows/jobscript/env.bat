@echo off

echo env.bat start.
echo ENV1=%ENV1%
echo ENV2=%ENV2%
ping localhost -n 2 > nul
echo env.bat end.

exit /b 2