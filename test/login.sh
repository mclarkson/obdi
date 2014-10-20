proto="https"
opts="-k -s" # don't check ssl cert, silent
ipport="127.0.0.1:443"
guid=`curl $opts -d '{"Login":"admin","Password":"admin"}' $proto://$ipport/api/login | grep -o "[a-z0-9][^\"]*"`
echo $guid

