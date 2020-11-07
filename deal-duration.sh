#!/bin/bash
#read file
rm -rf duration-$1
mkdir duration-$1
cp $1 duration-$1/
cd duration-$1
cut -d '"' -f 6,10,18,22,38,42 $1 > tmp1.log
sed -e '/cannot/d' -i tmp1.log
sed -e 's/"/ /g' -i tmp1.log
awk '{print $1" "$2}' tmp1.log > tmp2.log
sed -e '/default/d' -i tmp2.log
sed -e '/kube-system/d' -i tmp2.log
sed -e '/lens-metrics/d' -i tmp2.log
sed -e '/ingress-nginx/d' -i tmp2.log
awk '{print $1}' tmp2.log | sort | uniq > tmp3.log
for line in `cat tmp3.log`
do
    grep $line tmp1.log | awk '{print $1" "$3}' > dg-$line.log
    head -n 1 dg-$line.log >> dg-.log
    tail -n 1 dg-$line.log >> dg-.log
    rm dg-$line.log
done
rm tmp*.log
rm $1

go build ../../go/duration_submit.go
./duration_submit -s "$2"

for line in `cat result-onehour.log`
do
    grep $line dg-.log >> onehour-task.log
done
#g++ ../dg-duration.cpp -o ./duration
#./duration dg-.log $2 > result-duration.log
#g++ ../dg-submit.cpp -o ./submit
#./submit dg-.log $2 > result-submit.log
