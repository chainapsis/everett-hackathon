echo "Relay register interchain account packet"

printf "12345678\n" | evrtcli \
  --home ~/evrt-testnets/evrt0/n0/evrtcli/ \
  tx ibc channel pull interchainaccount chan4 \
  --node1 http://localhost:26657 \
  --node2 http://localhost:16657 \
  --chain-id2 evrt1 \
  --from n0 \
  --limit 1

echo "Relay transfer asset packet"

printf "12345678\n" | evrtcli \
  --home ~/evrt-testnets/evrt0/n0/evrtcli/ \
  tx ibc channel pull ibctransfer chan2 \
  --node1 http://localhost:26657 \
  --node2 http://localhost:16657 \
  --chain-id2 evrt1 \
  --from n0

echo "Relay run interchain tx packet"

printf "12345678\n" | evrtcli \
  --home ~/evrt-testnets/evrt0/n0/evrtcli/ \
  tx ibc channel pull interchainaccount chan4 \
  --node1 http://localhost:26657 \
  --node2 http://localhost:16657 \
  --chain-id2 evrt1 \
  --from n0 \
  --limit 1
