echo "Relay transfer packet"

printf "12345678\n" | evrtcli \
  --home ~/evrt-testnets/evrt1/n0/evrtcli/ \
  tx ibc channel pull ibctransfer chan3 \
  --node1 http://localhost:16657 \
  --node2 http://localhost:26657 \
  --chain-id2 evrt0 \
  --from n1
