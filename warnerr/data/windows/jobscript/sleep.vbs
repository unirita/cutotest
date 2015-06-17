Option Explicit

Dim args
Set args = WScript.Arguments

if args.Count < 1 then
WScript.sleep(1000)
else
WScript.sleep(args.item(0)*1000)
end if
