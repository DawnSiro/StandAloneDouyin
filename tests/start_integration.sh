lsof -i:30000 | awk '{print $2}' | grep -v PID | xargs kill
go mod tidy && go run test_init_data.go
cd .. && go mod tidy && go build .
echo '编译完成'
nohup ./douyin -c tests/integration/integration.yml&
while ! nc -z localhost 30000; do
    sleep 1
done
echo '启动完成'
cd tests && ginkgo integration
