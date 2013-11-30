#!/bin/sh

./build.sh || exit 1
ant -f android/build.xml clean debug install
