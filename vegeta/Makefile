# https://www.jesuisundev.com/load-testing-en-5-minutes-avec-vegeta/

load:
	vegeta attack -targets=targets.txt -header AppId:c7c734e9 -name=300qps -rate=300 -duration=25s > results.300qps.bin

report: load	
	cat results.300qps.bin | vegeta plot > plot.300qps.html


burst: 	
	vegeta attack -targets=targets.txt -header AppId:c7c734e9 -name=1500qps -rate=1500 -duration=5s  > results.1500qps.bin
	cat results.1500qps.bin | vegeta plot > plot.1500qps.html

long:
	vegeta attack -targets=targets.txt -header AppId:c7c734e9 -name=1500qps -rate=1500 -duration=600s > results.1500qps.bin
	cat results.1500qps.bin | vegeta plot > plot.1500qps.html
	cat results.1500qps.bin | vegeta report -type=text