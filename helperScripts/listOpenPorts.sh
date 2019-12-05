echo "If port 80 is not open change deploy gateway local to redirect port 89 to port 80 on docker."
echo "Testing port 80\n"
sudo lsof -i tcp:80
echo "Testing port 89"
sudo lsof -i tcp:89
echo "\nComplete!"
