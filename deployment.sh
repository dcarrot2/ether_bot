echo 'Removing source...'
cd ..
rm -rf ether_bot/
echo 'Cloning latest source'
git clone https://github.com/dcarrot2/ether_bot.git
cd ether_bot
echo 'Building new image...'
docker-compose up -d
