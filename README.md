Vecno-Core testnet.

How to build:

Ubuntu 20.04 blank (64-bit)
sudo apt-get install build-essential
sudo apt update && sudo apt upgrade
wget https://go.dev/dl/go1.21.4.linux-amd64.tar.gz -O go.tar.gz
sudo tar -xzvf go.tar.gz -C /usr/local
echo export PATH=$HOME/go/bin:/usr/local/go/bin:$PATH >> ~/.profile
source ~/.profile
apt install git
git clone https://github.com/Vecno-Foundation/vecno-node.git
