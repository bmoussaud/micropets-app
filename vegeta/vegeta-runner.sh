#!/bin/sh
COUNTER=0
echo "PAUSE_WAIT ${PAUSE_WAIT}"
echo "MAX_LOOP ${MAX_LOOP}"
echo "VEGETA_DURATION ${VEGETA_DURATION}"
echo "VEGETA_RATE ${VEGETA_RATE}"
vegeta -version
echo "TARGETS"
cat /vegeta-config/targets.txt
echo "/TARGETS"

while [ : ]
do
    echo "Start Vegeta ${COUNTER}"
    vegeta attack -targets=/vegeta-config/targets.txt -rate=${VEGETA_RATE} -duration=${VEGETA_DURATION} | vegeta plot > /vegeta-data/plot_${COUNTER}.html
    echo "Stop Vegeta"
    echo "Pause ${PAUSE_WAIT}"
    sleep ${PAUSE_WAIT}
    COUNTER=$(( COUNTER + 1 ))
    if test ${COUNTER} -eq ${MAX_LOOP}; 
    then 
        echo "Max has been reached ${COUNTER}/${MAX_LOOP}"
        exit 0
    fi
done