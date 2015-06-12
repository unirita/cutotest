#/bin/sh

. ../setparm.sh

if [ -d ./result ] ; then
  rm -r result
)
mkdir result

PATH=$CUTOROOT/bin;$PATH
cd ./result

NUM=0

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n NoServiceTask > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n NoServiceTask -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n NoStartEvent > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n NoStartEvent -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n NoEndEvent > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n NoEndEvent -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n MultiStartEvent > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n MultiStartEvent -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n MultiEndEvent > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n MultiEndEvent -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=12
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName1 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName1 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName2 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName2 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName3 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName3 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName4 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName4 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName5 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName5 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName6 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName6 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName7 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName7 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName8 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName8 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName9 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName9 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName10 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName10 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName11 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n ForbiddenJobName11 -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=35
master -c $CUTOROOT/bin/master.ini -n StartWithoutStartEvent -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n StartWithoutStartEvent -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n EndWithoutEndEvent -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n EndWithoutEndEvent -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n Isolation -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n Isolation -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n DuplicateID -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n DuplicateID -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n BranchWithoutGateway -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n BranchWithoutGateway -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n MergeWithoutGateway -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n MergeWithoutGateway -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n EndBeforeMerge -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n EndBeforeMerge -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n NestedBranch -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n NestedBranch -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n NotMerge -s > pf%NUM%.log
echo $? >> pf%NUM%.log

NUM=`expr $NUM + 1`
master -c $CUTOROOT/bin/master.ini -n NotMerge -s > pf%NUM%.log
echo $? >> pf%NUM%.log

exit 0
