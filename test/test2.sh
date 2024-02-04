echo "test2.sh" >/home/kehao/result2.txt
hostname -s >>/home/kehao/result2.txt
hostname -i >>/home/kehao/result2.txt
id >>/home/kehao/result2.txt
ls -l /home/kehao/test2.txt >>/home/kehao/result2.txt
cat /home/kehao/test2.txt >>/home/kehao/result2.txt
chown kehao: /home/kehao/result2.txt
echo "I am $(hostname -i) :)"
