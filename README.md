# Cocos Tool Kit

Simple tools for cocos creator

## Build

```bash
go get github.com/zhiruili/cckit
cd cckit
go build
# or use make command to build other platform version:
# > make win
```

## lsnode

Subcommand `lsnode` can list all nodes from root by parse prefab file.

Example usage:

```bash
cckit lsnode --sep="." ~/projects/mygame/assets/resources/example.prefab
```

Output:

```log
testRoot.pannel1.img11
testRoot.pannel1.img12
testRoot.pannel2.img21
testRoot.pannel2.panel3.img31
testRoot.pannel2.panel3.img32
```

```bash
cckit lsnode --cuthead=1 --maxdepth=3 --prefix='root' --suffix='.node' --nprefix='["' --nsuffix='"]' ~/projects/mygame/assets/resources/example.prefab
```

Output:

```log
root["pannel1"]["img11"].node
root["pannel1"]["img12"].node
root["pannel2"]["img21"].node
root["pannel2"]["panel3"].node
```

## findref

Subcommand `findref` can find references from all prefabs.

Example usage:

```bash
cckit findref -s ~/projects/mygame/assets ~/projects/mygame/assets/resources/img.png
```

Output:

```log
example.prefab: img11
example.prefab: img12
```
