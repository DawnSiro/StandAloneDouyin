lsof -i:30000 | awk '{print $2}' | grep -v PID | xargs kill -9
cd .. && go mod tidy && go run test_init_data.go -f benchmark/benchmark.sql
cd .. && go mod tidy && go build .
echo '编译完成'
nohup ./douyin -c tests/benchmark/benchmark.yml&
while ! nc -z localhost 30000; do
    sleep 1
done
echo '启动完成'
