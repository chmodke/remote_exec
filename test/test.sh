hostname -s >/home/kehao/result.txt
hostname -i >>/home/kehao/result.txt
id >>/home/kehao/result.txt
ls -l /home/kehao/test.txt >>/home/kehao/result.txt
chown kehao: /home/kehao/result.txt
echo "I am $(hostname -i) :)"
