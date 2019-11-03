#!/bin/zsh

echo "Killing existing evrtd instances..."

killall evrtd

echo "Generating configurations..."

cd ~ && mkdir -p evrt-testnets && cd evrt-testnets
echo -e "\n" | evrtd testnet -o evrt0 --v 1 --chain-id evrt0 --node-dir-prefix n --node-cli-home evrtcli --node-daemon-home evrtd --staking-denom uatom
echo -e "\n" | evrtd testnet -o evrt1 --v 1 --chain-id evrt1 --node-dir-prefix n --node-cli-home evrtcli --node-daemon-home evrtd --staking-denom uevrt

if [ "$(uname)" = "Linux" ]; then
  sed -i 's/"leveldb"/"goleveldb"/g' evrt0/n0/evrtd/config/config.toml
  sed -i 's/"leveldb"/"goleveldb"/g' evrt1/n0/evrtd/config/config.toml
  sed -i 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:16656"#g' evrt1/n0/evrtd/config/config.toml
  sed -i 's#"tcp://0.0.0.0:26657"#"tcp://0.0.0.0:16657"#g' evrt1/n0/evrtd/config/config.toml
  sed -i 's#"localhost:6060"#"localhost:6061"#g' evrt1/n0/evrtd/config/config.toml
  sed -i 's#"tcp://127.0.0.1:26658"#"tcp://127.0.0.1:16658"#g' evrt1/n0/evrtd/config/config.toml
else
  sed -i '' 's/"leveldb"/"goleveldb"/g' evrt0/n0/evrtd/config/config.toml
  sed -i '' 's/"leveldb"/"goleveldb"/g' evrt1/n0/evrtd/config/config.toml
  sed -i '' 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:16656"#g' evrt1/n0/evrtd/config/config.toml
  sed -i '' 's#"tcp://0.0.0.0:26657"#"tcp://0.0.0.0:16657"#g' evrt1/n0/evrtd/config/config.toml
  sed -i '' 's#"localhost:6060"#"localhost:6061"#g' evrt1/n0/evrtd/config/config.toml
  sed -i '' 's#"tcp://127.0.0.1:26658"#"tcp://127.0.0.1:16658"#g' evrt1/n0/evrtd/config/config.toml
fi;

evrtcli config --home evrt0/n0/evrtcli/ chain-id evrt0
evrtcli config --home evrt1/n0/evrtcli/ chain-id evrt1
evrtcli config --home evrt0/n0/evrtcli/ output json
evrtcli config --home evrt1/n0/evrtcli/ output json
evrtcli config --home evrt0/n0/evrtcli/ node http://localhost:26657
evrtcli config --home evrt1/n0/evrtcli/ node http://localhost:16657

echo "Importing keys..."

SEED0=$(jq -r '.secret' evrt0/n0/evrtcli/key_seed.json)
SEED1=$(jq -r '.secret' evrt1/n0/evrtcli/key_seed.json)
echo -e "12345678\n" | evrtcli --home evrt1/n0/evrtcli keys delete n0

echo "Seed 0: ${SEED0}"
echo "Seed 1: ${SEED1}"

echo "Enter seed 0:"
evrtcli --home evrt1/n0/evrtcli keys add n0 --recover

echo "Enter seed 1:"
evrtcli --home evrt0/n0/evrtcli keys add n1 --recover

echo "Enter seed 1:"
evrtcli --home evrt1/n0/evrtcli keys add n1 --recover

#echo -e "12345678\n12345678\n$SEED1\n" | evrtcli --home evrt0/n0/evrtcli keys add n1 --recover
#echo -e "12345678\n12345678\n$SEED0\n" | evrtcli --home evrt1/n0/evrtcli keys add n0 --recover
#echo -e "12345678\n12345678\n$SEED1\n" | evrtcli --home evrt1/n0/evrtcli keys add n1 --recover

echo "Keys should match:"

evrtcli --home evrt0/n0/evrtcli keys list | jq '.[].address'
evrtcli --home evrt1/n0/evrtcli keys list | jq '.[].address'

echo "Starting Everettd instances..."

nohup evrtd --home evrt0/n0/evrtd --log_level="*:debug" start > evrt0.log &
nohup evrtd --home evrt1/n0/evrtd --log_level="*:debug" start > evrt1.log &

sleep 10

echo "Starting rest server..."

nohup evrtcli --home ~/evrt-testnets/evrt0/n0/evrtcli/ rest-server --trust-node --laddr tcp://localhost:1317 &
nohup evrtcli --home ~/evrt-testnets/evrt1/n0/evrtcli/ rest-server --trust-node --laddr tcp://localhost:2317 &

echo "Creating clients..."

# client for chain ibc1 on chain ibc0
echo -e "12345678\n" | evrtcli --home evrt0/n0/evrtcli \
  tx ibc client create c0 \
  "$(evrtcli --home evrt1/n0/evrtcli q ibc client consensus-state)" \
  --from n0 -y -o text

# client for chain ibc0 on chain ibc1
echo -e "12345678\n" | evrtcli --home evrt1/n0/evrtcli \
  tx ibc client create c1 \
  "$(evrtcli --home evrt0/n0/evrtcli q ibc client consensus-state)" \
  --from n1 -y -o text

sleep 3

echo "Establishing a connection..."

evrtcli \
  --home evrt0/n0/evrtcli \
  tx ibc connection handshake \
  conn0 c0 "$(evrtcli --home evrt1/n0/evrtcli q ibc client path)" \
  conn1 c1 "$(evrtcli --home evrt0/n0/evrtcli q ibc client path)" \
  --chain-id2 evrt1 \
  --from1 n0 --from2 n1 \
  --node1 tcp://localhost:26657 \
  --node2 tcp://localhost:16657

echo "Establishing a channel..."

evrtcli \
  --home evrt0/n0/evrtcli \
  tx ibc channel handshake \
  ibcmocksend chan0 conn0 \
  ibcmockrecv chan1 conn1 \
  --node1 tcp://localhost:26657 \
  --node2 tcp://localhost:16657 \
  --chain-id2 evrt1 \
  --from1 n0 --from2 n1

evrtcli \
  --home evrt0/n0/evrtcli \
  tx ibc channel handshake \
  ibctransfer chan2 conn0 \
  ibctransfer chan3 conn1 \
  --node1 tcp://localhost:26657 \
  --node2 tcp://localhost:16657 \
  --chain-id2 evrt1 \
  --from1 n0 --from2 n1

evrtcli \
  --home evrt0/n0/evrtcli \
  tx ibc channel handshake \
  interchainaccount chan4 conn0 \
  interchainaccount chan5 conn1 \
  --node1 tcp://localhost:26657 \
  --node2 tcp://localhost:16657 \
  --chain-id2 evrt1 \
  --from1 n0 --from2 n1

#echo "Sending seq packets from evrt0..."

#evrtcli \
#  --home evrt0/n0/evrtcli \
#  tx ibcmocksend sequence \
#  chan0 "$(evrtcli --home evrt0/n0/evrtcli q ibcmocksend next chan0)" --from n0 -o text -y

#echo "Recieving seq packets on evrt1..."

#evrtcli \
#  --home evrt1/n0/evrtcli \
#  tx ibc channel pull ibcmockrecv chan1 \
#  --node1 tcp://localhost:16657 \
#  --node2 tcp://localhost:26657 \
#  --chain-id2 evrt0 \
#  --from n1
