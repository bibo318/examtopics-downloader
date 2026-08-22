@echo off
echo [ExamDemon] Building...
go build -ldflags="-s -w" -o examdemon.exe .
if %ERRORLEVEL% == 0 (
    echo [OK] Build thanh cong: examdemon.exe
    for %%F in (examdemon.exe) do echo     Kich thuoc: %%~zF bytes
) else (
    echo [LOI] Build that bai!
)
