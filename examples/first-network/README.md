1. Start etcdraft network
```
./byfn.sh generate -o etcdraft
```

2. Create solo channel
```
docker exec -it cli bash

./scripts/solochannel.sh
```

3. Create kafka channel
```
docker-compose -f docker-compose-kafka.yaml up -d

docker exec -it cli bash

./scripts/kafkachannel.sh
```
