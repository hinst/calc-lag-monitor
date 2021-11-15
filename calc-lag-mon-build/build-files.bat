cd /D "%~dp0"
rmdir /s /q app
del ..\calc_lag_monitor\compiled\calc_lag_monitor
rmdir /s /q ..\calc-lag-mon-ui\dist
mkdir app

cd /D "%~dp0"
echo "[ ] Build server..."
call ..\calc_lag_monitor\build-linux.bat

cd /D "%~dp0"
echo "[ ] Build web UI..."
set /p API_URL=<api-url.txt
call ..\calc-lag-mon-ui\build.bat

cd /D "%~dp0"
echo "[ ] Copying..."
xcopy /y "..\calc_lag_monitor\compiled\calc_lag_monitor" ".\app\calc_lag_monitor\compiled\*"
copy configuration.json .\app\calc_lag_monitor\configuration.json
xcopy /y /E "..\calc-lag-mon-ui\dist\*" ".\app\calc-lag-mon-ui\dist\*"
echo "Completed."