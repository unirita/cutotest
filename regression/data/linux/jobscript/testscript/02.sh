#!/bin/sh

if [ $RC -ne 0 ]; then
    echo "[err] Invalid RC : $RC"
fi
if [ $SD -eq "" ]; then
    echo "[err] Invalid SD : $SD"
fi
if [ $ED -eq "" ]; then
    echo "[err] Invalid ED : $ED"
fi
if [ $OUT -ne "" ]; then
    echo "[err] Invalid OUT : $OUT"
fi

exit 0
