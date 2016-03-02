#!/usr/bin/env bash

set -x

#cd src/
#rm -rf github.com
#go get -u github.com/bitly/go-simplejson
#go get -u github.com/garyburd/redigo/redis
#go get -u github.com/czhou/INSYNC-Futures-Lib

cd src/main/
go build main.go

while getopts "h:s:" opt; do  
  case $opt in  
    h)  
	    HOST=$OPTARG
	    ;;  
	s)  
	    STG=$OPTARG
	    ;;
 
    \?)  
      echo "Invalid option: -$OPTARG"   
      ;;  
  esac  
done 

./main -h $HOST -s $STG
