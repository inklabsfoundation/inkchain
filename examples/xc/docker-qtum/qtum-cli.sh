#! /bin/bash

echo '      ############################'
echo '     #     1: qtumd_node1       #'
echo '    #     2: qtumd_node2       #'
echo '   #     3: qtumd_node3       #'
echo '  ############################'
echo ' # (default:1) Enter [1~3]:'
read aNum
echo "result"
case $aNum in
    '')
    docker run -i --network container:qtumd_node1 \
    -v ${PWD}/node1_qtumd.conf:/home/qtum/qtum.conf:ro \
    -v ${PWD}/node1_data:/data \
    cryptominder/qtum:latest \
    qtum-cli $@
    ;;
    1) docker run -i --network container:qtumd_node1 \
    -v ${PWD}/node1_qtumd.conf:/home/qtum/qtum.conf:ro \
    -v ${PWD}/node1_data:/data \
    cryptominder/qtum:latest \
    qtum-cli $@
    ;;
    2) docker run -i --network container:qtumd_node2 \
    -v ${PWD}/node2_qtumd.conf:/home/qtum/qtum.conf:ro \
    -v ${PWD}/node2_data:/data \
    cryptominder/qtum:latest \
    qtum-cli $@
    ;;
    3) docker run -i --network container:qtumd_node3 \
    -v ${PWD}/node3_qtumd.conf:/home/qtum/qtum.conf:ro \
    -v ${PWD}/node3_data:/data \
    cryptominder/qtum:latest \
    qtum-cli $@
    ;;
    *)  echo '输入错误!!!'
    ;;
esac
