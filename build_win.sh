VERSION=$1
cd cats
./build_win.sh $VERSION
cd ..

cd dogs
./build_win.sh $VERSION
cd ..

cd pets
./build_win.sh $VERSION
cd ..

cd gui
./build_win.sh $VERSION
cd ..
