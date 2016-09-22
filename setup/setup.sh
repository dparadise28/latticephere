orig_path = $(pwd)

echo $orig_path #$(pwd)
cd $HOME/bin
echo $(pwd)
#wget https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz

#tar xzvf go1.6.linux-amd64.tar.gz

gr = "$HOME/bin/go"
gp = "$GOROOT/bin:/usr/local/go/bin"

if ["$GOROOT" != *$gr*]; then
	echo "setting GOROOT"
	export GOROOT = $gr
	echo "GOROOT set to: $GORROT"
else
	echo "GOROOT already set to: $GOROOT"
fi

if ["$PATH" != *$gp*]; then
	echo "adding go to PATH"
	export PATH = "$PATH:$gp"
	echo "PATH set to: $PATH"
else
	echo "GOPATH already in path as: $PATH"
fi

cd $orig_path
