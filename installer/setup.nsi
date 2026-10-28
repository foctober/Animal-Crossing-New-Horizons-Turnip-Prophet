Unicode True
!include "MUI2.nsh"
!include "x64.nsh"
!include "FileFunc.nsh"

Name "Turnip Prophet"
OutFile "..\TurnipProphet-Setup.exe"
InstallDir "$PROGRAMFILES64\Turnip Prophet"
InstallDirRegKey HKLM "Software\TurnipProphet" "InstallDir"
RequestExecutionLevel admin
SetCompressor /SOLID lzma

!define APP_VERSION "1.2.0"
VIProductVersion "1.2.0.0"
VIFileVersion "1.2.0.0"
VIAddVersionKey /LANG=1033 "ProductName" "Turnip Prophet"
VIAddVersionKey /LANG=1033 "ProductVersion" "1.2.0"
VIAddVersionKey /LANG=1033 "FileVersion" "1.2.0"
VIAddVersionKey /LANG=1033 "FileDescription" "Turnip Prophet 1.2 setup"
VIAddVersionKey /LANG=1033 "CompanyName" "foctober"

!define MUI_ICON "..\TurnipProphet.ico"
!define MUI_UNICON "..\TurnipProphet.ico"
!define MUI_ABORTWARNING
!define MUI_WELCOMEPAGE_TITLE "Turnip Prophet 1.2"
!define MUI_WELCOMEPAGE_TEXT "This installs Turnip Prophet 1.2 into Program Files.$\r$\n$\r$\nQuit Turnip Prophet first if it is already open. Closing the window quits it.$\r$\n$\r$\nYour weeks stay in your Windows profile. Uninstall does not delete them."
!define MUI_FINISHPAGE_TEXT "Turnip Prophet is installed in Program Files.$\r$\n$\r$\nIt is in the Start menu. A desktop icon is there too if you left that box checked.$\r$\n$\r$\nThere is no system tray icon. Closing the window quits the app.$\r$\n$\r$\nYour weeks stay in AppData\Roaming\TurnipProphet. Uninstall leaves that folder alone."
!define MUI_FINISHPAGE_RUN
!define MUI_FINISHPAGE_RUN_TEXT "Open Turnip Prophet"
!define MUI_FINISHPAGE_RUN_FUNCTION LaunchApp
!define MUI_COMPONENTSPAGE_SMALLDESC

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_COMPONENTS
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_LANGUAGE "English"

Function .onInit
  ${IfNot} ${RunningX64}
    MessageBox MB_OK|MB_ICONSTOP "Turnip Prophet needs 64-bit Windows 10 or 11."
    Abort
  ${EndIf}
  SetRegView 64
FunctionEnd

Function un.onInit
  SetRegView 64
FunctionEnd

Section "Turnip Prophet" SecCore
  SectionIn RO
  SetShellVarContext all
  SetOutPath "$INSTDIR"
  File /oname=TurnipProphet.exe "..\build\TurnipProphet.exe"
  File "..\predictions.js"
  File "..\TurnipProphet.ico"
  File "..\LICENSE"
  File "..\NOTICE"
  File "..\README.md"
  File "..\INSTALL.md"
  SetOutPath "$INSTDIR\www"
  File "..\www\index.html"
  File "..\www\app.js"
  File "..\www\styles.css"
  File "..\www\logo.png"
  File "..\www\field.jpg"
  File "..\www\wood.png"
  File "..\www\FinkHeavy.ttf"

  SetOutPath "$INSTDIR"
  WriteUninstaller "$INSTDIR\Uninstall.exe"

  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TurnipProphet" "DisplayName" "Turnip Prophet"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TurnipProphet" "DisplayVersion" "1.2.0"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TurnipProphet" "UninstallString" '"$INSTDIR\Uninstall.exe"'
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TurnipProphet" "DisplayIcon" "$INSTDIR\TurnipProphet.ico"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TurnipProphet" "Publisher" "foctober"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TurnipProphet" "InstallLocation" "$INSTDIR"
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TurnipProphet" "NoModify" 1
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TurnipProphet" "NoRepair" 1
  WriteRegStr HKLM "Software\TurnipProphet" "InstallDir" "$INSTDIR"
  ${GetSize} "$INSTDIR" "/S=0K" $0 $1 $2
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TurnipProphet" "EstimatedSize" "$0"

  CreateShortcut "$SMPROGRAMS\Turnip Prophet.lnk" "$INSTDIR\TurnipProphet.exe" "" "$INSTDIR\TurnipProphet.ico" 0
SectionEnd

Section "Desktop icon" SecDesktop
  SetShellVarContext all
  CreateShortcut "$DESKTOP\Turnip Prophet.lnk" "$INSTDIR\TurnipProphet.exe" "" "$INSTDIR\TurnipProphet.ico" 0
SectionEnd

LangString DESC_SecCore ${LANG_ENGLISH} "Install Turnip Prophet into Program Files and add a Start menu shortcut."
LangString DESC_SecDesktop ${LANG_ENGLISH} "Also put a Turnip Prophet icon on the desktop. Uncheck this if you do not want one. The Start menu shortcut is always installed."

!insertmacro MUI_FUNCTION_DESCRIPTION_BEGIN
  !insertmacro MUI_DESCRIPTION_TEXT ${SecCore} $(DESC_SecCore)
  !insertmacro MUI_DESCRIPTION_TEXT ${SecDesktop} $(DESC_SecDesktop)
!insertmacro MUI_FUNCTION_DESCRIPTION_END

Function LaunchApp
  Exec '"$INSTDIR\TurnipProphet.exe"'
FunctionEnd

Section "Uninstall"
  SetShellVarContext all
  Delete "$DESKTOP\Turnip Prophet.lnk"
  Delete "$SMPROGRAMS\Turnip Prophet.lnk"
  Delete "$INSTDIR\TurnipProphet.exe"
  Delete "$INSTDIR\TurnipProphet.ico"
  Delete "$INSTDIR\predictions.js"
  Delete "$INSTDIR\LICENSE"
  Delete "$INSTDIR\NOTICE"
  Delete "$INSTDIR\README.md"
  Delete "$INSTDIR\INSTALL.md"
  Delete "$INSTDIR\www\index.html"
  Delete "$INSTDIR\www\app.js"
  Delete "$INSTDIR\www\styles.css"
  Delete "$INSTDIR\www\logo.png"
  Delete "$INSTDIR\www\field.jpg"
  Delete "$INSTDIR\www\wood.png"
  Delete "$INSTDIR\www\FinkHeavy.ttf"
  RMDir "$INSTDIR\www"
  Delete "$INSTDIR\Uninstall.exe"
  RMDir "$INSTDIR"
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\TurnipProphet"
  DeleteRegKey HKLM "Software\TurnipProphet"
SectionEnd
