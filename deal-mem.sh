#!/bin/bash
#read file
rm -rf mem-$1
mkdir mem-$1
cp $1 mem-$1/
cd mem-$1
cut -d '"' -f 10,18,42 $1 > tmp1.log
sed -e '/cannot/d' -i tmp1.log
sed -e 's/"/ /g' -i tmp1.log
sed -e '/default/d' -i tmp1.log
sed -e '/kube-system/d' -i tmp1.log
sed -e '/lens-metrics/d' -i tmp1.log
sed -e '/ingress-nginx/d' -i tmp1.log
grep "i[[:space:]]*$" tmp1.log > tmp2.log
#awk '{print $2}' tmp2.log | sort | uniq > tmp3.log
#g++ ../cal-sum.cpp -o ./cal-sum
awk '{print $2" "$3}' tmp2.log > tmp3.log
go build ../../go/cal-mem.go
./cal-mem -n tmp3.log
rm tmp*.log
rm $1
