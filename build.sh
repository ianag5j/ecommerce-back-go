cd lambda/
for d in */; do
  cd $d

  outupDir=""
  if [ -d "cmd/" ]; then
    cd "cmd/"
    outupDir="../"
  fi
  outupDir="$outupDir../../bin/$d"
  go mod tidy
  env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $outupDir main.go
  echo "builded $d"
  cd ..
  if [ -d "cmd/" ]; then
    cd ..
  fi
done
