#!/bin/sh

OSNAME=`uname`
DATAFILE=dispres.sql
TODAY=`date +"%Y-%m-%d"`
if [ $OSNAME = "Darwin" ] ; then
        DAY1=`date -v-1d +"%Y-%m-%d"`
        DAY2=`date -v-2d +"%Y-%m-%d"`
        DAY3=`date -v-3d +"%Y-%m-%d"`
        DAY4=`date -v-4d +"%Y-%m-%d"`
else
        DAY1=`date -d "1days ago" +"%Y-%m-%d"`
        DAY2=`date -d "2days ago" +"%Y-%m-%d"`
        DAY3=`date -d "3days ago" +"%Y-%m-%d"`
        DAY4=`date -d "4days ago" +"%Y-%m-%d"`
fi

if [ ! -s ${DATAFILE}.org ] ; then
        echo "<error> Not found ${DATAFILE}.org"
        exit 1
fi
if [ -s ${DATAFILE} ] ; then
        rm ${DATAFILE}
fi
if [ -s ${DATAFILE}ite ] ; then
        rm ${DATAFILE}ite
fi
cp ${DATAFILE}.org ${DATAFILE}

REP=s/@TODAY/${TODAY}/g
echo $REP > .datachange
REP=s/@DAY1/${DAY1}/g
echo $REP >> .datachange
REP=s/@DAY2/${DAY2}/g
echo $REP >> .datachange
REP=s/@DAY3/${DAY3}/g
echo $REP >> .datachange
REP=s/@DAY4/${DAY4}/g
echo $REP >> .datachange

sed -f ./.datachange ${DATAFILE}.org > ${DATAFILE}.temp 
mv ${DATAFILE}.temp ${DATAFILE}

sqlite3 ${DATAFILE}ite < dbinit.sql

exit $?

