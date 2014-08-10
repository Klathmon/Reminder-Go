#!/bin/bash

#Taken from https://github.com/gopns/gopns/blob/master/test-coverage.sh and modified

echo "mode: set" > acc.out
for Dir in $(find ./* -maxdepth 10 -type d );
do
	if ls $Dir/*.go &> /dev/null;
	then
		returnval=`go test -v -parallel=8 -coverprofile=profile.out $Dir`
		echo ${returnval}
		if [[ ${returnval} != *FAIL* ]]
		then
    		if [ -f profile.out ]
    		then
        		cat profile.out | grep -v "mode: set" >> acc.out
    		fi
    	else
    		exit 1
    	fi
    fi
done

$HOME/gopath/bin/goveralls -coverprofile=acc.out -repotoken 04kTnPgwp65LWxqgAsuXnDazEahF5Wrht
