@echo off

cd /d "%~dp0"
call ..\setparm.bat

if exist result (
  rd /s /q result
)
mkdir result

setlocal
set path=%CUTOROOT%\bin;%path%
cd .\result

set num=0

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n NoServiceTask > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n NoServiceTask -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n NoStartEvent > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n NoStartEvent -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n NoEndEvent > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n NoEndEvent -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n MultiStartEvent > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n MultiStartEvent -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n MultiEndEvent > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n MultiEndEvent -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set num=12
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName1 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName1 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName2 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName2 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName3 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName3 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName4 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName4 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName5 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName5 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName6 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName6 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName7 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName7 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName8 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName8 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName9 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName9 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName10 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName10 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName11 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n ForbiddenJobName11 -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set num=35
master.exe -c %CUTOROOT%\bin\master.ini -n StartWithoutStartEvent -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n StartWithoutStartEvent -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n EndWithoutEndEvent -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n EndWithoutEndEvent -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n Isolation -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n Isolation -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n DuplicateID -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n DuplicateID -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n BranchWithoutGateway -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n BranchWithoutGateway -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n MergeWithoutGateway -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n MergeWithoutGateway -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n EndBeforeMerge -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n EndBeforeMerge -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n NestedBranch -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n NestedBranch -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n NotMerge -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

set /A num+=1
master.exe -c %CUTOROOT%\bin\master.ini -n NotMerge -s > pf%num%.log
echo %errorlevel% >> pf%num%.log

exit /b 0
