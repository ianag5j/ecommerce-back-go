cd lambda/
for d in */; do
  cd $d
  env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../bin/"$d" main.go
  echo "builded $d"
  cd ..
done
