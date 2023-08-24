go mod tidy && go run test_init_data.go
cd .. && go mod tidy && go build .
echo '编译完成'
nohup ./douyin -c tests/test_config.yml&
while ! nc -z localhost 30000; do
    sleep 1
done
echo '启动完成'
