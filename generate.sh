#!/bin/sh

SPEC=24501-f70

if [ ! -f nasConvert/GPRSTimer2.go ] ; then
    echo "BAD directory."
    exit 1
fi

if [ ! -f spec.csv ] ; then
    if [ ! -f ${SPEC}.zip ] ; then
        wget https://www.3gpp.org/ftp/Specs/archive/24_series/24.501/${SPEC}.zip
    fi
    if [ ! -f ${SPEC}.zip ] ; then
        echo "Download failed."
        exit 1
    fi
    python3 extract.py
fi
