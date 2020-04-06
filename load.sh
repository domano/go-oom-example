#!/usr/bin/env bash

for i in {1..10000}
do
   echo "Sending Request Number $i"
   curl -XPUT \
   -d "Blogpost number $i and some more text to increase mem usage and some more text to increase mem usage and some more text to increase mem usage and some more text to increase mem usage and some more text to increase mem usage and some more text to increase mem usage and some more text to increase mem usage and some more text to increase mem usage and some more text to increase mem usage and some more text to increase mem usage and some more text to increase mem usage and some more text to increase mem usage and some more text to increase mem usage and some more text to increase mem usage and some more text to increase mem usage and some more text to increase mem usage and some more text to increase mem usage and some more text to increase mem usage and some more text to increase mem usage and some more text to increase mem usage" \
   "localhost:8080/blog/$i"
done

