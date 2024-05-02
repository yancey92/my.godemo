#!/bin/bash

while read line
 do
    echo "${line}"
done < $(curl  10.121.112.49:30009/api/v1/registering/registered\?watch=1)