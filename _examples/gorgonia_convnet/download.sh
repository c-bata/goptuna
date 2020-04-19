#!/bin/bash

DST=$(cd $(dirname $0); pwd)/data
cd "$DST"

if [ ! -f $DST/train-images-idx3-ubyte ]; then
    curl -O http://yann.lecun.com/exdb/mnist/train-images-idx3-ubyte.gz
    gunzip t*-ubyte.goz
    mv train-images-idx3-ubyte train-images.idx3-ubyte
fi

if [ ! -f $DST/train-labels-idx1-ubyte ]; then
    curl -O http://yann.lecun.com/exdb/mnist/train-labels-idx1-ubyte.gz
    gunzip t*-ubyte.gz
    mv train-labels-idx1-ubyte train-labels.idx1-ubyte
fi

if [ ! -f $DST/t10k-images-idx3-ubyte ]; then
    curl -O http://yann.lecun.com/exdb/mnist/t10k-images-idx3-ubyte.gz
    gunzip t*-ubyte.gz
    mv t10k-images-idx3-ubyte t10k-images.idx3-ubyte
fi


if [ ! -f $DST/t10k-labels-idx1-ubyte ]; then
    curl -O http://yann.lecun.com/exdb/mnist/t10k-labels-idx1-ubyte.gz
    gunzip t*-ubyte.gz
    mv t10k-labels-idx1-ubyte t10k-labels.idx1-ubyte
fi
