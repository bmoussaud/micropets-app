#!/bin/sh
COUNTER=0
echo "PAUSE_WAIT ${PAUSE_WAIT}"
echo "MAX_LOOP ${MAX_LOOP}"
echo "VEGETA_DURATION ${VEGETA_DURATION}"
echo "VEGETA_RATE ${VEGETA_RATE}"
vegeta -version

while [ : ]
do
    echo "Start Vegeta ${COUNTER}"
    echo 'GET http://front.mytanzu.xyz/pets' | vegeta attack -rate=${VEGETA_RATE} -duration=${VEGETA_DURATION} | vegeta plot > /vegeta-data/plot_${COUNTER}.html
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