echo "Relay transfer packet"

printf "12345678\n" | evrtcli \
  --home ~/evrt-testnets/evrt0/n0/evrtcli/ \
  tx ibc channel pull ibctransfer chan2 \
  --node1 http://localhost:26657 \
  --node2 http://localhost:16657 \
  --chain-id2 evrt1 \
  --from n0
