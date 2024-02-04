echo "test1.sh" >/home/kehao/result1.txt
hostname -s >>/home/kehao/result1.txt
hostname -i >>/home/kehao/result1.txt
id >>/home/kehao/result1.txt
ls -l /home/kehao/test1.txt >>/home/kehao/result1.txt
cat /home/kehao/test1.txt >>/home/kehao/result1.txt
chown kehao: /home/kehao/result1.txt
echo "I am $(hostname -i) :)"
