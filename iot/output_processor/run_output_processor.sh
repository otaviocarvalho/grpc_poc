#!/bin/bash
i=0
ls -d */ | while read d
do
    echo $d
    go run output_processor.go -i ./$d -o ./$d/output_$i.csv
    ((i++))
done
