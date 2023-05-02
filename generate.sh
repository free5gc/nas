#!/bin/sh

SPEC=24501-f70

cd `dirname $0`

if [ ! -f spec.csv ] ; then
    if [ ! -f ${SPEC}.zip ] ; then
        wget https://www.3gpp.org/ftp/Specs/archive/24_series/24.501/${SPEC}.zip
    fi
    if [ ! -f ${SPEC}.zip ] ; then
        echo "Download failed."
        exit 1
    fi
    python3 internal/tools/extract.py
fi

rm -rf testdata/GmmMessage testdata/GsmMessage
mkdir -p testdata/GmmMessage testdata/GsmMessage
rm -f testdata/fuzz/FuzzGmmMessageDecode/msg* testdata/fuzz/FuzzGsmMessageDecode/msg*
ls nasMessage/*go | grep -v "_test" | grep -v "NAS_EPD" | grep -v "NAS_CommInfoIE" |  xargs rm -f
go run internal/tools/generator_sub.go
go run internal/tools/generator/cmd/cmd.go
