Vecno-Core testnet.

How to build:

Ubuntu 20.04 blank (64-bit) <br>
sudo apt-get install build-essential <br>
sudo apt update && sudo apt upgrade <br>
wget https://go.dev/dl/go1.21.4.linux-amd64.tar.gz -O go.tar.gz <br>
sudo tar -xzvf go.tar.gz -C /usr/local <br>
echo export PATH=$HOME/go/bin:/usr/local/go/bin:$PATH >> ~/.profile <br>
source ~/.profile <br>
apt install git <br>
git clone https://github.com/Vecno-Foundation/vecno-node.git <br>
