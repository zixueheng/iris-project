#！/bin/bash
PROJECT=iris-project
if [ ! -z "$1" ]
then
	PROJECT=$1
fi

if 
    read -p "Restart project ${PROJECT}? Type Y/N: " ifStart &&
    [ $ifStart == "Y" ]
then
    chmod +x ${PROJECT}
    rm -f nohup.out
    ps -ef | grep ${PROJECT} | grep -v grep | awk '{print $2}' | xargs kill -9
    nohup ./${PROJECT} &
    echo "Restart ${PROJECT} successful..."
else
    echo "Bye..."
fi